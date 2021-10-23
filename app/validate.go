package app

import "strings"

func (a *Update) ValidateFlags() []string {
	errors := []string{}

	if a.image != "" && a.tag != "" {
		errors = append(errors, "you can only either set the -image flag or the -tag flag")
	}

	if a.node == "" {
		errors = append(errors, "please provide a node")
	}

	if a.image != "" {
		imageParts := strings.Split(a.image, ":")
		if len(imageParts) != 2 || imageParts[0] == "" || imageParts[1] == "" {
			errors = append(errors, "invalid image format. Example: imageName:latest")
		}
	}

	if a.image == "" && a.tag == "" {
		errors = append(errors, "please provide either an image flag or a tag flag")
	}

	if a.container == "" {
		errors = append(errors, "please provide a container")
	}

	return errors
}
