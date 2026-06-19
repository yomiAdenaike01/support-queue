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

func (i Image) GetPath() string {
	return fmt.Sprintf("%s:%s", i.Name, i.Tag)
}

type ImagesList []Image

func (i ImagesList) AsString() string {
	builder := strings.Builder{}
	for _, image := range i {
		builder.WriteString(fmt.Sprintf("%s,", image.GetPath()))
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
			response, err := cli.ImagePull(timeoutCtx, image.GetPath(), client.ImagePullOptions{})
			if err != nil {
				log.Printf("[Support-installer] Failed to pull image name=%s error=%s", image.GetPath(), err.Error())
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
