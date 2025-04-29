//go:build !solution

package main

import (
	"os"

	"gitlab.com/slon/shad-go/gitfame/internal/analyse"
	"gitlab.com/slon/shad-go/gitfame/internal/formatter"
	"gitlab.com/slon/shad-go/gitfame/internal/git"
	"gitlab.com/slon/shad-go/gitfame/internal/options"
)

func main() {
	config := options.ParseOptions(os.Args[1:])

	gitProcessor := git.NewGit(config)
	analyzer := analyse.New()

	results := analyzer.AnalyzeGitFiles(gitProcessor)
	results.Sort(config.OrderBy)

	formatter.New(config.OutFormat, os.Stdout).Print(results)
}
