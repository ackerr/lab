# Lab

[![CI](https://github.com/Ackerr/lab/workflows/CI/badge.svg)](https://github.com/Ackerr/lab)
[![Go Report Card](https://goreportcard.com/badge/github.com/ackerr/lab)](https://goreportcard.com/report/github.com/ackerr/lab)
[![release](https://img.shields.io/github/v/release/ackerr/lab.svg)](https://github.com/ackerr/lab/releases)

[![ackerr/lab](https://res.cloudinary.com/marcomontalbano/image/upload/v1606925692/video_to_markdown/images/youtube--qqKW9SQqjF0-c05b58ac6eb4c4700831b2b3070cd403.jpg)](https://www.youtube.com/watch?v=qqKW9SQqjF0 "ackerr/lab")

A command-line tool for gitlab. [中文](./README-CN.md)

## Feature

```
lab sync        Sync gitlab projects
lab browser     Fuzzy find gitlab repo and open it in $BROWSER
lab clone       Fuzzy find gitlab repo and clone it
lab ws          Fuzzy find repo in your codespace
lab lint        Check .gitlab-ci.yml syntax
lab open        Open the current repo remote in $BROWSER
lab ci          View the pipeline jobs trace log, default view running job
lab config      Use $EDITOR open config file
```

For more information, please use `lab help`.

## Install

### homebrew

```bash
$ brew install ackerr/tap/lab
```

### scoop

```bash
$ scoop bucket add ackerr https://github.com/Ackerr/scoop-bucket
$ scoop install ackerr/lab
```

### go get

```bash
$ go get -u "github.com/ackerr/lab"
```

## Configure

Recommended use `lab config` to edit config file (`~/.config/lab/config.toml`), this command will open the config file use $EDITOR. If config file don't exist, it will use config template auto generate. 

All configuration is as follows: 

> Two variables are required, `base_url` and `token`. The way to get gitlab token, see [this](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html#creating-a-personal-access-token)

```
[gitlab]
# gitlab domain, like https://gitlab.com
base_url = $GITLAB_BASE_URL

# gitlab access token
token = $GITLAB_TOKEN

# If set, lab clone and lab ws will use this path as target path
# default empty
codespace = ""

# If set, lab clone will auto set user.name in repo gitconfig
# default empty
name = ""

# If set, lab clone will auto set user.email in repo gitconfig
# default empty
email = ""

[main]
# If set 1, it will use fzf as fuzzy finder, default use go-fuzzyfinder
# default 0
fzf = 0

# lab clone extra custom git clone config
# example `clone_opts="--origin ackerr --branch fix"`
# default empty
clone_opts = ""
```
