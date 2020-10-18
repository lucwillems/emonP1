package Handlers

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"os"
	"time"
)

type MqttBus struct {
	client mqtt.Client
}

func NewMqttBus(clientId string, user string, password string, url string) (*MqttBus, error) {
	var bus MqttBus
	opts := mqtt.NewClientOptions()
	opts.AddBroker(url)
	opts.SetUsername(user)
	opts.SetPassword(password)
	opts.SetClientID(clientId)
	bus.client = mqtt.NewClient(opts)
	token := bus.client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		bus.client.Disconnect(0)
		return nil, err
	}
	return &bus, nil
}

/* implement MsgBus interface for mqtt bus */
func (c *MqttBus) Close() {
	if c.IsConnected() {
		c.client.Disconnect(1)
		c.client = nil
	}
}

func (c *MqttBus) IsConnected() bool {
	return c.client != nil
}

func (c *MqttBus) Publish(id string, value interface{}) error {
	fmt.Fprintf(os.Stdout, "%s : %s", id, value)
	return nil
}

func (c *MqttBus) MsgBus() *MsgBus {
	x := MsgBus(c)
	return &x
}
