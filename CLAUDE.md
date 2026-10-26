# CLAUDE.md — kinetix-notary-service

Guidance for Claude Code in this repository. The notary is part of the **Kinetix** polyglot estate.
It is **Solidity + Foundry** on chain and **Go** off chain, in a strict hexagonal layout with
**`go.uber.org/dig`** as the only dependency-injection mechanism.

**`PLAN.md` is the agreed design** — the decisions the user took on 2026-10-11, the exact scheme, the
contract's rules, the key model and the sequence of work through zero-knowledge proofs. Read it
before changing anything structural. This file is the short version: the rules, and what enforces
them.

## What this service is

It commits to what the platform owes, anchors the commitment on a public chain, and lets every
account holder prove their balance is inside it without trusting Kinetix. From step 12 it does the
same for evidence: it reads the append-only logs of order (lifecycle), matching (deliveries) and review,
commits to every record with a salt of its own, anchors a batch per stream a day in the
`EvidenceRegistry`, and gives each record's proof to the parties to it. It never owns a balance,
an order or a review — it reads what the owner wrote and commits to it — and it never moves
money on chain.

## Everything runs in Docker

Nothing is installed on the laptop. `bin/dev` runs any command inside the pinned toolchain image,
with a local chain (anvil) started beside it and `KINETIX_NOTARY_TEST_RPC_URL` pointing at it.

```sh
bin/dev bin/ci                                  # the whole gate: bin/ci-contracts, then bin/ci-go
bin/dev bin/ci-contracts                        # one half: forge fmt/build/test/snapshot, Slither, Aderyn
bin/dev bin/ci-go                               # the other: generation, gofmt, vet, boundaries, lint, bin/test
bin/dev sh -c 'cd contracts && forge test'      # contracts only
bin/dev bin/test                                # Go suite with true coverage
bin/dev go test ./tests/fixtures -update        # regenerate the differential fixtures
bin/dev down                                    # stop the chain
```

The image built from `Dockerfile` carries six binaries and runs `notaryd` by default (the sixth, `notary-prove`, is
the one-shot scheme-2 prover notaryd starts per epoch, PLAN.md step 13-5):
`notary-migrate` is the one-shot migrator, `notary-watch` is the council's second member (PLAN.md step 8),
an account holder verifies with `docker run <image> notary-verify proof.json` (the image used to
start `notary-verify` itself), and sends a sponsored call with `notary-sponsor` (PLAN.md step 10).
`cli/Dockerfile` builds the auditors' image, `ghcr.io/wildanfrananda/kinetix-notary-cli` (decision 33): `notary-verify`
and `notary-sponsor` only, on distroless static, for amd64 and arm64 by cross-compiling in a build stage pinned to
`$BUILDPLATFORM`, never by emulation. `cmd/notary-verify-wasm` is the verifier for the browser (decision 31,
`GOOS=js GOARCH=wasm`): it sets `kinetixNotaryVerify(document, settings)`, a promise of `{exit, stdout, stderr}`
exactly as `notary-verify` would print and exit, where `settings` holds the four `KINETIX_NOTARY_*` variables the CLI
reads. `bin/build-verifier-wasm <version> <dir>` builds it reproducibly (`-trimpath`, `-buildvcs=false`, the version
in `main.version`) beside the `wasm_exec.js` of the same Go and a `SHA256SUMS`. The release workflow attaches all three
to the GitHub release, and `kinetix-docs` pins that release and its hashes (decision 35), so any change to the Go code
the verifier links changes the file the docs must pin.
`bin/generate` runs every generator — abigen, protoc through `bin/sync-contracts` (contracts pinned to a
tag and its commit), sqlc — and nothing it writes is ever committed.

Operations on a real network have their own scripts, all through the same image:

```sh
bin/new-keystore notary-deployer|notary-admin   # encrypted keystore outside every repository
bin/new-anchorer-key prod                       # hot key straight into the SOPS store
bin/deploy-ledger base-sepolia [--broadcast]    # dry run by default; reads deploy/<network>.env
bin/inspect-ledger base-sepolia                 # every on-chain fact checked against the config
bin/new-statement-key prod                      # the statement signer, straight into the SOPS store
bin/deploy-court base-sepolia [--broadcast]     # the statement court beside the ledger; dry run by default
bin/inspect-court base-sepolia                  # the court's facts, its signer terms and its faults
bin/new-answerer-key prod                       # the answerer, straight into the SOPS store
bin/fund-anchorer|fund-answerer base-sepolia 0.001   # the deployer pays a hot key's gas
bin/deploy-withholding base-sepolia [--broadcast]    # the withholding court beside the court; dry run by default
bin/inspect-withholding base-sepolia            # its facts, the answerer's balance and the recorded silences
bin/new-watch-key prod                          # the watch's council key, straight into the SOPS store; it never holds ETH
bin/deploy-council base-sepolia [--broadcast]   # the anchor council beside the ledger; dry run by default
bin/inspect-council base-sepolia                # its facts, its members and who holds the anchorer role
bin/fund-admin base-sepolia 0.00002             # the admin's gas for the two role changes below
bin/hand-over-anchoring base-sepolia [--broadcast]           # the admin grants the council the anchorer role
bin/retire-anchorer base-sepolia <council anchor tx> [--broadcast]   # then revokes the anchorer key's own
bin/deploy-directory base-sepolia [--broadcast] # the ledger directory, generation 0 = the ledger and its courts
bin/inspect-directory base-sepolia              # its facts, and that its current generation is deploy/<network>.env's
bin/seal-ledger base-sepolia <head> [--broadcast]   # admin; permanent; the successor's genesis is read at it
bin/register-successor base-sepolia <ledger> <court> <withholding> [--broadcast]   # admin; after a seal
bin/deploy-ledger base-sepolia-2 [--broadcast]  # generation 2: deploy/<network>-2.env, KINETIX_NOTARY_SCHEME=2, after the seal
bin/deploy-paymaster base-sepolia [--broadcast] # the paymaster; refuses unless the EntryPoint and account carry the pinned code hashes
bin/stake-paymaster base-sepolia                # the admin's addStake up to the file's stake; refuses another key or an unlocking stake
bin/fund-paymaster base-sepolia 0.1             # the deployer's EntryPoint.depositTo(paymaster); waits for latest to show it
bin/inspect-paymaster base-sepolia              # its facts, the courts' current signer, the stake, the deposit in operations at the ceiling
bin/dev bin/build-verifier-wasm v1.4.0 dist/verifier   # the browser verifier, reproducible, with wasm_exec.js and SHA256SUMS
bin/new-evidence-anchorer-key prod              # the seventh key, the evidence anchorer, straight into the SOPS store
bin/fund-evidence-anchorer base-sepolia 0.001   # the deployer pays the evidence anchorer's gas
bin/deploy-evidence-registry base-sepolia [--broadcast]   # the evidence registry; dry run by default
bin/inspect-evidence-registry base-sepolia      # its facts, its streams, who holds the anchorer role, Sourcify
bin/dev go run ./cmd/notary-srs -depth 8 -srs /workspace/dist/srs/ignition-d8.srs   # Aztec Ignition → checked SRS (PLAN.md step 13)
bin/dev go run ./cmd/notary-srs -depth 8 -from /workspace/dist/srs/ignition-d8.srs -verifier /workspace/dist/srs/PlonkVerifier.sol
bin/dev env KINETIX_NOTARY_IGNITION_SRS=/workspace/dist/srs/ignition-d8.srs go test ./tests/fixtures -run Production -update
bin/dev env KINETIX_NOTARY_IGNITION_SRS=/workspace/dist/srs/ignition-d8.srs go test ./tests/integration -run Ignition   # notaryd at D = 8, ~80 s
```

`notary-verify` reads `KINETIX_NOTARY_RPC_URL`, `KINETIX_NOTARY_CHAIN_ID`, `KINETIX_NOTARY_DIRECTORY_ADDRESS`
and `KINETIX_NOTARY_RPC_TIMEOUT`, every one required (decision 24, step 9). It asks the `LedgerDirectory`
which generation holds the claim's epoch and verifies against that generation's ledger and court; it
trusts that court, never the one a document names, and reports beside every verdict whether the
proof's statement is attributable to the epoch's signer and how many faults the epoch carries. An epoch
held on another chain is `moved` (exit 9), never `not_anchored`; an address that is not a directory is
`wrong_directory` (exit 8). The directory is read at `latest`, the epoch itself at `finalized`, so a
successor that is not final yet can never make an epoch `included`.

notaryd also needs `KINETIX_NOTARY_WITHHOLDING_COURT_ADDRESS`, `KINETIX_NOTARY_WITHHOLDING_FROM_BLOCK` (the block
the withholding court was deployed in), `KINETIX_NOTARY_ANSWERER_PRIVATE_KEY` and `KINETIX_NOTARY_RPC_LOG_RANGE`
(the most blocks the RPC endpoint answers in one `eth_getLogs`).

notaryd and notary-watch both need `KINETIX_NOTARY_RPC_RATE`, the most requests a second the process sends the RPC
endpoint (1 to 1000). Both refuse a `ws` or `wss` endpoint, because only HTTP requests can be paced. `ethereum.Dial`
sends every request of the process through one `rate.Limiter` with a burst of one, so no second holds more than
rate + 1 of them. A 429 that still comes is waited out for a second and sent again, up to three times, unless the
caller's deadline is nearer. The verifiers and `notary-sponsor` get the retry without a rate. Each process exports
`kinetix_notary_rpc_requests_total`, `kinetix_notary_rpc_throttled_total` and `kinetix_notary_rpc_paced_seconds_total`.

It refuses to start unless:
- the withholding court stands beside its court;
- it reads the notary's ledger;
- it hashes a registration the way the notary does;
- the endpoint answers one full log span ending at the head.

From step 8 notaryd anchors through the `AnchorCouncil` and also needs `KINETIX_NOTARY_COUNCIL_ADDRESS` and
`KINETIX_NOTARY_WATCH_GRPC_ADDR`. From step 9 it also needs `KINETIX_NOTARY_DIRECTORY_ADDRESS`, and refuses to
start (`wrong_directory`) unless the directory's current generation is its own ledger, court and withholding
court on its chain: the proofs it serves are verified through that directory. It refuses to start unless the council anchors its ledger, holds the
anchorer role there, counts the anchorer key as a member, needs no more than the two approvals notaryd
gathers, and hashes an approval as Go does.

From step 10 notaryd also needs `KINETIX_NOTARY_PAYMASTER_ADDRESS`, and refuses to start (`not_sponsored`)
unless that paymaster sponsors its own court and withholding court on its chain and accepts statements from
the court's current signer. It then reads the deposit every tick into `kinetix_notary_paymaster_deposit_wei`,
the locked stake into `kinetix_notary_paymaster_stake_wei` (0 from the moment the admin unlocks it), and sets
`kinetix_notary_paymaster_current` to 0 the moment the signer moves on without a new paymaster. A missing
stake does not stop notaryd: its own work needs no paymaster, so the stake is an alert, not a refusal.

From step 13 notaryd takes `KINETIX_NOTARY_SCHEME` (1 by default). With 2 it also needs `KINETIX_NOTARY_SRS_PATH`
and `KINETIX_NOTARY_SRS_SHA256`, hashes the SRS against that pin at start, and reads `D` from the ledger. It proves each
epoch by starting `notary-prove`, a one-shot child, sends `proveEpoch` before it asks the watch, and anchors only a
statement the ledger already holds as proven. The ledger adapter is built for the configured scheme, so a notary of one
scheme refuses a ledger of the other at start (`wrong_ledger`).

A successor keeps its predecessor's database (PLAN.md 13-6). At start notaryd reads every earlier generation on its
chain from the directory (`Directory.Lineage`, `provideEarlierGenerations`), and then:
- it signs each epoch's statement for the court of the generation that holds the epoch (`proving.Venue`);
- it answers every earlier withholding court, scanned from 100 000 blocks back, through one answer writer that
  shares a nonce state across courts (`AnswerWriter.At`);
- it sums proven faults and recorded silences over all of them.

The anchorer continues a database whose newest final epoch is on another ledger only when the configured ledger's
genesis is that epoch's successor and digest. Anything else is `ledger_diverged`, "does not start where it ends".
The image takes the SRS by `ADD --checksum` from this repository's release `srs-ignition-d8`. No image builds
until that release exists.

From step 12 notaryd also needs `KINETIX_NOTARY_ORDER_GRPC_ADDR`, `KINETIX_NOTARY_MATCHING_GRPC_ADDR`,
`KINETIX_NOTARY_REVIEW_GRPC_ADDR`, `KINETIX_NOTARY_EVIDENCE_REGISTRY_ADDRESS` and
`KINETIX_NOTARY_EVIDENCE_ANCHORER_PRIVATE_KEY`; `KINETIX_NOTARY_EVIDENCE_BATCH_INTERVAL` (default 24h,
decision 42) and `KINETIX_NOTARY_FEED_TIMEOUT` (default 30s) are optional. It refuses to start unless the
registry is one on its chain, commits with scheme 1, keeps streams 1–3 and grants the evidence key its
anchorer role, and unless that key is not the liabilities anchorer (decision 41). The evidence loop shares
the writer lease with anchoring, and keeps one batch in flight across all three streams, because one key
has one nonce sequence. A batch the registry refuses for good is abandoned and stops its own stream; the
other two go on. `/health/ready` also requires the evidence loop to have completed a pass.

From step 6b notaryd also needs `KINETIX_NOTARY_NOTIFICATION_GRPC_ADDR`, `KINETIX_NOTARY_STATEMENT_VERIFY_URL` (an
https URL, the public Verify page) and `KINETIX_NOTARY_STATEMENT_UTC_OFFSET` (`+07:00`); `KINETIX_NOTARY_NOTIFY_TIMEOUT`
(default 30s) is optional. Once a day, at the first final epoch whose snapshot falls on a new day at that offset, the
statement loop (`internal/delivering`) plans one statement for each holder whose amounts changed since the last
statement planned for them, and hands each to notification's `Notify` as `TEMPLATE_LIABILITY_STATEMENT` with only
`epoch` and `verify_url`: the user ruled on 2026-10-20 that no amount goes through the email provider. Suspense
accounts are left out, the first run starts at the newest final epoch, and a holder notification has no address for is
`refused` and asked again only once their balance changes. Anything else notification says is retried with a backoff
from 1 minute to 6 hours, never taken as delivered. `/health/ready` also requires the statement loop to have
completed a pass. The release order and the refusals are in infra `docs/NOTARY-STATEMENTS.md`.

A document whose `format` is `kinetix.notary.evidence.v1` is evidence: `notary-verify` and the browser
verifier then read `KINETIX_NOTARY_EVIDENCE_REGISTRY_ADDRESS` instead of the directory, with the same RPC,
chain id and timeout, and trust that registry, never the one the document names. A missing setting is
`ExitConfig` (78), the same as for a liabilities proof.

`notary-sponsor` reads the verifier's four variables, `KINETIX_NOTARY_PAYMASTER_ADDRESS`,
`KINETIX_NOTARY_BUNDLER_URL`, `KINETIX_NOTARY_HOLDER_KEYSTORE`, `KINETIX_NOTARY_HOLDER_PASSWORD_FILE` and,
optionally, `KINETIX_NOTARY_RECEIPT_WAIT` (default 3m, at most 30m). It signs as the holder, delegates the
holder's EOA to the pinned `Simple7702Account` on its first operation, and never sends a call the paymaster
would refuse: a registration for another withholding court, an epoch held by another generation, or a proof
whose statement is not attributable to the court's signer is `not_sponsored` (exit 2) before anything is
signed. Before it signs it also asks the paymaster's `sponsorship` view, the same checks the paymaster runs
in validation: a call the court would refuse is `call_refused`, one that is not due yet (an unanswered record
before the answer is overdue) is `not_due`, both exit 3, and neither is sent. A landed operation whose call
reverted is `reverted` (exit 1), and the paymaster paid for it; since decision 26 was revised that takes a
race, such as someone proving the same fault between the bundler's check and the bundle.

`notary-watch` reads the verifier's five variables, `KINETIX_NOTARY_COUNCIL_ADDRESS`,
`KINETIX_NOTARY_WATCH_PRIVATE_KEY`, `KINETIX_NOTARY_PAYMENT_GRPC_ADDR`, `KINETIX_PKI_DIR`,
`KINETIX_NOTARY_WATCH_GRPC_PORT` and `KINETIX_NOTARY_HTTP_PORT` (`KINETIX_TRUST_DOMAIN` optional). It keeps
no database, admits only the `notary` identity on its gRPC port, and refuses to start unless its key is a
member of the council and the council hashes an approval as Go does.

`deploy/<network>.env` holds only public values and is committed. `KINETIX_NOTARY_ANCHORER_ROLE_HOLDERS` there
names who should hold the ledger's anchorer role, and `bin/inspect-ledger` checks the anchorer key and the
council against it; change it in the same commit as the role changes. Keys never are: the deployer and
admin keystores live in `KINETIX_NOTARY_KEYSTORE_DIR` (refused when inside a git work tree), the
anchorer lives in SOPS. Why the split is in PLAN.md §8.

Every tool version lives in **one** place, `bin/install-toolchain`, installed by component
(`foundry`, `abigen`, `lint`, `analysis`): Foundry and Aderyn by release tarball and SHA-256, abigen
and golangci-lint by `go install @version`, Slither from a hash-locked requirements file. The
production image installs only `foundry abigen`. The toolchain image, the GitHub jobs and the GitLab
jobs all run that script, so local and CI cannot drift. Bump a version there and nowhere else.

GitHub CI runs the two halves of the gate as parallel jobs, `gate (contracts)` and `gate (go)`, each
with `sh bin/dev bin/ci-<half>` inside the toolchain image, so the laptop and GitHub run the same
Python, Go and tools. The image is built from `toolchain/Dockerfile` by `docker/build-push-action`
with its layers in the Actions cache (scope `toolchain`, written only on pushes to `main`), and
`KINETIX_NOTARY_TOOLCHAIN_PREBUILT=1` tells `bin/dev` not to build it again. A job named `verify`
passes only when both halves passed: it is the check the `main` ruleset requires, and it runs with
`if: always()` because GitHub counts a skipped required check as passed. GitLab runs the same halves
as `verify-contracts` and `verify-go`.

## Layout, and the rule it encodes

| Package | May import |
|---|---|
| `internal/notary` | the standard library. **Nothing else.** |
| `internal/verifying` | `internal/notary` only; `DirectoryVerifier` resolves a generation, then verifies through `Verifier` |
| `internal/keccak`, `prooffile`, `config`, `verifycli` | `internal/notary` + their own library |
| `internal/poseidon` | `internal/notary` + gnark-crypto: scheme 2's Poseidon2 tree, its own permutation held to gnark-crypto's (PLAN.md step 13) |
| `internal/plonk` | `internal/notary` + gnark: scheme 2's circuit, its prover and the Solidity verifier it exports |
| `internal/ignition` | `internal/notary` + gnark-crypto + `net/http`: reads Aztec Ignition, checks every contribution and builds the KZG SRS (`cmd/notary-srs`) |
| `internal/ethereum` | the above + go-ethereum + `internal/contractgen` — the only adapter that touches a chain |
| `internal/payment`, `internal/feeds`, `internal/notification`, `internal/mesh` | `internal/notary` + gRPC + `internal/protogen` — the only adapters that speak gRPC; `feeds` reads the three evidence logs, `notification` hands holders' statements over |
| `internal/postgres` | `internal/notary` + pgx + goose + `internal/storegen` — the only adapter that touches Postgres |
| `internal/anchoring`, `internal/proving`, `internal/registering`, `internal/answering`, `internal/approving`, `internal/sponsoring`, `internal/evidencing`, `internal/delivering` | `internal/notary` only, like `internal/verifying`; `evidencing` collects, batches, anchors and proves evidence; `delivering` plans and hands over holders' statements |
| `internal/bundler` | `internal/notary` + the standard library: the ERC-4337 bundler's JSON-RPC, nothing from go-ethereum |
| `internal/sponsorcli` | `internal/notary` + `prooffile`; turns a kind into a verdict and an exit code, like `verifycli` |
| `internal/epochwatch` | `internal/notary` + gRPC + `internal/protogen` + `internal/mesh` — the watch's server and notaryd's client |
| `internal/wiring` | everything; the composition root, and the only package outside `cmd/` allowed to import dig |
| `cmd/…` | everything |

Package names say what a package provides (`keccak.NewScheme`, `ethereum.NewLedger`,
`verifying.NewVerifier`); there is no `ports/` or `adapters/`. Ports are the interfaces in
`internal/notary` (`Scheme`, `Ledger`).

## How the rules are enforced

- **`bin/check-boundaries`** fails the build on any of: a non-stdlib dependency reachable from
  `internal/notary`; `internal/verifying` reaching anything but `internal/notary`; dig imported
  outside `cmd/` and `internal/wiring`; go-ethereum imported outside `internal/ethereum`,
  `internal/contractgen` and `cmd/`. Controlled both ways on 2026-10-11: four planted violations, four
  reports, exit 1; removed, clean. That control ran on the laptop; on GitHub the import rules checked
  nothing until 2026-10-16 (Traps: "On GitHub the checkout is not the toolchain's").
- **`bin/ci`** is the single gate, `bin/ci-contracts` and then `bin/ci-go`: `forge fmt --check`,
  `forge build` (forge-lint with `deny = "warnings"`), `forge test` on the CI profile,
  `forge snapshot --check`, the fork suite, `bin/rehearse-deploy` (both deploy scripts on a throwaway
  anvil with each network's chain id; nothing reaches a network), Slither and Aderyn; then binding
  generation, `gofmt`, `go mod verify`, `go vet`, the boundary check, golangci-lint, and `bin/test`.
  `bin/ci` and `bin/ci-go` refuse to run without a node, because a skipped integration test proves
  nothing. `bin/ci-go` installs the Soldeer dependencies itself, so it passes from a clean checkout
  without the other half.
- **`bin/analyze`** runs Slither and Aderyn over `contracts/src` and fails on anything not already
  accepted. Accepted, and why: Slither `timestamp` (refusing a future snapshot needs
  `block.timestamp`; a validator moves it by seconds), Aderyn `centralization-risk` (the admin and
  anchorer roles are the design; the admin key never goes online), Aderyn `unchecked-return` on the
  constructor's `_grantRole` (it cannot fail there, and the deployed ledger's source is frozen: an edit
  would make the repository disagree with the verified bytecode); for `StatementCourt`, Aderyn
  `centralization-risk` (the admin schedules signers, by design) and `reentrancy-state-change` (its only
  external calls are `view` reads of the immutable ledger, compiled to STATICCALL, which cannot re-enter;
  Slither, which tells a STATICCALL apart, reports nothing there); for `AnchorCouncil` and its scheme-2 copy
  `ConfidentialAnchorCouncil`, Aderyn `centralization-risk` (the admin schedules memberships, behind a 7-day delay, by
  design); for `LiabilityLedgerV2`, the same entries as v1's minus `unchecked-return`, which v2 fixes in code; for
  `LedgerDirectory`, Aderyn `centralization-risk` (only the admin registers a successor, decision 23); for
  `NotaryPaymaster`, Aderyn `centralization-risk` (only the admin withdraws the deposit, the kill switch of
  decision 26, and only the admin adds to, unlocks or withdraws the stake, because the EntryPoint relocks a
  stake whenever it grows and anyone else could cancel an unlocking). An accepted
  entry that stops being reported also fails, so the list cannot rot. forge-lint's `reentrancy-events` is off in `foundry.toml`
  for the same reason and with the same backstop: it flags any event after any external call.
- **Two lint settings are decisions, not silencings.** `wrapcheck` ignores `.RecvMsg(`: the gRPC
  client interceptor that counts calls must hand the stream's own error back unchanged, because callers
  compare `io.EOF` by identity and a wrapped end of stream reads as a failure. It also ignores
  `status.Error(`: the watch's handler must return a gRPC status as it is, because a wrapped one reaches
  the caller as `Unknown` and notaryd could no longer tell a refusal from an outage.
- **The differential fixture is a contract between the languages.** `tests/fixtures` regenerates
  `contracts/test/fixtures/scheme-v1.json` and `evidence-v1.json` and fails when a committed file differs;
  Foundry's `Differential.t.sol` and `EvidenceDifferential.t.sol` re-verify every tree, proof, rejection
  and digest in them. Proven on
  2026-10-11 by mutating Go's leaf tag: Foundry went red. Scheme 2 adds `scheme-v2.json`
  (`ConfidentialDifferential.t.sol`) and `plonk-v2.json` (`PlonkVerifier.t.sol`).
- **A generated verifier is gnark's output, byte for byte.** `tests/fixtures` exports the PLONK verifier of a
  depth-3 circuit over a fixed test SRS (`unsafekzg.WithToxicValue`, never cached to disk) into
  `contracts/test/support/TestPlonkVerifier.sol`, renamed from `PlonkVerifier`, and fails when the committed file
  differs. Its comments, Apache-2.0 header and layout are gnark's: it is in forge fmt's `ignore` list, and it is
  never edited, formatted or stripped by hand, or the check goes red. forge-lint and solc report nothing on it, so
  it has no other exemption. A PLONK proof is randomized, so `-update` proves again every time and changes the
  `proof` line of `plonk-v2.json`; the Go check compares the statement and the public inputs, and Foundry checks
  that the committed proof verifies. Proven on 2026-10-20: one flipped digit in the proof turned Foundry red, a
  64→63 range check in the circuit turned the Go check red.
- **The production verifier is gnark's output too, from the Aztec Ignition SRS.** `src/confidential/PlonkVerifier.sol`
  is the depth-8 verifier `notary-srs -from dist/srs/ignition-d8.srs -verifier` writes, byte for byte. It is in forge
  fmt's and forge-lint's `ignore`. It is left out of Slither (`filter_paths`) and Aderyn (`-x`) the way `dependencies/`
  is, because it is third-party generated code that may never be edited. forge-lint reports `uninitialized-local` and
  `require-revert-in-loop` in its assembly, though it reports nothing on the test copy in `test/support`. `internal/plonk/production.go` pins four SHA-256s: the
  circuit at `D = 8`, the SRS, the verifying key and the verifier file.
  - `tests/plonk` fails when the committed file or the compiled circuit stops matching its pin.
  - With `KINETIX_NOTARY_IGNITION_SRS` set, `tests/fixtures` derives everything from the SRS again and checks all
    four. With `-update` as well, it proves `plonk-production.json` afresh, which Foundry verifies with the
    production verifier.
  - A change to the circuit, or a gnark bump, therefore goes red until the verifier is derived again and every pin
    and fixture is updated together. The deployed generation-2 ledger keeps the verifier it was deployed with.
- **A copy of a deployed contract is held to its original by `tests/sources`.** `ConfidentialAnchorCouncil.sol` and
  `DeployConfidentialAnchorCouncil.s.sol` must equal `AnchorCouncil.sol` and `DeployAnchorCouncil.s.sol` after the
  listed substitutions. Change the original and every substitution's count is checked again. Change the copy
  anywhere else and the test goes red. The original stays frozen, because it is the verified source of the live
  council.
- **The ledger holds the generated `PlonkVerifier` type, never an interface of ours.** Slither reports an interface
  method named `Verify` as a naming violation, and the name is fixed by gnark's selector.
- **Every adapter carries a compile-time assertion** (`var _ notary.Ledger = (*Ledger)(nil)`).
- **Generated code is committed, and CI proves it is current.** `internal/contractgen/` (abigen, from the forge
  artifacts), `internal/protogen/` (protoc, from the pinned kinetix-contracts tag) and `internal/storegen/` (sqlc)
  are in git since 2026-10-26, at the user's choice: Dependabot runs `go mod tidy` on a bare clone and cannot run
  `bin/generate`, so with them ignored every Go update failed on `internal/storegen: no matching versions for query
  "latest"`. `bin/ci-go` regenerates them and fails when the tree differs from what is committed, so they cannot
  drift from their sources. The images still generate their own (`.dockerignore` keeps the three out), and
  `.contracts/` (the vendored `.proto`) stays ignored: no `.proto` is committed here.

## Conventions

- **Zero comments**, in Solidity and Go alike. The single exception is the SPDX line in each `.sol`
  file: the compiler warns without it and `deny = "warnings"` would then fail the build. Lint
  findings are fixed in code, never silenced with an inline directive — there is no `nolint`, no
  `forge-lint: disable`.
- **Directives are not comments.** The goose annotations in `internal/migrations/*.sql` and
  `//go:embed` are instructions to a tool, and stay.
- **Explicit `return`** in every Solidity function, named returns included, as in the Rust service.
- **One type per file** in both languages: each Solidity struct, interface, library and contract has
  its own file; each Go type in `internal/notary` has its own file.
- **Go layout is enforced, not remembered** (the user's style, restated 2026-10-12 for every repository):
  a call, composite literal or signature that does not fit on one line puts **every element on its own
  line** and the closing bracket alone — slog's key/value pairs included, a key and its value on
  separate lines. The one exception is the trailing-block idiom, `t.Run(name, func(...) {…})`, whose
  head fits on its first line. And **a blank line between contexts**: after a closing brace, before a
  return that ends a longer block, between groups of statements; an assignment stays cuddled with the
  `if err != nil` that checks it. Three tools hold this in `bin/ci`: `golines` (120 columns),
  `tools/arglayout` (a wrapped call that still shares lines; `go run ./tools/arglayout -fix …` rewrites
  it) and golangci-lint's `wsl_v5` (`golangci-lint run --fix` adds the blank lines).
- **Braces on the same line**; long calls one argument per line. For Solidity this is
  `forge fmt` with `multiline_func_header = "params_first_multi"`: a header that does not fit breaks
  its parameters one per line and keeps `) … {` together. The default (`attributes_first`) puts `{`
  alone on its own line, and that mismatch turned CI red on 2026-10-11 the moment the code was
  formatted by hand to the house style. Fix the formatter's configuration, never the house style.
  `prefer_compact = "none"` (2026-10-12) makes a wrapped call put every argument on its own line,
  the same rule the Go tools enforce. `src/**` is in forge fmt's `ignore`: those files are the
  verified source of the deployed ledger and stay byte-identical, so they keep the older compact
  wrapping until a successor ledger is written.
- **Go tests live in `tests/`**, never beside the code; use `bin/test`, because plain
  `go test -cover` reports 0% for every package under `internal/`. Foundry tests live in
  `contracts/test/`, Foundry's own convention, except the two suites that need their own profile:
  `contracts/fork/` (a network) and `contracts/rules/` (a tracer).
- **Hand-written fakes** (`tests/notarytest`), no mocking framework.
- **Errors carry a `notary.Kind`**, read back with `notary.KindOf`. Only the edge (`verifycli`) turns a
  kind into a verdict and an exit code.
- **dig composes through parameter and result objects, in `internal/wiring` only** (the user's
  choice, 2026-10-12). Every provider takes one `dig.In` struct naming what it needs and returns either
  a plain value or a `dig.Out` struct; a `dig.Out` field typed as a port (`Ledger notary.Ledger`) is how
  a concrete adapter is bound to its interface. The constructors those providers call stay plain
  functions in their own packages, so no package outside `cmd/` and `internal/wiring` imports dig
  (`bin/check-boundaries` enforces it) and every adapter is tested by calling its constructor. No
  container is passed around. `tests/wiring` resolves the whole graph so a wiring mistake fails a test
  rather than a start-up.

## The house failure mode, applied here

A failure must never be reported as a definite answer. In this service that means:

- `verifyInclusion` returns `false` for a bad proof and **reverts** for an epoch that does not exist.
- The verifier distinguishes `not_anchored`, `not_final`, `unknown` (could not ask) and
  `wrong_ledger` from `rejected`, and refuses to give a verdict when the contract and the Go scheme
  disagree.
- The verifier checks it is talking to a ledger on the right chain **before** it reads an epoch, so
  a wrong address can never produce "epoch not found".
- A contract that answers with nothing, or reverts, when asked for its interface is `wrong_ledger`,
  not `unknown`: the call succeeded and the answer was no.
- A JWKS outage is `identity_unavailable` (503), never 401. `tokens` records whether identity's last JWKS answer
  was readable (its transport and keyfunc's refresh error handler), and a token whose key could not be found
  (`jwt.ErrTokenUnverifiable`) while that answer failed is an outage. A key already read keeps verifying, and a bad
  signature stays 401. Only a 401 carries `WWW-Authenticate`.
- In `httpapi` only what the handler parsed itself is 400. A `KindInvalidArgument` from behind the request (the
  issuer, the registrar, the proof codec) is the notary's own fault: 500 `internal`, its message never shown.
- `registration_not_open` is 409, a state like `epoch_not_final`, not an outage.
- What the store could never hold is refused by the handler before any lookup, because behind the request it reads as
  the notary's own fault. An epoch is `latest` or a decimal from 0 to `proving.MaxEpochNumber` (9223372036854775807,
  the BIGINT the store keys by) without a leading zero, else 400 `invalid_epoch`. An evidence reference is 1 to 255
  bytes of UTF-8 without a NUL byte, else 400 `invalid_reference`; the store also refuses such a name as
  `KindInvalidArgument` before it asks Postgres. Until 2026-10-22 an epoch past the BIGINT answered 500 `internal`, and
  a `%00` or invalid UTF-8 reference reached Postgres, which refuses the text (SQLSTATE 22021), and answered 503
  `unavailable`.

## The HTTP API

Every answer is JSON, errors as `{"error", "message"}`. A path nothing serves is 404 `not_found`, a served path asked
with another method is 405 `method_not_allowed` with `Allow`. `httpapi.Routes` is the list the router serves and the
metrics seed; `tests/httpapi` holds the two to each other. Only GET and HEAD of an unknown path are counted, as
`unmatched`, because the method label of anything else is whatever the caller sends.

`GET /api/v1/proofs/registrations` (bearer token) answers the EIP-712 typed data a holder signs before they post a
registration: `types`, `primaryType`, `domain` (the withholding court's name, version, `chainId` as a JSON number,
`verifyingContract`) and `message.principal`, the token's subject. The client adds `holder` and `issuedAt` and signs it
as it is (`eth_signTypedData_v4`). The domain is `keccak`'s, which `CheckWithholding` proves equal to the court's at
start; `tests/httpapi` hashes the served template with go-ethereum's `apitypes` against `keccak`, and the integration
test signs from it and has the court accept the registration. The signature binds the principal: the handler puts the
token's subject into the request, so a request signed for another principal recovers to another address and is
`registration_refused`.

## Traps found while building the foundation

- **`vm.prank` is consumed by the first external call, including one made while evaluating an
  argument.** `vm.prank(admin); ledger.seal(ledger.head())` pranks `head()`, and `seal` runs as the
  test contract. Read the argument into a local first. The invariant suite caught this.
- **`vm.setEnv` is process-global and Foundry runs tests in parallel**, so two tests setting the same
  variable race. Only one test touches the environment; the deploy script separates `fromEnvironment`
  from `deploy` so everything else is tested without it.
- **forge-lint reads tests too**, and flags `block.timestamp` after `vm.warp`; tests use
  `vm.getBlockTimestamp()`. It also flags `(node, false)`-style returns as boolean-constant misuse,
  which is why the sum tree checks the path sum up front instead of returning an `ok` flag.
- **abigen v2's generated `UnpackError` slices `raw[:4]` without a length check** and panics on short
  revert data. `internal/ethereum` never calls it; it compares the selector itself after checking the
  length.
- **anvil's `finalized` block trails `latest` by 64 blocks.** The integration suite mines 65 blocks to
  finalize, and uses the gap to test `not_final`.
- **A call into a contract that is not a ledger can return empty data with no error.** Identity
  probes use `CallContract` directly and require exactly one ABI word back.
- **forge 1.8 migrated `soldeer.lock` to format v2** on the first install. The committed file is v2.
- **`docker compose run --build` without a TTY writes BuildKit's log to stdout.** Piped into
  `kinetix-secrets set`, the first log line would have been stored as the anchorer's private key, with
  no error. Every script now builds the image first (`build --quiet`, to stderr) and runs without
  `--build`, and `bin/new-anchorer-key` refuses anything that is not exactly `0x` + 64 hex. A stub
  store caught this on 2026-10-12, before the real store was ever touched.
- **A stub proves the pipe, not the tool.** The stub store above passed while the real
  `kinetix-secrets`, called from this repository, decrypted the prod store, could not find
  `.sops.yaml` (sops searches from the caller's working directory) and left every production secret
  in plaintext on 2026-10-12. Re-encrypted within a minute, verified identical to `HEAD`, nothing
  rotated; fixed in infra with an explicit `--config` and a bats test, and `bin/new-anchorer-key` now
  calls the tool from its own repository. Anything that writes a real store is rehearsed against the
  real tool with a throwaway store and age key before it touches one.
- **forge 1.8.3 broadcasts a salted `new` as plain CREATE on chain 84532, whatever the config says.**
  The first ledger went out that way on 2026-10-12, and `always_use_create_2_factory` was turned on as
  the fix; the first court went out the same way on 2026-10-13 with the setting on. forge routes a
  salted `new` through the factory on 31337 only: the same script, keystore or unlocked, sends CREATE2
  on `anvil` and CREATE on `anvil --chain-id 84532`, so every rehearsal on the lab chain looked right.
  The script's `== Return ==` still shows the CREATE2 address, where nothing ever lands. Both deploy
  scripts now call the factory themselves (`script/Create2Factory.sol`), and `bin/rehearse-deploy`
  replays them on anvil with each network's own chain id and fails unless the broadcast is one CREATE2
  through the factory. Old scripts swapped in turn it red. Both deployed contracts keep their CREATE
  addresses; only their successors are predictable.
- **forge reports verification optimistically.** It printed "All (1) contracts were verified!" over
  Blockscout's `Fail - Unable to verify` (Blockscout had no solc 0.8.37). Verification goes to
  Sourcify, and `bin/inspect-ledger` asks Sourcify itself.
- **A pinned Debian base is not a patched one.** Trivy found three fixable CRITICALs in
  `debian:bookworm-slim`, including its newest digest; the final stage runs `apt-get upgrade`. Scan
  the image locally before the first `build-and-push`, never after.
- **`go install` leaves its caches in the layer.** abigen, protoc-gen-* and golangci-lint left 1.1 GB
  of module cache and 1.1 GB of build cache behind: the toolchain image was 5.06 GB and GitHub CI
  rebuilt it from nothing on every run (3m12s of a 7-minute gate), while `build-and-push` spent
  109–140 s of each build uploading those layers to the Actions cache. Both Dockerfiles now run
  `go clean -cache -modcache` in the same `RUN` (toolchain 2.1 GB, its tool layer 572 MB), and the
  production `go build` keeps its build cache in a `--mount=type=cache`, never in a layer.
- **A `sh -c '…'` typed into the agent's shell tool can lose its script** to the command rewriter,
  and the rest then runs on the Mac. It happened with a multi-line script, and again on 2026-10-12 with
  a one-line script that had pipes inside the quotes: every Postgres test came back `SKIP`, because the
  `go test` that answered ran on the host, where the test database does not exist. Run commands
  directly (`bin/dev go test …`), keep pipes outside, and put anything longer in a mounted file.
- **`cast wallet new --json` changed shape in Foundry 1.8**: the wallet is under `data[0]` inside an
  envelope, not a bare list.
- **A replacement must not be estimated again.** On anvil (and any node that estimates against the
  pending block) `eth_estimateGas` for a second attempt at the same nonce sees the first attempt already
  applied and reverts `EpochAlreadyAnchored`. A stuck anchor is replaced by `Resign`: same calldata, gas
  and nonce, new fees.
- **Reading at `finalized` before the ledger itself is final returns "no contract code",** not
  `UnknownEpoch`. Both mean the epoch is not final yet.
- **A broadcast is not a receipt.** Even anvil mines asynchronously after `eth_sendRawTransaction`;
  tests wait for the receipt, and the daemon treats its absence as pending.
- **anvil's clock is shared and only moves forward.** `evm_increaseTime` from earlier tests and earlier
  runs leaves chain time hours ahead of the wall clock, so a test that steers the daemon's clock starts
  it at the latest block's time, never at `time.Now()`.
- **A dependency days old is not yet a dependency.** client_golang v1.25.0 was four days old on
  2026-10-12; v1.24.1 is pinned until it ages.
- **`vm.expectRevert` is consumed by the next external call, argument evaluation included**, exactly
  like `vm.prank`: `vm.expectRevert(...); court.proveBadProof(s, _sign(key, s), siblings)` expects the
  revert from `court.statementDigest` inside `_sign`. Seven court tests failed that way on 2026-10-13
  until every signature was computed into a local first.
- **A deploy through the CREATE2 factory reverts without data.** A constructor's custom error never
  leaves the factory, and neither does an address already taken: `Create2Factory.deploy` turns both into
  `DeploymentFailed(expected)`, and the unit suites check which constructor error it was. A taken
  address also burns all the gas the call was given, which is why a second deploy test costs a billion.
- **forge fmt's `ignore` is a list of the v1 files, not `src/**`.** The glob would have exempted the
  court too; the frozen files are named one by one.
- **The fork test has its own profile** (`[profile.fork]`, `test = "fork"`), so a plain `forge test`
  stays offline; `bin/ci-contracts` runs it with the public RPC and the pinned block from
  `deploy/base-sepolia.env`.
  `GenerationTwoFork.t.sol` forks a later block, `KINETIX_NOTARY_FORK_GENERATION_BLOCK`, where the live directory
  exists, and `plonk-production.json` is proven for that block's next epoch: moving the block means proving the
  fixture again with the SRS.
- **The ERC-7562 test has its own profile too** (`[profile.rules]`, `test = "rules"`, `verbosity = 3`, its
  own `out/rules` and `cache/rules`). `vm.startDebugTraceRecording` needs a tracer, which forge starts only
  at `-vvv`; the profile carries the verbosity so `FOUNDRY_PROFILE=rules forge test` is all `bin/ci-contracts`
  needs. It records every opcode the paymaster's validation runs for each of the four calls and fails on
  any ERC-7562 bans, and on any storage read of the paymaster or the sender. Controlled on 2026-10-16: a
  planted `block.timestamp` and a planted role read each turned it red.
- **`bin/analyze` never accepts a high finding**, only low ones. Aderyn's `reentrancy-state-change` on
  the court was answered by moving the ledger reads into private views, not by an entry.
- **gosec G115** flags every narrowing integer conversion; `UnixTime` is a sealed type whose
  constructors check the range, so the conversions happen after a visible bound.
- **sqlc cannot analyse a multi-argument `unnest`** (`function unnest(unknown, unknown) does not exist`).
  One `unnest` per column in the select list zips the same way and generates.
- **Aderyn's `unused-import` does not count a type used only in an error's parameters**, and its
  `reentrancy-state-change` flags a constructor that assigns immutables after an external call. The
  withholding court's errors carry no state (`requestOf` reads it), and its constructor reads the court
  through a private view and assigns `ledger` last. Both were fixed in code, nothing was accepted.
- **`forge … --root contracts` run from the repository root writes `cache/` there**, fuzz and invariant
  failure records included. Delete it; it is not `contracts/cache/`, which is ignored.
- **A zero account secret means the first derivation.** `Opening.Commitment` picks `H(0x04 ‖ kind ‖ salt ‖
  holder)` for an epoch built before secrets and `H(0x06 ‖ H(0x05 ‖ kind ‖ secret ‖ holder) ‖ salt)` after,
  the way an all-zero salt is refused everywhere: the source never draws one.
- **An RPC endpoint limits `eth_getLogs`, and nothing local shows it.** sepolia.base.org answers at most
  201 blocks per query and refuses more with HTTP 413. notaryd v1.0.0 asked for 2 000, and every pass
  over the withholding court failed in production from its first minute. anvil answers any range, and so
  does the fork test, because anvil serves a forked chain's logs itself. The span is now
  `KINETIX_NOTARY_RPC_LOG_RANGE`, and notaryd probes one full span before it starts. Measure a new
  endpoint's limit with `curl` before you point notaryd at it.
- **A request the notary cannot answer is never retried in silence.** It is logged as an error and
  counted in `kinetix_notary_withholding_unanswerable_requests`. An outage is retried. A refusal for
  one request never holds up the others, because each lost request becomes someone's recorded silence.

- **The laptop has a Foundry of its own, 1.7.1 in `~/.foundry`, and a `sh -c` the rewriter splits runs it.**
  On 2026-10-14 `forge fmt` meant for the toolchain ran on the Mac instead: it ignored the frozen-file
  list, rewrapped `src/SumTree.sol` and fifteen other tracked files, and reformatted the Soldeer
  dependencies. The tracked files were restored from `HEAD` and `dependencies/`, `out/` and `cache/`
  reinstalled in the image. Run forge as `docker compose -f toolchain/compose.yaml run --rm --no-deps -T
  -w /workspace/contracts toolchain forge …`, and read `git status` after any formatter.
- **Slither taints a whole storage struct when one field comes from `block.timestamp`.** `pending.members.length
  == 0` was reported as a dangerous strict equality because the same struct holds `activatesAt`. The council
  asks a private view, `_hasPending()`, which compares with `!=`.
- **Aderyn reads a storage array handed to a `memory` parameter as an edit that will not persist** — a high
  finding, even for a read. Copy it into a `memory` local first. Its `costly-loop` flags any storage write
  inside a loop, which is why the council keeps no membership mapping: it assigns the whole list and scans
  at most sixteen addresses.
- **forge-lint refuses a `revert` inside a loop.** The loop returns the first offending index and the revert
  follows it.
- **A deploy that reverts inside `vm.startBroadcast` leaves the broadcast open**, so a second expected revert
  in the same test fails with "a broadcast is active already". One refused deploy per test.

- **A client-streaming call never returns `io.EOF`.** The interceptor that counts gRPC calls counted on
  `io.EOF` or an error, so every successful `ApproveEpoch` went uncounted; the daemon test that asserts
  the watch call found it. A stream that is not server-streaming now counts on its one successful
  `RecvMsg`.
- **Asking the watch about an epoch already on chain abandons it.** The watch checks the proposal against
  the ledger's next epoch, which has moved on; notaryd would read the refusal as a disagreement. `submit`
  reads the ledger first and moves a landed epoch on without asking.

- **A keystore needs a terminal.** `cast send --keystore` asks for the password on the TTY, so a script
  that signs must not run its container with `-T`. The first `bin/hand-over-anchoring` did, and on
  2026-10-14 it stopped with `No such device or address (os error 6)` before signing anything. Read-only
  scripts keep `-T`; the deploy and funding scripts never had it.

- **`cast call f()(a,b)` and `cast call f()((a,b))` print differently.** The first decodes several return
  values and prints one per line; the second decodes one tuple and prints `(a, b)` on one line. The first
  draft of `bin/inspect-directory` cut a tuple by commas out of the per-line form and reported every field
  `WRONG`; the rehearsal on a fork caught it. It reads `current()` per line now.
- **`forge create --constructor-args` swallows every flag after it.** Put it last, or `--rpc-url` becomes a
  constructor argument and forge dials `localhost:8545`.
- **`cast call` reports a custom error as a bare selector.** `bin/register-successor` names the refusal by
  looking the selector up in `forge inspect LedgerDirectory errors --json`, so the operator reads
  `PredecessorNotSealed(address)`, not `custom error 0x92422182`.

- **Base Sepolia hands back a receipt before `latest` holds the transaction.** sepolia.base.org answers
  `cast send` with a preconfirmed receipt (`blockHash` all zeros), and a `cast call` made at once still
  reads the old state. The first `bin/hand-over-anchoring` checked once and reported a successful grant
  as failed on 2026-10-14; both admin scripts now wait up to 60 s for `latest` to show the change.

- **`type(I).interfaceId` leaves out what `I` inherits; the ABI does not.** `INotaryPaymaster is IPaymaster`,
  so its ABI carries `validatePaymasterUserOp` and `postOp`, and XORing every method in it gave an id the
  contract does not announce: the adapter called the deployed paymaster "not a notary paymaster". Every unit
  test passed, because the fake has no interface; the integration test against the deployed contract caught
  it. The adapter XORs out the methods of the generated `IPaymaster` binding.
- **A JSON-RPC quantity has no leading zero.** The bundler client sent the authorization's `yParity` as
  `0x01`; Pimlico's Alto parses anything hex and would have taken it, a geth-strict parser refuses it. The
  local test bundler decodes with geth's `hexutil` and refused the first sponsored request. `yParity`, `r` and
  `s` are now quantities (`0x1`), as the execution API specifies.
- **gosec G101 reads `pass` in a constant's name as a credential.** The constant for
  `KINETIX_NOTARY_HOLDER_PASSWORD_FILE` is `HolderUnlockFileKey`; the variable keeps its name.
- **A dependency's compiler warnings fail `deny = "warnings"`.** account-abstraction v0.8.0 warns in its own
  files; `ignored_warnings_from = ["dependencies/account-abstraction-0.8.0"]` excludes that path and nothing else.
- **Aderyn's high `contract-locks-ether` fires on any `payable` with no way out.** The paymaster's payable
  `deposit` was removed; funding goes straight to `EntryPoint.depositTo`, and `withdrawTo` takes it back.
  The stake cannot go that way, because the EntryPoint stakes its caller, so `addStake` is payable. Aderyn
  counts neither `entryPoint.addStake{value: …}` nor OpenZeppelin's `Address.sendValue` as a way out, only
  a native `call`, `send` or `transfer` in the contract itself; Slither calls such a send in an `onlyRole`
  function `arbitrary-send-eth`, because it does not read OpenZeppelin's role check, and a `call`
  `low-level-calls`. `withdrawBalance` checks `defaultAdmin()` inline and pays it with a checked `send`,
  clean in both, and it is also the only way out for ether forced in by a `selfdestruct`. Probed
  2026-10-16 on copies, eight shapes: only that one passed both.

- **`cast codehash` asks `eth_getProof`, and an anvil fork answers it for an account it has not loaded yet
  with the empty account's hash** (`0xc5d2…a470`). The first `bin/inspect-paymaster` on a fork reported the
  canonical EntryPoint `WRONG` until a transaction had touched it. The scripts hash `cast code` instead,
  which every endpoint answers.
- **A fork rehearsal writes `contracts/broadcast/<script>/84532/` exactly as a real deploy does**, and only
  `31337/` and `dry-run/` are ignored. Delete it after a rehearsal, or it records a deploy that never
  happened on the real chain.
- **forge's incremental build can keep a test's embedded creation code stale.** A test that deploys with
  `new NotaryPaymaster(...)` carries the paymaster's creation code. After an edit to the paymaster alone,
  forge recompiled two files and the tests still ran the old bytecode: a planted `block.timestamp` passed
  the opcode test until `--force`. A positive control, and any build whose bytecode decides an address,
  runs from `forge build --force`.
- **EntryPoint v0.8 takes an operation while `validAfter < now <= validUntil`.** v0.7 accepted `now ==
  validAfter`; v0.8 does not. The court records a silence only when `blockTime > dueAt`, so the paymaster's
  `validAfter` is `dueAt` itself, and a request's `validUntil` is the court's `closesAt`, which it still takes
  at that very second. The rules tests warp to each edge.
- **abigen v2's `UnpackError` reads `raw[:4]` without a length check.** `revertData` returns revert data only
  when it has a selector, and the paymaster adapter decodes nothing else.

- **A file named `*_wasm_test.go` is a `GOARCH=wasm` file.** Go reads the `_wasm` suffix as a build constraint,
  so the first browser-verifier test, `verify_wasm_test.go`, compiled into nothing and `go test -run` reported
  `no tests to run`. It is `browser_verifier_test.go`. The same goes for `_js`, `_linux` and every other OS or
  architecture suffix.
- **Go turns `fetch` off under Node.** `net/http` on `js/wasm` uses the Fetch API unless `process.argv0` starts
  with `node` (go.dev/issue/57613), and then falls back to a network stack that cannot dial. The integration test
  runs the WebAssembly under Bun, whose `argv0` is `bun`, through the same `wasm_exec.js` a browser loads; Bun is a
  toolchain component for that reason. `net/http` also skips `fetch` for any `Transport` with a `Dial` or
  `DialContext` set, so the RPC client in the verifier must keep the default transport.
- **The browser verifier is checked against the CLI, not against itself.** `tests/integration/browser_verifier_test.go`
  builds both binaries, runs them on the same documents and settings on anvil, and requires the same exit code,
  stdout and stderr for every verdict. A planted rewrite of the document inside the WebAssembly turned it red on
  2026-10-16. Two builds with separate caches gave the same SHA-256.
- **On GitHub the checkout is not the toolchain's.** The runner owns `/workspace` (uid 1001) and the toolchain runs
  as root, so git refuses the repository (`dubious ownership`, exit 128), and anything that stamps VCS on a main
  package fails: `go build` and `go list` both exit 1. Docker Desktop maps ownership, so the laptop never sees it.
  `go list` then prints nothing at all, internal packages included, and `bin/check-boundaries` piped it into
  `awk`, so on GitHub the dig, go-ethereum, gRPC and Postgres rules checked an empty list and printed
  `boundaries hold` from the day CI started. A planted dig import in `internal/ethereum` passed the old script in
  a copy owned by uid 1001 on 2026-10-16. Every listing is now captured before it is filtered, so a failed
  `go list` stops the script, and the one listing that includes `cmd/` passes `-buildvcs=false`; the planted
  import fails the new script, and the new script without that flag exits 1. The browser verifier test builds
  `notary-verify` with `-buildvcs=false` for the same reason: it was the first `go build` in the suite, and it went
  red on its first CI run.

- **Postgres keeps microseconds; .NET writes 100 ns ticks.** An evidence record's time is part of its
  digest, so a `timestamptz` column would round order's `occurred_at` and every later proof would
  reach another leaf. The record tables keep seconds and nanoseconds as two integers.
- **A batch built is a batch in flight, even when sending it failed.** The closer's first draft went on
  to the next stream after `begin` built a batch whose fee quote was refused; the next stream then hit
  the one-in-flight index and logged a collision. A built batch now ends the pass, success or not;
  the fee-ceiling test asserts that no second batch was tried.
- **Attempts with the same `sent_at` sort by hash.** `BatchAttempts` orders by `(sent_at, tx_hash)`, and
  keep-alive takes the last as the latest. A test that never moves its clock gets two attempts in an
  arbitrary order; real clocks never tie at microsecond precision.
- **A config error wrapped in a kind is a verdict, not a config error.** The evidence verifier's
  provider first wrapped the loader's error as `invalid_argument`, so a missing registry printed
  `unknown` and exited 4 where a missing directory exits 78. The browser test, which compares the two
  binaries case by case, caught it.
- **The Ignition SRS everyone uses is the sealed transcript, not the last participant's.** The manifest's `crs.t2`
  equals the τ of `MAIN IGNITION/sealed/transcript00.dat`. The last participant's (position 197) differs, because a
  sealing step multiplied τ once more. `Follows` holds from the seal back to the generator.
- **A curve point in a PLONK proof that is off the curve burns all the gas it is given.** gnark's verifier calls the
  ecMul and ecAdd precompiles with `gas()`, and a failing precompile consumes everything it was forwarded before the
  verifier reverts. The first tamper test exhausted Foundry's 1.07 billion gas. Tests cap each call at 1M. The ledger
  does not cap it, because a cap in an immutable contract breaks when a fork reprices the precompiles.
- **gnark logs to stdout.** Its zerolog console writer prints coloured lines (`compiling circuit`, …) on standard
  output, and the first child prover's proof came back as `invalid byte: U+001B`. `notary-prove` calls
  `logger.Disable()`: its stdout carries the proof and nothing else.
- **Binary from `curl` through the hook is text.** The rewriter turned a range read of a transcript into UTF-8
  replacement characters, and the header decoded as nonsense. Fetch with `rtk proxy curl -o <file>` and read the
  file.
- **A second daemon in one test waits out the first one's lease.** The writer's instance name ends in random bytes and
  the lease lasts 30 s, so a restarted notaryd ticks `nil` and builds nothing until the old lease lapses. The cutover
  test expires `writer_lease` itself, as a dead process would.
- **An epoch in flight at a seal blocks the successor.** The sealed ledger refuses it (`LedgerIsSealed` is
  `anchor_refused`), it is abandoned, and its row keeps the number the successor's first epoch needs. The cutover
  stops notaryd between two epochs first (`kinetix-infrastructure/docs/NOTARY-GENERATION-TWO.md`, step 0).
- **Tests freeze the clock that times the prover,** so `kinetix_notary_proving_seconds` reads 0 there. The Ignition
  test logs wall time instead.
- **`bin/inspect-directory` checks the generation that holds the file's ledger, not the current one.** Comparing with
  the current generation would have reported generation 1's own file `WRONG` from the moment it was succeeded.
- **`COPY --chmod` gives the directories it creates the file's mode.** v1.6.0 copied the SRS with `--chmod=0444`, and
  BuildKit made `/usr/local/share/kinetix-notary` `dr--r--r--` too: the user `notary` could not enter it, and the
  cutover's first start refused with `CheckSRSFile … permission denied`. Scheme 1 never opens the file, so v1.6.0 ran
  green on it for an hour. The final stage creates the directory with `install -d -m 0755` before the copy, and
  build-and-push opens the SRS in the built image as `notary` before it pushes.
- **A public RPC endpoint limits requests per IP, and every loop ticked on the same 15 s grid.** sepolia.base.org
  allows 25 a second. On 2026-10-20 notaryd and notary-watch, on one host, drew 49 `429`s: each idle tick sent about
  21 calls in under a second, the watch's readiness 6 more, and every restart re-read 100 000 blocks of the
  withholding court in 500 `eth_getLogs` back to back. Nothing paced or retried them. Four fixes, one release:
  - Every request is paced and a 429 retried (above).
  - The paymaster's identity is read once. All of it is `immutable`, and re-reading it was 8 of the 21 calls.
  - The withholding scan keeps its position and its open requests in Postgres (`withholding_scans`,
    `withholding_requests`, migration 00008), so a restart resumes where the last pass stopped. A request leaves the
    table only when it is closed on chain, past its deadline, or answered by a mined transaction. A cursor kept
    without the requests would have lost every request found before the restart.
  - The watch's readiness reads `nextEpoch()` once instead of the whole head.

  `ProofReader.Depth` reported a 429 as `wrong_ledger`. Only a revert or an address without code is that now; any
  other failure is `ledger_unavailable`.
- **A fake server that drops what it cannot place hides the reader's refusal.** The first evidence feed
  fake skipped a record with no position, so the reader's "without a position" refusal was never
  reached. Fakes send what they are given; the reader refuses.

## Not built yet

Steps 1–13 are live; PLAN.md records each one's release and proof. The last, scheme 2, has run in production since
2026-10-20 on notary v1.6.1: generation 2 from epoch 183, its council `0xB165…A45A` the ledger's only anchorer. What
the plan still holds is outside this repository:
- 6b's delivery itself is notification's (decision 14). The notary side, which decides who is told and when, has been
  live since 2026-10-20 in notary v1.7.0 with notification v1.1.1, on kinetix-contracts v1.7.0;
- step 14, the product passport, is a service of its own, and it starts with owner-side work in warehouse and buyer
  wallets.

Generation 2 holds 256 accounts. Growing past that is a new generation (decision 46), and
`kinetix-notary-capacity` warns at 75%. Contract deploys stay manual, never in CI.
