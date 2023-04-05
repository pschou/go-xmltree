package xmltree

// The Each method calls Func for each of the Element's children.  If the Func
// returns a non-nil error, Each will return it immediately.
func (el *Element) Each(fn func(*Element) error) (err error) {
	if el != nil {
		for i := 0; i < len(el.Children); i++ {
			if el.Children[i].Type == XML_Tag {
				if err = fn(&el.Children[i]); err != nil {
					return
				}
			}
		}
	}
	return
}

// The WalkDepthFunc method calls Func for each of the Element's children in a
// depth-first order.  If the Func returns true the children will
// continue to be considered, otherwise the depth is no longer searched.
func (el *Element) WalkDepthFunc(fn func(*Element) bool) {
	el.walkDepthDeep(fn, recursionLimit)
}
func (el *Element) walkDepthDeep(fn func(*Element) bool, n int) {
	if n--; n >= 0 {
		for i := 0; i < len(el.Children); i++ {
			if el.Children[i].Type == XML_Tag {
				if fn(&el.Children[i]) {
					el.Children[i].walkDepthDeep(fn, n)
				}
			}
		}
	}
}

// The WalkFunc method calls Func for each of the Element's children in a
// depth-first order.  If the Func returns a non-nil error, WalkFunc will
// return it immediately.
func (el *Element) WalkFunc(fn func(*Element) error) (err error) {
	return el.walkFuncDeep(fn, recursionLimit)
}
func (el *Element) walkFuncDeep(fn func(*Element) error, n int) (err error) {
	if n--; n >= 0 {
		for i := 0; i < len(el.Children); i++ {
			if el.Children[i].Type == XML_Tag {
				if err = fn(&el.Children[i]); err != nil {
					return
				}
				if err = el.Children[i].walkFuncDeep(fn, n); err != nil {
					return
				}
			}
		}
	}
	return
}

// walkFunc is the type of the function called for each of an Element's
// children.
//type walkFunc func(*Element)
