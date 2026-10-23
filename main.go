package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

func getExecutable() string {
	e, err := os.Executable()
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	return e
}

var executableDir = filepath.Dir(getExecutable())
var dbFile = filepath.Join(executableDir, ".explain.toml")

func printHelp() {
	fmt.Print("explain <command>		Print an explanation of the given command.\n" +
		"explain --edit           Open the database file.\n" +
		"explain --help           Print this message and exit.\n")
}

func main() {
	if len(os.Args) != 2 {
		printHelp()
		os.Exit(1)
	}
	switch os.Args[1] {
	case "-h", "--help":
		printHelp()
		os.Exit(0)
	case "-e", "--edit":
		err := open(dbFile)
		if err != nil {
			_, _ = fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	contents, err := os.ReadFile(dbFile)
	if err != nil {
		f, err := os.Create(dbFile)
		if err != nil {
			_, _ = fmt.Fprint(os.Stderr, err)
		}
		_ = f.Close()
	}
	decoded := make(map[string]Command)
	err = toml.Unmarshal(contents, &decoded)
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(decoded[os.Args[1]])
}

type Command struct {
	Name   string
	Desc   string
	Source string
}

func (c Command) String() string {
	return fmt.Sprintf("%s: %s\n%s", c.Name, c.Desc, c.Source)
}
