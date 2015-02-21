package dom

import (
	"encoding/xml"
	"strings"
	"testing"
)

type tc struct {
	name       string
	creator    func() *Document
	sample     string
	nameSpaces map[string]string
}

var testCases = []tc{
	tc{
		name: "EmptyDoc",
		creator: func() *Document {
			return CreateDocument()
		},
		sample: "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n",
	},
	tc{
		name: "OneEmptyNode",
		creator: func() *Document {
			doc := CreateDocument()
			root := ElementN("root")
			doc.SetRoot(root)
			return doc
		},
		sample: "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n<root/>\n",
	},
	tc{
		name: "MoreNodes",
		creator: func() *Document {
			doc := CreateDocument()
			root := ElementN("root")
			node1 := ElementN("node1")
			root.AddChild(node1)
			subnode := ElementN("sub")
			node1.AddChild(subnode)
			node2 := ElementN("node2")
			root.AddChild(node2)
			doc.SetRoot(root)
			return doc
		},
		sample: `<?xml version="1.0" encoding="utf-8"?>
<root>
 <node1>
  <sub/>
 </node1>
 <node2/>
</root>
`,
	},
	tc{
		name: "WithAttribs",
		creator: func() *Document {
			doc := CreateDocument()
			root := ElementN("root")
			node1 := ElementN("node1")
			node1.AddAttr(xml.Attr{Name: xml.Name{Local: "attr1"}, Value: "pouet"})
			root.AddChild(node1)
			doc.SetRoot(root)
			return doc
		},
		sample: `<?xml version="1.0" encoding="utf-8"?>
<root>
 <node1 attr1="pouet"/>
</root>
`,
	},
	tc{
		name: "WithContent",
		creator: func() *Document {
			doc := CreateDocument()
			root := ElementN("root")
			node1 := ElementN("node1")
			node1.Content = []byte("this is a text content")
			root.AddChild(node1)
			doc.SetRoot(root)
			return doc
		},
		sample: `<?xml version="1.0" encoding="utf-8"?>
<root>
 <node1>this is a text content</node1>
</root>
`,
	},
	tc{
		name: "WithNamespaces",
		creator: func() *Document {
			doc := CreateDocument()
			ns := "http://schemas.xmlsoap.org/ws/2004/08/addressing"
			root := ElementN("root")
			node1 := ElementN("node1")
			root.AddChild(node1)
			node1.Name.Space = ns
			node1.Content = []byte("this is a text content")
			doc.SetRoot(root)
			return doc
		},
		sample: `<?xml version="1.0" encoding="utf-8"?>
<root xmlns:ns0="http://schemas.xmlsoap.org/ws/2004/08/addressing">
 <ns0:node1>this is a text content</ns0:node1>
</root>
`,
	},
}

func TestCases(t *testing.T) {
	for _, testCase := range testCases {
		manualdoc := testCase.creator()
		parsedoc, err := Parse(strings.NewReader(testCase.sample))
		if err != nil {
			t.Fatalf("Cannot parse testcase %s sample %s\n\nGot error %v",
				testCase.name, testCase.sample, err)
		}
		if sample := manualdoc.String(); sample != testCase.sample {
			t.Fatalf("Manually created DOM for %s did not render.\nExpected: %s\n\nGot: %s\n",
				testCase.name, testCase.sample, sample)
		}
		if sample := parsedoc.String(); sample != testCase.sample {
			t.Fatalf("Parsed DOM for %s did not render.\nExpected: %s\n\nGot: %s\n",
				testCase.name, testCase.sample, sample)
		}
	}
}
