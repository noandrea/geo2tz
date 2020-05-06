package server

import (
	"fmt"
	"testing"
)

func Test_parseCoordinate(t *testing.T) {
	type c struct {
		val  string
		side string
	}
	tests := []struct {
		ll      c
		want    float32
		wantErr bool
	}{
		{c{"22", Latitude}, 22, false},
		{c{"78.312", Longitude}, 78.312, false},
		{c{"0x429c9fbe", Longitude}, 0, true}, // 78.312 in hex
		{c{"", Longitude}, 0, true},
		{c{"   ", Longitude}, 0, true},
		{c{"2e4", Longitude}, 0, true},
		{c{"not a number", Longitude}, 0, true},
		{c{"-90.1", Latitude}, 0, true},
		{c{"90.001", Latitude}, 0, true},
		{c{"-180.1", Longitude}, 0, true},
		{c{"180.001", Longitude}, 0, true},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprint(tt.ll), func(t *testing.T) {
			got, err := parseCoordinate(tt.ll.val, tt.ll.side)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseCoordinate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseCoordinate() = %v, want %v", got, tt.want)
			}
		})
	}
}
