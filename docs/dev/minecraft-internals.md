# Minecraft Java Edition Internals

## Unobfuscated Server Builds

Starting with 1.21.11, Mojang released separate unobfuscated server builds (titled `1.21.11_unobfuscated` in the launcher). This was an experimental step before removing obfuscation entirely in 26.1+.

- https://minecraft.wiki/w/Java_Edition_1.21.11
- https://minecraft.wiki/w/Tutorial:See_Minecraft%27s_code

### Data Extraction

The `tools/` package downloads unobfuscated server JARs and runs Java
extractors against them to produce all data files used by go-mc.
See [tools.md](tools.md) for the full extraction and generation pipeline.

## Packets

https://minecraft.wiki/w/Java_Edition_protocol/Packets#Level_Chunk_with_Light

Java Edition protocol for 1.21.11, protocol 774.
