package predicate

import (
	"context"
	"fmt"
	"log"
	"strings"

	"golang.org/x/net/html"
)

type AndPredicate struct {
	Predicates []Predicate
}

func And(predicates ...Predicate) *AndPredicate {
	return &AndPredicate{predicates}
}

func (ap *AndPredicate) Match(n *html.Node) bool {
	for _, p := range ap.Predicates {
		if !p.Match(n) {
			return false
		}
	}

	return true
}

func (ap *AndPredicate) Search(ctx context.Context, n *html.Node) []*html.Node {
	nodesSeen := make(map[*html.Node]*[]bool)
	if Debug {
		log.Printf("And: Node %s", htmlNodeToString(n))
	}

	for i, p := range ap.Predicates {
		// Stop if the context is signalled to stop.
		if ctx.Err() != nil {
			break
		}

		// For each node found by the current predicate...
		nodesFound := p.Search(ctx, n)
		if Debug {
			nodesFoundString := make([]string, len(nodesFound))
			for _, nodeFound := range nodesFound {
				nodesFoundString = append(nodesFoundString, htmlNodeToString(nodeFound))
			}

			log.Printf("And: Node %s, Predicate %v, nodesFound %s", htmlNodeToString(n), p, strings.Join(nodesFoundString, " "))
		}

		for _, nodeFound := range nodesFound {
			// See if we've seen this node before
			nodePredicateMatch, exists := nodesSeen[nodeFound]

			if !exists {
				// We haven't, so allocate the bit array for it.
				nodePredicateMatchStorage := make([]bool, len(ap.Predicates))
				nodesSeen[nodeFound] = &nodePredicateMatchStorage
				nodePredicateMatch = &nodePredicateMatchStorage
			}

			// Set the bit for the current predicate.
			(*nodePredicateMatch)[i] = true
		}
	}

	// Build the list of nodes that matched all predicates.
	results := make([]*html.Node, 0, len(nodesSeen))

nextNode:
	for node, nodePredicateMatch := range nodesSeen {
		for i, isMatch := range *nodePredicateMatch {
			// Fail if any predicate didn't match.
			if !isMatch {
				if Debug {
					log.Printf("And: Node %s rejected by predicate %s", htmlNodeToString(node), ap.Predicates[i])
				}

				continue nextNode
			}
		}

		// All predicates matched; add the node to the results.
		results = append(results, node)
	}

	return results
}

func (ap *AndPredicate) String() string {
	preds := make([]string, len(ap.Predicates))
	for i, p := range ap.Predicates {
		preds[i] = fmt.Sprintf("%v", p)
	}

	return fmt.Sprintf("And(%s)", strings.Join(preds, ", "))
}

type OrPredicate struct {
	Predicates []Predicate
}

func Or(predicates ...Predicate) *OrPredicate {
	return &OrPredicate{predicates}
}

func (op *OrPredicate) Match(n *html.Node) bool {
	for _, p := range op.Predicates {
		if p.Match(n) {
			return true
		}
	}

	return false
}

func (op *OrPredicate) Search(ctx context.Context, n *html.Node) []*html.Node {
	nodesSeen := make(map[*html.Node]bool)

	if Debug {
		log.Printf("Or: Node %s", htmlNodeToString(n))
	}

	for _, p := range op.Predicates {
		// Stop if the context is signalled to stop.
		if ctx.Err() != nil {
			break
		}

		// For each node found by the current predicate...
		nodesFound := p.Search(ctx, n)
		if Debug {
			nodesFoundString := make([]string, len(nodesFound))
			for _, nodeFound := range nodesFound {
				nodesFoundString = append(nodesFoundString, htmlNodeToString(nodeFound))
			}

			log.Printf("Or: Node %s, Predicate %v, nodesFound %s", htmlNodeToString(n), p, strings.Join(nodesFoundString, " "))
		}

		for _, nodeFound := range nodesFound {
			// Add this node to the set of nodes we've seen.
			nodesSeen[nodeFound] = true
		}
	}

	// Build the list of nodes that matched any predicates.
	results := make([]*html.Node, 0, len(nodesSeen))

	for node := range nodesSeen {
		results = append(results, node)
	}

	return results
}

func (op *OrPredicate) String() string {
	preds := make([]string, len(op.Predicates))
	for i, p := range op.Predicates {
		preds[i] = fmt.Sprintf("%v", p)
	}

	return fmt.Sprintf("Or(%s)", strings.Join(preds, ", "))
}
