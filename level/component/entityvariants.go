package component

import (
	pk "github.com/Tnze/go-mc/net/packet"
)

// Simple VarInt variant components for entity types.

var _ DataComponent = (*VillagerVariant)(nil)

type VillagerVariant struct{ pk.VarInt }

func (VillagerVariant) ID() string { return "minecraft:villager/variant" }

var _ DataComponent = (*WolfVariant)(nil)

type WolfVariant struct{ pk.VarInt }

func (WolfVariant) ID() string { return "minecraft:wolf/variant" }

var _ DataComponent = (*WolfSoundVariant)(nil)

type WolfSoundVariant struct{ pk.VarInt }

func (WolfSoundVariant) ID() string { return "minecraft:wolf/sound_variant" }

var _ DataComponent = (*WolfCollar)(nil)

type WolfCollar struct{ pk.VarInt }

func (WolfCollar) ID() string { return "minecraft:wolf/collar" }

var _ DataComponent = (*FoxVariant)(nil)

type FoxVariant struct{ pk.VarInt }

func (FoxVariant) ID() string { return "minecraft:fox/variant" }

var _ DataComponent = (*SalmonSize)(nil)

type SalmonSize struct{ pk.VarInt }

func (SalmonSize) ID() string { return "minecraft:salmon/size" }

var _ DataComponent = (*ParrotVariant)(nil)

type ParrotVariant struct{ pk.VarInt }

func (ParrotVariant) ID() string { return "minecraft:parrot/variant" }

var _ DataComponent = (*TropicalFishPattern)(nil)

type TropicalFishPattern struct{ pk.VarInt }

func (TropicalFishPattern) ID() string { return "minecraft:tropical_fish/pattern" }

var _ DataComponent = (*TropicalFishBaseColor)(nil)

type TropicalFishBaseColor struct{ pk.VarInt }

func (TropicalFishBaseColor) ID() string { return "minecraft:tropical_fish/base_color" }

var _ DataComponent = (*TropicalFishPatternColor)(nil)

type TropicalFishPatternColor struct{ pk.VarInt }

func (TropicalFishPatternColor) ID() string { return "minecraft:tropical_fish/pattern_color" }

var _ DataComponent = (*MooshroomVariant)(nil)

type MooshroomVariant struct{ pk.VarInt }

func (MooshroomVariant) ID() string { return "minecraft:mooshroom/variant" }

var _ DataComponent = (*RabbitVariant)(nil)

type RabbitVariant struct{ pk.VarInt }

func (RabbitVariant) ID() string { return "minecraft:rabbit/variant" }

var _ DataComponent = (*PigVariant)(nil)

type PigVariant struct{ pk.VarInt }

func (PigVariant) ID() string { return "minecraft:pig/variant" }

var _ DataComponent = (*CowVariant)(nil)

type CowVariant struct{ pk.VarInt }

func (CowVariant) ID() string { return "minecraft:cow/variant" }

// ChickenVariant (wire 86) is NOT here — it's a registryEntryHolder<string>, see chickenvariant.go

var _ DataComponent = (*FrogVariant)(nil)

type FrogVariant struct{ pk.VarInt }

func (FrogVariant) ID() string { return "minecraft:frog/variant" }

var _ DataComponent = (*HorseVariant)(nil)

type HorseVariant struct{ pk.VarInt }

func (HorseVariant) ID() string { return "minecraft:horse/variant" }

var _ DataComponent = (*LlamaVariant)(nil)

type LlamaVariant struct{ pk.VarInt }

func (LlamaVariant) ID() string { return "minecraft:llama/variant" }

var _ DataComponent = (*AxolotlVariant)(nil)

type AxolotlVariant struct{ pk.VarInt }

func (AxolotlVariant) ID() string { return "minecraft:axolotl/variant" }

var _ DataComponent = (*CatVariant)(nil)

type CatVariant struct{ pk.VarInt }

func (CatVariant) ID() string { return "minecraft:cat/variant" }

var _ DataComponent = (*CatCollar)(nil)

type CatCollar struct{ pk.VarInt }

func (CatCollar) ID() string { return "minecraft:cat/collar" }

var _ DataComponent = (*SheepColor)(nil)

type SheepColor struct{ pk.VarInt }

func (SheepColor) ID() string { return "minecraft:sheep/color" }

var _ DataComponent = (*ShulkerColor)(nil)

type ShulkerColor struct{ pk.VarInt }

func (ShulkerColor) ID() string { return "minecraft:shulker/color" }
