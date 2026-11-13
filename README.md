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
