package usecases

import "github.com/ethereum/go-ethereum/common"

// erc20 tokens
var wethAddr common.Address
var usdtAddress common.Address
var usdcAddress common.Address
var wbtcAddress common.Address

// factories
var uniSwapV2FactoryAddr common.Address
var uniSwapV3FactoryAddr common.Address
var sushiSwapV2FactoryAddr common.Address

var baseCurrencies []common.Address
var baseCurrenciesMap = map[common.Address]struct{}{}

func init() {
	wethAddr = common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")
	usdtAddress = common.HexToAddress("0xdac17f958d2ee523a2206206994597c13d831ec7")
	usdcAddress = common.HexToAddress("0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48")
	wbtcAddress = common.HexToAddress("0x2260FAC5E5542a773Aa44fBCfeDf7C193bc2C599")

	uniSwapV2FactoryAddr = common.HexToAddress("0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f")
	uniSwapV3FactoryAddr = common.HexToAddress("0x1F98431c8aD98523631AE4a59f267346ea31F984")
	sushiSwapV2FactoryAddr = common.HexToAddress("0xC0AEe478e3658e2610c5F7A4A2E1777cE9e4f2Ac")

	baseCurrencies = []common.Address{
		wethAddr,
		wbtcAddress,
		usdtAddress,
		usdcAddress,
	}

	for _, curr := range baseCurrencies {
		baseCurrenciesMap[curr] = struct{}{}
	}
}

func BaseCurrenciesList() []common.Address {
	return baseCurrencies
}

func IsBaseCurrency(address common.Address) bool {
	_, ok := baseCurrenciesMap[address]
	return ok
}

func IsWeth(addr common.Address) bool {
	return addr == wethAddr
}
