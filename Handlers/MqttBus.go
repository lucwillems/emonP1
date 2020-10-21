package Handlers

import (
	"emonP1/P1"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"os"
	"time"
)

type MqttBus struct {
	debug  bool
	client mqtt.Client
}

var publishKeys = [...]string{
	P1.OBISTypeElectricityTariffIndicator,
	P1.OBISTypePowerDelivered,
	P1.OBISTypePowerGenerated,
	P1.OBISTypeElectricityDeliveredTariff1,
	P1.OBISTypeElectricityDeliveredTariff2,
	P1.OBISTypeElectricityGeneratedTariff1,
	P1.OBISTypeElectricityGeneratedTariff2,
	P1.OBISTypeGasTempNotCorrectedDelivered,
	//	P1.OBISTypeInstantaneousPowerDeliveredL1,
	//	P1.OBISTypeInstantaneousPowerDeliveredL2,
	//	P1.OBISTypeInstantaneousPowerDeliveredL3,
	//	P1.OBISTypeInstantaneousPowerGeneratedL1,
	//	P1.OBISTypeInstantaneousPowerGeneratedL2,
	//	P1.OBISTypeInstantaneousPowerGeneratedL3,
	P1.OBISTypeInstantaneousVoltageL1,
	P1.OBISTypeInstantaneousVoltageL2,
	P1.OBISTypeInstantaneousVoltageL3,
	P1.OBISTypeInstantaneousCurrentL1,
	P1.OBISTypeInstantaneousCurrentL2,
	P1.OBISTypeInstantaneousCurrentL3,
}

func NewMqttBus(clientId string, user string, password string, url string) (*MqttBus, error) {
	var bus MqttBus
	opts := mqtt.NewClientOptions()
	opts.AddBroker(url)
	opts.SetAutoReconnect(true)
	opts.SetUsername(user)
	opts.SetPassword(password)
	opts.SetClientID(clientId)
	opts.SetCleanSession(true)
	opts.SetKeepAlive(60 * 4 * time.Second)
	bus.client = mqtt.NewClient(opts)
	fmt.Fprintf(os.Stdout, "connect %s ...\n", url)
	token := bus.client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		bus.client.Disconnect(0)
		return nil, err
	}
	fmt.Fprintf(os.Stdout, "connected to mqtt\n")
	return &bus, nil
}

/* implement MsgBus interface for mqtt bus */
func (c *MqttBus) Close() {
	if c.IsConnected() {
		c.client.Disconnect(1)
		c.client = nil
	}
}

func (c *MqttBus) Debug(debug bool) {
	c.debug = debug
}

func (c *MqttBus) IsConnected() bool {
	return c.client != nil
}

func (c *MqttBus) Publish(cnt uint64, telegram *P1.Telegram) error {
	//only send each 10 second
	if cnt%10 == 0 {
		c._publish(c._topic("cnt"), fmt.Sprintf("%d", cnt))
		//publish well known counters
		for _, k := range publishKeys {
			if i, ok := telegram.Get(k); ok == true {
				c._publishCOSEM(i)
				if c.debug {
					fmt.Fprintf(os.Stdout, "%v\n", i)
				}
			} else {
				if c.debug {
					fmt.Fprintf(os.Stderr, "key not found: {}\n", k)
				}
			}
		}
	}
	return nil
}

func (c *MqttBus) _topic(name string) string {
	topic := fmt.Sprintf("emon/%s/%s", "p1", name)
	return topic

}

func (c *MqttBus) _publishCOSEM(item *P1.COSEMInstance) mqtt.Token {
	ok := item.IsValid() && item.HasValue()
	if len(item.Queue()) > 0 && ok {
		t := c._publish(c._topic(item.Queue()), item.AsString())
		t.WaitTimeout(200 * time.Millisecond)
		if t.Error() != nil {
			fmt.Fprintf(os.Stderr, "mqtt: %s", t.Error())
		}
		return t
	}
	return nil
}

func (c *MqttBus) _publish(topic string, value string) mqtt.Token {
	t := c.client.Publish(topic, 0, false, value)
	t.WaitTimeout(200 * time.Millisecond)
	if t.Error() != nil {
		fmt.Fprintf(os.Stderr, "mqtt: %s", t.Error())
	}
	return t
}
