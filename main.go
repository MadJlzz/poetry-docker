package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const ConfigurationFilepath = "configuration.json"

// PlatformVersions struct which contains
// an array of platforms.
type PlatformVersions struct {
	Poetry12 Platform `json:"1.2"`
	Poetry11 Platform `json:"1.1"`
}

// Platform struct which contains a name
// a type and a list of social links
type Platform struct {
	PythonVersions []string `json:"pythons"`
	ImageVariants  []string `json:"variants"`
	Version        string   `json:"version"`
}

func openConfiguration(filepath string) (*os.File, error) {
	jsonFile, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	return jsonFile, nil
}

func unmarshalConfiguration(reader io.Reader) (*PlatformVersions, error) {
	// read our opened jsonFile as a byte array.
	byteValue, _ := io.ReadAll(reader)

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'platformVersions' which we defined above
	var platformVersions PlatformVersions
	err := json.Unmarshal(byteValue, &platformVersions)
	if err != nil {
		return nil, err
	}

	return &platformVersions, nil
}

func main() {
	confFile, err := openConfiguration(ConfigurationFilepath)
	defer confFile.Close()
	if err != nil {
		panic(err)
	}

	platformVersions, err := unmarshalConfiguration(confFile)
	if err != nil {
		panic(err)
	}
	fmt.Println(platformVersions)
}
