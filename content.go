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
