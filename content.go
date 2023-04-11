package xmltree

// Returns string of content if available
func (el *Element) GetContent() string {
	var matches []*Element
	if el != nil {
		return el.Content
	}
	return ""
}

// Sets the content if available
func (el *Element) SetContent(val string) {
	var matches []*Element
	if el != nil {
		el.Content = val
	}
}
