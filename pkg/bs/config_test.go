package bs_test

import (
	"bytes"
	"github.com/go-git/go-git/v5"
	"github.com/hundertmark/build-set/pkg/bs"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func stdTestBuildSetConfigYaml() string {
	return `
build-sets:
    a:
        hash-output: a/.hash
        include:
          - a/*
          - c/*
        exclude:
          - .gitignore
    b:
        hash-output: b/.hash
        include:
          - b/*
          - c/*
        exclude:
          - .gitignore
`[1:]
}

func stdBuildSetConfig() *bs.BuildSetConfig {
	return &bs.BuildSetConfig{BuildSets: map[string]*bs.BuildSet{
		"a": {
			Name:       "a",
			HashOutput: "a/.hash",
			Remote:     "",
			Include:    []string{"a/*", "c/*"},
			Exclude:    []string{".gitignore"},
		},
		"b": {
			Name:       "b",
			HashOutput: "b/.hash",
			Remote:     "",
			Include:    []string{"b/*", "c/*"},
			Exclude:    []string{".gitignore"},
		},
	}}
}

func TestReadBuildSetConfig(t *testing.T) {
	bsConfigFromString, err := bs.ReadBuildSetConfig(strings.NewReader(stdTestBuildSetConfigYaml()))
	assert.Nil(t, err)
	expectedBuildSetConfig := stdBuildSetConfig()
	assert.Equal(t, expectedBuildSetConfig, bsConfigFromString)
}

func TestBuildSetConfig_Write(t *testing.T) {
	var b bytes.Buffer
	bsConfig := stdBuildSetConfig()
	assert.Nil(t, bsConfig.Write(&b))
	assert.Equal(t, stdTestBuildSetConfigYaml(), b.String())
}

func TestReadBuildSetConfigFromIndex(t *testing.T) {
	r, _, _ := generateStdTestRepo(t)
	bsConfigFromStdRepo, err := bs.ReadBuildSetConfigFromIndex(r)
	assert.Nil(t, err)
	assert.Equal(t, stdBuildSetConfig(), bsConfigFromStdRepo)

	// We modify and save the build set config
	modifiedBuildSetConfig := stdBuildSetConfig()
	modifiedBuildSetConfig.BuildSets["a"].Include = []string{"a/*"}
	w, err := r.Worktree()
	assert.Nil(t, err)
	bsConfigFile, err := w.Filesystem.Create(bs.BuildSetConfigFileName)
	assert.Nil(t, err)
	assert.Nil(t, modifiedBuildSetConfig.Write(bsConfigFile))
	assert.Nil(t, bsConfigFile.Close())

	// Since the build set config is read from the index, we should still get the old config
	bsConfigFromModifiedRepo, err := bs.ReadBuildSetConfigFromIndex(r)
	assert.Nil(t, err)
	assert.Equal(t, stdBuildSetConfig(), bsConfigFromModifiedRepo)

	// After staging the modified build set config we should get the modified version from the repo
	assert.Nil(t, w.AddWithOptions(&git.AddOptions{Path: bs.BuildSetConfigFileName}))
	bsConfigFromStagedRepo, err := bs.ReadBuildSetConfigFromIndex(r)
	assert.Nil(t, err)
	assert.Equal(t, modifiedBuildSetConfig, bsConfigFromStagedRepo)
	assert.NotEqual(t, stdBuildSetConfig(), bsConfigFromStagedRepo)
}
