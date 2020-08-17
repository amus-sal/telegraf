package stomp

import (
	"net"

	"github.com/go-stomp/stomp"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/common/tls"
	"github.com/influxdata/telegraf/plugins/outputs"
	"github.com/influxdata/telegraf/plugins/serializers"
)

//STOMP ...
type STOMP struct {
	Host      string `toml:"host"`
	Username  string `toml:"username,omitempty"`
	Password  string `toml:"password,omitempty"`
	QueueName string `toml:"queueName"`
	SSL       bool   `toml:"ssl"`
	tls.ClientConfig
	Conn      *tls.Conn
	NetConn   net.Conn
	Stomp     *stomp.Conn
	serialize serializers.Serializer
}

//Connect ...
func (q *STOMP) Connect() error {
	var err error
	if q.SSL == true {
		tlsCongi, _ := q.ClientConfig.TLSConfig()

		q.Conn, err = tls.Dial("tcp", q.Host, tlsCongi)

		q.Stomp, err = stomp.Connect(q.Conn, stomp.ConnOpt.HeartBeat(0, 0), stomp.ConnOpt.Login(q.Username, q.Password))
	} else {
		q.NetConn, err = net.Dial("tcp", q.Host)
		q.Stomp, err = stomp.Connect(q.NetConn, stomp.ConnOpt.HeartBeat(0, 0), stomp.ConnOpt.Login(q.Username, q.Password))

	}
	if err != nil {
		println("cannot connect to server", err.Error())
		return err
	}

	if err != nil {
		println(err.Error())
		return err
	}
	println("STOMP Connected...")
	return nil
}

//SetSerializer ...
func (q *STOMP) SetSerializer(serializer serializers.Serializer) {
	q.serialize = serializer
}

//Write ...
func (q *STOMP) Write(metrics []telegraf.Metric) error {
	for _, metric := range metrics {
		values, err := q.serialize.Serialize(metric)
		if err != nil {
			return err
		}
		err = q.Stomp.Send(q.QueueName, "text/plain",
			[]byte(values), nil)
		if err != nil {
			panic(err)
			return err
		}
	}
	return nil
}

//Close ...
func (q *STOMP) Close() error {
	println("Closiong is starting .....")
	q.Stomp.Disconnect()
	q.Conn.Close()
	return nil
}

//SampleConfig ...
func (q *STOMP) SampleConfig() string {
	return `ok = true`
}

//Description ...
func (q *STOMP) Description() string {
	return "Telegraf Output Plugin For Stomp"
}
func init() {
	outputs.Add("stomp", func() telegraf.Output { return &STOMP{} })
}
