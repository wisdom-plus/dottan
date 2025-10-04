package main

import (
	"bufio"
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
)

type Settings struct {
	Theme     string `toml:"theme"`
	Branch    string `toml:"branch"`
	GithubURL string `toml:"github_url"`
}

type Config struct {
	Username string   `toml:"username"`
	Email    string   `toml:"email"`
	Settings Settings `toml:"settings"`
}

func WriteConfigFile(path string, conf *Config) (*Config, error) {
	file, err := os.Create("./config.toml")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	err1 := toml.NewEncoder(file).Encode(conf)
	if err1 != nil {
		panic(err1)
	}
	return conf, nil
}

func ReadInput() string {
	fmt.Println("your github url:")
	reader := bufio.NewScanner(os.Stdin)
	reader.Scan()
	fmt.Println(reader.Text())
	return reader.Text()
}

func main() {
	var conf Config
	_, err := toml.DecodeFile("./config.toml", &conf)
	if err != nil {
		panic(err)
	}
	conf.Settings.GithubURL = ReadInput()
	WriteConfigFile("./config.toml", &conf)
	fmt.Println("Username:", conf.Username)
	fmt.Println("Email:", conf.Email)
	fmt.Println("Settgins[theme]:", conf.Settings.Theme)
	fmt.Println("Settgins[branch]:", conf.Settings.Branch)
}
