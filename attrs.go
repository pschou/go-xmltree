package xmltree

import (
	"encoding/xml"
	"strings"
)

// Attr gets the value of the first attribute whose name matches the
// space and local arguments. If space is the empty string, only
// attributes' local names are considered when looking for a match.
// If an attribute could not be found, the empty string is returned.
func (el *Element) Attr(space, local string) string {
	if el != nil {
		for _, v := range el.StartElement.Attr {
			if v.Name.Local != local {
				continue
			}
			if space == "" || space == v.Name.Space {
				return v.Value
			}
		}
	}
	return ""
}

// RemoveAttr removes an XML attribute from an Element's existing Attributes.
// If the attribute does not exist, no operation is done.
func (el *Element) RemoveAttr(space, local string) {
	for i, a := range el.StartElement.Attr {
		if a.Name.Local != local {
			continue
		}
		if space == "" || a.Name.Space == space {
			el.StartElement.Attr = append(
				el.StartElement.Attr[:i],
				el.StartElement.Attr[i+1:]...)
			return
		}
	}
}

// SetAttr adds an XML attribute to an Element's existing Attributes.
// If the attribute already exists, it is replaced.
func (el *Element) SetAttr(space, local, value string) {
	for i, a := range el.StartElement.Attr {
		if a.Name.Local != local {
			continue
		}
		if space == "" || a.Name.Space == space {
			el.StartElement.Attr[i].Value = value
			return
		}
	}
	el.StartElement.Attr = append(el.StartElement.Attr, xml.Attr{
		Name:  xml.Name{space, local},
		Value: value,
	})
}

// RemoveClass adds a class attribute to an Element's existing class.
// If the class already exists, nothing is done
func (el *Element) RemoveClass(class string) {
	toRemove := strings.Split(class, " ")
	for i, a := range el.StartElement.Attr {
		if a.Name.Space == "" && a.Name.Local == "class" {
			classes := cleanupSlice(strings.Split(a.Value, " "), toRemove)
			el.StartElement.Attr[i].Value = strings.Join(classes, " ")
			return
		}
	}
}

// AddClass adds a class attribute to an Element's existing class.
// If the class already exists, nothing is done
func (el *Element) AddClass(class string) {
	toAdd := strings.Split(class, " ")
	for i, a := range el.StartElement.Attr {
		if a.Name.Space == "" && a.Name.Local == "class" {
			el.StartElement.Attr[i].Value = strings.Join(cleanupSlice(append(strings.Split(a.Value, " "), toAdd...), nil), " ")
			return
		}
	}
	el.StartElement.Attr = append(el.StartElement.Attr, xml.Attr{
		Name:  xml.Name{"", "class"},
		Value: class,
	})
}

func cleanupSlice(vals, omit []string) (ret []string) {
cleanup:
	for _, v := range vals {
		if v == "" {
			continue
		}
		for _, o := range omit {
			if o == v {
				continue cleanup
			}
		}
		for _, r := range ret {
			if r == v {
				continue cleanup
			}
		}
		ret = append(ret, v)
	}
	return
}
