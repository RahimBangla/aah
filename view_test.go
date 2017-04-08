// Copyright (c) Jeevanandam M. (https://github.com/jeevatkm)
// go-aah/view source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package atemplate

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"aahframework.org/config.v0"
	"aahframework.org/test.v0/assert"
)

func TestViewAddTemplateFunc(t *testing.T) {
	AddTemplateFunc(template.FuncMap{
		"join": strings.Join,
	})

	_, found := TemplateFuncMap["join"]
	assert.True(t, found)
}

func TestViewStore(t *testing.T) {
	err := AddEngine("go", &GoViewEngine{})
	assert.NotNil(t, err)
	assert.Equal(t, "view: engine name 'go' is already added, skip it", err.Error())

	err = AddEngine("custom", nil)
	assert.NotNil(t, err)
	assert.Equal(t, "view: engine value is nil", err.Error())

	engine, found := GetEngine("go")
	assert.NotNil(t, engine)
	assert.True(t, found)

	engine, found = GetEngine("myengine")
	assert.Nil(t, engine)
	assert.False(t, found)
}

func TestViewCommonTemplateInit(t *testing.T) {
	c := &CommonTemplate{}
	cfg, _ := config.ParseString(`view { }`)

	err := c.Init(cfg, filepath.Join(getTestdataPath(), "common-not-exists"))
	assert.True(t, strings.HasPrefix(err.Error(), "commontemplate: base dir is not exists"))
}

func getTestdataPath() string {
	wd, _ := os.Getwd()
	return filepath.Join(wd, "testdata")
}
