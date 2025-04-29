// package formatter
//
// import (
//
//	"encoding/csv"
//	"encoding/json"
//	"io"
//	"strconv"
//	"strings"
//	"text/tabwriter"
//
//	"gitlab.com/slon/shad-go/gitfame/internal/author"
//
// )
//
//	type Formatter struct {
//		w       io.Writer
//		printer func(slice author.AuthorSlice, w io.Writer)
//	}
//
//	func New(t string, w io.Writer) *Formatter {
//		return &Formatter{
//			printer: getFunc(t),
//			w:       w,
//		}
//	}
//
//	func getFunc(format string) func(slice author.AuthorSlice, w io.Writer) {
//		if format == "tabular" {
//			return tabular
//		} else if format == "csv" {
//			return csvPrint
//		} else if format == "json" {
//			return jsonPrint
//		} else if format == "json-lines" {
//			return jsonLines
//		}
//		panic("There is no suitable format")
//	}
//
//	func tabular(a author.AuthorSlice, out io.Writer) {
//		w := tabwriter.NewWriter(out, 0, 0, 1, ' ', 0)
//		_, err := io.Copy(w, strings.NewReader("Name\tLines\tCommits\tFiles\n"))
//		if err != nil {
//			return
//		}
//
//		for _, auth := range a.Authors {
//			_, err = io.Copy(w, strings.NewReader(
//				auth.Name+"\t"+
//					strconv.Itoa(auth.Lines)+"\t"+
//					strconv.Itoa(auth.Commits)+"\t"+
//					strconv.Itoa(auth.Files)+"\n"))
//			if err != nil {
//				return
//			}
//		}
//		err = w.Flush()
//		if err != nil {
//			return
//		}
//	}
//
//	func csvPrint(a author.AuthorSlice, out io.Writer) {
//		w := csv.NewWriter(out)
//		err := w.Write([]string{"Name", "Lines", "Commits", "Files"})
//		if err != nil {
//			return
//		}
//
//		for _, auth := range a.Authors {
//			err = w.Write(
//				[]string{auth.Name,
//					strconv.Itoa(auth.Lines),
//					strconv.Itoa(auth.Commits),
//					strconv.Itoa(auth.Files)},
//			)
//			if err != nil {
//				return
//			}
//		}
//
//		w.Flush()
//	}
//
//	func jsonLines(a author.AuthorSlice, out io.Writer) {
//		for _, auth := range a.Authors {
//			b, _ := json.Marshal(auth)
//			_, err := out.Write(b)
//			if err != nil {
//				return
//			}
//			_, err = out.Write([]byte{'\n'})
//			if err != nil {
//				return
//			}
//		}
//	}
//
//	func jsonPrint(a author.AuthorSlice, out io.Writer) {
//		b, err := json.Marshal(a.Authors)
//		if err != nil {
//			return
//		}
//		_, err = out.Write(b)
//		if err != nil {
//			return
//		}
//	}
//
//	func (f *Formatter) Print(a author.AuthorSlice) {
//		f.printer(a, f.w)
//	}
package formatter

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"text/tabwriter"

	"gitlab.com/slon/shad-go/gitfame/internal/author"
)

type Formatter struct {
	w       io.Writer
	printer func(slice author.AuthorSlice, w io.Writer)
}

func New(t string, w io.Writer) *Formatter {
	return &Formatter{
		printer: getFunc(t),
		w:       w,
	}
}

func getFunc(format string) func(slice author.AuthorSlice, w io.Writer) {
	switch format {
	case "tabular":
		return tabular
	case "csv":
		return csvPrint
	case "json":
		return jsonPrint
	case "json-lines":
		return jsonLines
	default:
		panic("There is no suitable format")
	}
}

func tabular(a author.AuthorSlice, out io.Writer) {
	w := tabwriter.NewWriter(out, 0, 0, 1, ' ', 0)
	defer w.Flush()

	fmt.Fprintln(w, "Name\tLines\tCommits\tFiles")

	for _, auth := range a.Authors {
		fmt.Fprintf(w, "%s\t%d\t%d\t%d\n",
			auth.Name,
			auth.Lines,
			auth.Commits,
			auth.Files,
		)
	}
}

func csvPrint(a author.AuthorSlice, out io.Writer) {
	w := csv.NewWriter(out)
	defer w.Flush()

	w.Write([]string{"Name", "Lines", "Commits", "Files"})

	for _, auth := range a.Authors {
		w.Write([]string{
			auth.Name,
			strconv.Itoa(auth.Lines),
			strconv.Itoa(auth.Commits),
			strconv.Itoa(auth.Files),
		})
	}
}

func jsonLines(a author.AuthorSlice, out io.Writer) {
	enc := json.NewEncoder(out)
	for _, auth := range a.Authors {
		enc.Encode(auth)
	}
}

func jsonPrint(a author.AuthorSlice, out io.Writer) {
	json.NewEncoder(out).Encode(a.Authors)
}

func (f *Formatter) Print(a author.AuthorSlice) {
	f.printer(a, f.w)
}
