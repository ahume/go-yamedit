package yamedit

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/qri-io/jsonpointer"
	"gopkg.in/yaml.v3"
)

// Get the string value from the YAML from the given path
func Get(YAML []byte, path string) (string, error) {
	var data map[string]interface{}
	if err := yaml.Unmarshal(YAML, &data); err != nil {
		return "", err
	}

	ptr, err := jsonpointer.Parse(path)
	if err != nil {
		return "", err
	}
	r, err := ptr.Eval(data)
	if err != nil {
		return "", err
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
	return result, nil
}

// Edit the string value in the YAML from the given path
func Edit(YAML []byte, path, newValue string) ([]byte, error) {
	currentValue, err := Get(YAML, path)
	if err != nil {
		return []byte{}, err
	}

	splitPath := strings.Split(path, "/")
	finalProp := unescapeJSONPointerString(splitPath[len(splitPath)-1])
	regexEscapedFinalProp := escapePathStringForRegex(finalProp)

	matchToken := "matchToken" + getUUID()
	matchIndex := 0

	reMatchAllKeyValuePairs, _ := regexp.Compile(regexEscapedFinalProp + " *: *\"?" + currentValue + "\"?")
	reMatchAllKeyTokenPairs, _ := regexp.Compile(regexEscapedFinalProp + " *: *" + matchToken + "\\d+")

	// BUG: This is a shortcut with a significant bug. We detect if the value we're
	// looking for is a value member of an array, and if it is, we swap it directly
	// and return. However, because this is a global replace, if there are any
	// other arrays which have the same value in the same position, they will also
	// be changed and returned in the updated YAML string.
	if targetIsArrayMember(splitPath) {
		r, _ := regexp.Compile("- *" + currentValue)
		return r.ReplaceAll(YAML, []byte("- "+newValue)), nil
	}

	tokenisedYAML := reMatchAllKeyValuePairs.ReplaceAllFunc(YAML, func(s []byte) []byte {
		matchIndex = matchIndex + 1
		return []byte(finalProp + ": " + matchToken + strconv.Itoa(matchIndex))
	})

	g, _ := Get(tokenisedYAML, path)
	split := strings.Fields(finalProp + " : .? " + g)
	newS := strings.Join(split, "\\s*")
	reMatchTargetKeyToken, _ := regexp.Compile(newS)

	// Switch matching token to the new value
	YAMLWithNewValue := reMatchTargetKeyToken.ReplaceAll(tokenisedYAML, []byte(finalProp+": "+newValue))

	// Switch all remaining matchTokens back to their original values
	YAMLWithNewValue = reMatchAllKeyTokenPairs.ReplaceAll(YAMLWithNewValue, []byte(finalProp+": "+currentValue))

	return YAMLWithNewValue, nil
}

func getUUID() string {
	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		panic("Da fuq")
	}
	return strings.TrimSpace(string(out))
}

func targetIsArrayMember(path []string) bool {
	lastValue := path[len(path)-1]
	if _, err := strconv.Atoi(lastValue); err == nil {
		return true
	}
	return false
}

// escapePathStringForRegex ensures that all regex metacharacters are escaped.
// This allows for keys with regex characters e.g. "a.b.c/hello" to sill be
// found for tokenisation.
func escapePathStringForRegex(value string) string {
	return regexp.QuoteMeta(value)
}

// unescapePointerString removes the special case escapeing used in
// the Pointer spec. https://stackoverflow.com/questions/31483170/purpose-of-tilde-in-json-pointer
func unescapeJSONPointerString(value string) string {
	newValue := value
	newValue = strings.ReplaceAll(newValue, "~0", "~")
	newValue = strings.ReplaceAll(newValue, "~1", "/")
	return newValue
}
