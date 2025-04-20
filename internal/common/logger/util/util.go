package util

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/text"
	"gopkg.in/yaml.v3"
)

func PrettyPrint(v interface{}) {
	yamlBytes, _ := yaml.Marshal(v)
	output := text.Colors{text.FgGreen}.Sprintf(string(yamlBytes))
	fmt.Println(output)
}

func PrettyYaml(v interface{}) string {
	yamlBytes, _ := yaml.Marshal(v)
	return text.Colors{text.FgGreen}.Sprintf("%s", string(yamlBytes))
}
