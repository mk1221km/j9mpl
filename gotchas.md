# NetRexx / Tcl / SQLite Toolchain Gotchas & Technical Insights

This document catalogs critical edge cases and "gotchas" discovered during the design, synthesis, and testing phases of the autonomous software factory.

---

## 1. Tcl Command Substitution in String Literals

### The Gotcha
In Tcl, square brackets `[` and `]` inside double-quoted string literals trigger **command substitution** (e.g., executing the contents as a Tcl command). 
If you write:
```tcl
puts "[ERROR] Compilation failed."
```
Tcl will attempt to evaluate a command named `ERROR` and crash with:
`invalid command name "ERROR"`

### The Fix
Always escape square brackets with backslashes when they are part of log headers or literal strings inside double quotes:
```tcl
puts "\[ERROR\] Compilation failed."
```
Alternatively, use curly braces `{}` which prevent command substitution, though they also prevent variable interpolation.

---

## 2. NetRexx Class Visibility & File Name Constraint

### The Gotcha
NetRexx compiles down to standard Java class structures. As a result:
1. A single NetRexx source file can contain multiple classes, but **only one class can be declared `public`**, and its name must exactly match the source file name (e.g., `TransactionRouter.nrx` can only contain `class TransactionRouter public`).
2. Declaring helper DTO classes (e.g., `TransactionRecord`) as `public` in the same source file will trigger translator errors:
   `Error: Public class name must be the same as the program name`

### The Fix
- Declare secondary classes in the same file as `shared` (NetRexx's package-private visibility) or `private` (private to the file):
  ```rexx
  class TransactionRecord shared
  ```
- If a secondary class must be `public`, extract it into its own `.nrx` source file.

---

## 3. NetRexx Classpath Dependency Resolution

### The Gotcha
When translating a NetRexx source file that references a class defined in another file (such as `TransactionRouterTest.nrx` referencing `TransactionRecord`), the translator (`nrc`) needs the referenced class to be on the classpath. 
If it is missing, `nrc` will fail to recognize the type and report a misleading error:
`Error: The method 'TransactionRecord()' cannot be found in class 'TransactionRouterTest'`

### The Fix
Ensure the output directory (e.g., `bin`) containing the compiled `.class` files, along with `lib/NetRexxF.jar`, are explicitly appended to the classpath during translation:
```bash
./bin/nrc -cp bin:lib/NetRexxF.jar ...
```

---

## 4. SQLite JDBC Filename Generation & Path Fuzzing

### The Gotcha
When running property-based tests or fuzzing databases, injecting SQL payload boundaries (e.g., `' ; DROP TABLE system_metrics; --`) into SQLite database connection path parameters can lead to host pollution. 
The SQLite JDBC driver will not execute the injection payload on connection startup; instead, it treats the entire string as a literal file path and physically creates a file named `./'; DROP TABLE system_metrics; --` in the workspace directory.

### The Fix
Implement a **Parameter Bound Partitioning** scheme in test generators:
- Separate database URI paths into a safe `dbPathBounds` pool (e.g., `["generated/test.db", ":memory:", "null"]`).
- Keep aggressive SQL injection vectors confined strictly to text-field input variables (`stringBounds`), never routing them to file/database initialization parameters.

---

## 5. Tcl `exec` Standard Error Redirection

### The Gotcha
In Tcl, the `exec` command treats any output to the standard error stream (`stderr`) as an application execution failure, even if the executed process terminates with a successful exit code `0`.
If a sandboxed JVM fuzzer execution outputs expected database warning stack traces to stderr, Tcl will catch the call as a failure and halt the pipeline.

### The Fix
Redirect `stderr` to `stdout` within the Tcl `exec` invocation using `2>@1`. This ensures `catch` only intercepts failures based on the actual process exit code:
```tcl
set status [catch {exec {*}$cmd 2>@1} result]
```

