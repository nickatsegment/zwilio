package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/nickatsegment/zwilio/zmachine"
	log "github.com/sirupsen/logrus"
)

var zMachines map[string]*zmachine.ZMachine = map[string]*zmachine.ZMachine{}

func main() {
	log.SetLevel(log.DebugLevel)

	datBuf, err := ioutil.ReadFile("zork1.dat")
	if err != nil {
		log.Fatalf("failed to load zork1.dat: %s", err)
	}

	var zHeader zmachine.ZHeader
	zHeader.Read(datBuf)

	if zHeader.Version != 3 {
		log.Fatalf("Only Version 3 files supported; is %d", zHeader.Version)
	}

	http.HandleFunc("/sms", func(w http.ResponseWriter, r *http.Request) {
		pn := r.FormValue("From")
		sessionID := pn
		if pn == "" {
			sessionID = "<anonymous>"
		}
		zm, zmExists := zMachines[sessionID]
		if !zmExists {
			log.Debugf("new zMachine for sessionID %s", sessionID)
			zm = &zmachine.ZMachine{}
			zm.Initialize(datBuf, zHeader)
			zMachines[sessionID] = zm
		}
		startIP := zm.IP()

		stdout := bytes.Buffer{}
		zm.Stdout = &stdout

		if zmExists {
			body := r.FormValue("Body")
			log.Debugf("Body: %q", body)
			stdin := bytes.NewBuffer([]byte(body))
			zm.Stdin = stdin
		} else {
			// TODO: else check "start" message
		}
		zm.InterpretTillNextZRead()
		endIP := zm.IP
		log.Debugf("start zm IP: %x; end zm IP %x", startIP, endIP)

		msgs := []string{}
		for _, l := range strings.Split(string(stdout.Bytes()), "\n") {
			lt := strings.Trim(l, "\r\n")
			if lt != "" && lt != ">" {
				msgs = append(msgs, lt)
			}
		}

		fmt.Fprintf(w, `<?xml version="1.0" encoding="UTF-8"?><Response><Message><Body>`)
		for _, msg := range msgs {
			msgEsc := bytes.Buffer{}
			err := xml.EscapeText(&msgEsc, []byte(msg))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "%s", err)
				return
			}
			fmt.Fprintf(w, "%s\n\n", string(msgEsc.Bytes()))
		}
		fmt.Fprintf(w, "</Body></Message></Response>")
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
