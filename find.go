package xmltree

// Find returns a slice of matching child Element(s)
// in a depth-first matching a search.
func (el *Element) Find(match *Selector) []*Element {
	var results []*Element

	n := match.Depth
	if n < 1 {
		n = recursionLimit
	}

	search := func(el *Element) error {
		if el.Name.Local == match.Name.Local &&
			(match.Name.Space == "" || match.Name.Space == el.Name.Space) {
			results = append(results, el)
		}
		return nil
	}
	el.walkFuncDeep(search, n)

	return results
}

// FindFunc traverses the Element tree in depth-first order and returns
// a slice of Elements for which the function fn returns true.
func (el *Element) FindFunc(fn func(*Element) bool) []*Element {
	var results []*Element

	search := func(el *Element) error {
		if fn(el) {
			results = append(results, el)
		}
		return nil
	}
	el.WalkFunc(search)

	return results
}

// Flatten produces a slice of Element pointers referring to
// the children of el, and their children, in depth-first order.
func (el *Element) Flatten() []*Element {
	return el.FindFunc(func(*Element) bool { return true })
}
