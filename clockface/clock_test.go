package clockface_test

import (
	"bytes"
	"encoding/xml"
	"math"
	"testing"
	"time"

	"github.com/ydsxiong/go-playground/clockface"
)

func TestClockHand(t *testing.T) {
	testcases := []struct {
		name  string
		time  clockface.ClockHandDef
		point clockface.Point
	}{
		{
			"second hand test 1",
			clockface.GetSecondHandDef(simpletimeInSeconds(0, 0, 30)),
			clockface.Point{clockface.ClockCentreX, clockface.ClockCentreY + clockface.SecondHandLength},
		},
		{
			"minute hand test 1",
			clockface.GetMinuteHandDef(simpletimeInMinutes(0, 30, 0)),
			clockface.Point{clockface.ClockCentreX, clockface.ClockCentreY + clockface.MinuteHandLength},
		},
		{
			"hour hand test 1",
			clockface.GetHourHandDef(simpletimeInHours(6, 0, 0)),
			clockface.Point{clockface.ClockCentreX, clockface.ClockCentreY + clockface.HourHandLength},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := clockface.TransformHand(tc.time)

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
			simpletimeInSeconds(0, 0, 0), clockface.TimeInHalfClock,
			0,
		},
		{
			"second hand angel test 2",
			simpletimeInSeconds(0, 0, 30), clockface.TimeInHalfClock,
			math.Pi,
		},
		{
			"minute hand angel test 1",
			simpletimeInMinutes(0, 45, 0), clockface.TimeInHalfClock,
			(math.Pi / 2) * 3,
		},
		{
			"hour hand angel test 1",
			simpletimeInHours(7, 0, 0), clockface.HoursInHalfClock,
			(math.Pi / clockface.HoursInHalfClock) * 7,
		},
		{
			"hour hand angel test 2",
			simpletimeInHours(11, 0, 0), clockface.HoursInHalfClock,
			(math.Pi / clockface.HoursInHalfClock) * 11,
		},
		{
			"hour hand angel test 3",
			simpletimeInHours(11, 50, 0), clockface.HoursInHalfClock,
			(math.Pi / clockface.HoursInHalfClock) * (11.0 + 50.0/60.0),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := clockface.TimeInRadians(float64(tc.time), tc.cycle)

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
		point clockface.Point
	}{
		{
			"second hand point test 1",
			simpletimeInSeconds(0, 0, 30), clockface.TimeInHalfClock,
			clockface.Point{0, -1},
		},
		{
			"second hand point test 2",
			simpletimeInSeconds(0, 0, 45), clockface.TimeInHalfClock,
			clockface.Point{-1, 0},
		},
		{
			"minute hand point test 1",
			simpletimeInMinutes(0, 30, 0), clockface.TimeInHalfClock,
			clockface.Point{0, -1},
		},
		{
			"hour hand point test 1",
			simpletimeInHours(9, 0, 0), clockface.HoursInHalfClock,
			clockface.Point{-1, 0},
		},
		{
			"hour hand point test 2",
			simpletimeInHours(12, 0, 0), clockface.HoursInHalfClock,
			clockface.Point{0, 1},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := clockface.ClockHandPoint(tc.time, tc.cycle)

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
	time := simpletime(hours, minutes, seconds)
	return float64(time.Second())
}

func simpletimeInMinutes(hours, minutes, seconds int) float64 {
	time := simpletime(hours, minutes, seconds)
	return float64(time.Minute())
}

func simpletimeInHours(hours, minutes, seconds int) float64 {
	time := simpletime(hours, minutes, seconds)
	return float64(time.Hour()) + float64(time.Minute())/60.0
}

func testName(t time.Time) string {
	return t.Format("15:04:05")
}

func roughlyEqualFloat64(a, b float64) bool {
	const equalityThreshold = 1e-8
	return math.Abs(a-b) < equalityThreshold
}

func roughlyEqualPoint(a, b clockface.Point) bool {
	return roughlyEqualFloat64(a.X, b.X) &&
		roughlyEqualFloat64(a.Y, b.Y)
}

func TestSVGWriterHand(t *testing.T) {
	cases := []struct {
		name string
		time clockface.ClockHandDef
		line clockface.Line
	}{
		{
			"second hand line test 1",
			clockface.GetSecondHandDef(simpletimeInSeconds(0, 0, 0)),
			clockface.Line{clockface.ClockCentreX, clockface.ClockCentreY, clockface.ClockCentreX, clockface.ClockCentreY - clockface.SecondHandLength},
		},
		{
			"second hand line test 2",
			clockface.GetSecondHandDef(simpletimeInSeconds(0, 0, 30)),
			clockface.Line{clockface.ClockCentreX, clockface.ClockCentreY, clockface.ClockCentreX, clockface.ClockCentreY + clockface.SecondHandLength},
		},
		{
			"minute hand line test 1",
			clockface.GetMinuteHandDef(simpletimeInMinutes(0, 30, 0)),
			clockface.Line{clockface.ClockCentreX, clockface.ClockCentreY, clockface.ClockCentreX, clockface.ClockCentreY + clockface.MinuteHandLength},
		},
		{
			"hour hand line test 1",
			clockface.GetHourHandDef(simpletimeInHours(6, 0, 0)),
			clockface.Line{clockface.ClockCentreX, clockface.ClockCentreY, clockface.ClockCentreX, clockface.ClockCentreY + clockface.HourHandLength},
		},
		{
			"hour hand line test 1",
			clockface.GetHourHandDef(simpletimeInHours(12, 0, 0)),
			clockface.Line{clockface.ClockCentreX, clockface.ClockCentreY, clockface.ClockCentreX, clockface.ClockCentreY - clockface.HourHandLength},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b := bytes.Buffer{}
			clockface.WriteSVG(&b, tc.time)

			svg := clockface.SVG{}
			xml.Unmarshal(b.Bytes(), &svg)

			if !containsLine(tc.line, svg.Line) {
				t.Errorf("Expected to find the second hand line %+v, in the SVG lines %+v", tc.line, svg.Line)
			}
		})
	}
}

func containsLine(line clockface.Line, lines []clockface.Line) bool {
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
