package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/fs"
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

	assert.NotEmpty(t, platformVersions.Poetry12.PythonVersions)
	assert.NotEmpty(t, platformVersions.Poetry12.ImageVariants)
	assert.Equal(t, platformVersions.Poetry12.Version, "1.2.0")
}

func TestUnmarshalConfigurationJSONSyntaxError(t *testing.T) {
	_, err := unmarshalConfiguration(strings.NewReader(""))
	var te *json.SyntaxError
	assert.ErrorAs(t, err, &te)
}
