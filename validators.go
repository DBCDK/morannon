package main

import (
	"errors"
	"github.com/dbcdk/go-smaug/smaug"
	"strings"
)

func validatePresenceOfHealthChecks(app App, identity smaug.Identity) error {
	if app.HealthChecks == nil || len(app.HealthChecks) == 0 {
		return errors.New("Health checks required")
	}

	return nil
}

func validateJobId(app App, identity smaug.Identity) error {
	if !strings.HasPrefix(app.Id, "/") {
		return errors.New("Job ID must start with a '/'")
	}

	if strings.Count(app.Id, "/") < 3 {
		return errors.New("Job ID must have a depth of at last 3, e.g. /foo/bar/my-app")
	}

	idSegments := strings.Split(app.Id, "/")
	println(idSegments[1])

	switch idSegments[1] {
	case "prod":
	case "staging":
	case "dev":
		break
	default:
		return errors.New("First segment of the job-ID must be one of prod|staging|dev")

	}

	return nil
}


func validateNetwork(app App, identity smaug.Identity) error {
	if app.Constraints != nil {
		netConfigFound := false
		for _, constraint := range app.Constraints {
			if len(constraint) == 3 && constraint[0] == "net" && constraint[1] == "CLUSTER" {
				netConfigFound = true
			}
		}

		if netConfigFound {
			return nil
		}
	}

	return errors.New("Missing network constraint (e.g. [\"net\", \"CLUSTER\", \"prod\"])")
}