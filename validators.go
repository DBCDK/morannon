package main

import (
	"errors"
	"github.com/dbcdk/go-smaug/smaug"
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
			if env == "testing" || env == "staging" || env == "prod" {
				return nil
			}
		}
	}

	return errors.New("env-label must be present and equal one of testing|staging|prod")
}
