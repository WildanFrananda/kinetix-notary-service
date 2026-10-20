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
| 11 | Licence (2026-10-13) | **MIT** for every new file. The v1 ledger's files keep `UNLICENSED`: they are the verified source of a deployed contract and stay byte-identical; new contracts import their types as they are. |
| 12 | Statement key (2026-10-13) | A **fourth key**, separate from the anchorer: it signs statements off chain, never holds ETH, lives in the SOPS store like the anchorer. A leak of one never hands over the other. |
| 13 | Signer changes (2026-10-13) | Through the admin, **after a 7-day delay**, and a signer is registered for a **range of epochs**, never a span of time — so no rotation can disown a statement already signed. |
| 14 | Delivering statements (2026-10-13) | Pushing a statement to its holder is **notification's job** (the estate rule: no service takes another's). Step 6 splits: **6a** in this repository, **6b** in notification. |
| 15 | Withholding (2026-10-13) | Step 7 as proposed, R1–R5: the notary signs registrations, a second commitment derivation from a switch epoch, the inclusion proof is the answer and anyone may post it, a 7-day request window and a 24-hour answer deadline, and a separate MIT `WithholdingCourt` beside the deployed court. |
| 16 | Answering key (2026-10-14) | A **fifth key**, the answerer, in SOPS like the anchorer. It sends `answer` transactions and nothing else. The anchorer judges its own anchor by its nonce, so a second sender on that key would make it misread its own epoch, and a leak of one key never hands over another. |
| 17 | Proving the statement address (2026-10-14) | The holder signs an EIP-712 `RegistrationRequest{holder, principal, issuedAt}` under the withholding court's domain with the key of the address they register, and posts it with their access token. The notary refuses a signature that does not recover to that address, a principal that is not the token's, or an `issuedAt` older than five minutes. |
| 18 | Council threshold (2026-10-14) | **2-of-2**: the notary and the watch. No spare key: a third Kinetix key would let Kinetix anchor without the watch, and the council would mean nothing. A stopped watch stops anchoring, and the three-hour finality alert says so. The contract refuses any threshold below 2. |
| 19 | Where the watch lives (2026-10-14) | **`cmd/notary-watch` in this repository**, with its own image entry, container, SPIFFE identity and key. `bin/check-boundaries` keeps it away from the anchoring and Postgres code. |
| 20 | How the watch sees the same snapshot (2026-10-14) | **Payment records each snapshot** — its time, currency, count, total and a SHA-256 digest of its rows — and serves that record to the watch. notaryd forwards the rows; the watch checks them against payment's record, not against anything notaryd says. Payment keeps three values per snapshot and needs no new key. |
| 21 | The notary's council key (2026-10-14) | **The anchorer key** signs the notary's approval and still sends the transaction and pays its gas. Only the watch gets a new key, the sixth. |
| 22 | What the directory records (2026-10-15) | A generation is a ledger **with its statement court and its withholding court**: chain id, the three addresses, genesis epoch and digest, scheme and registration time. On this chain the directory reads the genesis and the scheme from the ledger itself, and checks that the court judges that ledger and the withholding court stands beside that court. A court is replaced only together with a successor ledger. The council stays out: it is how epochs are anchored, not how they are verified. |
| 23 | Who registers a successor, and when it counts (2026-10-15) | **The admin only.** A successor on the same chain counts at once, because the directory checks it on chain: the current ledger is sealed, and the successor starts exactly at its sealed head and next epoch. A successor on another chain is proposed with what cannot be checked here (its chain id, addresses and scheme), takes its genesis from the sealed predecessor, and counts after **7 days**; the admin may cancel it until then, and anyone may activate it after. Once a generation lives on another chain, this directory takes no more: the next is registered in a directory on that chain. |
| 24 | How the off-chain side uses it (2026-10-15) | `notary-verify` is configured with **the directory alone**, beside the RPC URL, chain id and timeout. It resolves the generation that holds the claim's epoch and verifies against that ledger and court; its ledger and court variables go. An epoch whose generation lives on another chain gets a verdict of its own, never "not found". notaryd keeps its explicit addresses and refuses to start unless the directory's current generation is its ledger and courts. |
| 25 | Deploying it (2026-10-15) | Like every contract beside the ledger: the ledger's admin key, a 3-day admin transfer delay, the CREATE2 factory with a salt, Sourcify. Generation 0 is the v1 ledger with the live court and withholding court, registered by the constructor under the same checks. |
| 26 | Stake for the paymaster (2026-10-15, revised 2026-10-16) | **1 ETH, locked for one day once unlocked.** First ruled unstaked: ERC-7562 lets an unstaked paymaster read only the operation, the sender's code and its own bytecode, so a registered holder could still make it pay for calls the court refuses, about 0.0000006 ETH each, and the 1 ETH the bundlers ask for was more than a faucet gives. Revised once the user had 5 Ethereum Sepolia ETH to bridge: staked, the paymaster reads the courts and the ledger and sponsors only a call the court would take, at a time it would take it (`validAfter`, `validUntil`). It still writes no storage, so no per-holder counter is needed. The stake is the bundlers' guarantee, never spent and never slashed; only the admin adds to, unlocks or withdraws it, and the constructor refuses an unstake delay under one day. `notary-sponsor` asks the same checks through the paymaster's `sponsorship` view before it signs, so a refused or early call is reported and never sent. |
| 27 | What it sponsors (2026-10-15) | **A holder's calls only:** `requestStatement` and `recordUnanswered` on the withholding court, and `proveBadProof` and `proveEquivocation` on the statement court. The council's anchor stays the anchorer key's to pay. |
| 28 | EntryPoint and account (2026-10-15) | **EntryPoint v0.8 with EIP-7702.** The holder's own key delegates to the audited `Simple7702Account` of v0.8, so the holder stays `msg.sender`, which the immutable withholding court requires. The public bundlers of Candide and Pimlico serve v0.8 on Base Sepolia without an API key. |
| 29 | Who builds the operation (2026-10-15, revised 2026-10-16) | **A Go CLI in this repository, in the same image.** The holder runs it locally with their own key and a bundler URL. First written as "the browser verifier (step 11) will use the same paymaster"; revised by decision 32: the browser never sends anything. |
| 30 | The deposit (2026-10-16) | **0.1 ETH**, bridged from Ethereum Sepolia with the stake: 1,000 operations at the 0.0001 ETH ceiling, about 13,000 at the fees seen on 2026-10-15. The 0.0005 ETH warning stays. |
| 31 | How the browser verifies (2026-10-16) | **The Go verifier compiled to WebAssembly**, not a second implementation. `cmd/notary-verify-wasm` (`GOOS=js GOARCH=wasm`) wires the same `wiring.Verifier` graph as `notary-verify` and hands the same report lines back, so the page and the CLI cannot disagree. The browser's own `fetch` carries the JSON-RPC: every public Base Sepolia endpoint measured on 2026-10-16 (`sepolia.base.org`, publicnode, drpc) answers CORS. It costs size: 44 MB raw and 9.4 MB gzipped, because go-ethereum comes whole, so the page loads it only when asked to verify. |
| 32 | What the browser may do (2026-10-16) | **Verify only.** A browser wallet will not sign an arbitrary EIP-7702 delegation, and a page that takes a holder's key is a page that can lose it. When a proof is provably false, the page prints the exact `notary-sponsor prove` command for the CLI image. |
| 33 | The CLI image (2026-10-16) | **`ghcr.io/wildanfrananda/kinetix-notary-cli`**, separate from the service image. It is multi-arch (amd64 and arm64), so an auditor's Apple Silicon machine runs it natively. It holds `notary-verify` and `notary-sponsor` only, and is tagged with the notary's own release tags. This settles the open question in §11. Making the package public is the user's step. |
| 34 | Where the page lives (2026-10-16) | **A new section, Verify, on the public docs site** (`kinetix-docs`), for account holders and auditors. It holds two pages: the how-to with the verifier, and an explanation of what the proof proves and what it does not (§3). The notary's own service page stays deferred. |
| 35 | How the page gets the verifier (2026-10-16) | **Each notary release attaches `notary-verify.wasm`, the matching `wasm_exec.js`, and `SHA256SUMS` to its GitHub release.** The build is reproducible: `-trimpath`, `-buildvcs=false`, the pinned Go. The docs pin a release and its SHA-256, and refuse to build on a mismatch. The page states the version and the hash, so anyone can rebuild the file from the tag and compare. The RPC endpoint defaults to `https://sepolia.base.org`, and the visitor may replace it. The chain id and the directory come from `deploy/base-sepolia.env` and are shown on the page. |
| 36 | Step 12's scope (2026-10-16) | **All three evidence streams in one step**: order lifecycle, proof of delivery, reviews, each built in its owner and anchored by one registry. Chosen over POD first (the recommendation): nothing goes live until all three are coherent, so the step is long. |
| 37 | Where order's evidence comes from (2026-10-16) | **An append-only lifecycle log in order**, written in the same transaction as each change it records (status changes and the return lifecycle, with the actor), served by a new feed whose cursor cannot skip a transaction that commits late. `OrdersChangedSince` returns current state only, and returns never touch `updated_at`, so it cannot carry transitions. |
| 38 | Where the delivery photo's hash comes from (2026-10-16) | **The courier app hashes the bytes before upload, and storage refuses any other bytes**: the presigned PUT signs the SHA-256 and may write the object once. Matching records the photo's and the signature's hashes at `delivered`, which is terminal. A new app release is part of the step. |
| 39 | Reviews in step 12 (2026-10-16) | **Now**, chosen over deferring (the recommendation): review gets the estate's first PHP gRPC server, a `review.v1` contract and a feed. The app only creates reviews today; the log is what shows a review was not changed or removed in the database. |
| 40 | MinIO (2026-10-16) | **Built from source by our CI and published to GHCR.** quay.io now refuses `minio/minio` and `minio/mc`, Docker Hub has no such tag, and upstream is archived. The build takes the last upstream release, `RELEASE.2025-10-15T17-29-55Z`, reproducibly. No security fixes will follow; replacing MinIO is later work, owned by infrastructure. **2026-10-17:** the user chose to replace it after step 12 is released, by a lab spike of SeaweedFS, RustFS and VersityGW against decision 38's seven cases (infra `docs/OBJECT-STORAGE.md` §8). |
| 41 | Who anchors evidence (2026-10-17) | **A seventh key, the evidence anchorer**, in the SOPS store like the anchorer, funded by the deployer like the answerer. It holds the registry's anchorer role and nothing else. Chosen over the existing anchorer key, which judges its own epochs by its nonce (decision 16, `anchoring/anchorer.go`), and over the council, whose watch would have to read all three feeds. |
| 42 | How often evidence is anchored (2026-10-17) | **Daily, per stream, only when the stream has new records**: at most three transactions a day. Chosen over hourly (the recommendation): a proof is available the day after the event, and gas is a twenty-fourth. |
| 43–48 | Scheme v2 (2026-10-20, on my recommendation: *"pake rekomendasi kamu"*) | **gnark PLONK on Base**, inside the v1 interfaces the immutable directory requires; Aztec Ignition SRS; 256 accounts per generation pending the VPS measurement; the blinding kept for auditors; the watch re-checks the proof. Full table in §10 step 13. |

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
  canonical CREATE2 factory with a configured salt, so a successor's address is known before it
  exists. The deploy scripts call the factory themselves (`script/Create2Factory.sol`) and
  `bin/rehearse-deploy` proves it on each network's chain id. The first ledger and the first court are
  the exceptions — see step 2 and step 6.

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
  always CREATE2, which is exactly why it could not see this. **That fix did not hold** (2026-10-13):
  the court went out as plain CREATE with the setting on, because forge 1.8.3 uses the factory on
  chain 31337 only, and the rehearsal ran on 31337. The scripts now call the factory directly, and
  `bin/rehearse-deploy` replays them with chain id 84532 (step 6).
- **Blockscout cannot verify it yet**: its newest Solidity compiler is 0.8.36 and the ledger is built
  with 0.8.37. forge printed "All (1) contracts were verified!" over Blockscout's own
  `Fail - Unable to verify`. Verification moved to Sourcify, which needs no key; `bin/inspect-ledger`
  now asks Sourcify directly instead of trusting forge's summary. Submitting to Blockscout as well,
  once it lists 0.8.37, is optional.
- **The store write left the prod secrets in plaintext for about a minute.** Contained and fixed in
  `kinetix-infrastructure`; see its `docs/SECRETS.md` rule 6.

### Step 3 — the snapshot contract (owner-side work). **Done 2026-10-12** (contracts v1.0.23, payment v0.1.11).

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

### Step 4 — `notaryd`, the anchorer. **Live 2026-10-12** (notary v0.2.0; epoch 0 final on Base Sepolia).

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

**4h done 2026-10-12.** Everything listed for it — database and role, the `notary` leaf in `pki-prod`,
the SOPS secrets, compose entries, the image digest, the deploy dispatch in `release.yml`, three Grafana
rules, the handshake rows, the anchorer funded — is in `kinetix-infrastructure` and deployed; epoch 0
went final on Base Sepolia and `bin/inspect-ledger base-sepolia` reports every fact as expected. The
lab question was answered the same day: anvil with a ledger of its own (step 5, "Lab").

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

### Step 5 — the proof API. **Live 2026-10-13** (notary v0.3.0, gateway v1.0.0).

**Built 2026-10-12, deployed 2026-10-13 (v0.3.0).** `GET
/api/v1/proofs/liabilities/{epoch|latest}` on notaryd's HTTP port: identity's RS256 access token
(JWKS, issuer, audience, `token_use=access`), `sub` is the holder, and only that holder's accounts are
proven. The rebuilt tree must reproduce the anchored root or nothing is served. 404 unknown, 409 not
final, 410 abandoned, 503 store unreachable. Each document carries `account` (holder, kind, salt), and
`notary-verify` checks the account opens the commitment and verifies a whole API response. JWKS
refresh is limited to once per 10 s and never waits more than 100 ms (a forged `kid` held a request
10 s under keyfunc's default — measured). No gRPC `ProofService` until a service needs one.
**Live 2026-10-13:** notary v0.3.0 and the gateway's `notary-proofs-route` (v1.0.0) deployed after
infrastructure moved the notary's environment into one `x-notary-env` block, so production carries
the JWT issuer like every other service. From outside: no token and a malformed token are refused by
Kong with 401, `POST` finds no route.

**Lab** (the user's choice of 2026-10-12, `kinetix-infrastructure/docs/NOTARY-LAB.md`): anvil on chain
31337 with a block every 5 s, and a ledger deployed by a one-shot container through the canonical
CREATE2 factory at a fixed address, `0xF049C3811184CE357Ef1ba3FF4EcfA623f1f1d78`, anchorer anvil's
first default account. **Proven end to end on 2026-10-13**: a customer registered through Kong, their
wallet opened by payment's real balance read, epoch 0 built with that one account, anchored, final,
fetched through Kong with the customer's token and verified `included` by `notary-verify` against the
lab chain. Negative controls: the holder changed → `rejected`; chain 84532 configured → `wrong_ledger`;
a second customer → an empty `proofs` list. Every status seen through Kong: 404 no final epoch,
409 epoch 0 submitted, 404 epoch never built, 400 `abc`, 401 no token or an expired one.

**Found by that run, fixed, not yet released:** a 4xx body carried the whole error chain —
`proving.Prover.resolve: epoch_not_anchored: …`, and for a refused token the JWT library's own reason.
`notary.Error` now keeps its sentence apart from its cause, and the API sends only the sentence
(`notary.MessageOf`); `Error()` reads exactly as before, so logs are unchanged. The HTTP test pins
every public message and went red against the old code.

**Left:** a notary release carrying that fix and the split CI gate.

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

#### Step 6 — signed statements and fraud proofs (`StatementCourt`, part 1). **6a live 2026-10-13** (notary v0.4.0).

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

**Design settled 2026-10-13 (step 6a).**

*The statement.* EIP-712, domain `{name: "Kinetix Statement Court", version: "1", chainId,
verifyingContract: the court}`, type
`EpochStatement(address ledger,uint64 epoch,bytes32 commitment,uint64 amount,uint64 index,bytes32 siblingsHash)`,
where `siblingsHash = keccak256(hash₀ ‖ sum₀ ‖ hash₁ ‖ sum₁ ‖ …)`, each hash 32 bytes and each sum 8 bytes
big-endian, exactly as the siblings appear in the proof. Every field is already in the proof document,
so the document gains only `statement: {court, signature}`; the verifier rebuilds the statement from
the proof it is reading and cannot be shown a signature over something else. Signatures are 65 bytes,
`r ‖ s ‖ v` with `v` 27 or 28, low `s` (OpenZeppelin's `ECDSA` refuses the other half).

*Signers.* An append-only list of terms `{signer, fromEpoch}`; the signer of epoch `e` is the last
term with `fromEpoch ≤ e`. The first term is set at construction from the ledger's genesis epoch. The
admin schedules `{signer, fromEpoch}` with `fromEpoch ≥ ledger.nextEpoch()` and above the last term;
after the 7-day delay anyone may activate it, and activation refuses (`SignerChangeStale`) if the
ledger has meanwhile reached `fromEpoch`. So when a term becomes active no epoch it covers exists yet,
and nothing the previous signer could have signed changes hands. The admin may cancel a pending term.
The notary signs each epoch with the key the court names for it, and refuses to serve a proof it
cannot sign: a statement withheld selectively is exactly what a cheating notary would do.

*Faults*, for anchored epochs only (`genesisEpoch ≤ e < ledger.nextEpoch()`), permissionless, no bond:

- `proveBadProof(statement, signature, siblings)` — signed by the epoch's signer, the siblings hash to
  `siblingsHash`, and `ledger.verifyInclusion` answers `false`. Evidence: the statement's digest.
- `proveEquivocation(a, sigA, b, sigB)` — both signed by the epoch's signer, same ledger and epoch,
  not identical, and they share a commitment or an index: one leaf cannot have two amounts, two
  indexes or two paths, and one index cannot hold two leaves. Evidence: `keccak256(lo ‖ hi)` of the two
  digests, sorted, so the pair counts once whichever way round it is submitted.

Each piece of evidence is recorded once (`isProven`); `faultsOf(epoch)`, `faultCount(epoch)` and
`totalFaults()` read them back, `FaultProven` announces them. Nothing removes a fault.

*Reading side.* The court exposes `statementDigest`, `siblingsHash` and `recoverSigner` as views, so
`notary-verify` checks a signature with the same code the court would use to convict, and compares
the digest with its own — a disagreement is refused, as with inclusion. The verifier is configured
with the court's address (never the one a document names) and reports, beside every verdict, whether
the statement is attributable to the epoch's signer and how many faults the epoch carries. `notaryd`
exports `totalFaults()` as a gauge, so a proven fault reaches Grafana.

*Work order*, each part green before the next:

- **6a-1** the contracts: `StatementCourt`, its interface and types (one per file), `StatementHashing`;
  unit, fuzz and handler-based invariant suites; the deploy script through the CREATE2 factory.
- **6a-2** Go↔Solidity: the statement digest and siblings hash in `internal/keccak`, a fixture
  generated by Go with a fixed test key and re-verified by Foundry against a court deployed at the
  fixture's address.
- **6a-3** a fork test against the live Base Sepolia ledger: a court deployed beside it convicts a
  forged statement about epoch 0.
- **6a-4** `notaryd`: the statement key, the court adapter, a statement on every proof, the gauge.
- **6a-5** `notary-verify`: statement and faults in every report.
- **6a-6** the estate: `bin/deploy-court`, `bin/inspect-court`, `bin/new-statement-key`, the lab court,
  compose, the Grafana rule. Then the court on Base Sepolia, by the user's keystore.

**Status, 2026-10-13: 6a is live. v0.4.0 signs every proof in production, proven end to end.**
`bin/ci-contracts`: 104 tests on the CI profile (court: 30 unit, 5 fuzz at 20 000 runs, 5 invariants over
16 384 calls, 2 differential, 7 deploy), 4 fork tests at Base Sepolia block 48046426 (moved to 48063800 by step 7), Slither 0, Aderyn
nothing new. `bin/ci-go`: coverage 86.4%, the anvil suite included. What each part proved:

| Part | Proof |
|---|---|
| 6a-1 | two planted mutations — a true statement convicted, a stale signer activated — each turned unit, fuzz and invariant suites red |
| 6a-2 | 13 statements Go signed reproduce in Solidity: siblings hash, EIP-712 digest and recovered signer; renaming one EIP-712 field in Go turned Foundry red |
| 6a-3 | a court deployed beside the live ledger convicts a forged statement about epoch 0 against the real anchored root |
| 6a-4 | every proof carries a statement signed with the key the court names for its epoch; no key, no proof (503), never an unsigned one; notaryd checks court identity, digest agreement and the newest signer before it serves |
| 6a-5 | on anvil: a statement signed by the signer is attributable, any other is not; a signed false proof is reported `provable`, convicted by a real `proveBadProof`, and the next verification reports the fault |
| 6a-6 | lab court at `0x9501aaB59C31Edc7aB0dE9d9453B62c09893EB9a`, rehearsed twice on a throwaway chain; `bin/deploy-court`, `bin/inspect-court`, `bin/new-statement-key`; the `kinetix-notary-proven-fault` alert |

The production court, deployed by the user on 2026-10-13 and checked line by line by
`bin/inspect-court base-sepolia`:

| | |
|---|---|
| court | `0xE1Ffcb004Dd734F83ED0B194fdaFF5B0c7c5927B`, Base Sepolia (84532), 13,392 bytes of code |
| deployed | block 48052866, tx `0xa86e057eeb6fa0936302e413a50f68aee60bf41734fcf0bed86fd85faff2367b` |
| cost | 3,072,611 gas × 0.006 gwei = 0.000018435666 ETH |
| judges | the ledger `0x4E869A55031DE3B0d030AEBECc66E523686761a2` |
| admin | `0x2BfBa1D3d8595ad7d1E10a71BB4ec4Ac3671748E`, transfer delay 259200 s |
| signer | `0xBB58A2D2eE5B09d57a1d406A01451dd66854F2CA` from epoch 0, change delay 604800 s; its key is `KINETIX_NOTARY_STATEMENT_PRIVATE_KEYS` in the prod store |
| source | verified on **Sourcify**: `match` |
| faults | 0 |

Like the first ledger, **the court was deployed by plain CREATE** (deployer nonce 2), not through the
factory: forge 1.8.3 routes a salted `new` through the factory only on chain 31337, so the lab and every
rehearsal used it and Base Sepolia did not, even with `always_use_create_2_factory` on. The script's
`== Return ==` printed `0x88A89C885D790DD872A741C485d867bAd7bbD15b`, the CREATE2 address, where nothing
exists; the broadcast record and the chain say `0xE1Ff…927B`. The court is correct either way. Since
then both deploy scripts call the factory themselves through `script/Create2Factory.sol`, and
`bin/rehearse-deploy` (in `bin/ci-contracts`) replays them on anvil with each network's chain id and
fails unless the broadcast is one CREATE2 through the factory. The old ledger script swapped in turned
it red; a mutated address formula turned the deploy suites red. The lab addresses did not move: the same
init code and salt through the same factory land where they always did.

**Released 2026-10-13 as v0.4.0.** notaryd started against the court. A fresh customer's proof from
epoch 20, fetched through Kong, carried a statement that `notary-verify` reported `attributable` with
0 faults, and one changed signature byte made it `not_attributable`
(`kinetix-infrastructure/docs/NOTARY-COURT-DEPLOY.md`). The same check found that production's
Prometheus had never scraped the notary, so the stall, abandon, runway and proven-fault rules had been
silent since step 4. That is fixed in infrastructure, and `bin/prod-deploy` now refuses a service that
Prometheus does not scrape. **6b**, pushing statements to holders, is notification's work.

**6b, designed and built 2026-10-20 (uncommitted).** notification owns delivery, so notary decides only who is told
and when, and notification sends it.

*Ruled by the user on 2026-10-20:*
- **The message carries no amount.** It says that the holder's statement for an epoch is signed and links the public
  Verify page. The email provider reads whatever is in it, so the balance, salt and secret stay out of it, and the
  holder fetches the proof with their own token.
- **Daily, and only on a change.** The first epoch final after midnight WIB is the day's statement epoch. A holder is
  told only when their amounts differ from the statement last planned for them.

*Taken on my recommendation:*
- notaryd calls notification's `Notify` over gRPC with a new template, `TEMPLATE_LIABILITY_STATEMENT` (kinetix-contracts,
  additive). That is a new edge, notary → notification.
- **A store-backed queue.** Each planned statement is a row in `statement_deliveries`, so an outage of notification
  loses nothing. The idempotency key is `liability-statement:<epoch>:<holder>`, so a retry or a second notaryd never
  sends twice.
- No amount ever lives in notification's store either: its params are the epoch and the URL.
- All channels notification has for the holder: email today, and push once an app registers devices.

*What notary keeps:*
- `statement_progress`, the last epoch examined and its day;
- `statement_deliveries`, holding `pending`, `sent` (with notification's id), `refused` (no address, not asked again
  until the balance changes) and `superseded` (a newer day replaced it before it was handed over).

Unknown is not no: a refusal that is not about the holder (a missing grant, an unknown template, a transport error) is
retried with a backoff, from 1 minute to 6 hours, and `kinetix-notary-statements-overdue` warns after six hours.

*Proven:*
- unit tests on fakes for the day boundary, change detection, suspense, races and backoff;
- the store against Postgres;
- the adapter against a real mTLS server, for every refusal code and status;
- the daemon end to end on anvil and Postgres: four holders told once, one refused, nothing repeated on an unchanged
  epoch, and statements waiting through an outage and sent after it.

Three mutations (suspense let in, the same-day supersede, a refusal misread) each turned the suite red.

*Release order* (infra `docs/NOTARY-STATEMENTS.md`): contracts, then infra, then notification, then notary. The user
published kinetix-contracts v1.7.0 (`26984a9`) the same day. notary and notification are pinned to it, notary's full
gate passes on it (Go coverage 84.8%), and so does notification's suite (111 tests, Postgres included).

#### Step 7 — withholding (`WithholdingCourt`). **Ruled 2026-10-13 (decision 15); live 2026-10-14 (v1.0.1).**

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

**The design, ruled by the user on 2026-10-13 as proposed, with a 7-day window and a 24-hour deadline.**

The obstacle is the commitment. Today it is `H(holder ‖ kind ‖ salt)` with a fresh salt per epoch, so
nothing on chain can tell which leaf belongs to whom, and a court cannot check that an answer names
the right account without a signature from the notary. The proposal moves that knowledge into the
commitment and lets the inclusion proof be the answer.

- **R1. The notary attests that a holder exists, by signing a registration.** A holder authenticates
  through the proof API (identity says who they are). They present a statement address `A` and a
  signature from `A`, which proves they hold its key, so nobody can register an address that is not
  theirs. The notary returns an EIP-712 `Registration{ledger, holder: A, account, fromEpoch}`, signed
  by the court's signer for `fromEpoch`. It only signs for an account in its latest final epoch, which
  is payment's snapshot. Nothing about identity or payment goes on chain, and neither service gains a
  key. A refused registration is visible to the holder but cannot be proven on chain; that limit is
  written into §3.
- **R2. A new commitment derivation, from a switch epoch onward.** `account = H(0x05 ‖ kind ‖
  accountSecret ‖ holder)` uses a random secret per account that is kept by the notary and served to the
  holder with their proofs. The commitment becomes `H(0x06 ‖ account ‖ salt)`. Tag `0x04` is today's
  commitment, `H(0x04 ‖ kind ‖ salt ‖ holder)`, so both derivations stay apart from each other and from
  the tree's own tags. The secret must outlive every registration that names its account: a lost
  secret means unanswerable requests, so it lives in the notary's database and its backups. The tree and the deployed
  ledger do not change, because the commitment is opaque to them. The proof file gains the account
  secret and a new `format`, and `notary-verify` reads both formats. A principal never goes on chain,
  and `account` cannot be linked to one without identity's records and the secret.
- **R3. The answer is the inclusion itself, so anyone may post it.**
  `answer(registration, epoch, salt, amount, index, siblings)` recomputes the commitment from
  `registration.account` and `salt`, then needs `ledger.verifyInclusion` to return `true`. It takes no
  notary signature, so the holder can answer their own request with a proof the API already gave them.
  notaryd answers the requests nobody else does. Answering reveals that account's amount for that
  epoch, linked to `account` and not to a person, and only when its holder asked.
- **R4. Bounds that keep a free challenge from draining the anchorer.** Only `A` itself may call
  `requestStatement(registration, signature, epoch)`. It applies only to an epoch from `fromEpoch` on
  that was anchored at most 7 days earlier. Each registration may hold one open request at a time and
  one request per epoch. The answer is due 24 hours after the request. After that, anyone may call
  `recordUnanswered`, which records a permanent `Unanswered(epoch, account)` fault. That fault is a
  silence, kept apart from step 6's lies in its own kind and its own Grafana rule. Measured on
  2026-10-13: an answer costs 46.7k gas at the median (57.6k for a deeper tree) and is the notary's to
  pay; a request costs about 101k and a recorded silence about 170k, paid by whoever sends them.
- **R5. A separate contract, `WithholdingCourt` (MIT), beside the deployed `StatementCourt`.** It reads
  the court's signer terms and the ledger, and changes neither. Folding the game into a redeployed
  court would orphan the live court and its record. notaryd exports its fault count as a second
  gauge.

What it still does not prove: the notary knows who has registered, so an unregistered account can
still be left out unnoticed. Protection starts at `fromEpoch`, and registering early is the defence.
Letting a registration attest backwards, to the account's first epoch inside the request window, is a
possible extension and is not part of this proposal.

*Work order*, each part green before the next, as in 6a:

- **7-1** (**done 2026-10-13**: 29 unit, 6 fuzz, 4 invariants over 16 384 calls, 7 deploy; mutants for the
  deadline, the window and the inclusion check each turned a suite red) the contract: `WithholdingCourt`,
  its interface and types (one per file), `RegistrationHashing`
  and `AccountCommitments`; unit, fuzz and handler-based invariant suites; a deploy script through the
  CREATE2 factory, its deploy test, and `bin/rehearse-deploy`.
- **7-2** (**done 2026-10-13**: `cast` known-answer vectors for both new hashes; `withholding-v1.json` re-verified
  by Foundry, which also recomputes each account id from its documented preimage; Go's account tag changed
  to 0x07 turned Foundry red) Go↔Solidity: the account id, the second commitment and the registration
  digest in `internal/keccak`, a fixture written by Go and re-verified by Foundry.
- **7-3** (**done 2026-10-13**, fork block 48063800: the deployed court answers ERC165 for `IStatementCourt`,
  refuses a registration its signer did not sign, and an account the real tree lacks ends in a recorded
  silence) a fork test: a withholding court deployed beside the live ledger and court.
- **7-4** `notaryd`, in three parts:
  - **7-4a** (**done 2026-10-14**) every epoch built from now on uses the second derivation. A secret per
    account is drawn once, kept in `account_secrets` (migration 00002) in the same transaction as the
    epoch, and checked against what the database already holds. A row naming an account without its
    secret is refused by a foreign key. The prover refuses to serve an account whose opening does not
    reproduce its commitment, which also catches a changed salt that the root check never saw. Proofs
    travel as `kinetix.notary.inclusion.v2` with the secret, and v1 files still verify.
  - **7-4b** (**done 2026-10-14**) registration: `POST /api/v1/proofs/registrations` (decision 17) issues one
    registration per account the principal holds in the latest final epoch, signed by the court's
    signer for that epoch and kept in `registrations` (migration 00003; a foreign key refuses a
    registration naming an account the principal does not hold). It answers `registration_not_open`
    until a final epoch carries secrets. notaryd now needs `KINETIX_NOTARY_WITHHOLDING_COURT_ADDRESS` and
    refuses to start unless that contract stands beside its court and hashes a registration as it does.
    Every EIP-712 digest is also checked against go-ethereum's own `apitypes`. On anvil, a registration
    issued through the API was accepted by the real contract's `requestStatement`. The gateway needs
    a `POST` route for it (7-6).
  - **7-4c** (**done 2026-10-14**) answering. The answerer key (decision 16) sends answers with its own
    nonce, and a stuck one is replaced at the same nonce, an eighth dearer, under the fee ceiling.
    notaryd reads `StatementRequested` logs (100 000 blocks back from the head at start, never before
    `KINETIX_NOTARY_WITHHOLDING_FROM_BLOCK`, `KINETIX_NOTARY_RPC_LOG_RANGE` blocks per call; nothing is
    stored, because a request lives 24 hours), asks the court whether each is still open, and answers from its own records after
    rebuilding the epoch's root. A request it cannot answer, and a refusal or a disagreement for one
    request, is counted and logged as an error without holding up the others; an outage is retried.
    Gauges: `kinetix_notary_withholding_silences`, `_open_requests`, `_unanswerable_requests` and
    `kinetix_notary_answerer_balance_wei`. On anvil the whole game ran: a holder registered through
    the API and requested epoch 0, and notaryd answered on chain. Epoch 1 then left the holder out, the
    request stayed open with the unanswerable gauge at 1, and after 24 hours a stranger recorded the
    silence and the silence gauge read 1. Answering with a wrong salt turned that test red.
- **7-5** (**folded into 7-4a**) `notary-verify`: both proof formats.
- **7-6** operations. Done in this repository on 2026-10-14: `bin/deploy-withholding`,
  `bin/inspect-withholding`, `bin/new-answerer-key`, and `bin/fund-hot-key`, with `fund-anchorer` and
  `fund-answerer` as its wrappers. Both funding paths were rehearsed on anvil. The deploy was rehearsed
  on an anvil fork of Base Sepolia at block 48063800, beside the real court: one CREATE2 through the
  factory, `== Return ==` equal to the deployed address `0xe980Cd4199d491D635a2d8fd563FC7730716a601`, and
  every inspect line `ok`. Because the address follows from the salt and the court, the real deploy
  lands there too. **Deployed by the user on 2026-10-14**, exactly there: `0xe980…a601`, block
  48086256, one CREATE2, 1.92M gas, Sourcify `match`, every inspect line `ok`. The predicted address
  held on the real chain, which the ledger and the court never did. Infra carries the lab deploy, the
  prod and shared compose lines and three Grafana rules; the gateway carries the `POST` route. Left,
  all the user's: fund the answerer, the lab's answerer key, then the release
  (infra `docs/NOTARY-WITHHOLDING-DEPLOY.md`).
- **7-7** release. **Released by the user on 2026-10-14** as notary v1.0.0 and gateway v1.1.0. The
  deploy checked out on migrations, the court, the answerer and the Grafana rules.

  Every pass over the withholding court still failed, from the first minute. Each `eth_getLogs` asked
  sepolia.base.org for 2 000 blocks, and the endpoint answers at most 201
  (`eth_getLogs is limited to a 200 range`, HTTP 413). anvil and the fork test answer any range, so
  neither could have shown it. Nothing was missed, because registration had not opened.

  The fix, green in `bin/ci-go` (85.4%):
  - the span is `KINETIX_NOTARY_RPC_LOG_RANGE`, required, because the limit belongs to the endpoint;
  - `Answerer.CheckScan` asks for one full span ending at the head before notaryd starts, so an endpoint
    that refuses it fails the deploy, not every pass;
  - `kinetix_notary_withholding_passes_total{outcome}` counts each pass, and infra's
    `withholding-stalled` rule pages when thirty minutes pass without a served one.

  **Released 2026-10-14 as v1.0.1** (infra `a30c510` first), deployed 06:11 UTC, and proven on the server
  by 09:15 UTC:
  - `notaryd starting` logs `log_range 200`; passes 713 `served`, 0 `failed`; no gauge reads `NaN`;
  - epochs 37 onward carry account ids, and every account's secret is kept;
  - a fresh customer's epoch 41 proof through Kong is `inclusion.v2` with its `secret`, and
    `notary-verify` reports `included` and `attributable`;
  - a holder-signed registration through Kong returned one registration from epoch 42; the same body
    naming another holder was refused 403.

#### Step 8 — threshold anchoring (`AnchorCouncil`). **Ruled 2026-10-14 (decisions 18–21); live 2026-10-15** (the council the ledger's only anchorer).

One hot key can anchor anything today. The council holds `ANCHORER_ROLE` and anchors an epoch only
with k-of-n EIP-712 signatures over its header (domain-separated by chain id and ledger address).

- Members: the notary's signer and **`notary-watch`**, an independent process with its own key that
  checks the epoch against payment's own record of the snapshot, rebuilds the tree on its own, and
  signs only when the root matches. Long term the watcher belongs to a different operator.
- Anyone holding enough signatures may submit (a relayer, or step 10's sponsored transaction).
- Membership and threshold change through the admin with a delay, and never below 2-of-n.

**Design settled 2026-10-14.**

*The approval.* EIP-712, domain `{name: "Kinetix Anchor Council", version: "1", chainId,
verifyingContract: the council}`, type `EpochApproval(address ledger,uint64 number,bytes32 digest)`,
where `digest` is the epoch digest of §4. The digest already commits to the scheme, the number, the
previous digest, the root, the total, the leaf count and the snapshot time, so signing it signs the
whole header, and the council computes it with the ledger's own `EpochChain`.

*Anchoring.* `anchor(header, signatures)` is permissionless. Signatures come in ascending order of the
address they recover to, so a repeated signer is refused without a mapping; every one must recover to
a current member, and at least `threshold` must be given. Then the council calls `ledger.anchor`, and a
ledger refusal reaches the caller with its own revert data, so notaryd decodes `EpochAlreadyAnchored`,
`SnapshotInFuture` and the rest exactly as it does today. `EpochApproved(number, digest, approvers)`
records who approved.

*Membership.* A list of up to 16 distinct non-zero members and a threshold between 2 and the list's
length, set at construction. The admin schedules a whole new membership; after 7 days anyone may
activate it, and the admin may cancel it before. A change is whole, never a single add or remove, so
no intermediate state can drop below the threshold. Approvals are not stored: they are checked against
the membership at the moment of `anchor`, so an approval from a removed member stops counting at once.

*Admin.* `AccessControlDefaultAdminRules`, the ledger's admin keystore, a 3-day transfer delay — the
court's arrangement. The council has no other role and holds no money.

*No equivocation fault.* The first draft recorded two approvals for one epoch as a fault. Two things
argue against it. The ledger accepts one header per number and refuses the other (`EpochConflict`), so
a second approval can never become a second history; there is no harm left to punish. And an honest
watch can approve two different trees of one snapshot, because the notary draws fresh salts for every
build: convicting it would punish honesty. The court is immutable and its two fault kinds are fixed,
so there is nothing to record there either. What a stranger can check is in `EpochApproved`.

*The snapshot record* (decision 20, payment's work). The liability stream's header gains a
`snapshot_id`, its trailer a `digest`: SHA-256 over
`"kinetix.payment.liability-snapshot.v1" ‖ asOf seconds:u64 ‖ asOf nanos:u32 ‖ currency ‖ count:u64 ‖
total:u64`, then every account sorted by kind and then by the holder's UTF-8 bytes, each as
`kind:u8 ‖ holderLength:u8 ‖ holder ‖ amount:u64`. The order is canonical so the watch can rebuild it
from rows in any order. Payment writes the record before it sends the trailer, so a snapshot the notary
can use always has one, and serves it from a new `payment.v1.LiabilityRecordService` that the watch
alone may call. notaryd recomputes the digest on every snapshot and refuses one it cannot reproduce,
so a Java↔Go encoding drift stops the first epoch, not the first watch call.

*The watch.* Stateless: no database. It serves `notary.v1.EpochWatchService/ApproveEpoch`, a client
stream of one proposal (snapshot id, header) then batches of accounts (kind, holder, amount, secret,
salt), over the mesh, to the `notary` identity only. It approves only when:
- payment's record for the snapshot id exists, and its time, currency, count and total match the
  proposal, and the rows hash to its digest;
- the header's number is the ledger's `nextEpoch`, its previous digest the ledger's `head`, and its
  snapshot time the record's, in whole seconds;
- the tree it rebuilds — account ids from the secrets, commitments from the salts, leaves sorted by
  commitment — has the header's root, total and leaf count.

A disagreement is `FAILED_PRECONDITION` with the reason; an outage of the chain or of payment is
`UNAVAILABLE`.

*notaryd.* In the `built` state it signs its own approval with the anchorer key, asks the watch for
its approval, and sends `council.anchor` from the anchorer. Approvals are not stored: ECDSA here is
deterministic, so asking again after a restart returns the same signature. A watch refusal abandons the
epoch, because the notary and the watch disagreeing is an incident; an outage is retried. It refuses
to start unless the council reads its ledger, holds the anchorer role there, counts the anchorer key
as a member and hashes an approval as Go does. The epoch row keeps its `snapshot_id` (migration 00004).

*Cutover*, after the council is deployed: the ledger admin grants the council `ANCHORER_ROLE`, the
notary is released in council mode, its first epoch through the council goes final, and only then is
the anchorer key's own role revoked. From then on a stolen anchorer key anchors nothing alone.

*Work order*, each part green before the next:

- **8-1** the contract: `AnchorCouncil`, its interface, `PendingMembership`, `ApprovalHashing`; unit,
  fuzz and handler-based invariant suites; the deploy script through the CREATE2 factory and
  `bin/rehearse-deploy`; Slither and Aderyn.
- **8-2** Go↔Solidity: the approval digest in `internal/keccak`, a fixture signed in Go and verified by
  Foundry; a fork test where a council beside the live ledger anchors its next epoch with two
  approvals and refuses one.
- **8-3** `kinetix-contracts`: the stream's new fields, `payment.v1.LiabilityRecordService` and
  `notary.v1.EpochWatchService`, each new service in a file of its own. The user publishes.
- **8-4** payment: the digest, the record (a Liquibase changelog from the entity), the record service,
  the allow-list entry.
- **8-5** `notary-watch`: the use case, the gRPC server with a peer check, the payment record client,
  config, metrics; tested over real mTLS.
- **8-6** notaryd in council mode: the snapshot digest check, the approval round, the council writer,
  start-up checks, migration 00004, an end-to-end run on anvil with a real watch.
- **8-7** the estate: `bin/deploy-council`, `bin/inspect-council`, `bin/new-watch-key`,
  `bin/hand-over-anchoring`, the lab council and watch, compose, the `notary-watch` PKI leaf,
  allow-lists, Grafana rules, the deploy document. Then the council on Base Sepolia and the cutover,
  by the user's keystores.

**Status, 2026-10-14: 8-1, 8-2 and 8-3 built and green.** The contracts gate passes except the deploy
rehearsal, which waits for the watch's address: 203 tests on the CI profile, 12 fork tests, Slither 0,
Aderyn nothing new. `bin/ci-go` exits 0 at 85.5%. What each part proved:

| Part | Proof |
|---|---|
| 8-1 | 36 unit, 5 fuzz (20 000 runs on CI) and 5 invariants over 131 072 calls on CI. Three planted mutations each turned suites red: no threshold check (unit, fuzz, invariant), a repeated approver counted twice (unit, fuzz, invariant, once the handler and the fuzz learned to repeat one) and, below, a renamed EIP-712 field. `anchor` through the council measured at most 183k gas against ≈144k for the ledger alone. The deploy script went through the CREATE2 factory on anvil with chain id 84532, and refuses a network with no council members |
| 8-2 | Go signs four approvals, from epoch 0 to `uint64` max; Foundry rebuilds the same epoch digest, the same EIP-712 digest and the same signer. Renaming `number` to `epoch` in Go's type string turned Foundry red. On a fork of Base Sepolia at block 48063800, a council beside the live ledger takes the anchorer role from the admin, the old hot key is refused, two approvals anchor the next epoch, one approval anchors nothing, and the ledger's own refusals come back through the council unchanged |
| 8-3 | `kinetix-contracts`: `snapshot_id` and `digest` added to the liability stream, `payment/v1/liability_record.proto` and `notary/v1/watch.proto`, each a new service in a new file; `buf lint` with moneylint, `buf format` and `tools/breaking-gate` against `main` all pass in a container. Published by the user on 2026-10-15 as **v1.2.0** = `48d7806308461c877e81bede8f254664b76cbc09` |

Building it changed one shared piece: the EIP-712 domain in `internal/keccak` is now one value per
contract (`courtDomain`, `withholdingDomain`, `councilDomain`), because golangci's `unparam` saw every
caller pass the version `"1"`. The statement and withholding digests are unchanged, and their fixtures
still pass byte for byte.

**Status, 2026-10-14 later: 8-4, 8-5 and 8-6 built and green.** `bin/ci-go` exits 0 with golangci 0 and
coverage 86.0%; payment's suite passes 248 tests with 0 skipped against a freshly migrated Postgres 16.
Both run against the new protos copied into each repository's ignored `.contracts/`; the pins move once
the user tags `kinetix-contracts`. What each part proved:

| Part | Proof |
|---|---|
| 8-4 | payment records every snapshot before its trailer and serves the record to the watch alone. The snapshot digest has a known-answer vector computed by an independent Python implementation, and Java and Go both reproduce it; dropping the canonical sort in Java turned both of its known-answer tests red. A snapshot whose record cannot be written ends without a trailer. The table came from the entity through `diffChangeLog`; seven unrelated drift changesets the diff also proposed were left out |
| 8-5 | `notary-watch` approves only when payment's record, the ledger's head and its own rebuild of the tree all agree: fifteen kinds of disagreement are each refused, and an outage of payment or the chain is reported as one, never as a refusal. Over real mTLS a 1 201-account proposal crosses whole, a caller that is not `notary` and a server that is not `notary-watch` are both refused, and a malformed proposal is the caller's mistake |
| 8-6 | notaryd verifies every snapshot's digest itself, keeps its `snapshot_id` (migration 00004), and in `built` asks the watch, signs its own approval and sends `council.anchor`. A watch refusal abandons the epoch with the watch's reason; an outage leaves it built and asks again. An epoch already on chain is moved on without asking the watch, which a planted mutation proved necessary. On anvil, with a real `AnchorCouncil`, the anchorer role handed to it and a real watch in process, the daemon anchors, finalizes and every account verifies; a payment record that disagrees with what notaryd saw makes the real watch refuse and the epoch is abandoned with no transaction sent |

Found on the way: the gRPC client counter never counted a successful client-streaming call, because it
counted only on `io.EOF`, which such a call never returns. The daemon test that asserts the watch call
caught it.

**8-7, as far as it can go before the production council exists (2026-10-14).**
- **This repository:**
  - `bin/new-watch-key`, `bin/deploy-council`, `bin/inspect-council`, `bin/fund-admin` (the admin has
    never held gas), `bin/hand-over-anchoring` and `bin/retire-anchorer`.
  - `retire-anchorer` takes the transaction of an anchor through the council as its proof, because the
    public endpoint will not scan logs wide enough to find one.
  - The admin scripts were rehearsed read-only against Base Sepolia. That rehearsal caught the pinned
    `cast` printing a receipt's status as `true`, which the first draft compared with `1 (success)` and
    would always have refused.
  - shellcheck is clean at warning level.
- **Infrastructure:**
  - The lab script deploys the lab council at `0xF92c3A6eEC782530b5907315E152669256a237a4`, hands it the
    role and revokes the anchorer key's. Rehearsed twice on a throwaway chain: it landed at its address,
    and the second run changed nothing.
  - The lab runs `kinetix-notary-watch-service`. The name is the one the PKI's DNS SAN gives it, which
    the client's hostname check needs.
  - Payment's allow-list names `notary-watch@payment.v1.LiabilityRecordService`, and `bin/kinetix-pki`
    knows `notary-watch`.
  - `docs/NOTARY-COUNCIL-DEPLOY.md` and `docs/NOTARY-LAB.md` are written; `compose-check` passes, dev 47
    and prod 33.
- **The watch key (user, 2026-10-14):** `0x9ae6E5Eae21b41Dc0744Ea9FB3d2AB4D33F106C6`, in the prod store.
  The members are now in `deploy/base-sepolia.env`. On an anvil fork of Base Sepolia at block 48109915
  the deploy went out as one CREATE2, landed at `0xf6c33dE8F884fCcfE523DdC43d454AD7023D791b` for
  2,797,275 gas, read back the right ledger, members, threshold and admin, and took the anchorer role
  from the impersonated admin for 51,486 gas.
- **Deployed by the user on 2026-10-14:** the council at `0xf6c33dE8F884fCcfE523DdC43d454AD7023D791b`,
  exactly where the rehearsal said, and the anchorer role handed to it (tx `0xf745…191e`, block 48110555,
  51,486 gas). The anchorer key holds the role too until the retirement step.
- **The watch's mTLS leaves (2026-10-14):** the production leaf from the production CA, as
  `spiffe://kinetix.local/service/notary-watch`, copied to the server with the notary's owners; the dev
  leaf and the lab key in the dev store.
- **Contracts v1.2.0 (2026-10-15):** payment and this repository are pinned to it. The Go generated from
  the tag is byte for byte the pre-tag copy's; payment passes 248 tests with none skipped, with its
  protobuf generated again; `bin/ci-go` exits 0, golangci 0, coverage 85.9%.
- **Found while adding the production lines, and fixed:**
  - the watch serves gRPC but counted nothing it answered, and infra's deploy gate requires
    `kinetix_grpc_server_calls_total` of every gRPC server. It now counts every call by method and code,
    seeded at zero; taking the interceptor out turns both daemon tests that assert it red;
  - it named itself `kinetix-notary-watch` in `kinetix_build_info`, where every process in the estate
    uses its compose service name. It is `kinetix-notary-watch-service` now, as its Prometheus target is.
- **Production compose (infra):** `x-notary-chain-env` in `compose.yaml` holds the chain, the contracts and
  the mesh notaryd and the watch share, and `x-notary-chain-prod` overrides it once for both, with the
  production council. `compose-check` passes, dev 47 and prod 34.
- **The release order changed** to infra A, payment, infra B, notary v1.1.0
  (`docs/NOTARY-COUNCIL-DEPLOY.md` step 8). A deploy recreates only its own service on infra's `main`,
  so payment admits the watch only if the allow-list lands before payment's release; and every deploy
  runs the metrics gate over every target, so the watch's target, its conformance entry and its
  `bin/prod-deploy` sibling (B) wait until payment is live.

- **Live (2026-10-15).** The first two notary v1.1.0 deploys failed in the deploy pipeline, not in the
  notary (`kinetix-infrastructure/docs/NOTARY-COUNCIL-DEPLOY.md` step 8). The third started the watch.
  - The first anchor through the council: `0xe7ea8010…d0dd`, 02:56 UTC, block 48132336, final.
  - The admin revoked the anchorer key's own role at 05:30 UTC, tx
    `0xdbb3106c66f1dbd4ed424ec59fb7cf40ef8fc845df56c80e50d68d638b2f4515`, 29,571 gas.
  - `bin/inspect-council` reads `council anchors true`, `anchorer key anchors false`. `bin/inspect-ledger`
    now checks both against `KINETIX_NOTARY_ANCHORER_ROLE_HOLDERS`, and a wrong list turns both lines
    `WRONG`.
  - The first anchor after the revocation went through the council too: epoch 57, 05:56:46 UTC,
    `0x22cdbc372e178b298fd8e47e20d96ad8cc1bb874d890cbf254b7c0e51b65279a`; the ledger's next epoch is 58.

#### Step 9 — ledger directory (`LedgerDirectory`)

A ledger succeeds another by sealing it (§6), but today only a config file says which one is current.
The directory is an append-only, on-chain list of generations — ledger address, genesis epoch,
genesis digest, scheme, registration time. Registering a successor on the same chain requires the
predecessor sealed and the successor's genesis equal to the predecessor's sealed head and next epoch,
both checked on chain; a successor on another chain is recorded with the values it declares and marked
as such. The v1 ledger is generation 0. Verifiers resolve the current ledger from the directory, so
their configuration shrinks to one address that never changes.

**Ruled 2026-10-15 (decisions 22–25).** The contract, `src/directory/LedgerDirectory.sol` (MIT):

- `Generation{chainId, ledger, court, withholding, genesisEpoch, genesisDigest, scheme, registeredAt}`,
  kept in an append-only list. The constructor registers generation 0.
- `registerSuccessor(ledger, court, withholding)`, admin only. It refuses unless no foreign successor is
  pending, the current generation is on this chain and its ledger is sealed, all three contracts announce
  their interfaces, the court judges the ledger, the withholding court stands beside the court, the
  successor is not sealed, and its genesis epoch and digest are the predecessor's next epoch and head.
- `proposeForeignSuccessor(chainId, ledger, court, withholding, scheme)`, admin only, under the same
  predecessor checks; the genesis comes from the sealed predecessor. `cancelForeignSuccessor()` by the
  admin, `activateForeignSuccessor()` by anyone once the 7-day delay has passed. After that, both
  registrations refuse with `DirectoryMoved`.
- Views: `generationCount`, `generation(i)`, `current()`, `pendingSuccessor()`, and `generationOf(epoch)`,
  which returns the latest generation whose genesis epoch is at or below the epoch. A predecessor sealed
  before it anchored anything therefore owns no epoch. An epoch below generation 0's genesis reverts
  `EpochBeforeDirectory`.

**Work order:**
- **9-1** the contract: unit, fuzz, invariant (append-only, continuity of every local generation, at
  most one foreign generation and only last), the deploy script through the factory, the deploy
  rehearsal, and a fork test over the live ledger and courts that seals the ledger on the fork and
  registers a successor.
- **9-2** Go: the directory binding and adapter, `notary-verify` on the directory with a verdict for an
  epoch on another chain, notaryd's start-up check, integration tests on anvil.
- **9-3** the estate: `bin/deploy-directory`, `bin/inspect-directory`, `bin/register-successor`, the lab
  directory, compose, the deploy document. Then the directory on Base Sepolia by the user's deployer
  keystore, and the release.

**Status, 2026-10-15: 9-1 and 9-2 built and green; 9-3 written and rehearsed, waiting on the user's deploy.**

| Part | Proof |
|---|---|
| 9-1 | `LedgerDirectory` with 28 unit tests, 4 fuzz tests and a handler-based invariant suite. The suite has five invariants over 16 384 calls: the list only grows; every local successor continues its sealed predecessor; only the last generation may live elsewhere; every epoch has one generation; no action had an outcome the rules forbid. The deploy goes through the CREATE2 factory, and `bin/rehearse-deploy` proves it on chain id 84532. On a fork of Base Sepolia, the live ledger, sealed on the fork, was succeeded and the successor took the next epoch. Five planted mutations each turned suites red: no genesis check, an unsealed predecessor, `generationOf` off by one, a foreign successor a second early, and a court of another ledger. `bin/ci-contracts` exits 0: 243 tests on the CI profile, 15 fork tests, Slither 0, Aderyn's `centralization-risk` accepted. The one `unchecked-return` it raised was fixed in code |
| 9-2 | `notary-verify` reads the directory alone, resolves the generation holding each epoch and verifies through that ledger and court. It reports `moved` (exit 9) for an epoch held on another chain and `wrong_directory` (exit 8) for anything that is not a directory. notaryd's `DirectoryCheck` refuses to start unless the current generation is its own ledger and courts. On anvil: epoch 0 in a sealed ledger and epoch 1 in its registered successor both verify through the directory, a statement signed for the successor's court is attributable, and a proof naming the wrong generation is `wrong_ledger`. A successor that does not continue the sealed head is refused with `GenesisMismatch`. `bin/ci-go` exits 0, golangci 0, coverage 85.9% |
| 9-3 | `bin/deploy-directory`, `bin/inspect-directory` and `bin/register-successor`, rehearsed on a fork of Base Sepolia at block 48148939. The production directory lands at `0x6C71162b46620Aa835681ADCDC4DB4ab9bB84ac8` for 2,760,697 gas. The rehearsal caught two faults in the scripts: tuple output cut by commas, and refusals printed as bare selectors. The lab script deploys the lab directory at `0x57C3DF1081c5B42D6104fd160df9BEeBcd573C6c`; it was run twice on a throwaway chain, and the second run changed nothing. Compose names both directories; the rendered models differ from the last commit by that one variable alone. `kinetix-infrastructure/docs/NOTARY-DIRECTORY-DEPLOY.md` is the sequence |

**Live, 2026-10-15.** The user deployed the directory at `0x6C71162b46620Aa835681ADCDC4DB4ab9bB84ac8`, where the
rehearsal said, and released notary v1.2.0 at 14:54 UTC. Its notaryd passed the directory check at start
and anchored epoch 66 through the council at 14:59 UTC. The withholding check's proof of epoch 41 verifies
`included` and `attributable` through the directory alone; the ledger's address given as the directory
is `wrong_directory`.

Building it split `config.Verifier`: the daemon and the watch read `config.Chain` (ledger, court) and
`notary-verify` reads `config.Verifier` (directory). The address, chain id, URL and timeout parsing is
shared.

#### Step 10 — sponsored gas (`NotaryPaymaster`, ERC-4337)

A holder who wants to request a statement or prove a fault should not need testnet ETH first. The
paymaster sponsors user operations whose calls go only to the court's request and fault functions and
the council's submit, with a per-sender rate limit and a budget whose runway is monitored like the
anchorer's (§7). The EntryPoint version is pinned and its presence on the network checked before
deploy, the way the CREATE2 factory was.

**Ruled 2026-10-15 (decisions 26–29), which narrow the paragraph above:** no council submit, and no
per-sender rate limit, because an unstaked paymaster may keep no state.

**Revised 2026-10-16 (decisions 26 and 30):** the paymaster is staked with 1 ETH and funded with 0.1 ETH.
The rate limit stays out: a staked paymaster may read the courts, so instead of counting a holder's calls
it refuses every call the court would refuse, which leaves a holder nothing to waste the deposit on.

The contract, `src/paymaster/NotaryPaymaster.sol` (MIT), built on the interfaces of eth-infinitism's
`account-abstraction` v0.8.0 (a Soldeer git dependency pinned to commit `4cbc0607…`). Its validation reads
the operation, the sender's code, its own immutables and, since 2026-10-16, the storage of the two courts
and the ledger. It writes nothing and returns no context:
- the caller is the EntryPoint, and the operation's `maxCost` is at most the ceiling;
- the sender is an EOA delegated to the pinned `Simple7702Account`, checked by its code hash;
- the call is that account's single `execute(target, 0, data)`, to one of the two courts:
  - `requestStatement` and `recordUnanswered` need the registration's holder to be the sender and the
    registration to be signed by the court's statement signer. `recordUnanswered` carries no signature,
    so the holder puts it in `paymasterData`;
  - `proveBadProof` and `proveEquivocation` need every statement signed by that signer.
- the court would take the call, checked against the court's own state (2026-10-16):
  - a request needs an anchored epoch no earlier than the registration, no request of the registration
    still open, and none already made for that epoch; it is valid until the court's window closes;
  - an unanswered record needs an open request, and is valid only after its answer is due;
  - a bad proof must not hold against the anchored root, two statements must contradict each other, and
    neither fault may be proven already;
  - every epoch must be the signer's.
- `sponsorship(sender, callData, paymasterData)` runs the same checks as a view and returns the window,
  so `notary-sponsor` asks before it signs.
- The admin alone stakes it (`addStake`, with the unstake delay fixed at deploy), unlocks the stake and
  withdraws it, and takes back with `withdrawBalance` any ether forced into the paymaster itself.
- The signer and both EIP-712 domains are read from the courts once, at deploy, and the constructor
  refuses courts whose own digests differ from its. A signer rotation therefore means a new paymaster.
- Funding goes straight to `EntryPoint.depositTo(paymaster)`. The admin, the ledger's admin key, may
  `withdrawTo` at any moment, which is the kill switch.

Production pins the EntryPoint `0x4337084D…f108` (code hash `0xa06a…029c`) and `Simple7702Account`
`0x4Cd241E8…66c9` (code hash `0x82c1…3ec4`), and caps an operation at 0.0001 ETH. Pimlico's fast price
on 2026-10-15 was 0.0077 gwei, so a million-gas operation costs at most 0.0000077 ETH.

**Work order:**
- **10-1** the contract: unit, fuzz and invariant suites, deploy script and rehearsal, and a fork test on
  the canonical EntryPoint and account.
- **10-2** Go:
  - the user-operation hash, with a differential fixture against `EntryPoint.getUserOpHash`;
  - the EntryPoint, paymaster and bundler adapters, and EIP-7702 authorizations;
  - `notary-sponsor` (`request`, `unanswered`, `prove`, `equivocation`);
  - notaryd's paymaster check and deposit gauge.
- **10-3** the estate: `bin/deploy-paymaster`, `bin/inspect-paymaster`, `bin/fund-paymaster`,
  `bin/stake-paymaster` (2026-10-16), the alerts, the deploy document. Then the deploy, and one sponsored
  request on Base Sepolia through a public bundler.

**10-1, 2026-10-15: built and green.**

| Part | Proof |
|---|---|
| 10-1 | 20 unit, 4 fuzz and 3 invariants over 16 384 calls, run through a real `EntryPoint.handleOps` with the official v0.8 EntryPoint and `Simple7702Account`:<br>• a holder holding no ether requests a statement, records a silence, and proves a bad proof and an equivocation;<br>• nothing else is sponsored: another target, another call, value, a batch, another delegate, another holder, a forged registration or statement, or a cost over the ceiling;<br>• the deposit moves only for sponsored operations and the admin, and the holder never pays.<br>On a fork of Base Sepolia, the canonical contracts carry exactly the pinned code hashes, and a holder's first delegated operation, with the `0x7702` marker, is sponsored. Six planted mutations each turned suites red. Slither 0; Aderyn's high `contract-locks-ether` was answered by removing the paymaster's payable `deposit`, and its literals were named. The pinned dependency's own compiler warnings are excluded by path only (`ignored_warnings_from`) |

**10-2, 2026-10-15: built and green.** `bin/ci-go` exits 0: golangci 0 issues, every layout tool clean,
85.9% of statements covered, the integration suite on anvil and Postgres with nothing skipped.

| Part | Proof |
|---|---|
| hash | `keccak.OperationHash` is the EntryPoint v0.8 EIP-712 hash, the 7702 marker included; `operation-v1.json` is re-hashed by the real `EntryPoint.getUserOpHash` in `DifferentialOperation.t.sol`. Swapping Go's verification and call gas turned Foundry red, and restoring it brought the fixture back byte for byte |
| adapters | `ethereum.Paymaster` reads every fact of the deployed paymaster and calls the ledger, the account and an empty address `not_sponsored`; `bundler.Client` speaks ERC-7769 with the standard library only, and a planted go-ethereum import fails `bin/check-boundaries` (exit 1) |
| `notary-sponsor` | end to end on anvil through a test bundler that sends `handleOps`: a holder with no ether requests a statement for epoch 0; the request is open, the holder still holds 0, their code is the delegation and the deposit paid. An early `unanswered` lands and reverts (exit 1), the paymaster paying. A registration the notary never signed is refused (exit 2) and costs the deposit nothing |
| notaryd | refuses a paymaster for another notary's courts (`not_sponsored`), and serves `kinetix_notary_paymaster_deposit_wei 1e+18` and `kinetix_notary_paymaster_current 1` beside a real anchor |

Two faults the unit suites could not see, both caught by the integration suite against deployed code: the
paymaster's interface id XORed the methods it inherits from `IPaymaster`, and the authorization's `yParity`
went out as `0x01`, not the quantity `0x1` (CLAUDE.md, traps).

**The staked revision, 2026-10-16: live and proven in production. Step 10 is complete.**

| Part | Proof |
|---|---|
| contract | 20 unit, 12 court-rule, 6 stake, 7 fuzz and 4 invariants, through a real `EntryPoint.handleOps`. Each court rule is refused exactly as the court refuses it, and the court is asked the same call in a snapshot to prove it would refuse. The equivalence fuzz drives requests, records, answers and waits for 8 steps, and bad proofs and equivocations with tampered fields, and asserts the paymaster sponsors exactly what the court takes and that every sponsored call succeeds. A request is sponsored at `closesAt` and refused a second later, a record refused at `dueAt` and sponsored a second later. |
| ERC-7562 | `rules/`, its own profile with a tracer: for each of the four calls, validation runs none of the banned opcodes, `GAS` only before a call, and reads no storage of the paymaster or the sender |
| controls | a planted `block.timestamp` turned the opcode test red, a planted role read the storage test, and a dropped `AlreadyRequested` check both the equivalence fuzz and the invariants |
| analysis | Slither 0, Aderyn nothing new. The payable `addStake` first drew Aderyn's high `contract-locks-ether`, which `withdrawBalance` answers (CLAUDE.md, traps) |
| Go | `bin/ci-go` exits 0: golangci 0 issues, 85.7% covered, nothing skipped. On anvil: the early `unanswered` is now `not_due` (exit 3) and costs nothing, and lands once the answer is overdue; an unanswered record without a request and a second request are `call_refused` (exit 3) and cost nothing; the adapter names the court's refusal for an unanchored epoch, an epoch before the registration, a sound proof, a forged registration and another holder's. notaryd serves `kinetix_notary_paymaster_stake_wei 1e+18` |
| fork | at block 48176348 the deploy lands at **`0xD38824A16EaB632dD09aa36D3Bcf245097480aeb`**; inspect fails on the missing stake and deposit and passes once both are there; `bin/stake-paymaster` refuses a keystore that is not the admin, stakes, says so again without sending, and refuses an unlocking stake |
| estate | the lab paymaster `0xD8C9EE907c5822e5c6b80A30a4Ff633F9ff1dEA5`, funded and staked by `scripts/lab/notary-ledger.sh` (rehearsed twice, the second run changing nothing); production's address in `compose.prod.yaml`; the alert `kinetix-notary-paymaster-unstaked`; `bin/rehearse-deploy` stakes the rehearsed paymaster as its admin |
| deployed | the user bridged 1.01 ETH to the admin and 0.11 ETH to the deployer from Ethereum Sepolia. A dry run from a fresh `forge build --force` against the live chain returned `0xD388…0aeb`. The deploy landed there through the CREATE2 factory, tx `0xa4be6dd8…07ee8`, block 48181606, 3,605,435 gas, 06:18 UTC. The admin's `addStake` of 1 ETH, tx `0x98cd0e25…4b407`, went out at 06:20 UTC, and the deployer's `depositTo` of 0.1 ETH, tx `0x3a49324f…df617`, at 06:23 UTC. `bin/inspect-paymaster` exits 0 with every line `ok`: staked, 1 ETH, delay 86400 s, no withdraw time, a deposit of 1,000 operations at the ceiling, Sourcify `match`. The admin and the deployer each kept about 0.01 ETH for gas |
| released | notary v1.3.0 (`1beede12`, digest `73bdf567…ce40`) deployed 06:41 UTC; `notaryd starting` names the paymaster, no warning or error since, and Prometheus reads a deposit of 0.1 ETH, a stake of 1 ETH and `kinetix_notary_paymaster_current` 1 |
| proven | from the withholding check's customer `0x8482…41e8`, holding 0 ETH and no code, `notary-sponsor request` for epoch 81 through Pimlico's public bundler, run from the released image on the server: `landed`, exit 0. Tx `0x05cba39c…5753f` (type 4, `handleOps`, block 48182630, 06:52:28 UTC); `UserOperationEvent` names the paymaster with success true and an actual gas cost of 0.0000013 ETH, which came out of the deposit. The holder still holds 0 ETH, and its code is now the delegation to `Simple7702Account`. notaryd answered the request 6 s later, tx `0x5417512c…07cd`, and the court reads the request as `Answered`. The same request again was `call_refused` (exit 3, already requested), refused before anything was signed. So a public bundler accepts both the 7702 encoding and the staked paymaster's reads of the courts |

The unstaked paymaster below was never deployed, and its address will never hold anything.

**10-3, 2026-10-15: written and rehearsed; waiting on the user's deploy.** Superseded by the staked
revision above.
- `bin/deploy-paymaster`, `bin/fund-paymaster` and `bin/inspect-paymaster`, rehearsed on a fork of Base
  Sepolia at block 48156690: the deploy refuses a wrong EntryPoint code hash, lands at
  **`0xDf8A57D9Ff2a04ec2F2c603f1fE8f7FEf0001B3a`**, the inspection fails on an empty deposit and passes
  after `fund-paymaster … 0.002`, and a second deploy is refused.
- Code hashes are computed from `cast code`, never `cast codehash`: the latter asks `eth_getProof`, and the
  fork answered it with the empty account's hash for accounts it had not loaded yet.
- infra: the lab paymaster `0x0bdaE93FE03499758683B5921438883117Bbc1eA` from `scripts/lab/notary-ledger.sh`
  (rehearsed twice on a throwaway chain), production's address in `compose.prod.yaml`, the alerts
  `kinetix-notary-paymaster-deposit` and `kinetix-notary-paymaster-stale`, and
  `docs/NOTARY-PAYMASTER-DEPLOY.md`.
- Left for the user: deploy, fund, inspect, commit infra, release notary. Then one real sponsored request
  through a public bundler, from the withholding check's customer, who holds no ETH.

### Step 11 — verification without Kinetix

A browser verifier in the public docs: it reads the chain directly from the visitor's browser, so
checking a proof involves no Kinetix server at all. It resolves the current ledger through the
directory and shows any fault recorded against the epoch. The CLI image is published for auditors.

Decisions 31–35 (2026-10-16). Work order:
- **11-1, the verifier in WebAssembly.**
  - `verifycli` verifies a document from a reader, so no file system is needed.
  - `cmd/notary-verify-wasm` exposes `kinetixNotaryVerify(document, settings)` as a promise of the CLI's
    exit code and report lines.
  - `bin/build-verifier-wasm` builds it reproducibly.
  - Bun joins the toolchain, because Go turns `fetch` off under Node.
  - An integration test on anvil runs the WebAssembly under Bun beside the native CLI, for every verdict,
    and requires the same exit code and the same lines.
- **11-2, the CLI image and the release.**
  - `kinetix-notary-cli` is multi-arch and built for each platform by cross-compiling, never by emulation.
  - The release workflow tags it with the release, and attaches the WebAssembly, `wasm_exec.js` and
    `SHA256SUMS` to the GitHub release.
- **11-3, the docs.**
  - The Verify section: the how-to with the Svelte island that loads the pinned WebAssembly, and the
    explanation, both in English and Indonesian.
  - A build step that downloads the release asset and checks its hash.
  - Tests, and a check of the page at 390 px.
- **11-4, the user and production.**
  - The user releases the notary and makes the package public.
  - The docs pin that release.
  - A production proof is verified in a real browser against Base Sepolia.

**Step 11 done, 2026-10-16: released as v1.4.0 (`7d292fa`) and live in production; the docs commit is the user's.**

| Part | Proof |
|---|---|
| WebAssembly | `notary-verify.wasm` for `v1.4.0`: 44,203,656 bytes, SHA-256 `289ae0c0…1d04`; two builds with separate caches gave identical bytes. `browser_verifier_test.go` runs the CLI binary and the WebAssembly under Bun on the same documents on anvil and requires identical exit, stdout and stderr: attributable inclusion, a statement another key signed, a provably false statement, a bundle, never anchored, not final, a ledger named as the directory, a malformed document, a missing timeout; an endpoint that does not answer is `unknown` in both. A planted rewrite inside the WebAssembly turned the bundle case red. `bin/ci-go` exits 0: 0 lint issues, the `js/wasm` package linted on its own, 85.7%, nothing skipped |
| CLI image | `cli/Dockerfile` built here for arm64 and amd64 without emulation: each holds static `notary-verify` and `notary-sponsor` for its own architecture, user 65532. The arm64 image ran the docs' command on the production proof of epoch 41 from `/dev/stdin`: `included`, attributable, 0 faults, exit 0. The `cli` job of Build and push and the `cli` and `verifier` jobs of Release all passed on their first run (`7d292fa`). `kinetix-notary-cli:v1.4.0` is pulled anonymously, an index of linux/amd64 and linux/arm64, and the published image ran the docs' command on the epoch 41 proof: `included`, attributable, 0 faults, exit 0 |
| docs | the Verify section in `kinetix-docs`, two pages in English and Indonesian, and the island that loads the pinned WebAssembly. 103 tests, `astro check` and svelte-check clean, translations current, no leaks, build and links valid. In headless Chromium against Base Sepolia, the epoch 41 proof was `included` in 8.1 s and the same proof with its amount changed was `rejected`; at 390 px nothing overflowed |
| release | `v1.4.0` at `7d292fa`. The `SHA256SUMS` CI built on amd64 lists exactly the pin built on the laptop on arm64, and a fresh clone of the tag running the docs' own command printed the same two hashes. Deploy run 10:01 UTC recreated notaryd and notary-watch at the tag's linux/amd64 digest `ce2c06ff…`, healthy, every observability gate passed, both name `7d292fa` |
| CI | the first PR run went red in Gate go: the browser test was the suite's first `go build`, and git refused the runner-owned checkout. The same cause had emptied `bin/check-boundaries`' import rules on GitHub since CI began (CLAUDE.md, Traps). Fixed before the merge; the merged run has no `error obtaining VCS status` line |
| docs pin | `kinetix-docs` pins `v1.4.0`, downloads it from the release and checks the bytes; every notary permalink moved to `7d292fa` (the linked files are unchanged since `1beede1`), and the entry point, `RunDocument`, `bin/build-verifier-wasm` and the release job are linked. All docs gates green against the downloaded release |

### Step 12 — phase 2: evidence anchoring. **Live 2026-10-18** (notary v1.5.0, `EvidenceRegistry` 0xbF1F…F480).

Decisions 36–40 (2026-10-16). Three evidence streams, each built in the service that owns the records and
anchored by one registry here. Nothing goes live until all three are coherent.

What the survey found on 2026-10-16, before anything was designed:

- **Order** keeps no history. `OrdersChangedSince` returns each order's current row; two transitions
  between polls collapse into one, and the return lifecycle never touches `orders.updated_at`, so it is
  invisible in the feed. `SHIPPED` and `REFUNDED` are never assigned anywhere. No actor is recorded.
- **Any feed keyed on a value written before commit can skip a row for ever.** Reproduced on Postgres 16:
  a transaction that writes first and commits later lands behind a cursor that already moved past it.
  `OrdersChangedSince` keys on `updated_at` set by the application before commit, so search's order
  read model is exposed to the same loss. A cursor on `(tx xid8, seq)` that serves only
  `tx < pg_snapshot_xmin(pg_current_snapshot())` returned both rows in the same reproduction, in order.
- **Proof of delivery** has no hash anywhere. The app PUTs the picker's bytes to a presigned MinIO URL that
  can be replayed for 600 s; the bucket is versioned but has no object lock. The signature is a PNG data
  URI in matching's database. `delivered` is terminal. Only the backoffice reads POD.
- **MinIO at the pinned `RELEASE.2025-09-07T16-13-09Z`** (built from source, because quay.io now refuses
  the image) enforces what decision 38 needs: a presigned PUT that signs `x-amz-checksum-sha256` and
  `if-none-match: *` accepts the right bytes, answers a replay `412`, other bytes `400
  XAmzContentChecksumMismatch`, and a request without either header fails its signature; `HEAD` with
  `x-amz-checksum-mode: ENABLED` returns the stored SHA-256. Upstream is archived (decision 40).
- **Review** only ever creates reviews (no edit, hide or delete route), serves no gRPC, has no contract,
  and no PHP gRPC server exists in the estate. RoadRunner's gRPC plugin hands PHP only `:peer.address`
  and `:peer.auth-type`, never the client certificate, but it runs named Go interceptors
  (`grpc.interceptors`), so the SPIFFE allow list per RPC becomes a Go interceptor compiled into
  review's own RoadRunner binary. Review's HTTP side stays FrankenPHP.

The design:

- **Every owner writes an append-only evidence log in the same transaction as the change it records**,
  with `seq` (identity) and `tx xid8 default pg_current_xact_id()`, and serves it through a feed whose
  cursor is `(tx, seq)` and which returns only rows below the snapshot's `xmin`. Order logs status
  changes and the return lifecycle with the actor; matching logs a delivery (photo SHA-256, signature
  SHA-256, `delivered_at`, driver); review logs each review and driver rating as written.
- **Commitments, not content.** A leaf is a domain-separated keccak of the record and a salt the notary
  generates and keeps, as for liabilities: nothing on chain names a buyer, an order or a review's text,
  a guess of a low-entropy field cannot be confirmed, and erasing a record's salt erases what its anchor
  can say about it.
- **`EvidenceRegistry`**: immutable, MIT, one contract with a hash-chained sequence of batches per stream
  (order lifecycle, delivery, review), so an outage of one owner stalls only its own stream.
  `anchor(stream, header)` accepts exactly the stream's next batch, as the ledger accepts the next epoch
  (identical re-send `BatchAlreadyAnchored`, different `BatchConflict`); `verifyInclusion` is a view that
  reverts `UnknownBatch` for a batch that does not exist; `seal` hands every stream's head to a
  successor. The anchorer role is held by the existing anchorer key. The registry is its own trust root,
  configured like the directory.
- **Proofs** come from the notary's proof API to the parties of the record only (buyer, merchant, driver,
  author), with their access token, and `notary-verify` checks them against the registry. For a delivery
  the proof carries the photo's SHA-256, so anyone holding the photo can show it existed unchanged when
  the batch was anchored.
- **The courier app hashes before upload.** `pod_upload` carries the SHA-256; matching records it and
  signs it into the presigned PUT with `if-none-match: *`; at `delivered` matching compares `HEAD`'s
  checksum with the recorded one and hashes the signature itself. Matching refuses a `pod_upload`
  without a hash, so the app release reaches couriers before matching's.

Work order:

- **12-1 contracts**: `order.v1` lifecycle feed, a delivery-evidence feed served by matching, a new
  `review.v1` package with its feed. Additive only. The user publishes the version.
- **12-2 infrastructure**: `kinetix-minio` and `kinetix-minio-mc` built by CI from upstream
  `RELEASE.2025-10-15T17-29-55Z` and published to GHCR, re-checked against the four behaviours above;
  the three feeds granted to `notary` per RPC; review's gRPC service, its handshake rows, scrape and
  alerts.
- **12-3 order**: the lifecycle log, written at every transition and return step, with the actor; the
  feed; tests including a late commit.
- **12-4 matching and the courier app**: hash, checksum-locked write-once upload, `HEAD` check, the
  delivery log and its feed.
- **12-5 review**: RoadRunner gRPC with the SPIFFE interceptor, the log and its feed.
- **12-6 notary**: `EvidenceRegistry` (Foundry tests, Slither, Aderyn, rehearsed deploy), the three feed
  readers, salts, batching and anchoring, the proof API, `notary-verify` for evidence.
- **12-7 docs**: the Verify section gains evidence.
- **12-8 the user**, in this order, because each step's gates run against what is already live:
  1. publish contracts, and move every pin to the new tag, the notary's `bin/sync-contracts` included;
     then run the whole lab on the new images before anything is released: every image fetches its
     contracts inside the build from the pinned tag (the notary's even excludes a local `.contracts`), so
     no image of step 12 can be built before the publish, and the lab is the one place the real owners'
     feeds, `review-grpc`, the handshake matrix and the metrics conformance meet notaryd;
  2. release order (old infrastructure is fine);
  3. ship the courier app to couriers: the new matching refuses an upload without the photo's hash, and the new
     app works against the old matching;
  4. release matching, still on the old infrastructure: its new allow list parser and reflection rule must be
     live before the narrowed list is;
  5. land infrastructure's step 12 and deploy review straight after it, with no other deploy between: from that
     commit `metrics_conformance.sh` expects `kinetix-review-grpc`, which only review's deploy starts, and
     `trace_test.sh` asks matching for reflection as order under the narrowed list;
  6. deploy the evidence registry and release the notary, by `kinetix-infrastructure/docs/NOTARY-EVIDENCE-DEPLOY.md`:
     fund the evidence anchorer, deploy and inspect the registry, name it in `compose.prod.yaml`, commit
     infrastructure, then tag the notary.
  The evidence anchorer key goes into the prod store before step 5 (`bin/new-evidence-anchorer-key prod`):
  infrastructure's tree names it with `:?`, so a store without it stops every render. The MinIO images are
  independent of all of this and can go first (Waiting on the user).

**12-6 design (2026-10-17).** Evidence scheme v1, keccak over packed big-endian bytes like §4, with tags of
its own so no evidence hash can be read as a liabilities one:

| Item | Preimage | Bytes |
|---|---|---|
| record digest | `0x20 ‖ stream:u8 ‖ fields`, each string or byte field prefixed by its length as `u32` | — |
| leaf | `0x10 ‖ stream:u8 ‖ salt:bytes32 ‖ recordDigest:bytes32` | 66 |
| node | `0x11 ‖ left ‖ right` | 65 |
| empty (padding) | `0x12` | 1 |
| batch digest | `0x13 ‖ scheme:u16 ‖ stream:u8 ‖ number:u64 ‖ previousDigest ‖ root ‖ leafCount:u64 ‖ throughTransaction:u64 ‖ throughSequence:u64 ‖ closedAt:u64` | 108 |

Streams are `1` order lifecycle, `2` delivery, `3` review. Depth is canonical and padded as in §4. A batch
binds the feed position it reads through, so consecutive batches show the log was read without a gap. The record
digest is computed off chain (the Go builder and the verifier); the contract computes the leaf from the salt and
the record digest, binding the stream, and checks the path.

`EvidenceRegistry` (MIT, immutable): per stream a hash-chained sequence of batches. `anchor(header)`, anchorer
role only, accepts exactly the stream's next batch — previous digest equal to the head, at least one leaf, a close
time after the previous one and not in the future, a feed position after the previous one; an identical re-send
reverts `BatchAlreadyAnchored`, a different one `BatchConflict`. `verifyInclusion(stream, number, salt,
recordDigest, index, siblings)` is a view that returns false for a bad proof and reverts `UnknownBatch` for a batch
that does not exist. `seal(expectedHeads)`, admin only, stops every stream and emits each head for a successor,
whose constructor takes per-stream genesis values. ERC-165, `AccessControlDefaultAdminRules`, CREATE2 factory,
Sourcify, as every contract here.

Off chain: one gRPC reader per feed (order, matching, review-grpc) over the mesh; a collector that pulls each feed
from its stored cursor, draws a salt per record and stores record, salt and leaf in one transaction with the
cursor; a daily closer that builds each stream's tree, records the batch and sends the anchor with the evidence
anchorer key, replacing a stuck transaction at the same nonce as the liabilities anchorer does; a proof API that
gives a record's proof only to its parties (a lifecycle event to its buyer, merchant and actor; a delivery to its
courier and to the order's buyer and merchant, read from the order stream; a review to its author, the reviewed
merchant or courier); and `notary-verify` for the evidence document, against the registry configured beside the
directory.

Work order for 12-6: (a) the contract, its tests, analysis and rehearsed deploy; (b) the Go scheme and the
differential fixture; (c) the three readers and the store; (d) collector and closer with the evidence anchorer;
(e) the proof API; (f) the verifier and the browser verifier.

**12-1, 2026-10-16: written and green in `kinetix-contracts`, uncommitted and unpublished.** `common.v1.LogCursor`
(`transaction_id`, `sequence`); `order.v1.OrderLifecycleService/LifecycleEventsSince` in a new
`order/v1/lifecycle.proto` (event kind, from and to status, return number and status, actor kind and
principal, buyer, merchant, `occurred_at`); `fulfillment.v1.DeliveryEvidenceService/DeliveriesSince` in a new
`fulfillment/v1/delivery_evidence.proto` (order number, driver, photo and signature SHA-256 as bytes,
`delivered_at`); a new `review.v1.ReviewEvidenceService/ReviewsSince` (review id, product or driver, order id,
subject, author, merchant, rating, comment SHA-256, `created_at`). New services in new files, so nothing
existing changes: `buf lint` with moneylint, `buf format` and `tools/breaking-gate` against `main` all pass
(buf v1.74.0, the version CI installs), and a planted lint fault and a planted field rename each turned their
gate red. The services build against a local copy first; the user publishes the version once all three
owners have run against it, and every pin then moves to that tag.

**12-3, 2026-10-16: built and green in `kinetix-order-service`, uncommitted.** `order_lifecycle_events` with
`LifecycleLog` at every transition and return step (checkout placed and paid by the buyer, cancellation by the
system, packed by the merchant, delivered by the driver, completed by the system or the admin who did it,
return opened and goods received by the merchant, resolved by the system, rejected by the admin); deferred
constraint triggers refuse a status committed without its event; the log refuses update, delete and truncate;
the migration writes a `BASELINE` event for every existing order and return; `LifecycleEventsSince` serves
`(xid8, sequence)` below the snapshot's `xmin`. The contract gained `LIFECYCLE_EVENT_KIND_BASELINE`. 380 tests
on a fresh Postgres, 0 skipped; a planted feed without the `xmin` filter, a planted delivery that skips the log
and a planted trigger that never fires each turned their test red. The admin completion and rejection endpoints
now refuse a token that names no principal.
Found on the way: neither CI ran `dotnet test` at all, the ten Postgres tests failed on a fresh database
(concurrent migrations), and GitLab's order jobs never run `bin/sync-contracts`. Both CIs now run the suite
against Postgres, with a test that fails when `CI` is set and the database is not.

**12-4, 2026-10-16: built and green in `kinetix-matching-service` and `kinetix-courier-mobile`, uncommitted.** The
app reads the photo, hashes it and sends `pod_upload` with `photo_sha256`; matching refuses one without it, names
the object `pod/<job>/<sha256>.jpg` and signs the PUT over `x-amz-checksum-sha256` and `if-none-match: *`, and
hands the app those headers to send. At `delivered` matching reads the stored checksum (`HEAD` with checksum mode)
and accepts only the hash in the name; the signature must be a PNG `data:` URI and matching hashes the decoded
PNG. `delivery_evidence` (photo and signature SHA-256, courier principal, `delivered_at`, one row per job, append
only) is inserted before the status in the same transaction, and a trigger refuses a job marked delivered without
it; `DeliveriesSince` serves it on `(xid8, sequence)` below the snapshot's `xmin`. No backfill: earlier photos were
never hashed when taken. matching's allow list now takes `peer@Service/Method`, and reflection lists every service
(it had left out `FleetRegistryService`). Matching 390 tests including nine against MinIO `RELEASE.2025-09-07` built
from source, credo and format clean, and Dialyzer adds only the unknown `GRPC.Server.Stream.t` every server there already
reports; the app 105 tests and `flutter analyze` clean. Planted
faults each turned their test red: a feed without the `xmin` filter, a delivery that skips the evidence, a trigger
that never fires, an upload address that stops signing the checksum (MinIO enforces the headers even unsigned, so
only the signature stops a client dropping them), and in the app a dropped header, a wrong hash and a replay read as
a failure. Matching's presign and the app's uploader were also run against each other on that MinIO: the photo was
stored under its hash, a replay answered as already stored, other bytes were refused `400`.
Found on the way: `flutter analyze`, which the app's CI runs, failed since the POD commit (the integration test's fake
repository lacked `attachPodPhoto`); `Dockerfile.test` generated three of matching's six packages. Both fixed. Owed to
12-2: `notary@fulfillment.v1.DeliveryEvidenceService/DeliveriesSince`, and `order` is still a bare entry on matching's
list, so it reaches every service there, the evidence feed included. Matching's build needs the published contracts:
`bin/generate-proto` stops while the pin (v1.0.15) has no `delivery_evidence.proto`.

**12-5, 2026-10-16: built and green in `kinetix-review-service`, uncommitted.** `review_evidence` (review id, product
or driver, order, subject, author, merchant, rating, comment SHA-256, `created_at`) is written by Postgres itself: a
trigger on each review table inserts the row in the statement that inserts the review, so the log says exactly what
the row says and no code path can skip it; triggers refuse update, delete and truncate on the log; the migration
wrote every existing review in creation order. `ReviewsSince` serves it on `(xid8, sequence)` below the snapshot's
`xmin`. The feed runs in `review-grpc`, a second process from the same image: our own RoadRunner v2025.1.15 built in
`rr/` from its plugins plus `kinetix_metrics` (its own `/metrics`: `kinetix_grpc_server_calls_total` seeded per method,
the scrapes it serves, `kinetix_build_info`) and `kinetix_peers` (SPIFFE id from the verified chain, allow list per RPC, refuses to start without one), and the
PHP worker `grpc-worker.php`. The suite moved from SQLite to Postgres in both CIs, and a Go job gates `rr/`. Review
127 tests on Postgres, also inside the CI runtime image (0 skipped); Go tests and `govulncheck` clean after bumping
grpc 1.79.2 to 1.84.0, otel to 1.46.0 and x/net to 0.60.0 (seven reachable advisories; x/net's fix is the only one,
eight days old). Planted faults each turned their test red: a feed without the `xmin` filter, a review table whose
trigger fires on update instead of insert, a log that admits `UPDATE`, a comment hashed in another encoding. End to
end on the built image with a throwaway mesh CA: the granted notary read both reviews and a page of one, a cursor
past every row was `InvalidArgument`, search was `PermissionDenied`, a caller with no certificate never completed
the handshake, the metrics counted each, and with the TLS block set to stop verifying client certificates the
granted notary itself was `Unauthenticated` — the interceptor fails closed. An empty allow list stops the process.
Found on the way: `spiral/roadrunner-grpc` 3.8.0 requires the vulnerable `spiral/roadrunner ^3.0` (held at 3.6.0);
`spiral/goridge` declares `ext-sockets`, which the pipes relay never touches (ignored at install, documented); and
review's lock carried three advisories, one a ReDoS reachable through `Cache-Control` — laravel/framework 13.34.0,
symfony/http-foundation 8.1.8, league/commonmark 2.10.3, `composer audit` clean. Owed to 12-2: the `review-grpc`
service (port 50060, metrics 2112, `./pki/review`, `APP_VERSION`), `notary@review.v1.ReviewEvidenceService/ReviewsSince`,
handshake rows using the contract's proto (no reflection), scrape and alerts. Review's image build needs the published
contracts: the pin (v1.1.0) has no `review.proto`.

**12-2, 2026-10-17: written and checked in `kinetix-infrastructure`, uncommitted.** `minio/Dockerfile` builds
`kinetix-minio` (MinIO `RELEASE.2025-10-15T17-29-55Z`, commit `9e49d5e7`, with mc for the healthcheck) and
`kinetix-minio-mc` (`RELEASE.2025-08-13T08-35-41Z`, `7394ce0d`) from upstream source, versions stamped as
upstream's release build does, cross-compiled for amd64 and arm64; `.github/workflows/minio-image.yml` builds each
architecture, scans it, pushes it, joins the manifest list and commits the digests to `deploy/{prod,dev}/minio.env`,
which render appends. Re-checked on the images built here: the spike's seven cases unchanged (right bytes `200`,
replay `412`, other bytes `400 XAmzContentChecksumMismatch`, a missing header `400 AccessDenied`, `HEAD` returns the
stored checksum), compose's provisioning script runs unchanged through the mc image, and matching's storage and
mesh tests pass against it (43). `kinetix-review-grpc` is in both compose files (50060, metrics 2112, review's leaf,
which already names `kinetix-review-grpc`), a companion of review in `bin/prod-deploy`, scraped as
`service=kinetix-review-service, process=grpc` like warehouse's gRPC process, and in `metrics_conformance.sh`; the
generic down alert covers it. For that, review's metrics plugin became `kinetix_metrics`, its own `/metrics` with
the estate's HTTP families (RoadRunner's metrics plugin could not count its own scrapes); the e2e on the rebuilt
image passed again. Allow lists: `notary@order.v1.OrderLifecycleService/LifecycleEventsSince` on order;
on matching, `order` narrowed from bare to its two services and `notary@fulfillment.v1.DeliveryEvidenceService` —
service entries, because matching v1.0.0 parses no method entry; matching now lets a listed peer use reflection,
which the probes need. `handshake_matrix.sh` asserts the three grants and the narrowing, and probes `review-grpc`
with the contract's proto. compose-check renders both environments (48 and 37 services), `promtool check config`,
shellcheck and the bats suite pass (observability's five fail only on the Mac's BSD `touch`, 12/12 on Debian). Not
run: the full lab stack, so the handshake matrix and conformance have not run against a live `review-grpc` in compose.
Later the same day the workflow's first run failed its scan: the archived releases compile in `grpc` v1.71/v1.72
(CVE-2026-33186) and, in minio, `amqp091-go` v1.10 (CVE-2026-77405, -77408, -77411), all CRITICAL with a fix that
no upstream release will carry. I had not scanned the local build before the first run. `minio/Dockerfile`
now raises them to the first fixed versions (grpc v1.79.3, amqp091-go v1.13.0) with `go get` before compiling,
and stops unless the module graph then holds exactly those. Both local images then scan clean under CI's policy,
and the behaviour checks are unchanged: the seven cases, the provisioning script, matching's 43.

**12-6, 2026-10-17: built and green in this repository, uncommitted; lab and production wiring written in
`kinetix-infrastructure`.** (a) `contracts/src/evidence/`: `EvidenceRegistry`, its structs and its two libraries
(`EvidenceTree`, `BatchChain`), and `DeployEvidenceRegistry.s.sol`, which deploys through the CREATE2 factory and
gives a successor a sealed predecessor's heads as its genesis. Tests: 23 unit, 2 fuzz, 7 deploy, 5 differential
against `test/fixtures/evidence-v1.json`, 5 invariants. `bin/ci-contracts` passes (329 tests on the CI profile,
fork 18, Slither nothing, Aderyn only `centralization-risk`, accepted as for every contract here), and
`bin/rehearse-deploy` deploys the registry on anvil with each network's chain id and checks its streams and its
anchorer. Controls: without the head check two tests went red, without the feed-position check two, and Go's
leaf tag changed to `0x14` turned four differential tests red. Operations: `bin/new-evidence-anchorer-key`,
`bin/fund-evidence-anchorer`, `bin/deploy-evidence-registry` (dry run by default) and
`bin/inspect-evidence-registry`, which also checks the liabilities anchorer lacks the role. (b) the evidence
types in `internal/notary` and the scheme in `internal/keccak`, checked against known answers computed
independently in Python. (c) `internal/feeds`, one reader per owner over the mesh (each dials only the service
it names), which refuses a record the notary cannot commit to: no position, a missing field, an unspecified
kind, a hash that is not 32 bytes, a NUL, a field over 64 KiB, a time outside the epoch; migration
`00005_evidence.sql` (a cursor per stream, a leaf table, one typed table per stream keeping time as seconds and
nanoseconds, batches and attempts, one batch in flight across all streams). (d) `internal/evidencing`: the
collector checks every page follows the log before anything is kept, and keeps a page with the cursor in one
fenced transaction; the closer builds a stream's batch when its last one closed a day ago, or at once when a
full batch is waiting, closes it at the chain's latest block time, and anchors it with the evidence key
through the shared `transactor` (nonce, fees, re-signing at the same nonce, broadcast, receipts), which
`Writer` now uses too. Evidence metrics: `kinetix_notary_evidence_{feed_up,collected_total,pending_records,
last_finalized_batch,last_finalized_closed_at_seconds,anchorer_balance_wei,anchor_cost_wei,runway_seconds,
fee_ceiling_hits_total}` and `kinetix_notary_abandoned_evidence_batches`. (e) `/api/v1/proofs/evidence/orders/{order}`
and `/reviews/{review}`, through Kong's existing `/api/v1/proofs` route: a record's proof goes to its parties only,
a stranger gets 404, a record whose batch is not final is counted as pending, and a record the database changed
since it was kept is refused, never proven. (f) `prooffile` evidence documents, `verifying.EvidenceVerifier`,
`verifycli.RunEvidence`, and both `notary-verify` and the browser verifier route an evidence document to the
registry. Verified: `bin/ci-go` passes (84.9%); the integration suite on anvil and Postgres serves the three
feeds over mTLS, anchors every stream and has the buyer's proofs verify `included` on chain, a copy with the
courier changed `rejected`; the browser verifier answered exactly as the CLI in eight evidence cases. Two
faults found and fixed while testing: a built batch whose send failed let the next stream try a second one,
and a missing registry setting exited `unknown` rather than `ExitConfig`. Controls, each compiling and each
red: the cursor's compare-and-set removed, the delivery party rule widened, the page check skipped, the header
digest check skipped, the hash length check loosened, a kept entry's recheck skipped, and a planted `feeds`
import in `evidencing` (the boundary check). Infrastructure: `compose.yaml` gives notaryd the three feeds, the
lab registry `0x529cBEB0d9c329a1430671d7805c04442926462B` and the evidence key; `scripts/lab/notary-ledger.sh`
deploys the registry with anvil's fifth account as the anchorer (rehearsed twice on a throwaway chain: all seven
contracts at their constants, the second run changed nothing); four Grafana alerts (`evidence-abandoned`,
`-feed-down`, `-runway`, `-stalled`, expressions checked by promtool); `docs/NOTARY-EVIDENCE-DEPLOY.md` is the
production runbook. The protos are generated from the local contracts copy: `bin/sync-contracts` still pins
v1.2.0, so CI goes green only once the contracts are published and the pin moves (12-8 step 1).

**12-8 step 1, 2026-10-17: contracts v1.3.0 published (`edc6c26`, all four registries) and every pin moved;
the whole path proven in the lab.** order, matching, review and the notary now pin v1.3.0 in `bin/sync-contracts`
(uncommitted), and each passes against the published tag: the notary's `bin/ci-go` (84.9%), matching 392,
order 380, review 127 (one skipped where the host has no `ext-grpc`, as before). Review's `generated/` from the
tag is byte-identical to the one generated from the local copy. The lab ran on images built from those trees:
- `kinetix-notary-ledger` deployed the evidence registry at its constant;
- notaryd's registry check passed, and all three feeds answered over the real mesh and allow lists;
- one record per stream was written in the owners' lab databases (review's through its own trigger) for an
  order whose buyer is a throwaway account registered with the lab identity;
- notaryd collected them, anchored one batch per stream in turn, and each was final about five minutes later;
- with that buyer's real token the proof API gave 2 proofs for the order (its event and its delivery) and 1 for
  the review, 401 without a token and 404 for an order they are not in;
- the lab image's `notary-verify` answered `included` for all three against the lab registry, the delivery
  proof's photo hash equals the SHA-256 of the photo's bytes, and the same file with the courier changed answered
  `rejected` (exit 1).
The handshake matrix passed every step-12 row (notary to the three feeds served, `search` and `order` refused on
them, `review-grpc`'s three cases) and metrics conformance every `review-grpc` check. Its other failures were
services not started for this run, and two things in the lab that predate step 12: its identity image is from
2026-10-09, before identity's per-RPC allow list (6 rows), and a lab-built notary names no version. The lab chain
also came up empty while the lab database still remembers epoch 0, so liabilities anchoring in the lab stops on
`ledger_diverged` until both are reset (docs/NOTARY-LAB.md, "Resetting"). Next: step 2, the order release.

**12-8 steps 2–6, 2026-10-17 and 18.** Order v1.1.0 and matching v1.1.0 were released green. Review v1.1.0
failed its deploy and was rolled back: FrankenPHP segfaulted on ext-grpc state its threads shared (review
`CLAUDE.md`, "ext-grpc is built once per process, under a lock"). Review v1.1.1 went live at 16:11 UTC on
infrastructure `c43c86a`, the infrastructure of step 12, with `kinetix-review-grpc` healthy. On 2026-10-18 the
evidence anchorer was funded with 0.001 ETH and the `EvidenceRegistry` was deployed at
`0xbF1F15F40cA8Dd63c070d223243febf6731EF480`, block 48260906. `bin/inspect-evidence-registry` passed every
line: interface, scheme 1, streams `[1, 2, 3]`, the admin and its delay, the evidence anchorer holding the role
and the liabilities anchorer not, unsealed, and Sourcify `match`. The address is in
`deploy/base-sepolia.env` and infrastructure's `compose.prod.yaml`. v1.5.0 was released on 2026-10-18 at
03:46 UTC, green, and its registry check passed. Order and matching then refused its feeds on every tick:
their containers predate step 12's infrastructure, which carries the notary's grants, and `bin/prod-deploy`
recreates only what it deploys. Review's feed was up. The fix is to re-run order's and matching's v1.1.0
releases; infrastructure's `docs/NOTARY-EVIDENCE-DEPLOY.md` now lists it in the sequence.

**Step 12 complete, checked in production on 2026-10-19 at 14:54 UTC.**
- Order and matching were re-released at their live digests on 2026-10-18, and since then all three feeds have
  answered: `kinetix_notary_evidence_feed_up` is 1 for `order_lifecycle`, `delivery` and `review`.
- The first order-lifecycle batch (0) was finalized after closing at 2026-10-18 03:58 UTC.
- Delivery and review have no batch yet. Neither has produced a record since its log began: no delivery with a
  hashed photo, and no review, has been made in production. A batch closes as soon as one exists.
- The evidence anchorer holds 0.000999 ETH, about 282 days of runway. The liabilities ledger is final through
  epoch 161, with about 685 days of runway.
- 12-7 is done: the Verify section's evidence pages were deployed in `kinetix-docs` `ea16b63`.

Production now runs order v1.2.0 and matching v1.3.0, after the shipping-fee release, and review v1.1.1. The public
Verify pages link the owners' logs at those tags.

### Step 13 — scheme v2: zero-knowledge proof of liabilities. **Live 2026-10-20** (notary v1.6.1; generation 2 from epoch 183, the cutover record in infra `docs/NOTARY-GENERATION-TWO.md`).

The total stops being public, and the proof gets stronger. Built with gnark on Base, decided on 2026-10-20
after a research round on Midnight/Compact (report: `reports/Notary ZK di Midnight Compact.md` in the estate
root). The user accepted the recommendation and asked that Compact be remembered for a future service.

**What the deployed contracts force (read 2026-10-20).** `LedgerDirectory` is immutable. `registerSuccessor`
accepts a same-chain successor only when:
- the current ledger is sealed;
- the successor reports `type(ILiabilityLedger).interfaceId` through ERC-165, and `genesisEpoch`/`genesisDigest`
  equal to the sealed head;
- its court reports `IStatementCourt` and judges that ledger;
- its withholding court reports `IWithholdingCourt` and stands beside that court.

A directory cannot be replaced on the same chain: a foreign successor needs another chain id and is one-way.
So generation 2 implements the v1 interfaces, **with scheme-2 meaning carried by the same signatures**:
- **`verifyInclusion(number, commitment, amount, index, SumNode[] siblings)`** keeps its meaning exactly ("this
  commitment with this amount is at this index of this epoch"). In scheme 2 every `SumNode.sum` must be 0,
  because sums are hidden, and `hash` is a Poseidon2 node. A sibling carrying a sum is `false`.
- **The courts are scheme-agnostic.** `StatementCourt` and `WithholdingCourt` only call `verifyInclusion`, read
  epochs and check signatures. Generation 2 deploys **new instances of their unchanged source**, bound to the
  new ledger, so the audited logic and its tests carry over.
- **`AnchorCouncil` hashes with `SumTree.SCHEME` (1)**, a constant, and passes only an `EpochHeader`, with no room
  for a proof or a 32-byte commitment. Generation 2 has a council whose digest is scheme 2's. The directory does
  not check councils (decision 22).
- **`NotaryPaymaster` is bound to its courts**, so generation 2 needs a new instance, with its own stake and
  deposit (13-8, the user's).

**Scheme v2.** Hashes inside the tree use gnark-crypto's Poseidon2 for BN254 (width 2, 6 full and 50 partial
rounds, S-box x^5). `Compress(l, r) = P(l, r)[1] + r`, and `MD(x₁…xₙ)` is its Merkle–Damgård hasher (state 0,
one 32-byte block per field element), exactly as `ecc/bn254/fr/poseidon2` v0.22.0 defines them. The digest
chain stays keccak.

| Item | Definition | Permutations |
| --- | --- | --- |
| leaf | `MD(1·2^128 + epoch·2^64 + amount, commitment.hi128, commitment.lo128)` | 3 |
| node | `Compress(left, right)`. No sums: they are proven inside the circuit and never published | 1 |
| empty (padding) | `MD(2·2^128)` | 1 |
| total commitment | `C = MD(5·2^128 + total, blinding)`, with a fresh 248-bit blinding per epoch, kept in the notary's database for disclosure to an auditor | 2 |
| binding root | `root = keccak256(0x04 ‖ treeRoot ‖ C)`, the value an `EpochHeader.root` carries in scheme 2 | — |
| epoch digest | v1's preimage with scheme 2 and `total = 0`: `0x03 ‖ 2:u16 ‖ number ‖ previousDigest ‖ root ‖ 0 ‖ leafCount ‖ asOf` | — |

- **Packing:**
  - The tag, epoch and amount share one field element (136 bits). The account commitment stays v1's keccak
    `bytes32`, in two 128-bit limbs.
  - Withholding's second derivation (decision 15) therefore carries over unchanged, and the deployed courts'
    source works as it is.
- **A node needs no tag, unlike v1's**, because depth is fixed per generation. v1 tags nodes because its
  canonical depth varies with the leaf count, so a node could otherwise pass for a leaf one level down. Here every
  proof carries exactly `D` siblings, so a node can never be presented as a leaf. Leaf and empty are still
  separated from each other by their tags.
- **Cost:** 256 leaves are 3·256 + 255 = 1,023 permutations, about 0.31M constraints before range checks, against
  about 0.56M with a tag in every node.
- **Depth is fixed per generation**, unlike v1's canonical depth. The circuit has a fixed capacity of `2^D`
  leaves, every proof carries exactly `D` siblings, and leaves at `index ≥ leafCount` must be the empty leaf.
  `leafCount > 2^D` is refused when building (`TooManyLeaves`). Growing past it is a new generation.

**The circuit (gnark PLONK on BN254).**
- Public inputs: epoch number, `treeRoot`, `C`, `leafCount`.
- Private inputs: each leaf's commitment limbs and amount, the total, the blinding.
- It proves:
  - every amount is a 64-bit unsigned integer;
  - every leaf at `index ≥ leafCount` is empty with amount 0;
  - the leaves hash to `treeRoot`;
  - their sum is the total, with no wrap: the total is range-checked to 64 bits, as v1 refuses a sum that does
    not fit;
  - `C` opens to that total.
- **The SRS comes from the Aztec Ignition ceremony** (176 participants, ~100.8M points), read and checked with
  Consensys's `gnark-ignition-verifier` (Apache-2.0), so there is no per-circuit trusted setup.
  - The proving and verifying keys are derived deterministically from the circuit and that SRS.
  - The verifying key's hash is pinned in the repository.
- **gnark exports the Solidity verifier.** It is committed, because it is deployed and source-verified.
  - `bin/ci` regenerates it and fails on any difference.
  - It is generated, not written here, so it is not counted as this repository's Solidity.
  - Proofs are made with `solidity.WithProverTargetSolidityVerifier`: without it they verify in Go and revert on
    chain (measured in the research).

**LiabilityLedgerV2.** It implements `ILiabilityLedger` with scheme 2's meaning, plus:
- **`proveEpoch(number, treeRoot, C, leafCount, proof)`**: anyone may call it, because the proof authenticates
  itself.
  - It runs the generated verifier and records `binding root ↔ (treeRoot, C, leafCount)` for that epoch number.
  - `anchor(header)` (anchorer role, the council) then accepts the next epoch only if its `root` and
    `leafCount` were proven for its number, and `total` is 0.
  - Two steps keep the v1 selector honest. A header can never be anchored without a proof.
- **`treeRootOf(number)` and `totalCommitmentOf(number)`**, the two values scheme 2 adds.
- **`epochAt` returns `total = 0`** for every scheme-2 epoch. A verifier reads `scheme()` first and never reads
  that zero as a total; the v1 verifier already refuses a ledger whose scheme it does not speak (`wrong_ledger`).
- **Poseidon2 in Solidity** (`Poseidon2.sol`, BN254 `mulmod`), held to gnark-crypto by a differential fixture
  (`scheme-v2.json`), the same way `scheme-v1.json` holds keccak.

**Measured in the research (gnark v0.16.4, laptop arm64, 6 cores, 7 GiB limit; to be re-measured on the x86-64
VPS):**
- Poseidon2 costs 302 PLONK constraints, and the whole tree about 630 per leaf.
- 8 leaves prove in 0.155 s with 35 MiB. 1,024 leaves take 14–16 s and 1.77 GiB.
- The PLONK verifier costs 273,855 gas on chain.
- BN254 recursion costs 8.21M constraints per inner proof, so it is **not** how a large tree is split: this
  replaces the "recursive aggregation" line this step used to carry.
- The VPS has about 1.1 GiB free (`MemTotal` 11,960 MiB, 10.6 GiB allowed to containers). One proof there fits
  about 256 leaves.

**Decisions taken 2026-10-20 on my recommendation (the user: *"pake rekomendasi kamu"*):**

| # | Question | Decision |
| --- | --- | --- |
| 43 | Where v2 lives | **Base, gnark PLONK.** Midnight was researched and rejected for now: no verifier outside Midnight, resetting testnets, two critical advisories. Compact is noted for a future service. |
| 44 | How v2 meets the immutable directory | **The v1 interfaces with scheme-2 meaning**, as above. New instances of the unchanged courts, a scheme-2 council and a new paymaster. No new directory. |
| 45 | Trusted setup | **Aztec Ignition through `gnark-ignition-verifier`**. No per-circuit ceremony, so Groth16 is out. |
| 46 | Capacity of generation 2 | **`D = 8` (256 accounts)**, settled 2026-10-20 by the VPS measurement in 13-2: 36.9 s and 1.21 GiB under a 1.5 GB, 2-CPU cap. Production holds about five accounts. More is a later generation with a dedicated prover host or commitment-chained chunk proofs, never recursion. |
| 47 | Disclosing the total | **Kept by the notary**: the blinding stays in its database, so an auditor can be shown `(total, blinding)` opening `C`. No public reveal. |
| 48 | How the watch approves | **It recomputes `treeRoot` and `C`** from the rows payment recorded plus the salts and blinding notaryd forwards over the mesh. It verifies the proof with the same verifying key before it signs, so a council approval always means two parties checked the proof. |

**13-1, 2026-10-20: built and green (`bin/ci` both halves, Go coverage 85.2%), uncommitted.**
- **Core (`internal/notary`, standard library only):** `ConfidentialScheme`, `ConfidentialTree` (fixed depth,
  `MaxConfidentialDepth = 16`, sums kept private), `ConfidentialHeader`, `ConfidentialEpoch`,
  `VerifyConfidentialInclusion` (refuses a sibling with a sum and a non-canonical sibling), `Blinding` and
  `BlindingSource`.
- **`internal/poseidon`:** its own Poseidon2 permutation on gnark-crypto's `fr` and round keys, so there is no panic
  and no unchecked error, held to gnark-crypto's `Compress` and Merkle–Damgård hasher on random inputs.
  - Known answers are pinned: `Compress(0, 0) = 0x292c…9dfb`, empty leaf `0x09c6…885a`.
  - A one-line mutation of the internal matrix turned five tests red.
- **Solidity:**
  - `confidential/Poseidon2.sol` is MIT and fully unrolled; its 62 round keys were printed once from gnark-crypto.
  - `ConfidentialTree.sol` provides `wellFormed` and `computeRoot`, split because forge-lint's `boolean-cst`
    refuses a constant inside a tuple.
  - `ConfidentialEpochChain.sol` holds the scheme-2 digest.
- **The differential fixture:**
  - `tests/fixtures` writes `scheme-v2.json`: 24 compressions including `p−1`, 20 leaves, 6 binding roots, 10
    trees of depth 0–4 with every proof, 9 rejections and a 6-link digest chain.
  - `ConfidentialDifferential.t.sol` re-checks all of it: 9/9.
  - **The decisive rejection is a sibling aliased as `x + p`.** `mulmod` would silently accept it, and only the
    explicit `< PRIME` check refuses it, on both sides.
- **Dependencies:** gnark-crypto v0.18.1 → v0.22.0 (2026-10-02), go-ethereum v1.17.8 still builds, and the `go`
  directive moved 1.26.4 → 1.26.8 (the image is 1.26.9). `bin/check-boundaries` now keeps gnark and gnark-crypto
  in `internal/poseidon`, `internal/plonk` and the composition root; a planted import in `internal/keccak` was
  caught.

**The SRS, found while preparing 13-2.** Verifying all of Aztec Ignition means downloading every one of the 176
contributions, about 1 TB, which this laptop and its link cannot do. The circuit needs about 2^20 points, all in
the first transcript of the final contribution (~322 MB). The plan:
- that transcript, checked against Aztec's published hash;
- pairing consistency over every point used, which proves a powers-of-tau structure;
- `Follows` against the previous contribution's first transcript;
- `gnark-ignition-verifier`'s KZG sanity check.

What remains trusted, and is to be written down as such: that τ is unknown, which rests on the ceremony. The full
1 TB check can be run once on a cloud machine if the user wants it.

**Revised 2026-10-20 after reading the bucket (range reads of a few hundred bytes, nothing stored):**
- **The SRS the world uses is `MAIN IGNITION/sealed/transcript00.dat`**, not the last participant's transcript.
  - Its τ in G2 equals the manifest's `crs.t2`, the point barretenberg pins.
  - The last participant's (position 197) does not equal it: a sealing step multiplied τ once more.
  - `monomial/transcript00.dat` is the same points with the generator in front.
- **Every contribution can be checked, not only the last one**, because `Follows` reads only a transcript's
  first G1 point and its two G2 points. Two range reads per transcript (92 + 256 bytes) give all 176
  participants and the seal, about 60 KB in all.
- Only the sealed transcript is read whole (322,560,348 bytes): its BLAKE2b checksum, and the powers of τ for
  the points the SRS uses.
- The ceremony's file layout is confirmed on the real files: the header reads `(0, 20, 100800000, 1, 5040000,
  2, 0)`, each coordinate is four big-endian words with the least significant first, and the G1 points start at
  `τ·[1]`, so the importer puts the generator in front.

**`internal/ignition` and `cmd/notary-srs` (built 2026-10-20).**
- **Checks, in this order:**
  - the manifest names `MAIN IGNITION` with its point counts, and positions that strictly increase;
  - every transcript header is the first transcript's;
  - every coordinate is canonical (`< p`), and no point is at infinity or outside its group;
  - each contribution follows the one before it, the first following the generator: `e(Gᵢ, [1]₂) = e(Gᵢ₋₁, [zᵢ]₂)`;
  - each contribution's τ in G2 is its τ in G1;
  - the seal's τ is the manifest's `crs.t2`;
  - the sealed file matches its own checksum, and nothing follows the checksum;
  - its first point and τ are the ones its range reads gave;
  - its points are successive powers of τ (a random linear combination, one pairing);
  - a KZG commitment made with the result opens and verifies.
- **Failure mode:** a 404 is a refusal (the manifest names a transcript nobody published). Every other HTTP
  failure, and a file that ends early, is `source_unavailable`, never a verdict on the ceremony.
- **The tool:** `notary-srs -depth D -srs out.bin [-verifier PlonkVerifier.sol]`. It sizes the SRS from the
  compiled circuit (domain + 3), writes it with `O_EXCL`, and prints its SHA-256. With `-verifier` it also derives
  the Lagrange basis, runs the setup and writes the Solidity verifier for that depth. Both hashes are what gets
  pinned.
- **Tests (`tests/ignition`)** run against a synthetic ceremony served over HTTP with ranges:
  - an honest ceremony gives exactly the powers of its τ;
  - an imported SRS proves a depth-0 epoch end to end;
  - 11 forgeries are refused: a broken chain, a τ that disagrees with itself, a seal that is not the manifest's,
    a point that is not the next power, a bad checksum, a coordinate aliased by `p`, a point at infinity,
    another ceremony, participants out of order, a missing transcript, a transcript of another size;
  - 2 outages are never read as refusals.
- **Controls:** with the chain check off, its case went red. The powers check and the KZG sanity check each catch
  a wrong power alone: with both off it went red, with either one on it stayed green.
- **The real import, 2026-10-20 04:19–04:22 UTC** (the user said yes; `notary-srs -depth 8 -srs`):
  - every check held: the chain of all 176 participants and the seal, the sealed file's checksum, the powers of τ,
    the KZG opening;
  - the SRS is 524,291 points, 16.0 MB, SHA-256 `63cdfe7a…f784d`, kept in `dist/srs/` (ignored by git);
  - `notary-srs -from` then derived the depth-8 verifier twice, byte-identical: SHA-256 `79568558…64c9`,
    verifying key `9d9d7335…4ce0`, circuit `cd63b6b8…8944`;
  - the four hashes are pinned in `internal/plonk/production.go`, and the verifier is
    `contracts/src/confidential/PlonkVerifier.sol`;
  - how the image gets the SRS is open, for 13-7. A release asset fetched with a pinned checksum is the likely
    answer.

**13-2, 2026-10-20: done, uncommitted.** The circuit, the prover, on-chain verification, the Ignition importer, the
real import, the production verifier with its pins, and the VPS measurement that settled `D = 8`.
- **The circuit (`internal/plonk.LiabilityCircuit`).**
  - Public inputs, in this order: `Epoch`, `TreeRoot`, `TotalCommitment`, `LeafCount`. The order is the
    Solidity verifier's `public_inputs`, and the fixture holds it.
  - Range checks: epoch, every amount and the total to 64 bits, every commitment limb to 128.
  - Padding: a chained `IsZero` product marks the leaves at `index ≥ leafCount`, which must be empty with amount
    0. `leafCount ≤ 2^D` is asserted.
- **Tests (`tests/plonk`).**
  - Every honest epoch of depth 0–3 at every leaf count is solved.
  - Eight forgeries are refused: another tree root, another `C`, another epoch, a −1 that cancels another
    amount, a total that overflows, a limb wider than 128 bits, an amount hidden in the padding, and more
    accounts than the capacity.
  - A control proves the forgery builder agrees with the circuit when it is honest.
  - Each added account costs about 1,417 constraints, held under 1,500 by a test.
- **The prover.**
  - `Compile`, then `NewProver(depth, system, canonical, lagrange)`.
  - `Prove` returns the proof in the verifier's byte layout. It verifies the proof itself before returning it,
    and returns `disagreement` if that fails.
  - `ExportVerifier` writes the Solidity verifier with `pragma solidity 0.8.37`.
- **Measured on the laptop (arm64; `KINETIX_NOTARY_MEASURE_DEPTHS`, peak RSS from `VmHWM`):**

  | D | accounts | constraints | prove | peak RSS |
  | --- | --- | --- | --- | --- |
  | 6 | 64 | 92,857 | 1.9 s | 314 MB |
  | 7 | 128 | 183,546 | 4.1 s | 617 MB |
  | 8 | 256 | 359,789 | 6.8 s | 1.12 GB |

  **Measured on the VPS, 2026-10-20 04:15 UTC** (the same test binary, amd64, inside the live notary image with
  `--memory 1500m --memory-swap 1500m --cpus 2 --network none --read-only`, one container per depth, then
  removed; the host had 8.3 GB available before and after):

  | D | accounts | compile | setup | prove + verify | peak RSS |
  | --- | --- | --- | --- | --- | --- |
  | 6 | 64 | 0.2 s | 2.3 s | 9.1 s | 357 MiB |
  | 7 | 128 | 0.3 s | 6.8 s | 17.5 s | 589 MiB |
  | 8 | 256 | 0.6 s | 8.7 s | 36.9 s | 1.21 GiB |

  The "1.1 GiB free" above was the sum of the containers' memory limits, not their use. notaryd itself runs at
  27 MiB of its 192 MiB.

  **Decision 46 settled: `D = 8`.** It fits under a 1.5 GB cap with 0.3 GB to spare, and costs about 47 s of
  two cores per hourly epoch.

  **For 13-5, my recommendation: a one-shot prover process.** notaryd starts `notary-prove` for each epoch, and
  the child compiles, sets up, proves, writes the proof and exits.
  - Between epochs notaryd stays at its 27 MiB. The 1.3 GB is held only for those 47 s, and is returned in full.
  - An OOM kills the child, not notaryd, which reports the failure as one.
  - notaryd's `mem_limit` must still cover the child, because a cgroup counts its children: about 1.6 GB.
  - A limit is a ceiling, not a reservation.
- **On chain, against a proof made in Go.**
  - `tests/fixtures` builds a depth-3 prover over a fixed test SRS: τ = sha256(`kinetix.notary.plonk.test-srs`)
    mod r, through `unsafekzg.WithToxicValue`. It writes:
    - `contracts/test/support/TestPlonkVerifier.sol`, gnark's verifier renamed;
    - `plonk-v2.json`: the statement, the public inputs read from gnark's public witness, and a proof.
  - The check compares the verifier byte for byte, and the statement. A PLONK proof is randomized, so its bytes
    are left to Foundry.
  - `PlonkVerifier.t.sol`, 7/7:
    - the Go proof verifies, for **275,464 gas**, with an 864-byte proof (27 words);
    - it is refused for any other statement, any input aliased by `+ p`, a missing or an extra input, each of
      its 27 words with one bit flipped, and a proof one word short or long.
  - Controls: one flipped digit in the committed proof turned Foundry red (`error ec operation`); a 64 → 63 range
    check on the total turned the Go check red (`TestPlonkVerifier.sol is stale`).
- **Found: a malformed point burns all the gas it is given.**
  - gnark's verifier calls the ecMul and ecAdd precompiles with `gas()`. A point off the curve makes the precompile
    fail, which consumes everything forwarded, and only then does it revert `error ec operation`.
  - Uncapped, the first tamper test exhausted Foundry's 1.07 billion gas; the tests now give each call 1M.
  - For 13-3: `LiabilityLedgerV2` does **not** cap the verifier's gas. A cap in an immutable contract breaks the
    day a fork reprices the precompiles, and the burn falls only on whoever sent the bad proof. The ledger treats
    a revert and `false` alike: not proven. notaryd estimates before it sends (13-5), so its own bad proof never
    leaves the host.

**The production proof, 2026-10-20.**
- `tests/fixtures` proves an epoch of 5 accounts at `D = 8` with the Ignition SRS (85 s on the laptop, setup
  included) into `plonk-production.json`, with `KINETIX_NOTARY_IGNITION_SRS` set and `-update`.
- With the variable set, the same test also derives everything again and checks the four pins and the committed
  verifier byte for byte. Without it, CI compares the statement and the pins, and Foundry verifies the committed
  proof.
- `ProductionPlonkVerifier.t.sol`, 5/5:
  - the proof verifies with the production verifier: 275,464 gas, the same as at depth 3, because PLONK
    verification does not grow with the circuit;
  - it is refused for any other statement;
  - the test verifier's code is not the production one's;
  - a depth-8 ledger proves the epoch (388,774 gas), anchors it, and closes it to a second proof.
- Controls: one byte appended to the production verifier turned its pin test red.

**13-3, 2026-10-20: `LiabilityLedgerV2` built and green, uncommitted.**
- **Proofs are keyed by statement, not by epoch number.** `proveEpoch` is open to anyone, and anyone can prove a
  well-formed tree of their own. A slot per number could therefore be overwritten, blocking the notary's anchor.
  The ledger records `keccak(number, bindingRoot, leafCount)`, and `anchor(header)` accepts only a header whose
  number, `root` and `leafCount` were proven. A stranger's proof for the same number is stored beside the
  notary's, and never in its place.
- **Other rules:**
  - `header.total` must be 0 (`TotalNotHidden`);
  - `leafCount ≤ 2^D` on both calls (`TooManyLeaves`);
  - `treeRoot` and `C` must be field elements (`NotAFieldElement`);
  - an epoch below `nextEpoch` cannot be proven any more (`EpochClosed`), and a statement cannot be proven twice
    (`StatementAlreadyProven`);
  - an epoch ahead of its turn can be proven, so notaryd may prove before the previous anchor lands;
  - a verifier that reverts and one that answers `false` are the same `ProofRejected`;
  - the verifier is called without a gas cap, as decided in 13-2;
  - the constructor refuses a verifier with no code, a depth above 16, and a role grant that does not happen.
- **The verifier is held as gnark's generated `PlonkVerifier` type, with no interface of ours.** Slither calls
  an interface method named `Verify` a naming violation, and the name is gnark's selector. `verifier()` returns
  its address.
- **Size** (`forge build --sizes`): 15,749 bytes of runtime code, under EIP-170's 24,576, and 16,884 of initcode.
- **Tests:**
  - 40 unit tests on a scripted verifier, which accepts only the exact `(proof, inputs)` it was given, so the
    order of the public inputs is checked;
  - 5 against the depth-3 test verifier and a real Go proof: Go's scheme-2 digest equals Solidity's, every
    account's inclusion verifies, any other statement is refused, a broken point is refused rather than
    panicking, and a valid proof on file anchors no other header;
  - 4 fuzz tests, 1,024 runs each;
  - 6 invariants over 16,384 calls: every anchored epoch is a proven statement with its total hidden, the chain
    recomputes, the head is the last digest, sealing is permanent, and no action had an outcome the rules
    forbid.
- **Controls:**
  - keying by number alone, plus no total check, turned 4 unit tests red;
  - no proof check in `anchor` turned 5 invariants red.
- **Gas** (depth 3, real proof): `proveEpoch` about 355,000, `anchor` 183,399, `verifyInclusion` 121,320.
- **Analysis:**
  - Slither 0;
  - Aderyn accepts `centralization-risk` on the admin's seal and the anchorer's anchor, as for v1;
  - the generated verifier is left out of both analysers, as `dependencies/` is;
  - the three `unchecked-return` findings were fixed in code.

**13-4, 2026-10-20: the generation's contracts, built and green, uncommitted. The operator scripts remain.**
- **`ConfidentialAnchorCouncil` is `AnchorCouncil` with only its digest and its ledger check changed.**
  - It hashes with `ConfidentialEpochChain`, and its constructor also requires `IConfidentialLedger`.
  - `AnchorCouncil.sol` cannot be edited, because it is the verified source of the live council.
  - `tests/sources` rebuilds the copy from the original by its four substitutions, and fails on any other
    difference. A planted `MIN_THRESHOLD = 1` turned it red.
  - The same test holds `DeployConfidentialAnchorCouncil.s.sol` to `DeployAnchorCouncil.s.sol`: only the council's
    type differs.
  - It keeps `IAnchorCouncil`, so notaryd's and the watch's bindings serve generation 2 unchanged.
  - Its 6 tests cover what differs:
    - it refuses a scheme-1 ledger;
    - its digest is the ledger's;
    - two approvals anchor a proven epoch;
    - approvals of the scheme-1 digest are refused;
    - it cannot anchor what nobody proved.
  - Everything else is `AnchorCouncil`'s suite, by construction.
- **The unchanged courts and paymaster judge generation 2** (`GenerationTwo.t.sol`, 7/7):
  - the directory registers it at the sealed head with scheme 2, and routes epochs on either side of it;
  - a new `StatementCourt` records a false statement, refuses to fault a true one, and treats a sibling carrying a
    sum as a fault;
  - a new `WithholdingCourt` takes a scheme-2 answer, refuses a wrong one and records a silence;
  - a new `NotaryPaymaster` sponsors only what its court would take, and refuses generation 1's court.
- **Deploy scripts:**
  - `DeployLiabilityLedgerV2` deploys the production verifier and the ledger in one broadcast, both through the
    CREATE2 factory. It reads the genesis (next epoch, head, last snapshot) from the predecessor on chain, and
    refuses unless the predecessor is sealed (`PredecessorNotSealed`).
  - The genesis can only be read after the seal, so the order on the real chain is: seal v1, deploy generation
    2, hand anchoring to its council, register it. Anchoring stops between the seal and the first generation-2
    anchor.
  - The courts and the paymaster use their existing scripts.
  - `DeployGenerationTwo.t.sol`, 6/6: the predicted addresses, the genesis from the sealed head, the verifier's
    code, refusing an unsealed predecessor, a taken address, a missing factory, and a scheme-2 council beside a
    scheme-1 ledger.
- **Rehearsed on a Base Sepolia fork** (`fork/GenerationTwoFork.t.sol`, block 48352128, pinned as
  `KINETIX_NOTARY_FORK_GENERATION_BLOCK`), 2/2:
  - the real v1 ledger was sealed at its real head (epoch 176 next) by the real admin, impersonated;
  - generation 2 was deployed through its scripts and the real CREATE2 factory;
  - the live, immutable `LedgerDirectory` registered it;
  - the production verifier accepted an Ignition proof of epoch 176: `plonk-production.json` is proven for the
    fork's next epoch for exactly this reason;
  - the scheme-2 council anchored it on two approvals;
  - every account's inclusion verified through the directory's generation 1.
- **The operator layer, 2026-10-20:**
  - **`deploy/base-sepolia-2.env` is generation 2, one file per generation.** It holds only what a generation owns:
    `KINETIX_NOTARY_SCHEME=2`, the predecessor, `D`, fresh salts, the same admin, keys and council members, and
    empty addresses that fill in as each part lands. The directory, the evidence registry and the fork pins stay
    in generation 1's file. Every existing `deploy-*`, `inspect-*`, `fund-*`, `hand-over-anchoring` and
    `retire-anchorer` serves it unchanged.
  - **`bin/seal-ledger <network> <head> [--broadcast]`** is new. It is the admin's, a dry run by default, and takes
    the head to seal at, which must still be the head. It prints the successor's genesis and waits for `latest` to
    show the seal.
  - **`bin/deploy-ledger`** with scheme 2 runs `DeployLiabilityLedgerV2`. It refuses while the predecessor still
    anchors, and prints the genesis it will use.
  - **`bin/deploy-council`** picks the council by the ledger's own `scheme()`, and refuses when the file and the
    ledger disagree.
  - **`bin/inspect-ledger`** with scheme 2 checks:
    - the predecessor is sealed, and the genesis is its head;
    - `D`;
    - the verifier's code is `forge inspect PlonkVerifier deployedBytecode`.
  - **`bin/inspect-council`** checks that the council hashes a sample header with its ledger's scheme. That catches
    a scheme-1 council beside a scheme-2 ledger, which a scheme-1 constructor would accept. It was run against the
    live council: ok.
  - **`bin/rehearse-deploy`** rehearses a scheme-2 file on anvil:
    - a throwaway predecessor, sealed by its admin;
    - two CREATE2s for the verifier and the ledger;
    - the courts, the scheme-2 council and the staked paymaster;
    - the council's digest scheme is checked, for both files.

    A planted scheme-1 council turned it red.
  - **Read-only runs against the live chain:**
    - `deploy-ledger base-sepolia-2` refused: the predecessor still anchors;
    - `seal-ledger base-sepolia <head>` simulated the admin's seal and sent nothing;
    - a stale head was refused.

**Decision 48 revised, 2026-10-20, on my recommendation.**
- **The watch does not verify the PLONK proof in Go.** It recomputes `treeRoot` and `C` from payment's rows plus
  the salts, secrets and blinding notaryd forwards. Before it signs, it requires
  `isProven(number, root, leafCount)` on chain.
- **Two parties still check every anchored epoch.**
  - The watch checks that the statement is payment's.
  - The chain's verifier, the canonical check, checks the proof. The ledger refuses to anchor what it did not verify.
- **What this avoids:**
  - shipping the verifying key to the watch;
  - parsing gnark's Solidity proof layout back in Go;
  - a second verifier that could disagree with the deployed one.
- **It fixes the order:**
  1. prove;
  2. `proveEpoch` lands;
  3. the watch approves;
  4. the council anchors.

**13-5 work order (off chain), each slice gated by `bin/ci`:**
- **13-5a Core:**
  - `PlanConfidentialEpoch`: the same account opening as v1, with keccak commitments, then a `ConfidentialTree` of
    depth `D` and a fresh blinding;
  - the plan carries `treeRoot`, `C`, the blinding and `D`.
- **13-5b Prover:**
  - `notary.EpochProver`, implemented by `cmd/notary-prove`, a one-shot child (input on stdin, proof on stdout);
  - the child loads the SRS file, checks it against the pin, and exits;
  - an in-process prover serves the tests.
- **13-5c Store and state machine:**
  - migration 00006: the scheme-2 columns (scheme, tree root, C, blinding, depth, proof) and an attempt kind
    (`proof` or `anchor`);
  - states: `built`, then `proved` (proof stored), then `proof_sent`, then `proven` (seen at latest), then the
    existing `submitted`, `mined` and `finalized`;
  - a proof that leaves the chain sends the epoch back to `proved`;
  - `proveEpoch` and the anchor share the anchorer key's nonce sequence, one transaction in flight at a time.
- **13-5d Watch:**
  - `watch.proto` gains the blinding, an additive field in `kinetix-contracts`, which the user publishes;
  - the watch rebuilds the scheme-2 tree, compares `treeRoot`, `C` and the root, requires `isProven`, and signs the
    scheme-2 digest.
- **13-5e Proof API and format:**
  - `inclusion.v3`: siblings without sums, the epoch's depth, no total;
  - the proof API serves it for scheme-2 epochs.
- **13-5f Verifiers:**
  - `internal/verifying` picks the scheme by the generation the directory names;
  - scheme 2 checks against `treeRootOf` and `depth()` in Go and through `verifyInclusion` on chain;
  - `notary-verify` and the browser build follow, compared case by case as today.

**13-5 progress, 2026-10-20 (uncommitted; gate green, Go coverage 84.8%):**
- **13-5a, done:**
  - `notary.PlanConfidentialEpoch`, with the same account opening and ordering as v1, now shared;
  - `ConfidentialPlan` holds the scheme, `treeRoot`, `C`, the blinding, `D` and the proof;
  - `notary.RebuildEpochTree` rebuilds either scheme from stored accounts and checks it against the anchored
    statement: root and total for scheme 1, `treeRoot`, `C` and the binding root for scheme 2. The proof API and the
    withholding answerer both use it.
- **13-5b, done:**
  - `notary.EpochProver`, implemented in process by `plonk.EpochProver`, and one shot by `cmd/notary-prove`;
  - `notary-prove` reads a JSON request on stdin, loads the SRS (`KINETIX_NOTARY_SRS_PATH`, checked against
    `KINETIX_NOTARY_SRS_SHA256`) and writes the proof on stdout. Exit codes: 65 unprovable, 78 misconfigured.
  - `internal/proverproc` runs it from `PATH`; a fixed binary name keeps gosec G204 satisfied in code;
  - the tests build the real binary over a test SRS;
  - a child that dies is `unknown`, never a verdict;
  - the image carries `notary-prove`.
- **13-5c, done, without a new state:**
  - the epoch is proven while it is built, and stored with its proof;
  - whether the proof is on chain is read from the chain (`isProven` at latest), not kept as a state;
  - `submit` sends `proveEpoch` until the statement is proven, and only then asks the watch and anchors;
  - `keepAlive` and `replace` work per attempt kind, so a proof is outbid, re-broadcast and resent exactly as an
    anchor is;
  - after proving, the writer takes the lease again, because a proof outlasts the 30 s lease: the fence still holds,
    and another writer's lease aborts the save;
  - migration 00006 adds the scheme-2 columns, CHECKs that a scheme-2 row is whole and hides its total, and an
    attempt `kind`;
  - recent costs are summed per epoch, so the runway counts both transactions;
  - `ethereum.ProofWriter` signs `proveEpoch` with the anchorer key, sharing its nonce sequence;
  - `EpochNotProven` on an anchor is `deferred`; `ProofRejected`, `EpochClosed` and the other proof refusals abandon
    the epoch;
  - `StatementAlreadyProven` counts as done;
  - metrics `kinetix_notary_proving_seconds` and `kinetix_notary_proving_failures_total`.
- **Tests:**
  - 7 state-machine cases on fakes, including a stuck proof outbid, a refused proof, a failing prover, a proof
    reorged out before any approval, a proof longer than the lease, and a lease taken while proving;
  - 4 store cases on Postgres;
  - 3 proof-writer cases on anvil, a real Go proof landing on the test verifier;
  - 2 proof-API cases.
- **Controls:**
  - without the lease renewal, its test went red;
  - without the `isProven` gate, 6 cases went red.
- **13-5d, done in the notary; the contract is the user's to publish:**
  - `kinetix-contracts` `notary/v1/watch.proto` gains `EpochProposal.confidential` (field 8), a
    `ConfidentialProposal` with the scheme, `D`, `tree_root`, `total_commitment` and the blinding. It is additive:
    `buf breaking` against `main` is clean and `buf format` passes.
  - The notary builds against the edited proto in its ignored `.contracts` copy until the user publishes, and
    GitHub CI stays red on the missing field until the pin moves, as in step 12.
  - The anchorer sends the statement and the blinding, never the proof.
  - The watch rebuilds with `AssembleConfidentialEpoch` and compares the header, `treeRoot` and `C`. It reads `D`
    and `isProven` through a keyless `ethereum.ProofReader`, and signs the scheme-2 digest.
  - A statement not yet proven is `deferred`, not refused.
  - 3 cases: approval over the scheme-2 digest, recovered; waiting for an unproven statement; 7 refusals (another
    tree root, commitment, blinding or depth, a total shown, another ledger depth, a watch without scheme 2).
  - 2 wire cases: the confidential part crosses gRPC whole, and a zero blinding never reaches the watch.
- **13-5e, decided: no `inclusion.v3`.**
  - A scheme-2 document is `kinetix.notary.inclusion.v2` with `"scheme": 2` and every sibling's sum `"0"`, because
    the document already names its scheme. The browser verifier's format and the docs' pin do not move.
  - `notary.Claim` carries its scheme. The proof API stamps it from the record and encodes each document with it.
  - `prooffile.Decode` takes the newest scheme the caller speaks, accepts 1 up to it, and limits siblings by the
    scheme's depth (64 for scheme 1, 16 for scheme 2).
  - Every caller still passes 1, except `notary-verify`, the browser build and `notary-sponsor`, which pass 2.
- **13-5f, done:**
  - `DirectoryVerifier` refuses a proof whose scheme is not its generation's (`wrong_ledger`), and verifies a
    scheme-2 generation with `ConfidentialVerifier`;
  - `ConfidentialVerifier` checks the proof with Poseidon2 against `treeRootOf` and `depth()` at `finalized`, and
    through `verifyInclusion` on chain. It refuses a verdict when the two disagree, as v1 does;
  - `ethereum.Ledger` is built for its generation's scheme, so a scheme-2 generation's identity check passes;
  - 5 cases.
- **The scheme-2 daemon, done:**
  - `KINETIX_NOTARY_SCHEME` (1 by default, or 2). With 2, `KINETIX_NOTARY_SRS_PATH` and `KINETIX_NOTARY_SRS_SHA256`
    are required.
  - The SRS file is hashed against its pin at start.
  - `D` is read from the ledger's `depth()`.
  - The ledger adapter is built for the configured scheme, so a scheme-1 notary refuses a scheme-2 ledger at start
    (`wrong_ledger`), and the other way round.
  - Blindings come from `crypto/rand`.
  - The watch always has a keyless `ProofReader` on its ledger and needs no new setting.
- **End to end on anvil and Postgres** (`TestTheDaemonProvesAnchorsAndServesASchemeTwoEpochEveryHolderCanVerify`,
  16 s), through the full dig graph:
  - the real `notary-prove` child proved epoch 0 over a test SRS;
  - `proveEpoch` landed on a depth-3 `LiabilityLedgerV2` beside the test verifier;
  - the watch, over gRPC, rebuilt the tree and required `isProven`;
  - the scheme-2 council anchored it, and it finalized with a total of 0;
  - three holders each fetched their proof from the API, a scheme-2 document;
  - the directory-routed `notary-verify` path found it `included` and `attributable`.

  Control: a proof API that ignored the claim's scheme turned it red.
- **Left at the end of 13-5:**
  - 13-6, the lab on the Ignition SRS at `D = 8`;
  - 13-7, infra and docs;
  - the user's contracts release that carries `ConfidentialProposal` — done: contracts v1.6.0, commit
    `563391d4`, pinned in `bin/sync-contracts`; the full gate passed on it (Go coverage 84.9%).

**13-6, 2026-10-20: the cutover rehearsed off chain found that it could not happen.** Every test so far had
started generation 2 on a fresh database. Production cuts over with the database it has, and three parts of
notaryd assumed one ledger for life:
- **The anchorer refused its own history.** A database whose newest final epoch is on another ledger was
  `ledger_diverged`, so notaryd pointed at generation 2 would never have anchored again. It now accepts that
  database when the configured ledger's genesis is exactly that epoch's successor and digest
  (`Anchorer.succeeds`), and refuses anything else with "does not start where it ends".
- **The proof API signed every statement for the current court.** A proof of a generation-1 epoch would have
  carried a statement for generation 2's court, which the verifier calls unattributable and generation 1's
  court cannot judge. The prover now holds a `Venue` per generation, `{ledger, court, identity}`. It signs each
  epoch for the court of the ledger that holds it, and refuses an epoch on a ledger outside the lineage.
  `CheckCourt` checks every venue's identity, digest and newest signer at start.
- **The answerer read only the current withholding court.** Holders registered before the cutover can request
  statements about old epochs for 7 days. Unanswered, each becomes a recorded silence: a public fault. The
  answerer now serves a list of venues, one per local generation. Earlier courts are scanned from the lookback
  alone, because a request older than 100,000 blocks is past its 24-hour answer period. Silences are summed
  over all of them. One answer writer serves every court through `AnswerWriter.At`, sharing one nonce state,
  because two writers on one key would fight over a nonce.
- **The court watcher** counts proven faults over every generation's court.

The lineage comes from the directory at start (`Directory.Lineage`, the generations before the configured
ledger on this chain), in `provideEarlierGenerations`.

**Proof.**
- Unit tests, each with a control that turned it red:
  - the anchorer accepts a successor and refuses a wrong genesis epoch or digest;
  - an earlier epoch is signed for its own court and recovers to its signer there, not under the current court;
  - a stray ledger is refused;
  - every earlier court is checked;
  - an old court's request is answered on that court's writer;
  - faults and silences are summed.
- **End to end** (`TestTheDaemonCarriesItsHistoryIntoAConfidentialSuccessor`, 33 s, anvil and Postgres, full dig
  graph):
  - a scheme-1 notaryd finalized epoch 0, and a holder registered;
  - the ledger was sealed, and a scheme-2 generation was deployed and registered;
  - notaryd restarted on scheme 2 with the same database and finalized epoch 1;
  - both epochs verify through the directory as `included` and `attributable`, with totals 21,000,000 and 0;
  - the old registration's request about epoch 0 on the old withholding court was answered by the new notaryd;
  - a new registration names the new court, and `/metrics` shows 6 accounts against a capacity of 8.

  Controls: a prover without the earlier venues failed the epoch-0 proof, and an answerer without them left the
  request open.

**Found while rehearsing, now in the runbook:**
- **An epoch in flight at the seal blocks generation 2.** The sealed ledger refuses it (`LedgerIsSealed` is
  `anchor_refused`), and it is abandoned. Its row keeps the number, and generation 2's first epoch has the same
  number. So the cutover stops notaryd between two epochs, after an `epoch final`, and checks that nothing is in
  flight before the seal.
- **The paymaster must exist before the first anchor,** because notaryd refuses to start (`not_sponsored`)
  without it. The earlier order — first council anchor, then the paymaster — could not run.
- **`bin/inspect-directory` compared a file with the current generation,** so generation 1's own file would
  have read `WRONG` from the cutover on. It now finds the generation that holds the file's ledger, checks its
  courts and scheme, and shows whether it is current. Run read-only against the live directory:
  `base-sepolia` passes as generation 0, current; `base-sepolia-2` is refused, because its ledger is still
  empty.
- **The admin holds 0.010 ETH**, and generation 2's paymaster needs a 1 ETH stake. The generation-1 stake must
  stay locked while the old courts can still be used. This is the user's, in "Waiting on the user".

**The lab** (`kinetix-infrastructure/scripts/lab/notary-ledger.sh`):
- it seals the lab ledger at its genesis and deploys generation 2 at fixed CREATE2 addresses: the ledger
  `0x3C51…6Dd0`, the verifier `0xb049…b3b4`, the courts `0x967b…71a6` and `0x845a…0a33`, the council
  `0xc774…E940` and the paymaster `0x1100…393C`;
- the generation-1 steps became functions, which generation 2 calls again;
- on a throwaway chain, all thirteen contracts landed at their constants, the second run changed nothing, and
  a lab whose first ledger had anchored was refused with the reset instruction.

**13-7.**
- **The image** takes the SRS by `ADD --checksum` from the release `srs-ignition-d8` into
  `/usr/local/share/kinetix-notary/ignition-d8.srs`, read-only to everyone. `docker buildx build --check`
  passes with the built-in frontend.
- **Compose:** the lab runs scheme 2 with the SRS path and pin. Production stays on scheme 1, with
  `mem_limit: 192m`, until its cutover release.
- **Two new gauges,** `kinetix_notary_accounts` and `kinetix_notary_account_capacity` (`2^D`, 0 on scheme 1).
- **Three alert rules** (critical, warning, warning):
  - the prover failing three times in an hour;
  - a proof over five minutes;
  - accounts over three quarters of capacity.
- **`docs/ALERTS.md`** lists all 22 rules; seven were missing.
- **The runbook is `kinetix-infrastructure/docs/NOTARY-GENERATION-TWO.md`.**

**The production combination, through notaryd.** This is the part of 13-6 that runs without the lab.
`TestTheDaemonAnchorsAnEpochProvenOnTheIgnitionSRSAtProductionDepth` is skipped unless
`KINETIX_NOTARY_IGNITION_SRS` is set. It ran on 2026-10-20 on the laptop with:
- the real Ignition SRS at `D = 8`;
- the production `PlonkVerifier` (a test-only binding, `plonkverifier`);
- the real `notary-prove` child, the watch and the scheme-2 council.

The epoch was proven, its proof landed and its anchor was mined within 1 min 5 s. Two holders verified it
`included` and `attributable`, with a total of 0, and `/metrics` showed a capacity of 256 with no proving failure.

**Not run:** the full compose lab, with notaryd built from the image. It waits for the release asset
`srs-ignition-d8`, the user's to upload.

**Gate, 2026-10-20:** `bin/ci` green.
- Contracts: 424 + 45 + 20 tests, the deploy rehearsal, Slither and Aderyn.
- Go: lint 0 issues, coverage 85.0%.
- The first run was red on three lines over 120 columns and five `wsl_v5` blank lines; both were fixed in code.

**The Verify pages for scheme 2 wait for the release.** `kinetix-docs` links code by permalink at the commit
production runs, and none runs scheme 2 yet.

**Work order.**
- **13-1 Scheme v2:** Poseidon2 in Go, through gnark-crypto behind `notary.Scheme`, and in Solidity, with
  known-answer vectors and the differential fixture.
- **13-2 Circuit:**
  - the circuit itself;
  - tests: honest, overflow, padding, wrong sum, wrong commitment, wrong epoch;
  - constraint counts, the Ignition SRS import and check, the pinned keys;
  - the generated verifier with its CI check;
  - the x86-64 VPS measurement that settles `D`.
- **13-3 Ledger:** `LiabilityLedgerV2` with unit, fuzz, invariant and differential tests, gas, Slither and Aderyn.
- **13-4 Generation suite:**
  - the scheme-2 council;
  - court and withholding instances re-tested against the new ledger;
  - the paymaster instance;
  - the deploy scripts with `bin/rehearse-deploy`, and a seal → register rehearsal on a fork.
- **13-5 Off chain:**
  - notaryd's prove → `proveEpoch` → council anchor, as one state machine with its own metrics;
  - the watch's checks;
  - proof format `inclusion.v3`;
  - `notary-verify` and the browser build for scheme 2;
  - migrations: blinding, tree root, C, proof, proof transaction.
- **13-6 Lab:** end to end on anvil: seal v1, deploy generation 2, register, anchor, a holder verifies, both
  courts judge.
- **13-7 Infrastructure and docs:**
  - compose and memory limits;
  - Grafana rules for proving time and failures;
  - the Verify pages for scheme 2.
- **13-8 The user**, once 13-5 to 13-7 have released a notaryd that anchors scheme 2. The sequence, with its checks,
  is `kinetix-infrastructure/docs/NOTARY-GENERATION-TWO.md`. In short, after the SRS release asset and the
  funds exist:
  0. stop notaryd after an `epoch final`, with nothing in flight;
  1. `bin/seal-ledger base-sepolia <head> --broadcast`;
  2. `bin/deploy-ledger base-sepolia-2 --broadcast`, its address into the file, `bin/inspect-ledger base-sepolia-2`;
  3. the court, the withholding court (and its `FROM_BLOCK`) and the council on `base-sepolia-2`, each inspected;
  4. `bin/register-successor`, then `bin/inspect-directory` on both files;
  5. `bin/hand-over-anchoring base-sepolia-2 --broadcast`;
  6. the paymaster deployed, staked, funded and inspected on `base-sepolia-2`;
  7. the infra release: generation 2's contracts, no scheme override, `mem_limit: 1600m`;
  8. the first council anchor, then `bin/retire-anchorer base-sepolia-2 <tx> --broadcast` and
     `KINETIX_NOTARY_ANCHORER_ROLE_HOLDERS`.

### Step 14 — the next service: product passport

After the notary is complete: per-unit provenance for serialized goods (ERC-721), minted by a
verified merchant and transferred to the buyer on delivery. It needs serial-level inventory, which
warehouse does not have today (no serial column anywhere in its schema), and buyer wallets
(account abstraction — the step 10 paymaster is the first piece of that), so it starts with that
owner-side work.

## 11. Open questions

- ~~Publishing the verifier image to GHCR, and under which name.~~ Settled by decision 33.
