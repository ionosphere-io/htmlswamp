package predicate

import (
	"context"
	"log"

	"golang.org/x/net/html"
)

type HasSiblingWithPredicate struct {
	SiblingPredicate Predicate
}

func SiblingWith(p Predicate) *HasSiblingWithPredicate {
	return &HasSiblingWithPredicate{p}
}

func (hswp *HasSiblingWithPredicate) Match(n *html.Node) bool {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if hswp.SiblingPredicate.Match(c) {
			return true
		}
	}

	return false
}

func (hswp *HasSiblingWithPredicate) Search(ctx context.Context, n *html.Node) []*html.Node {
	var result []*html.Node

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if ctx.Err() != nil {
			return result
		}

		result = append(result, hswp.SiblingPredicate.Search(ctx, c)...)
	}

	return result
}

type HasChildWithPredicate struct {
	ChildPredicate Predicate
}

func ChildWith(p Predicate) *HasChildWithPredicate {
	return &HasChildWithPredicate{p}
}

func (hcwp *HasChildWithPredicate) Match(n *html.Node) bool {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if hcwp.ChildPredicate.Match(c) {
			return true
		}
	}

	return false
}

func (hcwp *HasChildWithPredicate) Search(ctx context.Context, n *html.Node) []*html.Node {
	var result []*html.Node

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if ctx.Err() != nil {
			return result
		}

		result = append(result, hcwp.ChildPredicate.Search(ctx, c)...)
	}

	return result
}

type HasDescendantWithPredicate struct {
	DescendantPredicate Predicate
}

func DescendantWith(p Predicate) *HasDescendantWithPredicate {
	return &HasDescendantWithPredicate{p}
}

func (hdwp *HasDescendantWithPredicate) Match(n *html.Node) bool {
	// Perform a breadth-first search for a match.
	var nodes []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		nodes = append(nodes, c)
	}

	for len(nodes) > 0 {
		node := nodes[0]
		nodes = nodes[1:]

		if hdwp.DescendantPredicate.Match(node) {
			return true
		}

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			nodes = append(nodes, c)
		}
	}

	return false
}

func (hdwp *HasDescendantWithPredicate) Search(ctx context.Context, n *html.Node) []*html.Node {
	// Perform a breadth-first search for a match.
	var result []*html.Node
	var nodes []*html.Node

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		nodes = append(nodes, c)
	}

	for ctx.Err() == nil && len(nodes) > 0 {
		node := nodes[0]
		nodes = nodes[1:]

		if Debug {
			log.Printf("HasDescendant: %s.Search(%s)", hdwp.DescendantPredicate, htmlNodeToString(node))
		}

		result = append(result, hdwp.DescendantPredicate.Search(ctx, node)...)

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if ctx.Err() != nil {
				return result
			}

			nodes = append(nodes, c)
		}
	}

	return result
}
