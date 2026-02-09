# Developer Documentation

Internal documentation for go-mc contributors and maintainers.

## Contents

- [tools.md](tools.md) — Code generation & MC data extraction pipeline
- [minecraft-internals.md](minecraft-internals.md) — MC server internals, unobfuscated builds

## Overview

go-mc uses a container-based pipeline to extract data from Minecraft server jars
and generate Go source code. All generated files carry a header comment identifying
the generator and data source.

### Data Flow

```
MC server jar (unobfuscated)
  ↓ container: eclipse-temurin:21-jdk
  ├── MC --all data generator → reports/ (blocks, packets, registries, items, ...)
  └── Java extractors → entities, components, block_entities, block_properties
  ↓
temp/jsons/<version>/*.json  (~8 MB, gitignored)
  ↓ go run ./tools/ ...
  ├── data/packetid/packetid.go          (262 packet IDs)
  ├── data/soundid/soundid.go            (1,838 sounds)
  ├── data/item/item.go                  (1,505 items)
  ├── data/entity/entity.go              (157 entities)
  ├── data/registryid/*.go               (95 registries, 6,726 entries)
  ├── level/block/blocks.go              (1,166 block types)
  ├── level/block/block_states.nbt       (29,671 states)
  ├── level/block/properties_enum.go     (29 enum types)
  ├── level/block/blockentity.go         (48 block entity type constants)
  ├── level/block/blockentities.go       (block entity → block mappings)
  └── level/component/components.go      (104 data component types)
```

### Updating to a New MC Version

```bash
cd tools && go run . --extract --version X.YY.ZZ
cd .. && go build ./... && go test ./...
git diff --stat
```

See [tools.md](tools.md) for the full pipeline documentation.

### What's Generated vs Hand-Written

**Generated** (do not edit manually):
- All files in `data/packetid/`, `data/soundid/`, `data/item/`, `data/entity/`
- `data/registryid/*.go` (95 files)
- `level/block/blocks.go`, `block_states.nbt`, `properties_enum.go`
- `level/block/blockentity.go` (constants), `blockentities.go` (mappings)
- `level/component/components.go` (component constructor switch)

**Hand-written** (not generated):
- `level/biome/list.go` — biome list (65 entries)
- `level/component/*.go` — individual component type implementations
- `bot/`, `server/`, `net/`, `nbt/`, `chat/`, `save/` — all framework code
- `data/packetid/legacy.go` — backward-compat aliases
- `data/lang/gen_lang.go` — language file generator (fetches from Mojang API)
