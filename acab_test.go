package main

import "testing"

func DummyTest(t *testing.T) {
	return
}
/*
func TestCheckDomainValidity(t *testing.T) {
	var tests = []struct {
		input  string
		expect bool
	}{
		{"example", true},
		{"example.com", true},
		{"-foo.-", false},
		{"420.test", false},
		{"reallylongdomainlikewaytoolongnoreallywaaaaaytoooooofreeeeeeeeeaaaaaakkkkingglonog.test", false},
	}
	for _, test := range tests {
		got := checkDomainValidity(test.input)
		if got != test.expect {
			t.Errorf("checkDomainValidity(%s) returned %t; expected %t", test.input, got, test.expect)
		}
	}
}
	*/
