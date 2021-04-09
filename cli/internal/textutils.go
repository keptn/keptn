package internal

import (
	"fmt"
	"strings"
)

type node interface {
	String() string
}

func newInnerNode() *innerNode {
	return &innerNode{
		sortedKeys: []string{},
		children:   make(map[string]node),
	}
}

type innerNode struct {
	sortedKeys []string
	children   map[string]node
}

func (i *innerNode) String() string {
	var sb strings.Builder
	var entries = []string{}
	sb.WriteString("{")
	for _, e := range i.sortedKeys {
		entries = append(entries, `"`+e+`"`+":"+i.children[e].String())
	}

	s := strings.Join(entries, ",")
	sb.WriteString(s)
	sb.WriteString("}")
	return sb.String()
}

func (i *innerNode) AddChild(name string, c node) {
	i.children[name] = c
	i.sortedKeys = append(i.sortedKeys, name)
}

func (i *innerNode) GetChild(name string) node {
	n := i.children[name]
	return n
}

func newLeafNode() *leafnode {
	return &leafnode{}
}

type leafnode struct {
	value string
}

func (l *leafnode) String() string {
	return `"` + l.value + `"`
}

// JSONPathToJSONObj parse a slice of strings of the form "a.b.c=v" and
// returns a nested JSON of the form {"a":{"b":{"c":"v"}}}
func JSONPathToJSONObj(input []string) (string, error) {
	if input == nil {
		return "", fmt.Errorf("input is nil")
	}
	root := newInnerNode()
	for _, p := range input {
		pathValuePair := strings.Split(p, "=")
		if len(pathValuePair) != 2 {
			return "", fmt.Errorf("unable to parse input. Expected input: a.b.c=v")
		}
		pathParts := strings.Split(pathValuePair[0], ".")
		var currentNode = root
		for i := 0; i < len(pathParts)-1; i++ {
			node := currentNode.GetChild(pathParts[i])
			if node == nil {
				child := newInnerNode()
				currentNode.AddChild(pathParts[i], child)
				currentNode = child
			} else {
				currentNode = node.(*innerNode)
			}
		}
		leaf := newLeafNode()
		leaf.value = pathValuePair[1]
		currentNode.AddChild(pathParts[len(pathParts)-1], leaf)
	}

	return root.String(), nil
}
