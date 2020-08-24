package main

import (
	"flag"
	"fmt"
	"log"
	"os"
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

	fmt.Println(*script, *packages, *repositories, *output)
}
