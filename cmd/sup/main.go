package main

import (
	"log"
	"os"

	"github.com/pressly/sup"
)

// usage prints help for an arg and exits.
func usage(conf *sup.Supfile, arg int) {
	log.Println("Usage: sup <network> <target/command>\n")
	switch arg {
	case 1:
		// <network> missing, print available hosts.
		log.Println("Available networks (from ./Supfile):")
		for name, network := range conf.Networks {
			log.Printf("- %v\n", name)
			for _, host := range network.Hosts {
				log.Printf("   - %v\n", host)
			}
			if network.HostInventory != "" {
				log.Printf("   - script: %v\n", network.HostInventory)
			}
		}
	case 2:
		// <target/command> not found or missing,
		// print available targets/commands.
		log.Println("Available targets (from Supfile):")
		for name, commands := range conf.Targets {
			log.Printf("- %v", name)
			for _, cmd := range commands {
				log.Printf("\t%v\n", cmd)
			}
		}
		log.Println()
		log.Println("Available commands (from Supfile):")
		for name, cmd := range conf.Commands {
			log.Printf("- %v\t%v", name, cmd.Desc)
		}
	}
	os.Exit(1)
}

// parseArgs parses os.Args and returns network and commands to be run.
// On error, it prints usage and exits.
func parseArgsOrDie(conf *sup.Supfile) (*sup.Network, []*sup.Command) {
	var commands []*sup.Command

	// Check for the first argument first
	if len(os.Args) < 2 {
		usage(conf, len(os.Args))
	}

	err := conf.LoadNetwork(os.Args[1])
	if err != nil {
		log.Println(err)
		usage(conf, 1)
	}

	// Does <network> have any hosts?
	if conf.NetworkHostCount(os.Args[1]) == 0 {
		log.Printf("No hosts specified for network \"%v\"", os.Args[1])
		usage(conf, 1)
	}

	// Check for the second argument
	if len(os.Args) < 3 {
		usage(conf, len(os.Args))
	}

	// Does the <target/command> exist?
	isEntryPoint := conf.HasEntryPoint(os.Args[2])
	if !isEntryPoint {
		log.Printf("Unknown target/command \"%v\"\n\n", os.Args[2])
		usage(conf, 2)
	}

	commands, err = conf.CollectCommands(os.Args[2])
	if err != nil {
		log.Println(err)
		usage(conf, 2)
	}

	// Check for extra arguments
	if len(os.Args) != 3 {
		usage(conf, len(os.Args))
	}

	network := conf.GetNetwork(os.Args[1])
	return &network, commands
}

func main() {
	// Parse configuration file in current directory.
	// TODO: -f flag to pass custom file.
	conf, err := sup.NewSupfile("./Supfile")
	if err != nil {
		log.Fatal(err)
	}

	// Parse network and commands to be run from os.Args.
	network, commands := parseArgsOrDie(conf)

	// Create new Stackup app.
	app, err := sup.New(conf)
	if err != nil {
		log.Fatal(err)
	}

	// Run all the commands in the given network.
	err = app.Run(network, commands...)
	if err != nil {
		log.Fatal(err)
	}
}
