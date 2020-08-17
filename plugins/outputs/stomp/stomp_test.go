package stomp

import (
	"testing"

	"github.com/influxdata/telegraf/plugins/serializers"
	"github.com/influxdata/telegraf/testutil"

	"github.com/stretchr/testify/require"
)

// TestiConnectAndWrite ...
func TestiConnectAndWrite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	var url = testutil.GetLocalHost() + ":61613"
	st := &STOMP{
		Host:      url,
		Username:  "",
		Password:  "",
		QueueName: "test_queue",
		SSL:       false,
		serialize: serializers.NewJsonSerializer(),
	}
	err := st.Connect()
	require.NoError(t, err)

	err = st.Write(testutil.MockMetrics())
	require.NoError(t, err)
}
