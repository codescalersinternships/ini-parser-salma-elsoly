package parser

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

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
			assertEqual(t, ans, test.wantValue)
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
			assertEqual(t, value, test.wantValue1)
			assertError(t, err, test.wantValue2)
		})
	}
	t.Run("Test5: Get value of non saved section", func(t *testing.T) {
		value, err := p.Get("section2", "key21")
		assertEqual(t, value, "")
		assertError(t, err, errNotFound)
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
		assertEqual(t, got, want)
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
		assertEqual(t, got, want)
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
		assertEqual(t, p.sections, want)
		assertError(t, err, nil)
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
		assertEqual(t, p.sections, want)
		assertError(t, err, nil)

	})

}
func TestString(t *testing.T) {
	p := NewParser()
	p.sections = map[string]map[string]string{
		"default": {"compression": "yes",
			"compressionlevel": "9", "forwardx11": "yes", "serveraliveinterval": "45",
		}, "forge.example": {
			"user": "hg",
		}, "topsecret.server.example": {
			"forwardx11": "no", "port": "50022",
		},
	}
	want := `[default]
compression=yes
compressionlevel=9
forwardx11=yes
serveraliveinterval=45

[forge.example]
user=hg

[topsecret.server.example]
forwardx11=no
port=50022

`
	t.Run("Test: Convert data to string", func(t *testing.T) {
		got := p.String()
		assertEqual(t, got, want)
	})
}

func TestSaveToFile(t *testing.T) {
	p := NewParser()
	p.sections = map[string]map[string]string{
		"default": {"compression": "yes",
			"compressionlevel": "9", "forwardx11": "yes", "serveraliveinterval": "45",
		}, "forge.example": {
			"user": "hg",
		}, "topsecret.server.example": {
			"forwardx11": "no", "port": "50022",
		},
	}
	currDirec, _ := os.Getwd()
	var tests = []struct {
		name            string
		inputPath       string
		outputPathCheck string
	}{
		{"Test: Save to file in default directory", currDirec+`/testdata/`, currDirec + `/testdata/config.ini`},
		//{"Test: Save to file in specfied directory", `/home/salmaelsoly/Codescalers-internship/`, `/home/salmaelsoly/Codescalers-internship/`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var err error
			if test.inputPath == "" {
				err = p.SaveToFile()
			} else {
				err = p.SaveToFile(test.inputPath)
			}
			if err != nil {
				t.Fatal(err)
			}
			got, err := os.Stat(test.outputPathCheck)
			if err != nil {
				t.Fatal(err)
			}

			assertEqual(t, got.ModTime().Format(time.UnixDate), time.Now().Format(time.UnixDate))
			assertError(t, err, nil)
		})
	}

}
func assertEqual(t *testing.T, got, want any) {
	t.Helper()
	if reflect.TypeOf(got) == reflect.TypeOf(" ") {
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	} else if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func assertError(t *testing.T, got, want error) {
	t.Helper()
	if !strings.Contains(fmt.Sprint(got), fmt.Sprint(want)) {
		t.Errorf("got %v want %v", got, want)
	}
}
