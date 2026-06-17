package comfactoryrouting

import (
	"fmt"
	"os"
)

func main() {
	dbPath := "./transactions.db"
	if len(os.Args) > 1 {
		dbPath = os.Args[1]
	}

	router := TransactionRouter{}

	if err := router.initRoutingTable(dbPath); err != nil {
		fmt.Fprintf(os.Stderr, "init error: %v\n", err)
		os.Exit(1)
	}

	testRecords := []TransactionRecord{
		{TxId: "tx001", Sender: "alice", Receiver: "bob", Amount: 100.0, Priority: "high"},
		{TxId: "tx002", Sender: "charlie", Receiver: "dave", Amount: 250.5, Priority: "normal"},
		{TxId: "tx003", Sender: "eve", Receiver: "frank", Amount: 75.0, Priority: "low"},
	}

	for _, rec := range testRecords {
		channel, err := router.routeTransaction(dbPath, rec)
		if err != nil {
			fmt.Fprintf(os.Stderr, "route error: %v\n", err)
			continue
		}
		fmt.Printf("%s -> %s\n", rec.TxId, channel)
	}

	total, err := router.getTransactionCount(dbPath, "done")
	if err != nil {
		fmt.Fprintf(os.Stderr, "count error: %v\n", err)
	}
	_ = total
}
