package main

import (
	"bytes"
	"encoding/json"
	"github.com/dbcdk/go-smaug/smaug"
	"github.com/vulcand/oxy/forward"
	"github.com/vulcand/oxy/roundrobin"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

var (
	app          = kingpin.New("timeattack", "Replays http requests").Version("1.0")
	httpPort     = kingpin.Flag("port", "Http port to listen on").Default("8080").Int()
	marathons    = kingpin.Flag("marathon", "url to marathon (repeatable for multiple instances of marathon)").Required().Strings()
	forwarder, _ = forward.New()
)

func main() {
	kingpin.Parse()

	marathon, _ := roundrobin.New(forwarder)
	for _, marathon_url := range *marathons {
		u, err := url.Parse(marathon_url)
		if err != nil {
			panic(err)
		}

		marathon.UpsertServer(u)
	}

	validators := []validatorFunc{
		validatePresenceOfHealthChecks,
		//validatePresenceOfEnvLabel,
	}

	extenders := []extenderFunc{
		ensureNetwork,
		setOwner,
	}

	redirect := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		token, err := smaug.TokenFromRequest(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		smaugUrl := url.URL{Scheme: "http", Host: "platform-i01.dbc.dk:3001"}
		identity, err := smaug.Authenticate(smaugUrl, token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		if req.Method == "POST" && req.URL.Path == "/v2/apps" {
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				panic(err)
			}

			var marathonApp App
			json.Unmarshal(body, &marathonApp)

			for _, validator := range validators {
				err := validator(marathonApp, *identity)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			}

			for _, extender := range extenders {
				marathonApp = extender(marathonApp, *identity)
			}

			newBody, err := json.MarshalIndent(marathonApp, "", "  ")
			if err != nil {
				panic(err)
			}

			req.Body = ioutil.NopCloser(bytes.NewReader(newBody))
			req.ContentLength = int64(len(newBody))
		}

		marathon.ServeHTTP(w, req)
	})

	s := &http.Server{
		Addr:    ":" + strconv.Itoa(*httpPort),
		Handler: redirect,
	}
	log.Fatal(s.ListenAndServe())
}
