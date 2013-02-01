package main

import (
	"encoding/xml"
	"time"
)

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

	// Required
	DcIdent   dcIdent `xml:"dc:identifier"`
	DcTitle   string  `xml:"dc:title"`
	DcCreator string  `xml:"dc:creator"`

	DcLang  string `xml:"dc:language,omitempty"`
	Metas   []meta
	Comment string `xml:",comment"`
}

type dcIdent struct {
	Id    string `xml:"id,attr"`
	Value string `xml:",innerxml"`
}

type meta struct {
	XMLname  xml.Name `xml:"meta"`
	Property string   `xml:"property,attr"`
	Value    string   `xml:",innerxml"`
}

type manifest struct {
	XMLName xml.Name `xml:"manifest"`
	Items   []item
}

type item struct {
	XMLName    xml.Name `xml:"item"`
	Href       string   `xml:"href,attr"`
	Id         string   `xml:"id,attr"`
	MediaType  string   `xml:"media-type,attr"`
	Properties string   `xml:"properties,attr,omitempty"`
}

type spine struct {
	XMLName  xml.Name `xml:"spine"`
	Itemrefs []itemref
}

type itemref struct {
	XMLName xml.Name `xml:"itemref"`
	IdRef   string   `xml:"idref,attr"`
}

func exampleOpfPackage() *OpfPackage {
	return &OpfPackage{
		XmlNs:   "http://www.idpf.org/2007/opf",
		Version: "3.0",
		UniqId:  "hi-epub",
		Metadata: metadata{
			XmlNs:     "http://purl.org/dc/elements/1.1",
			DcIdent:   dcIdent{"hi-epub", "hello-epub-0.1"},
			DcTitle:   "Hello 3Pub",
			DcCreator: "Thomas Jespersen",
			Metas: []meta{
				meta{Property: "dcterms:modified", Value: time.Now().String()},
			},
			Comment: "FIXME: The list of meta fields shouldn't be enclosed in the <Metas> tag",
		},
		Manifest: manifest{
			Items: []item{
				item{Href: "hello.xhtml", Id: "hello", MediaType: "application/xhtml+xml", Properties: "nav"},
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
