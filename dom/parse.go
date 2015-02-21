package dom

import (
	"bytes"
	"encoding/xml"
	"io"
)

func parseElement(decoder *xml.Decoder, tok xml.StartElement) (res *Element, err error) {
	res = CreateElement(tok.Name)
	for _, attr := range tok.Attr {
		res.AddAttr(attr)
	}

	for {
		newtok, err := decoder.Token()
		if err != nil {
			return nil, err
		}
		switch rt := newtok.(type) {
		case xml.EndElement:
			return res, nil
		case xml.CharData:
			content := bytes.TrimSpace([]byte(rt.Copy()))
			if len(content) > 0 {
				res.Content = content
			}
		case xml.StartElement:
			child, err := parseElement(decoder, rt)
			if err != nil {
				return nil, err
			}
			res.AddChild(child)
		}
	}
}

// Parse parses the XML document from the passed io.Reader and
// returns either a Document or an error if the io.Reader stream
// could not be parsed as well-formed XML.
func Parse(r io.Reader) (doc *Document, err error) {
	doc = CreateDocument()
	decoder := xml.NewDecoder(r)
	decoder.Strict = true
	var tok xml.Token
	for {
		tok, err = decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		switch rt := tok.(type) {
		case xml.StartElement:
			root, err := parseElement(decoder, rt)
			if err != nil {
				return nil, err
			}
			doc.SetRoot(root)
			return doc, nil
		}
	}
	return doc, nil
}
