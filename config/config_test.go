// Copyright (c) Jeevanandam M (https://github.com/jeevatkm)
// Source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package config

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"aahframe.work/essentials"
	"github.com/stretchr/testify/assert"
)

func TestKeysAndKeysByPath(t *testing.T) {
	cfg := initFile(t, join(testdataBaseDir(), "test.cfg"))
	keys := cfg.Keys()
	assert.True(t, ess.IsSliceContainsString(keys, "prod"))
	assert.True(t, ess.IsSliceContainsString(keys, "dev"))
	assert.True(t, ess.IsSliceContainsString(keys, "subsection"))

	keys = cfg.KeysByPath("prod")
	assert.True(t, ess.IsSliceContainsString(keys, "float32"))
	assert.True(t, ess.IsSliceContainsString(keys, "string"))
	assert.True(t, ess.IsSliceContainsString(keys, "truevalue"))

	// Key not exists
	keys = cfg.KeysByPath("prod.not.found")
	assert.True(t, len(keys) == 0)

	// key is not a section
	keys = cfg.KeysByPath("prod.truevalue")
	assert.True(t, len(keys) == 0)
}

func TestGetSubConfig(t *testing.T) {
	cfg := initFile(t, join(testdataBaseDir(), "test.cfg"))

	prod, found1 := cfg.GetSubConfig("prod")
	assert.NotNil(t, prod)
	assert.True(t, found1)

	setProfileForTest(t, cfg, "prod")
	subsection, found2 := cfg.GetSubConfig("subsection")
	assert.NotNil(t, subsection)
	assert.True(t, found2)

	value, f := subsection.Float32("sub_float")
	assert.True(t, f)
	assert.Equal(t, float32(100.5), value)

	assert.NotNil(t, NewEmpty())
	assert.True(t, len(cfg.ToJSON()) > 0)

	cfg.ClearProfile()

	// Key not exists
	_, prod2NotFound := cfg.GetSubConfig("prod.not.found")
	assert.False(t, prod2NotFound)

	// key is not a section
	_, keyNotASection := cfg.GetSubConfig("prod.truevalue")
	assert.False(t, keyNotASection)
}

func TestIsExists(t *testing.T) {
	cfg := initFile(t, join(testdataBaseDir(), "test.cfg"))
	found := cfg.IsExists("prod.string")
	assert.True(t, found)

	found = cfg.IsExists("prod.not.found")
	assert.False(t, found)

	setProfileForTest(t, cfg, "prod")
	found = cfg.IsExists("string")
	assert.True(t, found)
}

func TestStringValues(t *testing.T) {
	cfg := initFile(t, join(testdataBaseDir(), "test.cfg"))

	v1, _ := cfg.String("string")
	assert.Equal(t, "a string", v1)
	assert.Equal(t, cfg.StringDefault("string_not_exists", "nice 1"), "nice 1")
	assert.Equal(t, cfg.StringDefault("string", "nice 1"), "a string")

	setProfileForTest(t, cfg, "dev")
	dv1, _ := cfg.String("string")
	assert.Equal(t, "a string inside dev", dv1)
	assert.Equal(t, cfg.StringDefault("string_not_exists", "nice 2"), "nice 2")

	setProfileForTest(t, cfg, "prod")
	pv1, _ := cfg.String("string")
	assert.Equal(t, "a string inside prod", pv1)
	assert.Equal(t, cfg.StringDefault("string_not_exists", "nice 3"), "nice 3")

	cfg.SetString("string", "a string is inside prod")
	pv2, _ := cfg.String("string")
	assert.Equal(t, "a string is inside prod", pv2)
}

func TestIntValues(t *testing.T) {
	bytes, _ := ioutil.ReadFile(join(testdataBaseDir(), "test.cfg"))
	cfg := initString(t, string(bytes))

	v1, _ := cfg.Int("int")
	assert.Equal(t, 32, v1)

	v2, _ := cfg.Int64("int64")
	assert.Equal(t, int64(1), v2)

	v3, _ := cfg.Int64("int64not")
	assert.Equal(t, int64(0), v3)
	assert.Equal(t, cfg.IntDefault("int_not_exists", 99), 99)
	assert.Equal(t, cfg.IntDefault("int", 99), 32)

	cfg.SetInt("int", 30)
	v4, _ := cfg.Int("int")
	assert.Equal(t, 30, v4)

	setProfileForTest(t, cfg, "dev")
	dv1, _ := cfg.Int("int")
	assert.Equal(t, 500, dv1)

	dv2, _ := cfg.Int64("int64")
	assert.Equal(t, int64(2), dv2)
	assert.Equal(t, cfg.IntDefault("int_not_exists", 199), 199)

	setProfileForTest(t, cfg, "prod")
	pv1, _ := cfg.Int("int")
	assert.Equal(t, 1000, pv1)

	pv2, _ := cfg.Int64("int64")
	assert.Equal(t, int64(3), pv2)
	assert.Equal(t, cfg.IntDefault("int_not_exists", 299), 299)
}

func TestFloatValues(t *testing.T) {
	cfg := initFile(t, join(testdataBaseDir(), "test.cfg"))

	v1, _ := cfg.Float32("float32")
	assert.Equal(t, float32(32.2), v1)

	v2, _ := cfg.Float64("float64")
	assert.Equal(t, float64(1.1), v2)

	v3, _ := cfg.Float64("subsection.sub_float")
	assert.Equal(t, float64(10.5), v3)

	v4, _ := cfg.Float64("float64not")
	assert.Equal(t, float64(0.0), v4)
	assert.Equal(t, cfg.Float32Default("float_not_exists", float32(99.99)), float32(99.99))
	assert.Equal(t, cfg.Float32Default("float32", float32(99.99)), float32(32.2))

	cfg.SetFloat32("float32", float32(38.2))
	v5, _ := cfg.Float32("float32")
	assert.Equal(t, float32(38.2), v5)

	setProfileForTest(t, cfg, "dev")
	dv1, _ := cfg.Float32("float32")
	assert.Equal(t, float32(62.2), dv1)

	dv2, _ := cfg.Float64("float64")
	assert.Equal(t, float64(2.1), dv2)

	dv3, _ := cfg.Float64("subsection.sub_float")
	assert.Equal(t, float64(50.5), dv3)
	assert.Equal(t, cfg.Float32Default("float_not_exists", float32(199.99)), float32(199.99))

	setProfileForTest(t, cfg, "prod")
	pv1, _ := cfg.Float32("float32")
	assert.Equal(t, float32(122.2), pv1)

	pv2, _ := cfg.Float64("float64")
	assert.Equal(t, float64(3.1), pv2)

	pv3, _ := cfg.Float64("subsection.sub_float")
	assert.Equal(t, float64(100.5), pv3)
	assert.Equal(t, cfg.Float32Default("float_not_exists", float32(299.99)), float32(299.99))
}

func TestBoolValues(t *testing.T) {
	bytes, _ := ioutil.ReadFile(join(testdataBaseDir(), "test.cfg"))
	cfg := initString(t, string(bytes))

	v1, _ := cfg.Bool("truevalue")
	assert.True(t, v1)

	v2, _ := cfg.Bool("falsevalue")
	assert.False(t, v2)
	assert.Equal(t, cfg.BoolDefault("bool_not_exists", true), true)
	assert.Equal(t, cfg.BoolDefault("falsevalue", true), false)

	cfg.SetBool("truevalue", false)
	v3, _ := cfg.Bool("truevalue")
	assert.False(t, v3)
	cfg.SetBool("truevalue", true)

	setProfileForTest(t, cfg, "dev")
	assert.Equal(t, cfg.BoolDefault("truevalue", true), true)    // keys is found by fallback
	assert.Equal(t, cfg.BoolDefault("falsevalue", false), false) // keys is found by fallback

	setProfileForTest(t, cfg, "prod")
	pv1, _ := cfg.Bool("truevalue")
	assert.True(t, pv1)

	pv2, _ := cfg.Bool("falsevalue")
	assert.False(t, pv2)
	assert.Equal(t, cfg.BoolDefault("bool_not_exists", true), true)
}

func TestStringList(t *testing.T) {
	cfg, _ := ParseString(`
	# build config read from here during a command 'aah package' and 'aah run'
build {
  # Valid exclude patterns refer: https://golang.org/pkg/path/filepath/#Match
  excludes = ["*_test.go", ".*", "*.bak", "*.tmp", "vendor"]

	keys = [
		"X3pGTSOuJeEVw989IJ/cEtXUEmy52zs1TZQrU06KUKg=",
		"MHJYVThihUrJcxW6wcqyOISTXIsInsdj3xK8QrZbHec=",
		"GGekerhihUrJcxW6wcqyOISTXIsInsdj3xK8QrZbHec=",
	]
}`)

	lst1, found1 := cfg.StringList("build.excludes")
	assert.True(t, found1)
	assert.True(t, len(lst1) > 0)
	assert.Equal(t, "*.bak", lst1[2])

	lst2, found2 := cfg.StringList("name")
	assert.False(t, found2)
	assert.True(t, len(lst2) == 0)

	v, f := cfg.StringList("build.keys")
	assert.True(t, f)
	assert.True(t, len(v) == 3)
}

func TestIntAndInt64List(t *testing.T) {
	cfg, _ := ParseString(`
		int_list = [10, 20, 30, 40, 50]
		int64_list = [100000001, 100000002, 100000003, 100000004, 100000005]
	`)

	lst1, found1 := cfg.IntList("int_list")
	assert.True(t, found1)
	assert.True(t, len(lst1) > 0)
	assert.Equal(t, int(20), lst1[1])

	lst2, found2 := cfg.Int64List("int64_list")
	assert.True(t, found2)
	assert.True(t, len(lst2) > 0)
	assert.Equal(t, int64(100000005), lst2[4])

	lst3, found3 := cfg.IntList("name_not_found")
	assert.False(t, found3)
	assert.True(t, len(lst3) == 0)

	lst4, found4 := cfg.Int64List("int64_list")
	assert.True(t, found4)
	assert.True(t, len(lst4) > 0)
}

func TestProfile(t *testing.T) {
	cfg := initFile(t, join(testdataBaseDir(), "test.cfg"))

	t.Log(cfg.cfg.Keys())

	assert.Equal(t, "", cfg.Profile())
	assert.Equal(t, "profile doesn't exists: not_exists_profile",
		cfg.SetProfile("not_exists_profile").Error())

	cfg.ClearProfile()
	assert.Equal(t, true, cfg.profile == "")
}

func TestConfigLoadNotExists(t *testing.T) {
	_, err := LoadFile(join(testdataBaseDir(), "not_exists.cfg"))
	assert.True(t, strings.Contains(err.Error(), "does not exists:"))

	_, err = ParseString(`
  # Error configuration
  string = "a string"
  int = 32 # adding comment without semicolon will lead to error
  float32 = 32.2
  int64 = 1
  float64 = 1.1
	`)
	assert.Equal(t, true,
		strings.Contains(err.Error(), "adding comment without semicolon will lead to error"))
}

func TestMergeConfig(t *testing.T) {
	cfg1, _ := ParseString(`
global = "global value";

prod {
	value = "string value";
	integer = 500
	float = 80.80
	boolean = true
	negative = FALSE
	nothing = NULL
}
	`)

	cfg2, _ := ParseString(`
global = "global value";

newvalue = "I'm new value"
prod {
  value = "I'm prod value"
	nothing = 200
}
`)

	err := cfg1.Merge(cfg2)
	assert.NoErrorf(t, err, "merge failed")

	t.Log("Merge2Section test")
	if tocfg, found := cfg1.GetSubConfig("prod"); found {
		err = cfg2.Merge2Section("new.prod", tocfg)
		assert.Nil(t, err)
		assert.Equal(t, 500, cfg2.IntDefault("new.prod.integer", 0))
	}
	err = cfg1.Merge2Section("", nil)
	assert.Equal(t, errors.New("source is nil"), err)
	err = cfg1.Merge2Section("", cfg2)
	assert.Equal(t, errors.New("key is empty"), err)

	_ = cfg1.SetProfile("prod")

	v1 := cfg1.IntDefault("nothing", 0)
	assert.Equal(t, 200, v1)

	v2 := cfg1.StringDefault("value", "")
	assert.Equal(t, "I'm prod value", v2)

	v3 := cfg1.StringDefault("newvalue", "")
	assert.Equal(t, "I'm new value", v3)

	err = cfg1.Merge(nil)
	assert.Equal(t, "source is nil", err.Error())
}

func TestLoadFiles(t *testing.T) {
	testdataPath := testdataBaseDir()

	cfg, err := LoadFiles(
		join(testdataPath, "test-1.cfg"),
		join(testdataPath, "test-2.cfg"),
		join(testdataPath, "test-3.cfg"),
	)
	assert.NoErrorf(t, err, "loading failed")

	assert.Equal(t, float32(10.5), cfg.Float32Default("subsection.sub_float", 0.0))
	assert.Equal(t, float32(32.4), cfg.Float32Default("float32", 0.0))

	_ = cfg.SetProfile("dev")
	assert.Equal(t, "a string inside dev from test-2", cfg.StringDefault("string", ""))
	assert.Equal(t, float32(500.5), cfg.Float32Default("subsection.sub_float", 0.0))
	assert.Equal(t, float32(62.2), cfg.Float32Default("float32", 0.0))

	_ = cfg.SetProfile("prod")
	assert.Equal(t, "a string inside prod from test-3", cfg.StringDefault("string", ""))
	assert.Equal(t, float32(1000.5), cfg.Float32Default("subsection.sub_float", 0.0))
	assert.Equal(t, float32(222.2), cfg.Float32Default("float32", 0.0))
	assert.Equal(t, true, cfg.BoolDefault("falsevalue", false))

	// fail cases
	_, err = LoadFiles(join(testdataPath, "not_exists.cfg"))
	assert.Equal(t, true, strings.Contains(err.Error(), "does not exists:"))

	_, err = LoadFiles(
		join(testdataPath, "test-1.cfg"),
		join(testdataPath, "test-error.cfg"),
	)
	assert.Equal(t, true, strings.Contains(err.Error(), "source (STRING) and target (SECTION)"))
}

func TestNestedKeyWithProfile(t *testing.T) {
	testdataPath := testdataBaseDir()

	cfg, err := LoadFiles(join(testdataPath, "test-4.cfg"))
	assert.NoErrorf(t, err, "loading failed")

	releases, found := cfg.StringList("docs.releases")
	assert.True(t, found)
	assert.True(t, len(releases) > 0)

	err = cfg.SetProfile("env.dev")
	assert.Nil(t, err)

	assert.Equal(t, "Not in env section", cfg.StringDefault("not_in_env", "defaultValue"))

	releases, found = cfg.StringList("docs.releases")
	assert.True(t, found)
	assert.True(t, len(releases) > 0)

	// not a list type
	devStr, found := cfg.StringList("env.dev.string")
	assert.False(t, found)
	assert.True(t, len(devStr) == 0)

	value := cfg.StringDefault("env.path", "")
	assert.NotNil(t, value)
}

func TestConfigSetValues(t *testing.T) {
	testdataPath := testdataBaseDir()
	cfg, err := LoadFiles(join(testdataPath, "test-4.cfg"))
	assert.NoErrorf(t, err, "loading failed")

	// env.active
	assert.False(t, cfg.IsExists("env.active"))
	cfg.SetString("env.active", "dev")
	assert.True(t, cfg.IsExists("env.active"))
	assert.Equal(t, "dev", cfg.StringDefault("env.active", ""))

	cfg.SetString("env.active", "prod")
	assert.Equal(t, "prod", cfg.StringDefault("env.active", ""))

	// request.id.header
	assert.False(t, cfg.IsExists("request.id.header"))
	cfg.SetString("request.id.header", "My-Request-Hdr")
	assert.True(t, cfg.IsExists("request.id.header"))
	assert.Equal(t, "My-Request-Hdr", cfg.StringDefault("request.id.header", ""))
}

func initString(t *testing.T, configStr string) *Config {
	cfg, err := ParseString(configStr)
	if !assert.NoErrorf(t, err, "loading failed") {
		assert.FailNow(t, "parse error")
	}

	return cfg
}

func initFile(t *testing.T, file string) *Config {
	cfg, err := LoadFile(file)
	if !assert.NoErrorf(t, err, "loading failed") {
		assert.FailNow(t, "parse error")
	}
	return cfg
}

func setProfileForTest(t *testing.T, cfg *Config, profile string) {
	err := cfg.SetProfile(profile)
	if !assert.NoErrorf(t, err, "loading failed") {
		assert.FailNow(t, "parse error")
	}
}

func testdataBaseDir() string {
	wd, _ := os.Getwd()
	return filepath.Join(wd, "testdata")
}

func join(elem ...string) string {
	return filepath.Join(elem...)
}
