package main

type Config struct {
	GlobalConfig Globalconfig `toml:"global_config"`
	UserConfig   UserConfig   `toml:"user_config"`
}

type GlobalConfig struct {
	repo_name   string
	repo_url    string
	repo_owner  string
	ignore_file string
}

type UserConfig struct {
	username string
}
