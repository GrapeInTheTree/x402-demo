package main

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("https://sepolia.base.org")
	if err != nil {
		fmt.Println("RPC error:", err)
		return
	}
	ctx := context.Background()

	usdc := common.HexToAddress("0x036CbD53842c5426634e7929541eC2318f3dCF7e")
	erc20ABI, _ := abi.JSON(strings.NewReader(`[{"inputs":[{"name":"account","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"stateMutability":"view","type":"function"}]`))

	ethDiv := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
	usdcDiv := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(6), nil))

	fmt.Println("==========================================")
	fmt.Println("  Base Sepolia — After Payment")
	fmt.Println("==========================================")

	wallets := []struct {
		name string
		addr common.Address
	}{
		{"Facilitator", common.HexToAddress("0x23fbdE5A14dFB508502f5A2622f66c0D3B0ab37A")},
		{"Client", common.HexToAddress("0x47322Ca28a85B12a7EA64a251Cd8b9Ea1fac037b")},
		{"PAY_TO (you)", common.HexToAddress("0xDBCbC75772954F82d436700cDC4B7c8F434e07F5")},
	}

	for _, w := range wallets {
		ethBal, _ := client.BalanceAt(ctx, w.addr, nil)
		data, _ := erc20ABI.Pack("balanceOf", w.addr)
		result, _ := client.CallContract(ctx, ethereum.CallMsg{To: &usdc, Data: data}, nil)
		out, _ := erc20ABI.Unpack("balanceOf", result)

		fmt.Printf("\n%-14s %s\n", w.name+":", w.addr.Hex())
		fmt.Printf("  ETH:  %s\n", new(big.Float).Quo(new(big.Float).SetInt(ethBal), ethDiv).Text('f', 6))
		if len(out) > 0 {
			fmt.Printf("  USDC: %s\n", new(big.Float).Quo(new(big.Float).SetInt(out[0].(*big.Int)), usdcDiv).Text('f', 6))
		}
	}
}
