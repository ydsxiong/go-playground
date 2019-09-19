package clockface

import (
	"bytes"
	"encoding/xml"
	"math"
	"testing"
	"time"
)

func TestClockHand(t *testing.T) {
	testcases := []struct {
		name  string
		time  ClockHandDef
		point Point
	}{
		{
			"second hand test 1",
			GetSecondHandDef(simpletimeInSeconds(0, 0, 30)),
			Point{clockCentreX, clockCentreY + SecondHandLength},
		},
		{
			"minute hand test 1",
			GetMinuteHandDef(simpletimeInMinutes(0, 30, 0)),
			Point{clockCentreX, clockCentreY + MinuteHandLength},
		},
		{
			"hour hand test 1",
			GetHourHandDef(simpletimeInHours(6, 0, 0)),
			Point{clockCentreX, clockCentreY + HourHandLength},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := TransformHand(tc.time)

			if !roughlyEqualPoint(got, tc.point) {
				t.Errorf("Got %v, wanted %v", got, tc.point)
			}
		})
	}
}

func TestTimeInRadians(t *testing.T) {
	testcases := []struct {
		name  string
		time  float64
		cycle float64
		angle float64
	}{
		{
			"second hand angel test 1",
			simpletimeInSeconds(0, 0, 0), TimeInHalfClock,
			0,
		},
		{
			"second hand angel test 2",
			simpletimeInSeconds(0, 0, 30), TimeInHalfClock,
			math.Pi,
		},
		{
			"minute hand angel test 1",
			simpletimeInMinutes(0, 45, 0), TimeInHalfClock,
			(math.Pi / 2) * 3,
		},
		{
			"hour hand angel test 1",
			simpletimeInHours(7, 0, 0), HoursInHalfClock,
			(math.Pi / HoursInHalfClock) * 7,
		},
		{
			"hour hand angel test 2",
			simpletimeInHours(11, 0, 0), HoursInHalfClock,
			(math.Pi / HoursInHalfClock) * 11,
		},
		{
			"hour hand angel test 3",
			simpletimeInHours(11, 50, 0), HoursInHalfClock,
			(math.Pi / HoursInHalfClock) * (11.0 + 50.0/60.0),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := TimeInRadians(float64(tc.time), tc.cycle)

			if got != tc.angle {
				t.Errorf("Got %v, wanted %v", got, tc.angle)
			}
		})
	}
}

func TestClockHandPoint(t *testing.T) {
	testcases := []struct {
		name  string
		time  float64
		cycle float64
		point Point
	}{
		{
			"second hand point test 1",
			simpletimeInSeconds(0, 0, 30), TimeInHalfClock,
			Point{0, -1},
		},
		{
			"second hand point test 2",
			simpletimeInSeconds(0, 0, 45), TimeInHalfClock,
			Point{-1, 0},
		},
		{
			"minute hand point test 1",
			simpletimeInMinutes(0, 30, 0), TimeInHalfClock,
			Point{0, -1},
		},
		{
			"hour hand point test 1",
			simpletimeInHours(9, 0, 0), HoursInHalfClock,
			Point{-1, 0},
		},
		{
			"hour hand point test 2",
			simpletimeInHours(12, 0, 0), HoursInHalfClock,
			Point{0, 1},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := ClockHandPoint(float64(tc.time), tc.cycle)

			if !roughlyEqualPoint(got, tc.point) {
				t.Errorf("Got %v, wanted %v", got, tc.point)
			}
		})
	}
}

func simpletime(hours, minutes, seconds int) time.Time {
	return time.Date(312, time.October, 28, hours, minutes, seconds, 0, time.UTC)
}

func simpletimeInSeconds(hours, minutes, seconds int) float64 {
	time := time.Date(312, time.October, 28, hours, minutes, seconds, 0, time.UTC)
	return float64(time.Second())
}

func simpletimeInMinutes(hours, minutes, seconds int) float64 {
	time := time.Date(312, time.October, 28, hours, minutes, seconds, 0, time.UTC)
	return float64(time.Minute())
}

func simpletimeInHours(hours, minutes, seconds int) float64 {
	time := time.Date(312, time.October, 28, hours, minutes, seconds, 0, time.UTC)
	return float64(time.Hour()) + float64(time.Minute())/60.0
}

func testName(t time.Time) string {
	return t.Format("15:04:05")
}

func roughlyEqualFloat64(a, b float64) bool {
	const equalityThreshold = 1e-8
	return math.Abs(a-b) < equalityThreshold
}

func roughlyEqualPoint(a, b Point) bool {
	return roughlyEqualFloat64(a.X, b.X) &&
		roughlyEqualFloat64(a.Y, b.Y)
}

func TestSVGWriterHand(t *testing.T) {
	cases := []struct {
		name string
		time ClockHandDef
		line Line
	}{
		{
			"second hand line test 1",
			GetSecondHandDef(float64(simpletime(0, 0, 0).Second())),
			Line{150, 150, 150, 150 - SecondHandLength},
		},
		{
			"second hand line test 2",
			GetSecondHandDef(float64(simpletime(0, 0, 30).Second())),
			Line{150, 150, 150, 150 + SecondHandLength},
		},
		{
			"minute hand line test 1",
			GetMinuteHandDef(float64(simpletime(0, 30, 0).Minute())),
			Line{150, 150, 150, 150 + MinuteHandLength},
		},
		{
			"hour hand line test 1",
			GetHourHandDef(float64(simpletime(6, 0, 0).Hour())),
			Line{150, 150, 150, 150 + HourHandLength},
		},
		{
			"hour hand line test 1",
			GetHourHandDef((float64(simpletime(12, 0, 0).Hour()))),
			Line{150, 150, 150, 150 - HourHandLength},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b := bytes.Buffer{}
			WriteSVG(&b, tc.time)

			svg := SVG{}
			xml.Unmarshal(b.Bytes(), &svg)

			if !containsLine(tc.line, svg.Line) {
				t.Errorf("Expected to find the second hand line %+v, in the SVG lines %+v", tc.line, svg.Line)
			}
		})
	}
}

func containsLine(line Line, lines []Line) bool {
	for _, l := range lines {
		if line == l {
			return true
		}
	}
	return false
}

/**
That whole SecondHand function is super tied to being an SVG... without
mentioning SVGs or actually producing an SVG...
... while at the same time it's not testing any of my SVG code.

one option is to check the characters spewing out of the SVGWriter
contain things that look like the sort of SVG tag we're expecting for a particular time. For instance:

func TestSVGWriterAtMidnight(t *testing.T) {
    tm := time.Date(1337, time.January, 1, 0, 0, 0, 0, time.UTC)

    var b strings.Builder
    clockface.SVGWriter(&b, tm)
    got := b.String()

    want := `<line x1="150" y1="150" x2="150" y2="60"`

    if !strings.Contains(got, want) {
        t.Errorf("Expected to find the second hand %v, in the SVG output %v", want, got)
    }
}

a more sensible solution is to test the output as XML.
Parsing XML
encoding/xml is the Go package that can handle all things to do with simple XML parsing.
The function xml.Unmarshall takes a []byte of XML data and a pointer to a struct for it to get unmarshalled in to.
*/
