package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/moby/moby/client"
	"github.com/yomiAdenaike01/support-ops/installer/daemon"
	"github.com/yomiAdenaike01/support-ops/installer/images"
)

func main() {
	version := images.GetEnvVar("VERSION", "0.0.1-dev")

	imagesList := images.GetImagesList(version)
	log.Printf("[Support-installer] Support queue installer version-%s images=%s", version, imagesList.AsString())

	ctx := context.Background()
	sign, cancel := signal.NotifyContext(ctx, os.Kill, os.Interrupt, syscall.SIGINT)
	defer cancel()

	cli, err := client.New(
		client.FromEnv,
	)
	log.Printf("[Support-installer] Initialised client...")
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	err = daemon.Ping(cli, ctx)

	if err != nil {
		<-daemon.OpenBrowserAwaitInstallation(ctx, cli)
	}

	imagesList.Pull(ctx, cli)
	imagesList.RunImages(ctx, cli)

	<-sign.Done()

}
