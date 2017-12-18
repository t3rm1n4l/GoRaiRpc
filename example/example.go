package main

import (
	"fmt"
	xrpc "github.com/devmahno/GoRaiRpc"
)

func main() {
	rpc := xrpc.NewRaiRpc()
	fmt.Printf("Begin\n")

	fmt.Printf("AvailableSupply: %s\n", rpc.RpcAvailableSupply("raw"))
	showData(rpc.RpcVersion())

	fmt.Printf("-\n")
	showData(rpc.RpcAccountInfo("xrb_3cwsajkzdycgg6k4q7zea68t4h5tu3ui59rdriurtjjem7rrjqtoixu473bt", "Mrai", true, true, true))

	fmt.Printf("-RpcAccountsPending\n")
	//accounts, count = '4096', threshold = 0, unit = 'raw', source = false
	showData(rpc.RpcAccountsPending([]string{"xrb_1111111111111111111111111111111111111111111111111117353trpda",
	"xrb_3t6k35gi95xu6tergt6p69ck76ogmitsa8mnijtpxm9fkcm736xtoncuohr3"}, "1", 1, "Mrai", true))

	fmt.Printf("-RpcPending\n")
	//account, count = '4096', threshold = 0, unit = 'raw', source = false
	showData(rpc.RpcPending("xrb_1111111111111111111111111111111111111111111111111117353trpda", "1", 1, "Mrai", true))

	fmt.Printf("-RpcBlock\n")
	fmt.Println(rpc.RpcBlock("3BDCF72E7662DB9E2F38B15AF8FD1B633232928ED323C44E83D63573889E9BAE"))

	fmt.Printf("-RpcBlocks\n")
	hashes := []string{"3BDCF72E7662DB9E2F38B15AF8FD1B633232928ED323C44E83D63573889E9BAE"}
	showData(rpc.RpcBlocks(hashes))

	fmt.Printf("-RpcRepresentatives\n")
	fmt.Println(rpc.RpcRepresentatives("Mrai", "1048576", "false"))

	fmt.Printf("-RpcUcheckedKeys\n")
	fmt.Println(rpc.RpcUcheckedKeys("03C9B2BDEAD1F1982BAF16A8E66CFCD493B464533151BC01BA27963B0AAF6D81", "1"))

	fmt.Printf("-RpcMraiToRaw\n")
	fmt.Printf("MraiToRaw: %v\n", rpc.RpcMraiToRaw("1"))

	fmt.Printf("-ToUnit\n")
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
