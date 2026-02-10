// gen_lang generates data/lang/en-us/en_us.go from en_us.json.
//
// The en_us.json is extracted from the MC inner jar during --extract.
// Java format strings (%2$s) are converted to Go format strings (%[2]s).
package main

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var javaFmtArg = regexp.MustCompile(`%(\d+)\$s`)

// transJavaFormat converts Java positional format strings to Go format.
// e.g., "%2$s" → "%[2]s"
func transJavaFormat(s string) string {
	return javaFmtArg.ReplaceAllStringFunc(s, func(m string) string {
		var idx int
		fmt.Sscanf(m, "%%%d$s", &idx)
		return fmt.Sprintf("%%[%d]s", idx)
	})
}

func genLang(jsonDir, goMCRoot string) error {
	jsonPath := filepath.Join(jsonDir, "en_us.json")
	outPath := filepath.Join(goMCRoot, "data", "lang", "en-us", "en_us.go")

	var langMap map[string]string
	if err := readJSON(jsonPath, &langMap); err != nil {
		return fmt.Errorf("genLang: %w", err)
	}

	// Apply Java→Go format string conversion.
	for k, v := range langMap {
		if javaFmtArg.MatchString(v) {
			langMap[k] = transJavaFormat(v)
		}
	}

	// Sort keys for deterministic output.
	keys := make([]string, 0, len(langMap))
	for k := range langMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buf strings.Builder
	buf.WriteString(generatedHeader("gen_lang.go", "en_us.json"))
	buf.WriteString("\npackage en_us\n\n")
	buf.WriteString("var Map = map[string]string{\n")
	for _, k := range keys {
		fmt.Fprintf(&buf, "\t%q: %q,\n", k, langMap[k])
	}
	buf.WriteString("}\n")

	if err := writeFile(outPath, []byte(buf.String())); err != nil {
		return fmt.Errorf("genLang: %w", err)
	}
	logf("genLang: wrote %s (%d entries, %d bytes)", outPath, len(langMap), buf.Len())
	return nil
}
