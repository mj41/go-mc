// gen_biome generates level/biome/list.go from biomes.json.
package main

import (
	"fmt"
	"go/format"
	"path/filepath"
	"strings"
)

type biomeJSON struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func genBiome(jsonDir, goMCRoot string) error {
	jsonPath := filepath.Join(jsonDir, "biomes.json")
	outPath := filepath.Join(goMCRoot, "level", "biome", "list.go")

	var biomes []biomeJSON
	if err := readJSON(jsonPath, &biomes); err != nil {
		return fmt.Errorf("genBiome: %w", err)
	}

	var buf strings.Builder

	buf.WriteString(generatedHeader("gen_biome.go", "biomes.json"))
	buf.WriteString(`package biome

import (
	"errors"
	"hash/maphash"
	"math/bits"
)

// Type is the protocol ID of a biome.
type Type int

var hashSeed = maphash.MakeSeed()

func (t Type) MarshalText() (text []byte, err error) {
	if t >= 0 && int(t) < len(biomesNames) {
		return biomesNames[t], nil
	}
	return nil, errors.New("invalid type")
}

func (t *Type) UnmarshalText(text []byte) error {
	var ok bool
	*t, ok = biomesIDs[maphash.Bytes(hashSeed, text)]
	if ok {
		return nil
	}
	return errors.New("invalid type")
}

// String returns the biome id. Debugging purposes only.
func (t Type) String() string {
	if t >= 0 && int(t) < len(biomesNames) {
		return string(biomesNames[t])
	}
	return "<invalid biome type>"
}

`)

	// Write var block.
	buf.WriteString("var (\n")
	buf.WriteString("\t// BitsPerBiome reports how many bits are required to represent all possible biomes.\n")
	buf.WriteString("\tBitsPerBiome int\n")
	buf.WriteString("\tbiomesIDs    map[uint64]Type\n")
	buf.WriteString("\tbiomesNames  = [][]byte{\n")
	for _, b := range biomes {
		fmt.Fprintf(&buf, "\t\t[]byte(%q),\n", b.Name)
	}
	buf.WriteString("\t}\n")
	buf.WriteString(")\n\n")

	// Write init.
	buf.WriteString(`func init() {
	BitsPerBiome = bits.Len(uint(len(biomesNames)))
	biomesIDs = make(map[uint64]Type, len(biomesNames))
	for i, v := range biomesNames {
		h := maphash.Bytes(hashSeed, v)
		biomesIDs[h] = Type(i)
	}
}
`)

	// Format with gofmt.
	formatted, err := format.Source([]byte(buf.String()))
	if err != nil {
		return fmt.Errorf("genBiome: gofmt: %w", err)
	}

	if err := writeFile(outPath, formatted); err != nil {
		return fmt.Errorf("genBiome: %w", err)
	}
	logf("genBiome: wrote %s (%d biomes)", outPath, len(biomes))
	return nil
}
