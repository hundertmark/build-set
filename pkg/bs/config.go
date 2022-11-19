package bs

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
)

type BuildSet struct {
	Name       string   `yaml:"-"`
	HashOutput string   `yaml:"hash-output,omitempty"`
	Remote     string   `yaml:"remote,omitempty"`
	Include    []string `yaml:"include"`
	Exclude    []string `yaml:"exclude,omitempty"`
}

func (set *BuildSet) GetExcludeMatcher() gitignore.Matcher {
	var ps []gitignore.Pattern
	for _, s := range append(set.Exclude, set.HashOutput) {
		ps = append(ps, gitignore.ParsePattern(s, nil))
	}
	return gitignore.NewMatcher(ps)
}

type BuildSetConfig struct {
	BuildSets map[string]*BuildSet `yaml:"build-sets"`
}

func ReadBuildSetConfig(in io.Reader) (*BuildSetConfig, error) {
	var bsc BuildSetConfig
	bytes, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(bytes, &bsc)
	for name, set := range bsc.BuildSets {
		set.Name = name
	}
	return &bsc, err
}

func (bsc *BuildSetConfig) Write(w io.Writer) error {
	bscBytes, err := yaml.Marshal(bsc)
	if err != nil {
		return err
	}
	_, err = w.Write(bscBytes)
	return err
}

const BuildSetConfigFileName = "bsconfig.yml"

func ReadBuildSetConfigFromIndex(r *git.Repository) (*BuildSetConfig, error) {
	idx, err := r.Storer.Index()
	if err != nil {
		return nil, err
	}
	bsConfigEntry, err := idx.Entry(BuildSetConfigFileName)
	if err != nil {
		return nil, err
	}
	bsConfigObject, err := r.BlobObject(bsConfigEntry.Hash)
	if err != nil {
		return nil, err
	}
	bsConfigReader, err := bsConfigObject.Reader()
	if err != nil {
		return nil, err
	}
	return ReadBuildSetConfig(bsConfigReader)
}
