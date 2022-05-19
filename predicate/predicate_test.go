package predicate

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const testDocument = `
<html>
  <head>
    <title>Title of the page</title>
  </head>
  <body bgcolor=black style="flex;">
    <img src="http://example.com" class="class1 class2 class3">
	<table id="hello">
	  <thead>
	    <tr>
		  <th>Header 1</th>
		  <th>Header 2</th>
	    </tr>
	  </thead>
	  <tbody>
	    <tr>
		  <td class="label1 data">Row 1</td>
		  <td class="data label2">Cell 1</td>
		</tr>
	    <tr>
		  <td class="label1 data">Row 2</td>
		  <td class="data label2">Cell 2</td>
		</tr>
	  </tbody>
	</table>
  </body>
</html>
`

var doc *html.Node

func init() {
	var err error
	doc, err = html.Parse(strings.NewReader(testDocument))
	if err != nil {
		panic(err)
	}
}

func TestTags(t *testing.T) {
	assertions := assert.New(t)

	ctx := context.Background()
	results := DescendantWith(Tag("head")).Search(ctx, doc)
	if assertions.Len(results, 1, "Expected just one head") {
		assertions.Equal("head", results[0].Data)
	}

	results = DescendantWith(And(Tag("td"), Attribute("class", "label1"))).Search(ctx, doc)
	assertions.Len(results, 2, "Expected two label1 cells")
}

func TestRecursiveAbort(t *testing.T) {
	assertions := assert.New(t)

	root := &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Html,
		Data:     "html",
	}

	body := &html.Node{
		Parent:   root,
		Type:     html.ElementNode,
		DataAtom: atom.Body,
		Data:     "body",
	}
	root.FirstChild = body

	mainDiv := &html.Node{
		Parent:   body,
		Type:     html.ElementNode,
		DataAtom: atom.Div,
		Data:     "div",
		Attr:     []html.Attribute{{Key: "id", Val: "maindiv"}},
	}
	body.FirstChild = mainDiv

	innerDiv := &html.Node{
		Parent:   mainDiv,
		Type:     html.ElementNode,
		DataAtom: atom.Div,
		Data:     "div",
		Attr:     []html.Attribute{{Key: "id", Val: "innerdiv"}},
	}
	mainDiv.FirstChild = innerDiv

	// Loop innerDiv on itself
	innerDiv.FirstChild = innerDiv

	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(10*time.Millisecond))

	t.Cleanup(func() { cancel() })

	results := DescendantWith(And(Tag("div"), Attribute("id", "maindiv"))).Search(ctx, root)
	assertions.Len(results, 1, "Expected just one div")
}
