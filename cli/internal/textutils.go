package internal

import (
	"encoding/json"
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
				if cn, ok := node.(*innerNode); ok {
					currentNode = cn
				} else {
					return "", fmt.Errorf("unable to map path part %s as the value is already bound to %s", pathParts[i], node.(*leafnode).value)
				}
			}
		}
		leaf := newLeafNode()
		leaf.value = pathValuePair[1]
		currentNode.AddChild(pathParts[len(pathParts)-1], leaf)
	}
	return root.String(), nil
}

// UnfoldToMap takes a map of with keys of the form "a.b.c" and a value "v" each
// and returns an "unfold" map containing map[a[b[c]]] = v
func UnfoldMap(inMap map[string]string) (map[string]interface{}, error) {
	if inMap == nil {
		return map[string]interface{}{}, nil
	}
	var transformed []string
	for path, value := range inMap {
		transformed = append(transformed, path+"="+value)
	}
	s, err := JSONPathToJSONObj(transformed)
	if err != nil {
		return map[string]interface{}{}, err
	}

	var res map[string]interface{}
	err = json.Unmarshal([]byte(s), &res)
	if err != nil {
		return map[string]interface{}{}, err
	}

	return res, nil
}
