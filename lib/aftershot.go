package lib

import (
	"encoding/xml"
)

type AfterShotPreset struct {
	ToneCurve  AfterShotCombinedToneCurve
	Attributes map[string]string
}

func (self AfterShotPreset) MarshalXML(e *xml.Encoder, start xml.StartElement) error {

	attributes := []xml.Attr{}
	attributes = append(attributes, self.ToneCurve.ToXmlAttributes()...)
	for key, value := range self.Attributes {
		attributes = append(attributes, xml.Attr{
			Name:  xml.Name{Local: key},
			Value: value,
		})
	}

	e.EncodeToken(xml.StartElement{
		Name: xml.Name{Local: "x:xmpmeta"},
		Attr: []xml.Attr{
			{Name: xml.Name{Local: "xmlns:x"}, Value: "adobe:ns:meta/"},
			{Name: xml.Name{Local: "x:xmptk"}, Value: "XMP Core 4.4.0"},
		},
	})
	e.EncodeToken(xml.StartElement{
		Name: xml.Name{Local: "rdf:RDF"},
		Attr: []xml.Attr{
			{Name: xml.Name{Local: "xmlns:rdf"}, Value: "http://www.w3.org/1999/02/22-rdf-syntax-ns#"},
		},
	})
	e.EncodeToken(xml.StartElement{
		Name: xml.Name{Local: "rdf:Description"},
		Attr: []xml.Attr{
			{Name: xml.Name{Local: "rdf:about"}, Value: ""},
			{Name: xml.Name{Local: "xmlns:bib"}, Value: "http://www.bibblelabs.com/BibbleToplevel/5.0/"},
			{Name: xml.Name{Local: "xmlns:bset"}, Value: "http://www.bibblelabs.com/BibbleSettings/5.0/"},
			{Name: xml.Name{Local: "xmlns:blay"}, Value: "http://www.bibblelabs.com/BibbleLayers/5.0/"},
			{Name: xml.Name{Local: "xmlns:bopt"}, Value: "http://www.bibblelabs.com/BibbleOpt/5.0/"},
		},
	})
	e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "bib:settings"}})

	e.EncodeToken(xml.StartElement{
		Name: xml.Name{Local: "rdf:Description"},
		Attr: []xml.Attr{
			{Name: xml.Name{Local: "bset:settingsVersion"}, Value: "66"},
			{Name: xml.Name{Local: "bset:respectsTransfor"}, Value: "True"},
			{Name: xml.Name{Local: "bset:curLayer"}, Value: "0"},
		},
	})
	e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "bset:layers"}})
	e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "rdf:Seq"}})
	e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "rdf:li"}})
	e.EncodeToken(xml.StartElement{
		Name: xml.Name{Local: "rdf:Description"},
		Attr: []xml.Attr{
			{Name: xml.Name{Local: "blay:layerId"}, Value: "0"},
			{Name: xml.Name{Local: "blay:layerPos"}, Value: "0"},
			{Name: xml.Name{Local: "blay:name"}, Value: ""},
			{Name: xml.Name{Local: "blay:enabled"}, Value: "True"},
		},
	})

	e.EncodeElement("", xml.StartElement{
		Name: xml.Name{Local: "blay:options"},
		Attr: attributes,
	})

	e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "rdf:Description"}})
	e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "rdf:li"}})
	e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "rdf:Seq"}})
	e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "bset:layers"}})
	e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "rdf:Description"}})
	e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "bib:settings"}})
	e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "rdf:Description"}})
	e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "rdf:RDF"}})
	e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "x:xmpmeta"}})

	return nil
}
