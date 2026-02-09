// gen_blockentities generates level/block/blockentity.go and
// level/block/blockentities.go from block_entities.json.
package main

import (
	"fmt"
	"go/format"
	"path/filepath"
	"sort"
	"strings"
	"unicode/utf8"
)

type blockEntityJSON struct {
	Name        string   `json:"name"`
	ValidBlocks []string `json:"valid_blocks"`
}

func genBlockEntities(jsonDir, goMCRoot string) error {
	jsonPath := filepath.Join(jsonDir, "block_entities.json")
	outDir := filepath.Join(goMCRoot, "level", "block")

	var entities []blockEntityJSON
	if err := readJSON(jsonPath, &entities); err != nil {
		return fmt.Errorf("genBlockEntities: %w", err)
	}

	// Sort by registry name for stable output.
	sort.Slice(entities, func(i, j int) bool {
		return entities[i].Name < entities[j].Name
	})

	if err := generateBlockEntityFile(entities, outDir); err != nil {
		return fmt.Errorf("genBlockEntities: blockentity.go: %w", err)
	}

	if err := generateBlockEntitiesFile(entities, outDir); err != nil {
		return fmt.Errorf("genBlockEntities: blockentities.go: %w", err)
	}

	logf("genBlockEntities: wrote blockentity.go + blockentities.go (%d block entity types)", len(entities))
	return nil
}

// generateBlockEntityFile generates blockentity.go with Entity interface,
// struct types, EntityType, and EntityTypes map init().
func generateBlockEntityFile(entities []blockEntityJSON, outDir string) error {
	var buf strings.Builder

	buf.WriteString(generatedHeader("gen_blockentities.go", "block_entities.json"))
	buf.WriteByte('\n')
	buf.WriteString("package block\n\n")

	// Entity interface.
	buf.WriteString("type Entity interface {\n")
	buf.WriteString("\tID() string\n")
	buf.WriteString("\tIsValidBlock(block Block) bool\n")
	buf.WriteString("}\n\n")

	// Struct type definitions.
	buf.WriteString("type (\n")

	// Compute max name width for alignment.
	maxWidth := 0
	for _, e := range entities {
		goName := blockEntityGoTypeName(e.Name)
		if len(goName) > maxWidth {
			maxWidth = len(goName)
		}
	}

	for _, e := range entities {
		goName := blockEntityGoTypeName(e.Name)
		padding := strings.Repeat(" ", maxWidth-len(goName))
		fmt.Fprintf(&buf, "\t%s%s struct{}\n", goName, padding)
	}
	buf.WriteString(")\n\n")

	// EntityType type.
	buf.WriteString("type EntityType int32\n\n")

	// EntityTypes map + init.
	buf.WriteString("var EntityTypes map[string]EntityType\n\n")
	buf.WriteString("func init() {\n")
	buf.WriteString("\tEntityTypes = make(map[string]EntityType, len(EntityList))\n")
	buf.WriteString("\tfor i, v := range EntityList {\n")
	buf.WriteString("\t\tEntityTypes[v.ID()] = EntityType(i)\n")
	buf.WriteString("\t}\n")
	buf.WriteString("}\n")

	formatted, err := format.Source([]byte(buf.String()))
	if err != nil {
		return fmt.Errorf("gofmt: %w", err)
	}

	outPath := filepath.Join(outDir, "blockentity.go")
	return writeFile(outPath, formatted)
}

// generateBlockEntitiesFile generates blockentities.go with EntityList,
// ID(), and IsValidBlock() methods.
func generateBlockEntitiesFile(entities []blockEntityJSON, outDir string) error {
	var buf strings.Builder

	buf.WriteString(generatedHeader("gen_blockentities.go", "block_entities.json"))
	buf.WriteByte('\n')
	buf.WriteString("package block\n\n")

	// EntityList array.
	buf.WriteString("var EntityList = [...]Entity{\n")
	for _, e := range entities {
		fmt.Fprintf(&buf, "\t%s{},\n", blockEntityGoTypeName(e.Name))
	}
	buf.WriteString("}\n\n")

	// ID() methods.
	for _, e := range entities {
		goName := blockEntityGoTypeName(e.Name)
		fmt.Fprintf(&buf, "func (%s) ID() string { return %q }\n", goName, e.Name)
	}
	buf.WriteString("\n")

	// IsValidBlock() methods.
	for _, e := range entities {
		goName := blockEntityGoTypeName(e.Name)
		receiver := strings.ToLower(goName[:1])

		// Sort valid blocks for consistent output.
		blocks := make([]string, len(e.ValidBlocks))
		copy(blocks, e.ValidBlocks)
		sort.Strings(blocks)

		if len(blocks) == 1 {
			fmt.Fprintf(&buf, "func (%s %s) IsValidBlock(block Block) bool {\n", receiver, goName)
			fmt.Fprintf(&buf, "\treturn block.ID() == %q\n", blocks[0])
			buf.WriteString("}\n\n")
		} else {
			fmt.Fprintf(&buf, "func (%s %s) IsValidBlock(block Block) bool {\n", receiver, goName)
			buf.WriteString("\tswitch block.ID() {\n")
			buf.WriteString("\tcase ")
			for i, b := range blocks {
				if i > 0 {
					buf.WriteString(",\n\t\t")
				}
				fmt.Fprintf(&buf, "%q", b)
			}
			buf.WriteString(":\n")
			buf.WriteString("\t\treturn true\n")
			buf.WriteString("\tdefault:\n")
			buf.WriteString("\t\treturn false\n")
			buf.WriteString("\t}\n")
			buf.WriteString("}\n\n")
		}
	}

	formatted, err := format.Source([]byte(buf.String()))
	if err != nil {
		return fmt.Errorf("gofmt: %w", err)
	}

	outPath := filepath.Join(outDir, "blockentities.go")
	return writeFile(outPath, formatted)
}

// blockEntityGoTypeName converts a registry name like "minecraft:mob_spawner"
// to a Go type name like "MobSpawnerEntity".
func blockEntityGoTypeName(registryName string) string {
	name := stripMinecraftPrefix(registryName)

	// Convert snake_case to CamelCase.
	parts := strings.Split(name, "_")
	for i, p := range parts {
		if len(p) > 0 {
			r, size := utf8.DecodeRuneInString(p)
			parts[i] = strings.ToUpper(string(r)) + p[size:]
		}
	}

	return strings.Join(parts, "") + "Entity"
}
