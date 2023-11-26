## trace_trackSwap

```shell
curl http://localhost:8543/ \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{
        "jsonrpc":"2.0",
        "method":"trace_trackSwap",
        "params":[       
	        {
		        "transactions": [
		        {
                      "from":"0x1a3eb84BE3DdfAB270B04d423F3C493418caC938",
                      "to":"0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D",
	                  "data":"0x38ed17390000000000000000000000000000000000000000000000000c4770e98144000000000000000000000000000000000000000000000000000000180d31282e223900000000000000000000000000000000000000000000000000000000000000a00000000000000000000000001a3eb84be3ddfab270b04d423f3c493418cac93800000000000000000000000000000000000000000000000000000000653b57620000000000000000000000000000000000000000000000000000000000000002000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc20000000000000000000000005acd02940d2e56d9402b8d224e56bd800c544466",
                      "value": "0x0"
		        }
		    ]
		    },
            "0x1195E28",
            {
            
            }
        ],
        "id":1
    }'
```

### Response :
> {"jsonrpc":"2.0","id":1,"result":{"swaps":[[{"type":"v2","pair":"0xf6230de716e50bdbf4c5ea2b5e7d006bc6af603c","input":"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2","output":"0x5acd02940d2e56d9402b8d224e56bd800c544466","inputAmount":884800000000000000}]],"duration":21}}
 
