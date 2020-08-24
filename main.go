package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

const (
	SETUP_SCRIPT_FILE_DEFAULT    = "setup.sh"
	PACKAGE_LIST_FILE_DEFAULT    = "packages.txt"
	REPOSITORY_LIST_FILE_DEFAULT = "repositories.txt"
	OUTPUT_IMAGE_FILE_DEFAULT    = "alpine.qcow2"
	OUTPUT_IMAGE_SIZE_DEFAULT    = "250M"

	DOCKER_IMAGE_URL = "docker.io/library/alpine:edge"
	DOCKER_IMAGE     = "alpine:edge"

	WORKDIR = "/tmp"
)

func main() {
	setupScriptFile := flag.String("script", SETUP_SCRIPT_FILE_DEFAULT, "Setup script file")
	packageListFile := flag.String("packages", PACKAGE_LIST_FILE_DEFAULT, "Package list file")
	repositoryListFile := flag.String("repositories", REPOSITORY_LIST_FILE_DEFAULT, "Repository list file")
	outputImageFile := flag.String("output", OUTPUT_IMAGE_FILE_DEFAULT, "Output image file")
	outputImageSize := flag.String("size", OUTPUT_IMAGE_SIZE_DEFAULT, "Output image size")

	flag.Parse()

	if err := os.Chmod(*setupScriptFile, 0777); err != nil {
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

	filePaths := [][2]string{{*setupScriptFile, SETUP_SCRIPT_FILE_DEFAULT}, {*packageListFile, PACKAGE_LIST_FILE_DEFAULT}, {*repositoryListFile, REPOSITORY_LIST_FILE_DEFAULT}}
	cmds := [][]string{}
	for _, filePath := range filePaths {
		archive, err := archive.Tar(filePath[0], archive.Gzip)
		if err != nil {
			log.Fatal(err)
		}

		if err := cli.CopyToContainer(ctx, resp.ID, WORKDIR, archive, types.CopyToContainerOptions{}); err != nil {
			log.Fatal(err)
		}

		cmds = append(cmds, []string{"mv", path.Join(WORKDIR, path.Base(filePath[0])), path.Join(WORKDIR, filePath[1])})
	}

	cmds = append(cmds, []string{"apk", "add", "alpine-make-vm-image"}, []string{"sh", "-c", fmt.Sprintf("alpine-make-vm-image --image-format qcow2 --image-size %v --repositories-file %v --packages \"$(cat %v)\" --script-chroot %v %v", *outputImageSize, REPOSITORY_LIST_FILE_DEFAULT, PACKAGE_LIST_FILE_DEFAULT, OUTPUT_IMAGE_FILE_DEFAULT, SETUP_SCRIPT_FILE_DEFAULT)})
	for _, cmd := range cmds {
		exec, err := cli.ContainerExecCreate(ctx, resp.ID, types.ExecConfig{Cmd: cmd, WorkingDir: WORKDIR})
		if err != nil {
			log.Fatal(err)
		}

		if err := cli.ContainerExecStart(ctx, exec.ID, types.ExecStartCheck{}); err != nil {
			log.Fatal(err)
		}

		running := true
		for running {
			info, err := cli.ContainerExecInspect(ctx, exec.ID)
			if err != nil {
				log.Fatal(err)
			}

			if info.ExitCode != 0 {
				log.Fatal("failed to run command in Docker", cmd)
			}

			running = info.Running
		}
	}

	// if err := cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
	// log.Fatal(err)
	// }

	log.Println(*setupScriptFile, *packageListFile, *repositoryListFile, *outputImageFile, resp.ID)
}
