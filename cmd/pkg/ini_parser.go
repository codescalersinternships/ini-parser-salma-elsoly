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
	sectionL:=strings.ToLower(section)
	keyL:=strings.ToLower(key)
	if parser.sections[sectionL] == nil {
		parser.sections[sectionL] = make(map[string]string)
	}
	parser.sections[sectionL][keyL] = value
	return sectionL, keyL, value
}

func (parser *IniParser) Get(section_name, key string) (string, error) {
	if value, ok := parser.sections[strings.ToLower(section_name)][strings.ToLower(key)]; ok {
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
	strTrim:= strings.TrimSpace(str)
	strLower := strings.ToLower(strTrim)
	sectionRegex, _ := regexp.Compile(`^\[.+\]$`)
	sectionsNames := strings.Split(strLower, "\n")
	for _, section := range sectionsNames {
		if sectionRegex.Match([]byte(section)) {
			section := strings.TrimPrefix(strings.TrimSuffix(section, "]"), "[")
			fmt.Println(section)
			parser.sections[section] = make(map[string]string)
			parser.sectionsNameList = append(parser.sectionsNameList, string(section))
			parser.numberOfSections++
		} else if strings.Contains(section,"="){
			keysAndValues := strings.Split(section, "\n")
			for _, pair := range keysAndValues {
				if !strings.Contains(string(pair),"#"){
					fmt.Println(pair)
					pair := strings.Split(pair, "=")
					fmt.Println("---------------->",pair)
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

func (parser *IniParser) ToString() string{
	sections := parser.sectionsNameList
	if sections == nil {
		return ""
	}
	str := ""
	for _, section := range sections {
		str += "[" + section + "]\n"
		keys := maps.Keys(parser.sections[section])
		values:=maps.Values(parser.sections[section])
		for i, key:= range (keys){
			str+=key+"="+values[i]+"\n"
		}

	}
	return str
}

func (parser *IniParser) SaveToFile() error {
	str:=parser.ToString()
	if str==""{
		errNoSections := errors.New("No avaliable sections to save in file")
		return errNoSections
	}
	file, err := os.Create(`../config.ini`)
	if err != nil {
		return err
	}
	_, errWrite := file.WriteString(str)
	if errWrite != nil {
		return err
	}
	return nil
}
