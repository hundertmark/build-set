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

func main() {
	r, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{DetectDotGit: true})
	exitOnErr(err)
	buildSetConfig, err := bs.ReadBuildSetConfigFromIndex(r)
	exitOnErr(err)
	for _, set := range buildSetConfig.BuildSets {
		exitOnErr(bs.AddHashOutput(r, set))
	}
}
