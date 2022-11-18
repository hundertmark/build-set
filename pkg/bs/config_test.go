package bs

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestReadBuildSetConfig(t *testing.T) {
	bsConfigYaml := `
a:
 hash-output: "a/.hash"
 include: [
  "a/*",
  "c/*"
 ]
 exclude: [
  ".gitignore"
 ]

b:
 hash-output: "b/.hash"
 include: [
  "b/*",
  "c/*"
 ]
 exclude: [
  ".gitignore"
 ]
`[1:]
	bsConfig, err := ReadBuildSetConfig(strings.NewReader(bsConfigYaml))
	assert.Nil(t, err)
	a := &BuildSet{
		Name:       "a",
		HashOutput: "a/.hash",
		Remote:     "",
		Include:    []string{"a/*", "c/*"},
		Exclude:    []string{".gitignore"},
	}
	b := &BuildSet{
		Name:       "b",
		HashOutput: "b/.hash",
		Remote:     "",
		Include:    []string{"b/*", "c/*"},
		Exclude:    []string{".gitignore"},
	}
	assert.Equal(t, a, bsConfig.BuildSets["a"])
	assert.Equal(t, b, bsConfig.BuildSets["b"])
}
