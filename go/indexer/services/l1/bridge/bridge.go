package bridge

import (
	"context"
	"math/big"

	"github.com/ethereum-optimism/optimism/go/indexer/bindings/l1bridge"
	"github.com/ethereum-optimism/optimism/go/indexer/db"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type Bridge interface {
	Address() common.Address
	GetDepositsByBlockRange(uint64, uint64) (map[common.Hash][]db.Deposit, error)
}

var DEFAULT_L2_CONTRACT_ADDRESSES = map[string]string{
	"OVM_L2ToL1MessagePasser": "0x4200000000000000000000000000000000000000",
	"OVM_DeployerWhitelist":   "0x4200000000000000000000000000000000000002",
	"L2CrossDomainMessenger":  "0x4200000000000000000000000000000000000007",
	"OVM_GasPriceOracle":      "0x420000000000000000000000000000000000000F",
	"L2StandardBridge":        "0x4200000000000000000000000000000000000010",
	"OVM_SequencerFeeVault":   "0x4200000000000000000000000000000000000011",
	"L2StandardTokenFactory":  "0x4200000000000000000000000000000000000012",
	"OVM_L1BlockNumber":       "0x4200000000000000000000000000000000000013",

	// We're temporarily disabling OVM_ETH because the jury is still out on whether or not ETH as an
	// ERC20 is desirable.
	"OVM_ETH": "0xDeadDeAddeAddEAddeadDEaDDEAdDeaDDeAD0000",

	// We're also putting WETH9 at the old OVM_ETH address.
	"WETH9": "0x4200000000000000000000000000000000000006",
}

var CONTRACT_ADDRESSES = map[uint64]map[string]map[string]string{
	// Mainnet
	1: {
		"l1": {
			"AddressManager":            "0xdE1FCfB0851916CA5101820A69b13a4E276bd81F",
			"L1CrossDomainMessenger":    "0x25ace71c97B33Cc4729CF772ae268934F7ab5fA1",
			"L1StandardBridge":          "0x99C9fc46f92E8a1c0deC1b1747d010903E884bE1",
			"StateCommitmentChain":      "0xBe5dAb4A2e9cd0F27300dB4aB94BeE3A233AEB19",
			"CanonicalTransactionChain": "0x5E4e65926BA27467555EB562121fac00D24E9dD2",
			"BondManager":               "0xcd626E1328b41fCF24737F137BcD4CE0c32bc8d1",
		},
		"l2": DEFAULT_L2_CONTRACT_ADDRESSES,
	},
	// Kovan
	42: {
		"l1": {
			"AddressManager":            "0x100Dd3b414Df5BbA2B542864fF94aF8024aFdf3a",
			"L1CrossDomainMessenger":    "0x4361d0F75A0186C05f971c566dC6bEa5957483fD",
			"L1StandardBridge":          "0x22F24361D548e5FaAfb36d1437839f080363982B",
			"StateCommitmentChain":      "0xD7754711773489F31A0602635f3F167826ce53C5",
			"CanonicalTransactionChain": "0xf7B88A133202d41Fe5E2Ab22e6309a1A4D50AF74",
			"BondManager":               "0xc5a603d273E28185c18Ba4d26A0024B2d2F42740",
		},
		"l2": DEFAULT_L2_CONTRACT_ADDRESSES,
	},
	// Goerli
	5: {
		"l1": {
			"AddressManager":            "0x2F7E3cAC91b5148d336BbffB224B4dC79F09f01D",
			"L1CrossDomainMessenger":    "0xEcC89b9EDD804850C4F343A278Be902be11AaF42",
			"L1StandardBridge":          "0x73298186A143a54c20ae98EEE5a025bD5979De02",
			"StateCommitmentChain":      "0x1afcA918eff169eE20fF8AB6Be75f3E872eE1C1A",
			"CanonicalTransactionChain": "0x2ebA8c4EfDB39A8Cd8f9eD65c50ec079f7CEBD81",
			"BondManager":               "0xE5AE60bD6F8DEe4D0c2BC9268e23B92F1cacC58F",
		},
		"l2": DEFAULT_L2_CONTRACT_ADDRESSES,
	},
	// Hardhat local
	31337: {
		"l1": {
			"AddressManager":            "0x5FbDB2315678afecb367f032d93F642f64180aa3",
			"L1CrossDomainMessenger":    "0x8A791620dd6260079BF849Dc5567aDC3F2FdC318",
			"L1StandardBridge":          "0x610178dA211FEF7D417bC0e6FeD39F05609AD788",
			"StateCommitmentChain":      "0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9",
			"CanonicalTransactionChain": "0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9",
			"BondManager":               "0x5FC8d32690cc91D4c39d9d3abcBD16989F875707",
		},
		"l2": DEFAULT_L2_CONTRACT_ADDRESSES,
	},
}

func BridgesByChainID(chainID *big.Int, client bind.ContractFilterer, ctx context.Context) ([]Bridge, error) {
	l1StandardBridgeAddress := common.HexToAddress(CONTRACT_ADDRESSES[chainID.Uint64()]["l1"]["L1StandardBridge"])
	standardBridgeContract, err := l1bridge.NewL1StandardBridgeFilterer(l1StandardBridgeAddress, client)
	if err != nil {
		return nil, err
	}

	switch chainID.Uint64() {
	case 42:
		return []Bridge{
			&StandardBridge{
				ctx:      ctx,
				address:  l1StandardBridgeAddress,
				client:   client,
				filterer: standardBridgeContract,
			},
			&EthBridge{
				ctx:      ctx,
				address:  l1StandardBridgeAddress,
				client:   client,
				filterer: standardBridgeContract,
			},
		}, nil
	}
	return nil, nil
}
