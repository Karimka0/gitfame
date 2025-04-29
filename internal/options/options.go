package options

import (
	"github.com/spf13/pflag"
)

type Options struct {
	Repository   string
	Revision     string
	OrderBy      string
	UseCommitter bool
	OutFormat    string
	Extensions   []string
	Languages    []string
	Exclude      []string
	RestrictTo   []string
}

func ParseOptions(args []string) *Options {
	flagSet := pflag.NewFlagSet("gitfame", pflag.PanicOnError)

	var opt Options

	flagSet.StringVar(&opt.Repository, "repository", ".", "path")
	flagSet.StringVar(&opt.Revision, "revision", "HEAD", "commit")
	flagSet.StringVar(&opt.OrderBy, "order-by", "lines", "sorting-order[lines | commits | files]")
	flagSet.BoolVar(&opt.UseCommitter, "use-committer", false, "author-committer switching")
	flagSet.StringVar(&opt.OutFormat, "format", "tabular", "output format [tabular | csv | json | json-lines]")
	flagSet.StringSliceVar(&opt.Extensions, "extensions", nil, "Extensions")
	flagSet.StringSliceVar(&opt.Languages, "languages", nil, "Languages")
	flagSet.StringSliceVar(&opt.Exclude, "exclude", nil, "Exclude")
	flagSet.StringSliceVar(&opt.RestrictTo, "restrict-to", nil, "restrict to")

	if err := flagSet.Parse(args); err != nil {
		panic(err)
	}

	return &opt
}
