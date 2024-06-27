package parser

import (
	// "bytes"
	// "go/format"
	// "os"
	// "path/filepath"
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestNewParser(t *testing.T) {
	var want map[string]map[string]string
	res := NewParser()
	if reflect.DeepEqual(res, want) {
		t.Errorf("NewParser---> got %v want %v", res.sections, want)
	}
}

func TestSetTableDriven(t *testing.T) {
	p := NewParser()
	var tests = []struct {
		Name         string
		InputSection string
		InputKey     string
		InputValue   string
		wantValue    string
	}{
		{"test1: Set in new section", "section1", "key1", "value1", "value1"},
		{"test2: Set new key value in existing section", "section1", "key1-2", "value2", "value2"},
		{"test3: Set empty value", "section1", "key1-3", "", ""},
		{"test4: Set a value of existing key", "section1", "key1-3", "value4", "value4"},
		{"test5: Set value of uppercase", "section1", "key1-3", "VALUE4", "value4"},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			p.Set(test.InputSection, test.InputKey, test.InputValue)
			ans := p.sections[test.InputSection][test.InputKey]
			if ans != test.wantValue {
				t.Errorf("Set--->got %s, want %s", ans, test.wantValue)
			}
		})
	}
}
func TestGetTableDriven(t *testing.T) {
	p := NewParser()
	var errNotFound = errors.New("The given data is not valid")
	var tests = []struct {
		name         string
		inputSection string
		inputKey     string
		inputValue   string
		wantValue1   string
		wantValue2   error
	}{
		{"test1: Get value of key in specfic section", "section1", "key1", "value1", "value1", nil},
		{"test2: Get value of section with square bracket mistakely", "[section1]", "key1-2", "value2", "value2", nil},
		{"test3: Get empty value", "section1", "key1-3", "", "", nil},
		{"test4: Get value of uppercase", "section1", "key1-3", "VALUE4", "value4", nil},
		//{"test5: Get value of non saved section", "section2", "key21", "","",errNotFound},
	}
	for _, test := range tests {
		test.inputSection = strings.TrimPrefix(strings.TrimSuffix(test.inputSection, "]"), "[")
		p.sections[test.inputSection] = make(map[string]string)
		p.sections[test.inputSection][test.inputKey] = test.inputValue
		t.Run(test.name, func(t *testing.T) {
			value, err := p.Get(test.inputSection, test.inputKey)
			if value != test.wantValue1 && err != test.wantValue2 {
				t.Errorf("Get---> got %s %v want %s %v", value, err, test.wantValue1, test.wantValue2)
			}
		})
	}
	t.Run("test5: Get value of non saved section", func(t *testing.T) {
		value, err := p.Get("section2", "key21")
		if value != "" && err != errNotFound {
			t.Errorf("Get---> got %s %v want %s %v", value, err, "", errNotFound)
		}
	})

}

// func TestLoadFromFile(t *testing.T) {
// 	p := NewParser()
// 	paths, err := filepath.Glob(filepath.Join("testdata", "*.input"))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	for _, path := range paths {
// 		_, filename := filepath.Split(path)
// 		testname := filename[:len(filename)-len(filepath.Ext(path))]

// 		t.Run(testname, func(t *testing.T) {
// 			p.LoadFromFile(path)
// 	}

// }
