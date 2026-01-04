package enums

type World string

const (
	Nether    World = "nether"
	Overworld World = "overworld"
	End       World = "end"
)

var allWorlds = []World{Nether, Overworld, End}

func AvailableWorlds() []World {
	return allWorlds
}

func (w World) String() string {
	return string(w)
}

func (w World) IsValid() bool {
	switch w {
	case Nether, Overworld, End:
		return true
	default:
		return false
	}
}

func ParseWorld(s string) (World, bool) {
	w := World(s)
	if w.IsValid() {
		return w, true
	}
	return "", false
}
