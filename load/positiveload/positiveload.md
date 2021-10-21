# positive load testing 

Positive load testing of my REST API


## Test1  200 users / 200 ramp up

Execute command from console

```
locust -f positiveload.py -u 200 -r 200 --run-time 1m --headless --csv=test1 --host=http://18.222.141.251:8080

```

simulated users 200 users
ramp up rate 200
running time 1 m

### Conclusion:

1. Throughput 3983 requests per minute.

2. Latency (/block) average latency 5237 ms and Latency (/transaction) average 594 ms are quite high (should be improved). The possible reason might be that I perform Unmarshal costly operation [here](https://github.com/cyberbono3/infura-challenge/blob/main/rpc/client.go#L62) and [here](https://github.com/cyberbono3/infura-challenge/blob/main/rpc/client.go#L144) and return  transaction and block, respectively, to a user (should be improved).



### Possible improvements:

1. POST /transaction
To return
```
{"jsonrpc": "2.0","id": 1,
  "result": transactionObject}
```
instead of pure transactionObject to a user

2. POST /block
To return
```
{"jsonrpc": "2.0","id": 1,
  "result": blockObject }
```
instead of pure blockObject to a user


### Report

```
Name                                                          # reqs      # fails  |     Avg     Min     Max  Median  |   req/s failures/s
--------------------------------------------------------------------------------------------------------------------------------------------
 POST /block                                                     1898    13(0.68%)  |    5237     397   25120    4300  |   31.71    0.22
 POST /transaction                                               2085     0(0.00%)  |     594     131    4369     510  |   34.84    0.00
--------------------------------------------------------------------------------------------------------------------------------------------
 Aggregated                                                      3983    13(0.33%)  |    2806     131   25120    1100  |   66.55    0.22

Response time percentiles (approximated)
 Type     Name                                                              50%    66%    75%    80%    90%    95%    98%    99%  99.9% 99.99%   100% # reqs
--------|------------------------------------------------------------|---------|------|------|------|------|------|------|------|------|------|------|------|
 POST     /block                                                           4300   5600   6700   7400  10000  12000  14000  16000  23000  25000  25000   1898
 POST     /transaction                                                      510    580    620    660    860   1300   1700   1900   2900   4400   4400   2085
--------|------------------------------------------------------------|---------|------|------|------|------|------|------|------|------|------|------|------|
 None     Aggregated                                                       1100   3100   4100   4800   7300  10000  13000  14000  18000  25000  25000   3983

Error report
 # occurrences      Error                                                                                               
--------------------------------------------------------------------------------------------------------------------------------------------
 1                  POST /block: ChunkedEncodingError(ProtocolError('Connection broken: IncompleteRead(3597 bytes read, 6643 more expected)', IncompleteRead(3597 bytes read, 6643 more expected)))
 1                  POST /block: ChunkedEncodingError(ProtocolError('Connection broken: IncompleteRead(141 bytes read, 10099 more expected)', IncompleteRead(141 bytes read, 10099 more expected)))
 2                  POST /block: ChunkedEncodingError(ProtocolError('Connection broken: IncompleteRead(2933 bytes read, 59 more expected)', IncompleteRead(2933 bytes read, 59 more expected)))
 1                  POST /block: ChunkedEncodingError(ProtocolError('Connection broken: IncompleteRead(1733 bytes read, 1259 more expected)', IncompleteRead(1733 bytes read, 1259 more expected)))
 1                  POST /block: ChunkedEncodingError(ProtocolError('Connection broken: IncompleteRead(7629 bytes read, 2611 more expected)', IncompleteRead(7629 bytes read, 2611 more expected)))
 1                  POST /block: ChunkedEncodingError(ProtocolError('Connection broken: IncompleteRead(1485 bytes read, 1507 more expected)', IncompleteRead(1485 bytes read, 1507 more expected)))
 1                  POST /block: ChunkedEncodingError(ProtocolError('Connection broken: IncompleteRead(4693 bytes read, 5547 more expected)', IncompleteRead(4693 bytes read, 5547 more expected)))
 1                  POST /block: ChunkedEncodingError(ProtocolError('Connection broken: IncompleteRead(3037 bytes read, 7203 more expected)', IncompleteRead(3037 bytes read, 7203 more expected)))
 1                  POST /block: ChunkedEncodingError(ProtocolError('Connection broken: IncompleteRead(9141 bytes read, 1099 more expected)', IncompleteRead(9141 bytes read, 1099 more expected)))
 1                  POST /block: ChunkedEncodingError(ProtocolError('Connection broken: IncompleteRead(6493 bytes read, 3747 more expected)', IncompleteRead(6493 bytes read, 3747 more expected)))
 1                  POST /block: ChunkedEncodingError(ProtocolError('Connection broken: IncompleteRead(7485 bytes read, 2755 more expected)', IncompleteRead(7485 bytes read, 2755 more expected)))
 1                  POST /block: ChunkedEncodingError(ProtocolError('Connection broken: IncompleteRead(4589 bytes read, 5651 more expected)', IncompleteRead(4589 bytes read, 5651 more expected)))
----------------------------------------------------------------------------------------------------------------------------------
```














