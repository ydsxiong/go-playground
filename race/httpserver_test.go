package race

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMainHandler(t *testing.T) {

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("request generation error %s", err)
	}

	cases := []struct {
		name            string
		w               *httptest.ResponseRecorder
		r               *http.Request
		tokenHeader     string
		expectedrescode int
		expectedresbody []byte
		expectedlogs    []string
	}{
		{
			name:            "authorized",
			w:               httptest.NewRecorder(),
			r:               req,
			tokenHeader:     "magic",
			expectedrescode: http.StatusOK,
			expectedresbody: []byte("You have some magic in you\n"),
			expectedlogs: []string{
				"Allowed an access attempt\n",
			},
		},
		{
			name:            "unauthorized",
			w:               httptest.NewRecorder(),
			r:               req,
			tokenHeader:     "",
			expectedrescode: http.StatusForbidden,
			expectedresbody: []byte("You don't have enough magic in you\n"),
			expectedlogs:    []string{"Denied an access attempt\n"},
		},
	}

	chanLogs := make(chan error)
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.r.Header.Set("X-Access-Token", c.tokenHeader)

			//we create the Writer to pass in using io.Pipe. This will provide us with a Reader
			// that we can use to Read the subsequent Write calls made in the Logger.
			// We wrap the given PipeReader in a bufio.Reader so that we can easily read line-by-line
			// using a call to bufio.Reader’s ReadString method.
			logreader, logwriter := io.Pipe()
			buflogreader := bufio.NewReader(logreader)
			log.SetOutput(logwriter)

			//Note that in PipeWriter’s documentation it says:
			//Write implements the standard Write interface: it writes data to the pipe,
			// blocking until readers have consumed all the data or the read end is closed.
			go func() {
				for _, expectedline := range c.expectedlogs {
					msg, err := buflogreader.ReadString('\n')
					if err != nil {
						t.Errorf("Expected to be able to read from log but got error: %s", err)
					}
					if !strings.HasSuffix(msg, expectedline) {
						chanLogs <- fmt.Errorf("Log line didn't match suffix:\n\t%q\n\t%q", expectedline, msg)
					}
				}
			}()

			mainHandler(c.w, c.r)

			if c.w.Code != c.expectedrescode {
				t.Errorf("Status Code didn't match:\n\t%q\n\t%q", c.expectedrescode, c.w.Code)
			}
			if !bytes.Equal(c.expectedresbody, c.w.Body.Bytes()) {
				t.Errorf("Body didn't match:\n\t%q\n\t%q", string(c.expectedresbody), c.w.Body.String())
			}
			for range c.expectedlogs {
				select {
				case err := <-chanLogs:
					t.Error(err)
				default:
				}
			}
		})

	}
}
