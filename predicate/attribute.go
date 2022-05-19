package predicate

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

func Attribute(args ...interface{}) Predicate {
	if len(args) == 1 {
		if keyString, ok := args[0].(string); ok {
			return &HasAttributePredicate{keyString}
		}

		if keyRegexp, ok := args[0].(*regexp.Regexp); ok {
			return &HasAttributeRegexpPredicate{keyRegexp}
		}

		panic(fmt.Sprintf("Cannot create an attribute from type %T", args[0]))
	}

	if len(args) == 2 {
		if keyString, ok := args[0].(string); ok {
			if valueString, ok := args[1].(string); ok {
				return &HasAttributeKeyAnyValuesPredicate{keyString, []string{valueString}}
			}

			if valueRegexp, ok := args[1].(*regexp.Regexp); ok {
				keyRegexp := regexp.MustCompile(fmt.Sprintf("^$%s$", regexp.QuoteMeta(keyString)))
				return &HasAttributeKeyValueRegexpPredicate{keyRegexp, valueRegexp}
			}

			panic(fmt.Sprintf("Cannot create an attribute from types %T %T", args[0], args[1]))
		}

		if keyRegexp, ok := args[0].(*regexp.Regexp); ok {
			if valueString, ok := args[1].(string); ok {
				valueRegexp := regexp.MustCompile(fmt.Sprintf("^%s$", regexp.QuoteMeta(valueString)))
				return &HasAttributeKeyValueRegexpPredicate{keyRegexp, valueRegexp}
			}

			if valueRegexp, ok := args[1].(*regexp.Regexp); ok {
				return &HasAttributeKeyValueRegexpPredicate{keyRegexp, valueRegexp}
			}

			panic(fmt.Sprintf("Cannot create an attribute from types %T %T", args[0], args[1]))
		}

		panic(fmt.Sprintf("Cannot create an attribute from types %T %T", args[0], args[1]))
	}

	var key string
	values := make([]string, 0, len(args)-1)

	for i, arg := range args {
		argString, ok := arg.(string)

		if !ok {
			var types []string

			for _, arg := range args {
				types = append(types, fmt.Sprintf("%T", arg))
			}

			panic(fmt.Sprintf("Cannot create an attribute from types %s", strings.Join(types, ", ")))
		}

		if i == 0 {
			key = argString
		} else {
			values = append(values, argString)
		}
	}

	return &HasAttributeKeyAnyValuesPredicate{key, values}
}

type HasAttributePredicate struct {
	Key string
}

func (hap *HasAttributePredicate) Match(n *html.Node) bool {
	for _, attr := range n.Attr {
		if attr.Key == hap.Key {
			return true
		}
	}

	return false
}

func (hap *HasAttributePredicate) Search(_ context.Context, n *html.Node) []*html.Node {
	if hap.Match(n) {
		return []*html.Node{n}
	}

	return nil
}

func (hap *HasAttributePredicate) String() string {
	return fmt.Sprintf("HasAttribute(%s)", hap.Key)
}

type HasAttributeRegexpPredicate struct {
	KeyRegexp *regexp.Regexp
}

func (harp *HasAttributeRegexpPredicate) Match(n *html.Node) bool {
	for _, attr := range n.Attr {
		if harp.KeyRegexp.MatchString(attr.Key) {
			return true
		}
	}

	return false
}

func (harp *HasAttributeRegexpPredicate) Search(_ context.Context, n *html.Node) []*html.Node {
	if harp.Match(n) {
		return []*html.Node{n}
	}

	return nil
}

func (harp *HasAttributeRegexpPredicate) String() string {
	return fmt.Sprintf("HasAttributeRegexp(%s)", harp.KeyRegexp.String())
}

type HasAttributeKeyExactValuePredicate struct {
	Key   string
	Value string
}

func (hakvp *HasAttributeKeyExactValuePredicate) Match(n *html.Node) bool {
	for _, attr := range n.Attr {
		if attr.Key == hakvp.Key && attr.Val == hakvp.Value {
			return true
		}
	}

	return false
}

func (hakvp *HasAttributeKeyExactValuePredicate) Search(_ context.Context, n *html.Node) []*html.Node {
	if hakvp.Match(n) {
		return []*html.Node{n}
	}

	return nil
}

func (hakvp *HasAttributeKeyExactValuePredicate) String() string {
	return fmt.Sprintf("HasAttributeKeyValue(%s, %s)", hakvp.Key, hakvp.Value)
}

var whitespace *regexp.Regexp

func init() {
	whitespace = regexp.MustCompile(`\s+`)
}

type HasAttributeKeyAllValuesPredicate struct {
	Key    string
	Values []string
}

func (hakavp *HasAttributeKeyAllValuesPredicate) Match(n *html.Node) bool {
attributeLoop:

	for _, attr := range n.Attr {
		if attr.Key == hakavp.Key {
			nodeValues := whitespace.Split(attr.Val, -1)

		valueLoop:
			for _, wantedValue := range hakavp.Values {
				for _, nodeValue := range nodeValues {
					if wantedValue == nodeValue {
						continue valueLoop
					}
				}

				// Not found -- continue to next attribute.
				continue attributeLoop
			}

			// All values found in this attribute; return true.
			return true
		}
	}

	return false
}

func (hakavp *HasAttributeKeyAllValuesPredicate) Search(_ context.Context, n *html.Node) []*html.Node {
	if hakavp.Match(n) {
		return []*html.Node{n}
	}

	return nil
}

func (hakavp *HasAttributeKeyAllValuesPredicate) String() string {
	return fmt.Sprintf("HasAttributeKeyAllValues(%s, [%s])", hakavp.Key, strings.Join(hakavp.Values, ", "))
}

type HasAttributeKeyAnyValuesPredicate struct {
	Key    string
	Values []string
}

func (hakavp *HasAttributeKeyAnyValuesPredicate) Match(n *html.Node) bool {
	for _, attr := range n.Attr {
		if attr.Key == hakavp.Key {
			nodeValues := whitespace.Split(attr.Val, -1)

			for _, wantedValue := range hakavp.Values {
				for _, nodeValue := range nodeValues {
					if wantedValue == nodeValue {
						return true
					}
				}
			}
		}
	}

	return false
}

func (hakavp *HasAttributeKeyAnyValuesPredicate) Search(_ context.Context, n *html.Node) []*html.Node {
	if hakavp.Match(n) {
		return []*html.Node{n}
	}

	return nil
}

func (hakavp *HasAttributeKeyAnyValuesPredicate) String() string {
	return fmt.Sprintf("HasAttributeKeyAnyValues(%s, [%s])", hakavp.Key, strings.Join(hakavp.Values, ", "))
}

type HasAttributeKeyValueRegexpPredicate struct {
	KeyRegexp   *regexp.Regexp
	ValueRegexp *regexp.Regexp
}

func (hakvrp *HasAttributeKeyValueRegexpPredicate) Match(n *html.Node) bool {
	for _, attr := range n.Attr {
		if hakvrp.KeyRegexp.MatchString(attr.Key) && hakvrp.ValueRegexp.MatchString(attr.Val) {
			return true
		}
	}

	return false
}

func (hakvrp *HasAttributeKeyValueRegexpPredicate) Search(_ context.Context, n *html.Node) []*html.Node {
	if hakvrp.Match(n) {
		return []*html.Node{n}
	}

	return nil
}

func (hakvrp *HasAttributeKeyValueRegexpPredicate) String() string {
	return fmt.Sprintf("HasAttributeKeyValueRegexp(%s, %s)", hakvrp.KeyRegexp.String(), hakvrp.ValueRegexp.String())
}
