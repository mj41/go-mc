// packetdiff compares the packet wire schemas of two extracted MC versions.
//
//	cd tools && go run ./packetdiff 1.21.11 26.2
//	cd tools && go run ./packetdiff ../temp/jsons/1.21.11 ../temp/jsons/26.2
//
// It reads packet_schema.json (from GenPacketSchema) and packets.json (numeric
// IDs per state) from each directory and prints: packets added / removed per
// state and flow, ID moves, and per-packet wire-layout changes as a token diff.
// A token is one codec or buffer read in wire order; nested codecs appear
// inline as Owner.FIELD{...}, so a change inside a shared type is reported on
// every packet that carries it.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type schemaEntry struct {
	State  string   `json:"state"`
	Class  string   `json:"class"`
	Tokens []string `json:"tokens"`
}

type schema struct {
	Packets map[string]schemaEntry `json:"packets"`
}

// packets.json: state -> flow -> name -> {protocol_id}
type packetsReport map[string]map[string]map[string]struct {
	ProtocolID int `json:"protocol_id"`
}

type version struct {
	name   string
	schema schema
	ids    packetsReport
}

func main() {
	var args []string
	verbose := false
	for _, s := range os.Args[1:] {
		if s == "-v" {
			verbose = true
		} else {
			args = append(args, s)
		}
	}
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: packetdiff [-v] <versionA|dirA> <versionB|dirB>   (-v lists renumbered packets)")
		os.Exit(2)
	}
	a, err := load(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, "packetdiff:", err)
		os.Exit(1)
	}
	b, err := load(args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "packetdiff:", err)
		os.Exit(1)
	}
	fmt.Printf("packetdiff: %s (%d packets) -> %s (%d packets)\n\n", a.name, len(a.schema.Packets), b.name, len(b.schema.Packets))
	idChanges(a, b, verbose)
	layoutChanges(a, b)
}

func load(arg string) (*version, error) {
	dir := arg
	if _, err := os.Stat(dir); err != nil {
		root, rerr := goMCRoot()
		if rerr != nil {
			return nil, rerr
		}
		dir = filepath.Join(root, "temp", "jsons", arg)
	}
	v := &version{name: filepath.Base(dir)}
	if err := readJSON(filepath.Join(dir, "packet_schema.json"), &v.schema); err != nil {
		return nil, err
	}
	if err := readJSON(filepath.Join(dir, "packets.json"), &v.ids); err != nil {
		fmt.Fprintf(os.Stderr, "packetdiff: %s: no packets.json, IDs skipped (%v)\n", v.name, err)
		v.ids = packetsReport{}
	}
	return v, nil
}

func goMCRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			if _, err := os.Stat(filepath.Join(dir, "tools")); err == nil {
				return dir, nil
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go-mc root not found from %s", dir)
		}
		dir = parent
	}
}

func readJSON(path string, v any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// idChanges prints packets added and removed per state/flow; renumbered
// packets (the usual consequence of an insertion) are counted unless verbose.
func idChanges(a, b *version, verbose bool) {
	states := union(keys(a.ids), keys(b.ids))
	printed := false
	for _, st := range states {
		for _, flow := range []string{"clientbound", "serverbound"} {
			am, bm := a.ids[st][flow], b.ids[st][flow]
			names := union(keys(am), keys(bm))
			var lines, moved []string
			for _, n := range names {
				ai, aok := am[n]
				bi, bok := bm[n]
				switch {
				case aok && !bok:
					lines = append(lines, fmt.Sprintf("  - %-45s (was 0x%02X)", n, ai.ProtocolID))
				case !aok && bok:
					lines = append(lines, fmt.Sprintf("  + %-45s 0x%02X", n, bi.ProtocolID))
				case ai.ProtocolID != bi.ProtocolID:
					moved = append(moved, fmt.Sprintf("  ~ %-45s 0x%02X -> 0x%02X", n, ai.ProtocolID, bi.ProtocolID))
				}
			}
			if len(lines) == 0 && len(moved) == 0 {
				continue
			}
			printed = true
			fmt.Printf("## %s %s: %d added/removed, %d renumbered\n", st, flow, len(lines), len(moved))
			for _, l := range lines {
				fmt.Println(l)
			}
			if verbose {
				for _, l := range moved {
					fmt.Println(l)
				}
			}
			fmt.Println()
		}
	}
	if !printed {
		fmt.Println("## packet IDs: no changes")
		fmt.Println()
	}
}

// layoutChanges prints, per packet present in both versions, the token diff.
func layoutChanges(a, b *version) {
	names := union(keys(a.schema.Packets), keys(b.schema.Packets))
	changed := 0
	for _, n := range names {
		ae, aok := a.schema.Packets[n]
		be, bok := b.schema.Packets[n]
		if !aok || !bok {
			continue
		}
		if strings.Join(ae.Tokens, "\x00") == strings.Join(be.Tokens, "\x00") {
			continue
		}
		changed++
		fmt.Printf("## %s (%s, %s)\n", n, be.State, shortClass(be.Class))
		for _, l := range diffTokens(flatten(ae.Tokens), flatten(be.Tokens)) {
			fmt.Println(l)
		}
		fmt.Println()
	}
	both := 0
	for n := range a.schema.Packets {
		if _, ok := b.schema.Packets[n]; ok {
			both++
		}
	}
	fmt.Printf("packetdiff: %d of %d shared packets changed their wire layout\n", changed, both)
}

// diffTokens is an LCS diff over token lists; nested-only changes show as
// one removed and one added token with the differing nested part.
func diffTokens(a, b []string) []string {
	n, m := len(a), len(b)
	lcs := make([][]int, n+1)
	for i := range lcs {
		lcs[i] = make([]int, m+1)
	}
	for i := n - 1; i >= 0; i-- {
		for j := m - 1; j >= 0; j-- {
			if a[i] == b[j] {
				lcs[i][j] = lcs[i+1][j+1] + 1
			} else if lcs[i+1][j] >= lcs[i][j+1] {
				lcs[i][j] = lcs[i+1][j]
			} else {
				lcs[i][j] = lcs[i][j+1]
			}
		}
	}
	var out []string
	i, j := 0, 0
	for i < n && j < m {
		switch {
		case a[i] == b[j]:
			out = append(out, "    "+clip(a[i]))
			i++
			j++
		case lcs[i+1][j] >= lcs[i][j+1]:
			out = append(out, "  - "+clip(a[i]))
			i++
		default:
			out = append(out, "  + "+clip(b[j]))
			j++
		}
	}
	for ; i < n; i++ {
		out = append(out, "  - "+clip(a[i]))
	}
	for ; j < m; j++ {
		out = append(out, "  + "+clip(b[j]))
	}
	return out
}

// flatten opens a single top-level "λX.<init>{a, b, c}" (packets read by a
// buffer constructor) into its members so the diff points at the field that
// changed instead of reporting the whole constructor.
func flatten(tokens []string) []string {
	if len(tokens) != 1 || !strings.HasPrefix(tokens[0], "λ") {
		return tokens
	}
	t := tokens[0]
	open := strings.Index(t, "{")
	if open < 0 || !strings.HasSuffix(t, "}") {
		return tokens
	}
	inner := t[open+1 : len(t)-1]
	var out []string
	depth, start := 0, 0
	for i := 0; i < len(inner); i++ {
		switch inner[i] {
		case '{':
			depth++
		case '}':
			depth--
		case ',':
			if depth == 0 {
				out = append(out, strings.TrimSpace(inner[start:i]))
				start = i + 1
			}
		}
	}
	if rest := strings.TrimSpace(inner[start:]); rest != "" {
		out = append(out, rest)
	}
	return out
}

func clip(s string) string {
	const max = 200
	if len(s) <= max {
		return s
	}
	return s[:max] + "…"
}

func shortClass(c string) string {
	return c[strings.LastIndex(c, ".")+1:]
}

func keys[V any](m map[string]V) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func union(a, b []string) []string {
	set := map[string]bool{}
	for _, s := range a {
		set[s] = true
	}
	for _, s := range b {
		set[s] = true
	}
	return keys(set)
}
