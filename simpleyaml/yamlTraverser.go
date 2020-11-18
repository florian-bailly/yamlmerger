package simpleyaml

// TraverseUp returns the parent node
func TraverseUp(node *YamlNode) *YamlNode {
	return node.parent
}

// TraverseUpX returns the parent X node
func TraverseUpX(node *YamlNode, x uint) *YamlNode {
	nodeUp := node

	for i := 0; i < int(x); i++ {
		nodeUp = TraverseUp(nodeUp)
		if nodeUp == nil {
			break
		}
	}

	return nodeUp
}

// TraverseDown returns the first child node
func TraverseDown(node *YamlNode) *YamlNode {
	if len(node.children) == 0 {
		return nil
	}

	return node.children[0]
}

// TraverseNext returns the closest next node (either sibling or next parent child)
func TraverseNext(node *YamlNode, rewindDepth *uint) *YamlNode {
	var nextChild *YamlNode

	nodeNext := node
	*rewindDepth = 0

	for {
		nextChild = TraverseNextChild(nodeNext)
		if nextChild != nil {
			break
		}

		nodeNext = nodeNext.parent
		if nodeNext == nil {
			// Root
			break
		}

		*rewindDepth++
	}

	return nextChild
}

// TraverseNextChild returns the next child
func TraverseNextChild(node *YamlNode) *YamlNode {
	if node.parent == nil {
		// Root
		return nil
	}

	var nextChild *YamlNode
	c := len(node.parent.children)

	for i := 0; i < c; i++ {
		if node.parent.children[i] == node {
			if i == c-1 {
				// No next child
				break
			}

			nextChild = node.parent.children[i+1]
			break
		}
	}

	return nextChild
}

// TraverseFindChild returns the child node matching the given name
func TraverseFindChild(parentNode *YamlNode, name string) *YamlNode {
	var child *YamlNode
	c := len(parentNode.children)

	for i := 0; i < c; i++ {
		if parentNode.children[i].name == name {
			child = parentNode.children[i]
			break
		}
	}

	return child
}

// TraverseFindSibling returns the sibling node matching the given name
func TraverseFindSibling(node *YamlNode, name string) *YamlNode {
	if node.parent == nil {
		return nil
	}

	return TraverseFindChild(node.parent, name)
}
