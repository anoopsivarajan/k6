package dns

import (
	"context"
	"encoding/json"
	"net"
	"time"
)

//DNS type
type DNS struct{}

//Setup type
type Setup struct {
	Server  string `json:"server"`
	Timeout int    `json:"timeOut"`
}

//New creates DNS instance
func New() *DNS {
	return &DNS{}
}

func parseSetup(params map[string]interface{}) Setup {
	var setup Setup
	values, _ := json.Marshal(params)
	_ = json.Unmarshal(values, &setup)
	return setup
}

func updateResolver(setup Setup) {
	if setup.Timeout == 0 {
		setup.Timeout = 10000
	}
	resolver := &net.Resolver{
		PreferGo: false,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(setup.Timeout),
			}
			return d.DialContext(ctx, "udp", setup.Server)
		},
	}
	net.DefaultResolver = resolver
}

//Resolve resolves the dns request
func (d *DNS) Resolve(ctx context.Context, address string, params map[string]interface{}) ([]string, error) {
	//Parse params from the script
	setup := parseSetup(params)
	//check if there is a customer resolver ip
	if setup.Server != "" {
		updateResolver(setup)
	}
	//Return the IP address
	return net.LookupHost(address)
}
