# PLAN.md — kinetix-notary-service

The agreed design of the notary: what it is for, every decision taken, the exact cryptographic
scheme, and the sequence of work from the foundation to zero-knowledge proofs. `CLAUDE.md` is the
short version — the rules and how they are enforced. When the two disagree, fix whichever is wrong;
do not let them drift.

## 1. What this service is

**The notary makes the platform's commitments publicly verifiable.** Kinetix holds other people's
money — buyer funds in escrow until the return window closes, merchant and driver balances, customer
wallets — and today every one of those numbers is visible only inside `kinetix-payment-service`.
The notary takes a consistent snapshot of what the platform owes, commits to it with a Merkle sum
tree, anchors the commitment on a public chain, and hands each account holder a proof that their
own balance is inside the committed total. Anyone can check that proof against the chain without
trusting Kinetix, and nobody — Kinetix included — can rewrite an anchored epoch afterwards.

It owns: epochs, the trees built for them, the salts that make each account's commitment private,
the anchoring transactions and their receipts, and the inclusion proofs it issues.

It never owns: balances (payment), orders (order), reviews (review), identities (identity). It
reads a snapshot from the owner over gRPC and commits to it. If a balance is wrong in payment, the
notary anchors the wrong balance faithfully; correctness of the ledger stays payment's job.

It never moves money on chain. Indonesian law (UU 7/2011 on Currency, PBI 22/23/PBI/2020) does not
allow crypto as a means of payment, so a payment rail would be a feature that may not run in
production. The notary uses a chain for what a chain is uniquely good at: a public, append-only,
timestamped record that no single party controls.

## 2. Decisions (taken 2026-10-11)

| # | Question | Decision |
|---|---|---|
| 1 | Which web3 service | The notary (proof of liabilities first). The product passport (anti-counterfeit) is the **next** service, after this one is complete. |
| 2 | Off-chain language | **Go**, enterprise layout, **`go.uber.org/dig`** for dependency injection. search moves to dig later as well. |
| 3 | Network | A **testnet, permanently**, with production-grade code: switching networks is RPC URL + chain id + keys + ledger address, nothing else. Chosen: **Base Sepolia** (chain id 84532). Gas must not need constant faucet trips — see §7. |
| 4 | Upgradeability | **Immutable contracts**, versioned by redeploying, with a hash-chained history that lets a successor continue an epoch chain exactly. Built for the long term — see §6. |
| 5 | Publish the total? | Scheme v1 publishes the total (the money is sandbox money, so nothing real is disclosed). **Scheme v2 hides it with a zero-knowledge proof** — §10 step 8. |
| 6 | Key custody | Whatever is least effort: **the existing SOPS store**, three separate keys — §8. |
| 7 | Phase 1 scope | **Proof of liabilities only.** Evidence anchoring (orders, proof of delivery, reviews) follows as phase 2. |
| 8 | Comments | **Zero comments** in Solidity and Go. The reasoning lives here and in `CLAUDE.md`. |
| 9 | Timing | Start now; the user releases the pending identity/order work independently. |
| — | Toolchain | **Docker only.** Nothing is installed on the laptop; `bin/dev <command>` runs inside a pinned toolchain image with a local chain beside it. |

## 3. What the proof does and does not prove

Written down because a proof of liabilities is easy to oversell, and overselling it is the estate's
house failure mode wearing a cryptographic costume.

It proves:

- Kinetix committed, publicly and at a recorded time, to one specific set of account balances and
  their exact total.
- A given account's balance is inside that committed set, at that amount.
- The commitment cannot be changed afterwards, and epochs cannot be silently dropped or reordered:
  each digest commits to the one before it.

It does **not** prove:

- **Solvency.** The money backing those liabilities sits in Midtrans and a bank, off chain. A proof
  of reserves would need those custodians to attest, which they do not. Liabilities without reserves
  shows what is owed, not that it can be paid.
- **That every account is included.** An account holder who never checks can be left out unnoticed.
  The defence is that anyone *can* check, at any time, against any epoch.
- **That payment's numbers are right.** The tree commits to whatever the snapshot says.

## 4. Scheme v1

Every hash is keccak256 over `abi.encodePacked` big-endian bytes, so the Solidity verifier and the
Go builder produce identical digests. Each preimage starts with a one-byte domain tag so that no
leaf can ever be mistaken for a node, a node for padding, or any of them for an epoch digest.

| Item | Preimage | Bytes |
|---|---|---|
| leaf | `0x00 ‖ epoch:u64 ‖ commitment:bytes32 ‖ amount:u64` | 49 |
| node | `0x01 ‖ left.hash ‖ left.sum:u64 ‖ right.hash ‖ right.sum:u64` | 81 |
| empty (padding) | `0x02` | 1 |
| epoch digest | `0x03 ‖ scheme:u16 ‖ number:u64 ‖ previousDigest ‖ root ‖ total:u64 ‖ leafCount:u64 ‖ asOf:u64` | 123 |

- **A node carries both children's sums, not just their total.** Committing only to `left + right`
  is the published flaw in Maxwell's original summation tree: a prover can shift value between
  siblings without changing the parent. Hashing both sums closes it.
- **Amounts are unsigned 64-bit minor units.** A negative balance cannot be expressed, so it cannot
  be used to cancel someone else's. A sum that does not fit is refused when building and rejected
  when verifying — never wrapped.
- **Depth is canonical**: `0` for zero or one leaf, otherwise `bitlen(leafCount − 1)`. The tree is
  padded to `2^depth` leaves with the empty node. A proof must carry exactly `depth` siblings and an
  index below `leafCount`, so padding can never be proven as an account. Maximum depth 32.
- **The leaf binds the epoch number**, so a proof cannot be replayed against another epoch.
- **The epoch digest excludes the chain id, the contract address and the anchoring time.** It is a
  pure function of the ledger's content, which is what lets the chain move to another network or
  another contract without breaking continuity (§6).
- **Verification checks the path sum before any hash.** The root sum is exactly `amount + Σ sibling
  sums`, so comparing that sum (in 256-bit arithmetic on chain, carry-checked in Go) with the
  anchored total first rejects cheap forgeries and guarantees the hashing loop can never overflow.
  This is also why an overflowing proof returns `false` instead of reverting.
- An account's commitment is `bytes32`, opaque to the tree. In step 3 it becomes
  `H(principal ‖ accountKind ‖ salt)` with a fresh salt per account per epoch, so the public tree
  reveals neither who holds an account nor whether two epochs contain the same holder.

The scheme is implemented twice — `contracts/src/SumTree.sol` + `EpochChain.sol`, and
`internal/notary` + `internal/keccak` — and the two are held together by three independent checks:
known-answer vectors computed with `cast` (tests/keccak), a differential suite where Foundry
re-verifies every tree, proof, rejection and digest the Go builder writes
(`contracts/test/fixtures/scheme-v1.json`), and the integration suite, where proofs built in Go are
verified by the deployed contract on a real chain.

## 5. The ledger contract

`LiabilityLedger` is append-only. Its whole state is a map of epochs, the head digest, the next
epoch number, the last snapshot time and a sealed flag.

`anchor(header)` — anchorer role only — accepts exactly the next epoch, whose `previousDigest` must
equal the head, whose snapshot time must be strictly after the previous one and not in the future,
whose leaf count fits the scheme, and which, when empty, carries the canonical empty root and a zero
total. Re-sending an anchored epoch reverts `EpochAlreadyAnchored` when identical (the off-chain
writer reads that as "already applied") and `EpochConflict` when different (that is an incident).

`verifyInclusion(...)` is a `view`, so checking a proof costs nobody anything. It returns `false`
for a proof that does not prove inclusion and **reverts `UnknownEpoch`** for an epoch that does not
exist. Those are different answers and must never be collapsed: "your proof is wrong" and "there is
nothing to check against" lead a person to very different conclusions.

`seal(expectedHead)` — admin only, one way — stops the ledger for ever and emits the head, the next
epoch number and the last snapshot time: exactly the three values a successor's constructor needs.

Roles use OpenZeppelin `AccessControlDefaultAdminRules`: a single admin, transferable only in two
steps with a delay (the deploy sets it; three days is the intended value), and an anchorer role the
admin grants and revokes. There is no pause: revoking the anchorer stops anchoring, and sealing ends
the ledger.

The contract announces `ILiabilityLedger` through ERC-165, so a verifier can refuse an address that
is not a ledger instead of reporting "epoch not found" from a contract that never had epochs.

Measured gas (Foundry gas report, 2026-10-11): `anchor` **≈144k**, `verifyInclusion` **≈12.8k**
(free through `eth_call`), `seal` ≈29k.

## 6. Built for the long term

The contracts never change after deployment. Everything that could change over decades is handled
by redeploying, and continuity is a property of the data, not of an address.

- **Successor ledgers.** Seal the old ledger at its head, then deploy a new one with
  `genesisEpoch = nextEpoch`, `genesisDigest = head`, `genesisAsOf = lastAsOf` (all three are in the
  `LedgerSealed` event). The successor accepts exactly the epoch that would have come next. The two
  contracts together are one unbroken chain.
- **Moving networks.** Base Sepolia depends on Sepolia, and testnets are retired — Goerli and Holesky
  both were. Because the digest excludes the chain, the same procedure moves the ledger to another
  network: seal, deploy the successor elsewhere, keep anchoring. Proofs for older epochs stay
  verifiable for as long as the old network serves reads; the notary keeps every anchored header and
  receipt in its own database so that history survives the network.
- **Scheme versions.** `scheme()` is part of the interface and of every digest. A new scheme is a new
  contract version and a new `notary.Scheme` implementation; the verifier refuses a ledger whose
  scheme it does not speak.
- **Key compromise.** Revoke the anchorer at once. Epochs anchored by the attacker stay in the old
  contract (nothing can delete them), so deploy a successor whose genesis is the **last known-good**
  digest and epoch, publish which epochs are disowned and why, and point the notary at the
  successor.
- **Pinned toolchain.** Solidity 0.8.37 exactly, EVM version `cancun` (accepted by every L2 that
  matters, now and after future forks), no metadata hash in the bytecode, CREATE2 with a configured
  salt — the same sources always produce the same contract at a predictable address.

## 7. Network and gas

Base Sepolia, chain id **84532**. Finality on an OP Stack chain follows L1 finality, roughly 15–20
minutes; the verifier reads at the `finalized` block and says `not_final` until then (§9).

Hourly epochs are ~8,760 anchors a year, ≈1.26 billion gas at the measured 144k. Testnet gas prices
on an L2 are fractions of a gwei, so a single faucet grant covers a long time. The exact cost is not
written here because it is not known — it is computed from real receipts once the anchorer runs.

What "no constant faucet trips" means in practice, honestly:

- **Faucets cannot be automated legitimately.** Every public faucet sits behind a captcha, an account
  or a mainnet-balance check, and scripting around those breaks their terms. That is not built.
- **The anchorer measures its own runway.** Step 4 exports the signer balance and the cost of every
  anchor from its receipt (L2 execution plus the L1 data fee), computes days of runway, and alerts
  through the estate's existing Discord chain well before the balance runs out — the same alerting
  path already proven for backups.
- **Refills are rare and one-off**: one faucet grant or a bridge from Sepolia, triggered by the alert.

## 8. Keys

Three keys, three jobs, all in the SOPS store (`kinetix-secrets set <env> <KEY>`), generated with
`cast wallet new` inside the toolchain:

| Key | Used by | Holds |
|---|---|---|
| deployer | `forge script` only, at deploy time | enough gas to deploy |
| admin | `forge script` / `cast` only: rotate the anchorer, seal | the admin role; never online |
| anchorer | the running notary | the anchorer role and the gas budget |

The anchorer is the only hot key, and it can do exactly one thing: append the next epoch. A leaked
anchorer cannot rewrite history, cannot seal, and cannot grant itself anything. A Safe multisig as
admin is the upgrade path once the notary handles anything beyond sandbox money.

## 9. Architecture

```
cmd/notary-verify/        the independent verifier CLI. composition root.
internal/notary/          THE CORE. standard library only.
internal/keccak/          scheme v1 hashing.
internal/verifying/       use case: Verifier.
internal/ethereum/        adapter → the chain (go-ethereum, abigen v2 bindings).
internal/prooffile/       adapter → the public proof document, kinetix.notary.inclusion.v1.
internal/config/          environment reading, fail-closed.
internal/verifycli/       the CLI's argument, output and exit-code contract.
internal/wiring/          the dig graph. the only place outside cmd/ that knows dig.
internal/contractgen/     abigen output, generated by bin/generate-bindings, never committed.
contracts/                Foundry: src/, test/, script/.
tests/                    every Go test, by package; tests/integration needs a node.
toolchain/                the pinned development image and its compose file.
```

Errors carry a `notary.Kind`, and each kind is a distinct answer to the person asking. The verifier
CLI reports them as:

| Verdict | Exit | Meaning |
|---|---|---|
| `included` | 0 | the chain and the local scheme both confirm the account |
| `rejected` | 1 | both say the proof does not prove inclusion; the reason says why |
| `not_anchored` | 2 | the ledger has no such epoch |
| `not_final` | 3 | the epoch exists but its block is not final yet |
| `unknown` | 4 | the chain could not be asked — never a verdict |
| `wrong_ledger` | 5 | wrong chain, no contract, not a ledger, another scheme, or a proof issued against another ledger |
| `disagreement` | 6 | the contract and the Go scheme disagree: a defect, loudly |
| `malformed_proof` / `no_proof` | 65 / 66 | the document is unusable |
| configuration | 78 | a required variable is missing or invalid |

The verifier checks the ledger's identity first (chain id, code at the address, ERC-165, scheme),
then reads the epoch at the finalized block, verifies locally, asks the contract the same question,
and only reports a verdict when both agree.

## 10. Sequence of work

### Step 1 — foundation. **Done 2026-10-11.**

- `LiabilityLedger` with `SumTree` and `EpochChain`; deploy script (CREATE2, every value from the
  environment, SafeCast on every narrowing).
- Foundry suite: 38 unit, 7 fuzz (1,024 runs locally, 20,000 in CI), 7 invariants over a handler with
  ghost state (131k calls per CI run), 6 differential tests against Go's fixtures, 4 deploy tests.
  forge-lint with warnings denied; gas snapshot checked in CI.
- Go: the pure core, scheme v1, the verifier use case, the chain adapter reading at the finalized
  block, the proof document, fail-closed config, the dig graph, the `notary-verify` CLI and image.
- Tests: core, scheme known-answers, fuzzing of every tamper, codec, config, CLI, the dig graph, and
  9 integration tests against a real chain. Go coverage 91.2% across `internal/`.
- `bin/check-boundaries`, golangci-lint, `bin/ci` as the single gate, GitHub and GitLab pipelines
  that both run it. **Neither pipeline has run yet.**

### Owed next to step 2 — CI/CD on GitHub

The user runs CI/CD on GitHub; GitLab carries the paired CI only (`verify`, `secret-scan`,
`build-image`) and never publishes or deploys. Step 1 shipped `ci.yml` alone. Still to write, with
search's three workflows as the model:

- **`build-and-push.yml`** — on every push to `main`: build the root `Dockerfile`, push
  `ghcr.io/wildanfrananda/kinetix-notary-service:<sha>`, gate on Trivy (HIGH/CRITICAL) as search
  does, and print the digest. Today the image carries `notary-verify`; from step 4 the same image
  also carries `notaryd`, one image for every binary as in search.
- **`release.yml`** — on a `vMAJOR.MINOR.PATCH` tag: promote the digest already built from `main`,
  never rebuild. It must **not** dispatch a deploy to `kinetix-infrastructure` until step 4 has put
  the notary into the production model; a dispatch for a service infra does not know fails, and a
  digest pinned before the image exists blocks every deploy in the estate (the notification lesson).
- **Static analysis for the contracts** — Slither and Aderyn as their own job, failing on findings,
  with any accepted finding recorded in `CLAUDE.md` rather than silenced in the source.
- **Branch protection** — a ruleset requiring `verify` (and the analysis job), created after the
  first pull request has shown those checks; the picker only offers checks it has already seen.

Contract deployment is deliberately **not** CD. The admin key never sits in a CI secret store;
deploys and seals run by hand through `bin/dev` with keys from SOPS (§8).

### Step 2 — Base Sepolia deployment

Generate the three keys, store them in SOPS, fund deployer and anchorer once. Run
`DeployLiabilityLedger` through `bin/dev`, verify the source on the chain's explorer, and commit the
broadcast record (`contracts/broadcast/…/84532/`). Record the address and deployment block in
`deploy/base-sepolia.json`. Nothing about the contract changes for this step.

### Step 3 — the snapshot contract (owner-side work)

A new RPC on payment, built **in payment** under the estate's rule that a service never takes another
service's job: a consistent snapshot of every liability (one `REPEATABLE READ` transaction, so the
total matches a real instant), streamed or paged, with an account kind per row (customer wallet,
merchant available, merchant pending, driver available, driver pending, escrow held, suspense — the
exact list is payment's to define). Published in `kinetix-contracts` as an additive `payment/v1`
change, plus a new `notary/v1` package. The notary's allow-list names payment as its only source.

### Step 4 — `notaryd`, the anchorer

Postgres of its own (goose migrations, sqlc), and the epoch as a persisted state machine:
`Planned → Built → Submitted → Mined → Finalized`, plus `Abandoned` needing a person — the checkout
saga's pattern: the row is written before the transaction is broadcast, so a crash leaves a record
rather than an orphaned transaction.

- One writer, by lease with fencing, because nonces are per sender and two writers corrupt them.
- EIP-1559 fees with a configured ceiling; a stuck transaction is replaced at the same nonce, never
  duplicated at a new one.
- An epoch counts as anchored only at `finalized`. A reorg that drops the transaction re-submits it;
  `EpochAlreadyAnchored` on re-submission means it already landed.
- Salts generated per account per epoch from `crypto/rand`, stored encrypted, never logged.
- Runway metrics and alerts (§7). Metrics conformance, gRPC reflection and health like every service;
  SPIFFE mTLS with an explicit allow-list.

### Step 5 — the proof API

`notary.v1.ProofService` over gRPC and `GET /api/v1/proofs/liabilities/{epoch}` through Kong. The
caller's principal comes from identity's token; a principal receives the proofs for its own accounts
and nothing else, and the response is exactly the `kinetix.notary.inclusion.v1` document the CLI
already verifies. Unknown, not final and unavailable stay distinct all the way to the HTTP status.

### Step 6 — verification without Kinetix

A browser verifier in the public docs: it reads the chain directly from the visitor's browser, so
checking a proof involves no Kinetix server at all. The CLI image is published for auditors.

### Step 7 — phase 2: evidence anchoring

`EvidenceRegistry`, batched Merkle roots of:

- **order lifecycle transitions**, from order's existing `OrdersChangedSince` feed;
- **proof of delivery**: the fulfillment contract carries only `pod_photo_url` today, and a URL's
  content can be replaced. A `pod_content_sha256` must be computed where the photo is captured and
  stored (courier app and matching), then anchored, so a return dispute can show the delivery
  evidence existed at time T and is unchanged;
- **review integrity**: review content hashes, so a review cannot be quietly edited or removed.
  review serves no gRPC today; it needs a server and a `ChangedSince` first.

### Step 8 — scheme v2: zero-knowledge proof of liabilities

The total stops being public, and the proof gets stronger.

- **gnark** (Consensys, Go) — the same language as this service; PLONK over BN254 with a KZG
  structured reference string from a public ceremony, so there is no per-circuit trusted setup.
- The circuit proves: every leaf amount is in `[0, 2^64)`; the tree's root is built from those
  leaves with **Poseidon2** (cheap inside a circuit, unlike keccak); the leaves sum to a total `T`;
  and the ledger stores only `C = Poseidon(T, blinding)`.
- gnark exports a Solidity verifier; it becomes `LiabilityLedgerV2` (scheme 2), tested in Foundry
  the same way v1 is, with differential fixtures from the Go prover.
- Account holders still verify inclusion with a Merkle path — Poseidon instead of keccak — in the
  browser (WASM) or the CLI.
- Beyond one circuit's capacity: recursive aggregation of per-batch proofs.
- Budget to measure before committing: proving time and memory per 2^16 and 2^20 accounts on the
  production VPS, which shares 16 GB with the whole estate.

### Step 9 — the next service: product passport

After the notary is complete: per-unit provenance for serialized goods (ERC-721), minted by a
verified merchant and transferred to the buyer on delivery. It needs serial-level inventory, which
warehouse does not have today (no serial column anywhere in its schema), and buyer wallets
(account abstraction), so it starts with that owner-side work.

## 11. Open questions

- **Licence.** Every Solidity file says `UNLICENSED`. A public verifier contract is usually MIT so
  others may read and reuse the verifier; the choice is the user's.
- **Publishing the verifier image** to GHCR, and under which name.
