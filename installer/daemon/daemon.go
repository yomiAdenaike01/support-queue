package daemon

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"runtime"
	"time"

	"github.com/moby/moby/client"
	"github.com/pkg/browser"
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
