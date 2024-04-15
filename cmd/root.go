/*
Copyright © 2024 Frank Mayer
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tsukinoko-kun/drop/internal/git"
	"github.com/tsukinoko-kun/drop/internal/ui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "drop",
	Short: "drop is a replacement for the GNU `rm` command.",
	Long:  "drop is a replacement for the GNU `rm` command. It is a simple, fast, and safe way to delete files and directories.",
	RunE: func(_ *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("no files or directories specified")
		}
		for _, arg := range args {
			if files, err := filepath.Glob(arg); err == nil {
				if len(files) != 0 {
					for _, file := range files {
						if err := del(file); err != nil {
							return errors.Join(fmt.Errorf("failed to delete %q", file), err)
						}
					}
					continue
				}
			}
			if err := del(arg); err != nil {
				return errors.Join(fmt.Errorf("failed to delete %q", arg), err)
			}
		}
		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Example = "drop file1 file2 file3"
}

var ask bool = true

func del(file string) error {
	fmt.Println("del", file)
	if ask {
		gitWarn := make(chan string)
		go git.FindUncomittedGit(file, gitWarn)
	loop:
		for dir := range gitWarn {
			i, err := ui.Choose(dir, "Delete anyway", "Abort", "Delete all without asking")
			if err != nil {
				return err
			}
			switch i {
			case 1: // Abort
				return nil
			case 2: // Delete all without asking
				ask = false
				break loop
			}
		}
	}
	return os.RemoveAll(file)
}
