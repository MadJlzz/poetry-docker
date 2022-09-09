package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"text/template"
)

const (
	ConfigurationFilepath  = "configuration.json"
	DockerTemplateFilepath = "Dockerfile.gotmpl"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

// Platforms struct which contains
// an array of platforms.
type Platforms struct {
	Poetry120  Platform `json:"1.2.0"`
	Poetry1115 Platform `json:"1.1.15"`
}

func (p *Platforms) GetPlatforms() []Platform {
	return []Platform{
		p.Poetry120,
		p.Poetry1115,
	}
}

// Platform struct which contains a name
// a type and a list of social links
type Platform struct {
	PythonVersions []string `json:"pythons"`
	ImageVariants  []string `json:"variants"`
	Version        string   `json:"version"`
}

type pythonImageVariant struct {
	PythonVersion string
	ImageVariant  string
}

func (piv *pythonImageVariant) GetDockerfileNotation() string {
	return fmt.Sprintf("%s-%s", piv.PythonVersion, piv.ImageVariant)
}

func openConfiguration(filepath string) (*os.File, error) {
	jsonFile, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	return jsonFile, nil
}

func unmarshalConfiguration(reader io.Reader) (*Platforms, error) {
	// read our opened jsonFile as a byte array.
	byteValue, _ := io.ReadAll(reader)

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'platforms' which we defined above
	var platforms Platforms
	err := json.Unmarshal(byteValue, &platforms)
	if err != nil {
		return nil, err
	}

	return &platforms, nil
}

func getImageNamesFrom(pythonVersions []string, imageVariants []string) []pythonImageVariant {
	pairs := make([]pythonImageVariant, 0, len(pythonVersions)*len(imageVariants))
	for _, pv := range pythonVersions {
		for _, iv := range imageVariants {
			pair := pythonImageVariant{
				PythonVersion: pv,
				ImageVariant:  iv,
			}
			pairs = append(pairs, pair)
		}
	}
	return pairs
}

func getWritingPathFrom(platform Platform, pair pythonImageVariant) string {
	return fmt.Sprintf("%s/%s/%s/%s/Dockerfile", basepath, platform.Version, pair.PythonVersion, pair.ImageVariant)
}

func generateDockerfilesFrom(tmpl *template.Template, platform Platform, pairs []pythonImageVariant) {
	for _, image := range pairs {
		path := getWritingPathFrom(platform, image)
		f, err := os.Create(path)
		if err != nil {
			fmt.Printf("an error occured while trying to create Dockerfile for platform [%s] and image [%s]",
				platform.Version, image)
		}
		err = tmpl.Execute(f, map[string]string{"FromVersion": image.GetDockerfileNotation(), "PoetryVersion": platform.Version, "ImageVariant": image.ImageVariant})
		if err != nil {
			fmt.Printf("an error occured while trying to generate Dockerfile content for platform [%s] and image [%s].\ngot: %s\n",
				platform.Version, image, err)
		}
		if err = f.Close(); err != nil {
			fmt.Printf("couldn't close file for platform [%s] and image [%s]",
				platform.Version, image)
		}
	}
}

func generateDockerfilesFor(tmpl *template.Template, platforms *Platforms) {
	for _, platform := range platforms.GetPlatforms() {
		images := getImageNamesFrom(platform.PythonVersions, platform.ImageVariants)
		generateDockerfilesFrom(tmpl, platform, images)
	}
}

func main() {
	confFile, err := openConfiguration(ConfigurationFilepath)
	defer confFile.Close()
	if err != nil {
		panic(err)
	}

	platforms, err := unmarshalConfiguration(confFile)
	if err != nil {
		panic(err)
	}

	tmpl, err := template.New(DockerTemplateFilepath).ParseFiles(DockerTemplateFilepath)
	if err != nil {
		panic(err)
	}
	generateDockerfilesFor(tmpl, platforms)
}
