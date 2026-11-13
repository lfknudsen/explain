# Explain

Simple little `man`-alternative for your own programmes.

It will create a subdirectory to its own directory to store new
manual pages created through the app.\
When looking up a manual page, it will search in that subdirectory as well
as on your PATH.

## Usage

To look up a command, execute

```
explain <command name>
```

To add a new manual-page for a command, execute

```
explain --add <command name>
```

To edit the existing manual-page for a command, execute

```
explain --edit <command name>
```

## Setup

```
go build
```

## File Format

The manual files are .toml files with three named strings:
- Name
- Desc
- Source

You can execute `explain --edit explain` to see an example.

Markdown are also supported, and will be rendered.

## Example

Executing this:
```
$ explain explain
```

Would output this:
```
explain:
explain <command>            Print an explanation of the given command.
explain --edit <command>     Open the command's manual file for editing.
explain --add <command>      Create a manual page for the command.
explain --paths              Print a list of all directories where the programme will look for manual files.
explain --help               Print this message and exit.

Each man page is a TOML file. The file should have the following format:
Name="<name of the command>"
Desc="<description>"
Source="<where to download it>"
Multi-line strings start and end with """ or '''

Example:
Name="install"
Desc="Move the given files to the same directory as this executable."
Source="https://codeberg.org/lfknudsen/install"

Source: https://codeberg.org/lfknudsen/explain
```
