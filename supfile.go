package sup

import (
	"errors"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

var (
	ErrNetworkNotFound         = errors.New("Network not found")
	ErrInvalidCommandReference = errors.New("Invalid command or target")
)

// Supfile represents the Stackup configuration YAML file.
type Supfile struct {
	Networks map[string]Network  `yaml:"networks"`
	Commands map[string]Command  `yaml:"commands"`
	Targets  map[string][]string `yaml:"targets"`
	Env      map[string]string   `yaml:"env"`
}

// Network is group of hosts with extra custom env vars.
type Network struct {
	Hosts         []string          `yaml:"hosts"`
	HostInventory string            `yaml:"host_inventory"`
	Env           map[string]string `yaml:"env"`
}

// Command represents command(s) to be run remotely.
type Command struct {
	Name   string   `yaml:-`        // Command name.
	Desc   string   `yaml:"desc"`   // Command description.
	Run    string   `yaml:"run`     // Command(s) to be run remotelly.
	Script string   `yaml:"script"` // Load command(s) from script and run it remotelly.
	Upload []Upload `yaml:"upload"` // See below.
	Stdin  bool     `yaml:"stdin"`  // Attach localhost STDOUT to remote commands' STDIN?
}

// Upload represents file copy operation from localhost Src path to Dst
// path of every host in a given Network.
type Upload struct {
	Src string `yaml:"src"`
	Dst string `yaml:"dst"`
}

func (supFile *Supfile) NetworkHostCount(name string) int {
	network, ok := supFile.Networks[name]
	if !ok {
		return 0
	}
	return len(network.Hosts)
}

func (supFile *Supfile) GetNetwork(name string) Network {
	return supFile.Networks[name]
}

func (supFile *Supfile) LoadNetwork(name string) error {
	network, ok := supFile.Networks[name]

	if !ok {
		return ErrNetworkNotFound
	}

	if network.HostInventory == "" {
		return nil
	}

	rawHosts, err := CommandOutput(network.HostInventory)
	if err != nil {
		return err
	}

	hosts := strings.Split(rawHosts, "\n")
	network.Hosts = append(network.Hosts, hosts...)

	supFile.Networks[name] = network
	return nil
}

func (supFile *Supfile) HasEntryPoint(name string) bool {
	_, ok := supFile.Targets[name]
	if ok {
		return true
	}

	_, ok = supFile.Commands[name]
	return ok
}

// NewSupfile parses configuration file and returns Supfile or error.
func NewSupfile(file string) (*Supfile, error) {
	var conf Supfile
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}

func (supFile *Supfile) CollectCommands(name string) ([]*Command, error) {
	ret := []*Command{}

	cmd, isCommand := supFile.Commands[name]

	if isCommand {
		ret = append(ret, &cmd)
		return ret, nil
	}

	targetList, isTarget := supFile.Targets[name]
	if !isTarget {
		return ret, ErrInvalidCommandReference
	}

	for _, targetOrCmd := range targetList {
		cmdList, err := supFile.CollectCommands(targetOrCmd)
		if err != nil {
			return ret, err
		}
		for _, cmd := range cmdList {
			if cmd.Name == "" {
				cmd.Name = name
			}
			ret = append(ret, cmd)
		}
	}

	return ret, nil
}
