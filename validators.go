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

func validatePresenceOfEnvLabel(app App, identity smaug.Identity) error {
	if app.Labels != nil {
		if env, ok := app.Labels["env"]; ok {
			if env == "dev" || env == "staging" || env == "prod" {
				return nil
			}
		}
	}

	return errors.New("env-label must be present and equal one of dev|staging|prod")
}

func validateJobId(app App, identity smaug.Identity) error {
	if !strings.HasPrefix(app.Id, "/") {
		return errors.New("Job ID must start with a '/'")
	}

	if strings.Count(app.Id, "/") < 3 {
		return errors.New("Job ID must have a depth of at last 3, e.g. /foo/bar/my-app")
	}

	return nil
}
