package utils

import (
	"bytes"
	"encoding/xml"
	"io"
	"io/ioutil"
	"os/exec"
	"strings"
)

type Payload struct {
	Message string `xml:"message"`
}

func GetData(data io.Reader) string {
	var payload Payload
	decoder := xml.NewDecoder(data)
	decoder.Decode(&payload)

	return strings.ToUpper(payload.Message)
}

func GetXMLFromCommand() io.Reader {
	cmd := exec.Command("cat", "msg.xml")
	out, _ := cmd.StdoutPipe()

	cmd.Start()
	data, _ := ioutil.ReadAll(out)
	cmd.Wait()

	return bytes.NewReader(data)
}
