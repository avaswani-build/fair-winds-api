package service

import "github.com/avaswani-build/fair-winds-api/internal/domain"

func Recommend(f domain.Forecast) domain.Recommendation {
	switch {
	case f.GustKts > 22:
		return domain.Recommendation{
			Status: "no-go",
			Reason: "Conditions are not beginner-friendly due to strong gusts.",
		}
	case f.WindAvgKts < 6:
		return domain.Recommendation{
			Status: "marginal",
			Reason: "Wind is likely too light for a worthwhile session.",
		}
	case f.WindAvgKts >= 8 && f.WindAvgKts <= 16 && (f.GustKts-f.WindAvgKts) <= 6:
		return domain.Recommendation{
			Status: "go",
			Reason: "Conditions look suitable for an intermediate sailor.",
		}
	default:
		return domain.Recommendation{
			Status: "conditional",
			Reason: "Usable conditions, but caution is warranted based on wind variability.",
		}
	}
}
