package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func main() {
	script := flag.String("script", "setup.sh", "Configuration script file")
	packages := flag.String("packages", "packages.txt", "Package list file")
	repositories := flag.String("repositories", "repositories.txt", "Repostory list file")
	output := flag.String("output", "alpine.qcow2", "Output file")

	flag.Parse()

	if err := os.Chmod(*script, 0777); err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatal(err)
	}

	if _, err := cli.ImagePull(ctx, "docker.io/library/alpine:edge", types.ImagePullOptions{}); err != nil {
		log.Fatal(err)
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{Image: "alpine:edge", Cmd: []string{"tail", "-f", "/dev/null"}}, &container.HostConfig{Privileged: true}, nil, "")
	if err != nil {
		log.Fatal(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		log.Fatal(err)
	}

	log.Println(*script, *packages, *repositories, *output, resp.ID)
}
