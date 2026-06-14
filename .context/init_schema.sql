-- Project Relational Context Ledger Schema
-- Enforces logical constraints, type containment, and manual pages.

CREATE TABLE IF NOT EXISTS declarations (
    symbol_uri TEXT PRIMARY KEY, -- logical mapping (e.g., method://App/Telemetry/process)
    file_path TEXT,              -- physical filesystem file location
    start_line INTEGER,
    end_line INTEGER
);

CREATE TABLE IF NOT EXISTS containment (
    parent_uri TEXT,
    child_uri TEXT,
    PRIMARY KEY (parent_uri, child_uri)
);

CREATE TABLE IF NOT EXISTS symbol_uses (
    caller_uri TEXT,
    callee_uri TEXT,
    location_span TEXT
);

CREATE TABLE IF NOT EXISTS system_documentation (
    symbol_uri TEXT PRIMARY KEY,
    markdown_payload TEXT,       -- Inline specifications, requirements, and system man-pages
    invariants_json TEXT         -- Formally tracked architectural constraints
);
