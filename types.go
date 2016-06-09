package main

import "github.com/dbcdk/go-smaug/smaug"

type extenderFunc func(App, smaug.Identity) App
type validatorFunc func(App, smaug.Identity) error

type App struct {
	Id              string            `json:"id"`
	User            string            `json:"user"`
	Cmd             string            `json:"cmd"`
	Constraints     []Constraint      `json:"constraints"`
	Cpus            float32           `json:"cpus"`
	Mem             float32           `json:"mem"`
	Instances       int               `json:"instances"`
	Ports           []int             `json:"ports"`
	Uris            []string          `json:"uris"`
	Labels          map[string]string `json:"labels"`
	HealthChecks    []interface{}     `json:"healthChecks"`
	UpgradeStrategy interface{}       `json:"upgradeStrategy"`
}

type Constraint []string
