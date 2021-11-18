package app

import "strings"

func (u *Update) ValidateFlags() []string {
	errors := []string{}

	if u.config.Image != "" && u.config.Tag != "" {
		errors = append(errors, "you can only either set the -image flag or the -tag flag")
	}

	if u.config.NodeLocation == "" {
		errors = append(errors, "please provide u node")
	}

	if u.config.Image != "" {
		imageParts := strings.Split(u.config.Image, ":")
		if len(imageParts) != 2 || imageParts[0] == "" || imageParts[1] == "" {
			errors = append(errors, "invalid image format. Example: imageName:latest")
		}
	}

	if u.config.Image == "" && u.config.Tag == "" {
		errors = append(errors, "please provide either an image flag or u tag flag")
	}

	if u.config.Container == "" {
		errors = append(errors, "please provide u container")
	}

	return errors
}
