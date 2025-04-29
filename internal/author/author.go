package author

import (
	"sort"
	"strings"

	"gitlab.com/slon/shad-go/gitfame/internal/git"
)

func GetFiles(repo *git.Git) []string {
	filtered := repo.LsTree()
	filtered = repo.Exclude(filtered)
	filtered = repo.RestrictTo(filtered)
	return filtered
}

type Author struct {
	Name    string `json:"name"`
	Lines   int    `json:"lines"`
	Commits int    `json:"commits"`
	Files   int    `json:"files"`
}

type AuthorSlice struct {
	Authors []Author
	sortKey string
}

func (a AuthorSlice) Len() int      { return len(a.Authors) }
func (a AuthorSlice) Swap(i, j int) { a.Authors[i], a.Authors[j] = a.Authors[j], a.Authors[i] }

type priorityKey [3]int

func (a *AuthorSlice) getPriority(idx int) priorityKey {
	switch a.sortKey {
	case "lines":
		return priorityKey{a.Authors[idx].Lines, a.Authors[idx].Commits, a.Authors[idx].Files}
	case "files":
		return priorityKey{a.Authors[idx].Files, a.Authors[idx].Lines, a.Authors[idx].Commits}
	case "commits":
		return priorityKey{a.Authors[idx].Commits, a.Authors[idx].Lines, a.Authors[idx].Files}
	default:
		panic("invalid sort key: " + a.sortKey)
	}
}

func comparePriority(p1, p2 priorityKey, nameOrder bool) bool {
	for i := 0; i < len(p1); i++ {
		if p1[i] != p2[i] {
			return p1[i] > p2[i]
		}
	}
	return nameOrder
}

func (a AuthorSlice) Less(i, j int) bool {
	p1 := a.getPriority(i)
	p2 := a.getPriority(j)
	nameOrder := strings.ToLower(a.Authors[i].Name) < strings.ToLower(a.Authors[j].Name)
	return comparePriority(p1, p2, nameOrder)
}

func (a *AuthorSlice) Sort(order string) {
	a.sortKey = order
	sort.Sort(a)
}
