# Lab

[![CI](https://github.com/Ackerr/lab/workflows/CI/badge.svg)](https://github.com/Ackerr/lab)
[![Go Report Card](https://goreportcard.com/badge/github.com/ackerr/lab)](https://goreportcard.com/report/github.com/ackerr/lab)
[![release](https://img.shields.io/github/v/release/ackerr/lab.svg)](https://github.com/ackerr/lab/releases)

[![ackerr/lab](https://res.cloudinary.com/marcomontalbano/image/upload/v1606925692/video_to_markdown/images/youtube--qqKW9SQqjF0-c05b58ac6eb4c4700831b2b3070cd403.jpg)](https://www.youtube.com/watch?v=qqKW9SQqjF0 "ackerr/lab")

关于GitLab的命令行工具

## 功能

```
lab sync     同步gitlab项目至本地
lab browser  模糊搜索项目名, 回车后，默认浏览器中打开项目地址
lab ws       模糊搜索codespace中的项目，可配合cd，rm使用
lab clone    模糊搜索项目名, 如果设置了codespace, 会将项目clone至codespace，
             否则在当前目录，当然也可以通过--current(-c)，clone至当前路径
lab lint     校验.gitlab-ci.yml文件格式
lab open     快捷在默认浏览器中打开当前所在项目的web地址
lab ci       获取当前项目指定远端分支的ci输出
lab config   快捷打开lab的配置文件
```

> 通过 `lab help` 查看lab更多命令及其参数

## 安装

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

## 配置

推荐使用`lab config`编辑配置，此配置会通过$EDITOR编辑配置文件，如果配置文件不存在，则会使用[默认配置](https://github.com/Ackerr/lab/blob/master/config.toml)新建

其中base_url和token为必填项。token获取方式可参考[这里](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html#creating-a-personal-access-token)

其余配置项含义如下，支持环境变量:

```
 
[gitlab]
# gitlab域名, 比如https://gitlab.com
base_url = "$GITLAB_BASE_URL"

# gitlab access token
token = "$GITLAB_TOKEN"

# 通过 lab sync 同步的项目本地存放放置
# 默认地址为 $HOME/config/.lab/.projects
projects = ""

# 项目存放目录，设置后，lab clone / lab ws 会使用该值作为目标目录
# 默认为空
codespace = ""

# git config中的user.name，设置后会在使用lab clone自动设置
# 默认为空
name = ""

# git config中的user.email，设置后会在使用lab clone自动设置
# 默认为空
email = ""

[main]
# 是否开启fzf。如果开发者个人使用fzf，建议开启
# 默认不开启，使用go-fuzzyfinder
fzf = 0

# 通过lab clone 克隆项目时，支持自定义 git clone 参数
# 例如 `clone_opts="--origin ackerr --branch fix"`, 默认克隆时设置remote为ackerr，分支为fix
# 默认为空
clone_opts = ""
```
