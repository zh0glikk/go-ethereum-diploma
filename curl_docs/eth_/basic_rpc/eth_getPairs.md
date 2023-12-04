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
> {"jsonrpc":"2.0","id":1,"result":[{"pair":"0x0d4a11d5eeaac28ec3f61d100daf4d40471f1852","pair_version":2,"dex":"univ2"},{"pair":"0x0d4a11d5eeaac28ec3f61d100daf4d40471f1852","pair_version":2,"dex":"sushiv2"},{"pair":"0xc7bbec68d12a0d1830360f8ec58fa599ba1b0e9b","pair_version":3,"dex":"univ3_100"},{"pair":"0x11b815efb8f581194ae79006d24e0d814b7697f6","pair_version":3,"dex":"univ3_500"},{"pair":"0x4e68ccd3e89f51c3074ca5072bbac773960dfa36","pair_version":3,"dex":"univ3_3000"},{"pair":"0xc5af84701f98fa483ece78af83f11b6c38aca71d","pair_version":3,"dex":"univ3_10000"}]}
 
 

