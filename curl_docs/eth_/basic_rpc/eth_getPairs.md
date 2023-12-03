## eth_callMany for counter contract

### Curl :

```shell
curl http://localhost:8543/ \
-X POST \
-H "Content-Type: application/json" \
-d '{
   "jsonrpc":"2.0",
   "method":"eth_getPairs",
   "params":[
      {
          "token0": "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
          "token1": "0xdac17f958d2ee523a2206206994597c13d831ec7"
      },
      {},
      "latest",
      {}
   ],
   "id":1
}'
```

### Response :
> {"jsonrpc":"2.0","id":1,"result":[{"value":"0x0000000000000000000000000000000000000000000000000000000000000001"},{"value":"0x0000000000000000000000000000000000000000000000000000000000000002"},{"value":"0x0000000000000000000000000000000000000000000000000000000000000003"}]}
 

