# middleware

[![Build Status](https://travis-ci.org/Akagi201/middleware.svg)](https://travis-ci.org/Akagi201/middleware)
[![Coverage Status](https://coveralls.io/repos/github/Akagi201/middleware/badge.svg?branch=master)](https://coveralls.io/github/Akagi201/middleware?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/Akagi201/middleware)](https://goreportcard.com/report/github.com/Akagi201/middleware)
[![GoDoc](https://godoc.org/github.com/Akagi201/middleware?status.svg)](https://godoc.org/github.com/Akagi201/middleware)

A collection of net/http middlewares in Go.

It is designed to be fully compatible with http standard library, easy to customize and reuse.

## Middlewares
* [pressly/lg](https://github.com/pressly/lg): logrus.
* [vulcand/oxy](https://github.com/vulcand/oxy): proxy utils.
* [rs/cors](https://github.com/rs/cors): cors.
* [didip/tollbooth](https://github.com/didip/tollbooth): rate limit
* [sebest/xff](https://github.com/sebest/xff): handle X-Forwarded-For
* [rs/formjson](https://github.com/rs/formjson): transparently manage posted JSON
* [daaku/go.httpgzip](https://github.com/daaku/go.httpgzip): http gzip

## Middlewares management / framework
* light: <https://github.com/Akagi201/light>
* alice: <https://github.com/justinas/alice>
* catena: <https://github.com/codemodus/catena>
* chain: <https://github.com/codemodus/chain>
* gohttp: <https://github.com/gohttp/app>
* go-swagger: <https://github.com/go-swagger/go-swagger>
* negroni: <https://github.com/urfave/negroni>
* chi: <https://github.com/pressly/chi>
* [interpose](https://github.com/carbocation/interpose): a minimalist net/http middleware framework
