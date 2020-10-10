package dns

import (
	"context"
	"testing"

	"github.com/dop251/goja"
	"github.com/loadimpact/k6/js/common"
	"github.com/stretchr/testify/assert"
)

func TestDNS_Resolve(t *testing.T) {

	rt := goja.New()
	ctx := context.Background()
	ctx = common.WithRuntime(ctx, rt)
	rt.Set("dns", common.Bind(rt, New(), &ctx))

	t.Run("Custom dns", func(t *testing.T) {
		_, err := common.RunString(rt, `var params={server:'8.8.8.8'};
		var resp = dns.resolve("google.com",params);
		if(resp==null){
			throw new Error("resp is not expected"+resp);
		}`)
		assert.NoError(t, err)
	})

	t.Run("Default dns", func(t *testing.T) {
		_, err := common.RunString(rt, `
		var resp = dns.resolve("google.com");
		if(resp==null){
			throw new Error("resp is not expected"+resp);
		}`)
		assert.NoError(t, err)
	})
}
