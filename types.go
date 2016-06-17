package main

import "github.com/dbcdk/go-smaug/smaug"

type extenderFunc func(App, smaug.Identity) App
type validatorFunc func(App, smaug.Identity) error

type App struct {
	Id              *string            `json:"id,omitempty"`
	User            *string            `json:"user,omitempty"`
	Cmd             *string            `json:"cmd,omitempty"`
	Constraints     *[]Constraint      `json:"constraints,omitempty"`
	Cpus            *float32           `json:"cpus,omitempty"`
	Env             *interface{}       `json:"env,omitempty"`
	Mem             *float32           `json:"mem,omitempty"`
	Instances       *int               `json:"instances,omitempty"`
	Ports           *[]int             `json:"ports,omitempty"`
	Uris            *[]string          `json:"uris,omitempty"`
	Labels          *map[string]string `json:"labels,omitempty"`
	HealthChecks    *[]interface{}     `json:"healthChecks,omitempty"`
	Container       *interface{}       `json:"container,omitempty"`
	UpgradeStrategy *interface{}       `json:"upgradeStrategy,omitempty"`
}

type Constraint []string
