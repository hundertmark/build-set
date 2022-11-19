package main

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/hundertmark/build-set/pkg/bs"
	"os"
)

func exitOnErr(err error) {
	if err != nil {
		fmt.Printf("bs caused an error: %v\n", err)
		os.Exit(1)
	}
}

const BuildSetConfigFileName = "bsconfig.yml"

func main() {
	r, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{DetectDotGit: true})
	exitOnErr(err)
	idx, err := r.Storer.Index()
	exitOnErr(err)
	bsConfigEntry, err := idx.Entry(BuildSetConfigFileName)
	exitOnErr(err)
	bsConfigObject, err := r.BlobObject(bsConfigEntry.Hash)
	exitOnErr(err)
	bsConfigReader, err := bsConfigObject.Reader()
	exitOnErr(err)
	buildSetConfig, err := bs.ReadBuildSetConfig(bsConfigReader)
	exitOnErr(err)
	for _, set := range buildSetConfig.BuildSets {
		exitOnErr(bs.AddHashOutput(r, set))
	}
}
