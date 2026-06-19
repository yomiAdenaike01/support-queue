package daemon

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"runtime"
	"sync"
	"time"

	"github.com/moby/moby/client"
	"github.com/pkg/browser"
	"github.com/yomiAdenaike01/support-ops/installer/images"
)

func OpenBrowserAwaitInstallation(ctx context.Context, cli *client.Client) <-chan struct{} {
	log.Printf("[Support-installer] Failed to find docker installation opening docker in your browser...")
	baseUrl := "https://docs.docker.com/desktop/setup/install"
	extension := ""
	switch runtime.GOOS {
	case "windows":
		extension = "windows-install"
	case "darwin":
		extension = "mac-install"
	case "linux":
		extension = "linux"
	}
	if extension == "" {
		panic("Platform not found")
	}
	downloadUrlPage := fmt.Sprintf("%s/%s", baseUrl, extension)
	_, err := url.Parse(downloadUrlPage)
	if err != nil {
		panic(err)
	}
	err = browser.OpenURL(downloadUrlPage)
	if err != nil {
		panic(err)
	}
	installationawait := make(chan struct{})

	go func() {
		timout, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		defer close(installationawait)
		log.Printf("[Support-installer] Awaiting installation...")
		for {
			select {
			case <-timout.Done():
				// check installation?
				err := Ping(cli, ctx)
				if err != nil {
					continue
				}
				return
			default:
				continue
			}
		}

	}()

	return installationawait
}

func Ping(cli *client.Client, ctx context.Context) error {
	log.Printf("[Support-installer] Pinging docker deamon")
	timeoutCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_, err := cli.Ping(timeoutCtx, client.PingOptions{
		NegotiateAPIVersion: true,
	})
	return err
}

func RunImages(ctx context.Context, imagesList images.ImagesList, cli *client.Client) map[string]string {
	wg := sync.WaitGroup{}
	imagesListLen := len(imagesList)
	wg.Add(imagesListLen)
	errorsList := make([]error, 0, imagesListLen)
	timeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	containersById := make(map[string]string, imagesListLen)
	defer cancel()
	for _, image := range imagesList {
		go func(imageData images.Image) {
			defer wg.Done()
			tim, c := context.WithTimeout(timeout, 5*time.Second)
			defer c()

			containerId := findOrCreateContainer(cli, tim, image, &errorsList)
			if containerId == "" {
				return
			}
			startContainer(tim, cli, image, containerId, &errorsList)
			containersById[image.Name] = containerId
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

func createContainer(cli *client.Client, ctx context.Context, image images.Image, errorList *[]error) *client.ContainerCreateResult {
	container, err := cli.ContainerCreate(ctx, client.ContainerCreateOptions{
		Image: image.GetPath(),
		Name:  image.Shortname,
	})
	if err != nil {
		log.Printf("[Support-installer] Failed to create container name=%s reason=%s", image.Name, err.Error())
		*errorList = append(*errorList, err)
		return nil
	}
	return &container
}

func startContainer(ctx context.Context, cli *client.Client, image images.Image, containerId string, errorList *[]error) {
	log.Printf("[Support-installer] Starting container name=%s id=%s", image.Name, containerId)
	_, err := cli.ContainerStart(ctx, containerId, client.ContainerStartOptions{})
	if err != nil {
		log.Printf("[Support-installer] Failed to start container id=%s name=%s reason=%s", containerId, image.Name, err.Error())
		*errorList = append(*errorList, err)
		return
	}
}

func findOrCreateContainer(cli *client.Client, ctx context.Context, image images.Image, errorList *[]error) string {
	inspectResult, err := cli.ContainerInspect(ctx, image.Shortname, client.ContainerInspectOptions{})
	if err != nil {
		*errorList = append(*errorList, err)
		return ""
	}
	containerId := inspectResult.Container.ID

	if containerId != "" {
		return containerId
	}
	container := createContainer(cli, ctx, image, errorList)

	if container != nil {
		return container.ID
	}
	return ""

}
