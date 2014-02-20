package monk

// Config holds various configuration options used throughout a Context and its
// collaborators.
type Config struct {
	Fingerprint bool
	AssetRoot   string
}

func NewConfig() *Config {
	return &Config{
		Fingerprint: false,
		AssetRoot:   "/assets/",
	}
}
