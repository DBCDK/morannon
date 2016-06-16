package main

import "github.com/dbcdk/go-smaug/smaug"

func setOwner(app App, identity smaug.Identity) App {
	if *app.Labels == nil {
		*app.Labels = make(map[string]string)
	}

	(*app.Labels)["owner"] = identity.Id + "@" + identity.Agency

	return app
}
