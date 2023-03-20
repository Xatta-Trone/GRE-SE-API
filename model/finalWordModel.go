package model

type Combined struct {
	PartsOfSpeech string   `json:"partsOfSpeech"`
	Definitions   []string `json:"definitions"`
	Examples      []string `json:"examples"`
	SynonymsG     []string `json:"synonyms_gre"`
	SynonymsN     []string `json:"synonyms_normal"`
}

type CombinedWithWord struct {
	Word            string     `json:"word"`
	PartsOfSpeeches []Combined `json:"partsOfSpeeches"`
}
