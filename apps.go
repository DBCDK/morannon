package main

import (
	"github.com/Jeffail/gabs"
	"github.com/dbcdk/go-smaug/smaug"
)

func processApp(app *gabs.Container, identity *smaug.Identity, idPrefix string) (*gabs.Container, error) {
	for _, validator := range appValidators {
		err := validator(app, *identity, idPrefix)
		if err != nil {
			return nil, err
		}
	}

	for _, extender := range appExtenders {
		app = extender(app, *identity)
	}

	return app, nil
}
