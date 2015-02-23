package search

import (
	"encoding/xml"
	"github.com/VictorLowther/simplexml/dom"
	"strconv"
	"strings"
	"testing"
)

var testDoc string = `<?xml version="1.0" encoding="utf-8"?>
<a:root idx="0" xmlns:a="http://schemas.xmlsoap.org/ws/2004/08/addressing">
 <node1 foo="bar" idx="1">
  <sub idx="4"/>
 </node1>
 <node2 order="0" idx="2">I am Node 2
  <node2 order="2" idx="5">I am Groot</node2>
 </node2>
 <node2 order="1" idx="3">I am a different Node 2</node2>
</a:root>
`

func TestFindSimpleTag(t *testing.T) {
	doc, err := dom.Parse(strings.NewReader(testDoc))
	if err != nil {
		t.Fatalf("Cannot parse test document. Error: %v", err)
	}
	res := First(Tag("sub", ""), doc.Root().All())
	if res == nil {
		t.Fatal("Could not find sub element!")
	}
	if res.Name.Local != "sub" || res.Name.Space != "" {
		t.Fatalf("Looking for sub element gave me '%s' in namespace '%s'", res.Name.Local, res.Name.Space)
	}
}

func TestFindMultipleTags(t *testing.T) {
	doc, err := dom.Parse(strings.NewReader(testDoc))
	if err != nil {
		t.Fatalf("Cannot parse test document. Error: %v", err)
	}
	res := All(Tag("node2", ""), doc.Root().All())
	if len(res) != 3 {
		t.Fatalf("Expected to find 2 elements, found %d", len(res))
	}
	for i, e := range res {
		if e.Name.Local != "node2" || e.Name.Space != "" {
			t.Fatalf("Looking for node2 element gave me '%s' in namespace '%s'", e.Name.Local, e.Name.Space)
		}
		attr := e.Attributes[0]
		if attr.Name.Local != "order" {
			t.Fatal("Could not find expected order attribute on node2 element")
		}
		order, err := strconv.Atoi(attr.Value)
		if err != nil {
			t.Fatalf("Could not extract order attribute value: %v", err)
		}
		if order != i {
			t.Fatalf("Elements returned by All out of order! Expected %d, got %d", i, order)
		}
	}
}

func TestFindAttrs(t *testing.T) {
	doc, err := dom.Parse(strings.NewReader(testDoc))
	if err != nil {
		t.Fatalf("Cannot parse test document. Error: %v", err)
	}
	res := All(Attr("idx", "", "*"), doc.Root().All())
	if len(res) != 6 {
		t.Fatalf("Expected 6 elements, got %d", len(res))
	}

	for i, e := range res {
		var attr *xml.Attr
		for _, a := range e.Attributes {
			if a.Name.Local == "idx" {
				attr = &a
				break
			}
		}
		if attr == nil {
			t.Fatalf("Could not find idx addr on element %s", e.Name.Local)
		}
		idx, err := strconv.Atoi(attr.Value)
		if err != nil {
			t.Fatalf("Could not extract idx attribute value: %v", err)
		}
		if idx != i {
			t.Fatalf("Elements returned by attr search are out of order.  Expected %d, got %d", i, idx)
		}
	}
}

func TestAndCombinator(t *testing.T) {
	doc, err := dom.Parse(strings.NewReader(testDoc))
	if err != nil {
		t.Fatalf("Cannot parse test document. Error: %v", err)
	}
	expected := "I am a different Node 2"
	res := All(And(
		Attr("*","","1"),
		Tag("node2","")),
		doc.Root().All())
	if len(res) != 1 {
		t.Fatalf("Expected 1 element, got %d",len(res))
	}
	if string(res[0].Content) != expected {
		t.Fatalf("Expected node content not found!\nExpected: %s\n\nGot: %s",
			expected,
			string(res[0].Content))
	}
}

			
func TestOrCombinator(t *testing.T){
	doc, err := dom.Parse(strings.NewReader(testDoc))
	if err != nil {
		t.Fatalf("Cannot parse test document. Error: %v", err)
	}
	res := All(Or(
		Attr("idx","","0"),
		Attr("foo","","bar")),
		doc.Root().All())
	if len(res) != 2 {
		t.Fatalf("Expected 2 elements, got %d",len(res))
	}
	if (res[0].Name.Space != "http://schemas.xmlsoap.org/ws/2004/08/addressing" ||
		res[0].Name.Local != "root") {
		t.Fatalf("Expected first element to be root, not %s",res[0].Name.Local)
	}
	if res[1].Name.Local != "node1" {
		t.Fatalf("Expected second element to be node1, not %s",res[1].Name.Local)
	}
}
