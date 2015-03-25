package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	yaml "github.com/Clever/gitbot/Godeps/_workspace/src/gopkg.in/yaml.v2"
)

type Command struct {
	Path string
	Args []string
}

func (c Command) Validate() error {
	if c.Path == "" {
		return fmt.Errorf("command must specify 'path'")
	}
	return nil
}

type Config struct {
	Repos     []string  `yaml:"repos"`
	ChangeCmd Command   `yaml:"change_cmd"`
	PostCmds  []Command `yaml:"post_cmds"`
}

func (c Config) Validate() error {
	if len(c.Repos) == 0 {
		return fmt.Errorf("config must contain a non-empty 'repos' list")
	}
	if err := c.ChangeCmd.Validate(); err != nil {
		return err
	}
	if len(c.PostCmds) == 0 {
		return fmt.Errorf("config must contain a non-empty 'post_cmds' list")
	}
	for _, postcmd := range c.PostCmds {
		if err := postcmd.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("usage: gitbot [config]")
	}

	cfgfile, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer cfgfile.Close()
	cfgfiledata, err := ioutil.ReadAll(cfgfile)
	if err != nil {
		log.Fatal(err)
	}
	var cfg Config
	if err := yaml.Unmarshal(cfgfiledata, &cfg); err != nil {
		log.Fatal(err)
	}

	log.Printf("loaded config: %s", cfg)

}
