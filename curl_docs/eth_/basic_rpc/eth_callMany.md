## eth_callMany for counter contract

### Curl :

```shell
curl http://localhost:8543/ \
-X POST \
-H "Content-Type: application/json" \
-d '{
   "jsonrpc":"2.0",
   "method":"eth_callMany",
   "params":[
      [
            {
                  "to": "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
                  "data": "0xd09de08a"
            },
            {
                  "to": "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
                  "data": "0xd09de08a"
            },
            {
                  "to": "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
                  "data": "0xd09de08a"
            }
      ],
      "latest",
      {
        "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48": {
            "code": "0x608060405234801561001057600080fd5b50600436106100365760003560e01c80633fa4f2451461003b578063d09de08a14610059575b600080fd5b610043610077565b60405161005091906100ee565b60405180910390f35b61006161007d565b60405161006e91906100ee565b60405180910390f35b60005481565b6000600160005461008e9190610138565b6000819055507f20d8a6f5a693f9d1d627a598e8820f7a55ee74c183aa8f1a30e8d4e8dd9a8d846000546040516100c591906100ee565b60405180910390a1600054905090565b6000819050919050565b6100e8816100d5565b82525050565b600060208201905061010360008301846100df565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000610143826100d5565b915061014e836100d5565b925082820190508082111561016657610165610109565b5b9291505056fea264697066735822122067549d2a8a8993aed4ccac27769940697103360ff759392a48bcd15ac8cf282864736f6c63430008120033"
        }
      },
      {

      }
   ],
   "id":1
}'
```

### Response :
> {"jsonrpc":"2.0","id":1,"result":[{"value":"0x0000000000000000000000000000000000000000000000000000000000000001"},{"value":"0x0000000000000000000000000000000000000000000000000000000000000002"},{"value":"0x0000000000000000000000000000000000000000000000000000000000000003"}]}
 
