package cxbenchmark

import (
	"fmt"
	"testing"
	"time"

	"github.com/mit-dci/opencx/benchclient"
)

// BenchmarkEasyAuctionPlaceOrders places orders on an unauthenticated auction server with "pinky swear" settlement and full matching capabilites.
func BenchmarkEasyAuctionPlaceOrders(b *testing.B) {
	var err error

	b.Logf("Test start -- time: %s", time.Now())

	var client1 *benchclient.BenchClient
	var client2 *benchclient.BenchClient
	var offChan chan bool
	if client1, client2, offChan, err = setupEasyAuctionBenchmarkDualClient(false); err != nil {
		b.Errorf("Something is wrong with test, stopping benchmark")
		offChan <- true
		b.Fatalf("Could not start dual client benchmark: \n%s", err)
	}

	b.Logf("Test started client")
	// ugh when we run this benchmark and the server is noise then we basically crash the rpc server... need to figure out how to have that not happen, why is that fatal?
	runs := []int{1}
	for _, varRuns := range runs {
		placeFillTitle := fmt.Sprintf("AuctionPlaceAndFill%d", varRuns)
		b.Logf("Running %s", placeFillTitle)
		b.Run(placeFillTitle, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				AuctionPlaceFillX(client1, client2, varRuns)
			}
		})
		placeBuySellTitle := fmt.Sprintf("AuctionPlaceBuyThenSell%d", varRuns)
		b.Logf("Running %s", placeBuySellTitle)
		b.Run(placeBuySellTitle, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				AuctionPlaceBuySellX(client1, varRuns)
			}
		})
	}

	b.Logf("stop benchmark rpc")
	offChan <- true

	return
}