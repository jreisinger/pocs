package main

import "testing"

func TestParseOrder(t *testing.T) {
	for detail, wantErr := range map[string]bool{
		`{"customerId":"c1","orderId":"o1","amount":42}`: false,
		`{"orderId":"o1","amount":42}`:                   true,
		`{"customerId":"c1","amount":42}`:                true,
		`{"customerId":"c1","orderId":"o1","amount":0}`:  true,
		`{"customerId":"c1","orderId":"o1","amount":-1}`: true,
		`not json`: true,
	} {
		if _, err := parseOrder([]byte(detail)); (err != nil) != wantErr {
			t.Errorf("parseOrder(%s) err = %v, wantErr %v", detail, err, wantErr)
		}
	}
}
