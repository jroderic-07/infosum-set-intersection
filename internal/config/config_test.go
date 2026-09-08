package config

import (
	"flag"
	"os"
	"testing"
)

func TestArrayFlags_Set(t *testing.T) {
	var flags arrayFlags

	if err := flags.Set("a.csv"); err != nil {
		t.Fatalf("Set returned error: %v", err)
	}
	if err := flags.Set("b.csv"); err != nil {
		t.Fatalf("Set returned error: %v", err)
	}

	want := []string{"a.csv", "b.csv"}
	if len(flags) != len(want) {
		t.Fatalf("got %d values, want %d", len(flags), len(want))
	}
	for i := range want {
		if flags[i] != want[i] {
			t.Errorf("flags[%d] = %q, want %q", i, flags[i], want[i])
		}
	}
}

func TestArrayFlags_String(t *testing.T) {
	var flags arrayFlags
	if got := flags.String(); got != "array of flags" {
		t.Errorf("String() = %q, want %q", got, "array of flags")
	}
}

func resetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
}

func TestLoadConfigurationCLI(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantCPU int
		wantVer bool
		wantIn  []string
	}{
		{
			name:    "defaults",
			args:    []string{"compare"},
			wantCPU: 0,
			wantVer: false,
			wantIn:  nil,
		},
		{
			name:    "version flag",
			args:    []string{"compare", "-version"},
			wantVer: true,
		},
		{
			name:    "cpu and input files",
			args:    []string{"compare", "-cpuNum", "4", "-in", "a.csv", "-in", "b.csv"},
			wantCPU: 4,
			wantIn:  []string{"a.csv", "b.csv"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetFlags()

			oldArgs := os.Args
			t.Cleanup(func() { os.Args = oldArgs })
			os.Args = tt.args

			cfg := LoadConfigurationCLI()

			if cfg.CpuNum != tt.wantCPU {
				t.Errorf("CpuNum = %d, want %d", cfg.CpuNum, tt.wantCPU)
			}
			if cfg.ShowVersion != tt.wantVer {
				t.Errorf("ShowVersion = %v, want %v", cfg.ShowVersion, tt.wantVer)
			}
			if len(cfg.InputFilePaths) != len(tt.wantIn) {
				t.Fatalf("InputFilePaths = %v, want %v", cfg.InputFilePaths, tt.wantIn)
			}
			for i := range tt.wantIn {
				if cfg.InputFilePaths[i] != tt.wantIn[i] {
					t.Errorf("InputFilePaths[%d] = %q, want %q", i, cfg.InputFilePaths[i], tt.wantIn[i])
				}
			}
		})
	}
}
