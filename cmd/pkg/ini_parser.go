package parser

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"golang.org/x/exp/maps"
)

type IniParser struct {
	numberOfSections            int
	sections                    map[string]map[string]string
	sectionsNameList            []string
	errCouldNotLoadDataFromFile error
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
	newParsing.errCouldNotLoadDataFromFile = errors.New("Couldn't load data from file")
	return newParsing
}

func (parser *IniParser) Set(section, key, value string) (sectionV, keyV, valueV string) {
	if parser.sections[section] == nil {
		parser.sections[section] = make(map[string]string)
	}
	parser.sections[section][key] = value
	return section, key, value
}

func (parser *IniParser) Get(section_name, key string) (string, error) {

	if value, ok := parser.sections[section_name][key]; ok {
		return value, nil
	} else {
		errNotFound := errors.New("The given data is not valid")
		return value, errNotFound
	}
}

func (parser *IniParser) GetSectionNames() []string {
	return parser.sectionsNameList
}

func (parser *IniParser) GetSections() (map[string]map[string]string, error) {
	if parser.sections == nil {
		errNoSections := errors.New("No avaliable sections to return")
		return nil, errNoSections
	}
	return parser.sections, nil
}

func (parser *IniParser) LoadFromString(str string) error {
	strLower := strings.ToLower(str)
	sectionRegex, _ := regexp.Compile(`^\[.+\]$`)
	sectionsNames := strings.Split(strLower, "\n")
	for _, section := range sectionsNames {
		if sectionRegex.Match([]byte(section)) {
			section := strings.TrimPrefix(strings.TrimSuffix(section, "]"), "[")
			fmt.Println(section)
			parser.sections[section] = make(map[string]string)
			parser.sectionsNameList = append(parser.sectionsNameList, string(section[1:]))
			parser.numberOfSections++
		} else {
			keysAndValues := strings.Split(section, "\n")
			for _, pair := range keysAndValues {
				if string(pair) != "#" {
					pair := strings.Split(pair, "=")
					sec, _, _ := parser.Set(parser.sectionsNameList[parser.numberOfSections-1], pair[0], pair[1])
					if sec != parser.sectionsNameList[parser.numberOfSections-1] {
						errCouldNotParseData := errors.New("Couldn't parse value ")
						return errCouldNotParseData
					}
				}
			}

		}
	}
	return nil
}

func (parser *IniParser) LoadFromFile(path string) error {
	readFile, err := os.ReadFile(path)
	if err != nil {
		return parser.errCouldNotLoadDataFromFile
	}
	errParser := parser.LoadFromString(string(readFile))
	if errParser != nil {
		errCouldNotParseData := errors.New("Couldn't parse value")
		return errCouldNotParseData
	}
	return nil
}

func (parser *IniParser) SaveToFile() error {
	sections := maps.Keys(parser.sections)
	if sections == nil {
		errNoSections := errors.New("No avaliable sections to save in file")
		return errNoSections
	}
	str := ""
	for _, section := range sections {
		str += "[" + section + "]\n"
		keys := maps.Values(parser.sections[section])
		for _, key := range keys {
			str += key + "=" + parser.sections[section][key]
		}

	}
	file, err := os.Create(`../../config.ini`)
	if err != nil {
		return err
	}
	_, errWrite := file.WriteString(str)
	if errWrite != nil {
		return err
	}
	return nil
}
