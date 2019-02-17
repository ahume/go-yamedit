package yamedit

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/jmoiron/jsonq"
)

// Get the string value from the YAML from the given path
func Get(YAML []byte, path []string) string {
	var data map[string]interface{}
	yaml.Unmarshal(YAML, &data)

	jq := jsonq.NewQuery(data)

	r, err := jq.Interface(path...)
	if err != nil {
		log.Fatal("The requested path cannot be found.", err)
	}

	var result string
	switch v := r.(type) {
	case int:
		result = strconv.Itoa(r.(int))
	case float64:
		result = fmt.Sprintf("%g", r)
	case string:
		result = r.(string)
	default:
		log.Fatalf("Unexpected type found at path: %s", v)
	}
	return result
}

// Edit the string value in the YAML from the given path
func Edit(YAML []byte, path []string, newValue string) []byte {
	currentValue := Get(YAML, path)
	finalProp := path[len(path)-1]

	matchToken := "matchToken" + getUUID()
	matchIndex := 0

	reMatchAllKeyValuePairs, _ := regexp.Compile(finalProp + " *: *\"?" + currentValue + "\"?")
	reMatchAllKeyTokenPairs, _ := regexp.Compile(finalProp + " *: *" + matchToken + "\\d+")

	// IF 2nd last converts to int it's an array
	if targetIsArrayMember(YAML, path) {
		r, _ := regexp.Compile("- *" + currentValue)
		return r.ReplaceAll(YAML, []byte("- "+newValue))
	}

	tokenisedYAML := reMatchAllKeyValuePairs.ReplaceAllFunc(YAML, func(s []byte) []byte {
		matchIndex = matchIndex + 1
		return []byte(finalProp + ": " + matchToken + strconv.Itoa(matchIndex))
	})

	reMatchTargetKeyToken, _ := regexp.Compile(finalProp + " *: *" + Get(tokenisedYAML, path))

	// Switch matching token to the new value
	YAMLWithNewValue := reMatchTargetKeyToken.ReplaceAll(tokenisedYAML, []byte(finalProp+": "+newValue))

	// Switch all remaining matchTokens back to their original values
	YAMLWithNewValue = reMatchAllKeyTokenPairs.ReplaceAll(YAMLWithNewValue, []byte(finalProp+": "+currentValue))

	return YAMLWithNewValue
}

func getUUID() string {
	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		panic("Da fuq")
	}
	return strings.TrimSpace(string(out))
}

func targetIsArrayMember(YAML []byte, path []string) bool {
	var data map[string]interface{}
	yaml.Unmarshal(YAML, &data)

	newPath := path[:len(path)-1]

	jq := jsonq.NewQuery(data)

	_, err := jq.Array(newPath...)
	if err != nil {
		return false
	}
	return true
}
