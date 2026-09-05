# go-mc/tools — Code Generation & MC Data Extraction

Unified tool that extracts Minecraft data from the Mojang server jar and
generates all Go source files for go-mc.

## Quick Start

```bash
# Full pipeline: extract MC data (if needed) + generate Go code
cd tools && go run . --version 26.2

# Force re-extraction even if cached JSONs exist
cd tools && go run . --version 26.2 --extract

# Generate from a specific JSON directory
cd tools && go run . ../temp/jsons/26.2
```

## Requirements

- Go 1.22+
- `podman` or `docker` (for `--extract` mode; runs `eclipse-temurin:25-jdk` — MC 26.x needs Java 25,
  and JDK 25 still runs the 1.21.x jars)
- Internet access (first run downloads ~60 MB server jar; cached afterward)

## Usage

```
go run . [options] [json-dir]

Options:
  --version VER       MC version (extracts if cached JSONs missing)
  --extract           Force re-extraction even if cache exists
  --runtime RT        Container runtime: podman (default) or docker
  --dry-run           Print container command without running
  --gen-only          Extract only, skip generation

JSON dir can also be a positional arg: `go run . /path/to/jsons`
```

## What It Does

### Extraction (`--extract`)

Runs as a 3-phase pipeline:

**Phase 1 — Go host downloads (needs internet):**
1. Downloads the MC server jar (from `unobfuscated_versions.json` for
   pre-26.x, or from the Mojang version manifest for 26.x+)
2. Downloads all ~147 language files from the Mojang asset index CDN

**Phase 2 — Container extraction (no internet needed):**

Runs an `eclipse-temurin:25-jdk` container that:
1. Extracts the inner server jar from the bundler format
2. Runs the MC `--all` data generator → reports/*.json
3. Copies key reports: `blocks.json`, `packets.json`, `registries.json`,
   `commands.json`, `datapack.json` (and `items.json` up to 1.21.11 — 26.x no
   longer emits it, see `GenItems` below)
4. Compiles and runs 8 custom Java extractors:
   - **GenEntities** — entity types with dimensions (id, name, width, height)
   - **GenComponents** — data component types with networkability flags
   - **GenComponentSchema** — component wire format schema via reflection
   - **GenBlockEntities** — block entity types with valid blocks
   - **GenBlockProperties** — block state property definitions (boolean/integer/enum)
   - **GenBiomes** — biome protocol ordering via runtime registry introspection
   - **GenItems** — per-item `max_stack_size` / `item_name` in the shape of the old `items.json` report
   - **GenPacketSchema** — wire layout of every packet (and of shared wire structures such as
     `LevelChunkSection`, `ItemStack`, `ChatType$Bound`) read from bytecode with `java.lang.classfile`;
     `packet_schema.json` feeds `packetdiff`, not a generator

Output: `temp/jsons/<version>/*.json` + `lang/` (~83 MB total)

### Generation

Runs 10 Go generators that read the extracted JSON files and produce
Go source code:

| Generator | Input | Output |
|-----------|-------|--------|
| packetid | `packets.json` | `data/packetid/packetid.go` |
| soundid | `registries.json` | `data/soundid/soundid.go` |
| item | `items.json` (GenItems) + `registries.json` | `data/item/item.go` |
| blocks | `blocks.json` + `block_properties.json` | `level/block/blocks.go` + `block_states.nbt` + `properties_enum.go` |
| entity | `entities.json` | `data/entity/entity.go` |
| component | `components.json` + `component_schema.json` | `level/component/components.go` + `*_gen.go` |
| blockentities | `block_entities.json` | `level/block/blockentity.go` + `blockentities.go` |
| registryid | `registries.json` | `data/registryid/*.go` (95 files) |
| biome | `biomes.json` | `level/biome/list.go` |
| lang | `lang/*.json` | `data/lang/<locale>/<locale>.go` (147 languages) |

## Directory Layout

```
tools/
├── main.go              # unified entry point
├── download.go          # Go-side downloads (server jar + languages)
├── extract.go           # 3-phase extraction pipeline orchestration
├── helpers.go           # shared utilities
├── gen_packetid.go      # generator: packet IDs
├── gen_soundid.go       # generator: sound IDs
├── gen_item.go          # generator: items
├── gen_blocks.go        # generator: blocks + states + properties
├── gen_entity.go        # generator: entities
├── gen_component.go     # generator: data components
├── gen_blockentities.go # generator: block entities
├── gen_registryid.go    # generator: registry IDs (95 registries)
├── gen_biome.go         # generator: biomes
├── gen_lang.go          # generator: language translations (147 languages)
├── gen_component_types.go # generator: component type structs (*_gen.go)
├── hand-crafted/        # manually maintained config files
│   ├── component_schema.json
│   ├── naming_overrides.json
│   └── packet_phases.json
├── unobfuscated_versions.json  # hardcoded unobfuscated jar URLs (pre-26.x)
├── java/                # Java extractor sources (committed)
│   ├── ExtractAll.java  # container extraction (no internet needed)
│   ├── GenBiomes.java
│   ├── GenBlockEntities.java
│   ├── GenBlockProperties.java
│   ├── GenComponents.java
│   ├── GenComponentSchema.java
│   └── GenEntities.java
└── go.mod               # separate module (replace → parent go-mc)
```

Extracted data lives at the repo root (gitignored):

```
temp/                        # gitignored working data
├── cache/                   # cached server jars + libs
└── jsons/                   # extracted JSON files per version
    └── 26.2/
        ├── blocks.json
        ├── packets.json
        ├── registries.json
        ├── items.json
        ├── entities.json
        ├── components.json
        ├── block_entities.json
        ├── block_properties.json
        ├── component_schema.json
        ├── biomes.json
        ├── lang/                # 147 language JSON files
        │   ├── en_us.json
        │   ├── cs_cz.json
        │   └── ...
        └── ...
```

## Updating to a New MC Version

```bash
cd tools

# 0. Preview the scope without Java: registries, block states, item components
go run ./mcmeta versions            # what exists (releases and snapshots)
go run ./mcmeta diff 26.2 26.X

# 1. Extract (if needed) + generate
go run . --version 26.X
go run ./mcmeta check 26.X          # our registries.json == mcmeta's (sanity)

# 2. Wire changes: what the hand-written packet code in bot/ and server/ must follow
go run ./packetdiff 26.2 26.X        # -v also lists renumbered packets

# 3. Verify from go-mc root
cd .. && go build ./... && go test ./...

# 4. Review generated diffs
git diff --stat
```

### Registry preview (`mcmeta`)

`mcmeta` reads [misode/mcmeta](https://github.com/misode/mcmeta), the archive of the
game's data-generator output for every release and snapshot (tagged `<version>-summary`),
cached under `temp/mcmeta/<version>/`. `versions` lists ids with data/protocol versions;
`diff A B` prints per-registry entry counts with added/removed names, blocks whose state
properties changed, and items whose default components changed (with a per-component
tally); `check V` compares `temp/jsons/V/registries.json` with mcmeta's registries and
exits non-zero on a difference. Example, 26.2 → 26.3-pre-2: block 1196 → 1286, item
1537 → 1658, data_component_type 111 → 122 (+13 −2), `minecraft:attack_animation` added
on every item — known before the first extraction run.

### Packet wire diff (`packetdiff`)

`GenPacketSchema` writes `packet_schema.json`: for every packet (keyed `<flow>/<name>`,
with the state of the `*PacketTypes` class that declares it) the codec or buffer reads in
wire order — `ByteBufCodecs.VAR_INT`, `buf.readUUID`, nested codecs inline as
`Owner.FIELD{...}`, constructor readers as `λClass.<init>{...}`. Shared structures that
packets carry as byte blobs (`LevelChunkSection.read/write`, `ItemStack.STREAM_CODEC`,
`ChatType$Bound.STREAM_CODEC`, …) are listed under `struct/`. `packetdiff A B` joins two
such files with their `packets.json` and prints packets added/removed, renumbered counts,
and a token diff for every packet whose layout changed. Java type renames show up as
changed tokens with the same wire shape (e.g. `buf.readLpVec3` → `Vec3.LP_STREAM_CODEC{…}`);
read the tokens, they carry the primitive reads. For 1.21.11 → 26.2 it reports the login
`readBoolean`, the login-finished `UUIDUtil.STREAM_CODEC`, the new `interact` layout, the
chunk section `readShort` (fluid count), `set_time`'s clock structure and the
`set_player_team` parameter reorder — everything the live tests had to find by hand.

## See also

- [mcsrc.dev](https://mcsrc.dev) — Fabric's in-browser decompiled Minecraft source (Vineflower in
  WebAssembly); quickest way to look at one class of an unobfuscated 26.x jar.
- A local [Vineflower](https://vineflower.org) run over the client or server jar when a grep-able
  full source tree is needed (1.12.0 handles the 26.x class files on JDK 25).
- [misode/mcmeta](https://github.com/misode/mcmeta) — processed data-generator reports for every
  release and snapshot on tagged branches (`<version>-summary`), a Java-free way to preview what a
  new version changes in the registries.

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

## Unobfuscated Versions (Pre-26.x)

Versions before 26.1 ship obfuscated by default. Separate unobfuscated
builds were released for versions 25w45a through 1.21.11. These are not
listed in the standard Mojang version manifest, so their server jar URLs
are hardcoded in `tools/unobfuscated_versions.json`.

Versions 26.1+ are natively unobfuscated and use the standard manifest
automatically.

### Adding a new unobfuscated version entry

1. Open the Minecraft Wiki page for the version (e.g.,
   `https://minecraft.wiki/w/Java_Edition_1.21.11`)
2. Find the wikitext infobox via the wiki API:
   ```
   curl -s 'https://minecraft.wiki/api.php?action=parse&page=Java+Edition+VERSION&prop=wikitext&format=json&section=0' | jq -r '.parse.wikitext["*"]'
   ```
3. Look for the `serverdl` line with `{{dl|HASH|server|title=Unobfuscated}}`
4. The download URL is: `https://piston-data.mojang.com/v1/objects/HASH/server.jar`
5. Add an entry to `tools/unobfuscated_versions.json`:
   ```json
   {
     "VERSION": {
       "server_sha1": "HASH",
       "server_url": "https://piston-data.mojang.com/v1/objects/HASH/server.jar"
     }
   }
   ```

## Adding a New Java Extractor

1. Create `java/GenFoo.java` — a single-file Java 21 program that reads
   from the MC server jar classpath and writes JSON to the output directory
2. Add `"GenFoo"` to the `extractors` array in `java/ExtractAll.java`
3. Create the corresponding Go generator to consume the JSON output
