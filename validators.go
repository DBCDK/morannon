package main

import (
	"errors"
	"github.com/Jeffail/gabs"
	"github.com/dbcdk/go-smaug/smaug"
	"strings"
	"regexp"
)

var (
	matchMultiSlash = regexp.MustCompile("/+")
)

func validateId(app *gabs.Container, identity smaug.Identity, idPrefix string) error {
	id := app.Path("id").Data().(string)

	if !strings.HasPrefix(id, "/") {
		id = "/" + idPrefix + "/" + id
	}

	id = matchMultiSlash.ReplaceAllString(id, "/")

	if strings.Count(id, "/") < 3 {
		return errors.New("id must have a depth of at last 3, e.g. /foo/bar/my-app (was: " + id + ")")
	}

	idSegments := strings.Split(id, "/")

	switch idSegments[1] {
	case "prod":
	case "staging":
	case "dev":
		break
	default:
		return errors.New("First segment of the id must be one of prod|staging|dev")

	}

	return nil
}

func validateIsApp(app *gabs.Container, identity smaug.Identity, idPrefix string) error {
	isGroupErr := validateIsGroup(app, identity, idPrefix)
	isApp := isGroupErr != nil

	if !isApp {
		return errors.New("The app appears to be a group")
	}

	return nil
}

func validateIsGroup(group *gabs.Container, identity smaug.Identity, idPrefix string) error {
	_, groupsErr := group.Path("groups").Children()
	_, appsErr := group.Path("apps").Children()

	if groupsErr != nil && appsErr != nil {
		return errors.New("The group is missing apps or subgroups")
	}

	return nil
}
