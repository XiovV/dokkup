package docker

import (
	"io"
	"strings"

	"github.com/docker/docker/api/types"
	"go.uber.org/zap"
)

func (c *Controller) ImagePull(containerImage string) error {
	exists, err := c.ImageDoesExist(containerImage)
	if err != nil {
		return err
	}

	if exists {
		c.Logger.Debug("image already exists, exiting...")
		return nil
	}

	reader, err := c.cli.ImagePull(c.ctx, containerImage, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer reader.Close()
	io.Copy(io.Discard, reader)

	c.Logger.Debug("image pulled successfully", zap.String("image", containerImage))

	return nil
}

func (c *Controller) ImageDoesExist(containerImage string) (bool, error) {
	if !strings.Contains(containerImage, ":") {
		containerImage += ":latest"
	}

	images, err := c.ImageList(containerImage)
	if err != nil {
		return false, err
	}

	for _, image := range images {
		if len(image.RepoTags) > 0 && image.RepoTags[0] == containerImage {
			return true, nil
		}
	}

	return false, nil
}

func (c *Controller) ImageList(containerImage string) ([]types.ImageSummary, error) {
	images, err := c.cli.ImageList(c.ctx, types.ImageListOptions{All: true})
	if err != nil {
		return nil, err
	}

	return images, nil
}
