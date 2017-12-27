package gorairpc_test

import (
	"github.com/donutloop/GoRaiRpc"
	"testing"
)

const accountAddress string = "xrb_3cwsajkzdycgg6k4q7zea68t4h5tu3ui59rdriurtjjem7rrjqtoixu473bt"
const accountKey string = "AB994465F5F94E71242B97EC410DA13C7AD877019F0BC4378D462C997188DF55"

func TestRaiRpc_AvailableSupply(t *testing.T) {
	endpoints := gorairpc.New()

	availableSupply, err := endpoints.AvailableSupply()
	if err != nil {
		t.Fatalf("error requesting available supply (%v)", err)
	}

	if availableSupply == nil {
		t.Error("available supply is zero")
	}
}

func TestRaiRpc_AccountBalance(t *testing.T) {
	endpoints := gorairpc.New()
	_, err := endpoints.AccountBalance(accountAddress)
	if err != nil {
		t.Fatalf("error requesting account balance (%v)", err)
	}
}

func TestRaiRpc_AccountHistory(t *testing.T) {
	endpoints := gorairpc.New()
	_, err := endpoints.AccountHistory(accountAddress, "4096")
	if err != nil {
		t.Fatalf("error requesting account history (%v)", err)
	}
}

func TestRaiRpc_AccountKey(t *testing.T) {
	endpoints := gorairpc.New()
	_, err := endpoints.AccountKey(accountAddress)
	if err != nil {
		t.Fatalf("error requesting account key (%v)", err)
	}
}

func TestRaiRpc_AccountGet(t *testing.T) {
	endpoints := gorairpc.New()
	_, err := endpoints.AccountGet(accountKey)
	if err != nil {
		t.Fatalf("error requesting account get (%v)", err)
	}
}

func TestRaiRpc_Version(t *testing.T) {
	endpoints := gorairpc.New()
	_, err := endpoints.Version()
	if err != nil {
		t.Fatalf("error requesting account get (%v)", err)
	}
}


