# Lab

[![CI](https://github.com/Ackerr/lab/workflows/CI/badge.svg)](https://github.com/Ackerr/lab)
[![Go Report Card](https://goreportcard.com/badge/github.com/ackerr/lab)](https://goreportcard.com/report/github.com/ackerr/lab)
> [![codecov](https://codecov.io/gh/Ackerr/lab/branch/master/graph/badge.svg)](https://codecov.io/gh/Ackerr/lab)
[![release](https://img.shields.io/github/v/release/ackerr/lab.svg)](https://github.com/ackerr/lab/releases)


A cli tool with gitlab

## Feature

1. Fuzzy find your gitlab project, open or clone it, use `lab browser` or `lab clone`.
2. Open the current gitlab repository in web browser, use `lab open [remote]`

## Env

```bash
$ export GITLAB_BASE_URL=<GITLAB_BASE_URL>
$ export GITLAB_TOKEN=<GITLAB_TOKEN>
```

## Installation

### homebrew

```bash
$ brew install ackerr/tap/lab
```

### go get

```bash
$ go get -u "github.com/ackerr/lab"
```

## Usage

First, run `lab sync` to sync all gitlab projects, then you can fuzzy find your project and open, use `lab browser`

> Use`lab -h` to see more commands.
