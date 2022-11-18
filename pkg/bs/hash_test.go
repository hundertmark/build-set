package bs_test

import (
	"github.com/go-git/go-git/v5"
	"github.com/hundertmark/build-set/pkg/bs"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestBuildSetHashFromIndex_Clean(t *testing.T) {
	allSet := &bs.BuildSet{
		Include: []string{
			"*",
		},
	}
	aSet := &bs.BuildSet{
		Include: []string{
			"a/*",
			"c/*",
		},
	}
	bSet := &bs.BuildSet{
		Include: []string{
			"b/*",
			"c/*",
		},
	}
	r, _, expectedAllHash := generateStdTestRepo(t)
	_, _, expectedAHash := generateTestRepo(t, map[string]string{"a/hello": "hello a", "c/hello": "hello c"})
	_, _, expectedBHash := generateTestRepo(t, map[string]string{"b/hello": "hello b", "c/hello": "hello c"})

	allHash, err := bs.BuildSetHashFromIndex(r, allSet)
	assert.Nil(t, err)
	aHash, err := bs.BuildSetHashFromIndex(r, aSet)
	assert.Nil(t, err)
	bHash, err := bs.BuildSetHashFromIndex(r, bSet)
	assert.Nil(t, err)
	assert.Equal(t, expectedAllHash, allHash)
	assert.Equal(t, expectedAHash, aHash)
	assert.Equal(t, expectedBHash, bHash)
}

func getStatus(t *testing.T, r *git.Repository) string {
	w, err := r.Worktree()
	assert.Nil(t, err)
	s, err := w.Status()
	assert.Nil(t, err)
	return s.String()
}

func readFile(t *testing.T, r *git.Repository, filename string) string {
	w, err := r.Worktree()
	assert.Nil(t, err)
	f, err := w.Filesystem.Open(filename)
	assert.Nil(t, err)
	fbytes, err := ioutil.ReadAll(f)
	assert.Nil(t, err)
	return string(fbytes[:])
}

func TestAddHashOutput(t *testing.T) {
	aSet := &bs.BuildSet{
		Name:       "a",
		HashOutput: "a/.bs_hash",
		Include: []string{
			"a/*",
			"c/*",
		},
	}
	r, _, _ := generateStdTestRepo(t)
	_, _, expectedAHash := generateTestRepo(t, map[string]string{"a/hello": "hello a", "c/hello": "hello c"})
	err := bs.AddHashOutput(r, aSet)
	assert.Nil(t, err)
	status := getStatus(t, r)
	assert.Equal(t, "A  a/.bs_hash\n", status)
	assert.Equal(t, expectedAHash.String(), readFile(t, r, "a/.bs_hash"))
}
