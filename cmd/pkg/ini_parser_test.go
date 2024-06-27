package parser

import (
	// "bytes"
	// "go/format"
	// "os"
	// "path/filepath"
	"reflect"
	"testing"
)

func TestNewParser(t *testing.T) {
	want := &IniParser{
		sections: make(map[string]map[string]string),
	}
	res := NewParser()
	if reflect.DeepEqual(res, want) {
		t.Error("NewParser---> got", res.sections, "  want", want.sections)
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
		{"test1: in new section", "section1", "key1", "value1", "value1"},
		{"test2: add new key in existing section", "section1", "key1-2", "value2", "value2"},
		{"test3: add empty value", "section1", "key1-3", "", ""},
		{"test4: update a value of existing key", "section1", "key1-3", "value4", "value4"},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			ans := p.Set(test.InputSection, test.InputKey, test.InputValue)
			if ans != test.wantValue {
				t.Errorf("got %s, want %s", ans, test.wantValue)
			}
		})
	}
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
