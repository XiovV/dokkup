package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Job struct {
	Group     string      `yaml:"group,omitempty"`
	Node      string      `yaml:"node,omitempty"`
	Count     int         `yaml:"count"`
	Container []Container `yaml:"container"`
}

type Container struct {
	Name        string   `yaml:"name"`
	Image       string   `yaml:"image"`
	Ports       []Port   `yaml:"ports"`
	Networks    []string `yaml:"networks"`
	Volumes     []string `yaml:"volumes"`
	Environment []string `yaml:"environment"`
	Restart     string   `yaml:"restart"`
	Labels      []string `yaml:"labels"`
}

type Port struct {
	In  string `yaml:"in"`
	Out string `yaml:"out"`
}

func ReadJob(filename string) (*Job, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	job := Job{}
	err = yaml.Unmarshal(content, &job)
	if err != nil {
		return nil, err
	}

	job.setDefaults()

	return &job, nil
}

func (j *Job) setDefaults() {
	if len(j.Container[0].Networks) == 0 {
		j.Container[0].Networks = []string{"bridge"}
	}

	if j.Container[0].Restart == "" {
		j.Container[0].Restart = "always"
	}
}
