package config

type Config struct {
	DefaultProfile string             `toml:"default_profile"`
	Profiles       map[string]Profile `toml:"profiles`
}

type Profile struct {
	GitHubURL string `toml:"github_url"`
}
