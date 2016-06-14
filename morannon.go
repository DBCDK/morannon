package main

import (
	"bytes"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/dbcdk/go-smaug/smaug"
	"github.com/vulcand/oxy/forward"
	"github.com/vulcand/oxy/roundrobin"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

var (
	app            = kingpin.New("timeattack", "Replays http requests").Version("1.0")
	httpPort       = kingpin.Flag("port", "Http port to listen on").Default("8080").Int()
	marathons      = kingpin.Flag("marathon", "url to Marathon (repeatable for multiple instances of marathon)").Required().Strings()
	smaug_location = kingpin.Flag("smaug", "url to Smaug").Required().String()
	forwarder, _   = forward.New()
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	kingpin.Parse()
	smaugUrl, err := url.Parse(*smaug_location)
	if err != nil {
		log.Fatal(err)
	}

	log.WithFields(log.Fields{"app": "morannon", "event": "started"}).Info("morannon started")

	marathon, _ := roundrobin.New(forwarder)
	for _, marathon_url := range *marathons {
		u, err := url.Parse(marathon_url)
		if err != nil {
			log.Fatal(err)
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
			log.WithFields(log.Fields{
				"app":   "morannon",
				"event": "authentication_failed",
				"error": err.Error(),
			}).Info("request rejected")

			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		identity, err := smaug.Authenticate(*smaugUrl, token)
		if err != nil {
			log.WithFields(log.Fields{
				"app":   "morannon",
				"event": "authentication_failed",
				"token": *token,
				"error": err.Error(),
			}).Info("request rejected")

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

		log.WithFields(log.Fields{"app": "morannon", "event": "started"}).Info("forwarding request")

		marathon.ServeHTTP(w, req)
	})

	s := &http.Server{
		Addr:    ":" + strconv.Itoa(*httpPort),
		Handler: redirect,
	}
	log.Fatal(s.ListenAndServe())
}
