package main

import (
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/ydsxiong/go-playground/clockface"
)

/**
https://quii.gitbook.io/learn-go-with-tests/go-fundamentals/math

An Acceptance Test

<line x1="0" y1="00" x2="10.0" y2="10.0"
		style="fill:none;stroke:#000;stroke-width:7px;"/>

The centre of the clock (the attributes x1 and y1 for this line) is the same for each hand of the clock.
The numbers that need to change for each hand of the clock - the parameters to whatever builds the SVG
- are the x2 and y2 attributes. there are also the radius of the clockface circle, the size of the SVG,
the colours of the hands, their shape, etc... it's better to start off by solving a simple, concrete problem
with a simple, concrete solution, and then to start adding parameters to make it generalised.

so:
every clock has a centre of (x, y)
the hour hand is h long
the minute hand is m long
the second hand is s long.

A thing to note about SVGs: the origin - point (0,0) - is at the top left hand corner, not anywhere else

*/

//var clockFilePath = "./clock.svg"

func main() {
	var clockfaceTemplate, err = template.ParseFiles("clockface.html")
	if err != nil {
		log.Printf("problem upgrading connection to WebSockets %v\n", err)
	} else {
		http.HandleFunc("/clock", func(w http.ResponseWriter, req *http.Request) {
			clockfaceTemplate.Execute(w, nil)
		})
		http.HandleFunc("/ws", sendClockToBrowser)
	}
	http.ListenAndServe(":9080", nil)
}

func sendClockToBrowser(w http.ResponseWriter, req *http.Request) {
	clockfaceWS := clockface.NewWwebSocket(w, req)

	ticker := time.NewTicker(time.Second)
	for currTime := range ticker.C {
		clockface.WriteSVG(clockfaceWS,
			clockface.GetSecondHandDef(float64(currTime.Second())),
			clockface.GetMinuteHandDef(float64(currTime.Minute())),
			clockface.GetHourHandDef((float64(currTime.Hour() + currTime.Minute()/60))))
	}
}

// func sendClockToFile() {
// 	clockfile, err := os.Create(clockFilePath)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer clockfile.Close()
// 	now := time.Now()
// 	clockface.WriteSVG(clockfile,
// 		clockface.GetSecondHandDef(float64(now.Second())),
// 		clockface.GetMinuteHandDef(float64(now.Minute())),
// 		clockface.GetHourHandDef(float64(now.Hour()+now.Minute()/60)))
// }
