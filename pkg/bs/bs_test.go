package bs_test

import (
	"fmt"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func generateTestRepo(t *testing.T, files map[string]string) (*git.Repository, plumbing.Hash, plumbing.Hash) {
	require.NotNil(t, files)
	fs := memfs.New()
	r, err := git.Init(memory.NewStorage(), fs)
	require.Nil(t, err)

	for fName, fContent := range files {
		f, err := fs.Create(fName)
		require.Nil(t, err)
		_, err = fmt.Fprint(f, fContent)
		require.Nil(t, err)
		require.Nil(t, f.Close())
	}

	w, err := r.Worktree()
	require.Nil(t, err)

	err = w.AddWithOptions(&git.AddOptions{All: true})
	require.Nil(t, err)

	hash, err := w.Commit("Initial commit", &git.CommitOptions{})
	require.Nil(t, err)

	commitObject, err := r.CommitObject(hash)
	require.Nil(t, err)
	treeHash := commitObject.TreeHash

	return r, hash, treeHash
}

func stdTestRepoContent() map[string]string {
	return map[string]string{
		"a/hello": "hello a",
		"b/hello": "hello b",
		"c/hello": "hello c",
	}
}

func generateStdTestRepo(t *testing.T) (*git.Repository, plumbing.Hash, plumbing.Hash) {
	return generateTestRepo(t, stdTestRepoContent())
}

func TestTestRepoGeneration(t *testing.T) {
	r, _, treeHash := generateStdTestRepo(t)
	expectedTreeHash := plumbing.NewHash("3891d8802de47445f72489934e010cbc18e395fd")
	require.NotNil(t, r)
	assert.Equal(t, expectedTreeHash, treeHash)
}
