package main

type SlackClient struct {
	Channel string
}

func (c SlackClient) Notify(message string) error {
	return simulateError()
}

func (c SlackClient) Type() string {
	return "slack"
}

func (c SlackClient) Destination() string {
	return c.Channel
}
