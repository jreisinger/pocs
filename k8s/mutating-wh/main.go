package main

import (
	"errors"
	"fmt"
	"log/slog"
	"mutating-wh/internal/k8s"
	"mutating-wh/internal/server"
	"net/http"
	"os"
	"path/filepath"
)

var (
	tlsCertFile string
	tlsKeyFile  string
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		slog.Error(fmt.Sprintf("get user home dir: %v", err))
		os.Exit(1)
	}
	tlsCertFile = filepath.Join(homeDir, "tls.crt")
	tlsKeyFile = filepath.Join(homeDir, "tls.key")
}

func main() {
	flags, err := ParseFlags()
	if err != nil {
		slog.Error(fmt.Sprintf("parse flags: %v", err))
		os.Exit(1)
	}
	slog.Info(fmt.Sprintf("starting mutating admission webhook with flags: %s", flags))

	if err := os.WriteFile(tlsCertFile, []byte(flags.TLSCrt), 0640); err != nil {
		slog.Error(fmt.Sprintf("write tls.crt: %v", err))
		os.Exit(1)
	}
	if err := os.WriteFile(tlsKeyFile, []byte(flags.TLSKey), 0600); err != nil {
		slog.Error(fmt.Sprintf("write tls.key: %v", err))
		os.Exit(1)
	}

	mutateFn := func(body []byte) ([]byte, error) {
		maxRequests := k8s.MaxRequests{
			Memory:       32 * 1024 * 1024,
			MemoryString: "32Mi",
		}
		return k8s.Mutate(body, maxRequests)
	}

	if err := server.ListenAndServeTLS(mutateFn, tlsCertFile, tlsKeyFile); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error(err.Error())
		os.Exit(1)
	}
	slog.Info("server closed")
}
