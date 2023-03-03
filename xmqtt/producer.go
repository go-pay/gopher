package xmqtt

// Publish 推送消息
func (c *Client) Publish(topic string, qos QosType, payload any) error {
	token := c.Mqtt.Publish(topic, byte(qos), false, payload)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}
