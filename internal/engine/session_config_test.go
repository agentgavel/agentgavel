package engine

import (
	"testing"
)

func TestSessionConfigForMode_OracleInjectsBaseURL(t *testing.T) {
	t.Parallel()
	const oracleURL = "http://127.0.0.1:9876"
	cfg, err := SessionConfigForMode(ModeOracle, ModeEndpoints{
		OracleURL: oracleURL,
		ModelURL:  "https://api.openai.com/v1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.ModelBaseURL != oracleURL {
		t.Fatalf("ModelBaseURL = %q, want oracle URL %q", cfg.ModelBaseURL, oracleURL)
	}
	if cfg.RunMode != ModeOracle {
		t.Fatalf("RunMode = %q, want %q", cfg.RunMode, ModeOracle)
	}
}

func TestSessionConfigForMode_ModelUsesModelURL(t *testing.T) {
	t.Parallel()
	const modelURL = "https://api.openai.com/v1"
	cfg, err := SessionConfigForMode(ModeModel, ModeEndpoints{
		OracleURL: "http://127.0.0.1:9876",
		ModelURL:  modelURL,
	})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.ModelBaseURL != modelURL {
		t.Fatalf("ModelBaseURL = %q, want %q", cfg.ModelBaseURL, modelURL)
	}
	if cfg.RunMode != ModeModel {
		t.Fatalf("RunMode = %q, want %q", cfg.RunMode, ModeModel)
	}
}

func TestApplyModeConfig_OracleOnSession(t *testing.T) {
	t.Parallel()
	const oracleURL = "http://oracle.test:8080"
	sess := &Session{ID: "s1", Seed: 1, Mode: ModeOracle}
	if err := ApplyModeConfig(sess, ModeEndpoints{OracleURL: oracleURL}); err != nil {
		t.Fatal(err)
	}
	if sess.Config.ModelBaseURL != oracleURL {
		t.Fatalf("sess.Config.ModelBaseURL = %q, want %q", sess.Config.ModelBaseURL, oracleURL)
	}
	if sess.Config.RunMode != ModeOracle {
		t.Fatalf("sess.Config.RunMode = %q, want %q", sess.Config.RunMode, ModeOracle)
	}
}

func TestSessionConfigForMode_Errors(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name      string
		mode      string
		endpoints ModeEndpoints
	}{
		{name: "unknown mode", mode: "bogus", endpoints: ModeEndpoints{OracleURL: "http://x", ModelURL: "http://y"}},
		{name: "oracle missing url", mode: ModeOracle, endpoints: ModeEndpoints{}},
		{name: "model missing url", mode: ModeModel, endpoints: ModeEndpoints{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := SessionConfigForMode(tc.mode, tc.endpoints)
			if err == nil {
				t.Fatal("expected error")
			}
		})
	}
}
