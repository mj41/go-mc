# go-mc/tools — Code Generation & MC Data Extraction

Unified tool that extracts Minecraft data from an unobfuscated server jar and
generates all Go source files for go-mc.

## Quick Start

```bash
# Full pipeline: extract MC data + generate Go code (from go-mc root)
cd tools && go run . --extract --version 1.21.11

# Generate only (from previously extracted JSONs)
cd tools && go run . ../temp/jsons/1.21.11
```

## Requirements

- Go 1.22+
- `podman` or `docker` (for `--extract` mode)
- Internet access (first run downloads ~57 MB server jar; cached afterward)

## Usage

```
go run . [options] [json-dir]

Options:
  --extract           Run container extraction before generation
  --version VER       MC version to extract (required with --extract)
  --runtime RT        Container runtime: podman (default) or docker
  --dry-run           Print container command without running
  --gen-only          Extract only, skip generation

JSON dir can also be a positional arg: `go run . /path/to/jsons`
```

## What It Does

### Extraction (`--extract`)

Runs an `eclipse-temurin:21-jdk` container that:

1. Downloads the unobfuscated MC server jar from Minecraft Wiki API
2. Extracts the inner server jar from the bundler format
3. Runs the MC `--all` data generator → 6000+ data files
4. Copies key reports: `blocks.json`, `packets.json`, `registries.json`,
   `items.json`, `commands.json`, `datapack.json`
5. Compiles and runs 4 custom Java extractors:
   - **GenEntities** — entity types with dimensions (id, name, width, height)
   - **GenComponents** — data component types with networkability flags
   - **GenBlockEntities** — block entity types with valid blocks
   - **GenBlockProperties** — block state property definitions (boolean/integer/enum)

Output: `temp/jsons/<version>/*.json` (~8 MB total)

### Generation

Runs 8 Go generators that read the extracted JSON files and produce
Go source code:

| Generator | Input | Output |
|-----------|-------|--------|
| packetid | `packets.json` | `data/packetid/packetid.go` |
| soundid | `registries.json` | `data/soundid/soundid.go` |
| item | `items.json` + `registries.json` | `data/item/item.go` |
| blocks | `blocks.json` + `block_properties.json` | `level/block/blocks.go` + `block_states.nbt` + `properties_enum.go` |
| entity | `entities.json` | `data/entity/entity.go` |
| component | `components.json` | `level/component/components.go` |
| blockentities | `block_entities.json` | `level/block/blockentity.go` + `blockentities.go` |
| registryid | `registries.json` | `data/registryid/*.go` (95 files) |

## Directory Layout

```
tools/
├── main.go              # unified entry point
├── extract.go           # container extraction logic
├── helpers.go           # shared utilities
├── gen_packetid.go      # generator: packet IDs
├── gen_soundid.go       # generator: sound IDs
├── gen_item.go          # generator: items
├── gen_blocks.go        # generator: blocks + states + properties
├── gen_entity.go        # generator: entities
├── gen_component.go     # generator: data components
├── gen_blockentities.go # generator: block entities
├── gen_registryid.go    # generator: registry IDs (95 registries)
├── java/                # Java extractor sources (committed)
│   ├── ExtractAll.java  # container orchestrator
│   ├── GenEntities.java
│   ├── GenComponents.java
│   ├── GenBlockEntities.java
│   └── GenBlockProperties.java
└── go.mod               # separate module (replace → parent go-mc)
```

Extracted data lives at the repo root (gitignored):

```
temp/                        # gitignored working data
├── cache/                   # cached server jars
└── jsons/                   # extracted JSON files per version
    └── 1.21.11/
        ├── blocks.json
        ├── packets.json
        ├── registries.json
        ├── items.json
        ├── entities.json
        ├── components.json
        ├── block_entities.json
        ├── block_properties.json
        └── ...
```

## Updating to a New MC Version

```bash
cd tools

# 1. Extract + generate in one shot
go run . --extract --version 1.XX.YY

# 2. Verify from go-mc root
cd .. && go build ./... && go test ./...

# 3. Review generated diffs
git diff --stat
```

## Module Structure

`tools/` is a separate Go module (`github.com/Tnze/go-mc/tools`) with a
`replace` directive pointing to the parent go-mc module. This keeps the
generation tooling out of the main module's dependency graph.

```go
// tools/go.mod
module github.com/Tnze/go-mc/tools
go 1.22
require github.com/Tnze/go-mc v1.22.0
replace github.com/Tnze/go-mc => ../
```

## Adding a New Generator

1. Create `gen_foo.go` with `func genFoo(jsonDir, goMCRoot string) error`
2. Add `{"foo", genFoo}` to the `generators` slice in `main.go`
3. Use the shared helpers from `helpers.go` (`readJSON`, `writeFile`,
   `snakeToCamel`, etc.)

## Adding a New Java Extractor

1. Create `java/GenFoo.java` — a single-file Java 21 program that reads
   from the MC server jar classpath and writes JSON to `/jsons/<version>/`
2. Add `"GenFoo"` to the `extractors` array in `java/ExtractAll.java`
3. Create the corresponding Go generator to consume the JSON output
