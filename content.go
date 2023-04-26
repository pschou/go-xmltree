package xmltree

// Returns string of content if available
func (el *Element) GetContent() string {
	if el != nil {
		return el.Content
	}
	return ""
}

// Sets the content if available
func (el *Element) SetContent(val string) {
	if el != nil {
		el.Content = val
	}
}

// RemoveEmpty cleans up the tree of any empty elements.
func (el *Element) RemoveEmpty() {
	el.removeEmptyElms(make(map[*Element]struct{}))
}

func (el *Element) removeEmptyElms(visited map[*Element]struct{}) {
	if el == nil {
		return
	}
	if _, ok := visited[el]; ok {
		// We have a cycle.
		return
	}
	if len(el.Content) > 0 {
		return
	}
	var keep []Element
	visited[el] = struct{}{}
	for i := range el.Children {
		child := el.Children[i]
		child.removeEmptyElms(visited)
		if len(child.Content) > 0 || len(child.Children) > 0 || len(child.StartElement.Attr) > 0 {
			keep = append(keep, child)
		}
	}
	delete(visited, el)
	el.Children = keep
}
