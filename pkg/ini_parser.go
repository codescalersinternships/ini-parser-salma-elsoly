package parser

import (
	"errors"
	"sort"

	"os"
	"regexp"
	"strings"

	"golang.org/x/exp/maps"
)

type Stringer interface {
	String()
}

/*This struct is used to identify Iniparser type which contains methods for parsing and creating ini file getting and setting values*/
type IniParser struct {
	sections map[string]map[string]string
}

/* This acts as constructor for IniParser struct and returns Iniparser*/
func NewParser() IniParser {
	newSections := make(map[string]map[string]string)
	newParsing := IniParser{}
	newParsing.sections = newSections
	return newParsing
}

/*
This function set the given value of given section and key returns a string which is the value set if section doesn't exist

it creates it and set the value to key
*/
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

/* This function take section name and key and retruns the value corresponding to them. returns an error if section or key is not found*/
func (parser *IniParser) Get(section_name, key string) (string, error) {
	section_name = strings.TrimPrefix(strings.TrimSuffix(section_name, "]"), "[")
	if value, ok := parser.sections[strings.ToLower(section_name)][strings.ToLower(key)]; ok {
		return value, nil
	}
	return "", errors.New("The given data is not valid")
}

/*This function returns list of section names*/
func (parser *IniParser) GetSectionNames() []string {
	list := maps.Keys(parser.sections)
	sort.Strings(list)
	return list
}

/*This function returns the sections of ini in map[string]map[string]string and error if no sections was stored*/
func (parser *IniParser) GetSections() (map[string]map[string]string, error) {
	if parser.sections == nil {
		return nil, errors.New("No avaliable sections to return")
	}
	return parser.sections, nil
}

/*This function take string as parameter and parse it and store the  data returns an error if data couldn't be parsed*/
func (parser *IniParser) LoadFromString(str string) error {
	var currSection string
	str = strings.TrimSpace(str)
	str = strings.ToLower(str)
	sectionRegex, _ := regexp.Compile(`^\[.+\]$`)
	sectionsNames := strings.Split(str, "\n")
	for _, slice := range sectionsNames {
		slice = strings.TrimSpace(slice)
		if sectionRegex.Match([]byte(slice)) {
			section := strings.TrimPrefix(strings.TrimSuffix(slice, "]"), "[")
			currSection = section
			parser.sections[section] = make(map[string]string)
		} else if !strings.HasPrefix(string(slice), "#") {
			if !strings.Contains(slice, "=") {
				continue
			}
			pair := strings.Split(slice, "=")
			sec := parser.Set(currSection, pair[0], pair[1])
			if sec != pair[1] {
				return errors.New("Couldn't parse value ")
			}

		}
	}
	return nil
}

/*This function take path of file to load data from it returns an error file was not able to be read or couldn't parse the file*/
func (parser *IniParser) LoadFromFile(path string) error {
	if !(strings.HasSuffix(path, ".ini") || strings.HasSuffix(path, ".input")) {
		return errors.New("Not a INI file or input file")
	}
	readFile, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	errParser := parser.LoadFromString(string(readFile))
	if errParser != nil {
		return errParser
	}
	return nil
}

/* This function convert sections stored of .ini to string and return this string*/
func (parser *IniParser) String() string {
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

/*
The function doesn't take arguments, it opens file config.ini in the current working directory default or in the specfied directory

if the file doesn't exist it creates it and write to it, returns error if issue occured
*/
func (parser *IniParser) SaveToFile(path ...string) error {
	str := parser.String()
	if str == "" {
		errNoSections := errors.New("No avaliable sections to save in file")
		return errNoSections
	}
	currDire, _ := os.Getwd()
	var filePath = currDire + `/config.ini`
	if len(path) == 1 {
		filePath = path[0] + "config.ini"
	}
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	_, errWrite := file.WriteString(str)
	if errWrite != nil {
		return err
	}
	
	return nil
}
