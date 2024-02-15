package utils

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

func CalculateBalance(transactionValue *big.Int, decimals int) string {
	base := big.NewInt(10)
	exponent := big.NewInt(int64(decimals))
	exponentValue := new(big.Int).Exp(base, exponent, nil)

	transactionValueFloat := new(big.Float).SetInt(transactionValue)
	exponentValueFloat := new(big.Float).SetInt(exponentValue)

	resultFloat := new(big.Float).Quo(transactionValueFloat, exponentValueFloat)

	resultString := fmt.Sprintf("%.*f", decimals, resultFloat)

	resultString = strings.TrimRight(resultString, "0")

	resultString = strings.TrimRight(resultString, ".")

	return resultString
}

func HexStringToInt64(hexString string) (int64, error) {
	if hexString == "" {
		return 0, errors.New("hexString can not be empty")
	}
	intValue, err := strconv.ParseInt(hexString, 0, 64)
	if err != nil {
		return 0, err
	}
	return intValue, nil
}
