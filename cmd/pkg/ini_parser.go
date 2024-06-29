package parser

import (
	"errors"

	"slices"

	"os"
	"regexp"
	"strings"

	"golang.org/x/exp/maps"
)

type IniParser struct {
	sections         map[string]map[string]string
	sectionsNameList []string
}

func NewParser() IniParser {
	newSections := make(map[string]map[string]string)
	newParsing := IniParser{}
	newParsing.sections = newSections
	return newParsing
}

func (parser *IniParser) Set(section, key, value string) string {
	section = strings.TrimPrefix(strings.TrimSuffix(section, "]"), "[")
	section = strings.ToLower(section)
	key = strings.ToLower(key)
	value = strings.ToLower(value)
	if parser.sections[section] == nil {
		parser.sections[section] = make(map[string]string)
	}
	parser.sections[section][key] = value
	return parser.sections[section][key]
}

func (parser *IniParser) Get(section_name, key string) (string, error) {
	section_name = strings.TrimPrefix(strings.TrimSuffix(section_name, "]"), "[")
	if value, ok := parser.sections[strings.ToLower(section_name)][strings.ToLower(key)]; ok {
		return value, nil
	} else {
		errNotFound := errors.New("The given data is not valid")
		return value, errNotFound
	}
}

func (parser *IniParser) GetSectionNames() []string {
	if len(parser.sectionsNameList) != len(parser.sections) {
		parser.sectionsNameList = nil
		parser.sectionsNameList = append(parser.sectionsNameList, maps.Keys(parser.sections)...)

	}
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
	var sectionIndex int
	str = strings.TrimSpace(str)
	str = strings.ToLower(str)
	sectionRegex, _ := regexp.Compile(`^\[.+\]$`)
	sectionsNames := strings.Split(str, "\n")
	for _, slice := range sectionsNames {
		slice = strings.TrimSpace(slice)
		if sectionRegex.Match([]byte(slice)) {
			section := strings.TrimPrefix(strings.TrimSuffix(slice, "]"), "[")
			sectionIndex = slices.Index(parser.sectionsNameList, section)
			if sectionIndex != -1 {
				continue
			}
			parser.sections[section] = make(map[string]string)
			parser.sectionsNameList = append(parser.sectionsNameList, string(section))
		} else if !strings.HasPrefix(string(slice), "#") {
			if !strings.Contains(slice, "=") {
				continue
			}
			pair := strings.Split(slice, "=")
			if sectionIndex == -1 {
				sectionIndex = len(parser.sectionsNameList) - 1
			}
			sec := parser.Set(parser.sectionsNameList[sectionIndex], pair[0], pair[1])
			if sec != pair[1] {
				errCouldNotParseData := errors.New("Couldn't parse value ")
				return errCouldNotParseData
			}

		}
	}
	return nil
}

func (parser *IniParser) LoadFromFile(path string) error {
	readFile, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	errParser := parser.LoadFromString(string(readFile))
	if errParser != nil {
		errCouldNotParseData := errors.New("Couldn't parse value")
		return errCouldNotParseData
	}
	return nil
}

func (parser *IniParser) ToString() string {
	var str = ""
	sections := maps.Keys(parser.sections)
	for _, section := range sections {
		str += "[" + section + "]\n"
		keys := maps.Keys(parser.sections[section])
		values := maps.Values(parser.sections[section])
		for i, key := range keys {
			str += key + "=" + values[i] + "\n"
		}
		str += "\n"

	}
	return str
}

func (parser *IniParser) SaveToFile() error {
	str := parser.ToString()
	if str == "" {
		errNoSections := errors.New("No avaliable sections to save in file")
		return errNoSections
	}
	file, err := os.OpenFile(`../../config.ini`, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	_, errWrite := file.WriteString(str)
	if errWrite != nil {
		return err
	}
	defer file.Close()
	return nil
}
