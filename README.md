# [WIP] branehub

[![pipeline status](https://gitlab.com/braneproject/branehub/badges/master/pipeline.svg)](https://gitlab.com/braneproject/BitcoinPrimitives.jl/commits/master)

## About

RESTful API for Brane Project.

BraneHub actually retrieves information directly from cryptocurrency exchanges
(i.e. Kraken, Bitstamp and BL3P) to compute data. This is however a proof of
concept, in the future BraneHub will not compute any data but only forward
data computed by Brane's tools ([Bitcoin.jl](https://gitlab.com/braneproject/branehub/Bitcoin.jl)
to start with).

## Reference

### General Information

#### Base URL

HTTP

`https://api.brane.cc`

#### Variable Path
##### Definition

`/<version>/<callname>/<market>`

##### Description

Version of API (is currently: 0)

`<version> = 0`

Name of call (for example: "ticker")

`<callname> = $callname`

Market that the call will be applied to.

`<market> = "btceur"`

### Ticker

Passing any GET parameters, will result in your request being rejected.

Request | |
------------
GET | https://api.brane.cc/0/ticker/{currency_pair}
 | Supported values for currency_pair: any six letter combination supported by one or several of connected exchanges.

Response JSON |
----------------
market | Actual market
vwap | Volume weighted average price of exchanges' last reported price

#### Example

Input:
`URL: https://api.brane.cc/0/ticker/btcusd`

Result:
`{"market":"BTCUSD","vwap":10120.259}`

## Buy me a cup of coffee

[Donate Bitcoin](bitcoin:34nvxratCQcQgtbwxMJfkmmxwrxtShTn67)
