package main

import "github.com/dbcdk/go-smaug/smaug"

func setOwner(app App, identity smaug.Identity) App {
	if app.Labels == nil {
		labels := make(map[string]string)
		app.Labels = &labels
	}

	(*app.Labels)["owner"] = identity.Id + "@" + identity.Agency

	return app
}
