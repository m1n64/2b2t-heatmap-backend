package highways

type Line struct {
	X1, Z1 int
	X2, Z2 int
	Kind   string
}

const (
	mcMinI = -30_000_000
	mcMaxI = 30_000_000
	stride = 500_000
)

var ringRoads = []int{
	200, 500, 1000, 1500, 2000, 2500,
	5000, 7500, 10000, 15000, 20000,
	25000, 50000, 55000, 62500, 75000,
	100000, 125000, 250000, 500000,
	750000, 1000000, 1250000, 1568852,
	1875000, 2500000, 3750000,
}

var diamonds = []int{
	2500, 5000, 25000, 50000,
	125000, 250000, 500000,
	3750000,
}

type HighwayGenerator struct {
	Min    int
	Max    int
	Stride int
}

func NewHighwayGenerator() *HighwayGenerator {
	return &HighwayGenerator{
		Min:    -30_000_000,
		Max:    30_000_000,
		Stride: 500_000,
	}
}

func (g *HighwayGenerator) GenerateNether() []Line {
	lines := make([]Line, 0, 800)

	g.addCardinals(&lines)
	g.addDiagonals(&lines)
	g.addRingRoads(&lines)
	g.addDiamonds(&lines)
	g.addGrid50k(&lines)
	g.addExtensions50k(&lines)

	return lines
}

// func GenerateOverworld

func (g *HighwayGenerator) addCardinals(lines *[]Line) {
	for i := mcMinI; i < mcMaxI; i += stride {
		*lines = append(*lines,
			Line{i, 0, i + stride, 0, "cardinal"},
			Line{0, i, 0, i + stride, "cardinal"},
		)
	}
}

func (g *HighwayGenerator) addDiagonals(lines *[]Line) {
	for i := mcMinI; i < mcMaxI; i += stride {
		*lines = append(*lines,
			Line{i, i, i + stride, i + stride, "diagonal"},
			Line{-i, i, -i - stride, i + stride, "diagonal"},
		)
	}
}

func (g *HighwayGenerator) addRingRoads(lines *[]Line) {
	for _, r := range ringRoads {
		for i := -r; i < r; i += stride {
			end := min(i+stride, r)
			*lines = append(*lines,
				Line{i, -r, end, -r, "ring"},
				Line{i, r, end, r, "ring"},
				Line{-r, i, -r, end, "ring"},
				Line{r, i, r, end, "ring"},
			)
		}
	}
}

func (g *HighwayGenerator) addDiamonds(lines *[]Line) {
	for _, d := range diamonds {
		*lines = append(*lines,
			Line{d, 0, 0, d, "diamond"},
			Line{0, -d, d, 0, "diamond"},
			Line{0, -d, -d, 0, "diamond"},
			Line{-d, 0, 0, d, "diamond"},
		)
	}
}

func (g *HighwayGenerator) addGrid50k(lines *[]Line) {
	for i := -50000; i < 50000; i += 5000 {
		if i == 0 {
			continue
		}
		*lines = append(*lines,
			Line{i, -50000, i, 50000, "grid"},
			Line{-50000, i, 50000, i, "grid"},
		)
	}
}

func (g *HighwayGenerator) addExtensions50k(lines *[]Line) {
	*lines = append(*lines,
		Line{-125000, -50000, -50000, -50000, "extend"},
		Line{-125000, 50000, -50000, 50000, "extend"},
		Line{125000, -50000, 50000, -50000, "extend"},
		Line{125000, 50000, 50000, 50000, "extend"},
		Line{-50000, -125000, -50000, -50000, "extend"},
		Line{-50000, 125000, -50000, 50000, "extend"},
		Line{50000, -125000, 50000, -50000, "extend"},
		Line{50000, 125000, 50000, 50000, "extend"},
	)
}
