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
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
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
	baseDir     = ""
	indexFile   = ""
	defaultDirs = []string{"Design/Drawings", "Design/Specs", "PLC/Programs", "PLC/HMI", "PLC/Symbols", "PLC/Configs", "BOM/exports", "Docs/Notes", "Tests/Simulations", "Tests/Logs", "Tags", "Tools", "Archive"}
)

func init() {
	userHome, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Unable to determine home directory: %v", err)
	}
	baseDir = filepath.Join(userHome, "Projects")
	indexFile = filepath.Join(baseDir, "_index.yaml")
}

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
		createProject(Params{*id, *name, *desc, *status, *tags})
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

func timestamp() string {
	return time.Now().Format(time.RFC3339)
}

func validateId(id string) string {
	re := regexp.MustCompile(`[^A-Z0-9\-]`)
	upper := strings.ToUpper(id)
	clean := re.ReplaceAllString(upper, "")
	return clean
}

func cleanTags(tags string) []string {
	dirtytags := strings.Split(tags, ",")
	cleantags := []string{}
	for _, t := range dirtytags {
		t = strings.TrimSpace(t)
		if t != "" {
			cleantags = append(cleantags, t)
		}
	}
	return cleantags
}

type Params struct {
	id, name, description, status, tags string
}

func createProject(p Params) {
	id := validateId(p.id)
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

	project := Project{
		ID:          id,
		Name:        p.name,
		Status:      p.status,
		Tags:        cleanTags(p.tags),
		CreatedAt:   timestamp(),
		Description: p.description,
		Path:        projectPath,
	}

	projData, err := yaml.Marshal(&project)
	if err != nil {
		log.Fatalf("Failed to marshal project YAML: %v", err)
	}

	projFile := filepath.Join(projectPath, "project.yaml")
	if err := os.WriteFile(projFile, projData, 0644); err != nil {
		log.Fatalf("Failed to write project.yaml: %v", err)
	}

	fmt.Printf("✅ Created project %s at %s\n", id, projectPath)
}

func updateProject(p Params) {
	id := validateId(p.id)
	projectPath := filepath.Join(baseDir, id)
	projFile := filepath.Join(projectPath, "project.yaml")

	data, err := os.ReadFile(projFile)
	if err != nil {
		log.Fatalf("Failed to read project.yaml: %v", err)
	}

	var project Project
	if err := yaml.Unmarshal(data, &project); err != nil {
		log.Fatalf("Failed to parse project.yaml: %v", err)
	}

	// Update fields if provided
	if p.name != "" {
		project.Name = p.name
	}
	if p.description != "" {
		project.Description = p.description
	}
	if p.status != "" {
		project.Status = p.status
	}
	if p.tags != "" {
		project.Tags = cleanTags(p.tags)
	}

	updatedData, err := yaml.Marshal(&project)
	if err != nil {
		log.Fatalf("Failed to marshal updated project.yaml: %v", err)
	}

	if err := os.WriteFile(projFile, updatedData, 0644); err != nil {
		log.Fatalf("Failed to write updated project.yaml: %v", err)
	}

	fmt.Printf("✅ Updated project %s\n", id)
}

func listProjects() {
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		log.Fatalf("Failed to read base directory: %v", err)
	}

	fmt.Printf("%-12s %-25s %-10s %-20s\n", "ID", "Name", "Status", "Created")
	fmt.Println(strings.Repeat("-", 70))
	found := false
	for _, entry := range entries {
		if !entry.IsDir() || entry.Name() == "Archive" || strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		projPath := filepath.Join(baseDir, entry.Name(), "project.yaml")
		data, err := os.ReadFile(projPath)
		if err != nil {
			log.Printf("Warning: skipping %s (no project.yaml found)\n", entry.Name())
			continue
		}

		var p Project
		if err := yaml.Unmarshal(data, &p); err != nil {
			log.Printf("Warning: failed to parse %s: %v\n", projPath, err)
			continue
		}

		fmt.Printf("%-12s %-25s %-10s %-20s\n", p.ID, p.Name, p.Status, p.CreatedAt)
		found = true
	}
	if !found {
		fmt.Println("📭 No valid projects found.")
	}
}

func openProject(id string) {
	id = strings.ToUpper(id)
	projectPath := filepath.Join(baseDir, id)
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		log.Fatalf("Project %s not found", id)
	}

	var openCmd string
	switch os := runtime.GOOS; os {
	case "windows":
		openCmd = "explorer"
	case "darwin":
		openCmd = "open"
	default:
		openCmd = "xdg-open"
	}

	cmd := exec.Command(openCmd, projectPath)
	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to open project folder: %v", err)
	}
}

func showStatus(id string) {
	id = validateId(id)
	projectPath := filepath.Join(baseDir, id, "project.yaml")

	data, err := os.ReadFile(projectPath)
	if err != nil {
		log.Fatalf("Failed to read project.yaml: %v", err)
	}

	var project Project
	if err := yaml.Unmarshal(data, &project); err != nil {
		log.Fatalf("Failed to parse project.yaml: %v", err)
	}

	fmt.Println("📄 Project Status")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("ID:          %s", project.ID)
	fmt.Printf("Name:        %s", project.Name)
	fmt.Printf("Description: %s", project.Description)
	fmt.Printf("Status:      %s", project.Status)
	fmt.Printf("Created At:  %s", project.CreatedAt)
	fmt.Printf("Tags:        %s", strings.Join(project.Tags, ", "))
	fmt.Printf("Path:        %s", project.Path)
	fmt.Println(strings.Repeat("=", 50))
}

func archiveProject(id string) {}
