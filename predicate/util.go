package predicate

import (
	"fmt"

	"golang.org/x/net/html"
)

var Debug = false

func nodeTypeToString(nodeType html.NodeType) string {
	switch nodeType {
	case html.ErrorNode:
		return "ErrorNode"
	case html.TextNode:
		return "TextNode"
	case html.DocumentNode:
		return "DocumentNode"
	case html.ElementNode:
		return "ElementNode"
	case html.CommentNode:
		return "CommentNode"
	case html.DoctypeNode:
		return "DoctypeNode"
	case html.RawNode:
		return "RawNow"
	default:
		return fmt.Sprintf("UnknownNodeType(%d)", nodeType)
	}
}

func htmlNodeToString(n *html.Node) string {
	switch n.Type {
	case html.ElementNode:
		result := fmt.Sprintf("<%s", n.Data)
		for _, a := range n.Attr {
			result += fmt.Sprintf(" %s=\"%s\"", a.Key, a.Val)
		}
		result += ">"
		return result

	default:
		data := n.Data
		if len(data) > 20 {
			data = data[:20] + "..."
		}
		return fmt.Sprintf("%s(%q)", nodeTypeToString(n.Type), data)
	}
}
