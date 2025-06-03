// projman - A cross-platform automation project manager CLI
// Author: Daniel Thornburg (as design originator)
// Language: Go (Golang)
// No dependencies other than standard library

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/thornzero/projman/pkg/projman"
)

func main() {
	flag.Usage = func() {
		fmt.Println("Usage: projman <command> [options]")
		fmt.Println("Commands: new, list, open, status, update, archive")
		fmt.Println("Options:")
		flag.PrintDefaults()
	}

	id := flag.String("id", "", "Project ID (e.g. CP-1201)")
	name := flag.String("name", "", "Project name")
	desc := flag.String("desc", "", "Project description")
	status := flag.String("status", "active", "Project status")
	tags := flag.String("tags", "", "Comma-separated tags")
	baseDir := projman.GetDefaultBaseDir()

	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	cmd := args[0]

	switch cmd {
	case "new":
		if *id == "" || *name == "" {
			log.Fatal("Must provide -id and -name for new project")
		}
		projman.CreateProject(baseDir, projman.Params{
			ID:          *id,
			Name:        *name,
			Description: *desc,
			Status:      *status,
			Tags:        *tags,
		})

	//case "update":
	//	updateProject(Params{*id, *name, *desc, *status, *tags})
	case "list":
		projman.ListProjects(baseDir)

	//case "open":
	//  projman.openProject(baseDir,*id)

	case "status":
		if *id == "" {
			log.Fatal("Must provide -id for status")
		}
		projman.ShowStatus(baseDir, *id)

	//case "archive":
	//  projman.archiveProject(*id)

	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		flag.Usage()
		os.Exit(1)
	}
}
