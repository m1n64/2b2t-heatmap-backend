package highways

type FeatureCollection struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

type Feature struct {
	Type       string     `json:"type"`
	Properties Properties `json:"properties"`
	Geometry   Geometry   `json:"geometry"`
}

type Properties struct {
	Kind string `json:"kind"`
}

type Geometry struct {
	Type        string      `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
}

type Projector struct {
	McMin, McMax       float64
	WorldMin, WorldMax float64
	scaleFactor        float64
}

func NewDefaultProjector() *Projector {
	p := &Projector{
		McMin:    -30_000_000,
		McMax:    30_000_000,
		WorldMin: -30000,
		WorldMax: 30000,
	}
	p.scaleFactor = (p.WorldMax - p.WorldMin) / (p.McMax - p.McMin)
	return p
}

func (p *Projector) scale(v int) float64 {
	return float64(v) * p.scaleFactor
}

func (p *Projector) inWorld(v float64) bool {
	return v >= p.WorldMin && v <= p.WorldMax
}

func (p *Projector) ToGeoJSON(lines []Line) FeatureCollection {
	fc := FeatureCollection{Type: "FeatureCollection"}

	for _, l := range lines {
		x1, z1 := p.scale(l.X1), p.scale(l.Z1)
		x2, z2 := p.scale(l.X2), p.scale(l.Z2)

		if !p.inWorld(x1) && !p.inWorld(x2) &&
			!p.inWorld(z1) && !p.inWorld(z2) {
			continue
		}

		fc.Features = append(fc.Features, Feature{
			Type: "Feature",
			Properties: Properties{
				Kind: l.Kind,
			},
			Geometry: Geometry{
				Type: "LineString",
				Coordinates: [][]float64{
					{x1, -z1},
					{x2, -z2},
				},
			},
		})
	}

	return fc
}
