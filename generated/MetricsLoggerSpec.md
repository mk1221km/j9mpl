# System Metrics SQLite Logger Specification

## 1. Objective
Design a zero-dependency, data-oriented database transaction script in NetRexx that logs system metrics (timestamp, name, and value) directly to an embedded SQLite database and computes metric averages.

## 2. Invariants & Data Layout
* **Database Engine:** SQLite
* **Target Table:** `system_metrics`
* **Table Schema:**
  * `timestamp` TEXT
  * `name` TEXT
  * `value` REAL

## 3. Data Transfer Objects
A private helper class `MetricRecord` must represent individual data transactions:
* `timestamp` (String)
* `metricName` (String)
* `metricValue` (Rexx)

## 4. Required Interfaces
The primary public class `MetricsLogger` must implement:
1. `initDatabase(dbPath: String)`: Runs `CREATE TABLE IF NOT EXISTS system_metrics` on the target database.
2. `logMetric(dbPath: String, record: MetricRecord)`: Inserts a single metric record.
3. `getAverageMetric(dbPath: String, name: String) returns Rexx`: Queries and returns the mathematical average value for the named metric.
4. `main(args: String[])`: Command-line interface driver that runs basic verification (init database, log test points, print average CPU usage).
