package main

import (
	"log/slog"
	"os"
)

var logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

func main() {
	notifiers := []Notifier{
		SlackClient{Channel: "general"},
		SlackClient{Channel: "random"},
		EmailClient{Address: "john.doe@example.org"},
		EmailClient{Address: "jane.doe@example.net"},
	}
	sendNotification("hello", notifiers)
}

type Notifier interface {
	Notify(message string) error
	Type() string
	Destination() string
}

func sendNotification(message string, notifiers []Notifier) {
	for _, n := range notifiers {
		logger := logger.With("type", n.Type(), "destination", n.Destination())
		err := n.Notify(message)
		if err != nil {
			logger.Error("notifiaction failed", "error", err)
			continue
		}
		logger.Info("notification sent")
	}
}
