package parser

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"slices"
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
		{"Test1: Set in new section", "section1", "key1", "value1", "value1"},
		{"Test2: Set new key value in existing section", "section1", "key1-2", "value2", "value2"},
		{"Test3: Set empty value", "section1", "key1-3", "", ""},
		{"Test4: Set a value of existing key", "section1", "key1-3", "value4", "value4"},
		{"Test5: Set value of uppercase", "section1", "key1-3", "VALUE4", "value4"},
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
	errNotFound := errors.New("The given data is not valid")
	var tests = []struct {
		name         string
		inputSection string
		inputKey     string
		inputValue   string
		wantValue1   string
		wantValue2   error
	}{
		{"Test1: Get value of key in specfic section", "section1", "key1", "value1", "value1", nil},
		{"Test2: Get value of section with square bracket mistakely", "[section1]", "key1-2", "value2", "value2", nil},
		{"Test3: Get empty value", "section1", "key1-3", "", "", nil},
		{"Test4: Get value of uppercase", "section1", "key1-3", "VALUE4", "value4", nil},
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
	t.Run("Test5: Get value of non saved section", func(t *testing.T) {
		value, err := p.Get("section2", "key21")
		if value != "" && err != errNotFound {
			t.Errorf("Get---> got %s %v want %s %v", value, err, "", errNotFound)
		}
	})

}

func TestGetSectionNames(t *testing.T) {
	want := []string{"section1", "section2", "section3", "section4"}
	p := NewParser()
	for _, section := range want {
		p.sections[section] = make(map[string]string)
	}
	t.Run("Test: Get normal section name", func(t *testing.T) {
		got := p.GetSectionNames()
		fmt.Println(got)
		if !slices.Equal(want, got) {
			t.Errorf("GetSectionNames---> got %v want %v", got, want)
		}
	})
}
func TestGetSections(t *testing.T) {
	p := NewParser()
	want := map[string]map[string]string{
		"section1": {
			"key1": "value1", "key2": "value2",
		},
		"section2": {
			"key2-1": "value2-1", "key2-2": "value2-2",
		},
	}
	p.sections = want
	t.Run("Test: Get map sections", func(t *testing.T) {
		got, _ := p.GetSections()
		if !reflect.DeepEqual(got, want) {
			t.Errorf("GetSection--->got %v want %v", got, want)
		}
	})
}

func TestLoadFromString(t *testing.T) {
	p := NewParser()
	str := `
	[section1]
	key1=value1
	key2=value2
	#this is a comment
	[section2]
	key2-1=value2-1
	key2-2=value2-2
	[SEction1]
	keY1=value1-1
	`
	want := map[string]map[string]string{
		"section1": {
			"key1": "value1-1",
			"key2": "value2",
		},
		"section2": {
			"key2-1": "value2-1", "key2-2": "value2-2",
		},
	}
	t.Run("Test: Load data from string having section written twice to update value of key", func(t *testing.T) {
		err := p.LoadFromString(str)
		if err != nil && !reflect.DeepEqual(p.sections, want) {
			t.Errorf("LoadFromString---> got %v want %v", p.sections, want)
		}
	})
}

func TestLoadFromFile(t *testing.T) {
	p := NewParser()
	paths, err := filepath.Glob(filepath.Join("testdata", "*.input"))
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]map[string]string{
		"default": {"serveraliveinterval": "45", "compression": "yes",
			"compressionlevel": "9", "forwardx11": "yes",
		}, "forge.example": {
			"user": "hg",
		}, "topsecret.server.example": {
			"port": "50022", "forwardx11": "no",
		},
	}

	_, filename := filepath.Split(paths[0])
	testname := filename[:len(filename)-len(filepath.Ext(paths[0]))]
	t.Run(testname, func(t *testing.T) {
		err := p.LoadFromFile(paths[0])
		if err != nil && !reflect.DeepEqual(p.sections, want) {
			t.Errorf("LoadFromFile---> got %v want %v", p.sections, want)
		}
	})

}
func TestToString(t *testing.T) {
	p := NewParser()
	p.sections = map[string]map[string]string{
		"default": {"serveraliveinterval": "45", "compression": "yes",
			"compressionlevel": "9", "forwardx11": "yes",
		}, "forge.example": {
			"user": "hg",
		}, "topsecret.server.example": {
			"port": "50022", "forwardx11": "no",
		},
	}
	goldenfile := filepath.Join("testdata", "save_to_file"+".golden")
	bytes, err := os.ReadFile(goldenfile)
	if err != nil {
		t.Fatal("error reading golden file:", err)
	}
	want := string(bytes)
	t.Run("Test: Convert data to string", func(t *testing.T) {
		got := p.ToString()
		if reflect.DeepEqual(got, want) {
			t.Errorf("ToString---> got %s want %s", got, want)
		}
	})
}

func TestSaveToFile(t *testing.T) {
	p := NewParser()
	p.sections = map[string]map[string]string{
		"default": {"serveraliveinterval": "45", "compression": "yes",
			"compressionlevel": "9", "forwardx11": "yes",
		}, "forge.example": {
			"user": "hg",
		}, "topsecret.server.example": {
			"port": "50022", "forwardx11": "no",
		},
	}
	goldenfile := filepath.Join("testdata", "save_to_file"+".golden")
	want, err := os.ReadFile(goldenfile)
	if err != nil {
		t.Fatal("error reading golden file:", err)
	}

	t.Run("Test: Save to file", func(t *testing.T) {
		err := p.SaveToFile()
		got, _ := os.ReadFile(`../config.ini`)
		if err != nil && !bytes.Equal(want, got) {
			fmt.Println(err)
			t.Errorf("SaveToFile---> got %v want %v", got, want)
		}
	})

}
