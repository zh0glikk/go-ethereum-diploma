package unpacker

const (
	transferSign     = "0xa9059cbb"
	transferFromSign = "0x23b872dd"
	swapSigV2Sign    = "0x022c0d9f"
	swapSigV3Sign    = "0x128acb08"
)

type TxType int

const (
	Front TxType = iota
	Back
)
