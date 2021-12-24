package rc

import (
	"context"
	"io"
	"log"

	"github.com/docker/docker/api/types"
)

func (c *Client) CreateImageIfNotExists(ctx context.Context, image string) error {
	if c.IsImageExists(ctx, image) {
		log.Println("image already exists")
		return nil
	}

	return c.PullImage(ctx, image)
}

// IsImageExists get the image list and compare the images with the given image
// if the image is in the list return true otherwise return false
func (c *Client) IsImageExists(ctx context.Context, image string) bool {
	imageSummary, err := c.containerPort.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		log.Printf("image list failed: %v\n", err)
		return false
	}

	for idx := range imageSummary {
		repoTags := imageSummary[idx].RepoTags
		for _, tag := range repoTags {
			if tag == image {
				return true
			}
		}
	}

	return false
}

func (c *Client) PullImage(ctx context.Context, image string) error {
	log.Println("pulling image")
	reader, err := c.containerPort.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		log.Printf("failed while pulling the image, err: %v\n", err)
		return err
	}
	defer reader.Close()

	readerBytes, err := io.ReadAll(reader)
	if err != nil {
		log.Printf("failed when reading the bytes, err: %v\n", err)
		return err
	}
	log.Printf("image pull response: %v\n", string(readerBytes))

	return err
}
