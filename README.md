![Logo](http://svg.wiersma.co.za/hamba/project?title=cmd&tag=Go%20cmd%20helper)

[![Go Report Card](https://goreportcard.com/badge/github.com/hamba/cmd)](https://goreportcard.com/report/github.com/hamba/cmd)
[![Build Status](https://travis-ci.com/hamba/cmd.svg?branch=master)](https://travis-ci.com/hamba/cmd)
[![Coverage Status](https://coveralls.io/repos/github/hamba/cmd/badge.svg?branch=master)](https://coveralls.io/github/hamba/cmd?branch=master)
[![GoDoc](https://godoc.org/github.com/hamba/cmd?status.svg)](https://godoc.org/github.com/hamba/cmd)
[![GitHub release](https://img.shields.io/github/release/hamba/cmd.svg)](https://github.com/hamba/cmd/releases)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/hamba/cmd/master/LICENSE)

Go cmd helper. 

This provides helpers on top of `github.com/urfave/cli`.

## Overview

Install with:

```shell
go get github.com/hamba/cmd
```

## Example

```go
func yourAction(c *cli.Context) error {
    ctx, err := cmd.NewContext(c)
    if err != nil {
        return err
    }

    // Run your application here...

    <-cmd.WaitForSignals()
}
```
