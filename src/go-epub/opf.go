package main

import "encoding/xml"

type OpfPackage struct {
	XMLName xml.Name `xml:"package"`

	XmlNs   string `xml:"xmlns,attr"`
	Version string `xml:"version,attr"`
	UniqId  string `xml:"unique-identifier,attr"`

	Metadata metadata
	Manifest manifest
	Spine    spine
}

type metadata struct {
	XMLName xml.Name `xml:"metadata"`
	XmlNs   string   `xml:"xmlns:dc,attr"`
}

type manifest struct {
	XMLName xml.Name `xml:"manifest"`
	Items   []item
}

type item struct {
	XMLName   xml.Name `xml:"item"`
	Href      string   `xml:"href,attr"`
	Id        string   `xml:"id,attr"`
	MediaType string   `xml:"media-type,attr"`
	Properties string `xml:"properties,attr,omitempty"`
}

type spine struct {
	XMLName  xml.Name `xml:"spine"`
	Itemrefs []itemref
}

type itemref struct {
	XMLName xml.Name `xml:"itemref"`
	IdRef   string `xml:"idref,attr"`
}

func exampleOpfPackage() *OpfPackage {
	return &OpfPackage{
		XmlNs:   "http://www.idpf.org/2007/opf",
		Version: "3.0",
		UniqId:  "hi-epub",
		Metadata: metadata{
			XmlNs: "http://purl.org/dc/elements/1.1",
		},
		Manifest: manifest{
			Items: []item{
				item{Href: "hello.xhtml", Id: "hello", MediaType: "application/xhtml+xml"},
				item{Href: "lesson0.xhtml", Id: "lesson0", MediaType: "application/xhtml+xml"},
			},
		},
		Spine: spine{
			Itemrefs: []itemref{
				itemref{IdRef: "hello"},
				itemref{IdRef: "lesson0"},
			},
		},
	}
}
