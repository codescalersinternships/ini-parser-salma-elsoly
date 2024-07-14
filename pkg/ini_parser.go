package parser

import (
	"errors"
	"sort"

	"os"
	"strings"

	// "golang.org/x/exp/maps"
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
	if value, ok := parser.sections[strings.ToLower(section_name)][strings.ToLower(key)]; ok {
		return strings.ToLower(value), nil
	}
	return "", errors.New("The given data is not valid")
}

/*This function returns list of section names*/
func (parser *IniParser) GetSectionNames() []string {
	var list []string
	for section:= range(parser.sections){
		list = append(list, section)
	}
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
	lines := strings.Split(str, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line,"[") && strings.HasSuffix(line,"]") {
			section := strings.TrimPrefix(strings.TrimSuffix(line, "]"), "[")
			currSection = section
			if parser.sections[section] == nil {
				parser.sections[section] = make(map[string]string)
			}
		} else if !strings.HasPrefix(string(line), "#") {
			if !strings.Contains(line, "=") {
				continue
			}
			pair := strings.Split(line, "=")
			sec := parser.Set(currSection, strings.TrimSpace(pair[0]), strings.TrimSpace(pair[1]))
			if sec != strings.TrimSpace(pair[1]) {
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
	// sections := maps.Keys(parser.sections)
	// sort.Strings(sections)
	for section := range(parser.sections) {
		str += "[" + section + "]\n"
		// keys := maps.Keys(parser.sections[section])
		// sort.Strings(keys)
		for key := range (parser.sections[section]) {
			str += key + "=" + parser.sections[section][key] + "\n"
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
