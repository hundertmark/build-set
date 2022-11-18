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
	w, err := r.Worktree()
	exitOnErr(err)
	bsConfigFile, err := w.Filesystem.Open(BuildSetConfigFileName)
	exitOnErr(err)
	buildSetConfig, err := bs.ReadBuildSetConfig(bsConfigFile)
	exitOnErr(err)
	for _, set := range buildSetConfig.BuildSets {
		exitOnErr(bs.AddHashOutput(r, set))
	}
}
