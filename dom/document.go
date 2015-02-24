// Package dom implements a simple XML DOM that
// is a light wrapper on top of encoding/xml
package dom

import (
	"bytes"
)

// A Document represents an entire XML document.
type Document struct {
	root *Element
}

// CreateDocument creates a new XML document.
func CreateDocument() *Document {
	return &Document{}
}

// Root returns the root element of the document.
func (doc *Document) Root() (node *Element) {
	return doc.root
}

// SetRoot sets the new root element of the document.
func (doc *Document) SetRoot(node *Element) {
	node.parent = nil
	doc.root = node
}

// Encode encodes the entire Document using the passed-in Encoder.
// The output is a well-formed XML document.
func (doc *Document) Encode(e *Encoder) (err error) {
	_,err = e.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>")
	if err != nil {
		return err
	}
	if err = e.prettyEnd(); err != nil {
		return err
	}
	if doc.root != nil {
		return doc.root.Encode(e)
	}
	return nil
}

// Bytes returns the results of running Encoder against a byte array, which
// contains a well-formed XML document.
func (doc *Document) Bytes() []byte {
	b := bytes.Buffer{}
	encoder := NewEncoder(&b)
	encoder.Pretty()
	doc.Encode(encoder)
	encoder.Flush()
	return b.Bytes()
}

// String returns the result of stringifying the byte array that Bytes returns.
func (doc *Document) String() string {
	return string(doc.Bytes())
}
