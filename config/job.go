package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Job struct {
	Group     string      `yaml:"group"`
	Count     int         `yaml:"count"`
	Container []Container `yaml:"container"`
}

type Container struct {
	Name  string `yaml:"name"`
	Image string `yaml:"image"`
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

  return &job, nil
}
