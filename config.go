package main

type Config struct {
	Host      string
	Port      int
	StateFile string
	CertFile  string
	KeyFile   string
	APIKeys   []string
}

var (
	config = &Config{
		Host:      "",
		Port:      0,
		StateFile: "state.yaml",
        CertFile:  "",
		KeyFile:   "",
		APIKeys:   []string{
            "",
                   
        },
	}
)
