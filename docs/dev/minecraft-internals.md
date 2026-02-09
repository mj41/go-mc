# Minecraft Java Edition Internals

## Unobfuscated Server Builds

Starting with 1.21.11, Mojang released separate unobfuscated server builds (titled `1.21.11_unobfuscated` in the launcher). This was an experimental step before removing obfuscation entirely in 26.1+.

- https://minecraft.wiki/w/Java_Edition_1.21.11
- https://minecraft.wiki/w/Tutorial:See_Minecraft%27s_code

### Decompiling with mcsrc

Use `mcsrc` (in the w42-mc-cubes repo) to download the unobfuscated JAR and decompile to Java source:

```bash
cd w42-mc-cubes/cmd-dev
go run ./mcsrc --version 1.21.11           # → ../temp/mc-src-1.21.11/
go run ./mcsrc --version 1.21.11 --jar-only # JAR only, no decompile
```

Output: `temp/mc-src-1.21.11/net/minecraft/` (~4600 .java files, fully readable names).

The `tools/` package also uses unobfuscated JARs for its Java extractors — see
[tools.md](tools.md) for the extraction pipeline.

## Packets

https://minecraft.wiki/w/Java_Edition_protocol/Packets#Level_Chunk_with_Light

Java Edition protocol for 1.21.11, protocol 774.
