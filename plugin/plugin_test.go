package plugin_test

import (
	"github.com/hyperits/tlog/plugin"
	"testing"

	"github.com/stretchr/testify/assert"

)

type mockPlugin struct{}

func (p *mockPlugin) Type() string {
	return pluginType
}

func (p *mockPlugin) Setup(name string, decoder plugin.Decoder) error {
	return nil
}

func TestGet(t *testing.T) {
	plugin.Register(pluginName, &mockPlugin{})
	// test duplicate registration
	plugin.Register(pluginName, &mockPlugin{})
	p := plugin.Get(pluginType, pluginName)
	assert.NotNil(t, p)

	pNo := plugin.Get("notexist", pluginName)
	assert.Nil(t, pNo)
}
