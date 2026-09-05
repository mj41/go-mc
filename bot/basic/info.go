package basic

import (
	"io"
	"unsafe"

	"github.com/Tnze/go-mc/data/packetid"
	pk "github.com/Tnze/go-mc/net/packet"
)

// WorldInfo content player info in server.
type WorldInfo struct {
	DimensionType       int32
	DimensionNames      []string // Identifiers for all worlds on the server.
	DimensionName       string   // Name of the world being spawned into.
	HashedSeed          int64    // First 8 bytes of the SHA-256 hash of the world's seed. Used client side for biome noise
	MaxPlayers          int32    // Was once used by the client to draw the player list, but now is ignored.
	ViewDistance        int32    // Render distance (2-32).
	SimulationDistance  int32    // The distance that the client will process specific things, such as entities.
	ReducedDebugInfo    bool     // If true, a vanilla client shows reduced information on the debug screen. For servers in development, this should almost always be false.
	EnableRespawnScreen bool     // Set to false when the doImmediateRespawn gamerule is true.
	IsDebug             bool     // True if the world is a debug mode world; debug mode worlds cannot be modified and have predefined blocks.
	IsFlat              bool     // True if the world is a superflat world; flat worlds have different void fog and a horizon at y=0 instead of y=63.
	DoLimitCrafting     bool     // Whether players can only craft recipes they have already unlocked. Currently unused by the client.
	HasDeathLocation    bool     // True if DeathDimension/DeathLocation are set (player died before and has not respawned).
	DeathDimension      string   // Dimension the player died in.
	DeathLocation       pk.Position
	PortalCooldown      int32 // Ticks until the player can use a portal again.
	SeaLevel            int32 // Sea level of the current dimension (1.21.2+).
	EnforcesSecureChat  bool  // Whether the server requires signed chat (login packet only).
	OnlineMode          bool  // Whether the server runs in online mode (26.2+, login packet only).
}

type PlayerInfo struct {
	EID          int32 // The player's Entity ID (EID).
	Hardcore     bool  // Is hardcore
	Gamemode     byte  // Gamemode. 0: Survival, 1: Creative, 2: Adventure, 3: Spectator.
	PrevGamemode int8  // Previous Gamemode
}

// globalPos is the wire form of a GlobalPos: dimension identifier + packed block position.
type globalPos struct {
	Dimension pk.Identifier
	Pos       pk.Position
}

func (g *globalPos) ReadFrom(r io.Reader) (int64, error) {
	return pk.Tuple{&g.Dimension, &g.Pos}.ReadFrom(r)
}

func (g globalPos) WriteTo(w io.Writer) (int64, error) {
	return pk.Tuple{g.Dimension, g.Pos}.WriteTo(w)
}

func (p *Player) setDeathLocation(loc pk.Option[globalPos, *globalPos]) {
	p.HasDeathLocation = bool(loc.Has)
	if loc.Has {
		p.DeathDimension = string(loc.Val.Dimension)
		p.DeathLocation = loc.Val.Pos
	} else {
		p.DeathDimension = ""
		p.DeathLocation = pk.Position{}
	}
}

// handleLoginPacket parses ClientboundLogin (protocol 776 / MC 26.2 layout):
// entityID, hardcore, dimensionNames, maxPlayers, viewDistance, simulationDistance,
// reducedDebugInfo, enableRespawnScreen, doLimitedCrafting, CommonPlayerSpawnInfo
// (dimensionType, dimensionName, hashedSeed, gamemode, previousGamemode, isDebug, isFlat,
// deathLocation?, portalCooldown, seaLevel), onlineMode, enforcesSecureChat.
func (p *Player) handleLoginPacket(packet pk.Packet) error {
	var deathLocation pk.Option[globalPos, *globalPos]
	err := packet.Scan(
		(*pk.Int)(&p.EID),
		(*pk.Boolean)(&p.Hardcore),
		pk.Array((*[]pk.Identifier)(unsafe.Pointer(&p.DimensionNames))),
		(*pk.VarInt)(&p.MaxPlayers),
		(*pk.VarInt)(&p.ViewDistance),
		(*pk.VarInt)(&p.SimulationDistance),
		(*pk.Boolean)(&p.ReducedDebugInfo),
		(*pk.Boolean)(&p.EnableRespawnScreen),
		(*pk.Boolean)(&p.DoLimitCrafting),
		(*pk.VarInt)(&p.WorldInfo.DimensionType),
		(*pk.Identifier)(&p.DimensionName),
		(*pk.Long)(&p.HashedSeed),
		(*pk.UnsignedByte)(&p.Gamemode),
		(*pk.Byte)(&p.PrevGamemode),
		(*pk.Boolean)(&p.IsDebug),
		(*pk.Boolean)(&p.IsFlat),
		&deathLocation,
		(*pk.VarInt)(&p.PortalCooldown),
		(*pk.VarInt)(&p.SeaLevel),
		(*pk.Boolean)(&p.OnlineMode),         // 26.2+ (protocol 776)
		(*pk.Boolean)(&p.EnforcesSecureChat), // after online mode since 26.2
	)
	if err != nil {
		return Error{err}
	}
	p.setDeathLocation(deathLocation)
	err = p.c.Conn.WritePacket(pk.Marshal( // PluginMessage packet
		packetid.ServerboundCustomPayload,
		pk.Identifier("minecraft:brand"),
		pk.String(p.Settings.Brand),
	))
	if err != nil {
		return Error{err}
	}

	err = p.c.Conn.WritePacket(pk.Marshal(
		packetid.ServerboundClientInformation, // Client settings
		pk.String(p.Settings.Locale),
		pk.Byte(p.Settings.ViewDistance),
		pk.VarInt(p.Settings.ChatMode),
		pk.Boolean(p.Settings.ChatColors),
		pk.UnsignedByte(p.Settings.DisplayedSkinParts),
		pk.VarInt(p.Settings.MainHand),
		pk.Boolean(p.Settings.EnableTextFiltering),
		pk.Boolean(p.Settings.AllowListing),
		pk.VarInt(p.Settings.ParticleStatus), // 1.21.4+: 0=all, 1=decreased, 2=minimal
	))
	if err != nil {
		return Error{err}
	}

	p.resetKeepAliveDeadline()
	return nil
}

// handleRespawnPacket parses ClientboundRespawn: CommonPlayerSpawnInfo followed by
// the dataToKeep bit flags (0x01 keep attributes, 0x02 keep metadata).
func (p *Player) handleRespawnPacket(packet pk.Packet) error {
	var deathLocation pk.Option[globalPos, *globalPos]
	var dataToKeep pk.Byte
	err := packet.Scan(
		(*pk.VarInt)(&p.DimensionType),
		(*pk.Identifier)(&p.DimensionName),
		(*pk.Long)(&p.HashedSeed),
		(*pk.UnsignedByte)(&p.Gamemode),
		(*pk.Byte)(&p.PrevGamemode),
		(*pk.Boolean)(&p.IsDebug),
		(*pk.Boolean)(&p.IsFlat),
		&deathLocation,
		(*pk.VarInt)(&p.PortalCooldown),
		(*pk.VarInt)(&p.SeaLevel),
		&dataToKeep,
	)
	if err != nil {
		return Error{err}
	}
	p.setDeathLocation(deathLocation)
	return nil
}
