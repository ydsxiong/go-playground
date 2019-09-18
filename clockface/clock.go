package clockface

import (
	"encoding/xml"
	"fmt"
	"io"
	"math"
)

const svgStart = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
<svg xmlns="http://www.w3.org/2000/svg"
     width="100%"
     height="100%"
     viewBox="0 0 800 800"
     version="2.0">`

const bezel = `<circle cx="%.5f" cy="%.5f" r="%.5f" style="fill:#fff;stroke:#000;stroke-width:3px;"/>`

const handLineTemplate = `<line x1="%.5f" y1="%.5f" x2="%.5f" y2="%.5f" style="fill:none;stroke:%s;stroke-width:2px;"/>`

const svgEnd = `</svg>`

type SVG struct {
	XMLName xml.Name `xml:"svg"`
	Xmlns   string   `xml:"xmlns,attr"`
	Width   string   `xml:"width,attr"`
	Height  string   `xml:"height,attr"`
	ViewBox string   `xml:"viewBox,attr"`
	Version string   `xml:"version,attr"`
	Circle  Circle   `xml:"circle"`
	Line    []Line   `xml:"line"`
}

type Circle struct {
	Cx float64 `xml:"cx,attr"`
	Cy float64 `xml:"cy,attr"`
	R  float64 `xml:"r,attr"`
}

type Line struct {
	X1 float64 `xml:"x1,attr"`
	Y1 float64 `xml:"y1,attr"`
	X2 float64 `xml:"x2,attr"`
	Y2 float64 `xml:"y2,attr"`
}

type Point struct {
	X float64
	Y float64
}

const (
	clockCentreX     = 52.0
	clockCentreY     = 52.0
	clockR           = 50.0
	SecondHandLength = 45
	MinuteHandLength = 40
	HourHandLength   = 30

	HoursInHalfClock = 6
	TimeInHalfClock  = 30

	SecondHandColor = "#f00"
	MinuteHandColor = "#000"
	HourHandColor   = "#00f"
)

type ClockHandDef struct {
	TimeInUnit    float64
	HandLength    float64
	Color         string
	HalfClockTime float64
}

/**
1. Scale it to the length of the hand
2. Flip it over the X axis to account for the fact that a SVG is having an origin in the top left hand corner
3. Translate it to the right position (so that it's coming from an origin of (150,150))
*/
func TransformHand(hand ClockHandDef) Point {
	point := ClockHandPoint(hand.TimeInUnit, hand.HalfClockTime)
	// scaling it up to SVG spec
	point.X = point.X * hand.HandLength
	point.Y = point.Y * hand.HandLength

	// flip it
	point.Y = -point.Y

	// translate it from 0, 0 to SVG spec
	point.X += clockCentreX
	point.Y += clockCentreX

	return point
}

func TimeInRadians(timeInUnit float64, halfClockTime float64) float64 {
	//return (math.Pi / 30) * float64(tm.Second())
	// a simple way to get the accuracy back could be rearranging the equation
	// so that there's no dividing down and then multiplying up, it'd be just dividing all the way down.
	return math.Pi / (halfClockTime / float64(timeInUnit))
}

func ClockHandPoint(timeInUnit float64, timeInHalfClock float64) Point {
	angle := TimeInRadians(timeInUnit, timeInHalfClock)
	x := math.Sin(angle)
	y := math.Cos(angle)
	return Point{x, y}
}

func WriteSVG(writer io.Writer, hands ...ClockHandDef) {
	//lines := []string{svgStart, bezel}
	clockSVG := svgStart
	clockSVG += fmt.Sprintf(bezel, clockCentreX, clockCentreY, clockR)
	for _, hand := range hands {
		shp := TransformHand(hand)
		line := fmt.Sprintf(handLineTemplate, clockCentreX, clockCentreY, shp.X, shp.Y, hand.Color)
		//lines = append(lines, line)
		clockSVG += line
	}

	//lines = append(lines, svgEnd)
	clockSVG += svgEnd
	io.WriteString(writer, clockSVG)
	//writeClockDefinition(writer, lines)
}

// func writeClockDefinition(writer io.Writer, contents []string) {
// 	for _, content := range contents {
// 		io.WriteString(writer, content)
// 	}
// }

func GetSecondHandDef(timeInUnit float64) ClockHandDef {
	return ClockHandDef{timeInUnit, SecondHandLength, SecondHandColor, TimeInHalfClock}
}

func GetMinuteHandDef(timeInUnit float64) ClockHandDef {
	return ClockHandDef{timeInUnit, MinuteHandLength, MinuteHandColor, TimeInHalfClock}
}

func GetHourHandDef(timeInUnit float64) ClockHandDef {
	return ClockHandDef{timeInUnit, HourHandLength, HourHandColor, HoursInHalfClock}
}
