package main

import (
	"fmt"
	"github.com/devmahno/GoRaiRpc"
)

func main() {
	rpc := NewRaiRpc()
	fmt.Printf("Begin\n")
	//fmt.Printf("AvailableSupply: %s\n", rpc.RpcAvailableSupply("raw"))
	//showData(rpc.RpcVersion())
	//fmt.Printf("\n")
	//fmt.Println("RpcBlock:", rpc.RpcBlock("3BDCF72E7662DB9E2F38B15AF8FD1B633232928ED323C44E83D63573889E9BAE"))
	//fmt.Printf("\n")
	//hashes := []string{"3BDCF72E7662DB9E2F38B15AF8FD1B633232928ED323C44E83D63573889E9BAE"}
	//showData(rpc.RpcBlocks(hashes))
	//fmt.Printf("\n")

	//rpc.RpcAccountsPending([]string{"xrb_1111111111111111111111111111111111111111111111111117353trpda",
	//"xrb_3t6k35gi95xu6tergt6p69ck76ogmitsa8mnijtpxm9fkcm736xtoncuohr3"}, "4096", 0, "raw", true)

	rpc.RpcBlocksInfo([]string{"000D1BAEC8EC208142C99059B393051BAC8380F9B5A2E6B2489A277D81789F3F"}, "raw", false, false)

	fmt.Println("RpcRepresentatives:", rpc.RpcRepresentatives("Mrai", "1048576", "false"))
	fmt.Printf("\n")
	//show_data(RpcFrontiers(
	//	"xrb_1111111111111111111111111111111111111111111111111117353trpda",
	//	"100"))
	fmt.Printf("MraiToRaw: %v\n", rpc.RpcMraiToRaw("1"))
	fmt.Printf("%v\n", rpc.ToUnit("1", "Mrai", "raw"))

	fmt.Printf("End\n")
}

func showData(data map[string]interface{}) {
	for k, v := range data {
		switch vv := v.(type) {
		case string:
			fmt.Println(k, "is string", vv)
		case float64:
			fmt.Println(k, "is float64", vv)
		case []interface{}:
			fmt.Println(k, "is an array:")
			for i, u := range vv {
				fmt.Println("List", i, u)
			}
		case map[string]interface{}:
			fmt.Println(k, "is an array:")
			for i, u := range vv {
				fmt.Println("Map", i, u)
			}
		default:
			fmt.Println(k, "is of a type I don't know how to handle")
		}
	}
}
