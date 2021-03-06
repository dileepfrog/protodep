package repository

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/stormcat24/protodep/dependency"
	"github.com/stormcat24/protodep/helper"
	"github.com/stormcat24/protodep/logger"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

type GitRepository interface {
	Open() (*OpenedRepository, error)
	ProtoRootDir() string
}

type GitHubRepository struct {
	protodepDir  string
	dep          dependency.ProtoDepDependency
	authProvider helper.AuthProvider
}

func NewGitRepository(protodepDir string, dep dependency.ProtoDepDependency, authProvider helper.AuthProvider) GitRepository {
	return &GitHubRepository{
		protodepDir:  protodepDir,
		dep:          dep,
		authProvider: authProvider,
	}
}

type OpenedRepository struct {
	Repository *git.Repository
	Dep        dependency.ProtoDepDependency
	Hash       string
}

func (r *GitHubRepository) Open() (*OpenedRepository, error) {

	tag := r.dep.Tag
	branch := ""

	if tag == "" {
		branch = "master"
		if r.dep.Branch != "" {
			branch = r.dep.Branch
		}
	}

	revision := r.dep.Revision

	reponame := r.dep.Repository()
	repopath := filepath.Join(r.protodepDir, reponame)

	var rep *git.Repository

	if stat, err := os.Stat(repopath); err == nil && stat.IsDir() {
		spinner := logger.InfoWithSpinner("Fetching %s", reponame)

		rep, err = git.PlainOpen(repopath)
		if err != nil {
			return nil, errors.Wrap(err, "open repository is failed")
		}

		cmd := exec.Command("git", "fetch")
		cmd.Dir = repopath

		result, err := cmd.CombinedOutput()
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("fetch repository is failed: %s", string(result)))
		}

		spinner.Finish()
	} else {
		spinner := logger.InfoWithSpinner("Cloning %s", reponame)

		url := r.authProvider.GetRepositoryURL(reponame)
		cmd := exec.Command("git", "clone", url, repopath)

		result, err := cmd.CombinedOutput()
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("clone repository is failed: %s", string(result)))
		}

		rep, err = git.PlainOpen(repopath)
		if err != nil {
			return nil, errors.Wrap(err, "open repository is failed")
		}

		spinner.Finish()
	}

	wt, err := rep.Worktree()
	if err != nil {
		return nil, errors.Wrap(err, "get worktree is failed")
	}

	if revision == "" {
		var target *plumbing.Reference

		if branch != "" {
			target, err = rep.Storer.Reference(plumbing.ReferenceName(fmt.Sprintf("refs/remotes/origin/%s", branch)))
			if err != nil {
				return nil, errors.Wrapf(err, "change branch to %s is failed", branch)
			}
		} else {
			target, err = rep.Storer.Reference(plumbing.ReferenceName(fmt.Sprintf("refs/tags/%s", tag)))
			if err != nil {
				return nil, errors.Wrapf(err, "change tag to %s is failed", tag)
			}
		}

		if err := wt.Checkout(&git.CheckoutOptions{Hash: target.Hash()}); err != nil {
			return nil, errors.Wrapf(err, "checkout to %s is failed", revision)
		}

		head := plumbing.NewHashReference(plumbing.HEAD, target.Hash())
		if err := rep.Storer.SetReference(head); err != nil {
			return nil, errors.Wrapf(err, "set head to %s is failed", branch)
		}
	} else {
		hash := plumbing.NewHash(revision)
		if err := wt.Checkout(&git.CheckoutOptions{Hash: hash}); err != nil {
			return nil, errors.Wrapf(err, "checkout to %s is failed", revision)
		}

		head := plumbing.NewHashReference(plumbing.HEAD, hash)
		if err := rep.Storer.SetReference(head); err != nil {
			return nil, errors.Wrapf(err, "set head to %s is failed", revision)
		}
	}

	commiter, err := rep.Log(&git.LogOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "get commit is failed")
	}

	current, err := commiter.Next()
	if err != nil {
		return nil, errors.Wrap(err, "get commit current is failed")
	}

	return &OpenedRepository{
		Repository: rep,
		Dep:        r.dep,
		Hash:       current.Hash.String(),
	}, nil
}

func (r *GitHubRepository) ProtoRootDir() string {
	return filepath.Join(r.protodepDir, r.dep.Target)
}
