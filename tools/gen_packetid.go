package main

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
)

// PacketsJSON mirrors the structure of packets.json from MC --all.
// Top-level keys: phase names (login, status, configuration, play, handshake).
// Each phase has "clientbound" and/or "serverbound".
// Each direction maps packet name → {"protocol_id": N}.
type PacketsJSON map[string]map[string]map[string]struct {
	ProtocolID int `json:"protocol_id"`
}

type packetIDEntry struct {
	Name       string // e.g. "minecraft:add_entity"
	ProtocolID int
	GoName     string // e.g. "ClientboundAddEntity"
}

// phaseOrder defines which phases to generate and in what order.
// Handshake is excluded (only has serverbound intention, handled separately).
var phaseOrder = []string{"login", "status", "configuration", "play"}

// phaseAbbrev maps JSON phase name to Go name prefix.
// Play phase has no prefix (just Clientbound/Serverbound).
var phaseAbbrev = map[string]string{
	"login":         "Login",
	"status":        "Status",
	"configuration": "Config",
	"play":          "",
}

// phaseComment maps JSON phase name to section comment.
var phaseComment = map[string]string{
	"login":         "Login",
	"status":        "Status",
	"configuration": "Configuration",
	"play":          "Game",
}

func genPacketID(jsonDir, goMCRoot string) error {
	jsonPath := filepath.Join(jsonDir, "packets.json")
	outPath := filepath.Join(goMCRoot, "data", "packetid", "packetid.go")

	var packets PacketsJSON
	if err := readJSON(jsonPath, &packets); err != nil {
		return fmt.Errorf("genPacketID: %w", err)
	}

	var buf strings.Builder
	generatePacketID(&buf, packets)

	if err := writeFile(outPath, []byte(buf.String())); err != nil {
		return fmt.Errorf("genPacketID: %w", err)
	}
	logf("gen-packetid: wrote %s (%d bytes)", outPath, buf.Len())
	return nil
}

func generatePacketID(w *strings.Builder, packets PacketsJSON) {
	w.WriteString(generatedHeader("gen_packetid.go", "packets.json"))
	w.WriteByte('\n')
	w.WriteString("package packetid\n\n")
	w.WriteString("//go:generate stringer -type ClientboundPacketID\n")
	w.WriteString("//go:generate stringer -type ServerboundPacketID\n")
	w.WriteString("type (\n")
	w.WriteString("\tClientboundPacketID int32\n")
	w.WriteString("\tServerboundPacketID int32\n")
	w.WriteString(")\n")

	for _, phase := range phaseOrder {
		phaseData, ok := packets[phase]
		if !ok {
			continue
		}

		for _, dir := range []string{"clientbound", "serverbound"} {
			dirData, ok := phaseData[dir]
			if !ok {
				continue
			}

			abbrev := phaseAbbrev[phase]
			comment := phaseComment[phase]
			dirTitle := strings.ToUpper(dir[:1]) + dir[1:]

			w.WriteString(fmt.Sprintf("\n// %s %s\n", comment, dirTitle))
			w.WriteString("const (\n")

			pkts := sortPacketIDEntries(dirData, dir, abbrev)

			// Verify sequential IDs — if all sequential from 0, use iota.
			sequential := true
			for i, p := range pkts {
				if p.ProtocolID != i {
					sequential = false
					break
				}
			}

			for i, p := range pkts {
				if sequential {
					if i == 0 {
						typeName := dirTypeName(dir)
						w.WriteString(fmt.Sprintf("\t%s %s = iota\n", p.GoName, typeName))
					} else {
						w.WriteString(fmt.Sprintf("\t%s\n", p.GoName))
					}
				} else {
					typeName := dirTypeName(dir)
					w.WriteString(fmt.Sprintf("\t%s %s = %d\n", p.GoName, typeName, p.ProtocolID))
				}
			}

			// Add guard sentinel for play phase.
			if phase == "play" {
				prefix := dirPrefix(dir)
				w.WriteString(fmt.Sprintf("\t%sPacketIDGuard\n", prefix))
			}

			w.WriteString(")\n")
		}
	}
}

func sortPacketIDEntries(dirData map[string]struct {
	ProtocolID int `json:"protocol_id"`
}, dir, abbrev string) []packetIDEntry {
	var pkts []packetIDEntry
	prefix := dirPrefix(dir)

	for name, info := range dirData {
		goName := packetGoName(name, prefix, abbrev)
		pkts = append(pkts, packetIDEntry{
			Name:       name,
			ProtocolID: info.ProtocolID,
			GoName:     goName,
		})
	}

	sort.Slice(pkts, func(i, j int) bool {
		return pkts[i].ProtocolID < pkts[j].ProtocolID
	})

	return pkts
}

func packetGoName(mcName, dirPrefix, phaseAbbrev string) string {
	// Strip "minecraft:" prefix.
	name := strings.TrimPrefix(mcName, "minecraft:")

	// Replace "/" and "_" with space, then CamelCase.
	name = strings.ReplaceAll(name, "/", " ")
	name = strings.ReplaceAll(name, "_", " ")

	var sb strings.Builder
	sb.WriteString(dirPrefix)
	sb.WriteString(phaseAbbrev)

	for _, word := range strings.Fields(name) {
		if len(word) == 0 {
			continue
		}
		runes := []rune(word)
		runes[0] = unicode.ToUpper(runes[0])
		sb.WriteString(string(runes))
	}

	return sb.String()
}

func dirPrefix(dir string) string {
	if dir == "clientbound" {
		return "Clientbound"
	}
	return "Serverbound"
}

func dirTypeName(dir string) string {
	if dir == "clientbound" {
		return "ClientboundPacketID"
	}
	return "ServerboundPacketID"
}
