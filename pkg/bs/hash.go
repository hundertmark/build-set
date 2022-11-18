package bs

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/filemode"
	"github.com/go-git/go-git/v5/plumbing/format/index"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage"
	"path"
	"sort"
	"strings"
)

func AddHashOutput(r *git.Repository, set *BuildSet) error {
	hash, err := BuildSetHashFromIndex(r, set)
	if err != nil {
		return err
	}
	w, err := r.Worktree()
	if err != nil {
		return err
	}
	hashOutput, err := w.Filesystem.Create(set.HashOutput)
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(hashOutput, hash.String())
	if err != nil {
		return err
	}
	return w.AddWithOptions(&git.AddOptions{Path: set.HashOutput})
}

func BuildSetHashFromIndex(r *git.Repository, set *BuildSet) (plumbing.Hash, error) {
	s := r.Storer

	h := createBuildTreeHelper(s)
	idx, err := s.Index()
	if err != nil {
		return plumbing.ZeroHash, err
	}
	m := set.GetExcludeMatcher()
	for _, pattern := range set.Include {
		entries, err := idx.Glob(pattern)
		if err != nil {
			return plumbing.ZeroHash, err
		}
		for _, e := range entries {
			if !m.Match(strings.Split(e.Name, "/"), false) {
				if err := h.CommitIndexEntry(e); err != nil {
					return plumbing.ZeroHash, err
				}
			}
		}
	}
	return h.GetFingerprint()
}

// Based on github.com/go-git/worktree_commit.go:135-238
type buildTreeHelper struct {
	s storage.Storer

	trees   map[string]*object.Tree
	entries map[string]*object.TreeEntry
}

const rootNode = ""

func createBuildTreeHelper(s storage.Storer) *buildTreeHelper {
	return &buildTreeHelper{
		s:       s,
		trees:   map[string]*object.Tree{rootNode: {}},
		entries: map[string]*object.TreeEntry{},
	}
}

func (h *buildTreeHelper) GetFingerprint() (plumbing.Hash, error) {
	return h.copyTreeToStorageRecursive(rootNode, h.trees[rootNode], true)
}

func (h *buildTreeHelper) CommitIndexEntry(e *index.Entry) error {
	parts := strings.Split(e.Name, "/")

	var fullpath string
	for _, part := range parts {
		parent := fullpath
		fullpath = path.Join(fullpath, part)

		h.doBuildTree(e, parent, fullpath)
	}

	return nil
}

func (h *buildTreeHelper) doBuildTree(e *index.Entry, parent, fullpath string) {
	if _, ok := h.trees[fullpath]; ok {
		return
	}

	if _, ok := h.entries[fullpath]; ok {
		return
	}

	te := object.TreeEntry{Name: path.Base(fullpath)}

	if fullpath == e.Name {
		te.Mode = e.Mode
		te.Hash = e.Hash
	} else {
		te.Mode = filemode.Dir
		h.trees[fullpath] = &object.Tree{}
	}

	h.trees[parent].Entries = append(h.trees[parent].Entries, te)
}

type sortableEntries []object.TreeEntry

func (sortableEntries) sortName(te object.TreeEntry) string {
	if te.Mode == filemode.Dir {
		return te.Name + "/"
	}
	return te.Name
}
func (se sortableEntries) Len() int               { return len(se) }
func (se sortableEntries) Less(i int, j int) bool { return se.sortName(se[i]) < se.sortName(se[j]) }
func (se sortableEntries) Swap(i int, j int)      { se[i], se[j] = se[j], se[i] }

func (h *buildTreeHelper) copyTreeToStorageRecursive(parent string, t *object.Tree, onlyCalculateHash bool) (plumbing.Hash, error) {
	sort.Sort(sortableEntries(t.Entries))
	for i, e := range t.Entries {
		if e.Mode != filemode.Dir && !e.Hash.IsZero() {
			continue
		}

		ePath := path.Join(parent, e.Name)

		var err error
		e.Hash, err = h.copyTreeToStorageRecursive(ePath, h.trees[ePath], onlyCalculateHash)
		if err != nil {
			return plumbing.ZeroHash, err
		}

		t.Entries[i] = e
	}

	o := h.s.NewEncodedObject()
	if err := t.Encode(o); err != nil {
		return plumbing.ZeroHash, err
	}

	hash := o.Hash()
	if onlyCalculateHash {
		return hash, nil
	}
	if h.s.HasEncodedObject(hash) == nil {
		return hash, nil
	}
	return h.s.SetEncodedObject(o)
}
