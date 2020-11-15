package simpleyaml

// YAML Tokens
const (
	TkPostKey      = ":"
	TkPreListValue = "- "
	TkComment      = "#"
	TkStringDelim1 = "\""
	TkStringDelim2 = "'"
)

// YamlNode is a YAML node (duh)
type YamlNode struct {
	name     string
	values   []string
	ntype    uint
	children []*YamlNode // Slice of pointers to avoid pointer reset when appending (because of parent ref)
	parent   *YamlNode
}

// Node Types
const (
	NodeTypeChildren uint = iota
	NodeTypeScalar
	NodeTypeList
)

// NewChildNode returns the node pointer of the new child created.
func NewChildNode(node *YamlNode) *YamlNode {
	newNode := new(YamlNode)
	newNode.parent = node

	node.children = append(node.children, newNode)

	return node.children[len(node.children)-1]
}

// CreateRootNode returns the node pointer of a "virtual" root node.
func CreateRootNode() YamlNode {
	rootNode := YamlNode{}
	rootNode.name = "root"
	rootNode.ntype = NodeTypeChildren

	return rootNode
}

// RemoveChildNode removes a child from its parent.
func RemoveChildNode(child *YamlNode) {
	var newChildren []*YamlNode
	c := len(child.parent.children)

	for i := 0; i < c; i++ {
		child2 := child.parent.children[i]
		if child2 == child {
			continue
		}
		newChildren = append(newChildren, child2)
	}

	child.parent.children = newChildren
}

// CopyNode makes a deep copy of the given node.
func CopyNode(node *YamlNode, destNode *YamlNode) {
	destNode.name = node.name
	destNode.ntype = node.ntype
	destNode.values = node.values

	c := len(node.children)

	for i := 0; i < c; i++ {
		childNode := NewChildNode(destNode)
		CopyNode(node.children[i], childNode)
	}
}
