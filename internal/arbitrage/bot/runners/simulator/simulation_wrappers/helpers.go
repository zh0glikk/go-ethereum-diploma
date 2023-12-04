package simulation_wrappers

// func GetOverridesForEstimateSlippage(ctx context.Context, b ethapi.Backend, keyStorage *storage.StateDiffKeysStorage, data models.Data2Simulate) (*ethapi.StateOverride, error) {
// 	key := keyStorage.Get(data.Input)
// 	log.Info("GetOverridesForEstimateSlippage", "getting key from cache")
// 	if key == nil {
// 		log.Info("GetOverridesForEstimateSlippage", "key not found, add new one")
// 		bytes, err := packer.PackerObj.PackBalanceOf(data.Contract)
// 		if err != nil {
// 			return nil, err
// 		}
// 		tx := ethapi.TransactionArgs{To: &data.Input, Data: utils.Ptr(hexutil.Bytes(bytes))}
// 		log.Info(fmt.Sprintf("OverridesInput - %v", data.Input.String()))
// 		acl, _, _, err := ethapi.AccessList(ctx, b, data.BlockNumberOrHash, tx)
// 		if err != nil {
// 			return nil, err
// 		}
// 		key = &acl[0].StorageKeys[len(acl[0].StorageKeys)-1]
// 		keyStorage.Add(data.Input, key)
// 		log.Info("GetOverridesForEstimateSlippage", fmt.Sprintf("key %v added", key.String()))
// 	}
// 	log.Info("GetOverridesForEstimateSlippage", fmt.Sprintf("override key %v", key.String()))
// 	return &ethapi.StateOverride{
// 		data.Input: ethapi.OverrideAccount{
// 			StateDiff: &map[common.Hash]common.Hash{
// 				*key: common.HexToHash("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"),
// 			},
// 		},
// 	}, nil
// }
