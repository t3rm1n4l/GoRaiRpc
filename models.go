package gorairpc

import "math/big"

type Action struct {
	Name string `json:"action"`
}

type AvailableSupply struct {
	AvailableRaw string `json:"available"`
}

func (a *AvailableSupply) ConvertUnitTo(unit string) string {
	return toUnit(a.AvailableRaw, "raw", unit)
}

type AccountBalanceAction struct {
	Action
	Account string `json:"account"`
}

type AccountBalance struct {
	Balance string `json:"balance"`
	Pending string `json:"pending"`
}

type AccountHistoryAction struct {
	Action
	Account string `json:"account"`
	Count   string `json:"count"`
}

type AccountHistory struct {
	History string `json:"history"`
}

type AccountKeyAction struct {
	Action
	Account string `json:"account"`
}

type AccountKey struct {
	Key string `json:"key"`
}

type AccountGetAction struct {
	Action
	Key string `json:"key"`
}

type AccountGet struct {
	Account string `json:"account"`
}

type Version struct {
	RpcVersion string `json:"rpc_version"`
	StoreVersion string `json:"store_version"`
	NodeVendor string `json:"node_vendor"`
}

func toUnit(input, inputUnit, outputUnit string) string {
	var e big.Int
	var b10 = big.NewInt(10)
	val, _ := new(big.Int).SetString(input, 10)

	// Step 1: to RAW
	switch inputUnit {
	case "raw":
	case "XRB":
		val.Mul(val, e.Exp(b10, big.NewInt(30), nil))
	case "Trai":
		val.Mul(val, e.Exp(b10, big.NewInt(36), nil))
	case "Grai":
		val.Mul(val, e.Exp(b10, big.NewInt(33), nil))
	case "Mrai":
		val.Mul(val, e.Exp(b10, big.NewInt(30), nil))
	case "krai":
		val.Mul(val, e.Exp(b10, big.NewInt(27), nil))
	case "rai":
		val.Mul(val, e.Exp(b10, big.NewInt(24), nil))
	case "mrai":
		val.Mul(val, e.Exp(b10, big.NewInt(21), nil))
	case "urai":
		val.Mul(val, e.Exp(b10, big.NewInt(18), nil))
	case "prai":
		val.Mul(val, e.Exp(b10, big.NewInt(15), nil))
	default:
	}
	// Step 2: to output
	switch outputUnit {
	case "raw":
	case "XRB":
		val.Div(val, e.Exp(b10, big.NewInt(30), nil))
	case "Trai":
		val.Div(val, e.Exp(b10, big.NewInt(36), nil))
	case "Grai":
		val.Div(val, e.Exp(b10, big.NewInt(33), nil))
	case "Mrai":
		val.Div(val, e.Exp(b10, big.NewInt(30), nil))
	case "krai":
		val.Div(val, e.Exp(b10, big.NewInt(27), nil))
	case "rai":
		val.Div(val, e.Exp(b10, big.NewInt(24), nil))
	case "mrai":
		val.Div(val, e.Exp(b10, big.NewInt(21), nil))
	case "urai":
		val.Div(val, e.Exp(b10, big.NewInt(18), nil))
	case "prai":
		val.Div(val, e.Exp(b10, big.NewInt(15), nil))
	default:
	}
	return val.String()
}
