package xmltree

type Selector struct {
	Label, Space string
	Depth        int
}

// Match returns a slice of matching child Element(s)
// matching a search.
func (el *Element) Match(match *Selector) []*Element {
	var matches []*Element
	if el != nil {
		for i, child := range el.Children {
			if child.Type == XML_Tag &&
				child.Name.Local == match.Label &&
				(match.Space == "" || match.Space == child.Name.Space) {
				matches = append(matches, &el.Children[i])
			}
		}
	}
	return matches
}

// MatchOne returns a pointer to the first matching child of Element with a
// given match or nil if none matched.
func (el *Element) MatchOne(match *Selector) *Element {
	if el != nil {
		for i, child := range el.Children {
			if child.Type == XML_Tag &&
				child.Name.Local == match.Label &&
				(match.Space == "" || match.Space == child.Name.Space) {
				return &el.Children[i]
			}
		}
	}
	return nil
}

// FindOne returns a pointer to the first matching child of Element
// with a given match or nil if none matched.
func (el *Element) FindOne(match *Selector) *Element {
	if match.Depth > 0 {
		return el.matchAnyDeep(match, match.Depth)
	}
	return el.matchAnyDeep(match, recursionLimit)
}
func (el *Element) matchAnyDeep(match *Selector, d int) *Element {
	if el != nil {
		if e := el.MatchOne(match); e != nil {
			return e
		}
		if d > 0 {
			d--
			for i := range el.Children {
				if el.Children[i].Type == XML_Tag {
					if e := el.Children[i].matchAnyDeep(match, d-1); e != nil {
						return e
					}
				}
			}
		}
	}
	return nil
}
