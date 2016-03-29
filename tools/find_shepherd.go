package tools

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

// GetShepherd takes in an absolute path to a repo and returns a randomly chosen shepherd.
// If there are no shepherds for the repo, the empty string is returned.
func GetShepherd(pwd string) (string, error) {
	if pwd == "" {
		var pwdErr error
		pwd, pwdErr = os.Getwd()
		if pwdErr != nil {
			return "", pwdErr
		}
	}
	repo := filepath.Base(pwd)
	launchFile := fmt.Sprintf("launch/%s.yml", repo)
	ymlLoc := filepath.Join(pwd, launchFile)

	// If there's no launch file, there are no shepherds
	_, err := os.Stat(ymlLoc)
	if os.IsNotExist(err) {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	launchYML, ymlErr := ioutil.ReadFile(ymlLoc)
	if ymlErr != nil {
		log.Fatalf("Error reading yaml file: %s\n", ymlErr)
	}

	// config is used for unmarshaling shepherds if they exist
	var config struct {
		Shepherds []string `yaml:"shepherds,omitempty"`
	}
	if unmarshalErr := yaml.Unmarshal(launchYML, &config); unmarshalErr != nil {
		return "", unmarshalErr
	}

	shepherds := config.Shepherds
	if len(shepherds) == 0 {
		return "", nil
	}

	// Setup a random number generator seeded by time
	seed := time.Now().Unix()
	generator := rand.New(rand.NewSource(seed))
	randEl := generator.Intn(len(shepherds))
	return shepherds[randEl], nil
}
