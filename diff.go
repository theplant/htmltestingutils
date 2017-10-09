package htmltestingutils

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/andybalholm/cascadia"
	"github.com/pmezard/go-difflib/difflib"
	"golang.org/x/net/html"

	"github.com/yosssi/gohtml"
)

func PrettyHtmlDiff(actual io.Reader, actualCssSelector string, expected string) (r string) {
	buf := bytes.NewBuffer(nil)
	io.Copy(buf, actual)
	fexpected := gohtml.Format(expected)

	sel, err := cascadia.Compile(actualCssSelector)
	if err != nil {
		panic(err)
	}

	n, err := html.Parse(buf)
	if err != nil {
		panic(err)
	}
	mn := sel.MatchFirst(n)
	if mn == nil {
		panic(fmt.Sprintf("css selector '%s' not found in html:\n%s", actualCssSelector, buf.String()))
	}
	selBuf := bytes.NewBuffer(nil)
	html.Render(selBuf, mn)

	trimmedBuf := bytes.NewBuffer(nil)
	lines := strings.Split(selBuf.String(), "\n")
	for _, l := range lines {
		trimmedBuf.WriteString(strings.TrimSpace(l) + "\n")
	}

	factual := gohtml.Format(trimmedBuf.String())
	if fexpected != factual {
		diff := difflib.UnifiedDiff{
			A:        difflib.SplitLines(fexpected),
			B:        difflib.SplitLines(factual),
			FromFile: "Expected",
			ToFile:   "Actual",
			Context:  3,
		}
		r, _ = difflib.GetUnifiedDiffString(diff)
	}

	return
}
