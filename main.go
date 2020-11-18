package main

import (
	"crypto/tls"
	"flag"
	"os"
	"os/signal"
	"time"

	"github.com/sapcc/load-test/loadtest"
	log "github.com/sirupsen/logrus"
	vegeta "github.com/tsenart/vegeta/lib"
)

func main() {
	var (
		ok       bool
		targeter vegeta.Targeter
	)

	// main flags
	authToken := flag.String("token", "", "auth token")
	region := flag.String("region", "", "ccloud region")
	loadTestRate := flag.Int("rate", 5, "load test rate (per minute)")
	loadTestDuration := flag.Duration("duration", 5*time.Second, "duration")
	loadTestTimeout := flag.Duration("timeout", 30*time.Second, "time out of each api call")
	debugFlag := flag.Bool("debug", false, "debug flag")

	// subcommand barbican's flags
	barbicanCmd := flag.NewFlagSet("barbican", flag.ExitOnError)

	// subcommand manila's flags
	// manilaCmd := flag.NewFlagSet("manila", flag.ExitOnError)
	// manilaShareFilePath := manilaCmd.String("shares", "", "path to file that contains share ids on each line")

	// parse main flags
	flag.Parse()

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

	// initialize load test context
	ctx := &loadtest.Context{
		AuthToken: *authToken,
		Region:    *region,
		Rate:      *loadTestRate,
		Duration:  *loadTestDuration,
	}

	// set log level
	if *debugFlag {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	// parse subcommands
	subcommand := flag.Args()
	if len(subcommand) < 1 {
		log.Fatal("subcommand is expected")
	}

	switch subcommand[0] {

	case "barbican":
		barbicanCmd.Parse(subcommand[1:])
		ctx.Name = "barbican-get-secret"
		ctx.Url = "https://keymanager-3." + ctx.Region + ".cloud.sap/v1/secrets"
		targeter = loadtest.NewStaticTargeter(ctx)
		log.Debugf("Barbican Load Test Context: %v", ctx)

	default:
		log.Error("subcommand can be 'barbican'")
		os.Exit(1)
	}

	attacker := vegeta.NewAttacker(
		vegeta.KeepAlive(true),                                  // keep alive
		vegeta.TLSConfig(&tls.Config{InsecureSkipVerify: true}), // insecure tls
		vegeta.Timeout(*loadTestTimeout),                        // timeout
	)
	rate := vegeta.Rate{Freq: ctx.Rate, Per: time.Minute}
	loadTestResult := attacker.Attack(targeter, rate, ctx.Duration, ctx.Name)
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
