// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.8.2 <0.9.0;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

import "@uniswap/v2-core/contracts/interfaces/IUniswapV2Pair.sol";
import "@uniswap/v3-core/contracts/interfaces/IUniswapV3Pool.sol";

contract TokenSwapContract is Ownable {

    bytes4 private constant SELECTOR = bytes4(keccak256(bytes('transfer(address,uint256)')));

    function add(uint x, uint y) internal pure returns (uint z) {
        require((z = x + y) >= x, 'ds-math-add-overflow');
    }

    function sub(uint x, uint y) internal pure returns (uint z) {
        require((z = x - y) <= x, 'ds-math-sub-underflow');
    }

    function mul(uint x, uint y) internal pure returns (uint z) {
        require(y == 0 || (z = x * y) / y == x, 'ds-math-mul-overflow');
    }

    // given an input amount of an asset and pair reserves, returns the maximum output amount of the other asset
    function getAmountOut(uint amountIn, uint reserveIn, uint reserveOut) internal pure returns (uint amountOut) {
        require(amountIn > 0, 'UniswapV2Library: INSUFFICIENT_INPUT_AMOUNT');
        require(reserveIn > 0 && reserveOut > 0, 'UniswapV2Library: INSUFFICIENT_LIQUIDITY');
        uint amountInWithFee = mul(amountIn,997);
        uint numerator = mul(amountInWithFee,reserveOut);
        uint denominator =add( mul(reserveIn,1000),amountInWithFee);
        amountOut = numerator / denominator;
    }

    function swapV2(
        address pairAddr,
        address inputAddr,
        address outputAddr,
        uint256 inputAmount
    ) external returns (uint256 outputAmount) {
        IERC20 input = IERC20(inputAddr);
        IUniswapV2Pair pair = IUniswapV2Pair(pairAddr);

        (uint reserve0, uint reserve1,) = pair.getReserves();

        require(input.balanceOf(address(this))>inputAmount, "low input balance");
        _safeTransfer(inputAddr, pairAddr, inputAmount);

        outputAmount = 0;
        if (pair.token0() == inputAddr) {
            outputAmount = getAmountOut(inputAmount, reserve0, reserve1);

            pair.swap(0, outputAmount, address(this), new bytes(0));
        } else {
            outputAmount = getAmountOut(inputAmount, reserve1, reserve0);
            pair.swap(outputAmount, 0, address(this), new bytes(0));
        }

        return outputAmount;
    }

    function fixBalances(address input, address output) private view returns (uint256, uint256) {
        IERC20 tokenInput = IERC20(input);
        IERC20 tokenOutput = IERC20(output);

        uint256 inputBalance = tokenInput.balanceOf(address(this));
        uint256 outputBalance = tokenOutput.balanceOf(address(this));

        return (inputBalance, outputBalance);
    }

    function _safeTransfer(address token, address to, uint value) private {
        (bool success, bytes memory data) = token.call(abi.encodeWithSelector(SELECTOR, to, value));
        require(success && (data.length == 0 || abi.decode(data, (bool))), 'UniswapV2: TRANSFER_FAILED');
    }


    struct SwapCallbackData {
        bytes path;
        address payer;
    }
    struct Slot0 {
        // the current price
        uint160 sqrtPriceX96;
        // the current tick
        int24 tick;
        // the most-recently updated index of the observations array
        uint16 observationIndex;
        // the current maximum number of observations that are being stored
        uint16 observationCardinality;
        // the next maximum number of observations to store, triggered in observations.write
        uint16 observationCardinalityNext;
        // the current protocol fee as a percentage of the swap fee taken on withdrawal
        // represented as an integer denominator (1/x)%
        uint8 feeProtocol;
        // whether the pool is locked
        bool unlocked;
    }

    uint256 private constant ADDR_SIZE = 20;
    uint256 private constant FEE_SIZE = 3;
    uint256 private constant NEXT_OFFSET = ADDR_SIZE + FEE_SIZE;

    uint160 internal constant MIN_SQRT_RATIO = 4295128739;
    uint160 internal constant MAX_SQRT_RATIO = 1461446703485210103287273052203988822378723970342;
    bytes32 internal constant POOL_INIT_CODE_HASH = 0xe34f199b19b2b4f47f68442619d555527d244f78a3297ea89325f843f87b8b54;
    address internal constant UNISWAP_V3_FACTORY = 0x1F98431c8aD98523631AE4a59f267346ea31F984;

    function transfer(address token, address to, uint256 amount) internal  returns(uint256 res) {
        assembly {
            let emptyPtr := mload(0x40)
            mstore(emptyPtr, 0xa9059cbb00000000000000000000000000000000000000000000000000000000)
            mstore(add(emptyPtr, 0x4), to)
            mstore(add(emptyPtr, 0x24), amount)
            pop(call(gas(), token, 0, emptyPtr, 0x44, 0, 0))
        }
    }
    function balanceOf(address token, address acc) internal  returns(uint256 res) {
        assembly {
            let emptyPtr := mload(0x40)
            mstore(emptyPtr, 0x70a0823100000000000000000000000000000000000000000000000000000000)
            mstore(add(emptyPtr, 0x4), acc)
            pop(call(gas(), token, 0, emptyPtr, 0x24, emptyPtr, 0x20))
            res := mload(emptyPtr)
        }
    }
    function optimize(uint256 a) internal returns(uint a2) {
        assembly {
            a2 := a
            for {let m := 256 } gt( div(mul(mul(div(a2,m),m),100000),a) , 99998) { m := mul(m,256)  } {
                a2 := mul(div(a2, m),m)
            }
        }
    }
    function toAddress(bytes memory _bytes, uint256 _start) internal pure returns (address) {
        require(_start + 20 >= _start, 'toAddress_overflow');
        require(_bytes.length >= _start + 20, 'toAddress_outOfBounds');
        address tempAddress;

        assembly {
            tempAddress := div(mload(add(add(_bytes, 0x20), _start)), 0x1000000000000000000000000)
        }

        return tempAddress;
    }
    function toUint24(bytes memory _bytes, uint256 _start) internal pure returns (uint24) {
        require(_start + 3 >= _start, 'toUint24_overflow');
        require(_bytes.length >= _start + 3, 'toUint24_outOfBounds');
        uint24 tempUint;

        assembly {
            tempUint := mload(add(add(_bytes, 0x3), _start))
        }

        return tempUint;
    }
    function decodeFirstPool(bytes memory path)
    internal
    pure
    returns (
        address tokenA,
        address tokenB,
        uint24 fee
    )
    {
        tokenA = toAddress(path,0);
        fee = toUint24(path,ADDR_SIZE);
        tokenB = toAddress(path,NEXT_OFFSET);
    }

    function computeAddress(address factory, address token0, address token1, uint24 fee) internal pure returns (address pool) {
        (token0, token1) = (token0 < token1) ? (token0, token1) : (token1,token0);
        pool = address(
            uint160(uint256(
                keccak256(
                    abi.encodePacked(
                        hex'ff',
                        factory,
                        keccak256(abi.encode(token0, token1, fee)),
                        POOL_INIT_CODE_HASH
                    )
                )
            ))
        );
    }

    function verifyCallback(address factory, address token0, address token1, uint24 fee)
    internal
    view
    {
        address pool = computeAddress(factory, token0, token1, fee);
        require(msg.sender == address(pool));
    }


    function uniswapV3SwapCallback(
        int256 amount0Delta,
        int256 amount1Delta,
        bytes calldata _data
    ) external {
        SwapCallbackData memory data = abi.decode(_data, (SwapCallbackData));
        (address tokenIn, address tokenOut, uint24 fee) = decodeFirstPool(data.path);

        (bool isExactInput, uint256 amountToPay) =
        amount0Delta > 0
        ? (tokenIn < tokenOut, uint256(amount0Delta))
        : (tokenOut < tokenIn, uint256(amount1Delta));
        if (isExactInput) {
            transfer(tokenIn, msg.sender, amountToPay);
        } else {
            tokenIn = tokenOut; // swap in/out because exact output swaps are reversed
            transfer(tokenIn, msg.sender, amountToPay);
            assembly {
                sstore(1,amountToPay)
            }
        }
    }

    function swapV3(uint256 inputAmount,
        address input,
        address output,
        address pair,
        uint160 sqrtPriceLimitX96,
        address receiver) public{
        if(inputAmount == 0) {
            assembly {
                inputAmount := sload(0)
                if gt(inputAmount,0) {inputAmount := sub(inputAmount,1)}
            }
            inputAmount = (sqrtPriceLimitX96 == 0) ? inputAmount : optimize(inputAmount);
        }
        {
            uint256 tokensBefore = balanceOf(output,receiver);
            assembly{ sstore(0,tokensBefore) }
        }
        bool zeroForOne = input < output;
        SwapCallbackData memory data = SwapCallbackData({path: abi.encodePacked(input,
            uint24(0), output), payer: address(this)});

        (int256 amount0, int256 amount1) = IUniswapV3PoolActions(pair).swap(
            receiver,
            zeroForOne,
            int256(inputAmount),
        // uint160(sqrtLimit),
            sqrtPriceLimitX96 == 0
            ? (zeroForOne ? MIN_SQRT_RATIO + 1 : MAX_SQRT_RATIO - 1)
            : sqrtPriceLimitX96,
            abi.encode(data)
        );

        (uint160 limit , , , , , , ) = IUniswapV3PoolState(pair).slot0();

        uint256 received = balanceOf(output,receiver);
        assembly {
            received := sub(received, sload(0))
            if gt(received,0) { received := sub(received,1)}
            sstore(0,received)
        }
        uint256 outputAmount = uint256(-(zeroForOne ? amount1 : amount0));

        if(sqrtPriceLimitX96 != 0) {
            outputAmount = optimize(outputAmount);
            received = optimize(received);
        }

        assembly {
            let emptyPtr := mload(0x40)
            mstore( emptyPtr, outputAmount)
            mstore( add(emptyPtr,0x20), received)
            mstore( add(emptyPtr,0x40),limit)
            return (emptyPtr,0x60)
        }
    }

}
