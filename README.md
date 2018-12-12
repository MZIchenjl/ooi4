# OOI4

[![CircleCI](https://circleci.com/gh/MZIchenjl/ooi4/tree/master.svg?style=svg)](https://circleci.com/gh/MZIchenjl/ooi4/tree/master)

Online Objects Integration version 4 based on pure go

Inspired from [ooi3](https://github.com/acgx/ooi3)

## Feature

* Not depend on nginx
* Support https protocol

## Build

Require go >= 1.10.0 or latest (development)

Clone the repository

```shell
git clone https://github.com/MZIchenjl/ooi4.git
```

Build the project

```shell
go build -o ooi4
```

## Run

Place the executable file with the `static/` folder and the config file `app.toml`

Usage

```shell
Usage of ooi4:
  -config string
    	Set the config file(toml) (default "app.toml")
```

## License

[BSD-3-Clause](https://github.com/MZIchenjl/ooi4/blob/master/LICENSE)
