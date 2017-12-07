package gorairpc

import (
	"fmt"
	"encoding/json"
	"net/http"
	"log"
	"strings"
	"io/ioutil"
	"math/big"
)

// Struct
type RaiRpc struct {
	url string
}

// RaiRpc constructor  with default url
func NewRaiRpc() RaiRpc {
	rb := RaiRpc{"http://localhost:7076"}
	return rb
}

func (r *RaiRpc) GetUrl() string {
	return r.url
}

func (r *RaiRpc) SetUrl(url string) {
	r.url = url
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

func (r *RaiRpc) RpcAccountBalance(account string) map[string]interface{} {
	params := map[string]interface{}{"action": "account_balance", "account": account}
	mapRes := r.callRpc(params)
	return mapRes
}

func (r *RaiRpc) RpcAccountBlockCount(account string) string {
	params := map[string]interface{}{"action": "account_block_count", "account": account}
	mapRes := r.callRpc(params)
	return mapRes["block_count"].(string)
}

func (r *RaiRpc) RpcAccountCreate(wallet string, work bool) string {
	// work = true
	params := map[string]interface{}{"action": "account_create", "wallet": wallet, "work": work}
	mapRes := r.callRpc(params)
	return mapRes["account"].(string)
}

func (r *RaiRpc) RpcAccountInfo(account, unit string, representative, weight, pending bool) map[string]interface{} {
	// account, unit = 'raw', representative = false, weight = false, pending = false
	params := map[string]interface{}{"action": "account_info", "account": account,
		"representative": representative, "weight": weight, "pending": pending}
	mapRes := r.callRpc(params)
	mapRes["balance"] = r.ToUnit(mapRes["balance"].(string), "raw", unit)
	if weight {
		mapRes["weight"] = r.ToUnit(mapRes["weight"].(string), "raw", unit)
	}
	if pending {
		mapRes["pending"] = r.ToUnit(mapRes["pending"].(string), "raw", unit)
	}
	return mapRes
}

func (r *RaiRpc) RpcAccountHistory(account, count string) string {
	// count = '4096'
	params := map[string]interface{}{"action": "account_history", "account": account, "count": count}
	mapRes := r.callRpc(params)
	return mapRes["history"].(string)
}

func (r *RaiRpc) RpcAccountGet(key string) string {
	params := map[string]interface{}{"action": "account_get", "key": key}
	mapRes := r.callRpc(params)
	return mapRes["account"].(string)
}

func (r *RaiRpc) RpcAccountKey(account string) string {
	params := map[string]interface{}{"action": "account_key", "account": account}
	mapRes := r.callRpc(params)
	return mapRes["key"].(string)
}

func (r *RaiRpc) RpcAccountList(wallet string) string {
	params := map[string]interface{}{"action": "account_list", "wallet": wallet}
	mapRes := r.callRpc(params)
	return mapRes["accounts"].(string)
}

func (r *RaiRpc) RpcAccountMove(wallet, source, accounts string) string {
	params := map[string]interface{}{"action": "account_move", "wallet": wallet, "source": source, "accounts": accounts}
	mapRes := r.callRpc(params)
	return mapRes["moved"].(string)
}

func (r *RaiRpc) RpcAccountRemove(wallet, account string) string {
	params := map[string]interface{}{"action": "account_remove", "wallet": wallet, "account": account}
	mapRes := r.callRpc(params)
	return mapRes["removed"].(string)
}

func (r *RaiRpc) RpcAccountRepresentative(account string) string {
	params := map[string]interface{}{"action": "account_representative", "account": account}
	mapRes := r.callRpc(params)
	return mapRes["representative"].(string)
}

func (r *RaiRpc) RpcAccountRepresentativeSet(wallet, account, representative, work string) string {
	// work = '0000000000000000'
	params := map[string]interface{}{"action": "account_representative_set", "wallet": wallet,
		"account": account, "representative": representative, "work": work}
	mapRes := r.callRpc(params)
	return mapRes["block"].(string)
}

func (r *RaiRpc) RpcAccountWeight(account, unit string) string {
	// unit = 'raw'
	params := map[string]interface{}{"action": "account_weight", "account": account}
	mapRes := r.callRpc(params)
	accountWeight := r.ToUnit(mapRes["weight"].(string), "raw", unit)
	return accountWeight
}

func (r *RaiRpc) RpcAccountsBalances(accounts string) string {
	params := map[string]interface{}{"action": "accounts_balances", "accounts": accounts}
	mapRes := r.callRpc(params)
	return mapRes["balances"].(string)
}

func (r *RaiRpc) RpcAccountsCreate(wallet, count, work string) string {
	// wallet, count = 1, work = true
	params := map[string]interface{}{"action": "accounts_create", "wallet": wallet, "count": count, "work": work}
	mapRes := r.callRpc(params)
	return mapRes["accounts"].(string)
}

func (r *RaiRpc) RpcAccountsFrontiers(accounts string) string {
	params := map[string]interface{}{"action": "accounts_frontiers", "accounts": accounts}
	mapRes := r.callRpc(params)
	return mapRes["frontiers"].(string)
}

func (r *RaiRpc) RpcAccountsPending(accounts []string, count string, threshold int, unit string, source bool) map[string]interface{} {
	// accounts, count = '4096', threshold = 0, unit = 'raw', source = false
	params := map[string]interface{}{"action": "accounts_pending", "accounts": accounts, "count": count,
		"threshold": threshold, "source": source}
	mapRes := r.callRpc(params)

	blocks := mapRes["blocks"].(map[string]interface{})
	if source {
		for acc, accv := range blocks {
			b := accv.(map[string]map[string]interface{})
			for _, hashv := range b {
				b[acc]["account"] = r.ToUnit(hashv["amount"].(string), "raw", unit)
				blocks[acc] = b
			}
		}
	} else if threshold != 0 {
		for acc, accv := range blocks {
			b := accv.(map[string]interface{})
			for hash, hashv := range b {
				b[hash] = r.ToUnit(hashv.(string), "raw", unit)
				blocks[acc] = b
			}
		}
	}
	return blocks
}

func (r *RaiRpc) RpcAvailableSupply(unit string) string {
	// unit = 'raw'
	params := map[string]interface{}{"action": "available_supply"}
	mapRes := r.callRpc(params)
	availableSupply := r.ToUnit(mapRes["available"].(string), "raw", unit)
	return availableSupply
}

func (r *RaiRpc) RpcBlock(hash string) map[string]interface{} {
	params := map[string]interface{}{"action": "block", "hash": hash}
	mapRes := r.callRpc(params)
	var block map[string]interface{}
	json.Unmarshal([]byte(mapRes["contents"].(string)), &block)
	return block
}

func (r *RaiRpc) RpcBlocks(hashes []string) map[string]interface{} {
	params := map[string]interface{}{"action": "blocks", "hashes": hashes}
	mapRes := r.callRpc(params)
	blocks := make(map[string]interface{})
	for k, v := range mapRes["blocks"].(map[string]interface{}) {
		var val interface{}
		json.Unmarshal([]byte(v.(string)), &val)
		blocks[k] = val
	}
	return blocks
}

func (r *RaiRpc) RpcBlocksInfo(hashes []string, unit string, pending, source bool) map[string]interface{} {
	// unit = 'raw', pending = false, source = false
	params := map[string]interface{}{"action": "blocks_info", "hashes": hashes, "pending": pending, "source": source}
	mapRes := r.callRpc(params)
	blocks := mapRes["blocks"].(map[string]interface{})
	for k, v := range blocks {
		val := v.(map[string]map[string]interface{})
		json.Unmarshal([]byte(val["contents"]["amount"].(string)), val["contents"]["amount"])
		blocks[k] = val
	}
	return blocks
}

func (r *RaiRpc) RpcBlockAccount(hash string) string {
	// unit = 'raw'
	params := map[string]interface{}{"action": "block_account", "hash": hash}
	mapRes := r.callRpc(params)
	return mapRes["account"].(string)
}

func (r *RaiRpc) RpcBlockCount() map[string]interface{} {
	params := map[string]interface{}{"action": "block_count"}
	mapRes := r.callRpc(params)
	return mapRes
}

func (r *RaiRpc) RpcBlockCountType() map[string]interface{} {
	params := map[string]interface{}{"action": "block_count_type"}
	mapRes := r.callRpc(params)
	return mapRes
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
func (r *RaiRpc) RpcBlockCreate(params map[string]interface{}) map[string]interface{} {
	params["action"] = "block_create"
	mapRes := r.callRpc(params)
	var block map[string]interface{}
	json.Unmarshal([]byte(mapRes["block"].(string)), &block)
	return block
}

func (r *RaiRpc) RpcBootstrap(address, port string) string {
	// address = '::ffff:138.201.94.249', port = '7075'
	params := map[string]interface{}{"action": "bootstrap", "address": address, "port": port}
	mapRes := r.callRpc(params)
	return mapRes["success"].(string)
}

func (r *RaiRpc) RpcBootstrapAny() string {
	params := map[string]interface{}{"action": "bootstrap_any"}
	mapRes := r.callRpc(params)
	return mapRes["success"].(string)
}

func (r *RaiRpc) RpcChain(block, count string) string {
	// count = '4096'
	params := map[string]interface{}{"action": "chain", "block": block, "count": count}
	mapRes := r.callRpc(params)
	return mapRes["blocks"].(string)
}

func (r *RaiRpc) RpcDelegators(account, unit string) map[string]string {
	// unit = 'raw'
	params := map[string]interface{}{"action": "delegators", "account": account}
	mapRes := r.callRpc(params)
	delegators := mapRes["delegators"].(map[string]string)
	if unit != "raw" {
		for k, v := range delegators {
			delegators[k] = r.ToUnit(v, "raw", unit)
		}
	}
	return delegators
}

func (r *RaiRpc) RpcDelegatorsCount(account string) string {
	params := map[string]interface{}{"action": "delegators_count", "account": account}
	mapRes := r.callRpc(params)
	return mapRes["count"].(string)
}

func (r *RaiRpc) RpcDeterministicKey(seed, index string) map[string]interface{} {
	// index = 0
	params := map[string]interface{}{"action": "deterministic_key", "seed": seed, "index": index}
	mapRes := r.callRpc(params)
	return mapRes
}

func (r *RaiRpc) RpcFrontiers(account, count string) string {
	// account = 'xrb_1111111111111111111111111111111111111111111111111117353trpda', count = '1048576'
	params := map[string]interface{}{"action": "frontiers", "account": account, "count": count}
	mapRes := r.callRpc(params)
	return mapRes["frontiers"].(string)
}

func (r *RaiRpc) RpcFrontierCount() string {
	params := map[string]interface{}{"action": "frontier_count"}
	mapRes := r.callRpc(params)
	return mapRes["count"].(string)
}

func (r *RaiRpc) RpcHistory(hash, count string) string {
	// count = '4096'
	params := map[string]interface{}{"action": "history", "hash": hash, "count": count}
	mapRes := r.callRpc(params)
	return mapRes["amount"].(string)
}

func (r *RaiRpc) RpcMraiFromRaw(amount string) string {
	params := map[string]interface{}{"action": "mrai_from_raw", "amount": amount}
	mapRes := r.callRpc(params)
	return mapRes["amount"].(string)
}

// Use ToUnit instead of this function
func (r *RaiRpc) RpcMraiToRaw(amount string) string {
	params := map[string]interface{}{"action": "mrai_to_raw", "amount": amount}
	mapRes := r.callRpc(params)
	return mapRes["amount"].(string)
}

func (r *RaiRpc) RpcKraiFromRaw(amount string) string {
	params := map[string]interface{}{"action": "krai_from_raw", "amount": amount}
	mapRes := r.callRpc(params)
	return mapRes["amount"].(string)
}

func (r *RaiRpc) RpcKraiToRaw(amount string) string {
	params := map[string]interface{}{"action": "krai_to_raw", "amount": amount}
	mapRes := r.callRpc(params)
	return mapRes["amount"].(string)
}

func (r *RaiRpc) RpcRaiFromRaw(amount string) string {
	params := map[string]interface{}{"action": "rai_from_raw", "amount": amount}
	mapRes := r.callRpc(params)
	return mapRes["amount"].(string)
}

func (r *RaiRpc) RpcRaiToRaw(amount string) string {
	params := map[string]interface{}{"action": "rai_to_raw", "amount": amount}
	mapRes := r.callRpc(params)
	return mapRes["amount"].(string)
}

func (r *RaiRpc) RpcKeepalive(address, port string) map[string]interface{} {
	// address = '::ffff:192.168.1.1', port = '7075'
	params := map[string]interface{}{"action": "keepalive", "address": address, "port": port}
	mapRes := r.callRpc(params)
	return mapRes
}

func (r *RaiRpc) RpcKeyCreate() map[string]interface{} {
	params := map[string]interface{}{"action": "key_create"}
	mapRes := r.callRpc(params)
	return mapRes
}

func (r *RaiRpc) RpcKeyExpand(key string) map[string]interface{} {
	params := map[string]interface{}{"action": "key_expand", "key": key}
	mapRes := r.callRpc(params)
	return mapRes
}

func (r *RaiRpc) RpcLedger(account, count string, representative, weight, pending, sorting bool) string {
	// account = 'xrb_1111111111111111111111111111111111111111111111111117353trpda', count = '1048576',
	// representative = false, weight = false, pending = false, sorting = false
	params := map[string]interface{}{"action": "ledger", "account": account, "count": count,
		"representative": representative, "weight": weight, "pending": pending, "sorting": sorting}
	mapRes := r.callRpc(params)
	return mapRes["accounts"].(string)
}

func (r *RaiRpc) RpcPasswordChange(wallet, password string) string {
	params := map[string]interface{}{"action": "password_change", "wallet": wallet, "password": password}
	mapRes := r.callRpc(params)
	return mapRes["changed"].(string)
}

func (r *RaiRpc) RpcPasswordEnter(wallet, password string) string {
	params := map[string]interface{}{"action": "password_enter", "wallet": wallet, "password": password}
	mapRes := r.callRpc(params)
	return mapRes["valid"].(string)
}

func (r *RaiRpc) RpcPasswordValid(wallet, password string) string {
	params := map[string]interface{}{"action": "password_valid", "wallet": wallet}
	mapRes := r.callRpc(params)
	return mapRes["valid"].(string)
}

func (r *RaiRpc) RpcPaymentBegin(wallet, password string) string {
	params := map[string]interface{}{"action": "payment_begin", "wallet": wallet}
	mapRes := r.callRpc(params)
	return mapRes["account"].(string)
}

func (r *RaiRpc) RpcPaymentInit(wallet string) string {
	params := map[string]interface{}{"action": "payment_init", "wallet": wallet}
	mapRes := r.callRpc(params)
	return mapRes["status"].(string)
}

func (r *RaiRpc) RpcPaymentEnd(account, wallet string) map[string]interface{} {
	params := map[string]interface{}{"action": "payment_end", "account": account, "wallet": wallet}
	mapRes := r.callRpc(params)
	return mapRes
}

func (r *RaiRpc) RpcPaymentWait(account, amount, timeout string) string {
	params := map[string]interface{}{"action": "payment_wait", "account": account, "amount": amount, "timeout": timeout}
	mapRes := r.callRpc(params)
	return mapRes["status"].(string)
}

func (r *RaiRpc) RpcProcess(block string) string {
	params := map[string]interface{}{"action": "process", "block": block}
	mapRes := r.callRpc(params)
	return mapRes["hash"].(string)
}

func (r *RaiRpc) RpcPeers() string {
	params := map[string]interface{}{"action": "peers"}
	mapRes := r.callRpc(params)
	return mapRes["peers"].(string)
}

func (r *RaiRpc) RpcPending(account, count string, threshold int, unit string, source bool) map[string]interface{} {
	// count = '4096', threshold = 0, unit = 'raw', source = false
	params := map[string]interface{}{"action":"pending","account":account,"count":count,"threshold":threshold,"source":source}
	mapRes := r.callRpc(params)
	blocks := mapRes["blocks"].(map[string]interface{})
	if source {
		for k, v := range blocks {
			val := v.(map[string]string)
			val["amount"] = r.ToUnit(val["amount"], "raw", unit)
			blocks[k] = val
		}
	} else if threshold != 0 {
		for hash, v := range blocks {
			blocks[hash] = r.ToUnit(v.(string), "raw", unit)
		}
	}
	return blocks
}

func (r *RaiRpc) RpcPendingExists(hash string) string {
	params := map[string]interface{}{"action": "pending_exists", "hash": hash}
	mapRes := r.callRpc(params)
	return mapRes["exists"].(string)
}

func (r *RaiRpc) RpcReceive(wallet, account, block, work string) string {
	// work = '0000000000000000'
	params := map[string]interface{}{"action":"receive","wallet":wallet,"account":account,"block":block,"work":work}
	mapRes := r.callRpc(params)
	return mapRes["block"].(string)
}

func (r *RaiRpc) RpcReceiveMinimum(unit string) string {
	// unit = 'raw'
	params := map[string]interface{}{"action": "receive_minimum"}
	mapRes := r.callRpc(params)
	amount := r.ToUnit(mapRes["amount"].(string), "raw", unit)
	return amount
}

func (r *RaiRpc) RpcReceiveMinimumSet(amount, unit string) string {
	// unit = 'raw'
	rawAmount := r.ToUnit(amount, unit, "raw")
	params := map[string]interface{}{"action": "receive_minimum_set", "amount": rawAmount}
	mapRes := r.callRpc(params)
	return mapRes["success"].(string)
}

func (r *RaiRpc) RpcRepresentatives(unit, count, sorting string) map[string]interface{} {
	// unit = 'raw', count = '1048576', sorting = false
	params := map[string]interface{}{"action": "representatives", "count": count, "sorting": sorting}
	mapRes := r.callRpc(params)
	representatives := mapRes["representatives"].(map[string]interface{})
	if unit != "raw" {
		for k, v := range representatives {
			representatives[k] = r.ToUnit(v.(string), "raw", unit)
		}
	}
	return representatives
}

func (r *RaiRpc) RpcRepublish(hash, count, sources string) string {
	// count = 1024, sources = 2
	params := map[string]interface{}{"action": "republish", "hash": hash, "count": count, "sources": sources}
	mapRes := r.callRpc(params)
	return mapRes["blocks"].(string)
}

func (r *RaiRpc) RpcSearchPending(wallet string) string {
	params := map[string]interface{}{"action": "search_pending", "wallet": wallet}
	mapRes := r.callRpc(params)
	return mapRes["started"].(string)
}

func (r *RaiRpc) RpcSearchPendingAll() string {
	params := map[string]interface{}{"action": "search_pending_all"}
	mapRes := r.callRpc(params)
	return mapRes["success"].(string)
}

func (r *RaiRpc) RpcSend(wallet, source, destination, amount, unit string) string {
	// unit = 'raw'
	rawAmount := r.ToUnit(amount, unit, "raw")
	params := map[string]interface{}{"action": "send", "wallet": wallet, "source": source, "destination": destination, "amount": rawAmount}
	mapRes := r.callRpc(params)
	return mapRes["block"].(string)
}

func (r *RaiRpc) RpcStop() string {
	params := map[string]interface{}{"action": "stop"}
	mapRes := r.callRpc(params)
	return mapRes["success"].(string)
}

func (r *RaiRpc) RpcSuccessors(block, count string) string {
	// count = '4096'
	params := map[string]interface{}{"action": "successors", "block": block, "count": count}
	mapRes := r.callRpc(params)
	return mapRes["blocks"].(string)
}

func (r *RaiRpc) RpcUnchecked(count string) map[string]interface{} {
	// count = '4096'
	params := map[string]interface{}{"action": "unchecked", "count": count}
	mapRes := r.callRpc(params)
	blocks := mapRes["blocks"].(map[string]interface{})
	for k, v := range blocks {
		var val interface{}
		json.Unmarshal([]byte(v.(string)), &val)
		blocks[k] = val
	}
	return blocks
}

func (r *RaiRpc) RpcUncheckedClear() string {
	params := map[string]interface{}{"action": "unchecked_clear"}
	mapRes := r.callRpc(params)
	return mapRes["success"].(string)
}

func (r *RaiRpc) RpcUncheckedGet(hash string) string {
	params := map[string]interface{}{"action": "unchecked_get", "hash": hash}
	mapRes := r.callRpc(params)
	return mapRes["contents"].(string)
}

func (r *RaiRpc) RpcUcheckedKeys(key, count string) map[string]interface{} {
	// key = '0000000000000000000000000000000000000000000000000000000000000000', count = '4096'
	params := map[string]interface{}{"action": "unchecked_keys", "key": key, "count": count}
	mapRes := r.callRpc(params)
	//var unchecked = unchecked_keys.unchecked;
	//for(let key in unchecked){
	//unchecked[key].contents = JSON.parse(unchecked[key].contents);
	//}
	return mapRes
}

func (r *RaiRpc) RpcValidateAccountNumber(account string) string {
	params := map[string]interface{}{"action": "validate_account_number", "account": account}
	mapRes := r.callRpc(params)
	return mapRes["valid"].(string)
}

func (r *RaiRpc) RpcVersion() map[string]interface{} {
	params := map[string]interface{}{"action": "version"}
	mapRes := r.callRpc(params)
	return mapRes
}

func (r *RaiRpc) RpcWalletAdd(wallet, key string) string {
	params := map[string]interface{}{"action": "wallet_add", "wallet": wallet, "key": key}
	mapRes := r.callRpc(params)
	return mapRes["account"].(string)
}

func (r *RaiRpc) RpcWalletBalanceTotal(wallet, unit string) map[string]interface{} {
	// unit = 'raw'
	params := map[string]interface{}{"action": "wallet_balance_total", "wallet": wallet}
	mapRes := r.callRpc(params)
	//var wallet_balance_total = { balance: this.unit(rpc_wallet_balance.balance, 'raw', unit), pending: this.unit(rpc_wallet_balance.pending, 'raw', unit) };
	//return wallet_balance_total;
	return mapRes
}

func (r *RaiRpc) RpcWalletBalances(wallet string) string {
	params := map[string]interface{}{"action": "wallet_balances", "wallet": wallet}
	mapRes := r.callRpc(params)
	return mapRes["balances"].(string)
}

func (r *RaiRpc) RpcWalletChangeSeed(wallet, seed string) string {
	params := map[string]interface{}{"action": "wallet_change_seed", "wallet": wallet, "seed": seed}
	mapRes := r.callRpc(params)
	return mapRes["success"].(string)
}

func (r *RaiRpc) RpcWalletContains(wallet, account string) string {
	params := map[string]interface{}{"action": "wallet_contains", "wallet": wallet, "account": account}
	mapRes := r.callRpc(params)
	return mapRes["exists"].(string)
}

func (r *RaiRpc) RpcWalletCreate() string {
	params := map[string]interface{}{"action": "wallet_create"}
	mapRes := r.callRpc(params)
	return mapRes["wallet"].(string)
}

func (r *RaiRpc) RpcWalletDestroy(wallet string) map[string]interface{} {
	params := map[string]interface{}{"action": "wallet_destroy", "wallet": wallet}
	mapRes := r.callRpc(params)
	return mapRes
}

func (r *RaiRpc) RpcWalletExport(wallet string) string {
	params := map[string]interface{}{"action": "wallet_export", "wallet": wallet}
	mapRes := r.callRpc(params)
	return mapRes["json"].(string)
}

func (r *RaiRpc) RpcWalletFrontiers(wallet string) string {
	params := map[string]interface{}{"action": "wallet_frontiers", "wallet": wallet}
	mapRes := r.callRpc(params)
	return mapRes["frontiers"].(string)
}

func (r *RaiRpc) RpcWalletPending(wallet, count, threshold, unit string) map[string]interface{} {
	// count = '4096', threshold = 0, unit = 'raw'
	if threshold != "0" {
		threshold = r.ToUnit(threshold, unit, "raw")
	}
	params := map[string]interface{}{"action": "wallet_pending", "wallet": wallet, "count": count, "threshold": threshold}
	mapRes := r.callRpc(params)

	if threshold != "0" {
		//for (let account in wallet_pending.blocks) {
		//for (let hash in wallet_pending.blocks[account]) {
		//wallet_pending.blocks[account][hash] = this.unit(wallet_pending.blocks[account][hash], 'raw', unit);
		//}
		//}
	}
	//return wallet_pending.blocks;
	return mapRes
}

func (r *RaiRpc) RpcWalletRepresentative(wallet string) string {
	params := map[string]interface{}{"action": "wallet_representative", "wallet": wallet}
	mapRes := r.callRpc(params)
	return mapRes["representative"].(string)
}

func (r *RaiRpc) RpcWalletRepresentativeSet(wallet, representative string) string {
	params := map[string]interface{}{"action": "wallet_representative_set", "wallet": wallet, "representative": representative}
	mapRes := r.callRpc(params)
	return mapRes["set"].(string)
}

func (r *RaiRpc) RpcWalletRepublish(wallet, count string) string {
	//  count = 2
	params := map[string]interface{}{"action": "wallet_republish", "wallet": wallet, "count": count}
	mapRes := r.callRpc(params)
	return mapRes["blocks"].(string)
}

func (r *RaiRpc) RpcWalletWorkGet(wallet string) string {
	params := map[string]interface{}{"action": "wallet_work_get", "wallet": wallet}
	mapRes := r.callRpc(params)
	return mapRes["works"].(string)
}

func (r *RaiRpc) RpcWorkCancel(hash string) map[string]interface{} {
	params := map[string]interface{}{"action": "work_cancel", "hash": hash}
	mapRes := r.callRpc(params)
	return mapRes
}

func (r *RaiRpc) RpcWorkGenerate(hash string) string {
	params := map[string]interface{}{"action": "work_generate", "hash": hash}
	mapRes := r.callRpc(params)
	return mapRes["work"].(string)
}

func (r *RaiRpc) RpcWorkGet(wallet, account string) string {
	params := map[string]interface{}{"action": "work_get", "wallet": wallet, "account": account}
	mapRes := r.callRpc(params)
	return mapRes["work"].(string)
}

func (r *RaiRpc) RpcWorkSet(wallet, account, work string) string {
	params := map[string]interface{}{"action": "work_set", "wallet": wallet, "account": account, "work": work}
	mapRes := r.callRpc(params)
	return mapRes["success"].(string)
}

func (r *RaiRpc) RpcWorkValidate(work, hash string) string {
	params := map[string]interface{}{"action": "work_validate", "work": work, "hash": hash}
	mapRes := r.callRpc(params)
	return mapRes["valid"].(string)
}

func (r *RaiRpc) RpcWorkPeerAdd(address, port string) string {
	// address = '::1', port = '7076'
	params := map[string]interface{}{"action": "work_peer_add", "address": address, "port": port}
	mapRes := r.callRpc(params)
	return mapRes["success"].(string)
}

func (r *RaiRpc) RpcWorkPeers() string {
	params := map[string]interface{}{"action": "work_peers"}
	mapRes := r.callRpc(params)
	return mapRes["work_peers"].(string)
}

func (r *RaiRpc) RpcWorkPeersClear() string {
	params := map[string]interface{}{"action": "work_peers_clear"}
	mapRes := r.callRpc(params)
	return mapRes["success"].(string)
}

func (r *RaiRpc) callRpc(params map[string]interface{}) map[string]interface{} {
	// Prepare json POST request
	reqString, err := json.Marshal(params)

	res, err := http.Post(r.url, "application/json", strings.NewReader(string(reqString)))
	if err != nil {
		log.Fatal(err)
	}

	var byteTab []byte
	byteTab, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	var dataMap map[string]interface{}
	err = json.Unmarshal(byteTab, &dataMap)
	if err != nil {
		//error handling goes here
		fmt.Printf("err %s", err)
	}
	return dataMap
}
