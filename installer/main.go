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

type DockerComposefileTemplateInput struct {
	Version string
}

func main() {
	version := images.GetruntimeVersion()

	imagesList := images.GetImagesList(version)
	log.Printf("[Support-installer] Support queue installer version-%s images=%s", version, imagesList.String())

	ctx := context.Background()
	signalCtx, cancel := signal.NotifyContext(ctx, os.Kill, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	runner := daemon.New(signalCtx)

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

	runner.Up(version)

	<-signalCtx.Done()
	log.Println("Signal detected, tearing down containers...")
	runner.Down()
	log.Println("Successfully closed containers")

}
