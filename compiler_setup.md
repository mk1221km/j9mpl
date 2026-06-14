# Resilient Compiler & Automation Setup: NetRexx, ECJ, and OpenJ9

We have successfully configured and verified the **Eclipse Compiler for Java (ECJ)** alongside the **NetRexx 5.10 GA** distribution on top of the **IBM Semeru Runtime (OpenJ9)**.

Below is the updated system specification and guide for the autonomous compilation pipeline.

---

## 1. Architectural Details & Classpath Resolution

NetRexx compiles code via a two-stage pipeline:
1. **Translation:** Translates `.nrx` source files into Java code.
2. **Compilation:** Compiles the Java code in-memory into JVM bytecode using the Java Compiler API (JSR-199).

> [!IMPORTANT]
> **ECJ JSR-199 In-Memory Compatibility**
> A standard, unpatched modern ECJ jar (such as `ecj-3.46.0.jar`) on the classpath of the NetRexx translator causes JSR-199 in-memory compilation to fail with `File /<name>.java is missing (javac failed)`.
> This is because standard ECJ fails to correctly resolve the custom, in-memory URI scheme representing translated Java source code. NetRexx provides a custom-patched ECJ version (`lib/ecj-I20201218-1800-NRX-18.jar`) packaged inside `lib/NetRexxF.jar` to solve this.

To achieve maximum reliability:
1. **NetRexx Compiler Scripts** ([NetRexxC.sh](file:///home/me/code/j9mpl/bin/NetRexxC.sh), [nrc](file:///home/me/code/j9mpl/bin/nrc), etc.) are configured to use the bundled, custom-patched ECJ for translating and compiling NetRexx files.
2. **Standalone Java Compilations** are driven via a new custom wrapper script ([ecj](file:///home/me/code/j9mpl/bin/ecj)) that invokes the modern, headless **ECJ 3.46.0** directly on top of the OpenJ9 JRE.

---

## 2. Configured Compiler Utilities

### 1. Standalone Eclipse Java Compiler ([bin/ecj](file:///home/me/code/j9mpl/bin/ecj))
We created a custom executable wrapper `/home/me/code/j9mpl/bin/ecj` mapping to the latest **ECJ 3.46.0** downloaded from Maven Central:

```bash
#!/bin/sh
thisdir=$(dirname "$0")
java -jar "$thisdir/../ecj-3.46.0.jar" "$@"
```

#### Example Standalone Compile:
To compile standard Java source files with error tolerance ("Proceed on Errors"):
```bash
./bin/ecj -proceedOnError MyClass.java
```

---

### 2. NetRexx Command-Line Compiler ([bin/nrc](file:///home/me/code/j9mpl/bin/nrc))
We modified the NetRexx wrapper scripts to improve path robustness and default to the Eclipse Compiler:

* **Default Compiler:** Set `-Dnrx.compiler=ecj` system property by default in [NetRexxC.sh](file:///home/me/code/j9mpl/bin/NetRexxC.sh).
* **Relative Invocation:** Fixed [nrc](file:///home/me/code/j9mpl/bin/nrc) and [nr](file:///home/me/code/j9mpl/bin/nr) to resolve paths relative to their own script location, allowing them to be run from any directory.

---

## 3. Relational Context Ledger (SQLite)

We initialized the SQLite schema for the codebase retrieval indexing pipeline in [project_context.db](file:///home/me/code/j9mpl/.context/project_context.db) (schema defined in [init_schema.sql](file:///home/me/code/j9mpl/.context/init_schema.sql)):

* **Declarations Table:** Maps logical URIs to filesystem paths and line boundaries.
* **Containment Table:** Maps parent-to-child containment hierarchies.
* **Symbol Uses Table:** Tracks caller-to-callee dependencies.
* **System Documentation Table:** Stores markdown-formatted system man-pages and architectural invariants.

---

## 4. Automated Synthesis Loop ([factory-loop.sh](file:///home/me/code/j9mpl/factory-loop.sh))

The automated orchestration script [factory-loop.sh](file:///home/me/code/j9mpl/factory-loop.sh) executes the full translation, compilation, and validation cycle:

1. **Workspace Monitoring:** Stands by for incoming intent changes.
2. **Translation:** Uses the NetRexx compiler with `-nocompile -keepasjava -sourcedir -replace -format` flags to generate clean, standard `.java` files from NetRexx sources.
3. **Compilation:** Uses the standalone `ecj` tool with `-cp "${PROJECT_DIR}/lib/NetRexxF.jar"` and `-proceedOnError` flags to compile generated Java code while resolving NetRexx runtime types (such as `netrexx.lang.Rexx`).
4. **Validation:** Confirms the existence of compiled bytecode inside the `bin/` directory, ready to run on OpenJ9.

#### To Execute the Loop:
```bash
./factory-loop.sh
```

---

## 5. Relational RAG: Rascal Context Extractor ([src/ContextExtractor.rsc](file:///home/me/code/j9mpl/src/ContextExtractor.rsc))

We implemented a production-grade Rascal Meta-Programming Language (MPL) script at [src/ContextExtractor.rsc](file:///home/me/code/j9mpl/src/ContextExtractor.rsc) to perform static code analysis and populate the SQLite ledger database:

* **M3 Model Generation:** Generates an M3 model from the target Java source directories.
* **Declarations Extraction:** Maps logical URIs to physical files and line coordinates.
* **Containment Hierarchy:** Maps parent structures to child structures.
* **Symbol References Resolution:** Resolves code references (uses) and maps them back to their enclosing logical caller context by checking physical span overlaps.
* **Structured Output:** Emits the extracted facts as a highly optimized bulk transaction SQL file `.context/extracted.sql`.

#### To Execute the Extractor:
First, download JDT dependencies via Maven (runs against the unblocked HTTP repository configured in [pom.xml](file:///home/me/code/j9mpl/pom.xml)):
```bash
mvn dependency:copy-dependencies -DoutputDirectory=target/dependency/
```

Run the extractor module:
```bash
java -cp "rascal-shell-stable.jar:target/dependency/*" org.rascalmpl.shell.RascalShell ContextExtractor /home/me/code/j9mpl /home/me/code/j9mpl/.context/extracted.sql
```

Load the data into the context ledger:
```bash
sqlite3 .context/project_context.db < .context/extracted.sql
```

---

## 6. The Go Context Resolver ([bin/context_resolver.go](file:///home/me/code/j9mpl/bin/context_resolver.go))

To prevent **Syntax Contamination** and **Cognitive Translation Drag** (where Java syntax leaks back into synthesized NetRexx blocks) without adding Python runtime overhead, the workspace uses a compiled **Go Context Resolver**:

1. **Structural Analysis**: Rascal parses the compiled Java files to populate the logical relations in SQLite (`project_context.db`).
2. **Context Delivery**: The Go-based orchestration engine resolves the logical symbol URI using the code inside [bin/context_resolver.go](file:///home/me/code/j9mpl/bin/context_resolver.go). It maps the intermediate Java path back to the source `.nrx` file and performs a straight-line keyword block scan to extract pure NetRexx code.

#### Compilation:
To compile the resolver into a native static binary:
```bash
go build -o bin/context_resolver bin/context_resolver.go
```

#### Example Usage:
To retrieve the clean NetRexx source for a method context:
```bash
./bin/context_resolver "|java+method:///com/factory/telemetry/TelemetryEngine/processData(com.sun.net.httpserver.HttpExchange,com.factory.telemetry.TelemetryRecord)|"
```

---

## 7. Environment Summary

* **OS:** Linux
* **JVM Platform:** IBM Semeru Runtime Open Edition 26.0.1.0 (Eclipse OpenJ9 VM)
* **NetRexx Version:** 5.10-GA (Build 18-20260320)
* **Headless ECJ Version:** 3.46.0 (v20260528-0407)
