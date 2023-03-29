// Package xmltree converts XML documents into a tree of Go values.
//
// The xmltree package provides types and routines for accessing
// and manipulating XML documents as trees, along with
// functionality to resolve XML namespace prefixes at any point
// in the tree.
package xmltree // import "github.com/pschou/go-xmltree"

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"sort"
	"strings"

	"golang.org/x/net/html/charset"
)

const (
	xmlNamespaceURI = "http://www.w3.org/2000/xmlns/"
	xmlLangURI      = "http://www.w3.org/XML/1998/namespace"
	recursionLimit  = 3000
)

type byXMLName []xml.Name

func (x byXMLName) Len() int { return len(x) }
func (x byXMLName) Less(i, j int) bool {
	return x[i].Space+x[i].Local < x[j].Space+x[j].Local
}
func (x byXMLName) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

var errDeepXML = errors.New("xmltree: xml document too deeply nested")

type Kind uint8

const (
	XML_Tag = Kind(iota)
	XML_Blob
	XML_Comment
)

// An Element represents a single element in an XML document. Elements
// may have zero or more children. The byte array used by the Content
// field is shared among all elements in the document, and should not
// be modified. An Element also captures xml namespace prefixes, so
// that arbitrary QNames in attribute values can be resolved.
type Element struct {
	// What is this element's kind
	Type Kind
	// Details about the Element if is a labeled tag
	xml.StartElement
	// The XML namespace scope at this element's location in the
	// document.
	Scope
	// The raw content contained within this element's start and
	// end tags. Uses the underlying byte array passed to Parse.
	Content []byte
	// Sub-elements contained within this element.
	Children []Element
}

// SetContent sets the value of the content escaping any HTML entities
func (el *Element) SetContent(text string) {
	el.Content = []byte(htmlEscaper.Replace(text))
}

var htmlEscaper = strings.NewReplacer(
	`<`, "&lt;",
	`>`, "&gt;",
)

// GetContent gets the value of the content while unescaping any HTML entities
func (el *Element) GetContent() string {
	if el != nil {
		return html.UnescapeString(string(el.Content))
	}
	return ""
}

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

// The JoinScope method joins two Scopes together. When resolving
// prefixes using the returned scope, the prefix list in the argument
// Scope is searched before that of the receiver Scope.
func (outer *Scope) JoinScope(inner *Scope) *Scope {
	return &Scope{append(outer.ns, inner.ns...)}
}

// Unmarshal parses the XML encoding of the Element and stores the result
// in the value pointed to by v. Unmarshal follows the same rules as
// xml.Unmarshal, but only parses the portion of the XML document
// contained by the Element.
func Unmarshal(el *Element, v interface{}) error {
	return xml.Unmarshal(Marshal(el), v)
}

// A Scope represents the xml namespace scope at a given position in
// the document.
type Scope struct {
	ns []xml.Name
}

// Resolve translates an XML QName (namespace-prefixed string) to an
// xml.Name with a canonicalized namespace in its Space field.  This can
// be used when working with XSD documents, which put QNames in attribute
// values. If qname does not have a prefix, the default namespace is used.If
// a namespace prefix cannot be resolved, the returned value's Space field
// will be the unresolved prefix. Use the ResolveNS function to detect when
// a namespace prefix cannot be resolved.
func (scope *Scope) Resolve(qname string) xml.Name {
	name, _ := scope.ResolveNS(qname)
	return name
}

// The ResolveNS method is like Resolve, but returns false for its second
// return value if a namespace prefix cannot be resolved.
func (scope *Scope) ResolveNS(qname string) (xml.Name, bool) {
	var prefix, local string
	parts := strings.SplitN(qname, ":", 2)
	if len(parts) == 2 {
		prefix, local = parts[0], parts[1]
	} else {
		prefix, local = "", parts[0]
	}
	switch prefix {
	case "xml":
		return xml.Name{Space: xmlLangURI, Local: local}, true
	case "xmlns":
		return xml.Name{Space: xmlNamespaceURI, Local: local}, true
	}
	for i := len(scope.ns) - 1; i >= 0; i-- {
		if scope.ns[i].Local == prefix {
			return xml.Name{Space: scope.ns[i].Space, Local: local}, true
		}
	}
	return xml.Name{Space: prefix, Local: local}, false
}

// ResolveDefault is like Resolve, but allows for the default namespace to
// be overridden. The namespace of strings without a namespace prefix
// (known as an NCName in XML terminology) will be defaultns.
func (scope *Scope) ResolveDefault(qname, defaultns string) xml.Name {
	if defaultns == "" || strings.Contains(qname, ":") {
		return scope.Resolve(qname)
	}
	return xml.Name{defaultns, qname}
}

// SimplifyNS will try to find a namespace which is already declared and is
// used majorly in the file and use that namespace as the default instead of
// using prefix for everies in the XML file.  Note: make sure a name is defined
// for every prefix used in the file, or deep lying "xmlns=" may be added.
func (el *Element) SimplifyNS() {
	var same, other int
	el.WalkDepthFunc(func(e *Element) bool {
		if len(e.Scope.ns) > 0 && e.Scope.ns[len(e.Scope.ns)-1].Local == "" {
			return false
		}
		if el.Name.Space == e.Name.Space {
			same++
		} else {
			other++
		}
		return true
	})
	if same > other {
		if len(el.Scope.ns) > 0 && el.Scope.ns[len(el.Scope.ns)-1].Local == "" {
			// Do nothing, as things are simple as it is
		} else {
			var foundDefault bool
			for i := range el.Scope.ns {
				if el.Scope.ns[i].Local == "" {
					// Overwrite the default
					el.Scope.ns[i].Space = el.Name.Space
					foundDefault = true
				}
			}
			if !foundDefault {
				el.Scope.ns = append([]xml.Name{xml.Name{Space: el.Name.Space}}, el.Scope.ns...)
			}
			el.WalkDepthFunc(func(e *Element) bool {
				// Locally defined new default
				if len(e.Scope.ns) > 0 && e.Scope.ns[len(e.Scope.ns)-1].Local == "" {
					return false
				}
				for i := range e.Scope.ns {
					if e.Scope.ns[i].Local == "" {
						// Overwrite the default
						e.Scope.ns[i].Space = el.Name.Space
						return true
					}
				}
				e.Scope.ns = append([]xml.Name{xml.Name{Space: el.Name.Space}}, e.Scope.ns...)
				//el.Scope.ns = append(el.Scope.ns, xml.Name{Space: el.Name.Space})
				//el.Scope.ns = append(el.Scope.ns, xml.Name{Space: el.Name.Space})
				return true
			})
		}
	} else {
		// Try to remove locally set ns
		el.RemoveLocalNS()
	}
}

// RemoteLocalNS will try to find a namespace which is already declared and use
// that namespace prefix instead of a locally defined one.
func (el *Element) RemoveLocalNS() error {
	if len(el.Scope.ns) > 0 && el.Scope.ns[len(el.Scope.ns)-1].Local == "" {
		var found bool
		for i := len(el.Scope.ns) - 2; i >= 0; i-- {
			if el.Scope.ns[i].Space == el.Name.Space {
				found = true
				break
			}
		}
		if found {
			el.Scope.ns = el.Scope.ns[:len(el.Scope.ns)-1]
			return nil
		}
		return errors.New("Could not find NS to prefix with")
	}
	return errors.New("No local defined, nothing done")
}

// Prefix is the inverse of Resolve. It uses the closest prefix
// defined for a namespace to create a string of the form
// prefix:local. If the namespace cannot be found, or is the
// default namespace, an unqualified name is returned.
func (scope *Scope) Prefix(name xml.Name) (qname string) {
	switch name.Space {
	case "":
		return name.Local
	case xmlLangURI:
		return "xml:" + name.Local
	case xmlNamespaceURI:
		return "xmlns:" + name.Local
	}
	for i := len(scope.ns) - 1; i >= 0; i-- {
		if scope.ns[i].Space == name.Space {
			if scope.ns[i].Local == "" {
				// Favor default NS if there is an extra
				// qualified NS declaration
				qname = name.Local
			} else if len(qname) == 0 {
				qname = scope.ns[i].Local + ":" + name.Local
			}
		}
	}
	return qname
}

func (scope *Scope) pushNS(tag xml.StartElement) []xml.Attr {
	var ns []xml.Name
	var newAttrs []xml.Attr
	for _, attr := range tag.Attr {
		if attr.Name.Space == "xmlns" {
			ns = append(ns, xml.Name{attr.Value, attr.Name.Local})
		} else if attr.Name.Local == "xmlns" {
			ns = append(ns, xml.Name{attr.Value, ""})
		} else {
			newAttrs = append(newAttrs, attr)
		}
	}
	// Within a single tag, all ns declarations are sorted. This reduces
	// differences between xmlns declarations between tags when
	// modifying the xml tree.
	sort.Sort(byXMLName(ns))
	if len(ns) > 0 {
		scope.ns = append(scope.ns, ns...)
		// Ensure that future additions to the scope create
		// a new backing array. This prevents the scope from
		// being clobbered during parsing.
		scope.ns = scope.ns[:len(scope.ns):len(scope.ns)]
	}
	return newAttrs
}

// Save some typing when scanning xml
type scanner struct {
	*xml.Decoder
	tok xml.Token
	err error
}

func (s *scanner) scan() bool {
	if s.err != nil {
		return false
	}
	s.tok, s.err = s.Token()
	return s.err == nil
}

// Parse builds a tree of Elements by reading an XML document.  The
// byte slice passed to Parse is expected to be a valid XML document
// with a single root element.
func Parse(doc []byte) (*Element, error) {
	d := xml.NewDecoder(bytes.NewReader(doc))

	// The xmltree package, when constructing the tree, takes slices
	// of the source document for chardata (data between tags). To do
	// this, it takes the position of the Decoder in the utf-8 input
	// stream. If the source document is not utf8, the position may be
	// incorrect and cause invalid data or a run-time panic. So we copy
	// the utf8 conversion to an internal buffer.
	utf8buf := bytes.NewBuffer(doc[:0])
	d.CharsetReader = func(label string, r io.Reader) (io.Reader, error) {
		utf8input, err := charset.NewReaderLabel(label, r)
		if err != nil {
			return nil, err
		}
		// At this point, the encoding/xml package has already
		// parsed the <?xml?> header. To be able to index
		// into the document, we need to account for this.
		padding := make([]byte, int(d.InputOffset()))
		utf8buf.Write(padding)

		_, err = io.Copy(utf8buf, utf8input)
		if err != nil {
			return nil, err
		}
		return bytes.NewReader(utf8buf.Bytes()[len(padding)+1:]), nil
	}
	scanner := scanner{Decoder: d}
	root := new(Element)

	for scanner.scan() {
		if start, ok := scanner.tok.(xml.StartElement); ok {
			root.StartElement = start
			break
		}
	}
	if scanner.err != nil {
		return nil, scanner.err
	}
	if err := root.parse(&scanner, utf8buf.Bytes(), 0); err != nil {
		return nil, err
	}
	return root, nil
}

func (el *Element) parse(scanner *scanner, data []byte, depth int) error {
	if depth > recursionLimit {
		return errDeepXML
	}
	el.StartElement.Attr = el.pushNS(el.StartElement)

	begin := scanner.InputOffset()
	end := begin
walk:
	for scanner.scan() {
		switch tok := scanner.tok.(type) {
		case xml.StartElement:
			child := Element{StartElement: tok.Copy(), Scope: el.Scope}
			if err := child.parse(scanner, data, depth+1); err != nil {
				return err
			}
			el.Children = append(el.Children, child)
		case xml.EndElement:
			if tok.Name != el.Name {
				return fmt.Errorf("Expecting </%s>, got </%s>", el.Prefix(el.Name), el.Prefix(tok.Name))
			}
			if len(el.Children) == 0 {
				el.Content = data[int(begin):int(end)]
			}
			break walk
		}
		end = scanner.InputOffset()
	}
	return scanner.err
}

// The Each method calls Func for each of the Element's children.  If the Func
// returns a non-nil error, Each will return it immediately.
func (el *Element) Each(fn func(*Element) error) (err error) {
	if el != nil {
		for i := 0; i < len(el.Children); i++ {
			if err = fn(&el.Children[i]); err != nil {
				return
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
			if fn(&el.Children[i]) {
				el.Children[i].walkDepthDeep(fn, n)
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
			if err = fn(&el.Children[i]); err != nil {
				return
			}
			if err = el.Children[i].walkFuncDeep(fn, n); err != nil {
				return
			}
		}
	}
	return
}

// Flatten produces a slice of Element pointers referring to
// the children of el, and their children, in depth-first order.
func (el *Element) Flatten() []*Element {
	return el.FindFunc(func(*Element) bool { return true })
}

// DelAttr removes an XML attribute from an Element's existing Attributes.
// If the attribute does not exist, no operation is done.
func (el *Element) DelAttr(space, local string) {
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

// walkFunc is the type of the function called for each of an Element's
// children.
//type walkFunc func(*Element)

// Match returns a slice of matching child Element(s)
// matching a search.
func (el *Element) Match(match *Selector) []*Element {
	var matches []*Element
	if el != nil {
		for i, child := range el.Children {
			if child.Name.Local == match.Label &&
				(match.Space == "" || match.Space == child.Name.Space) {
				matches = append(matches, &el.Children[i])
			}
		}
	}
	return matches
}

type Selector struct {
	Label, Space string
	Depth        int
}

// MatchOne returns a pointer to the first matching child of Element with a
// given match or nil if none matched.
func (el *Element) MatchOne(match *Selector) *Element {
	if el != nil {
		for i, child := range el.Children {
			if child.Name.Local == match.Label &&
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
				if e := el.Children[i].matchAnyDeep(match, d-1); e != nil {
					return e
				}
			}
		}
	}
	return nil
}

// Find returns a slice of matching child Element(s)
// in a depth-first matching a search.
func (el *Element) Find(match *Selector) []*Element {
	var results []*Element

	n := match.Depth
	if n < 1 {
		n = recursionLimit
	}

	search := func(el *Element) error {
		if el.Name.Local == match.Label &&
			(match.Space == "" || match.Space == el.Name.Space) {
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
