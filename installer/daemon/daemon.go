package daemon

import (
	"bufio"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"time"

	"github.com/moby/moby/client"
	"github.com/pkg/browser"
)

//go:embed docker-compose.template.yml
var dockerComposeTemplateFile string

type ImageRunner struct {
	ctx                 context.Context
	tmpl                *template.Template
	templateFilePath    string
	readyDockerFilepath string
}

type templateInput struct {
	Version string
}

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
		defer close(installationawait)
		delay := 5 * time.Second
		ticker := time.NewTicker(delay)

		defer ticker.Stop()
		log.Printf("[Support-installer] Awaiting installation...")
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				err := Ping(cli, ctx)
				if err == nil {
					return
				}
				delay = min(delay*2, 60*time.Second)

				ticker.Reset(delay)
				log.Printf("Awaiting docker installation waiting=%d", delay)
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

func (d *ImageRunner) Up(version string) {
	log.Printf("[Daemon] Creating docker file with version %s", version)
	readyDockerFile, err := os.Create(d.readyDockerFilepath)
	if err != nil {
		panic(err)
	}
	err = d.tmpl.Execute(readyDockerFile, templateInput{
		Version: version,
	})
	if err != nil {
		panic(err)
	}
	cmd := exec.Command("docker", "compose", "-f", d.readyDockerFilepath, "up")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(stdout)

	if err = cmd.Start(); err != nil {
		if errors.Is(err, context.Canceled) {
			log.Println("Received cancel ctx")
			return
		}
	}

	for scanner.Scan() {
		log.Printf("[docker-compose up] %s", scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		panic(err)
	}

	if err = cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 130 {
				log.Println("Exiting")
				return
			}
		}
		log.Println(err.Error())
	}

	response, err := exec.CommandContext(context.TODO(), "docker", "exec", "ollama", "list").Output()
	if strings.Contains(string(response), "llama.2:latest") {
		return
	}
	if err = exec.CommandContext(context.TODO(), "docker", "exec", "ollama", "pull", "llama3.2").Run(); err != nil {
		panic("Failed to install model llama3.2")
	}

}

func (d *ImageRunner) Down() {
	cmd := exec.Command("docker", "compose", "-f", d.readyDockerFilepath, "down")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	if err = cmd.Start(); err != nil {
		if errors.Is(err, context.Canceled) {
			log.Println("Received cancel ctx")
			return
		}
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		downText := scanner.Text()
		log.Printf("[docker compose down]: %s", downText)
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	if err = cmd.Wait(); err != nil {
		panic(err)
	}
}

func New(ctx context.Context) *ImageRunner {
	tmpl := template.Must(template.New("docker-compose-mock").Parse(dockerComposeTemplateFile))
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	templateFilepath := filepath.Join(wd, "./docker-compose.template.yml")
	readyDockerFilepath := filepath.Join(wd, "../../", "docker-compose.yml")
	return &ImageRunner{ctx, tmpl, templateFilepath, readyDockerFilepath}
}
