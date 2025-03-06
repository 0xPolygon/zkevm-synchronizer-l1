package etherman_test

import (
	"context"
	"os"
	"testing"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestMainnet(t *testing.T) {
	t.Skip("exploratory test")
	cfg := etherman.Config{
		L1URL:           os.Getenv("L1URL_MAINNET"),
		ForkIDChunkSize: 100,
		Contracts: etherman.ContractConfig{
			ZkEVMAddr:                 common.HexToAddress("0x7ff0b5ff6eb8b789456639ac2a02487c338c1789"),
			RollupManagerAddr:         common.HexToAddress("0x5132A183E9F3CB7C848b0AAC5Ae0c4f0491B7aB2"),
			GlobalExitRootManagerAddr: common.HexToAddress("0x580bda1e7a0cfae92fa7f6c20a3794f169ce3cfb"),
		},
	}

	sut, err := etherman.NewClient(cfg)
	require.NoError(t, err)
	ctx := context.TODO()
	block, err := sut.GetL1BlockUpgradeLxLy(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, block)
	toBlock := uint64(6794610)
	blocks, order, err := sut.GetRollupInfoByBlockRange(ctx, uint64(6794610), &toBlock)
	require.NoError(t, err)
	require.NotNil(t, blocks)
	require.NotNil(t, order)
}

func TestSepolia(t *testing.T) {
	t.Skip("exploratory test")
	cfg := etherman.Config{
		L1URL:           os.Getenv("L1URL_SEPOLIA"),
		ForkIDChunkSize: 100,
		Contracts: etherman.ContractConfig{
			ZkEVMAddr:                 common.HexToAddress("0xA13Ddb14437A8F34897131367ad3ca78416d6bCa"),
			RollupManagerAddr:         common.HexToAddress("0x32d33D5137a7cFFb54c5Bf8371172bcEc5f310ff"),
			GlobalExitRootManagerAddr: common.HexToAddress("0xAd1490c248c5d3CbAE399Fd529b79B42984277DF"),
		},
	}

	sut, err := etherman.NewClient(cfg)
	require.NoError(t, err)
	ctx := context.TODO()
	block, err := sut.GetL1BlockUpgradeLxLy(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, block)
	// This block on Sepolia contains a EIP7702 transaction on an event L1InfoTree
	toBlock := uint64(7836383)
	blocks, order, err := sut.GetRollupInfoByBlockRange(ctx, uint64(7836383), &toBlock)
	require.NoError(t, err)
	require.NotNil(t, blocks)
	require.NotNil(t, order)
}
