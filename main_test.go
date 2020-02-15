package main

import "testing"


func TestProducer(t *testing.T) {
	testCases := []struct {
		desc string
	}{
		{
			desc: "success",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

		})
	}
}
