package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"
	vegeta "github.com/tsenart/vegeta/lib"
)

type Share struct {
	ShareId string
}

type Snapshot struct {
	ShareId     string `json:"share_id"`
	Force       bool   `json:"force"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Payload struct {
	Snapshot `json:"snapshot"`
}

func shareGenerator(shares []Share) <-chan *Share {
	c := make(chan *Share)
	go func() {
		defer close(c)
		for i := 0; i < len(shares); i++ {
			c <- &shares[i]
		}
	}()
	return c
}

// create snapshots
func NewSnapshotTargeter(shareCh <-chan *Share) vegeta.Targeter {
	header := http.Header{}
	header.Add("X-Auth-Token", authtoken)
	header.Add("Content-Type", "application/json")

	return func(target *vegeta.Target) error {
		target.Method = "POST"
		target.URL = baseURL + "/snapshots"
		target.Header = header
		log.Debug(target)

		s, ok := <-shareCh
		if ok {
			payload := &Payload{Snapshot{s.ShareId, false, "", ""}}
			buf := new(bytes.Buffer)
			json.NewEncoder(buf).Encode(payload)
			log.Debug(buf.String())
			target.Body = buf.Bytes()
			return nil
		} else {
			return errors.New("No more shares")
		}
	}
}
