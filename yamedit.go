package yamedit

import (
	"fmt"
	"log"
	"os/exec"
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

func Edit(YAML []byte, path []string, newValue string) []byte {
	var data map[string]interface{}

	fmt.Printf("%s\n", YAML)
	if err := yaml.Unmarshal(YAML, &data); err != nil {
		return YAML
	}

	fmt.Println(data)

	pointer := data

	for _, s := range path {
		if d, ok := pointer[s].(map[string]interface{}); ok {
			pointer = d
		}
	}

	pointer[path[len(path)-1]] = newValue

	b, err := yaml.Marshal(data)

	if err != nil {
		return YAML
	}

	return b
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
