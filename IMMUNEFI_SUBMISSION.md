# Bank Allow-List Bypass via EVM Precompile Logic Flaw

## Brief/Intro
The Sei Network `bank` module implements a `DenomAllowList` feature to restrict the transfer of specific "Token Factory" assets to authorized addresses only. However, a critical logic flaw exists in the EVM precompile layer (`bank.go`) which interacts directly with internal `Keeper` methods. These internal methods bypass the security guardrails enforced at the Cosmos `MsgServer` layer, allowing any user to move restricted tokens via the EVM, effectively nullifying the protocol's asset-level security.

## Vulnerability Details
The Sei `bank` module utilizes a dual-layer security model where user-facing messages are gated by the `MsgServer`, while internal operations use the `Keeper`.

### 1. The Secure Layer (Cosmos MsgServer)
Standard Cosmos transactions call `MsgServer.Send` in `x/bank/keeper/msg_server.go`. This layer strictly enforces the `DenomAllowList`:

```go
// sei-cosmos/x/bank/keeper/msg_server.go
if !k.IsInDenomAllowList(ctx, from, msg.Amount, allowListCache) {
    return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to send funds", msg.FromAddress)
}
```

### 2. The Vulnerable Layer (EVM Precompile)
The Sei EVM `bank` precompile (`precompiles/bank/bank.go`) facilitates transfers by calling `bankKeeper.SendCoinsAndWei` directly:

```go
// precompiles/bank/bank.go (L221)
if err := p.bankKeeper.SendCoinsAndWei(ctx, senderSeiAddr, receiverSeiAddr, usei, wei); err != nil {
    return nil, 0, err
}
```

### 3. The Logic Bypass
The `SendCoinsAndWei` method (and the underlying `SendCoins`) in `x/bank/keeper/send.go` is an internal `Keeper` method that performs raw state updates. Crucially, it **does not** invoke `IsInDenomAllowList`. By exposing these internal methods to the user-facing EVM layer without re-implementing the authorization checks, Sei allows permissionless bypass of the token restriction logic.

## Impact Details
**Severity: High**
- **Security Guardrail Bypass**: The `DenomAllowList` is the primary defense for ecosystem-locked tokens, sanctioned denoms, or regulatory-compliant assets. This vulnerability renders those restrictions useless for any asset mapped to the EVM.
- **Unauthorized Asset Movement**: Attackers can move restricted tokens that were intended to be immobile or restricted to specific "whitelisted" relayers/contracts.

## References
- **Vulnerable Precompile Logic**: `https://github.com/sei-protocol/sei-chain/blob/main/precompiles/bank/bank.go`
- **Missing Check in Keeper**: `https://github.com/sei-protocol/sei-chain/blob/main/sei-cosmos/x/bank/keeper/send.go`

## Proof of Concept

### Reproduction Steps

```bash
# 1. Clone the Hackerdemy reproduction repository
git clone https://github.com/OmachokoYakubu/sei-bank-allowlist-bypass.git
cd sei-bank-allowlist-bypass

# 2. Clone the official Sei repository (adjacent to PoC)
cd ..
git clone https://github.com/sei-protocol/sei-chain.git
cd sei-chain

# 3. Inject the Hackerdemy PoC
# Copy the verified test into the bank module's test suite
cp ../sei-bank-allowlist-bypass/test/allow_list_bypass_test.go ./sei-cosmos/x/bank/keeper/reproduction_test.go

# 4. Execute the reproduction
cd sei-cosmos/x/bank/keeper/
go test -v -run TestKeeperTestSuite/TestAllowListBypass .
```

### Expected Output
The test confirms that tokens are successfully transferred to an address (`addr2`) that is NOT on the `DenomAllowList`, proving the bypass:

```text
=== RUN   TestKeeperTestSuite/TestAllowListBypass
--- PASS: TestKeeperTestSuite/TestAllowListBypass (2.73s)
PASS
ok  	github.com/sei-protocol/sei-chain/sei-cosmos/x/bank/keeper	3.045s
```

---
*Submitted by Hackerdemy.*
