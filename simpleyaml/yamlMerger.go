package simpleyaml

import (
	"errors"
	"unsafe"
)

// Mapped list types
const (
	mlTypeValueAsKey = iota
	mlTypeKeyValue
)

// YamlMerger is the struct for merging YAML files
type YamlMerger struct {
	yamls       []*YamlNode               // YAMLs to merge
	finalYaml   *YamlNode                 // Result of YAMLs merge
	delTk       string                    // Deletion token
	dplMap      map[string]string         // Delimiter Per List map [list name => delimiter]
	mappedLists map[int]map[string]string // Mapped lists of the base file [node pointer => mapped list]
	strictMode  bool                      // Merge in strict mode
}

// NewMerger returns a new YamlMerger to merge X YAMLs.
//
// yamls		     YAMLs to merge (usually root nodes)
// deletionToken     Token to delete a node. e.g.: nil
// delimPerListMap   Delimiter per list map to identify key and value.
// strictMode        Merge in strict mode (Do not allow different node types)
func NewMerger(
	yamls []*YamlNode,
	deletionToken string,
	delimPerListMap map[string]string,
	strictMode bool,
) *YamlMerger {
	ym := new(YamlMerger)

	ym.yamls = yamls

	// Use the first YAML as the merge base
	ym.finalYaml = new(YamlNode)
	CopyNode(yamls[0], ym.finalYaml)

	ym.delTk = deletionToken
	ym.dplMap = delimPerListMap
	ym.mappedLists = make(map[int]map[string]string)
	ym.strictMode = strictMode

	return ym
}

// Merge returns the merged YAML.
func (ym *YamlMerger) Merge() (*YamlNode, error) {
	c := len(ym.yamls)

	for i := 1; i < c; i++ {
		childX := TraverseDown(ym.yamls[i])

		err := ym.mergeNodes(ym.finalYaml, childX)
		if err != nil {
			return nil, err
		}
	}

	return ym.finalYaml, nil
}

// mergeNodes merges recursively the childX node into the parent0 node.
func (ym *YamlMerger) mergeNodes(parent0 *YamlNode, childX *YamlNode) error {
	child0 := TraverseFindChild(parent0, childX.name)

	if child0 == nil {
		// Matching node not found in parent0, append childX
		newChild0 := NewChildNode(parent0)
		CopyNode(childX, newChild0)

		return ym.mergeNextNode(parent0, childX)
	}

	if childX.ntype == NodeTypeScalar && childX.values[0] == ym.delTk {
		// Deletion token as value, remove node
		RemoveChildNode(child0)

		return ym.mergeNextNode(parent0, childX)
	}

	if child0.ntype != childX.ntype {
		return ym.mergeDifferentNodeType(parent0, child0, childX)
	}

	if childX.ntype == NodeTypeScalar {
		child0.values = []string{childX.values[0]}
	} else if childX.ntype == NodeTypeList {
		err := ym.mergeNodesList(child0, childX)
		if err != nil {
			return err
		}
	}

	if childX.ntype == NodeTypeChildren {
		nextChildX := TraverseDown(childX)

		if nextChildX != nil {
			// Continue with first child of current node
			return ym.mergeNodes(child0, nextChildX)
		}
	}

	return ym.mergeNextNode(parent0, childX)
}

// mergeDifferentNodeType merges two different types of nodes.
func (ym *YamlMerger) mergeDifferentNodeType(
	parent0 *YamlNode,
	child0 *YamlNode,
	childX *YamlNode,
) error {
	if ym.strictMode {
		err := errors.New("Fatal error: [Strict Mode] Different node type found: " +
			child0.name + " / " + childX.name)
		return err
	}

	if child0.ntype == NodeTypeList {
		nptr := nodePointerToInt(child0)
		delete(ym.mappedLists, nptr)
	}

	// Overwrite child0 with childX
	CopyNode(childX, child0)

	return ym.mergeNextNode(parent0, childX)
}

// mergeNextNode gets the next node to continue the merge.
func (ym *YamlMerger) mergeNextNode(parent0 *YamlNode, childX *YamlNode) error {
	nextParent0 := parent0

	rewindDepth := new(uint)
	nextNodeX := TraverseNext(childX, rewindDepth)

	if *rewindDepth > 0 {
		nextParent0 = TraverseUpX(parent0, *rewindDepth)
	}

	if nextNodeX == nil {
		return nil
	}

	return ym.mergeNodes(nextParent0, nextNodeX)
}

// nodePointerToInt returns the pointer casted to integer for the given node.
func nodePointerToInt(node *YamlNode) int {
	ptr := unsafe.Pointer(node)

	return *((*int)(ptr))
}
