package main

type EmailClient struct {
	Address string
}

func (c EmailClient) Notify(message string) error {
	return simulateError()
}

func (c EmailClient) Type() string {
	return "email"
}

func (c EmailClient) Destination() string {
	return c.Address
}
