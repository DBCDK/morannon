package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/dbcdk/go-smaug/smaug"
	"github.com/julienschmidt/httprouter"
	"github.com/vulcand/oxy/forward"
	"github.com/vulcand/oxy/roundrobin"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	"net/url"
	"strconv"
)

var (
	app              = kingpin.New("morannon", "validates Marathon requests before forwarding them").Version("1.0")
	httpPort         = kingpin.Flag("port", "Http port to listen on").Default("8080").Int()
	marathons        = kingpin.Flag("marathon", "url to Marathon (repeatable for multiple instances of marathon)").Required().Strings()
	marathonUsername = kingpin.Flag("marathon-username", "username for marathon").String()
	marathonPassword = kingpin.Flag("marathon-password", "password for marathon").String()
	smaug_location   = kingpin.Flag("smaug", "url to Smaug").Required().String()
	sslCertFile      = kingpin.Flag("cert", "location of ssl certificate file").String()
	sslKeyFile       = kingpin.Flag("cert-key", "location of ssl certificate key file").String()
	forwarder, _     = forward.New()
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	kingpin.Parse()
}

func main() {
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

	redirect := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		token, err := smaug.TokenFromRequest(req)
		if err != nil {
			log.WithFields(log.Fields{
				"app":   "morannon",
				"event": "authentication_failed",
				"error": err.Error(),
			}).Info("request rejected")

			if req.URL.Path == "/" {
				http.Redirect(w, req, "/login", 302)
			} else {
				http.Error(w, err.Error(), http.StatusUnauthorized)
			}
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

			if req.URL.Path == "/" {
				http.Redirect(w, req, "/login", 302)
			} else {
				http.Error(w, err.Error(), http.StatusForbidden)
			}
			return
		}

		log.WithFields(log.Fields{
			"app":   "morannon",
			"event": "forward_request",
			"token": *token,
			"user":  identity.String(),
		}).Info("forwarding request")

		if len(*marathonUsername) > 0 && len(*marathonPassword) > 0 {
			req.SetBasicAuth(*marathonUsername, *marathonPassword)
		}

		marathon.ServeHTTP(w, req)
	})

	router := httprouter.Router{RedirectTrailingSlash: false, RedirectFixedPath: false, NotFound: redirect}
	router.HandlerFunc("GET", "/login", showLogin)
	router.HandlerFunc("POST", "/login", performLogin)

	enableSsl := len(*sslCertFile) > 0 && len(*sslKeyFile) > 0
	if enableSsl {
		log.Fatal(http.ListenAndServeTLS(":"+strconv.Itoa(*httpPort), *sslCertFile, *sslKeyFile, &router))
	} else {
		log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*httpPort), &router))
	}
}
