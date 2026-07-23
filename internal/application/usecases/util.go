package usecases

func makeEmptyUnknown(inp *string) {
	if *inp == "" {
		*inp = "unknown"
	}
}
