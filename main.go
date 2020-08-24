package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

const (
	SETUP_SCRIPT_FILE_NAME    = "setup.sh"
	PACKAGE_LIST_FILE_NAME    = "packages.txt"
	REPOSITORY_LIST_FILE_NAME = "repositories.txt"
	OUTPUT_IMAGE_FILE_NAME    = "alpine.qcow2"

	DOCKER_IMAGE_URL = "docker.io/library/alpine:edge"
	DOCKER_IMAGE     = "alpine:edge"

	WORKDIR = "/tmp"
)

func main() {
	script := flag.String("script", SETUP_SCRIPT_FILE_NAME, "Setup script file")
	packages := flag.String("packages", PACKAGE_LIST_FILE_NAME, "Package list file")
	repositories := flag.String("repositories", REPOSITORY_LIST_FILE_NAME, "Repository list file")
	output := flag.String("output", OUTPUT_IMAGE_FILE_NAME, "Output image file")

	flag.Parse()

	if err := os.Chmod(*script, 0777); err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatal(err)
	}

	if _, err := cli.ImagePull(ctx, DOCKER_IMAGE_URL, types.ImagePullOptions{}); err != nil {
		log.Fatal(err)
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{Image: DOCKER_IMAGE, Cmd: []string{"tail", "-f", "/dev/null"}}, &container.HostConfig{Privileged: true}, nil, "")
	if err != nil {
		log.Fatal(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		log.Fatal(err)
	}

	filePaths := []string{*script, *packages, *repositories}
	for _, filePath := range filePaths {
		archive, err := archive.Tar(filePath, archive.Gzip)
		if err != nil {
			log.Fatal(err)
		}

		if err := cli.CopyToContainer(ctx, resp.ID, WORKDIR, archive, types.CopyToContainerOptions{}); err != nil {
			log.Fatal(err)
		}
	}

	if err := cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
		log.Fatal(err)
	}

	log.Println(*script, *packages, *repositories, *output, resp.ID)
}
