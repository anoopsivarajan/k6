package dns

import (
	"context"
	"testing"

	"github.com/dop251/goja"
	"github.com/loadimpact/k6/js/common"
	miekg "github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
)

func TestDNS_Resolve(t *testing.T) {

	rt := goja.New()
	ctx := context.Background()
	ctx = common.WithRuntime(ctx, rt)
	rt.Set("dns", common.Bind(rt, New(), &ctx))

	dns := New()

	t.Run("Parse data", func(t *testing.T) {
		resp, err := common.RunString(rt, `dns.sendUDP("google.com","A","8.8.8.8:53")`)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("UDP", func(t *testing.T) {
		resp, err := dns.SendUDP(ctx, "google.com", "A", "8.8.8.8:53")
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("Parse data", func(t *testing.T) {
		question := make(map[string]interface{})
		question["domain"] = "google.com"
		question["qType"] = "A"

		data := parseData(question)
		assert.Equal(t, "google.com", data.Domain)
		assert.Equal(t, "A", data.QType)
	})

	t.Run("Build message", func(t *testing.T) {
		question := make(map[string]interface{})
		question["domain"] = "google.com"
		question["qType"] = "A"

		resp := dns.BuildMessage(ctx, question)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp)
	})

	t.Run("Get type from string", func(t *testing.T) {
		qt := getTypeFromString("CNAME")
		assert.Equal(t, miekg.TypeCNAME, qt)

		qt = getTypeFromString("ahsss")
		assert.Equal(t, miekg.TypeA, qt)
	})

	t.Run("Create message", func(t *testing.T) {
		msg := createMessage("bbc.co.uk", "CNAME")
		assert.Equal(t, 1, len(msg.Question))
		q := msg.Question[0]
		assert.Equal(t, "bbc.co.uk.", q.Name)
		assert.Equal(t, uint16(5), q.Qtype)
	})

	t.Run("Unpack message", func(t *testing.T) {

		question := make(map[string]interface{})
		question["domain"] = "bbc.co.uk"
		question["qType"] = "A"

		resp := dns.BuildMessage(ctx, question)

		msg := dns.UnpackMessage(ctx, resp)

		assert.Equal(t, 1, len(msg.Question))
		q := msg.Question[0]
		assert.Equal(t, "bbc.co.uk.", q.Name)
		assert.Equal(t, uint16(1), q.Qtype)

	})
}
