package server

import (
	"fmt"
	"reflect"
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

func Test_hash(t *testing.T) {
	tests := []struct {
		name string
		args []interface{}
		want []byte
	}{
		{
			"one element",
			[]interface{}{
				"test1",
			},
			[]byte{229, 104, 55, 204, 215, 163, 141, 103, 149, 211, 10, 194, 171, 99, 236, 204, 140, 43, 87, 18, 137, 166, 45, 196, 6, 187, 98, 118, 126, 136, 176, 108},
		},
		{
			"two elements",
			[]interface{}{
				"test1",
				"test2",
			},
			[]byte{84, 182, 224, 44, 5, 184, 19, 24, 41, 163, 6, 53, 242, 3, 167, 200, 192, 113, 61, 137, 208, 241, 141, 225, 134, 61, 78, 124, 88, 254, 117, 159},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hash(tt.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("hash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isEq(t *testing.T) {
	type args struct {
		expectedTokenHash []byte
		actualToken       string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"PASS: token matches",
			args{
				[]byte{229, 104, 55, 204, 215, 163, 141, 103, 149, 211, 10, 194, 171, 99, 236, 204, 140, 43, 87, 18, 137, 166, 45, 196, 6, 187, 98, 118, 126, 136, 176, 108},
				"test1",
			},
			true,
		},
		{
			"FAIL: token mismatch",
			args{
				[]byte{84, 182, 224, 44, 5, 184, 19, 24, 41, 163, 6, 53, 242, 3, 167, 200, 192, 113, 61, 137, 208, 241, 141, 225, 134, 61, 78, 124, 88, 254, 117, 159},
				"test1",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isEq(tt.args.expectedTokenHash, tt.args.actualToken); got != tt.want {
				t.Errorf("isEq() = %v, want %v", got, tt.want)
			}
		})
	}
}
