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
| 5 | Publish the total? | Scheme v1 publishes the total (the money is sandbox money, so nothing real is disclosed). **Scheme v2 hides it with a zero-knowledge proof** — §10 step 13. |
| 6 | Key custody | Whatever is least effort: **the existing SOPS store**, three separate keys — §8. |
| 7 | Phase 1 scope | **Proof of liabilities only.** Evidence anchoring (orders, proof of delivery, reviews) follows as phase 2. |
| 8 | Comments | **Zero comments** in Solidity and Go. The reasoning lives here and in `CLAUDE.md`. |
| 9 | Timing | Start now; the user releases the pending identity/order work independently. |
| — | Toolchain | **Docker only.** Nothing is installed on the laptop; `bin/dev <command>` runs inside a pinned toolchain image with a local chain beside it. |
| 10 | Where the code lives (2026-10-12) | **On chain first.** The user wants most of this repository to be Solidity. The rule that serves it without padding: whatever a stranger must be able to check or enforce without trusting Kinetix lives in a contract; whatever is expensive per account or private (building trees, salts, pulling snapshots) stays in Go. Steps 6–10 move accountability, quorum, discovery and gas sponsorship on chain. |

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
  The defence is that anyone *can* check, at any time, against any epoch — and steps 6 and 7 make
  a lie or a silence about it provable on chain.
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
- An account's commitment is `bytes32`, opaque to the tree. From step 4 it is
  `H(holder ‖ accountKind ‖ salt)` with a fresh salt per account per epoch, so the public tree
  reveals neither who holds an account nor whether two epochs contain the same holder. The holder is a
  principal, or the suspense purpose for money that is owed to nobody yet (step 3).

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
  matters, now and after future forks), no metadata hash in the bytecode, and deployment through the
  canonical CREATE2 factory with a configured salt (`always_use_create_2_factory`), so a successor's
  address is known before it exists. The first ledger is the exception — see step 2.

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

Three keys, three jobs. Corrected on 2026-10-12, before any key existed: the first version put all
three in the SOPS store, but `bin/prod-deploy` runs `kinetix-secrets render prod`, which decrypts the
**whole** prod store into the server's `.env` on every deploy. Anything in that store is online. So
only the key that is meant to be online lives there:

| Key | Lives in | Used by | Holds |
|---|---|---|---|
| anchorer | SOPS `prod` as `KINETIX_NOTARY_ANCHORER_PRIVATE_KEY` (`bin/new-anchorer-key`) | the running notary, from step 4 | the anchorer role and the gas budget |
| deployer | encrypted Foundry keystore `notary-deployer` on the laptop (`bin/new-keystore`) | `bin/deploy-ledger` only | enough gas to deploy |
| admin | encrypted Foundry keystore `notary-admin` on the laptop | rotating the anchorer, sealing | the admin role; never online |

The keystores sit in `KINETIX_NOTARY_KEYSTORE_DIR`, a directory the scripts refuse when it is inside
any git work tree — the same arrangement as the production CA keys in `pki-prod/`. Each is
scrypt-encrypted with its own password, typed at a prompt and never passed as an argument. **The
directory and both passwords must be backed up outside the laptop**: losing the admin keystore means
the anchorer can never be rotated and the ledger can never be sealed.

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

### Owed next to step 2 — CI/CD on GitHub, as strict as every other service

The user runs CI/CD on GitHub and asked for it to be **exactly as strict as the other services** —
"jangan lupa itu". GitLab carries the paired CI only (`verify`, `secret-scan`, `build-image`) and never
publishes or deploys. Step 1 shipped `ci.yml` alone; every other service repository has three
workflows and two rulesets. The bar, taken from search and the estate's
`kinetix-infrastructure/docs/BRANCH-PROTECTION.md` on 2026-10-12:

- [x] **`ci.yml`** — `verify` runs `bin/ci` on every pull request and every push to `main`. First run
  green on 2026-10-11.
- [x] **`build-and-push.yml`** (written 2026-10-12; its first run is the next merge to `main`) — on every push to `main` and on `workflow_dispatch`, one run at a time
  per ref: build the root `Dockerfile`; **Trivy gate on `CRITICAL`, `ignore-unfixed`, exit code 1, before
  anything is pushed**; an SPDX SBOM uploaded as an artifact; push
  `ghcr.io/wildanfrananda/kinetix-notary-service:<sha>` only after the scan passes; print the digest.
  Permissions `contents: read`, `packages: write`, nothing more. The image carries `notary-verify` now
  and `notaryd` from step 4, one image for every binary as in search.
- [x] **`release.yml`** (written 2026-10-12; first run on the first tag) — on `vMAJOR.MINOR.PATCH` tags only: refuse any other tag shape; refuse a tag
  not reachable from `main`; refuse a tag whose commit has no image already built and scanned —
  **promote, never rebuild**; point the release tag at that digest. The deploy dispatch to
  `kinetix-infrastructure` stays **out** until step 4 has put the notary into the production model:
  a dispatch for a service infra does not know fails, and a digest pinned before the image exists
  blocks every deploy in the estate (the notification lesson).
- [x] **Contract analysis** — Slither 0.11.6 (hash-locked, `toolchain/slither-requirements.txt`) and
  Aderyn 0.6.8 (SHA-256 pinned) run as `bin/analyze`, a step of `bin/ci`, so `verify` covers them and
  the ruleset needs no second check. Any new finding fails; the accepted ones are per detector and
  file in `contracts/analysis-accepted.json` and `contracts/slither.config.json`, with reasons in
  `CLAUDE.md`. Proven red twice: a planted `tx.origin` contract (Slither, exit 255) and a removed
  exception (Aderyn, `NEW`, exit 1). Stricter than the other services because a deployed contract
  cannot be patched.
- [x] **Ruleset `main`** — active; target the default branch; restrict deletions; block force pushes;
  require a pull request with **0** approvals; require the status check `verify` (and the analysis
  job once it exists); "up to date before merging" left off. Same as BRANCH-PROTECTION.md §2.
- [x] **Ruleset `release-tags`** — active; target **all tags**; restrict deletions and block force
  pushes, nothing else. A released tag is never moved or deleted; a bad release gets a patch tag.
  Same as §7. Both rulesets confirmed active through the API on 2026-10-12.
- [x] **`kinetix-infrastructure/docs/BRANCH-PROTECTION.md`** — the notary's row in §3 (`verify`) and
  its name in the §7 list.

The image passed the same Trivy gate locally before the workflow ever ran: the pinned
`debian:bookworm-slim` carried three fixable CRITICALs (`libpcre2-8-0`, `perl-base`) — even its newest
digest did — so the final stage runs `apt-get upgrade`, and Go moved to 1.26.9 for three fixable
HIGHs in the standard library. Result: no fixable CRITICAL or HIGH in either layer. search pins the
same Debian digest without an upgrade, so its next `build-and-push` is expected to fail the same gate.

Contract deployment is deliberately **not** CD. The admin key never sits in a CI secret store;
deploys and seals run by hand through `bin/dev` with keys from SOPS (§8).

### Step 2 — Base Sepolia deployment. **Done 2026-10-12.**

Tooling built and rehearsed end to end on anvil on 2026-10-12 (dry run, broadcast, inspection, a
second deploy refused, a wrong chain id caught). The operator's sequence:

1. `bin/new-keystore notary-deployer` and `bin/new-keystore notary-admin` — password prompts; each
   prints only its address.
2. `bin/new-anchorer-key prod` — generates the key inside the toolchain, checks it is exactly one
   32-byte key, and pipes it into `kinetix-secrets set prod`; prints only the address.
3. Both public addresses go into `deploy/base-sepolia.env`.
4. Fund the deployer once from a Base Sepolia faucet. Measured: the deploy uses **2,314,843 gas** and
   11,098 bytes of calldata; at the 0.006 gwei seen on 2026-10-12 that is about 0.000014 ETH plus the
   L1 data fee, so 0.01 ETH is ample.
5. `bin/deploy-ledger base-sepolia` — a dry run against the real chain: chain id checked, the
   contract simulated, nothing sent.
6. `bin/deploy-ledger base-sepolia --broadcast` — deploys through CREATE2 and verifies the source on
   Blockscout (no API key needed).
7. The ledger's address goes into `deploy/base-sepolia.env`; `bin/inspect-ledger base-sepolia` must
   report every line `ok`.
8. Commit `deploy/base-sepolia.env` and `contracts/broadcast/DeployLiabilityLedger.s.sol/84532/`;
   the store change in `kinetix-infrastructure` goes through its own pull request.

The result, every line checked by `bin/inspect-ledger` against the chain:

| | |
|---|---|
| ledger | `0x4E869A55031DE3B0d030AEBECc66E523686761a2`, Base Sepolia (84532) |
| deployed | block 48005190, tx `0x79e3a73fff93e5ce23a751d9c202fea44abced2c6d47cb5454e98ff34df14231` |
| cost | 2,309,865 gas × 0.006 gwei = 0.00001385919 ETH |
| admin | `0x2BfBa1D3d8595ad7d1E10a71BB4ec4Ac3671748E`, transfer delay 259200 s |
| anchorer | `0x51B2941aE55FC91D32D5595D7241612Bc3a166cc` — unfunded; fund it before step 4 |
| genesis | epoch 0, digest zero, no epochs anchored |
| source | verified on **Sourcify**: `match` for creation and runtime bytecode |

Three things went differently from the plan, and all three are now handled:

- **This ledger was deployed by plain CREATE, not CREATE2.** Foundry 1.8 sends `new X{salt: …}` in a
  broadcast as an ordinary contract creation unless `always_use_create_2_factory` is set, so the
  configured salt was not used and the address follows from the deployer and its nonce. The ledger is
  correct either way; only the predictability claim did not hold for it. The setting is now on, and a
  rehearsal on anvil showed the transaction going to the factory and the address equal to the one
  computed from the salt. The deploy test passed throughout — inside `forge test` a salted `new` is
  always CREATE2, which is exactly why it could not see this.
- **Blockscout cannot verify it yet**: its newest Solidity compiler is 0.8.36 and the ledger is built
  with 0.8.37. forge printed "All (1) contracts were verified!" over Blockscout's own
  `Fail - Unable to verify`. Verification moved to Sourcify, which needs no key; `bin/inspect-ledger`
  now asks Sourcify directly instead of trusting forge's summary. Submitting to Blockscout as well,
  once it lists 0.8.37, is optional.
- **The store write left the prod secrets in plaintext for about a minute.** Contained and fixed in
  `kinetix-infrastructure`; see its `docs/SECRETS.md` rule 6.

### Step 3 — the snapshot contract (owner-side work). **Built 2026-10-12; waiting on release.**

Built **in payment**, under the estate's rule that a service never takes another service's job. The
notary never reads payment's database: this RPC is the whole interface. The full reasoning and evidence
are in `kinetix-payment-service/CLAUDE.md`, "The liability snapshot, for the notary".

| Piece | Where | State |
|---|---|---|
| `payment.v1.LiabilityService/StreamLiabilitySnapshot` | `kinetix-contracts/proto/payment/v1/liability.proto`, a new file | `buf lint`, `format` and `breaking` exit 0; moneylint shown red on a `double` balance |
| the snapshot | payment: `LiabilitySnapshotService`, `PostgresLiabilitySnapshot`, `LiabilityGrpcServerService` | 231 tests, 0 skipped, against Postgres 16 |
| caller scoping | payment's `PeerAuthorizationInterceptor` now reads `caller@package.Service` | proven over real mTLS with a throwaway CA |
| allow-list | `kinetix-infrastructure/compose.yaml`: `order,notary@payment.v1.LiabilityService` | |

**The shape a reader must hold to** (step 4 builds this side):

- The stream is one header (`as_of`, currency), then batches of up to 500 accounts, then one trailer
  (`account_count`, `total`). A snapshot is complete only when the trailer arrives **and** the call ends
  `OK`. Anything else is no snapshot: never anchor a prefix. The reader recounts and re-sums the accounts
  itself and refuses a mismatch.
- An account is a `LiabilityKind` and a holder: a `principal_id`, or a `suspense_purpose` for money held
  for nobody yet (`UNASSIGNED_DRIVER_SHIPPING_FEE`). The kinds are customer balance, merchant
  available/pending escrow, driver available/pending escrow and suspense available/pending escrow, with
  zero balances included. **Escrow holds are not in it, on purpose:** the money of a `HELD` order already
  sits in the merchant's and driver's pending balances, and counting the hold too would double it.
- `as_of` is the instant the database took the snapshot (`REPEATABLE READ`, read-only, verified by the
  adapter). It becomes the epoch's `asOf`, which the ledger requires to strictly increase.
- Errors: `INVALID_ARGUMENT` (a currency that is not ISO 4217), `FAILED_PRECONDITION
  LIABILITY_SNAPSHOT_REFUSED` (payment's own data is inconsistent: an incident, alert on it),
  `UNAVAILABLE LIABILITY_SNAPSHOT_UNAVAILABLE` (retry), `DEADLINE_EXCEEDED` (the reader stopped reading
  for 30 s).

**Measured on 2026-10-12.** Under load from concurrent escrow traffic, `READ COMMITTED` let 295 of 1024
snapshots see money in flight; `REPEATABLE READ` let none. Over real mTLS, a notary leaf streamed 2,232
accounts in 7 messages with a matching count and sum, and was refused by `PaymentService`; a catalog leaf
was refused by `LiabilityService`.

**Release order:**

1. ✅ `kinetix-contracts` **v1.0.23** = `dfbe0221f783070dd65fdd868a15108d97db59cc`, published by the user
   2026-10-12. The tag adds exactly `liability.proto`, byte-identical to the file tested; npm, PyPI,
   RubyGems and Packagist all serve 1.0.23.
2. ✅ payment's `scripts/sync-contracts` pinned to v1.0.23 + that commit, `.contracts` re-fetched from the
   remote, and the CI sequence re-run on it (fresh Postgres, migrate profile, 231 tests, 0 skipped).
   Releasing payment is the user's.
3. Release `kinetix-infrastructure` with the allow-list. The order of 2 and 3 does not matter: an older
   payment treats `notary@…` as a caller that does not exist.

The `notary/v1` package moved to step 5, where its first message is written: an empty package is a
placeholder, and the contracts repository does not keep placeholders.

**Deferred, on purpose.** Scoping `order` to `order@payment.v1.PaymentService,order@grpc.reflection.v1.ServerReflection`
was run end to end and works, but a payment image older than this change would refuse every escrow call
from order under that value, so it waits until no such image can be rolled back to.

### Step 4 — `notaryd`, the anchorer

The daemon that turns payment's snapshot into an anchored epoch, once an hour, and can prove afterwards
exactly what it committed to. Designed on 2026-10-12, after step 3 fixed the shape of the snapshot.

**Processes.** `cmd/notaryd` (the daemon: anchoring loop, HTTP health and metrics) and
`cmd/notary-migrate` (goose, migrations embedded, a one-shot like every estate migrator). No gRPC server
yet: the first gRPC surface is `notary.v1.ProofService` in step 5, and reflection and gRPC health arrive
with it.

**One epoch in flight.** The next epoch is built only after the previous one is final. That keeps one
nonce in play at a time, makes `previousDigest` simply the ledger's head, and costs nothing at an
hourly cadence.

**The epoch, persisted.** `built → submitted → mined → finalized`, plus `abandoned`, which needs a
person. There is no `planned` state: pulling the snapshot and building the tree happen in one tick and
leave nothing worth recording when they fail; the next tick starts again. A `built` row is written in
one transaction with every leaf, so a header never exists without the accounts it commits to.

| State | Means | Leaves it when |
|---|---|---|
| `built` | snapshot read, salts drawn, tree built, header and digest stored | an attempt is written and broadcast |
| `submitted` | at least one signed attempt is broadcast | the ledger at `latest` holds the epoch with our digest |
| `mined` | it is on chain, not yet final | the ledger at `finalized` holds it (→ `finalized`), a reorg removed it (→ `submitted`), or the chain finalized a different digest (→ `abandoned`) |
| `finalized` | the ledger at the `finalized` block holds our digest | never |
| `abandoned` | the chain refused it for a reason no retry fixes | a person decides |

**Finality is a fact about the ledger, not about a receipt.** An epoch is final when `epochAt(n)`, read
at the `finalized` block, returns our digest. Which attempt landed does not matter; that is what makes
replacements, reorgs and re-broadcasts safe to reason about, and it is the same rule the verifier uses.

**Attempts.** Each signed transaction is written to the database (raw bytes, nonce, fees) **before** it
is broadcast, so a crash leaves a record and the next writer sees the pending nonce instead of
reusing it. Gas is `eth_estimateGas` + 20%, and the estimate doubles as a dry run: a revert there is
decoded before any gas is spent. EIP-1559: tip from `eth_maxPriorityFeePerGas`, fee cap
`2 × baseFee + tip`, both under `KINETIX_NOTARY_MAX_FEE_PER_GAS_WEI` (required, no default). An attempt
still unmined after five minutes is replaced at the **same nonce** with both fees raised 12.5% (nodes
require 10%); at the ceiling it is re-broadcast unchanged and counted, never raised past it.

**What a revert means.** `EpochAlreadyAnchored` with our digest: it landed, carry on watching.
`SnapshotInFuture`: the chain's clock is behind ours, try again later. `EpochConflict`,
`PreviousDigestMismatch`, `EpochOutOfSequence`, `LedgerIsSealed`, or a missing anchorer role:
`abandoned`, because another writer, a sealed ledger or a revoked key is an incident, not a retry.

**One writer, by lease with fencing.** A single-row lease (`holder`, `token`, `expires_at`; 30 s, renewed
every 10 s). Every write opens its transaction by re-reading the lease `FOR SHARE` and aborts unless this
process still holds the token it was given, so a writer that lost the lease mid-tick cannot record or
broadcast anything after losing it.

**The database must agree with the chain before anything is built.** The last finalized epoch in the
database must be the ledger's `nextEpoch − 1` with the ledger's `head`; for the very first epoch, the
database is empty and the ledger is at its genesis. A ledger holding epochs this database never built
stops the daemon: someone else holds the anchorer key.

**Accounts and commitments.** An account is a holder (a principal id, or the suspense purpose for
money owed to nobody yet), one of payment's seven kinds, and an amount in minor units. Its commitment is
`keccak256(0x04 ‖ kind:u8 ‖ salt:bytes32 ‖ holder)`: tag 0x04 is the next free domain tag, the holder
comes last because it is the only variable-length field, and it is capped at 255 bytes. Kinds keep
payment's numbers (1 customer balance … 7 suspense pending escrow). Leaves are sorted by commitment,
which is a fresh random permutation every epoch: a proof reveals its siblings' sums, and payment's
row order (customers first, by id) must not leak through them.

**Salts.** 32 bytes per account per epoch from `crypto/rand`, never logged. **Stored as they are, not
encrypted** — a change from the first version of this plan, for a reason: a leaf row has to hold the
holder and the amount to serve proofs later, so a column-encrypted salt beside a plaintext holder,
amount and commitment protects nothing the row does not already give away. The protection is the same
as payment's: a database of its own, reachable by this service alone, and backups encrypted with age
before they leave the server.

**Snapshot client.** mTLS to `payment.v1.LiabilityService` with the estate's SPIFFE leaf, the contract
of step 3 enforced to the letter: one header first, batches, one trailer last, then `OK`; the client
recounts and re-sums and refuses any difference; the currency must be `KINETIX_NOTARY_CURRENCY`.

**Metrics** (estate contract on `GET /metrics`, plus the notary's own):

| Metric | Type | For |
|---|---|---|
| `kinetix_build_info`, `kinetix_http_requests_total`, `kinetix_http_request_duration_seconds` | estate | conformance |
| `kinetix_grpc_client_calls_total{peer="payment"}` | counter | the snapshot calls |
| `kinetix_notary_last_finalized_epoch` | gauge | progress |
| `kinetix_notary_last_finalized_as_of_seconds` | gauge, unix time | **alert**: nothing final for 3 hours |
| `kinetix_notary_abandoned_epochs` | gauge | **alert**: any |
| `kinetix_notary_anchorer_balance_wei` | gauge | runway |
| `kinetix_notary_anchor_cost_wei` | gauge | mean of the last 24 finalized epochs, L2 execution + L1 data fee from the receipt |
| `kinetix_notary_runway_seconds` | gauge | **alert**: under 14 days (§7) |
| `kinetix_notary_fee_ceiling_hits_total` | counter | a ceiling set too low |

**Configuration** (fail fast, no defaults for anything that is a credential, an address or a bound):
`KINETIX_NOTARY_DATABASE_URL`, the four ledger variables the verifier already reads,
`KINETIX_NOTARY_ANCHORER_PRIVATE_KEY`, `KINETIX_NOTARY_PAYMENT_GRPC_ADDR`, `KINETIX_PKI_DIR`,
`KINETIX_NOTARY_CURRENCY`, `KINETIX_NOTARY_MAX_FEE_PER_GAS_WEI`, `KINETIX_NOTARY_HTTP_PORT`. Cadences may
default, as the estate allows: `KINETIX_NOTARY_EPOCH_INTERVAL` (1h), `KINETIX_NOTARY_TICK` (15s). The
daemon refuses to start when its key is not the ledger's anchorer.

**Status, 2026-10-12: 4a–4g built and green; 4h half done.** `bin/dev bin/ci` exits 0: 56 Foundry tests,
Slither 0, Aderyn nothing new, generation, `go mod verify`, vet, boundaries, golangci-lint 0, every Go
test including anvil and Postgres, Go coverage 85.6%. What each part proved:

| Part | Proof |
|---|---|
| 4a | protoc 36.2, sqlc 1.31.1 (SHA-256 from GitHub's own asset digests), protoc-gen-go v1.36.12, protoc-gen-go-grpc v1.6.2; Postgres 16 in the toolchain compose |
| 4b | commitment vectors equal `cast keccak`; every planned account proves inclusion; leaf order differs from the source and between epochs |
| 4c | over real mTLS with a CA made in the test: 20 kinds of broken stream all refused, never anchored; only a server whose SPIFFE id is `payment` is believed |
| 4d | lease race: 16 writers, 1 winner; a writer whose lease lapsed writes nothing; an epoch and its accounts are stored together or not at all |
| 4e | on anvil: anchor, receipt, finality, cost equal to the balance drop; every ledger refusal decoded before gas is spent; a stuck anchor replaced at the same nonce |
| 4f | 13 scenarios against a hand-written ledger simulator: replacement, fee ceiling, reorg, a dropped transaction, conflict, revert, a foreign nonce, divergence |
| 4g | end to end through the production wiring: payment's snapshot over mTLS, Postgres, the anvil ledger, final, and **every account verified by the contract**; a daemon killed mid-epoch is finished by the next with **one** transaction |
| 4h (repo) | one image with `notaryd`, `notary-migrate` and `notary-verify`, uid and gid 10001, Trivy: 0 fixable CRITICAL or HIGH; build-and-push passes the commit as the version |

**Left in 4h, all in `kinetix-infrastructure` and most of it the user's hands** (written up as a first-deploy
document there, the way search's was): the `kinetix_notary` database and role on the production cluster
(and the cluster healthcheck's database count), the `notary` leaf in `pki-prod` (CA password),
`KINETIX_NOTARY_DB_PASSWORD` in the SOPS store, the compose entries for the daemon and its migrator, the
image digest in `deploy/prod/versions.env`, the deploy dispatch in `release.yml`, three Grafana rules (no
final epoch for 3 hours, any abandoned epoch, runway under 14 days), the handshake rows, and funding the
anchorer with Base Sepolia ETH. The lab question is open: the lab has no chain, and the production
ledger must never receive lab data, so the lab either runs anvil with a ledger of its own or does not
run the notary at all.

**Work order**, each part green before the next:

- **4a** toolchain and generation: protoc and sqlc pinned by SHA-256, protoc-gen-go and
  protoc-gen-go-grpc by version, `bin/sync-contracts` pinned to contracts `v1.0.23`, Postgres 16 beside
  anvil in the toolchain compose. Generated code (abigen, protobuf, sqlc) is never committed.
- **4b** the core: accounts, kinds, commitments (known-answer vectors from `cast`), epoch states and
  their transitions, fee arithmetic, runway arithmetic — all in `internal/notary`, standard library only.
- **4c** the snapshot client (`internal/payment`), against a fake server over real mTLS.
- **4d** the store (`internal/postgres`): migrations, queries, lease and fencing, tested on Postgres.
- **4e** the chain writer (`internal/ethereum`): sign, send, replace, receipts with the L1 fee, tested
  on anvil.
- **4f** the use case (`internal/anchoring`): one tick of the state machine, every row of the tables
  above driven by hand-written fakes.
- **4g** `cmd/notaryd`, `cmd/notary-migrate`, HTTP, metrics, dig wiring; an end-to-end run on anvil
  with a fake payment that anchors, finalizes and survives a restart mid-epoch.
- **4h** the estate: images, compose entries, `notary` in `bin/kinetix-pki`, the handshake rows, Grafana
  rules for the three alerts, CD dispatch. Then the first epoch on Base Sepolia, once the anchorer is
  funded.

### Step 5 — the proof API

`notary.v1.ProofService` over gRPC and `GET /api/v1/proofs/liabilities/{epoch}` through Kong. The
caller's principal comes from identity's token; a principal receives the proofs for its own accounts
and nothing else, and the response is exactly the `kinetix.notary.inclusion.v1` document the CLI
already verifies. Unknown, not final and unavailable stay distinct all the way to the HTTP status.

### Steps 6–10 — accountability on chain

Five contracts that sit **beside** the deployed v1 ledger and never require changing it: the court and
the directory only read the ledger's public views, and the council becomes the ledger's anchorer when
the admin grants it `ANCHORER_ROLE` and revokes the notary's own key. Each gets the same treatment as
v1 — unit, fuzz and handler-based invariant suites, Go↔Solidity differential fixtures for every hash
and EIP-712 digest the two languages share, a fork test against the live Base Sepolia ledger,
Slither and Aderyn clean, immutable deployment through the CREATE2 factory, Sourcify verification.

#### Step 6 — signed statements and fraud proofs (`StatementCourt`, part 1)

Every proof the notary issues becomes a binding, publicly punishable claim.

- With each proof (step 5), and pushed to the holder when the epoch is final (through notification,
  which owns delivery), the notary issues an EIP-712 `EpochStatement`: ledger, epoch, commitment,
  amount, index, and `keccak256` of the siblings. Signed by a statement key the court knows.
- `proveBadProof(statement, signature, siblings)`: the signature is a registered signer's, the
  siblings hash to the signed value, the epoch exists — and `ledger.verifyInclusion` returns `false`.
  The notary signed a proof that does not prove anything: fault.
- `proveEquivocation(a, sigA, b, sigB)`: two signed statements for one epoch that cannot both be true
  — one commitment at two amounts or two indexes, or two commitments at one index: fault.
- Faults are permanent, per epoch, readable by anyone (`faultsOf(epoch)`, `FaultProven` events). The
  verifier and the browser page show them beside every verdict. There is no way to clear one.
- Signer rotation goes through the admin with a delay; a statement signed by a key at a time it was
  not registered proves nothing, so rotation cannot be used to disown old statements.

#### Step 7 — withholding (`StatementCourt`, part 2)

Step 6 punishes a lie the notary signed; this step punishes the silence that would hide an omitted
account — the limit PLAN §3 admits today. A holder registers a statement address (a key in their own
browser or wallet) through the authenticated proof API; on chain the holder may then
`requestStatement(epoch)` within a window after the epoch is final, and the notary must `answer` with
a statement naming that holder before a deadline, or anyone can record the epoch as
`Unanswered(epoch, holder)`. **Design to settle before building**, because each choice moves a
boundary: who attests that a holder exists (identity owns who a person is; payment owns whether they
hold anything), how a registration is itself made provable, how requests stay private (the holder
address, never the principal, goes on chain), and the rate limit that keeps a free challenge from
becoming a way to exhaust the notary's gas.

#### Step 8 — threshold anchoring (`AnchorCouncil`)

One hot key can anchor anything today. The council holds `ANCHORER_ROLE` and anchors an epoch only
with k-of-n EIP-712 signatures over its header (domain-separated by chain id and ledger address).

- Members: the notary's signer and **`notary-watch`**, an independent process with its own key that
  pulls the same snapshot from payment, rebuilds the tree on its own, and signs only when the root
  matches. Long term the watcher belongs to a different operator.
- Anyone holding enough signatures may submit (a relayer, or step 10's sponsored transaction).
- A member signing two different headers for one epoch is provable equivocation, recorded through
  the court's fault list.
- Membership and threshold change through the admin with a delay, and never below 2-of-n once the
  watcher exists.

#### Step 9 — ledger directory (`LedgerDirectory`)

A ledger succeeds another by sealing it (§6), but today only a config file says which one is current.
The directory is an append-only, on-chain list of generations — ledger address, genesis epoch,
genesis digest, scheme, registration time. Registering a successor on the same chain requires the
predecessor sealed and the successor's genesis equal to the predecessor's sealed head and next epoch,
both checked on chain; a successor on another chain is recorded with the values it declares and marked
as such. The v1 ledger is generation 0. Verifiers resolve the current ledger from the directory, so
their configuration shrinks to one address that never changes.

#### Step 10 — sponsored gas (`NotaryPaymaster`, ERC-4337)

A holder who wants to request a statement or prove a fault should not need testnet ETH first. The
paymaster sponsors user operations whose calls go only to the court's request and fault functions and
the council's submit, with a per-sender rate limit and a budget whose runway is monitored like the
anchorer's (§7). The EntryPoint version is pinned and its presence on the network checked before
deploy, the way the CREATE2 factory was.

### Step 11 — verification without Kinetix

A browser verifier in the public docs: it reads the chain directly from the visitor's browser, so
checking a proof involves no Kinetix server at all. It resolves the current ledger through the
directory and shows any fault recorded against the epoch. The CLI image is published for auditors.

### Step 12 — phase 2: evidence anchoring

`EvidenceRegistry`, batched Merkle roots of:

- **order lifecycle transitions**, from order's existing `OrdersChangedSince` feed;
- **proof of delivery**: the fulfillment contract carries only `pod_photo_url` today, and a URL's
  content can be replaced. A `pod_content_sha256` must be computed where the photo is captured and
  stored (courier app and matching), then anchored, so a return dispute can show the delivery
  evidence existed at time T and is unchanged;
- **review integrity**: review content hashes, so a review cannot be quietly edited or removed.
  review serves no gRPC today; it needs a server and a `ChangedSince` first.

### Step 13 — scheme v2: zero-knowledge proof of liabilities

The total stops being public, and the proof gets stronger.

- **gnark** (Consensys, Go) — the same language as this service; PLONK over BN254 with a KZG
  structured reference string from a public ceremony, so there is no per-circuit trusted setup.
- The circuit proves: every leaf amount is in `[0, 2^64)`; the tree's root is built from those
  leaves with **Poseidon2** (cheap inside a circuit, unlike keccak); the leaves sum to a total `T`;
  and the ledger stores only `C = Poseidon(T, blinding)`.
- gnark exports a Solidity verifier. The handwritten part is `LiabilityLedgerV2` (scheme 2) around it,
  registered in the directory as a new generation and judged by a court that speaks scheme 2, tested
  in Foundry the same way v1 is, with differential fixtures from the Go prover. The generated
  verifier is not counted as this repository's own Solidity.
- Account holders still verify inclusion with a Merkle path — Poseidon instead of keccak — in the
  browser (WASM) or the CLI.
- Beyond one circuit's capacity: recursive aggregation of per-batch proofs.
- Budget to measure before committing: proving time and memory per 2^16 and 2^20 accounts on the
  production VPS, which shares 16 GB with the whole estate.

### Step 14 — the next service: product passport

After the notary is complete: per-unit provenance for serialized goods (ERC-721), minted by a
verified merchant and transferred to the buyer on delivery. It needs serial-level inventory, which
warehouse does not have today (no serial column anywhere in its schema), and buyer wallets
(account abstraction — the step 10 paymaster is the first piece of that), so it starts with that
owner-side work.

## 11. Open questions

- **Licence.** Every Solidity file says `UNLICENSED`. A public verifier contract is usually MIT so
  others may read and reuse the verifier; the choice is the user's.
- **Publishing the verifier image** to GHCR, and under which name.
