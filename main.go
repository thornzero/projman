// projman - A cross-platform automation project manager CLI
// Author: Daniel Thornburg (as design originator)
// Language: Go (Golang)
// No dependencies other than standard library

package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
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
	archiveDir  = ""
	defaultDirs = []string{
		"Design/Drawings",
		"Design/Specs",
		"PLC/Programs",
		"PLC/HMI",
		"PLC/Symbols",
		"PLC/Configs",
		"BOM/exports",
		"Docs/Notes",
		"Tests/Simulations",
		"Tests/Logs",
		"Tags",
		"Tools",
		"Archive",
	}
)

func init() {
	userHome, err := os.UserHomeDir()
	errorLog("Unable to determine home directory: %v", err)
	baseDir = filepath.Join(userHome, "Projects")
	archiveDir = filepath.Join(baseDir, "Archive")
	err = os.MkdirAll(archiveDir, 0755)
	errorLog("Failed to create archive directory", err)
	indexFile = filepath.Join(baseDir, "_index.yaml")
}

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

	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	cmd := args[0]

	if DEBUG {
		fmt.Printf("Base Directory: %v\n", baseDir)
		fmt.Printf("Index: %v\n", indexFile)
		fmt.Printf("Default Directories: \n%v\n\n", defaultDirs)
	}

	switch cmd {
	case "new":
		if *id == "" || *name == "" {
			log.Fatal("Must provide -id and -name for new project")
		}
		createProject(Params{*id, *name, *desc, *status, *tags})
	case "update":
		if *id == "" {
			log.Fatal("Must provide -id to identify the project you wish to update")
		}
		updateProject(Params{*id, *name, *desc, *status, *tags})
	case "list":
		listProjects()
	case "open":
		openProject(*id)
	case "status":
		showStatus(*id)
	case "archive":
		archiveProject(*id)
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		flag.Usage()
		os.Exit(1)
	}
}

func warningLog(context string, err error) bool {
	if err != nil {
		log.Printf("⚠️  Warning in %s: %v\n", context, err)
		return true // signal caller to continue
	}
	return false
}

func errorLog(message string, err error) {
	if err != nil {
		log.Fatalf(message+": %s", err)
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

func readProjectFile(id string) Project {
	projectPath := filepath.Join(baseDir, id)
	projectFile := filepath.Join(projectPath, "project.yaml")
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		log.Fatalf("Project %s does not exist", id)
	}

	data, err := os.ReadFile(projectFile)
	errorLog("Failed to read project.yaml", err)

	var project Project
	err = yaml.Unmarshal(data, &project)
	errorLog("Failed to parse project.yaml", err)
	return project
}

func writeProjectFile(project Project) {
	data, err := yaml.Marshal(&project)
	errorLog("Failed to marshal project YAML", err)

	file := filepath.Join(project.Path, "project.yaml")
	err = os.WriteFile(file, data, 0644)
	errorLog("Failed to write project.yaml", err)
}

func createProject(p Params) {
	id := validateId(p.id)
	projectPath := filepath.Join(baseDir, id)

	_, err := os.Stat(projectPath)
	if !os.IsNotExist(err) {
		log.Fatalf("Project %s already exists", id)
	}

	err = os.MkdirAll(projectPath, 0755)
	errorLog("Failed to create project directory", err)

	for _, sub := range defaultDirs {
		subPath := filepath.Join(projectPath, sub)
		err := os.MkdirAll(subPath, 0755)
		errorLog("Failed to create subdirectory "+sub, err)
	}

	writeProjectFile(Project{
		ID:          id,
		Name:        p.name,
		Status:      p.status,
		Tags:        cleanTags(p.tags),
		CreatedAt:   timestamp(),
		Description: p.description,
		Path:        projectPath,
	})

	fmt.Printf("✅ Created project %s at %s\n", id, projectPath)
}

func updateProject(p Params) {
	id := validateId(p.id)
	project := readProjectFile(id)

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
	writeProjectFile(project)

	fmt.Printf("✅ Updated project %s\n", id)
}

func listProjects() {
	entries, err := os.ReadDir(baseDir)
	errorLog("Failed to read base directory: %v", err)

	fmt.Printf("%-12s %-25s %-10s %-20s\n", "ID", "Name", "Status", "Created")
	fmt.Println(strings.Repeat("-", 70))
	found := false
	for _, entry := range entries {
		if !entry.IsDir() || entry.Name() == "Archive" || strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		projPath := filepath.Join(baseDir, entry.Name(), "project.yaml")
		data, err := os.ReadFile(projPath)
		if warningLog(entry.Name(), err) {
			continue
		}

		var p Project
		if warningLog(projPath, yaml.Unmarshal(data, &p)) {
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
	id = validateId(id)
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
	err := cmd.Start()
	errorLog("Failed to open project folder: %v", err)
}

func showStatus(id string) {
	id = validateId(id)
	project := readProjectFile(id)

	fmt.Println("󱖫 Project Status")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("ID:          %s\n", project.ID)
	fmt.Printf("Name:        %s\n", project.Name)
	fmt.Printf("Description: %s\n", project.Description)
	fmt.Printf("Status:      %s\n", project.Status)
	fmt.Printf("Created At:  %s\n", project.CreatedAt)
	fmt.Printf("Tags:        %s\n", strings.Join(project.Tags, ", "))
	fmt.Printf("Path:        %s\n", project.Path)
	fmt.Println(strings.Repeat("=", 50))
}

func archiveProject(id string) {
	id = validateId(id)
	project := readProjectFile(id)

	archivePath := filepath.Join(archiveDir, id+".zip")
	zipFile, err := os.Create(archivePath)
	errorLog("Failed to create archive file", err)
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.Walk(project.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(project.Path, path)
		if err != nil {
			return err
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		fWriter, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(fWriter, file)
		return err
	})
	errorLog("Failed to archive project: %v", err)

	project.Status = "archived"

	writeProjectFile(project)
	fmt.Printf("📦 Archived project %s to %s\n", id, archivePath)

	getArchiveSize(id)
}

func getArchiveSize(id string) {
	archivePath := filepath.Join(archiveDir, id+".zip")
	info, err := os.Stat(archivePath)
	if err == nil {
		fmt.Printf("📦 Archive size: %.2f KB\n", float64(info.Size())/1024)
	}
}
