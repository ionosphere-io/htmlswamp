package predicate

import (
	"context"
	"fmt"
	"regexp"

	"golang.org/x/net/html"
)

type TagPredicate struct {
	Tag string
}

func Tag(tag string) *TagPredicate {
	return &TagPredicate{tag}
}

func TagRegexp(tagRegexp *regexp.Regexp) *TagRegexpPredicate {
	return &TagRegexpPredicate{tagRegexp}
}

func (tp *TagPredicate) Match(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == tp.Tag
}

func (tp *TagPredicate) Search(_ context.Context, n *html.Node) []*html.Node {
	if tp.Match(n) {
		return []*html.Node{n}
	}

	return nil
}

func (tp *TagPredicate) String() string {
	return fmt.Sprintf("Tag(%s)", tp.Tag)
}

type TagRegexpPredicate struct {
	TagRegexp *regexp.Regexp
}

func (trp *TagRegexpPredicate) Match(n *html.Node) bool {
	return n.Type == html.ElementNode && trp.TagRegexp.MatchString(n.Data)
}

func (trp *TagRegexpPredicate) Search(_ context.Context, n *html.Node) []*html.Node {
	if trp.Match(n) {
		return []*html.Node{n}
	}

	return nil
}

func (trp *TagRegexpPredicate) String() string {
	return fmt.Sprintf("TagRegexp(%s)", trp.TagRegexp.String())
}
