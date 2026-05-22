package main

import (
	"flag"
	"fmt"
	"os"
)

type Flags struct {
	TLSCrt string
	TLSKey string
}

func ParseFlags() (Flags, error) {

	f := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	tlsCrt := f.String("tls-crt", getStringEnv("MWH_TLS_CRT", ""),
		"tls certificate to be used by this service")
	tlsKey := f.String("tls-key", getStringEnv("MWH_TLS_KEY", ""),
		"tls key to be used by this service")
	if err := f.Parse(os.Args[1:]); err != nil {
		return Flags{}, fmt.Errorf("parse args: %w", err)
	}

	flags := Flags{
		TLSCrt: stringValue(tlsCrt),
		TLSKey: stringValue(tlsKey),
	}

	return flags, nil
}

func getStringEnv(envName string, defaultValue string) string {

	env, ok := os.LookupEnv(envName)
	if !ok {
		return defaultValue
	}
	return env
}

func stringValue(v *string) string {

	if v == nil {
		return ""
	}
	return *v
}
