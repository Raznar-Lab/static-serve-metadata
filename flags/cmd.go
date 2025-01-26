package flags

import "flag"

type Flags struct {
	ConfigPath string
	Debug bool
}

func (c *Flags) Parse() {
	flag.StringVar(&c.ConfigPath, "config", "config.yml", "Path to config file")
	flag.BoolVar(&c.Debug, "debug", false, "Enable debug mode")
	flag.Parse()
}

func New() (cmdFlags *Flags) {
	cmdFlags = &Flags{}

	return
}
