// projman - A cross-platform automation project manager CLI
// Author: Daniel Thornburg (as design originator)
// Language: Go (Golang)
// No dependencies other than standard library

package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"strings"
)

type Project struct {
	ID          string   `yaml:"id"`
	Name        string   `yaml:"name"`
	Status      string   `yaml:"status"`
	Tags        []string `yaml:"tags"`
	CreatedAt   string   `yaml:"created_at"`
	Description string   `yaml:"description"`
	Path        string   `yaml:"path"`
}

var (
	baseDir     = "./Projects"
	indexFile   = filepath.Join(baseDir, "_index.yaml")
	defaultDirs = []string{"Design/Drawings", "Design/Specs", "PLC/Programs", "PLC/HMI", "PLC/Symbols", "PLC/Configs", "BOM/exports", "Docs/Notes", "Tests/Simulations", "Tests/Logs", "Tags", "Tools", "Archive"}
)

func main() {
	cmd := flag.String("cmd", "", "Command: new, list, open, status, archive")
	id := flag.String("id", "", "Project ID (e.g. CP-1201)")
	name := flag.String("name", "", "Project name")
	desc := flag.String("desc", "", "Project description")
	status := flag.String("status", "active", "Project status")
	tags := flag.String("tags", "", "Comma-separated tags")
	flag.Parse()

	switch *cmd {
	case "new":
		if *id == "" || *name == "" {
			log.Fatal("Must provide -id and -name for new project")
		}
		createProject(*id, *name, *desc, *status, strings.Split(*tags, ","))
	case "list":
		listProjects()
	case "open":
		openProject(*id)
	case "status":
		showStatus(*id)
	case "archive":
		archiveProject(*id)
	default:
		fmt.Println("Unknown command. Use -cmd=new|list|open|status|archive")
	}
}

// Additional functions (createProject, listProjects, openProject, etc.) to be implemented below.
