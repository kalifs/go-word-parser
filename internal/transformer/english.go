package transformer

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"
)

type lineRegexps struct {
	word       *regexp.Regexp
	definition *regexp.Regexp
	tags       *regexp.Regexp
}

// TransformEngilishDefinitions transofrm english dictionary to 5-letter word list
func TransformEngilishDefinitions() {
	regexps := lineRegexps{
		word:       regexp.MustCompile(`<p><ent>(\w+)\<`),
		definition: regexp.MustCompile(`.*<pos>n\..+<def>(.+)</def>`),
		tags:       regexp.MustCompile(`(<\w+>)|(</\w+>)`),
	}

	files, err := filepath.Glob("data/english/CIDE.*")
	check(err)

	results := *processFiles(files, &regexps)

	fmt.Println(len(results))

	file, err := os.Create("output/english.tsv")
	check(err)
	defer file.Close()

	fmt.Println("Dumping results to output/english.tsv")
	for word, definition := range results {
		line := fmt.Sprintf("%s\t%s\n", word, definition)
		file.WriteString(line)
	}
}

func processFiles(files []string, regexps *lineRegexps) *map[string]string {
	results := make(map[string]string)
	for _, path := range files {
		fmt.Printf("Processing: %s", path)
		letterResults := filterWords(parseFile(path, &*regexps), 5)
		fmt.Printf(" found %d \n", len(letterResults))

		for word, def := range letterResults {
			results[word] = def
		}
	}
	return &results
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

func parseFile(path string, regexps *lineRegexps) map[string]string {
	file, err := os.Open(path)
	check(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)

	var line string
	var found bool = false
	var word, definition string = "", ""
	isWordLine := false
	results := make(map[string]string)

	for scanner.Scan() {
		line = scanner.Text()
		if !isWordLine {
			word, isWordLine = matchWord(line, &*regexps)
		} else {
			definition, found = matchDefinition(line, &*regexps)
			if found {
				results[word] = removeTags(definition, &*regexps)
			}
			isWordLine = false
		}
	}
	return results
}

func matchWord(line string, regexps *lineRegexps) (string, bool) {
	matches := regexps.word.FindStringSubmatch(line)
	if matches != nil {
		return strings.ToLower(matches[1]), true
	}
	return "", false
}

func matchDefinition(line string, regexps *lineRegexps) (string, bool) {
	matches := regexps.definition.FindStringSubmatch(line)
	if matches != nil {
		return matches[1], true
	}
	return "", false
}

func removeTags(text string, regexps *lineRegexps) string {
	return string(regexps.tags.ReplaceAll([]byte(text), []byte("")))
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
