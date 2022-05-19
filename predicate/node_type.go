package predicate

import (
	"context"

	"golang.org/x/net/html"
)

type IsNodeTypePredicate struct {
	NodeType html.NodeType
}

func (isntp *IsNodeTypePredicate) Match(n *html.Node) bool {
	return n.Type == isntp.NodeType
}

func (isntp *IsNodeTypePredicate) Search(_ context.Context, n *html.Node) []*html.Node {
	if isntp.Match(n) {
		return []*html.Node{n}
	}

	return nil
}
