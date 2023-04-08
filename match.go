package xmltree

import (
	"encoding/xml"
)

type Selector struct {
	Label, Space string
	Depth        int
	Attr         []xml.Attr
}

func doMatch(el Element, match *Selector) bool {
	if el.Type == XML_Tag &&
		el.Name.Local == match.Label &&
		(match.Space == "" || match.Space == el.Name.Space) {
	matchAttr:
		for _, a := range match.Attr {
			for _, b := range el.StartElement.Attr {
				if a.Name.Local == b.Name.Local && a.Name.Space == b.Name.Space {
					if a.Value != b.Value {
						return false // Found but wrong value
					}
					continue matchAttr
				}
			}
			return false // Missing attr
		}
		return true // No more attrs to consider
	}
	return false // Name is not matchable
}

// Match returns a slice of matching child Element(s)
// matching a search.
func (el *Element) Match(match *Selector) []*Element {
	var matches []*Element
	if el != nil {
		for i, child := range el.Children {
			if doMatch(child, match) {
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
			if doMatch(child, match) {
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

// First returns a pointer to the first child of Element
func (el *Element) First() *Element {
	if el != nil && len(el.Children) > 0 {
		return &el.Children[0]
	}
	return nil
}

// Last returns a pointer to the first child of Element
func (el *Element) Last() *Element {
	if el != nil && len(el.Children) > 0 {
		return &el.Children[len(el.Children)-1]
	}
	return nil
}
