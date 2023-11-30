package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

var RefPrefix = "VPA-PR-SUBMITTER"

func main() {
	now := time.Now()
	dir, _ := os.Getwd()
	fmt.Printf("dir: %s\n", dir)
	repo, err := git.PlainOpen("/Users/i2tsuki/Repo/i2tsuki/k8s-expt")
	if err != nil {
		log.Printf("Err: %v, %s", err, fmt.Sprintf("cannot open the repository: %s", dir))
		return
	}

	// get the reference name as `defaultRefName` of the default branch.
	defaultRefName := plumbing.NewBranchReferenceName("main")
	refIter, err := repo.References()
	if err != nil {
		log.Printf("err: %w", err)
	}
	refIter.ForEach(func(r *plumbing.Reference) error {
		if r.Name().IsBranch() && !strings.Contains(r.Name().String(), RefPrefix) {
			defaultRefName = r.Name()
			return nil
		}
		return nil
	})

	wt, _ := repo.Worktree()
	err = wt.Pull(&git.PullOptions{RemoteName: "origin"})
	if err != nil {
		log.Printf("err: failed to pull: %w", err)
		return
	}

	// create the branch based on HEAD in the current branch.
	refname := plumbing.NewBranchReferenceName(fmt.Sprintf("%s/%s", RefPrefix, now.Format("20060102")))
	headRef, err := repo.Head()
	if err != nil {
		log.Printf("Err: %v, %s", fmt.Sprintf("failed to get hash of the head: %s", dir))
		return
	}
	ref := plumbing.NewHashReference(refname, headRef.Hash())
	err = repo.Storer.SetReference(ref)
	if err != nil {
		log.Printf("Err: %v, %s", err, fmt.Sprintf("failed to set the reference: %s", ref))
		return
	}

	// commit the change to the repository.
	hash, err := wt.Commit("Test Commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "test",
			Email: "test@example.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		log.Printf("Err: %v, %s", err, "failed to commit the change")
		return
	}
	ref = plumbing.NewHashReference(refname, hash)
	err = repo.Storer.SetReference(ref)
	if err != nil {
		log.Printf("Err: %v, %s", err, fmt.Sprintf("failed to set the reference: %s", ref))
		return
	}
	log.Printf("defaultRefName: %s\n", defaultRefName)

	if err := repo.Push(&git.PushOptions{
		RemoteName: "origin",
		RefSpecs: []config.RefSpec{
			config.RefSpec(refname + ":" + refname),
		},
		Progress: os.Stdout,
	}); err != git.NoErrAlreadyUpToDate && err != nil {
		log.Printf("Err: failed to push to the remote: %v", err)
		return
	}

	// reset the change.
	ref = plumbing.NewHashReference(defaultRefName, headRef.Hash())
	_ = repo.Storer.SetReference(ref)

}
