package simpleyaml

import (
	"os"
	"strings"
)

// YamlWriter is the struct for writing YAML files.
type YamlWriter struct {
	file     *os.File
	rootNode YamlNode
}

// Output settings
const (
	OutputIndent uint = 2
)

// NewWriter returns a new yamlWriter to be used for YAML writing.
func NewWriter(file *os.File) *YamlWriter {
	yw := new(YamlWriter)

	yw.file = file

	return yw
}

// Write formats the given YAML tree into the output file.
func (yw *YamlWriter) Write(yaml *YamlNode) {
	yw.rootNode = *yaml
	// Write directly root children (as the root is "virtual")
	yw.writeNodeChildren(yaml, 0)
}

// writeNode formats the given node (recursively) into the output file.
// indent Indentation level to use (herited from recursivity)
func (yw *YamlWriter) writeNode(node *YamlNode, indent uint) {
	var data string

	data = strings.Repeat(" ", int(OutputIndent*indent)) +
		node.name + TkPostKey

	if node.ntype == NodeTypeScalar {
		data += " " + node.values[0]
	} else if node.ntype == NodeTypeList {
		for j := 0; j < len(node.values); j++ {
			data += "\n" +
				strings.Repeat(" ", int(OutputIndent*(indent+1))) +
				TkPreListValue + node.values[j]
		}
	}
	data += "\n"

	yw.file.Write([]byte(data))

	yw.writeNodeChildren(node, indent+1)
}

// writeNodeChildren formats the child nodes of the given node (recursively).
// indent Indentation level to use (herited from recursivity)
func (yw *YamlWriter) writeNodeChildren(node *YamlNode, indent uint) {
	var childNode *YamlNode

	for i := 0; i < len(node.children); i++ {
		childNode = node.children[i]

		yw.writeNode(childNode, indent)
	}
}
