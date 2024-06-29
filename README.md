# INI Parser

## Project Description
This project implements an INI parser using GoLang, to help parsing and storing key and values and create an ini file from
the stored key and values

## How To use the project
- to use the provided package:
    ```
    --> import github.com/codescalersinternships/salmaelsoly-inigo
    ```
    first import the package in your Go file
    ```
    --> go get github.com/codescalersinternships/salmaelsoly-inigo
    ```
    then get the package by ruuning this command in cli
    ```
    p:= parser.NewParser()
    ```
    use the constructor to create object on IniParser struct
## Usage
The struct have 8 methods
1. Set(section,key,value):
    - used for setting given value to the passed section and key returns the value set
    ```
    value:=p.Set("section1","key1","value1")
    fmt.Println(value)
    ```
    ```
    value1
    ```
2. Get(section, key):
    - used for returning the value of passed section and key, return error if section doesn't exist
    ```
    value, err := p.Get("section1","key1)
    fmt.println(value)
    ```
    ```
    value1
    ```
3. GetSectionNames( ):
    - returns list of section names
    ```
    list:= p.GetSectionNames()
    fmt.Println(list)
    ```
    ```
    section1 section2
    ```
4. GetSections( ):
    - used to get the full map of ini
    ```
    iniMap:=p.GetSections()
    ```
5. LoadFromString(str):
    - used to load configuration from passed string, returns an error if parsing the string was failed
    ```
    err:=p.LoadFromString("[Default]\nkey=value")
    ```
6. LoadFromFile(path):
    - Used to load configuration from INI file, returns an error if file doesn't exist or parsing the file failed or the file is not an ini file
    ```
    err:=p.LoadFromFile("config.ini")
    ```
7. ToString( ):
    - used to convert configuration saved to string and returns it
    ```
    str:=p.ToString()
    ```
8. SaveToFile(...path):
    - used to save configuration to config.ini, if path is passed it checks the config.ini in that path, if no path passed it checks inside the default directory (current directory)
    ```
    err:=p.SaveToFile()
    ```
    
