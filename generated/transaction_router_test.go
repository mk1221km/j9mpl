package comfactoryrouting

import (
	"testing"
)

func FuzzRouteTransaction(f *testing.F) {
	// Baseline valid corpus seed
	f.Add("test.db", "tx-001", "alice", "bob", 500.0, "low")

	// Auto-generated seeds from type-driven boundary induction
	seeds := SeedCorpusData()
	type seedCase struct {
		name   string
		values []interface{}
	}
	var seedCases []seedCase
	for name, vals := range seeds {
		seedCases = append(seedCases, seedCase{name, vals})
	}

	// Map seed names to fuzz parameter positions
	// f.Add(dbPath, txId, sender, receiver string, amount float64, priority string)
	for _, sc := range seedCases {
		dbPath := "test.db"
		txId := "tx-auto"
		sender := "auto-sender"
		receiver := "auto-receiver"
		amount := 100.0
		priority := "normal"

		// Parse seed name to determine which field to substitute
		switch sc.name {
		case "txId_empty", "txId_path_traversal", "txId_sql_injection":
			if len(sc.values) > 0 {
				txId = sc.values[0].(string)
			}
		case "sender_path_traversal", "sender_sql_injection":
			if len(sc.values) > 0 {
				sender = sc.values[0].(string)
			}
		case "receiver_empty", "receiver_path_traversal", "receiver_sql_injection":
			if len(sc.values) > 0 {
				receiver = sc.values[0].(string)
			} else {
				receiver = ""
			}
		case "amount_negative", "amount_nan":
			if len(sc.values) > 0 {
				amount = sc.values[0].(float64)
			}
		case "priority_empty", "priority_path_traversal", "priority_sql_injection":
			if len(sc.values) > 0 {
				priority = sc.values[0].(string)
			} else {
				priority = ""
			}
		}
		f.Add(dbPath, txId, sender, receiver, amount, priority)
	}

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
