package signer

import (
	"fmt"
	"math/big"
)

func getStringFromInterface(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case fmt.Stringer:
		return val.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}

func getBigIntFromInterface(v interface{}) *big.Int {
	switch val := v.(type) {
	case *big.Int:
		return val
	case int64:
		return big.NewInt(val)
	case float64:
		return big.NewInt(int64(val))
	case string:
		bi, _ := new(big.Int).SetString(val, 0)
		return bi
	default:
		return nil
	}
}
