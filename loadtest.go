package main

import (
	"crypto/tls"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/sapcc/load-test/loadtest"
	log "github.com/sirupsen/logrus"
	vegeta "github.com/tsenart/vegeta/lib"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	var ok bool
	var targeter vegeta.Targeter
	var testName string
	var baseURL string
	var testType string
	var endpoint string

	app := kingpin.New("loadtest", "Load test openstack API in CCloud")
	authToken := app.Flag("token", "openstack auth token (or env OS_AUTH_TOKEN)").String()
	region := app.Flag("region", "ccloud region (or env OS_REGION)").Short('r').String()
	ratePerMinute := app.Flag("rate", "load test rate (per minute)").Default("5").Int()
	duration := app.Flag("duration", "load test duration").Default("5s").Duration()
	timeout := app.Flag("timeout", "load test timeout").Default("30s").Duration()
	debug := app.Flag("debug", "Enable debug mode").Short('d').Bool()

	barbicanCmd := app.Command("barbican", "Service Barbican")
	barbicanEndpoint := barbicanCmd.Arg("endpoint", "API endpoint").Required().String()

	glanceCmd := app.Command("glance", "Service Glance")
	glanceEndpoint := glanceCmd.Arg("endpoint", "API endpoint").Required().String()

	subcmd := kingpin.MustParse(app.Parse(os.Args[1:]))

	if len(*authToken) == 0 {
		*authToken, ok = os.LookupEnv("OS_AUTH_TOKEN")
		if !ok {
			log.Error("Auth token should be set either via -token or env variable OS_AUTH_TOKEN")
			os.Exit(1)
		}
	}
	if len(*region) == 0 {
		*region, ok = os.LookupEnv("OS_REGION")
		if !ok {
			log.Error("Region should be set either via -region flag or env variable OS_REGION")
			os.Exit(1)
		}
	}

	switch subcmd {
	case barbicanCmd.FullCommand():
		testType = "static"
		testName = "barbican"
		baseURL = fmt.Sprintf("http://keymanager-3.%s.cloud.sap", *region)
		endpoint = *barbicanEndpoint
	case glanceCmd.FullCommand():
		testType = "static"
		testName = "glance"
		baseURL = fmt.Sprintf("https://image-3.%s.cloud.sap/", *region)
		endpoint = *glanceEndpoint
	}

	if *debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	// targeter
	if testType == "static" {
		var err error
		targeter, err = loadtest.NewStaticTargeter(baseURL, endpoint, *authToken)
		if err != nil {
			log.Error(err)
		}
	} else {
		log.Error("Unknown targeter type %s", testType)
		os.Exit(1)
	}

	// attacker
	attacker := vegeta.NewAttacker(
		vegeta.KeepAlive(true),                                  // keep alive
		vegeta.TLSConfig(&tls.Config{InsecureSkipVerify: true}), // insecure tls
		vegeta.Timeout(*timeout),                                // timeout
	)
	rate := vegeta.Rate{Freq: *ratePerMinute, Per: time.Minute}
	loadTestResult := attacker.Attack(targeter, rate, *duration, testName)
	log.Debug("Attack rate %v", rate)

	// Interrupt signal
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	// Output encoder
	enc := vegeta.NewEncoder(os.Stdout)

	// Loop over load test result until finish or interrupted
	for {
		select {
		case <-sig:
			attacker.Stop()
			return
		case r, ok := <-loadTestResult:
			if !ok {
				return
			}
			if err := enc.Encode(r); err != nil {
				log.Error(err)
				os.Exit(1)
			}
		}
	}
}
