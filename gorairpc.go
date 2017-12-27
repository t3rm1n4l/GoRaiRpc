package gorairpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
)

// Struct
type RaiRpc struct {
	url string
}

// RaiRpc constructor  with default url
func New(url ...string) RaiRpc {
	rb := RaiRpc{
		url: "http://localhost:7076",
	}

	if len(url) != 0 && url[0] != "" {
		rb.url = url[0]
	}

	return rb
}

func (r *RaiRpc) ToUnit(input, inputUnit, outputUnit string) string {
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

func (r *RaiRpc) AccountBalance(account string) (*AccountBalance, error) {
	resp := new(AccountBalance)
	if err := r.call(AccountBalanceAction{Action: Action{Name: "account_balance"}, Account: account}, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *RaiRpc) RpcAccountBlockCount(account string) (string, error) {
	params := map[string]interface{}{"action": "account_block_count", "account": account}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["block_count"].(string), nil
}

func (r *RaiRpc) RpcAccountCreate(wallet string, work bool) (string, error) {
	params := map[string]interface{}{"action": "account_create", "wallet": wallet, "work": work}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["account"].(string), nil
}

func (r *RaiRpc) RpcAccountInfo(account, unit string, representative, weight, pending bool) (map[string]interface{}, error) {
	params := map[string]interface{}{"action": "account_info", "account": account,
		"representative": representative, "weight": weight, "pending": pending}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return nil, err
	}

	mapRes["balance"] = r.ToUnit(mapRes["balance"].(string), "raw", unit)
	if weight {
		mapRes["weight"] = r.ToUnit(mapRes["weight"].(string), "raw", unit)
	}
	if pending {
		mapRes["pending"] = r.ToUnit(mapRes["pending"].(string), "raw", unit)
	}
	return mapRes, nil
}

// todo is count always a number?
func (r *RaiRpc) AccountHistory(account, count string) (string, error) {
	resp := new(AccountHistory)
	err := r.call(AccountHistoryAction{Action:Action{Name:"account_history"}, Count: count, Account:account}, resp)
	if err != nil {
		return "", err
	}
	return resp.History, nil
}

func (r *RaiRpc) AccountGet(key string) (string, error) {
	resp := new(AccountGet)
	err := r.call(AccountGetAction{Action:Action{Name:"account_get"}, Key:key}, resp)
	if err != nil {
		return "", err
	}
	return resp.Account, nil
}

func (r *RaiRpc) AccountKey(account string) (string, error) {
	resp := new(AccountKey)
	err := r.call(AccountKeyAction{Action: Action{Name:"account_key"}, Account:account}, resp)
	if err != nil {
		return "", err
	}
	return resp.Key, nil
}

func (r *RaiRpc) RpcAccountList(wallet string) (string, error) {
	params := map[string]interface{}{"action": "account_list", "wallet": wallet}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["accounts"].(string), nil
}

func (r *RaiRpc) RpcAccountMove(wallet, source, accounts string) (string, error) {
	params := map[string]interface{}{"action": "account_move", "wallet": wallet,
		"source": source, "accounts": accounts}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}

	return mapRes["moved"].(string), nil
}

func (r *RaiRpc) RpcAccountRemove(wallet, account string) (string, error) {
	params := map[string]interface{}{"action": "account_remove", "wallet": wallet, "account": account}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["removed"].(string), nil
}

func (r *RaiRpc) RpcAccountRepresentative(account string) (string, error) {
	params := map[string]interface{}{"action": "account_representative", "account": account}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["representative"].(string), nil
}

func (r *RaiRpc) RpcAccountRepresentativeSet(wallet, account, representative, work string) (string, error) {
	// work = "0000000000000000"
	params := map[string]interface{}{"action": "account_representative_set", "wallet": wallet,
		"account": account, "representative": representative, "work": work}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["block"].(string), nil
}

func (r *RaiRpc) RpcAccountWeight(account, unit string) (string, error) {
	// unit = "raw"
	params := map[string]interface{}{"action": "account_weight", "account": account}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	accountWeight := r.ToUnit(mapRes["weight"].(string), "raw", unit)
	return accountWeight, nil
}

func (r *RaiRpc) RpcAccountsBalances(accounts string) (string, error) {
	params := map[string]interface{}{"action": "accounts_balances", "accounts": accounts}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["balances"].(string), nil
}

func (r *RaiRpc) RpcAccountsCreate(wallet, count, work string) (string, error) {
	// wallet, count = 1, work = true
	params := map[string]interface{}{"action": "accounts_create", "wallet": wallet, "count": count, "work": work}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["accounts"].(string), nil
}

func (r *RaiRpc) RpcAccountsFrontiers(accounts string) (string, error) {
	params := map[string]interface{}{"action": "accounts_frontiers", "accounts": accounts}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["frontiers"].(string), nil
}

func (r *RaiRpc) RpcAccountsPending(accounts []string, count string, threshold int, unit string, source bool) (map[string]interface{}, error) {
	// accounts, count = "4096", threshold = 0, unit = "raw", source = false
	params := map[string]interface{}{"action": "accounts_pending", "accounts": accounts, "count": count,
		"threshold": threshold, "source": source}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return nil, err
	}

	blocks := mapRes["blocks"].(map[string]interface{})
	if source {
		for acc, accv := range blocks {
			for hkey, hashv := range accv.(map[string]interface{}) {
				blocks[acc].(map[string]interface{})[hkey].(map[string]interface{})["amount"] = r.ToUnit(hashv.(map[string]interface{})["amount"].(string), "raw", unit)
			}
		}
	} else if threshold != 0 {
		for acc, accv := range blocks {
			for hash, hashv := range accv.(map[string]interface{}) {
				blocks[acc].(map[string]interface{})[hash] = r.ToUnit(hashv.(string), "raw", unit)
			}
		}
	}
	return blocks, nil
}

func (r *RaiRpc) AvailableSupply() (*AvailableSupply, error) {
	// unit = "raw"
	data := new(AvailableSupply)
	err := r.call(Action{Name: "available_supply"}, data)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (r *RaiRpc) RpcBlock(hash string) (map[string]interface{}, error) {
	params := map[string]interface{}{"action": "block", "hash": hash}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return nil, err
	}
	var block map[string]interface{}
	if err := json.Unmarshal([]byte(mapRes["contents"].(string)), &block); err != nil {
		return nil, err
	}
	return block, nil
}

func (r *RaiRpc) RpcBlocks(hashes []string) (map[string]interface{}, error) {
	params := map[string]interface{}{"action": "blocks", "hashes": hashes}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return nil, err
	}
	blocks := make(map[string]interface{})
	for k, v := range mapRes["blocks"].(map[string]interface{}) {
		var val interface{}
		if err := json.Unmarshal([]byte(v.(string)), &val); err != nil {
			return nil, err
		}
		blocks[k] = val
	}
	return blocks, nil
}

func (r *RaiRpc) RpcBlocksInfo(hashes []string, unit string, pending, source bool) (map[string]interface{}, error) {
	// unit = "raw", pending = false, source = false
	params := map[string]interface{}{"action": "blocks_info", "hashes": hashes, "pending": pending, "source": source}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return nil, err
	}
	blocks := mapRes["blocks"].(map[string]interface{})
	for k, v := range blocks {
		var val interface{}
		if err := json.Unmarshal([]byte(v.(map[string]interface{})["contents"].(string)), &val); err != nil {
			return nil, err
		}
		blocks[k].(map[string]interface{})["contents"] = val
		if unit != "raw" {
			v.(map[string]interface{})["amount"] = r.ToUnit(v.(map[string]interface{})["amount"].(string), "raw", unit)
		}
	}
	return blocks, nil
}

func (r *RaiRpc) RpcBlockAccount(hash string) (string, error) {
	// unit = "raw"
	params := map[string]interface{}{"action": "block_account", "hash": hash}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", nil
	}
	return mapRes["account"].(string), nil
}

func (r *RaiRpc) RpcBlockCount() (map[string]interface{}, error) {
	params := map[string]interface{}{"action": "block_count"}
	return r.callRpc(params)
}

func (r *RaiRpc) RpcBlockCountType() (map[string]interface{}, error) {
	params := map[string]interface{}{"action": "block_count_type"}
	return r.callRpc(params)
}

/*	Sample block creation:
	blockData := map[string]interface{}
	blockData["type"] = "open"
	blockData["key"] = "0000000000000000000000000000000000000000000000000000000000000001"
	blockData["account"] = xrb_3kdbxitaj7f6mrir6miiwtw4muhcc58e6tn5st6rfaxsdnb7gr4roudwn951"
	blockData["representative"] = "xrb_1hza3f7wiiqa7ig3jczyxj5yo86yegcmqk3criaz838j91sxcckpfhbhhra1"
	blockData["source"] = "19D3D919475DEED4696B5D13018151D1AF88B2BD3BCFF048B45031C1F36D1858"
	block := rpc.RpcBlockCreate(blockData)
*/
func (r *RaiRpc) RpcBlockCreate(params map[string]interface{}) (map[string]interface{}, error) {
	params["action"] = "block_create"
	mapRes, err := r.callRpc(params)
	if err != nil {
		return nil, err
	}

	var block map[string]interface{}
	if err := json.Unmarshal([]byte(mapRes["block"].(string)), &block); err != nil {
		return nil, err
	}

	return block, nil
}

func (r *RaiRpc) RpcBootstrap(address, port string) (string, error) {
	// address = "::ffff:138.201.94.249", port = "7075"
	params := map[string]interface{}{"action": "bootstrap", "address": address, "port": port}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}

	return mapRes["success"].(string), nil
}

func (r *RaiRpc) RpcBootstrapAny() (string, error) {
	params := map[string]interface{}{"action": "bootstrap_any"}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["success"].(string), nil
}

func (r *RaiRpc) RpcChain(block, count string) (string, error) {
	// count = "4096"
	params := map[string]interface{}{"action": "chain", "block": block, "count": count}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["blocks"].(string), nil
}

func (r *RaiRpc) RpcDelegators(account, unit string) (map[string]string, error) {
	// unit = "raw"
	params := map[string]interface{}{"action": "delegators", "account": account}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return nil, err
	}

	delegators := mapRes["delegators"].(map[string]string)
	if unit != "raw" {
		for k, v := range delegators {
			delegators[k] = r.ToUnit(v, "raw", unit)
		}
	}
	return delegators, nil
}

func (r *RaiRpc) RpcDelegatorsCount(account string) (string, error) {
	params := map[string]interface{}{"action": "delegators_count", "account": account}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["count"].(string), nil
}

func (r *RaiRpc) RpcDeterministicKey(seed, index string) (map[string]interface{}, error) {
	// index = 0
	params := map[string]interface{}{"action": "deterministic_key", "seed": seed, "index": index}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return nil, err
	}
	return mapRes, nil
}

func (r *RaiRpc) RpcFrontiers(account, count string) (string, error) {
	// account = "xrb_1111111111111111111111111111111111111111111111111117353trpda", count = "1048576"
	params := map[string]interface{}{"action": "frontiers", "account": account, "count": count}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["frontiers"].(string), nil
}

func (r *RaiRpc) RpcFrontierCount() (string, error) {
	params := map[string]interface{}{"action": "frontier_count"}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["count"].(string), nil
}

func (r *RaiRpc) RpcHistory(hash, count string) (string, error) {
	// count = "4096"
	params := map[string]interface{}{"action": "history", "hash": hash, "count": count}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["amount"].(string), nil
}

func (r *RaiRpc) RpcMraiFromRaw(amount string) (string, error) {
	params := map[string]interface{}{"action": "mrai_from_raw", "amount": amount}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["amount"].(string), nil
}

// Use ToUnit instead of this function
func (r *RaiRpc) RpcMraiToRaw(amount string) (string, error) {
	params := map[string]interface{}{"action": "mrai_to_raw", "amount": amount}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["amount"].(string), nil
}

func (r *RaiRpc) RpcKraiFromRaw(amount string) (string, error) {
	params := map[string]interface{}{"action": "krai_from_raw", "amount": amount}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["amount"].(string), nil
}

func (r *RaiRpc) RpcKraiToRaw(amount string) (string, error) {
	params := map[string]interface{}{"action": "krai_to_raw", "amount": amount}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["amount"].(string), nil
}

func (r *RaiRpc) RpcRaiFromRaw(amount string) (string, error) {
	params := map[string]interface{}{"action": "rai_from_raw", "amount": amount}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["amount"].(string), nil
}

func (r *RaiRpc) RpcRaiToRaw(amount string) (string, error) {
	params := map[string]interface{}{"action": "rai_to_raw", "amount": amount}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["amount"].(string), nil
}

func (r *RaiRpc) RpcKeepalive(address, port string) (map[string]interface{}, error) {
	// address = "::ffff:192.168.1.1", port = "7075"
	params := map[string]interface{}{"action": "keepalive", "address": address, "port": port}
	return r.callRpc(params)
}

func (r *RaiRpc) RpcKeyCreate() (map[string]interface{}, error) {
	params := map[string]interface{}{"action": "key_create"}
	return r.callRpc(params)

}

func (r *RaiRpc) RpcKeyExpand(key string) (map[string]interface{}, error) {
	params := map[string]interface{}{"action": "key_expand", "key": key}
	return r.callRpc(params)
}

func (r *RaiRpc) RpcLedger(account, count string, representative, weight, pending, sorting bool) (string, error) {
	// account = "xrb_1111111111111111111111111111111111111111111111111117353trpda", count = "1048576",
	// representative = false, weight = false, pending = false, sorting = false
	params := map[string]interface{}{"action": "ledger", "account": account, "count": count,
		"representative": representative, "weight": weight, "pending": pending, "sorting": sorting}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["accounts"].(string), nil
}

func (r *RaiRpc) RpcPasswordChange(wallet, password string) (string, error) {
	params := map[string]interface{}{"action": "password_change", "wallet": wallet, "password": password}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["changed"].(string), nil
}

func (r *RaiRpc) RpcPasswordEnter(wallet, password string) (string, error) {
	params := map[string]interface{}{"action": "password_enter", "wallet": wallet, "password": password}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["valid"].(string), nil
}

func (r *RaiRpc) RpcPasswordValid(wallet, password string) (string, error) {
	params := map[string]interface{}{"action": "password_valid", "wallet": wallet}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["valid"].(string), nil
}

func (r *RaiRpc) RpcPaymentBegin(wallet, password string) (string, error) {
	params := map[string]interface{}{"action": "payment_begin", "wallet": wallet}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["account"].(string), nil
}

func (r *RaiRpc) RpcPaymentInit(wallet string) (string, error) {
	params := map[string]interface{}{"action": "payment_init", "wallet": wallet}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["status"].(string), nil
}

func (r *RaiRpc) RpcPaymentEnd(account, wallet string) (map[string]interface{}, error) {
	params := map[string]interface{}{"action": "payment_end", "account": account, "wallet": wallet}
	return r.callRpc(params)
}

func (r *RaiRpc) RpcPaymentWait(account, amount, timeout string) (string, error) {
	params := map[string]interface{}{"action": "payment_wait", "account": account, "amount": amount, "timeout": timeout}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}

	return mapRes["status"].(string), nil
}

func (r *RaiRpc) RpcProcess(block string) (string, error) {
	params := map[string]interface{}{"action": "process", "block": block}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["hash"].(string), nil
}

func (r *RaiRpc) RpcPeers() (string, error) {
	params := map[string]interface{}{"action": "peers"}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["peers"].(string), nil
}

func (r *RaiRpc) RpcPending(account, count string, threshold int, unit string, source bool) (map[string]interface{}, error) {
	// count = "4096", threshold = 0, unit = "raw", source = false
	params := map[string]interface{}{"action": "pending", "account": account,
		"count": count, "threshold": threshold, "source": source}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return nil, err
	}
	blocks := mapRes["blocks"].(map[string]interface{})
	if source {
		for k, v := range blocks {
			blocks[k].(map[string]interface{})["amount"] = r.ToUnit(v.(map[string]interface{})["amount"].(string), "raw", unit)
		}
	} else if threshold != 0 {
		for hash, v := range blocks {
			blocks[hash] = r.ToUnit(v.(string), "raw", unit)
		}
	}
	return blocks, nil
}

func (r *RaiRpc) RpcPendingExists(hash string) (string, error) {
	params := map[string]interface{}{"action": "pending_exists", "hash": hash}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}

	return mapRes["exists"].(string), nil
}

func (r *RaiRpc) RpcReceive(wallet, account, block, work string) (string, error) {
	// work = "0000000000000000"
	params := map[string]interface{}{"action": "receive", "wallet": wallet,
		"account": account, "block": block, "work": work}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["block"].(string), nil
}

func (r *RaiRpc) RpcReceiveMinimum(unit string) (string, error) {
	// unit = "raw"
	params := map[string]interface{}{"action": "receive_minimum"}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	amount := r.ToUnit(mapRes["amount"].(string), "raw", unit)
	return amount, nil
}

func (r *RaiRpc) RpcReceiveMinimumSet(amount, unit string) (string, error) {
	// unit = "raw"
	rawAmount := r.ToUnit(amount, unit, "raw")
	params := map[string]interface{}{"action": "receive_minimum_set", "amount": rawAmount}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["success"].(string), nil
}

func (r *RaiRpc) RpcRepresentatives(unit, count, sorting string) (map[string]interface{}, error) {
	// unit = "raw", count = "1048576", sorting = false
	params := map[string]interface{}{"action": "representatives", "count": count, "sorting": sorting}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return nil, err
	}
	representatives := mapRes["representatives"].(map[string]interface{})
	if unit != "raw" {
		for k, v := range representatives {
			representatives[k] = r.ToUnit(v.(string), "raw", unit)
		}
	}
	return representatives, nil
}

func (r *RaiRpc) RpcRepublish(hash, count, sources string) (string, error) {
	// count = 1024, sources = 2
	params := map[string]interface{}{"action": "republish", "hash": hash, "count": count, "sources": sources}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["blocks"].(string), nil
}

func (r *RaiRpc) RpcSearchPending(wallet string) (string, error) {
	params := map[string]interface{}{"action": "search_pending", "wallet": wallet}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["started"].(string), nil
}

func (r *RaiRpc) RpcSearchPendingAll() (string, error) {
	params := map[string]interface{}{"action": "search_pending_all"}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["success"].(string), nil
}

func (r *RaiRpc) RpcSend(wallet, source, destination, amount, unit string) (string, error) {
	// unit = "raw"
	rawAmount := r.ToUnit(amount, unit, "raw")
	params := map[string]interface{}{"action": "send", "wallet": wallet, "source": source,
		"destination": destination, "amount": rawAmount}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["block"].(string), nil
}

func (r *RaiRpc) RpcStop() (string, error) {
	params := map[string]interface{}{"action": "stop"}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["success"].(string), nil
}

func (r *RaiRpc) RpcSuccessors(block, count string) (string, error) {
	// count = "4096"
	params := map[string]interface{}{"action": "successors", "block": block, "count": count}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["blocks"].(string), nil
}

func (r *RaiRpc) RpcUnchecked(count string) (map[string]interface{}, error) {
	// count = "4096"
	params := map[string]interface{}{"action": "unchecked", "count": count}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return nil, err
	}
	blocks := mapRes["blocks"].(map[string]interface{})
	for k, v := range blocks {
		var val interface{}
		if err := json.Unmarshal([]byte(v.(string)), &val); err != nil {
			return nil, err
		}
		blocks[k] = val
	}
	return blocks, nil
}

func (r *RaiRpc) RpcUncheckedClear() (string, error) {
	params := map[string]interface{}{"action": "unchecked_clear"}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["success"].(string), nil
}

func (r *RaiRpc) RpcUncheckedGet(hash string) (string, error) {
	params := map[string]interface{}{"action": "unchecked_get", "hash": hash}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["contents"].(string), nil
}

func (r *RaiRpc) RpcUcheckedKeys(key, count string) (interface{}, error) {
	// key = "0000000000000000000000000000000000000000000000000000000000000000", count = "4096"
	params := map[string]interface{}{"action": "unchecked_keys", "key": key, "count": count}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	unchecked := mapRes["unchecked"].(interface{})
	for k, v := range unchecked.([]interface{}) {
		var contents interface{}
		if err := json.Unmarshal([]byte(v.(map[string]interface{})["contents"].(string)), &contents); err != nil {
			return nil, err
		}
		unchecked.([]interface{})[k].(map[string]interface{})["contents"] = contents
	}
	return unchecked, nil
}

func (r *RaiRpc) RpcValidateAccountNumber(account string) (string, error) {
	params := map[string]interface{}{"action": "validate_account_number", "account": account}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["valid"].(string), nil
}

func (r *RaiRpc) Version() (*Version, error) {
	resp := new(Version)
	err := r.call(Action{Name:"version"}, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *RaiRpc) RpcWalletAdd(wallet, key string) (string, error) {
	params := map[string]interface{}{"action": "wallet_add", "wallet": wallet, "key": key}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["account"].(string), nil
}

func (r *RaiRpc) RpcWalletBalanceTotal(wallet, unit string) (map[string]interface{}, error) {
	// unit = "raw"
	params := map[string]interface{}{"action": "wallet_balance_total", "wallet": wallet}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return nil, err
	}
	walletBalanceTotals := map[string]interface{}{"balance": r.ToUnit(mapRes["balance"].(string), "raw", unit),
		"pending": r.ToUnit(mapRes["pending"].(string), "raw", unit)}
	return walletBalanceTotals, nil
}

func (r *RaiRpc) RpcWalletBalances(wallet, unit string, threshold int) (map[string]interface{}, error) {
	// unit = "raw", threshold = 0
	if threshold != 0 {
		threshold, _ = strconv.Atoi(r.ToUnit(strconv.Itoa(threshold), unit, "raw"))
	}
	params := map[string]interface{}{"action": "wallet_balances", "wallet": wallet, "threshold": threshold}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return nil, err
	}
	walletBalances := mapRes["balances"].(map[string]interface{})
	for k, v := range walletBalances {
		balance := r.ToUnit(v.(map[string]interface{})["balance"].(string), "raw", unit)
		pending := r.ToUnit(v.(map[string]interface{})["pending"].(string), "raw", unit)
		walletBalances[k].(map[string]interface{})["balance"] = balance
		walletBalances[k].(map[string]interface{})["pending"] = pending
	}
	return walletBalances, nil
}

func (r *RaiRpc) RpcWalletChangeSeed(wallet, seed string) (string, error) {
	params := map[string]interface{}{"action": "wallet_change_seed", "wallet": wallet, "seed": seed}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["success"].(string), nil
}

func (r *RaiRpc) RpcWalletContains(wallet, account string) (string, error) {
	params := map[string]interface{}{"action": "wallet_contains", "wallet": wallet, "account": account}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["exists"].(string), nil
}

func (r *RaiRpc) RpcWalletCreate() (string, error) {
	params := map[string]interface{}{"action": "wallet_create"}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["wallet"].(string), nil
}

func (r *RaiRpc) RpcWalletDestroy(wallet string) (map[string]interface{}, error) {
	params := map[string]interface{}{"action": "wallet_destroy", "wallet": wallet}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return nil, err
	}
	return mapRes, nil
}

func (r *RaiRpc) RpcWalletExport(wallet string) (string, error) {
	params := map[string]interface{}{"action": "wallet_export", "wallet": wallet}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["json"].(string), nil
}

func (r *RaiRpc) RpcWalletFrontiers(wallet string) (string, error) {
	params := map[string]interface{}{"action": "wallet_frontiers", "wallet": wallet}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["frontiers"].(string), nil
}

func (r *RaiRpc) RpcWalletPending(wallet, count string, threshold int, unit string, source bool) (map[string]interface{}, error) {
	//count = "4096", threshold = 0, unit = "raw", source = false
	thresholdStr := "0"
	if threshold != 0 {
		thresholdStr = r.ToUnit(strconv.Itoa(threshold), unit, "raw")
	}
	params := map[string]interface{}{"action": "wallet_pending", "wallet": wallet,
		"count": count, "threshold": thresholdStr, "source": source}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return nil, err
	}
	if mapRes["blocks"] == "" {
		return nil, err
	}
	blocks := mapRes["blocks"].(map[string]interface{})
	if source {
		for k, v := range blocks {
			val := v.(map[string]interface{})
			val["amount"] = r.ToUnit(val["amount"].(string), "raw", unit)
			blocks[k] = val
		}
	} else if threshold != 0 {
		for hash, v := range blocks {
			blocks[hash] = r.ToUnit(v.(string), "raw", unit)
		}
	}
	return blocks, err
}

func (r *RaiRpc) RpcWalletRepresentative(wallet string) (string, error) {
	params := map[string]interface{}{"action": "wallet_representative", "wallet": wallet}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["representative"].(string), nil
}

func (r *RaiRpc) RpcWalletRepresentativeSet(wallet, representative string) (string, error) {
	params := map[string]interface{}{"action": "wallet_representative_set",
		"wallet": wallet, "representative": representative}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["set"].(string), nil
}

func (r *RaiRpc) RpcWalletRepublish(wallet, count string) (string, error) {
	//  count = 2
	params := map[string]interface{}{"action": "wallet_republish", "wallet": wallet, "count": count}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["blocks"].(string), nil
}

func (r *RaiRpc) RpcWalletWorkGet(wallet string) (string, error) {
	params := map[string]interface{}{"action": "wallet_work_get", "wallet": wallet}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["works"].(string), nil
}

func (r *RaiRpc) RpcWorkCancel(hash string) (map[string]interface{}, error) {
	params := map[string]interface{}{"action": "work_cancel", "hash": hash}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return nil, err
	}
	return mapRes, nil
}

func (r *RaiRpc) RpcWorkGenerate(hash string) (string, error) {
	params := map[string]interface{}{"action": "work_generate", "hash": hash}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["work"].(string), nil
}

func (r *RaiRpc) RpcWorkGet(wallet, account string) (string, error) {
	params := map[string]interface{}{"action": "work_get", "wallet": wallet, "account": account}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["work"].(string), nil
}

func (r *RaiRpc) RpcWorkSet(wallet, account, work string) (string, error) {
	params := map[string]interface{}{"action": "work_set", "wallet": wallet, "account": account, "work": work}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["success"].(string), nil
}

func (r *RaiRpc) RpcWorkValidate(work, hash string) (string, error) {
	params := map[string]interface{}{"action": "work_validate", "work": work, "hash": hash}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["valid"].(string), nil
}

func (r *RaiRpc) RpcWorkPeerAdd(address, port string) (string, error) {
	// address = "::1", port = "7076"
	params := map[string]interface{}{"action": "work_peer_add", "address": address, "port": port}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["success"].(string), nil
}

func (r *RaiRpc) RpcWorkPeers() (string, error) {
	params := map[string]interface{}{"action": "work_peers"}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["work_peers"].(string), nil
}

func (r *RaiRpc) RpcWorkPeersClear() (string, error) {
	params := map[string]interface{}{"action": "work_peers_clear"}
	mapRes, err := r.callRpc(params)
	if err != nil {
		return "", err
	}
	return mapRes["success"].(string), nil
}

func (r *RaiRpc) callRpc(params map[string]interface{}) (map[string]interface{}, error) {
	// Prepare json POST request
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	res, err := http.Post(r.url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error sending request (%v)", err)
	}
	defer res.Body.Close()

	data := make(map[string]interface{})
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("error unmarschaling response body (%v)", err)
	}
	return data, nil
}

func (r *RaiRpc) call(params interface{}, dest interface{}) error {
	// Prepare json POST request
	body, err := json.Marshal(params)
	if err != nil {
		return err
	}

	resp, err := http.Post(r.url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error sending request (%v)", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("rai blockes node responded with status code (%d)", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&dest); err != nil {
		return fmt.Errorf("error unmarschaling response body (%v)", err)
	}
	return nil
}
