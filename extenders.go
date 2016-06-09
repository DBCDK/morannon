package main

import "github.com/dbcdk/go-smaug/smaug"

func ensureNetwork(app App, identity smaug.Identity) App {
	if app.Constraints == nil {
		app.Constraints = make([]Constraint, 0)
	}

	setNetwork := true
	for _, constraint := range app.Constraints {
		if len(constraint) > 0 && constraint[0] == "net_id" {
			setNetwork = false
		}
	}

	if setNetwork {
		app.Constraints = append(app.Constraints, []string{"net_id", "CLUSTER", "prod"})
	}

	return app
}

func setOwner(app App, identity smaug.Identity) App {
	if app.Labels == nil {
		app.Labels = make(map[string]string)
	}

	app.Labels["owner"] = identity.Id + "@" + identity.Agency

	return app
}
