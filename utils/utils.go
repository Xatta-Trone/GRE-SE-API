package utils

import (
	"strings"
	"unicode"

	"github.com/fatih/color"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func PrintS(str string) {
	c := color.New(color.FgCyan).Add(color.Underline)
	c.Println(str)
}

func PrintG(str string) {
	c := color.New(color.FgGreen).Add(color.Underline)
	c.Println(str)
}

func PrintR(str string) {
	color.Red(str)
}

var normalizer = transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)

func ProcessWord(str string) []string {

	str = strings.TrimSpace(strings.Join(strings.Fields(str), " "))

	s, _, err := transform.String(normalizer, str)

	if err != nil {
		return []string{}
	}

	str = strings.ToLower(s)
	// replace underscores and slashes
	processedWord := strings.Replace(str, "'", "", -1)
	processedWord = strings.Replace(processedWord, "\\", "", -1)
	processedWord = strings.Replace(processedWord, "_", "-", -1)

	matchWords := match(processedWord)

	return matchWords
}

func match(s string) []string {
	str := []string{}
	i := strings.Index(s, "(")
	if i >= 0 {
		j := strings.Index(s, ")")
		if j >= 0 {
			//return s[i+1 : j]
			str = append(str, s[i+1:j])
			str = append(str, s[:i])
		}
	} else {
		str = append(str, s)
	}

	return str
}
