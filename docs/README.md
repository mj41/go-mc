# go-mc Documentation

Go libraries for Minecraft Java Edition — bot framework, server framework,
world I/O, NBT, and protocol implementation.

**Current version**: Minecraft 26.2 (protocol 776, data version 4903, Java 25)

## Data Coverage

All data is extracted directly from the Minecraft server jar via a container-based
pipeline (`tools/`). No third-party data sources (PrismarineJS, Burger, etc.) are
used. The extracted data is consumed by Go generators that produce type-safe Go code.

### MC 26.2 Data Summary

| Category | Count | Package | Source |
|----------|------:|---------|--------|
| Block types | 1,196 | `level/block` | `blocks.json` (MC `--all` report) |
| Block states | 32,366 | `level/block` | `blocks.json` → `block_states.nbt` |
| Block property enums | 29 | `level/block` | `block_properties.json` (Java extractor) |
| Block entity types | 49 | `level/block` | `block_entities.json` (Java extractor) |
| Packet IDs | 256 | `data/packetid` | `packets.json` (MC `--all` report) |
| Entity types | 158 | `data/entity` | `entities.json` (Java extractor) |
| Item types | 1,537 | `data/item` | `items.json` (Java extractor `GenItems`; the `--all` report was dropped in 26.x) |
| Sound IDs | 1,968 | `data/soundid` | `registries.json` (MC `--all` report) |
| Registries | 95 | `data/registryid` | `registries.json` (MC `--all` report) |
| Biomes | 66 | `level/biome` | `biomes.json` (Java extractor) |
| Data components | 111 | `level/component` | `components.json` (Java extractor) |
| Languages | 142 | `data/lang` | Mojang asset index CDN |

### Implementation Status

| Feature | Status | Notes |
|---------|--------|-------|
| Network protocol (bot) | Complete | Login, config, play phases. E2E tested against Paper 26.2 build 121 and vanilla 26.2 |
| Network protocol (server) | Complete | Per-registry RegistryData, all phases |
| Chunk I/O (network) | Complete | BitStorage format (1.21.5+ no length prefix) |
| Chunk I/O (save/region) | Complete | Anvil format, tested with MC 1.21.4 + 1.21.11; 26.1+ moved region dirs under `dimensions/<ns>/<dim>/` and `level.dat` fields (`difficulty_settings`, `spawn`) |
| Chat signing | Complete | HistoryUpdate checksum, globalIndex, PackedSignature |
| Slot / inventory | Complete | Post-1.20.5 format with component data |
| Data components | Complete | All 111 wire protocol types (IDs 0–110) |
| Block state mapping | Complete | 32,366 states, validated via E2E |
| NBT codec | Complete | Full spec, SNBT, RawMessage |
| RCON | Complete | Client and server |

### Data Accuracy

All generated data comes directly from Minecraft's own server jar — no third-party
data sources. The extraction + generation pipeline is fully repeatable:
`cd tools && go run . --extract --version 26.2` produces identical output
each run. See [dev/tools.md](dev/tools.md) for the full pipeline documentation.

## Packages

| Package | Description |
|---------|-------------|
| `bot/` | Bot (client) framework — connect, login, keepalive, events |
| `bot/basic/` | Core bot features — position, settings, keepalive, tags |
| `bot/msg/` | Chat message handling with signing support |
| `bot/screen/` | Inventory / slot management |
| `bot/world/` | Chunk tracking for connected bots |
| `server/` | Server framework — handshake, login, config, gameplay |
| `level/` | Chunk data structures, bit storage, palettes |
| `level/block/` | Block types, state IDs, property enums |
| `level/component/` | Data component types (111 types) |
| `level/biome/` | Biome list |
| `nbt/` | NBT codec (binary + SNBT) |
| `net/` | Low-level network (connection, encryption, compression) |
| `net/packet/` | Packet primitives (VarInt, String, etc.) |
| `chat/` | Chat message format (JSON + NBT) |
| `save/` | Anvil world I/O (region files, chunks) |
| `registry/` | Registry codec + network serialization |
| `data/packetid/` | Protocol packet ID constants |
| `data/entity/` | Entity type IDs |
| `data/item/` | Item type IDs |
| `data/soundid/` | Sound event IDs |
| `data/registryid/` | All 95 registry ID mappings |

## Developer Documentation

See [docs/dev/](dev/) for:

- [tools.md](dev/tools.md) — Code generation & MC data extraction pipeline
