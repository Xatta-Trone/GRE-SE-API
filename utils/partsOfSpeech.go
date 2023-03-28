package utils

func GetPos(pos string) string {
	switch pos {
	case "verb", "verb.", "v", "v.":
		return "verb"
	case "adjective", "adjective.", "adj", "adj.":
		return "adjective"
	case "noun", "noun.", "n", "n.":
		return "noun"
	case "adverb", "adverb.", "adv", "adv.":
		return "adverb"
	default:
		return ""
	}
}
