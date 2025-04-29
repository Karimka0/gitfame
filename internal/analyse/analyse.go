package analyse

import (
	"strconv"
	"strings"

	"gitlab.com/slon/shad-go/gitfame/internal/author"
	"gitlab.com/slon/shad-go/gitfame/internal/git"
)

type Analysis struct {
	authors      map[string]*author.Author
	commitCounts map[string]int
	fileCounts   map[string]int
	authorHashes map[string]map[string]struct{}
}

func New() *Analysis {
	return &Analysis{
		authors:      make(map[string]*author.Author),
		commitCounts: make(map[string]int),
		fileCounts:   make(map[string]int),
		authorHashes: make(map[string]map[string]struct{}),
	}
}

func (a *Analysis) AnalyzeGitFiles(g *git.Git) author.AuthorSlice {
	for _, file := range author.GetFiles(g) {
		a.processFile(g, file)
	}
	a.normalize()
	return a.compileResults()
}

func (a *Analysis) processFile(g *git.Git, filename string) {
	authorCommits, lineCounts := blameAnalysis(g, filename)

	for hash, count := range lineCounts {
		a.commitCounts[hash] += count
	}

	for name, hashes := range authorCommits {
		a.fileCounts[name]++
		if a.authorHashes[name] == nil {
			a.authorHashes[name] = make(map[string]struct{})
		}
		for _, hash := range hashes {
			a.authorHashes[name][hash] = struct{}{}
		}
	}
}

func blameAnalysis(g *git.Git, filename string) (map[string][]string, map[string]int) {
	blame := g.Blame(filename)
	commitsCnt := make(map[string]int)
	authorsHashes := make(map[string][]string)

	var (
		lastHash  string
		waitAuth  bool
		isNext    = true
		lineCount int
	)

	for _, text := range blame {
		parts := strings.Split(text, " ")

		if isNext {
			if lineCount == 0 {
				waitAuth = true
				lineCount, _ = strconv.Atoi(parts[len(parts)-1])
				lastHash = parts[0]
				commitsCnt[lastHash] += lineCount
			}
			lineCount--
			isNext = false
			continue
		}

		if text[0] != '\t' && waitAuth && g.CommitterOpt(parts[0]) {
			name := strings.Join(parts[1:], " ")
			authorsHashes[name] = append(authorsHashes[name], lastHash)
			waitAuth = false
		}

		if text[0] == '\t' {
			isNext = true
		}
	}

	if len(blame) == 0 {
		hash, author := g.Log(filename)
		commitsCnt[hash] = 0
		authorsHashes[author] = append(authorsHashes[author], hash)
	}

	return authorsHashes, commitsCnt
}

func (a *Analysis) normalize() {
	for name, hashes := range a.authorHashes {
		a.authors[name] = &author.Author{
			Name:    name,
			Lines:   0,
			Commits: len(hashes),
			Files:   a.fileCounts[name],
		}
		for hash := range hashes {
			a.authors[name].Lines += a.commitCounts[hash]
		}
	}
}

func (a *Analysis) compileResults() author.AuthorSlice {
	result := make([]author.Author, 0, len(a.authors))
	for _, author := range a.authors {
		result = append(result, *author)
	}
	return author.AuthorSlice{Authors: result}
}
