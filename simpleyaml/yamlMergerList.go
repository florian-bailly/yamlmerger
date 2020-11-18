package simpleyaml

import (
	"fmt"
	"strings"
)

// Tokens of raw "Delimiter Per List" format
const (
	tkRawDelimPerListPostKey   = ':'
	tkRawDelimPerListPostValue = ','
)

// mergeNodesList merges the lists of the given two nodes.
func (ym *YamlMerger) mergeNodesList(node *YamlNode, node2 *YamlNode) error {
	var mlistMerged map[string]string
	var err error

	delim, isDelimited := ym.dplMap[node2.name]

	nptr := nodePointerToInt(node)
	mlist, _ := ym.mappedLists[nptr]

	if isDelimited {
		mlistMerged, err = mergeNodesListDelimited(node, node2, delim, &mlist, ym.delTk)
	} else {
		mlistMerged = mergeNodesListBasic(node, node2, &mlist, ym.delTk)
	}

	node.values = mappedListToList(mlistMerged, delim)

	ym.mappedLists[nptr] = mlistMerged

	return err
}

// mergeNodesListDelimited returns a map corresponding to the merge
// of the delimited lists for the given two nodes.
// delim          Delimiter to use
// baseMappedList Mapped list to use as a base for the merge
// delTk          Deletion token
func mergeNodesListDelimited(
	node *YamlNode,
	node2 *YamlNode,
	delim string,
	baseMappedList *map[string]string,
	delTk string,
) (map[string]string, error) {
	var mlist map[string]string

	if len(*baseMappedList) == 0 {
		var err error
		mlist, err = delimitedListToMappedList(node.values, delim)
		if err != nil {
			return nil, err
		}
	} else {
		mlist = *baseMappedList
	}

	mlist2, err := delimitedListToMappedList(node2.values, delim)
	if err != nil {
		return nil, err
	}

	return mergeMappedListsDelimited(mlist, mlist2, delim, delTk), nil
}

// mergeNodesListBasic returns a map corresponding to the merge of the lists
// for the given two nodes.
// baseMappedList Mapped list to use as a base for the merge
// delTk          Deletion token
func mergeNodesListBasic(
	node *YamlNode,
	node2 *YamlNode,
	baseMappedList *map[string]string,
	delTk string,
) map[string]string {
	var mlist map[string]string

	if len(*baseMappedList) == 0 {
		mlist = listToMappedList(node.values)
	} else {
		mlist = *baseMappedList
	}

	mlist2 := listToMappedList(node2.values)

	return mergeMappedListsBasic(mlist, mlist2, delTk)
}

// mergeMappedListsDelimited returns a map corresponding to the merge of the lists
// for the given two delimited mapped lists.
// delim Delimiter to use
// delTk Deletion token
func mergeMappedListsDelimited(
	mlist map[string]string,
	mlist2 map[string]string,
	delim string,
	delTk string,
) map[string]string {
	mlistMerged := mlist

	for k2, v2 := range mlist2 {
		if v2 == delTk {
			delete(mlistMerged, k2)
			continue
		}

		mlistMerged[k2] = v2
	}

	return mlistMerged
}

// mergeMappedListsBasic returns a map corresponding to the merge of the lists
// for the given two mapped lists.
// delTk Deletion token. If the value ends with a colon followed by your token
//                       then it will be deleted.
func mergeMappedListsBasic(
	mlist map[string]string,
	mlist2 map[string]string,
	delTk string,
) map[string]string {
	mlistMerged := mlist
	tkLen := len(delTk)

	for v2 := range mlist2 {
		v2Len := len(v2)
		if delTk != "" && (v2Len > tkLen+1 && v2[v2Len-(tkLen+1):] == ":"+delTk) {
			delete(mlistMerged, v2[0:v2Len-(tkLen+1)])
			continue
		}

		mlistMerged[v2] = ""
	}

	return mlistMerged
}

// delimitedListToMappedList returns a mapped list which is essentially
// the list items splitted by the given delimiter.
// delim Delimiter to use
func delimitedListToMappedList(list []string, delim string) (map[string]string, error) {
	mp := make(map[string]string)

	c := len(list)
	for i := 0; i < c; i++ {
		split := strings.SplitN(list[i], delim, 2)
		if len(split) < 2 {
			return nil, fmt.Errorf("Malformed list item `%s`, delimiter `%s` not found",
				list[i], delim)
		}

		mp[split[0]] = split[1]
	}

	return mp, nil
}

// listToMappedList returns a mapped list which is essentially
// the list items as key with no value.
func listToMappedList(list []string) map[string]string {
	mp := make(map[string]string)

	c := len(list)
	for i := 0; i < c; i++ {
		mp[list[i]] = ""
	}

	return mp
}

// mappedListToList returns a list from the given mapped list.
// delim Delimiter to use; if basic list set it to empty
func mappedListToList(mlist map[string]string, delim string) []string {
	list := []string{}

	for k, v := range mlist {
		list = append(list, k+delim+v)
	}

	return list
}

// RawDelimPerListToMap returns the map of the given raw "Delimiter Per List" format.
// Raw format: listName1:delim1,listName2:delim2[,...]
func RawDelimPerListToMap(str string) map[string]string {
	dplMap := make(map[string]string)

	if str == "" {
		return dplMap
	}

	if str[len(str)-1] != tkRawDelimPerListPostValue {
		str += ","
	}

	var key, value string
	var prevChar rune
	var keyStartIndex, keyEndIndex int

	for i, char := range str {
		if keyEndIndex == 0 {
			// Get key
			if char == tkRawDelimPerListPostKey && prevChar != '\\' {
				key = str[keyStartIndex:i]
				// Handle token escape (if any)
				key = strings.Replace(key,
					string('\\'+tkRawDelimPerListPostKey),
					string(tkRawDelimPerListPostKey), -1)

				keyEndIndex = i + 1
			}
		} else {
			// Get value
			if char == tkRawDelimPerListPostValue && prevChar != '\\' {
				value = str[keyEndIndex:i]
				// Handle token escape (if any)
				value = strings.Replace(value,
					string('\\'+tkRawDelimPerListPostValue),
					string(tkRawDelimPerListPostValue), -1)

				dplMap[key] = value

				keyStartIndex = i + 1
				keyEndIndex = 0
			}
		}

		prevChar = char
	}

	return dplMap
}
