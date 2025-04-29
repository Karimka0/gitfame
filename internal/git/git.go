package git

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gitlab.com/slon/shad-go/gitfame/internal/options"
)

var (
	logCmd      = []string{"git", "log", "--pretty=format:%H %an"}
	blameCmd    = []string{"git", "blame", "--porcelain"}
	treeCmd     = []string{"git", "ls-tree", "-r", "--name-only"}
	langExtPath = "../../configs/language_extensions.json"
)

type Git struct {
	opt *options.Options
}

func NewGit(opt *options.Options) *Git {
	return &Git{opt: opt}
}

func (g *Git) cmdSplit(cmd *exec.Cmd) []string {
	cmd.Dir = g.opt.Repository
	b, err := cmd.Output()
	s := string(b)
	if err != nil {
		panic(err)
	}
	return strings.FieldsFunc(s, func(r rune) bool {
		return r == '\n'
	})
}

func (g *Git) Log(fileName string) (string, string) {
	split := g.cmdSplit(exec.Command(logCmd[0], logCmd[1], logCmd[2], g.opt.Revision, "--", fileName))
	hash := strings.Split(split[0], " ")
	return hash[0], strings.Join(hash[1:], " ")
}

func (g *Git) Blame(name string) []string {
	return g.cmdSplit(exec.Command(blameCmd[0], blameCmd[1], blameCmd[2], g.opt.Revision, name))
}

func (g *Git) LsTree() []string {
	ans := g.cmdSplit(exec.Command(treeCmd[0], treeCmd[1], treeCmd[2], treeCmd[3], g.opt.Revision))
	languages := g.opt.Languages
	extensions := g.opt.Extensions
	if len(languages) != 0 {
		extensions = extUpdate(languages, extensions)
	}
	return extFilter(ans, extensions)
}

func (g *Git) CommitterOpt(auth string) bool {
	return (!g.opt.UseCommitter && auth == "author") || (g.opt.UseCommitter && auth == "committer")
}

func FileFiltering(fs []string, filt func(file string) bool) []string {
	newFs := make([]string, 0)
	for _, f := range fs {
		if filt(f) {
			newFs = append(newFs, f)
		}
	}
	return newFs
}

func (g *Git) RestrictTo(files []string) []string {
	restrictTo := g.opt.RestrictTo
	if len(restrictTo) == 0 {
		return files
	}
	return FileFiltering(files, func(file string) bool {
		for _, rest := range restrictTo {
			if b, _ := filepath.Match(rest, file); b {
				return true
			}
		}
		return false
	})
}

func (g *Git) Exclude(files []string) []string {
	exclude := g.opt.Exclude
	if len(exclude) == 0 {
		return files
	}
	return FileFiltering(files, func(file string) bool {
		for _, excl := range exclude {
			if b, _ := filepath.Match(excl, file); b {
				return false
			}
		}
		return true
	})
}

func readLanguages() map[string][]string {
	type Lang struct {
		Name       string
		Type       string
		Extensions []string
	}

	var ls []Lang
	b, _ := os.ReadFile(langExtPath)

	err := json.Unmarshal(b, &ls)
	if err != nil {
		panic(err)
	}

	ans := make(map[string][]string)

	for _, l := range ls {
		ans[strings.ToLower(l.Name)] = l.Extensions
	}
	return ans
}

func extUpdate(languages []string, extensions []string) []string {
	langExtMp := readLanguages()

	for _, lang := range languages {
		if exts, ok := langExtMp[strings.ToLower(lang)]; ok {
			extensions = append(extensions, exts...)
		}
	}
	return extensions
}

func setify(extensions []string) (set map[string]struct{}) {
	set = make(map[string]struct{})
	for _, ext := range extensions {
		set[ext] = struct{}{}
	}
	return
}

func extFilter(oldAns, extensions []string) []string {
	if len(extensions) == 0 {
		return oldAns
	}

	set := setify(extensions)

	return FileFiltering(oldAns, func(file string) bool {
		_, ok := set[filepath.Ext(file)]
		return ok
	})
}
