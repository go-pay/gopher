package micro

type Broker interface {
	Address() []string
}

//func mqttBroker(c Broker) broker.Broker {
//	return mqtt.NewBroker(func(options *broker.Options) {
//		options.Addrs = c.Address()
//		options.Secure = true
//	})
//}
