package dto

type WorkoutGeneratorBlock struct {
	BlockName    string   `json:"block_name"`
	BlockType    string   `json:"block_type"`
	Instructions string   `json:"instructions"`
	Exercises    []string `json:"exercises"`
}

type WorkoutGeneratorResponse struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Difficulty  string                  `json:"difficulty"`
	Blocks      []WorkoutGeneratorBlock `json:"blocks"`
}
