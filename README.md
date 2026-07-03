# zdcli
my cli tool

## Depends on

- lua [[windows download](https://sourceforge.net/projects/luabinaries/files/5.4.2/Tools%20Executables/lua-5.4.2_Win64_bin.zip/download) | [linux download](https://sourceforge.net/projects/luabinaries/files/5.4.2/Tools%20Executables/lua-5.4.2_Linux54_64_bin.tar.gz/download)]

## Install 

### From Source

```sh
git clone https://github.com/zerodoctor/zdcli.git && \
cd zdcli && \
make build install
```

### From Github

```sh
mkdir ~/scripts || true && \
curl -o zd.tar.xz -L https://github.com/ZeroDoctor/zdcli/releases/download/v1.1.0/zd-amd64-unix.tar.xz && \
tar -xvf zd.tar.xz
```

if its the first install its advise to run the command below:

```sh
zd setup && zd setup --ls
```

## Usage

```
NAME:
   zd - A new cli application

USAGE:
   zd [global options] command [command options] 

COMMANDS:
   alert       notifies user when an event happens
   edit, e     edits a lua script
   list, ls    list current lua scripts
   lite        interacts with a sqlite database
   new, n      create a new lua script
   paste       common commands to interact with pastebin.com. May need to login via this cli before use.
   remove, rm  remove a lua script or a directory
   setup       setup lua, editor, and dir configs
   ui          opens a custom terminal emulator
   vault, v    commands that communicates with a vault server
   version     current version of cli
   weed, fs    
   help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```
To run a script simply run:

```sh
   zd test func_unix
```
where "test" is the name of the lua file in './lua/scripts' and "func_unix" is the name of the function inside test.lua.

### Languages

`zd <name> [func] [args]` (also usable as `zdcli`) dispatches by where `<name>` lives, lua first:

- **lua** — `lua/scripts/<name>.lua`, runs `func` via `build-app.lua`. Dotted names map to subdirs (`prod.vpn`).
- **python** — `python/<name>.py`, runs `python3 main.py -s <name> -f <func|main> -a '<json-kwargs>'`.
- **bash** — `sh/<name>.sh`, runs `bash <name>.sh <args>`.

Interpreters and script dirs are set in `zdconfig.toml` (`LuaCmd`, `PythonCmd`, `Root*ScriptDir`).
The seed script dirs (`lua/ python/ sh/`) ship with the repo; a real deployment may replace
`lua/` with its own library (e.g. the `zdscripts` repo) — `make install` won't clobber existing dirs.
