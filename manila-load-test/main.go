package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/sirupsen/logrus"

	vegeta "github.com/tsenart/vegeta/lib"
)

var (
	baseURL       string
	authtoken     string
	rate          int
	duration      time.Duration
	sharefilepath string
	tlsc          tls.Config
	debug         bool
)

func init() {
	flag.StringVar(&authtoken, "token", "", "auth token")
	flag.IntVar(&rate, "rate", 5, "rate")
	flag.DurationVar(&duration, "duration", time.Second, "duration")
	flag.StringVar(&sharefilepath, "shares", "", "path to file that contains share ids each line")
	flag.StringVar(&baseURL, "url", "", "url")
	flag.BoolVar(&debug, "debug", false, "debug")
	flag.Parse()
	if sharefilepath == "" || baseURL == "" {
		usage()
		os.Exit(1)
	}

	tlsc = tls.Config{InsecureSkipVerify: true}

	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func main() {
	shareCh, err := shareGeneratorFromFile(sharefilepath)
	if err != nil {
		fmt.Println(err)
		return
	}

	attacker := vegeta.NewAttacker(
		vegeta.KeepAlive(true),         // keep alive
		vegeta.TLSConfig(&tlsc),        // insecure tls
		vegeta.Timeout(60*time.Second), // timeout
	)

	targeter := NewSnapshotTargeter(shareCh)
	r := vegeta.Rate{Freq: rate, Per: time.Second} // default rate=5/second
	res := attacker.Attack(targeter, r, duration, "snapshot")

	// encode
	enc := vegeta.NewEncoder(os.Stdout)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	for {
		select {
		case <-sig:
			attacker.Stop()
			return
		case r, ok := <-res:
			if !ok {
				return
			}
			if err = enc.Encode(r); err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func usage() {
	flag.Usage()
}

func shareGeneratorFromFile(filepath string) (<-chan *Share, error) {
	var shares []Share
	sf, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer sf.Close()

	scanner := bufio.NewScanner(sf)
	for scanner.Scan() {
		shares = append(shares, Share{ShareId: scanner.Text()})
	}
	if err := scanner.Err(); err != nil {
		e := fmt.Errorf("reading file: %v", err)
		return nil, e
	}
	shareCh := shareGenerator(shares)
	return shareCh, nil
}
