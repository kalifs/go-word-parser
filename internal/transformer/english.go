package transformer

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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

	results := *processEnlishDataFiles(files, &regexps)

	fmt.Println(len(results))

	dumpFile(&results, "english")
}

func processEnlishDataFiles(files []string, regexps *lineRegexps) *map[string]string {
	results := make(map[string]string)
	for _, path := range files {
		fmt.Printf("Processing: %s", path)
		letterResults := filterWords(parseEnglishDataFile(path, &*regexps), 5)
		fmt.Printf(" found %d \n", len(letterResults))

		for word, def := range letterResults {
			results[word] = def
		}
	}
	return &results
}

func parseEnglishDataFile(path string, regexps *lineRegexps) map[string]string {
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
