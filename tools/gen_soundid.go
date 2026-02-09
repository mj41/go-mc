// gen_soundid generates data/soundid/soundid.go from registries.json.
package main

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

type soundEntry struct {
	Name string // e.g. "entity.allay.ambient_with_item"
	ID   int
}

func genSoundID(jsonDir, goMCRoot string) error {
	jsonPath := filepath.Join(jsonDir, "registries.json")
	outPath := filepath.Join(goMCRoot, "data", "soundid", "soundid.go")

	var regs registriesJSON
	if err := readJSON(jsonPath, &regs); err != nil {
		return fmt.Errorf("gen-soundid: %w", err)
	}

	soundEvent, ok := regs["minecraft:sound_event"]
	if !ok {
		return fmt.Errorf("gen-soundid: registry 'minecraft:sound_event' not found in registries.json")
	}

	var sounds []soundEntry
	for name, entry := range soundEvent.Entries {
		cleanName := strings.TrimPrefix(name, "minecraft:")
		cleanName = strings.ReplaceAll(cleanName, ":", ".")
		sounds = append(sounds, soundEntry{Name: cleanName, ID: entry.ProtocolID})
	}

	sort.Slice(sounds, func(i, j int) bool {
		return sounds[i].ID < sounds[j].ID
	})

	var buf strings.Builder
	generateSoundID(&buf, sounds)

	if err := writeFile(outPath, []byte(buf.String())); err != nil {
		return fmt.Errorf("gen-soundid: %w", err)
	}
	logf("gen-soundid: wrote %s (%d sounds, %d bytes)", outPath, len(sounds), buf.Len())
	return nil
}

func generateSoundID(w *strings.Builder, sounds []soundEntry) {
	w.WriteString(generatedHeader("gen_soundid.go", "registries.json"))
	w.WriteByte('\n')
	w.WriteString("package soundid\n\n")
	w.WriteString("// SoundID represents a sound ID used in the minecraft protocol.\n")
	w.WriteString("type SoundID int32\n\n")
	w.WriteString("// SoundNames - map of ids to names for sounds.\n")
	w.WriteString("var SoundNames = map[SoundID]string{\n")

	for _, s := range sounds {
		w.WriteString(fmt.Sprintf("\t%-5d: %q,\n", s.ID, s.Name))
	}

	w.WriteString("}\n\n")
	w.WriteString("// GetSoundNameByID helper method\n")
	w.WriteString("func GetSoundNameByID(id SoundID) (string, bool) {\n")
	w.WriteString("\tname, ok := SoundNames[id]\n")
	w.WriteString("\treturn name, ok\n")
	w.WriteString("}\n")
}
