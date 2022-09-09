package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
)

func TestOpenConfigurationOk(t *testing.T) {
	r, err := openConfiguration(ConfigurationFilepath)
	defer r.Close()
	assert.Nil(t, err, "no error should be returned when file exists. got %s", err)
	assert.NotNilf(t, r, "file descriptor should not be nil")
}

func TestOpenConfigurationKo(t *testing.T) {
	r, err := openConfiguration("")
	defer r.Close()
	assert.Nil(t, r, "no file descriptor should be return if there is an os error")
	assert.ErrorIs(t, err, fs.ErrNotExist)
}

func TestUnmarshalConfigurationOk(t *testing.T) {
	r, err := openConfiguration(ConfigurationFilepath)
	defer r.Close()

	platformVersions, err := unmarshalConfiguration(r)
	assert.Nil(t, err, "no error should be returned when configuration reading succeed. got %s", err)
	assert.NotEmptyf(t, platformVersions, "struct should be filled with informations.")

	assert.NotEmpty(t, platformVersions.Poetry120.PythonVersions)
	assert.NotEmpty(t, platformVersions.Poetry120.ImageVariants)
	assert.Equal(t, platformVersions.Poetry120.Version, "1.2.0")
}

func TestUnmarshalConfigurationJSONSyntaxError(t *testing.T) {
	_, err := unmarshalConfiguration(strings.NewReader(""))
	var te *json.SyntaxError
	assert.ErrorAs(t, err, &te)
}

func TestPlatforms_GetPlatforms(t *testing.T) {
	platforms := Platforms{
		Poetry120:  Platform{},
		Poetry1115: Platform{},
	}
	assert.Equal(t, 2, len(platforms.GetPlatforms()))
}

func TestPythonImageVariant_GetDockerfileNotation(t *testing.T) {
	pair := pythonImageVariant{
		PythonVersion: "a.b12a.c1235",
		ImageVariant:  "anNameWi123=-0123",
	}
	r := pair.GetDockerfileNotation()
	assert.Equal(t, "a.b12a.c1235-anNameWi123=-0123", r)
}

// Not the best way to test it, but it's enough I would say.
func TestBasepathPointsToRootOfProject(t *testing.T) {
	assert.FileExists(t, filepath.Join(basepath, "go.mod"))
	assert.FileExists(t, filepath.Join(basepath, "main.go"))
}

func TestGetWritingPathFromGeneratesFilepath(t *testing.T) {
	platform := Platform{
		Version: "1.0.0",
	}
	pair := pythonImageVariant{
		PythonVersion: "3.11.0-rc1",
		ImageVariant:  "bullseye",
	}
	r := getWritingPathFrom(platform, pair)
	assert.Equal(t, filepath.Join(basepath, platform.Version, pair.PythonVersion, pair.ImageVariant, "Dockerfile"), r)
}

func TestGetImagesFromPythonVersionsAndImageVariants(t *testing.T) {
	pythonVersions := []string{"3.10.0", "3.9.0"}
	imageVariants := []string{"bullseye", "slim-buster"}

	r := getImageNamesFrom(pythonVersions, imageVariants)

	assert.Equal(t, 4, len(r))
	assert.Equal(t, "3.10.0-bullseye", r[0].GetDockerfileNotation())
	assert.Equal(t, "3.10.0-slim-buster", r[1].GetDockerfileNotation())
	assert.Equal(t, "3.9.0-bullseye", r[2].GetDockerfileNotation())
	assert.Equal(t, "3.9.0-slim-buster", r[3].GetDockerfileNotation())
}
