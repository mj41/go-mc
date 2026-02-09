// tools generates all Go source files for go-mc from MC data JSONs.
//
// Usage:
//
//	go run ./tools/ <json-dir>                       # generate only
//	go run ./tools/ --extract --version 1.21.11      # extract + generate
//	go run ./tools/ --extract --version 1.21.11 --dry-run
//
// With --extract, the tool runs a Docker/Podman container to download the
// MC server jar, run its data generator, and run Java extractors. The output
// JSONs are placed in temp/jsons/<version>/ and then used for generation.
//
// The go-mc root is auto-detected by walking up from the working directory.
// Output files are written to the standard locations within go-mc.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type generator struct {
	name string
	fn   func(jsonDir, goMCRoot string) error
}

var generators = []generator{
	{"packetid", genPacketID},
	{"soundid", genSoundID},
	{"item", genItem},
	{"blocks", genBlocks},
	{"entity", genEntity},
	{"component", genComponent},
	{"blockentities", genBlockEntities},
	{"registryid", genRegistryID},
	{"biome", genBiome},
}

func main() {
	args := os.Args[1:]

	extract := flagBool(args, "--extract")
	version := flagValue(args, "--version")
	runtime := flagValue(args, "--runtime")
	dryRun := flagBool(args, "--dry-run")
	genOnly := flagBool(args, "--gen-only")

	goMCRoot, err := detectGoMCRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "tools: %v\n", err)
		os.Exit(1)
	}

	var jsonDir string

	if extract {
		if version == "" {
			fmt.Fprintln(os.Stderr, "tools: --extract requires --version")
			printUsage()
			os.Exit(1)
		}
		if runtime == "" {
			runtime = detectRuntime()
			if runtime == "" {
				fmt.Fprintln(os.Stderr, "tools: neither podman nor docker found in PATH")
				os.Exit(1)
			}
		}

		dir, err := runExtract(goMCRoot, version, runtime, dryRun)
		if err != nil {
			fmt.Fprintf(os.Stderr, "tools: extract: %v\n", err)
			os.Exit(1)
		}
		jsonDir = dir

		if dryRun || genOnly {
			return
		}
	} else {
		// Positional arg.
		jsonDir = positionalArg(args)
		if jsonDir == "" {
			printUsage()
			os.Exit(1)
		}
	}

	if _, err := os.Stat(jsonDir); err != nil {
		fmt.Fprintf(os.Stderr, "tools: JSON directory not found: %s\n", jsonDir)
		os.Exit(1)
	}

	jsonVersion = filepath.Base(jsonDir)
	fmt.Fprintf(os.Stderr, "tools: go-mc=%s json-dir=%s\n\n", goMCRoot, jsonDir)

	for _, g := range generators {
		start := time.Now()
		if err := g.fn(jsonDir, goMCRoot); err != nil {
			fmt.Fprintf(os.Stderr, "tools: %s: %v\n", g.name, err)
			os.Exit(1)
		}
		elapsed := time.Since(start).Round(time.Millisecond)
		fmt.Fprintf(os.Stderr, "  %s: done (%s)\n", g.name, elapsed)
	}
	fmt.Fprintf(os.Stderr, "\ntools: all %d generators done!\n", len(generators))
}

// detectGoMCRoot finds the go-mc repository root by walking up from the
// working directory looking for go.mod with the go-mc module path.
func detectGoMCRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("getting working directory: %v", err)
	}
	for d := wd; d != "/" && d != "."; d = filepath.Dir(d) {
		data, err := os.ReadFile(filepath.Join(d, "go.mod"))
		if err != nil {
			continue
		}
		if strings.Contains(string(data), "module github.com/Tnze/go-mc\n") {
			return d, nil
		}
	}
	return "", fmt.Errorf("cannot find go-mc root from %s (need go.mod with module github.com/Tnze/go-mc)", wd)
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr, "  go run ./tools/ <json-dir>                          # generate from existing JSONs")
	fmt.Fprintln(os.Stderr, "  go run ./tools/ --extract --version 1.21.11          # extract + generate")
	fmt.Fprintln(os.Stderr, "  go run ./tools/ --extract --version 1.21.11 --dry-run")
}

// positionalArg returns the first arg that doesn't look like a flag.
func positionalArg(args []string) string {
	skip := false
	for _, a := range args {
		if skip {
			skip = false
			continue
		}
		if a == "--version" || a == "--runtime" {
			skip = true
			continue
		}
		if strings.HasPrefix(a, "--") {
			continue
		}
		return a
	}
	return ""
}
