package predicate

import (
	"context"

	"golang.org/x/net/html"
)

type Predicate interface {
	Match(n *html.Node) bool
	Search(ctx context.Context, n *html.Node) []*html.Node
}
