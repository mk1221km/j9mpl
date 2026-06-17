package comfactoryrouting

import (
	"testing"
)

func FuzzRouteTransaction(f *testing.F) {
	// Seed corpus: first entry uses a valid temp db path
	f.Add("test.db", "tx-001", "alice", "bob", 500.0, "low")
	f.Add("test.db", "'; DROP TABLE routing_rules; --", "", "receiver", -10.0, "")
	f.Add("test.db", "../../etc/passwd", "normal", "receiver", 0.0, "medium")
	f.Add("test.db", "", "\x00null", "receiver", 100.0, "high")

	f.Fuzz(func(t *testing.T, dbPath, txId, sender, receiver string, amount float64, priority string) {
		router := TransactionRouter{}

		if err := router.initRoutingTable(dbPath); err != nil {
			t.Skip()
		}

		rec := TransactionRecord{
			TxId:     txId,
			Sender:   sender,
			Receiver: receiver,
			Amount:   amount,
			Priority: priority,
		}

		_, err := router.routeTransaction(dbPath, rec)
		_ = err
	})
}
