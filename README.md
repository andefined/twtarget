# twtarget

twtarget (or Twitter Target) is a CLI tool, written in [Go](https://golang.org/), that collects data from twitter API for a given User.


## Installation
You can download the binaries from the [releases](/releases) section, or you can install it with Go.


```bash
go install github.com/andefined/twtarget
```

## How to use
```
NAME:
   twtarget

USAGE:
   twtarget [global options] command [command options] [arguments...]

VERSION:
   0.1.0

DESCRIPTION:
   twtarget (or Twitter Target) is a CLI tool that collects data from twitter API for a given User.

COMMANDS:
     init     Initialize a new Target (User). The command will create a target folder and subfolders with the configuration files and collected data.
     fetch    Fetch data for a given Target (User). 
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

## Initialize a Target

Before you get started you need to add your Twitter API credentials into the configuration file [configuration](conf/default.yml).
```bash
twtarget init -c conf/default.yaml [target screen name]
```

## Get Friends

```bash
twtarget fetch --user --friends [target screen name]
```
## Get Followers
```bash
twtarget fetch --user --followers [target screen name]
```