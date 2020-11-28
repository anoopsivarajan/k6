package dns

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	miekg "github.com/miekg/dns"
)

//DNS type
type DNS struct{}

//Setup type
type Setup struct {
	Server  string `json:"server"`
	Timeout int    `json:"timeOut"`
}

type qData struct {
	Domain string `json:"domain"`
	QType  string `json:"qType"`
}

type optRecords struct {
	ID    int    `json:"id"`
	Value string `json:"value"`
}

//New creates DNS instance
func New() *DNS {
	return &DNS{}
}

func parseData(params map[string]interface{}) qData {
	var data qData
	values, _ := json.Marshal(params)
	_ = json.Unmarshal(values, &data)
	return data
}

//BuildMessage builds dns message
func (d *DNS) BuildMessage(ctx context.Context, data map[string]interface{}) []byte {
	values := parseData(data)

	if values.Domain == "" || values.QType == "" {
		return nil
	}

	msg := createMessage(values.Domain, values.QType)
	mm, _ := msg.Pack()
	return mm
}

//UnpackMessage unpacks dnsmessage
func (d *DNS) UnpackMessage(ctx context.Context, data []byte) *miekg.Msg {
	var msg = new(miekg.Msg)

	msg.Unpack(data)
	return msg
}

//Response info miekg
type Response struct {
	Data     *miekg.Msg `json:"data"`
	Duration int64      `json:"duration"`
}

//SendUDP function
func (*DNS) SendUDP(ctx context.Context, domain string, qType string, server string) (Response, error) {
	var m = createMessage(domain, qType)
	start := time.Now()
	client := new(miekg.Client)
	r, _, err := client.Exchange(m, server)
	if err != nil {
		return Response{}, err
	}
	end := time.Since(start)
	return Response{Data: r, Duration: end.Milliseconds()}, err
}

//CreateMessage message
func createMessage(domain string, qType string) *miekg.Msg {
	var mm = new(miekg.Msg)
	mm.SetQuestion(domain+".", getTypeFromString(qType))
	return mm
}

var types = map[string]uint16{
	"A":     miekg.TypeA,
	"AAAA":  miekg.TypeAAAA,
	"MX":    miekg.TypeMX,
	"Any":   miekg.TypeANY,
	"NS":    miekg.TypeNS,
	"CNAME": miekg.TypeCNAME,
	"PTR":   miekg.TypePTR,
	"TXT":   miekg.TypeTXT,
	"SRV":   miekg.TypeSRV,
	"OPT":   miekg.TypeOPT,
}

// GetTypeFromString returns the request type equivalent
func getTypeFromString(qType string) uint16 {
	requestType := miekg.TypeA
	typeFromMap := types[strings.ToUpper(qType)]
	if typeFromMap != 0 {
		requestType = typeFromMap
	}
	return requestType
}
