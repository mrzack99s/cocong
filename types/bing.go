package types

import "encoding/xml"

type BingImageAPIResponse struct {
	XMLName xml.Name `xml:"images"`
	Image   struct {
		URLBase   string `xml:"urlBase"`
		Copyright string `xml:"copyright"`
	} `xml:"image"`
}
