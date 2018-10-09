package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/CognitionFoundry/gohfc"
	"github.com/google/uuid"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Routine []Routine `yaml:"routines"`
}

type Routine struct {
	Name       string `yaml:"name"`
	YAMLClient string `yaml:"yamlClient"`
	YAMLCA     string `yaml:"yamlCA"`
	Peer       string `yaml:"peer"`
	Orderer    string `yaml:"orderer"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	Prefix     string `yaml:"prefix"`
	Channel    string `yaml:"channel"`
	Chaincode  string `yaml:"chaincode"`
	Version    string `yaml:"version"`
	MSPId      string `yaml:"mspid"`
	MsgSize    int    `yaml:"msgSizeBytes"`
	Duration   int    `yaml:"duration"`
}

func NewRoutines(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := new(Config)
	err = yaml.Unmarshal([]byte(data), config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func connectAndSend(routine Routine, results chan<- string) {
	c, err := gohfc.NewFabricClient(routine.YAMLClient)
	if err != nil {
		log.Fatalf("error loading file: %v", err)
	}

	ca, err := gohfc.NewCAClient(routine.YAMLCA, nil)
	if err != nil {
		log.Fatalf("NewCAClient failed: %v", err)
	}

	enrollRequest := gohfc.CaEnrollmentRequest{EnrollmentId: routine.Username, Secret: routine.Password}
	identity, _, err := ca.Enroll(enrollRequest)
	if err != nil {
		log.Fatal("enroll failed: %v", err)
	}

	count := 0
	lEnd := time.Duration(routine.Duration) * time.Second
	for lStart := time.Now(); time.Since(lStart) < lEnd; count++ {
		key := uuid.New().String()
		val := time.Now().Format(time.RFC3339)

		chaincode := &gohfc.ChainCode{
			ChannelId: routine.Channel,
			Type:      gohfc.ChaincodeSpec_GOLANG,
			Name:      routine.Chaincode,
			Version:   routine.Version,
			Args:      []string{"invoke", key, val},
		}

		_, err := c.Invoke(*identity, *chaincode, []string{routine.Peer}, routine.Orderer)
		if err != nil {
			log.Fatal("invoke failed: ", err)
		}

		log.Printf("%8d routine %s: sent %s / %s\n", count, routine.Prefix, key, val)
	}
	finishMsg := fmt.Sprintf("routine %s complete", routine.Prefix)
	results <- finishMsg
}

func main() {
	cfgFile := flag.String("config", "", "yaml input file")

	flag.Parse()

	if *cfgFile == "" {
		fmt.Println("input file required")
		flag.PrintDefaults()
		log.Fatal("no input file")
	}

	cfg, err := NewRoutines(*cfgFile)
	if err != nil {
		log.Fatalf("failed to load file %s: %v", *cfgFile, err)
	}

	cfgLength := len(cfg.Routine)
	results := make(chan string, cfgLength)
	for _, routine := range cfg.Routine {
		go connectAndSend(routine, results)
	}

	for i := 0; i < cfgLength; i++ {
		log.Printf("%v\n", <-results)
	}
	log.Print("complete")
}
