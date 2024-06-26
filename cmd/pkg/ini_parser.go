package parser

import (
	"fmt"
	"regexp"
	"strings"
)

type IniParser struct {
	numberOfSections int
	sections         map[string]map[string]string
	sectionsNameList []string
	// LoadFromString  error
	// LoadFromFile    error
	// GetSectionNames []string
	// GetSections     map[string]map[string]string
	// Get(section_name, key)
	// Set(section_name, key, value)
	// ToString
	// SaveToFile
}

func NewParser() IniParser {
	newSections := make(map[string]map[string]string)
	newParsing := IniParser{}
	newParsing.sections = newSections
	newParsing.numberOfSections = 0
	return newParsing
}

func (parser *IniParser) Set(section, key, value string) [3]string {
	if parser.sections[section] == nil {
		parser.sections[section] = make(map[string]string)
	}
	parser.sections[section][key] = value
	res := [...]string{section, key, value}
	return res
}

func (parser *IniParser) LoadFromString(str string) {
	sectionRegex,_:=regexp.Compile(`^\[.+\]$`)
	sectionsNames := strings.Split(str, "\n")
	for _, section := range sectionsNames {
		if sectionRegex.Match([]byte(section)) {
			section:=strings.TrimPrefix(strings.TrimSuffix(section,"]"),"[")
			fmt.Println(section)
			parser.sections[section] = make(map[string]string)
			parser.sectionsNameList = append(parser.sectionsNameList, string(section[1:]))
			parser.numberOfSections++
		} else {
			keysAndValues := strings.Split(section, "\n")
			for _, pair := range keysAndValues {
				if string(pair) != "#" {
					pair := strings.Split(pair, "=")
					parser.Set(parser.sectionsNameList[parser.numberOfSections-1], pair[0],pair[1])
				}
			}

		}
	}
}


