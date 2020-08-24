package main

import (
	"archive/tar"
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

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

	DOCKER_IMAGE_URL = "docker.io/library/alpine:edge"

	WORKDIR = "/tmp"
)

func main() {
	setupScriptFile := flag.String("script", SETUP_SCRIPT_FILE_DEFAULT, "Setup script file")
	packageListFile := flag.String("packages", PACKAGE_LIST_FILE_DEFAULT, "Package list file")
	repositoryListFile := flag.String("repositories", REPOSITORY_LIST_FILE_DEFAULT, "Repository list file")
	outputImageFile := flag.String("output", OUTPUT_IMAGE_FILE_DEFAULT, "Output image file")

	flag.Parse()

	ctx := context.Background()

	filePaths := [][2]string{{*setupScriptFile, SETUP_SCRIPT_FILE_DEFAULT}, {*packageListFile, PACKAGE_LIST_FILE_DEFAULT}, {*repositoryListFile, REPOSITORY_LIST_FILE_DEFAULT}}
	for _, filePath := range filePaths {
		info, err := os.Stat(filePath[0])
		if err != nil {
			log.Fatal("could not check if file exists", filePath[0], err)
		}

		if info == nil {
			log.Fatal("file doesn't exist", filePath[0], err)
		}
	}

	log.Println("connecting to Docker daemon")
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatal("could not connect to Docker daemon", err)
	}

	log.Println("pulling Alpine Linux image")
	out, err := cli.ImagePull(ctx, DOCKER_IMAGE_URL, types.ImagePullOptions{})
	if err != nil {
		log.Fatal("could not pull Alpine Linux image", err)
	}
	defer out.Close()
	if _, err := ioutil.ReadAll(out); err != nil {
		log.Fatal("could not write pulled Alpine Linux image", err)
	}

	log.Println("creating Alpine Linux container")
	resp, err := cli.ContainerCreate(ctx, &container.Config{Image: DOCKER_IMAGE_URL, Cmd: []string{"tail", "-f", "/dev/null"}}, &container.HostConfig{Privileged: true, DNS: []string{"8.8.8.8"}}, nil, "")
	if err != nil {
		log.Fatal("could not create Alpine Linux container", err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		log.Fatal("could not start Alpine Linux container", err)
	}
	defer func() {
		log.Println("stopping and removing Alpine Linux container")

		if err := cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
			log.Fatal("could not remove Alpine Linux container", err)
		}
	}()

	log.Println("copying files to Alpine Linux container")
	cmds := [][]string{}
	for _, filePath := range filePaths {
		archive, err := archive.Tar(filePath[0], archive.Gzip)
		if err != nil {
			log.Fatal("could not create tar archive for file to copy into Alpine Linux container", filePath, err)
		}

		if err := cli.CopyToContainer(ctx, resp.ID, WORKDIR, archive, types.CopyToContainerOptions{}); err != nil {
			log.Fatal("could not copy tar archive for file to copy into Alpine Linux container", filePath, err)
		}

		cmds = append(cmds, []string{"mv", WORKDIR + "/" + filepath.Base(filePath[0]), WORKDIR + "/" + filePath[1]})
	}

	log.Println("building image in Alpine Linux container")
	cmds = append(cmds, []string{"chmod", "+x", WORKDIR + "/" + SETUP_SCRIPT_FILE_DEFAULT}, []string{"apk", "add", "alpine-make-vm-image"}, []string{"sh", "-c", fmt.Sprintf("alpine-make-vm-image --image-format qcow2 --repositories-file %v --packages \"$(cat %v)\" --script-chroot %v %v", REPOSITORY_LIST_FILE_DEFAULT, PACKAGE_LIST_FILE_DEFAULT, OUTPUT_IMAGE_FILE_DEFAULT, SETUP_SCRIPT_FILE_DEFAULT)})
	for _, cmd := range cmds {
		exec, err := cli.ContainerExecCreate(ctx, resp.ID, types.ExecConfig{Cmd: cmd, WorkingDir: WORKDIR})
		if err != nil {
			log.Fatal("could not create exec", exec.ID, err)
		}

		if err := cli.ContainerExecStart(ctx, exec.ID, types.ExecStartCheck{}); err != nil {
			log.Fatal("could not start exec", exec.ID, err)
		}

		running := true
		for running {
			info, err := cli.ContainerExecInspect(ctx, exec.ID)
			if err != nil {
				log.Fatal("could not inspect exec", exec.ID, err)
			}

			if info.ExitCode != 0 {
				log.Fatal("could not run command in Alpine Linux container, exited with non-zero exit code", info.ExitCode, cmd)
			}

			running = info.Running
		}
	}

	log.Println("copying image from Alpine Linux container to host")
	localFile, err := os.Create(*outputImageFile)
	if err != nil {
		log.Fatal("could not create output file", *outputImageFile, err)
	}

	tarStream, _, err := cli.CopyFromContainer(ctx, resp.ID, WORKDIR+"/"+OUTPUT_IMAGE_FILE_DEFAULT)
	if err != nil {
		log.Fatal("could not request tar stream from Docker daemon", err)
	}

	remoteFile := tar.NewReader(tarStream)
	if _, err := remoteFile.Next(); err != nil {
		log.Fatal("could not read tar archive", err)
	}

	if _, err := io.Copy(localFile, remoteFile); err != nil {
		log.Fatal("could not write to output file", *outputImageFile, err)
	}
}
