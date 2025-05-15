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
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
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
	DEBUG       = true
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

	if DEBUG {
		fmt.Printf("Base Directory: %v\n", baseDir)
		fmt.Printf("Index: %v\n", indexFile)
		fmt.Printf("Default Directories: \n%v\n\n", defaultDirs)
	}

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
func createProject(id string, name string, description string, status string, tags []string) {
	id = strings.ToUpper(id)
	projectPath := filepath.Join(baseDir, id)
	if _, err := os.Stat(projectPath); !os.IsNotExist(err) {
		log.Fatalf("Project %s already exists", id)
	}

	if err := os.MkdirAll(projectPath, 0755); err != nil {
		log.Fatalf("Failed to create project directory: %v", err)
	}

	for _, sub := range defaultDirs {
		subPath := filepath.Join(projectPath, sub)
		if err := os.MkdirAll(subPath, 0755); err != nil {
			log.Fatalf("Failed to create subdirectory %s: %v", sub, err)
		}
	}

	p := Project{
		ID:          id,
		Name:        name,
		Status:      status,
		Tags:        tags,
		CreatedAt:   time.Now().Format(time.RFC3339),
		Description: description,
		Path:        projectPath,
	}

	projData, err := yaml.Marshal(&p)
	if err != nil {
		log.Fatalf("Failed to marshal project YAML: %v", err)
	}

	projFile := filepath.Join(projectPath, "project.yaml")
	if err := os.WriteFile(projFile, projData, 0644); err != nil {
		log.Fatalf("Failed to write project.yaml: %v", err)
	}

	fmt.Printf("✅ Created project %s at %s\n", id, projectPath)
}

func listProjects() {
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		log.Fatalf("Failed to read base directory: %v", err)
	}

	fmt.Printf("%-12s %-25s %-10s %-20s", "ID", "Name", "Status", "Created")
	fmt.Println(strings.Repeat("-", 70))

	for _, entry := range entries {
		if !entry.IsDir() || entry.Name() == "Archive" || strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		projPath := filepath.Join(baseDir, entry.Name(), "project.yaml")
		data, err := os.ReadFile(projPath)
		if err != nil {
			log.Printf("Warning: skipping %s (no project.yaml found)", entry.Name())
			continue
		}

		var p Project
		if err := yaml.Unmarshal(data, &p); err != nil {
			log.Printf("Warning: failed to parse %s: %v", projPath, err)
			continue
		}

		fmt.Printf("%-12s %-25s %-10s %-20s", p.ID, p.Name, p.Status, p.CreatedAt)
	}
}

func openProject(id string) {

}

func showStatus(id string) {}

func archiveProject(id string) {}
