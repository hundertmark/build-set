package bs

import (
	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
)

type BuildSet struct {
	Name       string   `yaml:"name"`
	HashOutput string   `yaml:"hash-output"`
	Remote     string   `yaml:"remote"`
	Include    []string `yaml:"include"`
	Exclude    []string `yaml:"exclude"`
}

func (set *BuildSet) GetExcludeMatcher() gitignore.Matcher {
	var ps []gitignore.Pattern
	for _, s := range append(set.Exclude, set.HashOutput) {
		ps = append(ps, gitignore.ParsePattern(s, nil))
	}
	return gitignore.NewMatcher(ps)
}

type BuildSetConfig struct {
	BuildSets map[string]*BuildSet
}

func ReadBuildSetConfig(in io.Reader) (*BuildSetConfig, error) {
	bsConfig := &BuildSetConfig{BuildSets: map[string]*BuildSet{}}
	bytes, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(bytes, &bsConfig.BuildSets)
	for name, set := range bsConfig.BuildSets {
		set.Name = name
	}
	return bsConfig, err
}
