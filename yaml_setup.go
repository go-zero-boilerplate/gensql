package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type YamlSetup map[string]*YamlEntitySetup

type YamlEntitySetup struct {
	Dialect string
	Fields  []string
	Uniques [][]string
}

func (y *YamlSetup) Validate() error {
	//TODO: Validation
	return nil
}

func LoadGeneratorSetup(inFile string) (*GeneratorSetup, error) {
	yamlBytes, err := ioutil.ReadFile(inFile)
	if err != nil {
		return nil, err
	}

	orderedYamlMap := yaml.MapSlice{}
	err = yaml.Unmarshal(yamlBytes, &orderedYamlMap)
	if err != nil {
		return nil, err
	}

	y := &YamlSetup{}
	err = yaml.Unmarshal(yamlBytes, y)
	if err != nil {
		return nil, err
	}

	err = y.Validate()
	if err != nil {
		return nil, err
	}

	orderedEntityNames := []string{}
	for _, kv := range orderedYamlMap {
		orderedEntityNames = append(orderedEntityNames, kv.Key.(string))
	}

	gs := GeneratorSetupFromYamlSetup(orderedEntityNames, y)

	return gs, nil
}
