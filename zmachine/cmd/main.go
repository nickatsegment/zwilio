package main

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"

	"github.com/nickatsegment/zwilio/zmachine"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)

	buffer, err := ioutil.ReadFile("zork1.dat")
	if err != nil {
		panic(err)
	}

	var header zmachine.ZHeader
	header.Read(buffer)

	if header.Version != 3 {
		panic("Only Version 3 files supported")
	}

	var zm zmachine.ZMachine

	/*
		zm.Stdin = os.Stdin
		zm.Stdout = os.Stdout
	*/
	stdout := bytes.Buffer{}
	zm.Stdout = &stdout

	stdinReader := bufio.NewReader(os.Stdin)
	zm.Initialize(buffer, header)

	for i := 0; !zm.Done(); i++ {
		log.Debugf("next batch %d", i)
		// first instruction shouldn't be ZRead, so we don't want to wait for input
		if i != 0 {
			input, err := stdinReader.ReadString('\n')
			if err != nil && err != io.EOF {
				log.Fatalf("failed to read input: %s", err)
			}
			log.Debugf("input: %q", input)
			zm.Stdin = bytes.NewReader([]byte(input))
		}

		zm.InterpretTillNextZRead()
		stdout.WriteTo(os.Stdout)
	}
}
