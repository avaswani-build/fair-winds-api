package weather

func degreesToCardinal(deg float64) string {
	directions := []string{
		"N", "NE", "E", "SE",
		"S", "SW", "W", "NW",
	}

	idx := int((deg + 22.5) / 45.0)
	return directions[idx%8]
}
