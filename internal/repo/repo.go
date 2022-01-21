package repo

import (
	"io/ioutil"

	"github.com/cockroachdb/errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

const SourceRepo = `https://github.com/toptal/gitignore.git`

func New(path, repoSource string) (*git.Repository, error) {
	repo, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:      repoSource,
		Progress: ioutil.Discard,
	})

	switch err {
	case nil:
		return repo, nil
	case git.ErrRepositoryAlreadyExists:
	default:
		return nil, errors.Wrap(err, "repo: failed to clone")
	}

	repo, err = git.PlainOpen(path)

	return repo, errors.Wrap(err, "repo: failed open repo")
}

type repoer interface {
	Worktree() (*git.Worktree, error)
	Head() (*plumbing.Reference, error)
}

func Checkout(r repoer, commitHash string) (string, error) {
	wt, err := r.Worktree()
	if err != nil {
		return "", errors.Wrap(err, "repo: getting worktree")
	}

	opts := git.CheckoutOptions{Force: true}

	if commitHash != "" {
		opts.Hash = plumbing.NewHash(commitHash)
	}

	if err := wt.Checkout(&opts); err != nil {
		return "", errors.Wrap(err, "repo: checkout")
	}

	ref, err := r.Head()
	if err != nil {
		return "", errors.Wrap(err, "repo: get head")
	}

	return ref.Hash().String(), nil
}
