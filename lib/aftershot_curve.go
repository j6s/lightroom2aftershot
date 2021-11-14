package lib

import (
	"encoding/xml"
	"fmt"
	"log"
	"strconv"
	"strings"
)

/*
 * Aftershot curves XML format explained:
 * There are 5 attributes that are relevant for building curves in aftershot. They contain
 * comma separated versions of the curve information for all available curves (RGB, R, G, B).
 *
 * When luminance values are used they are 16bit integers (0-65535)
 *
 * bopt:curves_m_cn: Number of points in the curve.
 *		- Is prefixed with `4,1` (Not sure what these values mean)
 *      - Has the number of points in the curve as integers following that (minimum 2 for black & white)
 *      - Example: `4,1,4,2,2,2`
 *          -> Has a curve with 2 additional points (black + white + 2 others) for RGB
 *
 * bopt:curves_m_cx: Point locations on X axis (input)
 * bopt:curves_m_cy: Point locations on Y axis (output)
 *      - Prefixed with `4,20` (Note sure what 4 means, 20 seems to be the number of points per channel [see more below])
 *      - Each channel has a list of 20 values (presumably the `20` in the prefix) describing up to 20 points in the curve.
 *        only the first `curves_m_cn` are read from that list for each channel. So if `curves_m_cn` is set to 4 for the
 *        RGB channel (as is the case in the example above), then the first 4 of 20 values for RGB will be read, the rest will
 *        be ignored.
 *      - Examples with explanation:
 *            ######################|                     RGB                       |                   R                       |                    G                      |                        B                  |
 *            bopt:curves_m_cx="4,20,0,30351,65535,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,65535,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,65535,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,65535,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0"
 *            bopt:curves_m_cy="4,20,0,40799,65535,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,65535,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,65535,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,65535,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0"
 *
 * bopt:curves_m_olo: Output Black point for channels (black triangle left)
 * bopt:curves_m_ohi: Output White point for channels (white triangle left)
 * bopt:curves_m_ilo: Input Black point for channels (black triangle bottom)
 * bopt:curves_m_imid: Input Mid poitn for channels (gray triangle bottom) [gamma]
 * bopt:curves_m_ihi: Input White point for channels (white triangle bottom)
 * 		- Prefixed by `4,1` (Not sure what these values mean)
 *      - Followed by one value per channel
 *      - Note the following difference: You can drag the last point of the curve down (or the first up) to effectively change
 *        the black & white points. That however will only change the X and Y positions of the points. The black and white points
 *        in aftershot correspond to the black and white triangles on the left side.
 *
 *
 */

type AfterShotToneCurvePoint struct {
	In  int
	Out int
}

type AfterShotToneCurveChannel struct {
	Points []AfterShotToneCurvePoint
}

func (self *AfterShotToneCurveChannel) serializePoints(points []int, maxNumberOfPoints int) string {
	serialized := make([]string, maxNumberOfPoints)

	// Fill with empty values
	for index, _ := range serialized {
		serialized[index] = "0"
	}

	// Every valid curve must have at least 2 points.
	if len(points) < 2 {
		serialized[0] = "0"
		serialized[1] = strconv.Itoa(AFTERSHOT_CURVE_MAX)
	} else {
		for index, value := range points {
			serialized[index] = strconv.Itoa(value)
		}
	}

	return strings.Join(serialized, ",")
}

func (self *AfterShotToneCurveChannel) SerializeNumberOfPoints() string {
	if len(self.Points) < 2 {
		return "2"
	}

	return strconv.Itoa(len(self.Points))
}

func (self *AfterShotToneCurveChannel) SerializePointsIn(maxNumberOfPoints int) string {
	points := make([]int, len(self.Points))

	for index, point := range self.Points {
		points[index] = point.In
	}

	return self.serializePoints(points, maxNumberOfPoints)
}

func (self *AfterShotToneCurveChannel) SerializePointsOut(maxNumberOfPoints int) string {
	points := make([]int, len(self.Points))

	for index, point := range self.Points {
		points[index] = point.Out
	}

	return self.serializePoints(points, maxNumberOfPoints)
}

// ---

type AfterShotCombinedToneCurve struct {
	Rgb   AfterShotToneCurveChannel
	Red   AfterShotToneCurveChannel
	Green AfterShotToneCurveChannel
	Blue  AfterShotToneCurveChannel
}

func (self *AfterShotCombinedToneCurve) SerializeNumberOfPoints() string {
	return fmt.Sprintf(
		"4,1,%s,%s,%s,%s",
		self.Rgb.SerializeNumberOfPoints(),
		self.Red.SerializeNumberOfPoints(),
		self.Green.SerializeNumberOfPoints(),
		self.Blue.SerializeNumberOfPoints(),
	)
}

func (self *AfterShotCombinedToneCurve) SerializePointsIn(maxNumberOfPoints int) string {
	return fmt.Sprintf(
		"4,%d,%s,%s,%s,%s",
		maxNumberOfPoints,
		self.Rgb.SerializePointsIn(maxNumberOfPoints),
		self.Red.SerializePointsIn(maxNumberOfPoints),
		self.Green.SerializePointsIn(maxNumberOfPoints),
		self.Blue.SerializePointsIn(maxNumberOfPoints),
	)
}

func (self *AfterShotCombinedToneCurve) SerializePointsOut(maxNumberOfPoints int) string {
	return fmt.Sprintf(
		"4,%d,%s,%s,%s,%s",
		maxNumberOfPoints,
		self.Rgb.SerializePointsOut(maxNumberOfPoints),
		self.Red.SerializePointsOut(maxNumberOfPoints),
		self.Green.SerializePointsOut(maxNumberOfPoints),
		self.Blue.SerializePointsOut(maxNumberOfPoints),
	)
}

func (self *AfterShotCombinedToneCurve) ToXmlAttributes() []xml.Attr {
	maxList := fmt.Sprintf("4,1,%[1]d,%[1]d,%[1]d,%[1]d", AFTERSHOT_CURVE_MAX)
	log.Print(maxList)

	return []xml.Attr{
		{Name: xml.Name{Local: "bopt:curves_m_cn"}, Value: self.SerializeNumberOfPoints()},
		{Name: xml.Name{Local: "bopt:curves_m_cx"}, Value: self.SerializePointsIn(AFTERSHOT_NUM_POINTS)},
		{Name: xml.Name{Local: "bopt:curves_m_cy"}, Value: self.SerializePointsOut(AFTERSHOT_NUM_POINTS)},
		{Name: xml.Name{Local: "bopt:curves_m_olo"}, Value: "4,1,0,0,0,0"},
		{Name: xml.Name{Local: "bopt:curves_m_ohi"}, Value: maxList},
		{Name: xml.Name{Local: "bopt:curves_m_ilo"}, Value: "4,1,0,0,0,0"},
		{Name: xml.Name{Local: "bopt:curves_m_imid"}, Value: "4,1,1,1,1,1"},
		{Name: xml.Name{Local: "bopt:curves_m_ihi"}, Value: maxList},
	}
}
