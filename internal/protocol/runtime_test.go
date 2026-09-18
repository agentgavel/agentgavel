package protocol

import "testing"

func TestNormalizeRuntime(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"", RuntimeStub},
		{"stub", RuntimeStub},
		{"LIVE", RuntimeLive},
		{" harness ", RuntimeHarness},
		{"weird", RuntimeStub},
	}
	for _, tc := range cases {
		if got := NormalizeRuntime(tc.in); got != tc.want {
			t.Fatalf("NormalizeRuntime(%q)=%q want %q", tc.in, got, tc.want)
		}
	}
}

func TestValidRuntime(t *testing.T) {
	if ValidRuntime("") {
		t.Fatal("empty should not be ValidRuntime")
	}
	if !ValidRuntime("live") {
		t.Fatal("live should be valid")
	}
}
