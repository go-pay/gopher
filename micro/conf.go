package micro

type EtcdRegistry struct {
	Addrs []string
}

func (e *EtcdRegistry) Address() []string {
	return e.Addrs
}

// #==============Broker==============
type NsqBroker struct {
	Addrs []string
}

func (n *NsqBroker) Address() []string {
	return n.Addrs
}

type KafkaBroker struct {
	Addrs []string
}

func (k *KafkaBroker) Address() []string {
	return k.Addrs
}

type RedisBroker struct {
	Addrs []string
}

func (r *RedisBroker) Address() []string {
	return r.Addrs
}

type MqttBroker struct {
	Addrs []string
}

func (m *MqttBroker) Address() []string {
	return m.Addrs
}
