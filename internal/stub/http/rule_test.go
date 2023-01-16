package http

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestReadRule(t *testing.T) {
	data := `
URL: /some-index
#---#
DELAY: 2s
#---#
HEADERS:
Content-Type: application/json
Connection: keep-alive
Server: nginx
#---#
BODY:
{"host": "rOne"}
`
	rule := &Rule{}
	parseOwnRule([]byte(data), rule)
	assert.Equal(t, []byte("/some-index"), rule.Url, "Rule url must be the same")
	assert.Equal(t, time.Second*2, rule.Delay, "Rule delay must be equal")
	assert.Equal(t, []byte("{\"host\": \"rOne\"}"), rule.Body, "Rule body must be equal")
	assert.EqualValues(t,
		map[string]string{
			"Content-Type": "application/json",
			"Connection":   "keep-alive",
			"Server":       "nginx",
		},
		rule.Headers,
		"Rule header must be equal")
}
