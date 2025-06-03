package projman

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type TagRow struct {
	Category string
	Subcat   string
	Name     string
}

type TagSpec struct {
	Format string `yaml:"format"` // e.g. "{category}-{subcat}-{id}"
	Start  int    `yaml:"start"`
}

type TagAssignment struct {
	ID       string `yaml:"id"`
	Category string `yaml:"category"`
	Subcat   string `yaml:"subcat"`
	Name     string `yaml:"name"`
}

func GenerateTags(csvPath, specPath, outputPath string) error {
	specData, err := os.ReadFile(specPath)
	if err != nil {
		return fmt.Errorf("reading spec: %w", err)
	}
	var spec TagSpec
	if err := yaml.Unmarshal(specData, &spec); err != nil {
		return fmt.Errorf("parsing spec: %w", err)
	}

	f, err := os.Open(csvPath)
	if err != nil {
		return fmt.Errorf("open csv: %w", err)
	}
	defer f.Close()
	reader := csv.NewReader(f)

	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("read csv: %w", err)
	}

	var assignments []TagAssignment
	counter := make(map[string]int)

	for i, row := range records {
		if i == 0 {
			continue // skip header
		}
		if len(row) < 3 {
			continue
		}

		r := TagRow{
			Category: strings.TrimSpace(row[0]),
			Subcat:   strings.TrimSpace(row[1]),
			Name:     strings.TrimSpace(row[2]),
		}

		key := r.Category + "-" + r.Subcat
		counter[key]++
		id := spec.Start + counter[key] - 1

		tagID := strings.ReplaceAll(spec.Format, "{category}", r.Category)
		tagID = strings.ReplaceAll(tagID, "{subcat}", r.Subcat)
		tagID = strings.ReplaceAll(tagID, "{id}", fmt.Sprintf("%02d", id))

		assignments = append(assignments, TagAssignment{
			ID:       tagID,
			Category: r.Category,
			Subcat:   r.Subcat,
			Name:     r.Name,
		})
	}

	outData, err := yaml.Marshal(assignments)
	if err != nil {
		return fmt.Errorf("marshal yaml: %w", err)
	}

	if err := os.WriteFile(outputPath, outData, 0644); err != nil {
		return fmt.Errorf("write output: %w", err)
	}

	return nil
}
