package transformer

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Root struct {
	XMLName xml.Name `xml:"TEI"`
	Body    Entries  `xml:"body"`
}
type Entries struct {
	XMLName xml.Name `xml:"body"`
	Entries []Entry  `xml:"entry"`
}

type Entry struct {
	XMLName    xml.Name `xml:"entry"`
	Identifier string   `xml:"id,attr"`
	Form       string   `xml:"form"`
	Senses     []Sense  `xml:"sense"`
}

type Sense struct {
	XMLName    xml.Name `xml:"sense"`
	Identifier string   `xml:"id,attr"`
	Def        string   `xml:"def"`
}

// TransformLatvianDefinitions parses latvian dicitionary and create output
func TransformLatvianDefinitions() {
	xmlFile, err := os.Open("data/latvian/tezaurs.xml")
	check(err)

	defer xmlFile.Close()

	fmt.Println("Loading dictionary")
	byteValue, _ := ioutil.ReadAll(xmlFile)

	fmt.Println("Marshaling XML")
	var root Root
	xml.Unmarshal(byteValue, &root)
	entries := root.Body.Entries
	fmt.Println("Data loaded!")

	results := make(map[string]string)
	word := ""
	for i := 0; i < len(entries); i++ {
		word = entries[i].Form
		parts := strings.Split(word, "")
		lastLetter := parts[len(parts)-1]
		if word == strings.ToLower(word) && (lastLetter == "š" || lastLetter == "s" || lastLetter == "a" || lastLetter == "e") {
			results[word] = entries[i].Senses[0].Def
		}
	}

	results = filterWords(results, 5)
	fmt.Println(len(results))
	dumpFile(&results, "latvian")
}
