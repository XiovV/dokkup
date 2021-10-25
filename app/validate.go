package app

import "strings"

func (a *Update) ValidateFlags() []string {
	errors := []string{}

	if a.config.Image != "" && a.config.Tag != "" {
		errors = append(errors, "you can only either set the -image flag or the -tag flag")
	}

	if a.config.NodeLocation == "" {
		errors = append(errors, "please provide a node")
	}

	if a.config.Image != "" {
		imageParts := strings.Split(a.config.Image, ":")
		if len(imageParts) != 2 || imageParts[0] == "" || imageParts[1] == "" {
			errors = append(errors, "invalid image format. Example: imageName:latest")
		}
	}

	if a.config.Image == "" && a.config.Tag == "" {
		errors = append(errors, "please provide either an image flag or a tag flag")
	}

	if a.config.Container == "" {
		errors = append(errors, "please provide a container")
	}

	return errors
}
