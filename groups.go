package main

import (
	"github.com/Jeffail/gabs"
	"github.com/dbcdk/go-smaug/smaug"
	"strings"
)

func processGroup(group *gabs.Container, identity *smaug.Identity, idPrefix string) (*gabs.Container, error) {
	groupId := group.Path("id").Data().(string)

	if !strings.HasPrefix(groupId, "/") {
		groupId = "/" + idPrefix + "/" + groupId
	}

	groupId = matchMultiSlash.ReplaceAllString(groupId, "/")

	for _, validator := range groupValidators {
		err := validator(group, *identity, idPrefix)
		if err != nil {
			return nil, err
		}
	}

	groups, groupsErr := group.Path("groups").Children()
	apps, appsErr := group.Path("apps").Children()

	if groupsErr == nil {
		for _, marathonGroup := range groups {
			_, processErr := processGroup(marathonGroup, identity, groupId)
			if processErr != nil {
				return nil, processErr
			}
		}
	}

	if appsErr == nil {
		for _, marathonApp := range apps {
			_, processErr := processApp(marathonApp, identity, groupId)
			if processErr != nil {
				return nil, processErr
			}
		}
	}

	return group, nil
}
