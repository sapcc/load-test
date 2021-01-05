package loadtest

import "time"

type Context struct {
	Name string

	// parameters that are necessary for calling the Openstack APIs
	Url       string
	Region    string
	AuthToken string

	// load test parameters
	Rate     int
	Duration time.Duration
}
