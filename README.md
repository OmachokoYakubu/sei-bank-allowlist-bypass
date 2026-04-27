# Sei Bank Precompile: Allow-List Bypass PoC

## Overview
This repository contains the official Proof of Concept (PoC) for a high-severity security guardrail bypass in the Sei Network `bank` precompile. The flaw allows unauthorized transfer of restricted tokens by circumventing the `DenomAllowList` enforcement.

## Project Structure
- `IMMUNEFI_SUBMISSION.md`: Formal bug report.
- `test/allow_list_bypass_test.go`: Go-based reproduction script.
- `exploit_results_final.txt`: Verified test output proving the bypass.

## Reproduction Steps

### 1. Clone this Repository
Clone the Hackerdemy reproduction package.
```bash
git clone https://github.com/OmachokoYakubu/sei-bank-allowlist-bypass.git
cd sei-bank-allowlist-bypass
```

### 2. Prepare the Target Environment
Clone the official Sei repository.
```bash
cd ..
git clone https://github.com/sei-protocol/sei-chain.git
cd sei-chain
```

### 3. Inject the PoC
Inject the Hackerdemy reproduction script into the local bank keeper directory.
```bash
cp ../sei-bank-allowlist-bypass/test/allow_list_bypass_test.go ./sei-cosmos/x/bank/keeper/reproduction_test.go
```

### 4. Run the Test
Execute the test suite in the target directory.
```bash
cd sei-cosmos/x/bank/keeper/
go test -v -run TestKeeperTestSuite/TestAllowListBypass .
```

### 5. Verify Results
The test will pass, confirming that restricted tokens were moved to an unauthorized address:
```text
=== RUN   TestKeeperTestSuite/TestAllowListBypass
--- PASS: TestKeeperTestSuite/TestAllowListBypass (2.73s)
PASS
```

---
*Developed by Hackerdemy.*
