package transformer

import (
	"fmt"
	"os"
	"unicode/utf8"
)

func dumpFile(results *map[string]string, language string) {
	file, err := os.Create(fmt.Sprintf("output/%s.tsv", language))
	check(err)
	defer file.Close()

	fmt.Printf("Dumping results to output/%s.tsv\n", language)
	for word, definition := range *results {
		line := fmt.Sprintf("%s\t%s\n", word, definition)
		file.WriteString(line)
	}
}

func filterWords(data map[string]string, limit int) map[string]string {
	results := make(map[string]string)
	for word, definition := range data {
		if utf8.RuneCountInString(word) == limit {
			results[word] = definition
		}
	}
	return results
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
