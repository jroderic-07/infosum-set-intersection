package config

import "flag"

type arrayFlags []string

// String and Set exist to implement a custom type of flag. Required by Value interface.
// Used for help text.
func (i *arrayFlags) String() string {
	return "array of flags"
}

// Allows arrayFlags to be used multiple times. Called once per -in value supplied to the CLI. Appends each value to a slice.
func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

type ConfigurationCLI struct {
	CpuNum         int
	ShowVersion    bool
	InputFilePaths arrayFlags
}

func LoadConfigurationCLI() *ConfigurationCLI {
	config := ConfigurationCLI{}

	flag.IntVar(&config.CpuNum, "cpuNum", 0, "number of CPU to use")
	flag.BoolVar(&config.ShowVersion, "version", false, "print compare version")
	flag.Var(&config.InputFilePaths, "in", "input file paths")

	flag.Parse()

	return &config
}
