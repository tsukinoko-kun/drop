package git

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
)

func FindGit(path string, gitDirs chan<- string) error {
	dirStack := []string{path}
	defer close(gitDirs)

stack:
	for len(dirStack) != 0 {
		dir := dirStack[0]
		dirStack = dirStack[1:]

		content, err := os.ReadDir(dir)
		if err != nil {
			return err
		}
		subfolders := make([]string, 0)
		for _, entry := range content {
			if !entry.IsDir() {
				continue
			}
			if entry.Name() == ".git" {
				gitDirs <- dir
				continue stack
			}
			p, err := filepath.Abs(filepath.Join(dir, entry.Name()))
			if err != nil {
				return err
			}
			subfolders = append(subfolders, p)
		}
		dirStack = append(subfolders, dirStack...)
	}

	return nil
}

func FindUncomittedGit(path string, uncomittedGitDirs chan<- string) {
	gitDirs := make(chan string)
	defer close(uncomittedGitDirs)
	go func() {
		err := FindGit(path, gitDirs)
		if err != nil {
			fmt.Println(err)
		}
	}()
	for dir := range gitDirs {
		repo, err := git.PlainOpen(dir)
		if err != nil {
			continue
		}
		// has remote?
		r, err := repo.Remotes()
		if err != nil {
			continue
		}
		if len(r) == 0 {
			uncomittedGitDirs <- fmt.Sprintf("Repository %s has no remote", dir)
			continue
		}
		// is clean?
		w, err := repo.Worktree()
		if err != nil {
			continue
		}
		s, err := w.Status()
		if err != nil {
			continue
		}
		if !s.IsClean() {
			uncomittedGitDirs <- fmt.Sprintf("Repository %s has uncomitted changes", dir)
			continue
		}
	}
}
