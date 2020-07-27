# Lab

A cli tool with gitlab

## Feature

  - command line fuzzy find gitlab projects, and open

## Env

```bash
export GITLAB_BASE_URL=<GITLAB_BASE_URL>
export GITLAB_TOKEN=<GITLAB_TOKEN>
```

## Installation

```bash
go get -u "github.com/ackerr/lab"
```

## Usage

First, run `lab sync` to sync all gitlab projects, then you can fuzzy find your project and open, use `lab browser`

> Use`lab -h` to see more commands.
