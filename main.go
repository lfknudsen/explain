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
var dbDir = filepath.Join(filepath.Dir(executable), "man")
var manDirs = []string { dbDir }

const description = `explain <command>        Print an explanation of the given command.
explain --edit           Open the database file.
explain --help           Print this message and exit.

Each man page is a TOML file. The file should have the following format:
Name="<name of the command>"
Desc="<description>"
Source="<where to download it>"

Example:
Name="install"
Desc="Move the given files to the same directory as this executable."
Source="https://codeberg.org/lfknudsen/install"
`

func printHelp() {
	fmt.Print(description)
}

func initDB() {
	_ = os.MkdirAll(dbDir, 644)
	dbFile := filepath.Join(dbDir, "explain.toml")
	f, err := os.OpenFile(dbFile, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return
	}
	base := filepath.Base(executable)
	stem := strings.TrimSuffix(base, ".exe")
	out := fmt.Sprintf("Name=\"explain\"\n", stem)
	out += fmt.Sprintf("Desc = '''\n%s'''\n", description)
	out += fmt.Sprintf("Source=\"https://codeberg.org/lfknudsen/explain\"\n")
	_, err = f.WriteString(out)
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
	}
	_ = f.Close()
}

func next() string {
	out := os.Args[0]
	os.Args = os.Args[1:]
	return out
}

func hasNext() bool {
	return len(os.Args) != 0
}

func main() {
	next()
	if !hasNext() {
		printHelp()
		os.Exit(0)
	}

	env, found := os.LookupEnv("PATH")
	if found {
		paths := strings.Split(env, ";")
		if len(paths) == 1 {
			paths = strings.Split(env, ":")
		}
		manDirs = append(manDirs, paths...)
	}

	initDB()

	var err error
	var path string

	for hasNext() {
		arg := next()
		switch arg {
		case "-h", "--help":
			printHelp()
			os.Exit(0)
		case "-e", "--edit": 
			if !hasNext() {
				os.Exit(1)
			}
			name := next()
			path = findFile(name)
			if path == "" {
				fmt.Printf("Could not find %s. Try executing 'explain --add %s' to create the file.", name, name) 
				os.Exit(1)
			}
			err = open(path)
			if err != nil {
				_, _ = fmt.Fprint(os.Stderr, err)
				os.Exit(1)
			}
		case "-a", "--add":
			if !hasNext() {
				os.Exit(1)
			}
			path = filepath.Join(dbDir, next()+".toml")
			_, _ = os.OpenFile(path, os.O_CREATE|os.O_EXCL, 644)
			err = open(path)
			if err != nil {
				_, _ = fmt.Fprint(os.Stderr, err)
				os.Exit(1)
			}
		case "--paths":
			for _, p := range manDirs {
				fmt.Println(p)
			}
		default:
			getCommand(arg)
		}
	}
}

func findFile(cmdName string) string {
	for _, dir := range manDirs {
		path := filepath.Join(dir, cmdName)
		_, err := os.Stat(path)
		if err == nil {
			return path
		}
	}
	return ""
}

func findAndRead(cmdName string) *[]byte {
	cmdName += ".toml"
	for _, dir := range manDirs {
		path := filepath.Join(dir, cmdName)
		buffer, err := os.ReadFile(path)
		if err == nil {
			return &buffer
		}
	}
	return nil
}

func getCommand(name string) {
	buffer := findAndRead(name)
	if buffer == nil {
		fmt.Printf("Command '%s' not found.\n", name)
		os.Exit(1)
	}
	decoded := new(Command)
	err := toml.Unmarshal(*buffer, &decoded)
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(*decoded)
}

type Command struct {
	Name   string
	Desc   string
	Source string
}

func (c Command) String() string {
	name := c.Name
	if c.Name == "" {
		name = "No name"
	}
	desc := c.Desc
	if desc == "" {
		desc = "No description."
	}
	source := c.Source
	if source != "" {
		source = "\nSource: " + source
	}
	return fmt.Sprintf("%s:\n%s%s", name, desc, source)
}
