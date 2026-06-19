package images

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/moby/moby/client"
)

type Image struct {
	Name        string
	Tag         string
	Shortname   string
	Environment map[string]string
}

func (i Image) getPath() string {
	return fmt.Sprintf("%s:%s", i.Name, i.Tag)
}

type ImagesList []Image

func (i ImagesList) AsString() string {
	builder := strings.Builder{}
	for _, image := range i {
		builder.WriteString(fmt.Sprintf("%s,", image.getPath()))
	}
	return builder.String()
}

func (i ImagesList) Pull(ctx context.Context, cli *client.Client) {
	log.Printf("[Support-installer] Pulling images...")
	wg := sync.WaitGroup{}
	wg.Add(len(i))

	errorsList := make([]error, 0, len(i))
	for _, image := range i {
		go func(image Image) {
			timeoutCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
			defer cancel()
			defer wg.Done()
			log.Printf("[Support-installer] Pulling image name=%s", image.Name)
			response, err := cli.ImagePull(timeoutCtx, image.getPath(), client.ImagePullOptions{})
			if err != nil {
				log.Printf("[Support-installer] Failed to pull image name=%s error=%s", image.getPath(), err.Error())
				errorsList = append(errorsList, err)
				return
			}
			defer response.Close()

		}(image)
	}
	wg.Wait()
	if len(errorsList) > 0 {
		log.Printf("[Support-installer] Failed to pull images!")
		return
	}
	log.Printf("[Support-installer] Successfully pulled images!")
}

func (i ImagesList) RunImages(ctx context.Context, cli *client.Client) map[string]string {
	wg := sync.WaitGroup{}
	imagesListLen := len(i)
	wg.Add(imagesListLen)
	errorsList := make([]error, 0, imagesListLen)
	timeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	containersById := make(map[string]string, imagesListLen)
	defer cancel()
	for _, image := range i {
		go func(imageData Image) {
			defer wg.Done()
			tim, c := context.WithTimeout(timeout, 5*time.Second)
			defer c()
			container, err := cli.ContainerCreate(tim, client.ContainerCreateOptions{
				Image: image.getPath(),
				Name:  image.Shortname,
			})
			if err != nil {
				log.Printf("[Support-installer] Failed to create container name=%s reason=%s", image.Name, err.Error())
				errorsList = append(errorsList, err)
				return
			}
			tim, c = context.WithTimeout(timeout, 10*time.Second)
			defer c()
			log.Printf("[Support-installer] Starting container name=%s id=%s", image.Name, container.ID)
			_, err = cli.ContainerStart(tim, container.ID, client.ContainerStartOptions{})
			if err != nil {
				log.Printf("[Support-installer] Failed to start container id=%s name=%s reason=%s", container.ID, image.Name, err.Error())
				errorsList = append(errorsList, err)
			}
			containersById[image.Name] = container.ID
		}(image)
	}
	wg.Wait()
	if len(errorsList) > 0 {
		log.Printf("[Support-installer]: Failed to run containers")
		return nil
	}
	log.Printf("[Support-installer]: Successfully ran containers!")
	return containersById

}

func getAppName(appName string) string {
	imgBasename := "adenaikeyomi/support-queue"
	return fmt.Sprintf("%s-%s", imgBasename, appName)
}

func GetImagesList(version string) ImagesList {
	return ImagesList{
		{
			Name:      "pgvector/pgvector",
			Tag:       "pg17",
			Shortname: "db",
		},
		{
			Name:      "redis",
			Tag:       "latest",
			Shortname: "redis",
		},
		{
			Name:      getAppName("ml-worker"),
			Tag:       version,
			Shortname: "ml-worker",
		},
		{
			Name:      getAppName("frontend"),
			Tag:       version,
			Shortname: "frontend",
		},
		{
			Name:      getAppName("srv"),
			Tag:       version,
			Shortname: "srv",
		},
	}

}
