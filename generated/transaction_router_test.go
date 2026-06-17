package comfactoryrouting

import (
	"testing"
)

func FuzzRouteTransaction(f *testing.F) {
	// Seed corpus: valid baseline
	f.Add("test.db", "tx-001", "alice", "bob", 500.0, "low")
	// Seed: empty receiver (targets coverage gap line 99-101)
	f.Add("test.db", "tx-002", "charlie", "", 200.0, "high")
	// Seed: empty priority (targets coverage gap line 102-104)
	f.Add("test.db", "tx-003", "dave", "eve", 100.0, "")
	// Seed: empty dbPath (targets coverage gap line 90-92)
	f.Add("", "tx-004", "frank", "grace", 50.0, "low")
	// Seed: path traversal in dbPath (targets coverage gap line 107-109)
	f.Add("../../etc/passwd", "tx-005", "henry", "ivy", 300.0, "medium")
	f.Add("/etc/shadow", "tx-006", "jack", "kate", 400.0, "high")
	// Seed: SQL injection patterns in string fields
	f.Add("test.db", "'; DROP TABLE routing_rules; --", "liam", "mia", 150.0, "low")
	f.Add("test.db", "noah", "' OR '1'='1", "olivia", 250.0, "high")
	f.Add("test.db", "peter", "quinn", "select * from users", 350.0, "medium")

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
