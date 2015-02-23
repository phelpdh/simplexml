// Package search contains searching routines for the simplexml/dom package.
package search

import (
	"github.com/VictorLowther/simplexml/dom"
	"regexp"
)

// Match is the basic type of a search function.
// It takes a single element, and returns a boolean
// indicating whether the element matched the func.
type Match func(*dom.Element) bool

// And takes any number of Match, and returns another
// Match that will match if all of passed Match functions
// match.
func And(funcs ...Match) Match {
	return func(e *dom.Element) bool {
		for _, fn := range funcs {
			if !fn(e) {
				return false
			}
		}
		return true
	}
}

// Or takes any number of Match, and returns another Match
// that will match of any of the passed Match functions match.
func Or(funcs ...Match) Match {
	return func(e *dom.Element) bool {
		for _, fn := range funcs {
			if fn(e) {
				return true
			}
		}
		return false
	}
}

// Not takes a single Match, and returns another Match
// that matches if fn does not match.
func Not(fn Match) Match {
	return func(e *dom.Element) bool {
		return !fn(e)
	}
}

// All returns all the nodes that fn matches
func All(fn Match, nodes []*dom.Element) []*dom.Element {
	res := make([]*dom.Element, 0, 0)
	for _, n := range nodes {
		if fn(n) {
			res = append(res, n)
		}
	}
	return res
}

// First returns the first element that fn matches
func First(fn Match, nodes []*dom.Element) *dom.Element {
	for _, n := range nodes {
		if fn(n) {
			return n
		}
	}
	return nil
}

// Tag is a helper function for matching against a specific tag.
// It takes a name and a namespace URL to match against.
// If either name or space are "*", then they will match
// any value.
// Return is a Match.
func Tag(name, space string) Match {
	return func(e *dom.Element) bool {
		if space != "*" && space != e.Name.Space {
			return false
		}
		if name != "*" && name != e.Name.Local {
			return false
		}
		return true
	}
}

// TagRE is a helper function for matching against a specific tag
// using regular expressions.  It follows roughly the same rules as
// search.Tag
// Return is a Match
func TagRE(name, space *regexp.Regexp) Match {
	return func(e *dom.Element) bool {
		if space != nil && !space.MatchString(e.Name.Space) {
			return false
		}
		if name != nil && name.MatchString(e.Name.Local) {
			return false
		}
		return true
	}
}

// Attr creates a Match against the attributes of an element.
// It follows the same rules as Tag
func Attr(name, space, value string) Match {
	return func(e *dom.Element) bool {
		for _, a := range e.Attributes {
			if (space == "*" || space == a.Name.Space) &&
				(name == "*" || name == a.Name.Local) &&
				(value == "*" || value == a.Value) {
				return true
			}
		}
		return false
	}
}

// AttrRE creates a Match against the attributes of an element.
// It follows the same rules as MatchRE
func AttrRE(name, space, value *regexp.Regexp) Match {
	return func(e *dom.Element) bool {
		for _, a := range e.Attributes {
			if (space == nil || space.MatchString(a.Name.Space)) &&
				(name == nil || name.MatchString(a.Name.Local)) &&
				(value == nil || value.MatchString(a.Value)) {
				return true
			}
		}
		return false
	}
}

// ContentExists creates a Match against an element that has non-empty
// Content.
func ContentExists() Match {
	return func(e *dom.Element) bool {
		return len(e.Content) > 0
	}
}

// ContentRE creates a Match against the Content of am element
// that passes if the regex matches the content.
func ContentRE(regex *regexp.Regexp) Match {
	return func(e *dom.Element) bool {
		return regex.Match(e.Content)
	}
}
