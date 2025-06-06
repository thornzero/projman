package app

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// app project metadata
type Project struct {
	ID          string   `yaml:"id"`
	Name        string   `yaml:"name"`
	Status      string   `yaml:"status"`
	Tags        []string `yaml:"tags"`
	CreatedAt   string   `yaml:"created_at"`
	Description string   `yaml:"description"`
	Path        string   `yaml:"path"`
}

// Struct for CLI/TUI parameters
type Params struct {
	ID, Name, Description, Status, Tags string
}

func Timestamp() string {
	return time.Now().Format(time.RFC3339)
}

func ValidateID(id string) string {
	re := regexp.MustCompile(`[^A-Z0-9\-]`)
	upper := strings.ToUpper(id)
	return re.ReplaceAllString(upper, "")
}

func CleanTags(tags string) []string {
	split := strings.Split(tags, ",")
	clean := []string{}
	for _, t := range split {
		t = strings.TrimSpace(t)
		if t != "" {
			clean = append(clean, t)
		}
	}
	return clean
}

func ReadProjectFile(baseDir, id string) (Project, error) {
	var p Project
	id = ValidateID(id)
	path := filepath.Join(baseDir, id, "project.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return p, err
	}
	err = yaml.Unmarshal(data, &p)
	return p, err
}

func WriteProjectFile(p Project) error {
	data, err := yaml.Marshal(&p)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(p.Path, "project.yaml"), data, 0644)
}

// 🔓 Public API

func GetDefaultBaseDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to determine home directory: %v", err)
	}
	return filepath.Join(home, "Projects")
}

func CreateProject(baseDir string, p Params) {
	id := ValidateID(p.ID)
	path := filepath.Join(baseDir, id)

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		log.Fatalf("Project %s already exists", id)
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		log.Fatalf("Failed to create project directory: %v", err)
	}

	defaultDirs := []string{
		"Docs", "Planning", "Logs", "Exports",
	}
	for _, d := range defaultDirs {
		sub := filepath.Join(path, d)
		if err := os.MkdirAll(sub, 0755); err != nil {
			log.Fatalf("Failed to create subdirectory %s: %v", sub, err)
		}
	}

	proj := Project{
		ID:          id,
		Name:        p.Name,
		Description: p.Description,
		Status:      p.Status,
		Tags:        CleanTags(p.Tags),
		CreatedAt:   Timestamp(),
		Path:        path,
	}

	if err := WriteProjectFile(proj); err != nil {
		log.Fatalf("Failed to write project file: %v", err)
	}

	fmt.Printf("✅ Created project %s at %s\n", id, path)
}

func ListProjects(baseDir string) {
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

		p, err := ReadProjectFile(baseDir, entry.Name())
		if err != nil {
			log.Printf("⚠️  Skipping %s: %v", entry.Name(), err)
			continue
		}

		fmt.Printf("%-12s %-25s %-10s %-20s\n", p.ID, p.Name, p.Status, p.CreatedAt)
		found = true
	}

	if !found {
		fmt.Println("📭 No valid projects found.")
	}
}

func ShowStatus(baseDir, id string) {
	p, err := ReadProjectFile(baseDir, id)
	if err != nil {
		log.Fatalf("Failed to read project: %v", err)
	}

	fmt.Println("󱖫 Project Status")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("ID:          %s\n", p.ID)
	fmt.Printf("Name:        %s\n", p.Name)
	fmt.Printf("Description: %s\n", p.Description)
	fmt.Printf("Status:      %s\n", p.Status)
	fmt.Printf("Created At:  %s\n", p.CreatedAt)
	fmt.Printf("Tags:        %s\n", strings.Join(p.Tags, ", "))
	fmt.Printf("Path:        %s\n", p.Path)
	fmt.Println(strings.Repeat("=", 50))
}
