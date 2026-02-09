// gen_entity generates data/entity/entity.go from entities.json.
package main

import (
	"fmt"
	"go/format"
	"math"
	"path/filepath"
	"strings"
)

type entityJSON struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Width    float64 `json:"width"`
	Height   float64 `json:"height"`
	Category string  `json:"category"`
}

func genEntity(jsonDir, goMCRoot string) error {
	jsonPath := filepath.Join(jsonDir, "entities.json")
	outPath := filepath.Join(goMCRoot, "data", "entity", "entity.go")

	var entities []entityJSON
	if err := readJSON(jsonPath, &entities); err != nil {
		return fmt.Errorf("genEntity: %w", err)
	}

	var buf strings.Builder

	buf.WriteString(generatedHeader("gen_entity.go", "entities.json"))
	buf.WriteString(`// Package entity stores information about entities in Minecraft.
package entity

// ID describes the numeric ID of an entity.
type ID uint32

// Entity describes information about a type of entity.
type Entity struct {
	ID          ID
	InternalID  uint32
	DisplayName string
	Name        string
	Width       float64
	Height      float64
	Type        string
}

`)

	// Write var block with entity declarations.
	buf.WriteString("var (\n")
	for _, e := range entities {
		shortName := stripMinecraftPrefix(e.Name)
		goName := snakeToCamel(shortName)
		displayName := snakeToTitle(shortName)
		w := formatEntityFloat(e.Width)
		h := formatEntityFloat(e.Height)

		fmt.Fprintf(&buf, "\t%s = Entity{\n", goName)
		fmt.Fprintf(&buf, "\t\tID:          %d,\n", e.ID)
		fmt.Fprintf(&buf, "\t\tInternalID:  %d,\n", e.ID)
		fmt.Fprintf(&buf, "\t\tDisplayName: %q,\n", displayName)
		fmt.Fprintf(&buf, "\t\tName:        %q,\n", shortName)
		fmt.Fprintf(&buf, "\t\tWidth:       %s,\n", w)
		fmt.Fprintf(&buf, "\t\tHeight:      %s,\n", h)
		fmt.Fprintf(&buf, "\t\tType:        %q,\n", e.Category)
		buf.WriteString("\t}\n")
	}
	buf.WriteString(")\n\n")

	// Write ByID map.
	buf.WriteString("// ByID is an index of minecraft entities by their ID.\n")
	buf.WriteString("var ByID = map[ID]*Entity{\n")

	// Determine max ID width for alignment.
	maxID := 0
	for _, e := range entities {
		if e.ID > maxID {
			maxID = e.ID
		}
	}
	idWidth := len(fmt.Sprintf("%d", maxID))

	for _, e := range entities {
		shortName := stripMinecraftPrefix(e.Name)
		goName := snakeToCamel(shortName)
		fmt.Fprintf(&buf, "\t%-*d: &%s,\n", idWidth+1, e.ID, goName)
	}
	buf.WriteString("}\n")

	// Format with gofmt.
	formatted, err := format.Source([]byte(buf.String()))
	if err != nil {
		return fmt.Errorf("genEntity: gofmt: %w", err)
	}

	if err := writeFile(outPath, formatted); err != nil {
		return fmt.Errorf("genEntity: %w", err)
	}
	logf("genEntity: wrote %s (%d entities)", outPath, len(entities))
	return nil
}

func roundSig(f float64, digits int) float64 {
	if f == 0 {
		return 0
	}
	d := math.Ceil(math.Log10(math.Abs(f)))
	pow := math.Pow(10, float64(digits)-d)
	return math.Round(f*pow) / pow
}

func formatEntityFloat(f float64) string {
	rounded := roundSig(f, 7)
	return fmt.Sprintf("%g", rounded)
}
