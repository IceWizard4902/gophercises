package link

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

// Link represents <a href="..."> in a HTML document
type Link struct {
	Href string
	Text string
}

// NodeType in x/net/html
// TextNode         any node that only contains text
// DocumentNode     the parent root node (immediate child is the html tag)
// ElementNode      HTML tag <h1>, <div>, ...
// CommentNode		for HTML comments
// DoctypeNode		<!DOCTYPE html>

// Parse will take in a HTML document and will return a slice
// of links parsed from it.
func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	nodes := linkNodes(doc)
	var links []Link
	for _, node := range nodes {
		links = append(links, buildLink(node))
	}

	return links, nil
}

func text(n *html.Node) string {
	if n.Type == html.TextNode {
		// probably no children so just return node
		return n.Data
	}
	if n.Type != html.ElementNode {
		return ""
	}
	var ret string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		// This is not the optimal way to create strings
		// Can look into byte buffer to build strings in optimized way
		ret += text(c)
	}
	// Same as .split and .join in Python syntax
	return strings.Join(strings.Fields(ret), " ")
}

func linkNodes(n *html.Node) []*html.Node {
	// Base case
	if (n.Type == html.ElementNode) && (n.Data == "a") {
		return []*html.Node{n}
	}
	var ret []*html.Node

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		// ... to expand the variadic parameter, as the return type is a slice
		// but we only want a single Node at a time to append
		ret = append(ret, linkNodes(c)...)
	}
	return ret
}

func buildLink(n *html.Node) Link {
	var ret Link
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			ret.Href = attr.Val
			break
		}
	}
	ret.Text = text(n)
	return ret
}
