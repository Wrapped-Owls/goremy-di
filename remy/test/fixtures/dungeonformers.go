package fixtures

type Dungeonformer interface {
	Name() string
	Class() string
}

type (
	DnDBase struct {
		ClassName string // only required field for all
	}
	GenericDungeonformer struct {
		DnDBase
		Name   string
		Height uint32
	}
)

func NewGenericDungeonformer(name string, height uint32) GenericDungeonformer {
	return GenericDungeonformer{Name: name, Height: height}
}

func (b *DnDBase) Class() string {
	return b.ClassName
}

// Bumblelf = Bumblebee + Elf
type Bumblelf struct {
	DnDBase
}

func NewBumblelf() *Bumblelf {
	return &Bumblelf{
		DnDBase: DnDBase{ClassName: "Elf Scout"},
	}
}

func (t *Bumblelf) Name() string {
	return "Bumblelf"
}

// OptimadinPrime = Optimus Prime + Paladin
type OptimadinPrime struct {
	DnDBase
}

func NewOptimadinPrime() *OptimadinPrime {
	return &OptimadinPrime{
		DnDBase: DnDBase{ClassName: "Paladin Commander"},
	}
}

func (t *OptimadinPrime) Name() string {
	return "OptimadinPrime"
}

// MegadwarfTron = Megatron + Dwarf
type MegadwarfTron struct {
	DnDBase
}

func NewMegadwarfTron() *MegadwarfTron {
	return &MegadwarfTron{
		DnDBase: DnDBase{ClassName: "Dwarven Warlord"},
	}
}

func (t *MegadwarfTron) Name() string {
	return "MegadwarfTron"
}

// Sorcerscream = Starscream + Sorcerer
type Sorcerscream struct {
	DnDBase
}

func NewSorcerscream() *Sorcerscream {
	return &Sorcerscream{
		DnDBase: DnDBase{ClassName: "Screaming Sorcerer"},
	}
}

func (t *Sorcerscream) Name() string {
	return "Sorcerscream"
}

// Soundbard = Soundwave + Bard
type Soundbard struct {
	DnDBase
}

func NewSoundbard() *Soundbard {
	return &Soundbard{
		DnDBase: DnDBase{ClassName: "Shadow Bard"},
	}
}

func (t *Soundbard) Name() string {
	return "Soundbard"
}

// Ironknight = Ironhide + Knight
type Ironknight struct {
	DnDBase
}

func NewIronknight() *Ironknight {
	return &Ironknight{
		DnDBase: DnDBase{ClassName: "Armored Knight"},
	}
}

func (t *Ironknight) Name() string {
	return "Ironknight"
}

// Ratcheric = Ratchet + Cleric
type Ratcheric struct {
	DnDBase
}

func NewRatcheric() *Ratcheric {
	return &Ratcheric{
		DnDBase: DnDBase{ClassName: "Battle Cleric"},
	}
}

func (t *Ratcheric) Name() string {
	return "Ratcheric"
}

// Jazogue = Jazz + Rogue
type Jazogue struct {
	DnDBase
}

func NewJazogue() *Jazogue {
	return &Jazogue{
		DnDBase: DnDBase{ClassName: "Agile Rogue"},
	}
}

func (t *Jazogue) Name() string {
	return "Jazogue"
}

// Wheelificer = Wheeljack + Artificer
type Wheelificer struct {
	DnDBase
}

func NewWheelificer() *Wheelificer {
	return &Wheelificer{
		DnDBase: DnDBase{ClassName: "Arcane Artificer"},
	}
}

func (t *Wheelificer) Name() string {
	return "Wheelificer"
}

// Shocklock = Shockwave + Warlock
type Shocklock struct {
	DnDBase
}

func NewShocklock() *Shocklock {
	return &Shocklock{
		DnDBase: DnDBase{ClassName: "Eldritch Warlock"},
	}
}

func (t *Shocklock) Name() string {
	return "Shocklock"
}
