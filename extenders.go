package main

import (
	"github.com/Jeffail/gabs"
	"github.com/dbcdk/go-smaug/smaug"
)

func setOwner(app *gabs.Container, identity smaug.Identity) *gabs.Container {
	app.Set(identity.Id+"@"+identity.Agency, "labels", "owner")

	return app
}
