/**
 * GenPacketSchema — Extracts the wire layout of every network packet from the
 * unobfuscated server jar by reading bytecode (java.lang.classfile, JDK 24+).
 *
 * Output: packet_schema.json in the current directory:
 *   { "packets": {
 *       "<flow>/<name>": {            // e.g. "serverbound/minecraft:interact"
 *         "state":  "play",           // from the *PacketTypes class that declares it
 *         "class":  "net.minecraft.network.protocol.game.ServerboundInteractPacket",
 *         "tokens": ["ByteBufCodecs.VAR_INT", "InteractionHand.STREAM_CODEC{...}", ...]
 *   } } }
 *
 * Tokens are the codec/read calls in wire order. Two packet styles are handled:
 *   - STREAM_CODEC built in <clinit> (StreamCodec.composite & friends): the
 *     GETSTATIC/INVOKE chain that builds the field, split per PUTSTATIC.
 *   - Packet.codec(write, read) / buffer constructors: the read side is walked
 *     (readVarInt, readUUID, nested X.read(buf) / new X(buf), codec.decode).
 * Codec fields of other MC classes are expanded inline as `Owner.FIELD{...}` up
 * to MAX_DEPTH, so a change inside a shared type (Vec3, ChatType$Bound, ItemStack)
 * shows on every packet that carries it. Anything unrecognised becomes
 * `opaque:<owner>.<method>` so a diff still flags it.
 *
 * Numeric IDs and per-state membership come from packets.json (joined by
 * tools/packetdiff on the Go side); this file only needs the jar.
 */
import java.lang.classfile.*;
import java.lang.classfile.instruction.*;
import java.lang.constant.*;
import java.lang.reflect.*;
import java.nio.charset.StandardCharsets;
import java.nio.file.*;
import java.util.*;

public class GenPacketSchema {
    static final int MAX_DEPTH = 4;
    static final String[][] TYPES_CLASSES = {
        {"net.minecraft.network.protocol.handshake.HandshakePacketTypes", "handshake"},
        {"net.minecraft.network.protocol.status.StatusPacketTypes", "status"},
        {"net.minecraft.network.protocol.login.LoginPacketTypes", "login"},
        {"net.minecraft.network.protocol.configuration.ConfigurationPacketTypes", "configuration"},
        {"net.minecraft.network.protocol.game.GamePacketTypes", "play"},
        {"net.minecraft.network.protocol.common.CommonPacketTypes", "common"},
        {"net.minecraft.network.protocol.cookie.CookiePacketTypes", "cookie"},
        {"net.minecraft.network.protocol.ping.PingPacketTypes", "ping"},
    };

    /**
     * Wire structures that packets carry as opaque byte blobs or that are shared by
     * many packets; walked like a packet and emitted as "struct/<Class>.<member>" so
     * a change inside them is reported even when no packet's own layout moved.
     * {internal class name, member}: a static codec field or a read/write method.
     */
    static final String[][] STRUCTS = {
        {"net/minecraft/world/level/chunk/LevelChunkSection", "read"},
        {"net/minecraft/world/level/chunk/LevelChunkSection", "write"},
        {"net/minecraft/world/level/chunk/PalettedContainer", "read"},
        {"net/minecraft/network/protocol/game/ClientboundLevelChunkPacketData", "<init>"},
        {"net/minecraft/network/protocol/game/ClientboundLightUpdatePacketData", "<init>"},
        {"net/minecraft/network/protocol/game/CommonPlayerSpawnInfo", "<init>"},
        {"net/minecraft/network/protocol/game/ClientboundPlayerInfoUpdatePacket$Entry", "<init>"},
        {"net/minecraft/world/item/ItemStack", "STREAM_CODEC"},
        {"net/minecraft/world/item/ItemStack", "OPTIONAL_STREAM_CODEC"},
        {"net/minecraft/core/component/DataComponentPatch", "STREAM_CODEC"},
        {"net/minecraft/network/chat/ComponentSerialization", "STREAM_CODEC"},
        {"net/minecraft/network/chat/ChatType$Bound", "STREAM_CODEC"},
        {"net/minecraft/network/chat/RemoteChatSession$Data", "<init>"},
        {"net/minecraft/network/chat/SignedMessageBody$Packed", "<init>"},
        {"net/minecraft/network/chat/LastSeenMessages$Update", "<init>"},
        {"net/minecraft/core/RegistrySynchronization$PackedRegistryEntry", "STREAM_CODEC"},
        {"net/minecraft/world/entity/EntityType", "STREAM_CODEC"},
    };

    record Entry(String flow, String name, String state, String className, List<String> tokens) {}

    static final Map<String, ClassModel> classCache = new HashMap<>();
    static final Map<String, Map<String, List<String>>> clinitCache = new HashMap<>();

    public static void main(String[] args) throws Exception {
        List<Entry> entries = new ArrayList<>();
        int missing = 0;
        for (String[] tc : TYPES_CLASSES) {
            Class<?> types;
            try {
                types = Class.forName(tc[0]);
            } catch (ClassNotFoundException e) {
                System.err.println("GenPacketSchema: no " + tc[0] + " (skipped)");
                continue;
            }
            for (Field f : types.getDeclaredFields()) {
                if (!Modifier.isStatic(f.getModifiers())) continue;
                if (!f.getType().getName().equals("net.minecraft.network.protocol.PacketType")) continue;
                Object pt = f.get(null);
                String flow = String.valueOf(pt.getClass().getMethod("flow").invoke(pt)).toLowerCase(Locale.ROOT);
                String name = String.valueOf(pt.getClass().getMethod("id").invoke(pt));
                Class<?> packetClass = packetClassOf(f);
                List<String> tokens;
                if (packetClass == null) {
                    tokens = List.of("opaque:no-packet-class");
                    missing++;
                } else {
                    tokens = packetTokens(packetClass.getName().replace('.', '/'));
                }
                entries.add(new Entry(flow, name, tc[1], packetClass == null ? null : packetClass.getName(), tokens));
            }
        }
        for (String[] st : STRUCTS) {
            String owner = st[0], member = st[1];
            String label = shortOwner(owner) + "." + member;
            if (classModel(owner) == null) {
                System.err.println("GenPacketSchema: no class " + owner + " (struct skipped)");
                continue;
            }
            List<String> tokens;
            if (member.equals("<init>") || Character.isLowerCase(member.charAt(0))) {
                tokens = walkMethodByName(owner, member);
            } else {
                List<String> raw = clinitFields(owner).get(member);
                tokens = raw == null ? List.of("opaque:no-field") : expand(raw, owner, 0, new HashSet<>());
            }
            entries.add(new Entry("struct", label, "struct", owner.replace('/', '.'), tokens));
        }
        entries.sort(Comparator.comparing((Entry e) -> e.state).thenComparing(e -> e.flow).thenComparing(e -> e.name));

        StringBuilder sb = new StringBuilder("{\n  \"packets\": {\n");
        for (int i = 0; i < entries.size(); i++) {
            Entry e = entries.get(i);
            sb.append("    ").append(json(e.flow + "/" + e.name)).append(": {\n");
            sb.append("      \"state\": ").append(json(e.state)).append(",\n");
            sb.append("      \"class\": ").append(e.className == null ? "null" : json(e.className)).append(",\n");
            sb.append("      \"tokens\": [");
            for (int j = 0; j < e.tokens.size(); j++) {
                if (j > 0) sb.append(", ");
                sb.append(json(e.tokens.get(j)));
            }
            sb.append("]\n    }").append(i + 1 < entries.size() ? ",\n" : "\n");
        }
        sb.append("  }\n}\n");
        Files.writeString(Path.of("packet_schema.json"), sb.toString(), StandardCharsets.UTF_8);
        System.err.println("GenPacketSchema: " + entries.size() + " packets written to packet_schema.json"
            + (missing > 0 ? " (" + missing + " without a packet class)" : ""));
    }

    static Class<?> packetClassOf(Field f) {
        Type t = f.getGenericType();
        if (t instanceof ParameterizedType p && p.getActualTypeArguments().length == 1) {
            Type a = p.getActualTypeArguments()[0];
            if (a instanceof Class<?> c) return c;
            if (a instanceof ParameterizedType ap && ap.getRawType() instanceof Class<?> c) return c;
            if (a instanceof WildcardType w && w.getUpperBounds().length == 1 && w.getUpperBounds()[0] instanceof Class<?> c) return c;
        }
        return null;
    }

    // ---- per-packet entry point ------------------------------------------------

    static List<String> packetTokens(String internalName) {
        Map<String, List<String>> fields = clinitFields(internalName);
        List<String> codec = fields.get("STREAM_CODEC");
        if (codec == null) {
            // Some packets expose the codec under another name; take the first StreamCodec-typed field.
            for (Map.Entry<String, List<String>> e : fields.entrySet()) {
                if (e.getKey().endsWith("STREAM_CODEC") || e.getKey().endsWith("CODEC")) { codec = e.getValue(); break; }
            }
        }
        if (codec == null) return List.of("opaque:no-stream-codec");
        return expand(codec, internalName, 0, new HashSet<>());
    }

    // ---- <clinit> segmentation ---------------------------------------------------

    static ClassModel classModel(String internalName) {
        return classCache.computeIfAbsent(internalName, n -> {
            try (var in = GenPacketSchema.class.getClassLoader().getResourceAsStream(n + ".class")) {
                if (in == null) return null;
                return ClassFile.of().parse(in.readAllBytes());
            } catch (Exception e) {
                return null;
            }
        });
    }

    /** Raw tokens per static field, from the class initializer, keyed by field name. */
    static Map<String, List<String>> clinitFields(String internalName) {
        return clinitCache.computeIfAbsent(internalName, n -> {
            Map<String, List<String>> out = new LinkedHashMap<>();
            ClassModel cm = classModel(n);
            if (cm == null) return out;
            for (MethodModel m : cm.methods()) {
                if (!m.methodName().stringValue().equals("<clinit>")) continue;
                List<String> cur = new ArrayList<>();
                for (CodeElement el : m.code().map(c -> (Iterable<CodeElement>) c).orElse(List.of())) {
                    if (el instanceof FieldInstruction fi) {
                        if (fi.opcode() == Opcode.PUTSTATIC && fi.owner().asInternalName().equals(n)) {
                            out.put(fi.name().stringValue(), cur);
                            cur = new ArrayList<>();
                        } else if (fi.opcode() == Opcode.GETSTATIC && isCodecField(fi)) {
                            cur.add(shortOwner(fi.owner().asInternalName()) + "." + fi.name().stringValue());
                        }
                    } else if (el instanceof InvokeInstruction ii) {
                        String tok = invokeToken(ii);
                        if (tok != null) cur.add(tok);
                    } else if (el instanceof InvokeDynamicInstruction idi) {
                        cur.add(lambdaToken(idi));
                    }
                }
            }
            return out;
        });
    }

    static boolean isCodecField(FieldInstruction fi) {
        String desc = fi.typeSymbol().descriptorString();
        return desc.contains("Codec") || fi.owner().asInternalName().endsWith("ByteBufCodecs");
    }

    /** Tokens for codec-building calls; null for noise (accessor method refs, boxing, ...). */
    static String invokeToken(InvokeInstruction ii) {
        String owner = ii.owner().asInternalName();
        String name = ii.name().stringValue();
        String ret = ii.typeSymbol().returnType().descriptorString();
        String ownerShort = shortOwner(owner);
        if (owner.endsWith("codec/StreamCodec") || owner.endsWith("ByteBufCodecs") || ownerShort.equals("Packet")
                || ret.contains("StreamCodec") || (ret.contains("Codec") && owner.startsWith("net/minecraft/"))) {
            return ownerShort + "." + name;
        }
        return null;
    }

    static String lambdaToken(InvokeDynamicInstruction idi) {
        for (ConstantDesc cd : idi.bootstrapArgs()) {
            if (cd instanceof DirectMethodHandleDesc dmh) {
                return "λ" + shortOwner(dmh.owner().descriptorString().replaceAll("^L|;$", "")) + "." + dmh.methodName()
                    + "|" + dmh.owner().descriptorString().replaceAll("^L|;$", "") + "|" + dmh.methodName() + "|" + dmh.lookupDescriptor();
            }
        }
        return "λ?";
    }

    // ---- expansion --------------------------------------------------------------

    /**
     * Expands raw tokens: codec fields of other MC classes → "Owner.FIELD{...}",
     * Packet.codec(write, read) → the read lambda's method walk, lambda tokens → walked
     * when they point at an MC method taking a FriendlyByteBuf.
     */
    /** Combinators taking (encoder, decoder) lambdas: only the decoder carries wire order. */
    static final Set<String> ENCODE_DECODE = Set.of("Packet.codec", "StreamCodec.of", "StreamCodec.ofMember");
    /** Structural markers that carry no wire information of their own. */
    static final Set<String> MARKERS = Set.of("StreamCodec.composite", "StreamCodec.of", "StreamCodec.ofMember",
        "Packet.codec", "StreamCodec.cast");

    static List<String> expand(List<String> raw, String self, int depth, Set<String> visiting) {
        List<String> out = new ArrayList<>();
        // X.of(encode, decode) / Packet.codec(encode, decode): keep only the decode side.
        for (int i = 0; i < raw.size(); i++) {
            if (!ENCODE_DECODE.contains(raw.get(i))) continue;
            List<Integer> lambdas = new ArrayList<>();
            for (int j = 0; j < i; j++) if (raw.get(j).startsWith("λ")) lambdas.add(j);
            if (lambdas.size() >= 2) {
                int decode = lambdas.get(lambdas.size() - 1);
                List<String> pruned = new ArrayList<>();
                for (int j = 0; j < raw.size(); j++) {
                    if (j < i && raw.get(j).startsWith("λ") && j != decode) continue;
                    pruned.add(raw.get(j));
                }
                raw = pruned;
                break;
            }
        }
        for (String t : raw) {
            if (MARKERS.contains(t)) continue;
            out.addAll(expandOne(t, self, depth, visiting));
        }
        return out;
    }

    static List<String> expandOne(String t, String self, int depth, Set<String> visiting) {
        if (t.startsWith("λ")) {
            String[] parts = t.split("\\|", 4);
            String implOwner = parts.length > 1 ? parts[1] : "";
            String implName = parts.length > 2 ? parts[2] : "";
            String implDesc = parts.length > 3 ? parts[3] : "";
            // accessor / constructor refs used by composite() are noise; buffer readers are walked
            if (implDesc.contains("FriendlyByteBuf") || implDesc.contains("io/netty/buffer/ByteBuf")) {
                return List.of(walkLambda(t, depth, visiting));
            }
            return List.of();
        }
        int dot = t.lastIndexOf('.');
        if (dot > 0 && !t.contains("{")) {
            String ownerShort = t.substring(0, dot);
            String field = t.substring(dot + 1);
            String owner = resolveOwner(ownerShort, self);
            if (owner != null && !owner.endsWith("ByteBufCodecs") && Character.isUpperCase(field.charAt(0))
                    && depth < MAX_DEPTH && !visiting.contains(owner + "." + field)) {
                Map<String, List<String>> fields = clinitFields(owner);
                List<String> nested = fields.get(field);
                if (nested != null && !nested.isEmpty()) {
                    visiting.add(owner + "." + field);
                    List<String> inner = expand(nested, owner, depth + 1, visiting);
                    visiting.remove(owner + "." + field);
                    return List.of(t + "{" + String.join(", ", inner) + "}");
                }
            }
        }
        return List.of(t);
    }

    static final Map<String, String> ownerIndex = new HashMap<>();

    /** Maps a short owner ("Vec3", "ChatType$Bound") back to an internal name seen in bytecode. */
    static String resolveOwner(String shortName, String self) {
        return ownerIndex.get(shortName);
    }

    static String shortOwner(String internalName) {
        String s = internalName.substring(internalName.lastIndexOf('/') + 1);
        ownerIndex.putIfAbsent(s, internalName);
        return s;
    }

    // ---- method walking (read side) -----------------------------------------------

    static String walkLambda(String lambdaToken, int depth, Set<String> visiting) {
        String[] parts = lambdaToken.split("\\|", 4);
        if (parts.length < 4) return lambdaToken;
        String owner = parts[1], name = parts[2], desc = parts[3];
        List<String> toks = walkMethod(owner, name, desc, depth, visiting);
        return "λ" + shortOwner(owner) + "." + name + "{" + String.join(", ", toks) + "}";
    }

    /** Walks the overload of owner.name that takes a (Registry)FriendlyByteBuf. */
    static List<String> walkMethodByName(String owner, String name) {
        ClassModel cm = classModel(owner);
        if (cm == null) return List.of("opaque:" + shortOwner(owner) + "." + name);
        String desc = "";
        for (MethodModel m : cm.methods()) {
            if (!m.methodName().stringValue().equals(name)) continue;
            String d = m.methodType().stringValue();
            if (d.contains("FriendlyByteBuf")) { desc = d; break; }
        }
        if (desc.isEmpty()) return List.of("opaque:" + shortOwner(owner) + "." + name + "(no buffer overload)");
        return walkMethod(owner, name, desc, 0, new HashSet<>());
    }

    static List<String> walkMethod(String owner, String name, String desc, int depth, Set<String> visiting) {
        String key = owner + "." + name + desc;
        if (depth > MAX_DEPTH || visiting.contains(key)) return List.of("…");
        ClassModel cm = classModel(owner);
        if (cm == null) return List.of("opaque:" + shortOwner(owner) + "." + name);
        MethodModel target = null;
        for (MethodModel m : cm.methods()) {
            if (!m.methodName().stringValue().equals(name)) continue;
            if (desc.isEmpty() || m.methodType().stringValue().equals(desc)) { target = m; break; }
        }
        if (target == null) return List.of("opaque:" + shortOwner(owner) + "." + name);
        visiting.add(key);
        List<String> out = new ArrayList<>();
        String pendingCodec = null;
        for (CodeElement el : target.code().map(c -> (Iterable<CodeElement>) c).orElse(List.of())) {
            if (el instanceof FieldInstruction fi) {
                if (fi.opcode() == Opcode.GETSTATIC && isCodecField(fi)) {
                    pendingCodec = shortOwner(fi.owner().asInternalName()) + "." + fi.name().stringValue();
                }
            } else if (el instanceof InvokeInstruction ii) {
                String o = ii.owner().asInternalName();
                String n = ii.name().stringValue();
                String d = ii.typeSymbol().descriptorString();
                String oShort = shortOwner(o);
                if (o.endsWith("FriendlyByteBuf") || o.endsWith("RegistryFriendlyByteBuf") || o.endsWith("io/netty/buffer/ByteBuf")
                        || o.endsWith("codec/VarInt") || o.endsWith("codec/VarLong") || o.endsWith("codec/Utf8String")) {
                    if (n.startsWith("read") || n.startsWith("write")) out.add("buf." + n);
                } else if (n.equals("decode") || n.equals("encode")) {
                    if (pendingCodec != null) {
                        out.addAll(expandOne(pendingCodec, owner, depth + 1, visiting));
                    } else {
                        out.add(oShort + "." + n);
                    }
                } else if (o.startsWith("net/minecraft/") && (d.contains("FriendlyByteBuf") || d.contains("io/netty/buffer/ByteBuf"))) {
                    List<String> nested = walkMethod(o, n, d, depth + 1, visiting);
                    out.add(oShort + "." + n + "{" + String.join(", ", nested) + "}");
                } else if (ii.typeSymbol().returnType().descriptorString().contains("StreamCodec")) {
                    pendingCodec = oShort + "." + n;
                }
                if (!n.equals("decode") && !n.equals("encode") && !ii.typeSymbol().returnType().descriptorString().contains("StreamCodec")) {
                    pendingCodec = null;
                }
            } else if (el instanceof InvokeDynamicInstruction idi) {
                String lt = lambdaToken(idi);
                String[] parts = lt.split("\\|", 4);
                if (parts.length == 4 && parts[1].startsWith("net/minecraft/")
                        && (parts[3].contains("FriendlyByteBuf") || parts[3].contains("io/netty/buffer/ByteBuf"))) {
                    out.add(walkLambda(lt, depth + 1, visiting));
                }
            }
        }
        visiting.remove(key);
        return out;
    }

    static String json(String s) {
        StringBuilder sb = new StringBuilder("\"");
        for (char c : s.toCharArray()) {
            switch (c) {
                case '"' -> sb.append("\\\"");
                case '\\' -> sb.append("\\\\");
                case '\n' -> sb.append("\\n");
                default -> { if (c < 0x20) sb.append(String.format("\\u%04x", (int) c)); else sb.append(c); }
            }
        }
        return sb.append('"').toString();
    }
}
