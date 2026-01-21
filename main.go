package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// Command represents a command to run.
type Command struct {
	Path string   `yaml:"path"`
	Args []string `yaml:"args"`
}

// Validate the command is usable.
func (c Command) Validate() error {
	if c.Path == "" {
		return fmt.Errorf("command must specify 'path'")
	}
	return nil
}

// Config for gitbot.
type Config struct {
	Repos      []string  `yaml:"repos"`
	BasePath   string    `yaml:"base_path"`
	ChangeCmds []Command `yaml:"change_cmds"`
	PostCmds   []Command `yaml:"post_cmds"`
}

// NormalizePath will return an absolute path from a relative one, with
// the assumption that is is relative to the location of the configuration file.
func NormalizePath(configPath, commandPath string) string {
	if !strings.HasPrefix(commandPath, ".") {
		return commandPath
	}
	return filepath.Join(filepath.Dir(configPath), commandPath)
}

// Validate that the config is usable.
func (c Config) Validate() error {
	if len(c.Repos) == 0 {
		return fmt.Errorf("config must contain a non-empty 'repos' list")
	}
	for _, changecmd := range c.ChangeCmds {
		if err := changecmd.Validate(); err != nil {
			return err
		}
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
	version := flag.Bool("version", false, "Shows version and exits")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [config]:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Flags: \n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *version {
		fmt.Printf("gitbot %s\n", Version)
		os.Exit(0)
	}

	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	cfgfilePath := flag.Args()[0]
	cfgfile, err := os.Open(cfgfilePath)
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
	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}

	basePath := cfg.BasePath

	for _, repo := range cfg.Repos {
		// clone repo to temp directory
		tempdir, err := ioutil.TempDir(basePath+os.TempDir(), "")
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%s: cloning to %s", repo, tempdir)
		if os.Getenv("GITBOT_LEAVE_TEMPDIRS") == "" {
			defer os.RemoveAll(tempdir)
		}
		clonecmd := exec.Command("git", "clone", "--depth", "1", repo, tempdir)
		clonecmd.Dir = tempdir
		clonecmd.Stdout = os.Stdout
		clonecmd.Stderr = os.Stderr
		if err := clonecmd.Run(); err != nil {
			log.Fatalf("%s: error cloning: %s", repo, err)
		}

		// make changes
		log.Printf("%s: making changes", repo)
		noChangesMade := true
		for _, changecmd := range cfg.ChangeCmds {
			commandPath := NormalizePath(cfgfilePath, changecmd.Path)
			changecmd := exec.Command(commandPath, append(changecmd.Args, tempdir)...)
			var changecmdstdout bytes.Buffer
			changecmd.Stdout = io.MultiWriter(os.Stdout, &changecmdstdout)
			changecmd.Stderr = os.Stderr
			if err := changecmd.Run(); err != nil {
				log.Printf("%s: error running change command: %s"", repo, err)
				log.Printf("%s: no changes to make", repo)
				continue
			}
			noChangesMade = false

			// commit changes
			log.Printf("%s: committing changes", repo)
			gitaddcmd := exec.Command("git", "add", "-A")
			gitaddcmd.Dir = tempdir
			gitaddcmd.Stdout = os.Stdout
			gitaddcmd.Stderr = os.Stderr
			if err := gitaddcmd.Run(); err != nil {
				log.Fatalf("%s: error adding: %s", repo, err)
			}
			commitcmd := exec.Command("git", "commit", "-m", changecmdstdout.String())
			commitcmd.Dir = tempdir
			commitcmd.Stdout = os.Stdout
			commitcmd.Stderr = os.Stderr
			if err := commitcmd.Run(); err != nil {
				log.Fatalf("%s: error committing: %s", repo, err)
			}
		}

		// don't run post commands if none of the change commands made a change
		if noChangesMade {
			continue
		}

		// run post commands
		log.Printf("%s: running post commands", repo)
		for _, postcmd := range cfg.PostCmds {
			postcmdPath := NormalizePath(cfgfilePath, postcmd.Path)
			cmd := exec.Command(postcmdPath, postcmd.Args...)
			cmd.Dir = tempdir
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				log.Fatalf("%s: error running post command: %s", repo, err)
			}
		}
	}
}
