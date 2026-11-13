package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

var executable = getExecutable()
var executableDir = filepath.Dir(executable)
var dbFile = filepath.Join(executableDir, ".explain.toml")

const description = `explain <command>        Print an explanation of the given command.
explain --edit           Open the database file.
explain --help           Print this message and exit.

The database is a TOML file. Each entry in the file should have the following format:
[<key>]
Name="<name of the command>"
Desc="<description>"
Source="<where to download it>"

Example:
[install]
Name="install"
Desc="Move the given files to the same directory as this executable."
Source="https://codeberg.org/lfknudsen/install"
`

func printHelp() {
	fmt.Print(description)
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
		f, err := os.OpenFile(dbFile, os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			_, _ = fmt.Fprint(os.Stderr, err)
		}
		base := filepath.Base(executable)
		stem := strings.TrimSuffix(base, ".exe")
		out := fmt.Sprintf("[%s]\nName=\"explain\"\n", stem)
		out += fmt.Sprintf("Desc = '''\n%s'''\n", description)
		out += fmt.Sprintf("Source=\"https://codeberg.org/lfknudsen/explain\"\n")
		_, err = f.WriteString(out)
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
	cmd, ok := decoded[os.Args[1]]
	if ok {
		fmt.Println(cmd)
	} else {
		fmt.Printf("Command '%s' not found.\n", os.Args[1])
	}
}

type Command struct {
	Name   string
	Desc   string
	Source string
}

func (c Command) String() string {
	return fmt.Sprintf("%s:\n%s\nSource: %s", c.Name, c.Desc, c.Source)
}
