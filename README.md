# zdcli
my cli tool

## Depends on

- lua [[windows download](https://sourceforge.net/projects/luabinaries/files/5.4.2/Tools%20Executables/lua-5.4.2_Win64_bin.zip/download) | [linux download](https://sourceforge.net/projects/luabinaries/files/5.4.2/Tools%20Executables/lua-5.4.2_Linux54_64_bin.tar.gz/download)]

## Install 

### From Go

```sh
go install github.com/zerodoctor/zdcli
```

### From Source

```sh
git clone https://github.com/zerodoctor/zdcli.git
```

```sh
make build install
```

## Usage

```
NAME:
   zd - A new cli application

USAGE:
   zd [global options] command [command options] [arguments...]

COMMANDS:
   alert       notifies user when an event happens
   edit, e     edits a lua script
   list, ls    list current lua scripts
   new, n      create a new lua script
   paste       common commands to interact with pastebin.com. May need to login via this cli before use.
   remove, rm  remove a lua script or a directory
   setup       setup lua, editor, and dir configs
   ui          opens a custom terminal emulator
   vault, v    commands that communicates with a vault server
   help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help (default: false)
```
To run a script simply run:

```sh 
   zd test func_unix
```
where "test" is the name of the lua file in './lua/scripts' and "func_unix" is the name of the function inside test.lua
