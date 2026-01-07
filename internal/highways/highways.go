package highways

type HighwayService struct {
	highwayGenerator *HighwayGenerator
	projector        *Projector
}

func NewHighwayService(highwayGenerator *HighwayGenerator, projector *Projector) *HighwayService {
	return &HighwayService{
		highwayGenerator: highwayGenerator,
		projector:        projector,
	}
}

func (h *HighwayService) GenerateNether() FeatureCollection {
	lines := h.highwayGenerator.GenerateNether()

	return h.projector.ToGeoJSON(lines)
}
