# Lab

[![CI](https://github.com/Ackerr/lab/workflows/CI/badge.svg)](https://github.com/Ackerr/lab)
[![Go Report Card](https://goreportcard.com/badge/github.com/ackerr/lab)](https://goreportcard.com/report/github.com/ackerr/lab)
[![release](https://img.shields.io/github/v/release/ackerr/lab.svg)](https://github.com/ackerr/lab/releases)

[![ackerr/lab](https://res.cloudinary.com/marcomontalbano/image/upload/v1606925692/video_to_markdown/images/youtube--qqKW9SQqjF0-c05b58ac6eb4c4700831b2b3070cd403.jpg)](https://www.youtube.com/watch?v=qqKW9SQqjF0 "ackerr/lab")

A fuzzy finder command line tool for gitlab. [中文文档](./README-CN.md)

## Feature

```
lab sync        Sync gitlab projects
lab browser     Fuzzy find gitlab repo and open it in $BROWSER
lab clone       Fuzzy find gitlab repo and clone it
lab cs          Fuzzy find repo in your codespace
lab lint        Check .gitlab-ci.yml syntax
lab open        Open the current repo remote in $BROWSER
lab config      Use $EDITOR open config file, support custom config path, use --config filepath
lab ci (Beta)   View the pipeline jobs trace log, default view running job
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

Recommend use `lab config` to edit config file (`~/.config/lab/config.toml`), this command will open the config file use `$EDITOR`. If the file don't exist, it will auto generate by [config template](https://github.com/Ackerr/lab/blob/master/config.toml).

> Two variables are required, `base_url` and `token`. The way to get gitlab token, see [this](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html#creating-a-personal-access-token)
