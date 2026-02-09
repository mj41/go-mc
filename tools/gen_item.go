// gen_item generates data/item/item.go from items.json + registries.json.
package main

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

// --- JSON models ---

// itemsJSONData: map of "minecraft:<name>" -> { components: { ... } }
type itemsJSONData map[string]struct {
	Components itemComponents `json:"components"`
}

type itemComponents struct {
	MaxStackSize int      `json:"minecraft:max_stack_size"`
	ItemName     itemName `json:"minecraft:item_name"`
}

type itemName struct {
	Translate string `json:"translate"`
}

type itemEntry struct {
	VarName     string // Go identifier, e.g. "PolishedGranite"
	ID          int    // protocol_id from registries.json
	DisplayName string // Title Case, e.g. "Polished Granite"
	Name        string // short name, e.g. "polished_granite"
	StackSize   int    // max_stack_size, defaults to 64
}

func genItem(jsonDir, goMCRoot string) error {
	itemsPath := filepath.Join(jsonDir, "items.json")
	registriesPath := filepath.Join(jsonDir, "registries.json")
	outPath := filepath.Join(goMCRoot, "data", "item", "item.go")

	// --- Parse items.json ---
	var items itemsJSONData
	if err := readJSON(itemsPath, &items); err != nil {
		return fmt.Errorf("genItem: %w", err)
	}

	// --- Parse registries.json ---
	var regs registriesJSON
	if err := readJSON(registriesPath, &regs); err != nil {
		return fmt.Errorf("genItem: %w", err)
	}

	itemReg, ok := regs["minecraft:item"]
	if !ok {
		return fmt.Errorf("genItem: registry 'minecraft:item' not found in registries.json")
	}

	// --- Merge data ---
	var entries []itemEntry
	for fullName, regEntry := range itemReg.Entries {
		shortName := stripMinecraftPrefix(fullName)

		stackSize := 64 // default
		if itemData, ok := items[fullName]; ok {
			if itemData.Components.MaxStackSize > 0 {
				stackSize = itemData.Components.MaxStackSize
			}
		}

		entries = append(entries, itemEntry{
			VarName:     snakeToCamel(shortName),
			ID:          regEntry.ProtocolID,
			DisplayName: snakeToTitle(shortName),
			Name:        shortName,
			StackSize:   stackSize,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].ID < entries[j].ID
	})

	var buf strings.Builder
	generateItem(&buf, entries)

	if err := writeFile(outPath, []byte(buf.String())); err != nil {
		return fmt.Errorf("genItem: %w", err)
	}
	logf("genItem: wrote %s (%d items, %d bytes)", outPath, len(entries), buf.Len())
	return nil
}

func generateItem(w *strings.Builder, entries []itemEntry) {
	w.WriteString(generatedHeader("gen_item.go", "items.json", "registries.json"))
	w.WriteByte('\n')
	w.WriteString("// Package item stores information about items in Minecraft.\n")
	w.WriteString("package item\n\n")
	w.WriteString("// ID describes the numeric ID of an item.\n")
	w.WriteString("type ID uint32\n\n")
	w.WriteString("// Item describes information about a type of item.\n")
	w.WriteString("type Item struct {\n")
	w.WriteString("\tID          ID\n")
	w.WriteString("\tDisplayName string\n")
	w.WriteString("\tName        string\n")
	w.WriteString("\tStackSize   uint\n")
	w.WriteString("}\n\n")

	// Item variables
	w.WriteString("var (\n")
	for _, e := range entries {
		fmt.Fprintf(w, "\t%s = Item{\n", e.VarName)
		fmt.Fprintf(w, "\t\tID:          %d,\n", e.ID)
		fmt.Fprintf(w, "\t\tDisplayName: %q,\n", e.DisplayName)
		fmt.Fprintf(w, "\t\tName:        %q,\n", e.Name)
		fmt.Fprintf(w, "\t\tStackSize:   %d,\n", e.StackSize)
		w.WriteString("\t}\n")
	}
	w.WriteString(")\n\n")

	// ByID map
	w.WriteString("// ByID is a map of all items indexed by their protocol ID.\n")
	w.WriteString("var ByID = map[ID]*Item{\n")
	for _, e := range entries {
		fmt.Fprintf(w, "\t%d: &%s,\n", e.ID, e.VarName)
	}
	w.WriteString("}\n")
}
