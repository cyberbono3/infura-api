# Infura Ethereum API and load testing

This application written in Go and Python exposes REST API to retrieve Ethereum Mainnet transaction and block data via the INFURA JSON-RPC API.

## Description:
This applications contains the following packages:

1.**client-api/data** takes care about payload validation

2.**client-api/handlers** contains handlers for APIs and middleware

3.**config** has config information to access Infura

4.**load** contains load tests code and findings

5.**rpc** contains core functionality of this application and tests

6.**main.go** entry point of this application


## Methods:
1. eth_getTransactionByBlockNumberAndIndex gets transaction by block number and transaction index.

2. eth_getBlockByNumber gets block by block number.

## Installation

### Install using Docker (recommended)

It fetchesthe latest Docker image of this application from my repo and runs application at your local computer inside Docker container.

```
make docker-run
```

### Install locally

```
make install
make run
```

The http server will be started at 8080 port.


### How to use?

1. You can use my Infura-ID to access Infura JSON RPC since my app is for experimental purposes only.

2. Perform POST request.


**eth_getTransactionByBlockNumberAndIndex**


Request:

```json
curl localhost:8080/transaction -X POST -d '{"block":11246512,"index": 10}'
```
Response:

```json
{"Hash":"0x918bc0aaad8f393f84793beb0d3a7a041d5103921fd0a06875fa414de2af7932","Nonce":2860,"BlockHash":"0x9ab73d86b959120525467b91fc3602f3958eba47932c1c63de8c85de7a8727a0","BlockNumber":11246512,"TransactionIndex":10,"From":"0x6893593c695d23f002f9278fa75ae1367ce78d96","To":"0xdac17f958d2ee523a2206206994597c13d831ec7","Value":0,"Gas":84313,"GasPrice":145000000000,"Input":"0xa9059cbb000000000000000000000000825829f75103258c71c85b71d479d5aa867b273900000000000000000000000000000000000000000000000000000000caaefd35"}
```

***eth_getBlockByNumber***

Request:

```json

curl localhost:8080/block -X POST -d '{"block":112465, "show": true}'

```

Response:

```json

{"Number":112465,"Hash":"0xd28b3e21fb530c795513da08e21d3314057060973a3dce152fd080983a033833","ParentHash":"0x811241925469126c3d615b0ff14cb0d9320f619be90e548161b0bbce94a513bd","Nonce":"0x2112a0d8688ed909","Sha3Uncles":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347","LogsBloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","TransactionsRoot":"0x74fdf310083cea3337c57a73cf4ac068c9df58f3dde243a7a29751fd0a3e79d6","StateRoot":"0x0b12afaec57ef065837f196c5c88a68733c677c239686c942f48c9cdd274d31d","Miner":"0xe6a7a1d47ff21b6321162aea7c6cb457d5476bca","Difficulty":4421775782084,"TotalDifficulty":222548530267764322,"ExtraData":"0x476574682f76312e302e312d38326566323666362f6c696e75782f676f312e34","Size":663,"GasLimit":3141592,"GasUsed":21000,"Timestamp":1440007852,"Uncles":[],"Transactions":[{"Hash":"0xfcc490f19eaebd5bcf3b5bc94ea3172414dd7fafffd90d685f003f8026956e32","Nonce":4,"BlockHash":"0xd28b3e21fb530c795513da08e21d3314057060973a3dce152fd080983a033833","BlockNumber":112465,"TransactionIndex":0,"From":"0xbfa71aba804c2b986d969e175acfdac7266cde9c","To":"0xa9f5e7bf6cbf6fec55b30850eb492b47a376d050","Value":106000000000000000000,"Gas":30000,"GasPrice":57107945920,"Input":"0x"}]}

```



## Docker

Dockerfile is used to create Docker image for this application. It copies application files into app folder(inside container) and runs main.go inside app folder.

To create an image 
```
make docker
```

Run the latest image from my repo:
```
make docker-run
```


## How to run tests?

This commands runs all tests of this application.

```
make test
```


# Load testing

Key goals of load testing:

1. Identify breaking points of the system.

2. Find out API throughput.

3. Find out latency.

I have used [Locust](https://locust.io/) to perform load testing. It is powerful opensource tool that allows to define user behavior with Python code and swarm your system with large amount of simultaneous users.

I have performed multiple testing strategies:

1. **Positive testing strategy**. I have used JSON payload with valid parameters. You can view  [code](load/positiveload/positiveload.py) and [findings](load/positiveload/positiveload.md) 

2. **Negative transaction api strategy**. I have used JSON payload with invalid parametersin transaction API. You can view  [code](load/negativetransload/negativetransload.py) and [findings](load/negativetransload/negativetransload.md) 


3. **negative block api strategy**. I have used JSON payload with invalid parameters in block API. You can view [code](load/negativeblockload/negativetransload.py) and [findings](load/negativeblockload/negativeblockload.md).


### Future work on load testing
Locust provides functionality to perform [distributed load testing]( https://docs.locust.io/en/stable/running-locust-distributed.html) using multiple nodes.
I can setup multiple master and worker machines at server machines distributed geographically(China, Japan ,etc) in order to get more comprehensive picture of load testing from multiple regions.















