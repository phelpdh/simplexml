# Simplexml Dom library for Go

This is a naive and simple Go library to build a XML DOM to be able to produce
XML content, and parse simple XML into an in-memory DOM.

It started as a fork of https://github.com/masterzen/simplexml, but has
since been massively refactored to make it work more closely with encoding/xml,
and to include a set of useful functions for doing simple searches against the
element tree.

## Contact

- Bugs: https://github.com/VictorLowther/simplexml/issues


### Building

You can build the library from source:

```sh
git clone https://github.com/VictorLowther/simplexml
cd simplexml
go build
```

## Usage

