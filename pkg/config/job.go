package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Job struct {
	Group     string `yaml:"group,omitempty"`
	Node      string `yaml:"node,omitempty"`
	Count     int    `yaml:"count"`
	Name      string
	Container []Container `yaml:"container"`
}

type Container struct {
	Image       string   `yaml:"image"`
	Ports       []Port   `yaml:"ports"`
	Network     string   `yaml:"network"`
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
	content, err := os.ReadFile(filename)
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
	if len(j.Container[0].Network) == 0 {
		j.Container[0].Network = "bridge"
	}

	if j.Container[0].Restart == "" {
		j.Container[0].Restart = "always"
	}
}
