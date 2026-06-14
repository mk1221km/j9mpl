# Enterprise Transaction Routing Service Specification

## 1. Objective
Design a zero-dependency, data-oriented database transaction script in NetRexx that routes financial transactions (sender, receiver, amount, priority) to target processing channels based on SQLite routing rules, logs routing events, and compiles statistics.

## 2. Invariants & Data Layout
* **Database Engine:** SQLite
* **Target Table 1:** `routing_rules`
  * `min_amount` REAL (Primary Key)
  * `priority` TEXT (Primary Key)
  * `channel` TEXT
* **Target Table 2:** `transaction_log`
  * `tx_id` TEXT (Primary Key)
  * `sender` TEXT
  * `receiver` TEXT
  * `amount` REAL
  * `channel` TEXT
  * `status` TEXT

## 3. Data Transfer Objects
A private helper class `TransactionRecord` must represent individual data transactions:
* `txId` (String)
* `sender` (String)
* `receiver` (String)
* `amount` (Rexx)
* `priority` (String)

## 4. Required Interfaces
The primary public class `TransactionRouter` must implement:
1. `initRoutingTable(dbPath: String)`: Creates tables `routing_rules` and `transaction_log` if they do not exist, and populates default routing rules.
2. `routeTransaction(dbPath: String, record: TransactionRecord) returns String`: Determines the processing channel (e.g. "ACH", "WIRE", "SWIFT") based on `routing_rules`, inserts the logged transaction status, and returns the channel name.
3. `getTransactionCount(dbPath: String, status: String) returns Rexx`: Queries and returns the total transaction count for a given status.
4. `main(args: String[])`: Command-line driver that populates default rules, routes test transactions, and verifies log totals.
