package lib

import (
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type LightroomToneCurvePoint struct {
	In  int
	Out int
}

// Example tone curve point:
// <rdf:li>228, 205</rdf:li>
func (self *LightroomToneCurvePoint) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// The next element will always be a text node
	text, err := d.Token()
	if err != nil {
		return err
	}

	parts := strings.Split(fmt.Sprintf("%s", text), ",")
	inStr := strings.Trim(parts[0], " ")
	outStr := strings.Trim(parts[1], " ")

	in, err := strconv.Atoi(inStr)
	if err != nil {
		return err
	}

	out, err := strconv.Atoi(outStr)
	if err != nil {
		return err
	}

	self.In = in
	self.Out = out

	return d.Skip()
}

type LightroomToneCurve struct {
	Points []LightroomToneCurvePoint
}

// Example ToneCurve:
// <crs:ToneCurvePV2012>
// 	<rdf:Seq>
// 		<rdf:li>0, 18</rdf:li>
// 		<rdf:li>39, 28</rdf:li>
// 		<rdf:li>83, 57</rdf:li>
// 		<rdf:li>119, 103</rdf:li>
// 		<rdf:li>228, 205</rdf:li>
// 		<rdf:li>255, 240</rdf:li>
// 	</rdf:Seq>
// </crs:ToneCurvePV2012>
func (self *LightroomToneCurve) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for {
		tok, err := d.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		switch tok.(type) {
		case xml.StartElement:
			element := tok.(xml.StartElement)
			if element.Name.Local == "li" {
				point := LightroomToneCurvePoint{}
				d.DecodeElement(&point, &element)
				self.Points = append(self.Points, point)
			}
		}
	}

	return d.Skip()
}

type LightroomCombinedToneCurve struct {
	Rgb   LightroomToneCurve
	Red   LightroomToneCurve
	Green LightroomToneCurve
	Blue  LightroomToneCurve
}

type LightroomPreset struct {
	Attributes map[string]string
	ToneCurve  LightroomCombinedToneCurve
}

func (self *LightroomPreset) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for {
		tok, err := d.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		switch tok.(type) {
		case xml.StartElement:
			element := tok.(xml.StartElement)
			switch element.Name.Local {
			case "Description":
				for _, attribute := range element.Attr {
					self.Attributes[attribute.Name.Local] = attribute.Value
				}
				break
			case "ToneCurvePV2012":
				d.DecodeElement(&self.ToneCurve.Rgb, &element)
				break
			case "ToneCurvePV2012Red":
				d.DecodeElement(&self.ToneCurve.Red, &element)
				break
			case "ToneCurvePV2012Green":
				d.DecodeElement(&self.ToneCurve.Green, &element)
				break
			case "ToneCurvePV2012Blue":
				d.DecodeElement(&self.ToneCurve.Blue, &element)
				break
			}
		}
	}

	return nil
}

func NewLightroomPreset() LightroomPreset {
	preset := LightroomPreset{}
	preset.Attributes = make(map[string]string)
	return preset
}
