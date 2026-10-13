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
account holder prove their balance is inside it without trusting Kinetix. It never owns a balance,
an order or a review — it reads a snapshot from the owner and commits to it — and it never moves
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

The image built from `Dockerfile` carries three binaries and runs `notaryd` by default:
`notary-migrate` is the one-shot migrator, and an account holder verifies with
`docker run <image> notary-verify proof.json` (the image used to start `notary-verify` itself).
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
```

`notary-verify` reads `KINETIX_NOTARY_RPC_URL`, `KINETIX_NOTARY_CHAIN_ID`, `KINETIX_NOTARY_LEDGER_ADDRESS`,
`KINETIX_NOTARY_COURT_ADDRESS` and `KINETIX_NOTARY_RPC_TIMEOUT`, every one required. It trusts the court
it is configured with, never the one a document names, and reports beside every verdict whether the
proof's statement is attributable to the epoch's signer and how many faults the epoch carries.

`deploy/<network>.env` holds only public values and is committed. Keys never are: the deployer and
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
| `internal/verifying` | `internal/notary` only |
| `internal/keccak`, `prooffile`, `config`, `verifycli` | `internal/notary` + their own library |
| `internal/ethereum` | the above + go-ethereum + `internal/contractgen` — the only adapter that touches a chain |
| `internal/payment`, `internal/mesh` | `internal/notary` + gRPC + `internal/protogen` — the only adapters that speak gRPC |
| `internal/postgres` | `internal/notary` + pgx + goose + `internal/storegen` — the only adapter that touches Postgres |
| `internal/anchoring`, `internal/proving` | `internal/notary` only, like `internal/verifying` |
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
  reports, exit 1; removed, clean.
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
  Slither, which tells a STATICCALL apart, reports nothing there). An accepted entry that stops being
  reported also fails, so the list cannot rot. forge-lint's `reentrancy-events` is off in `foundry.toml`
  for the same reason and with the same backstop: it flags any event after any external call.
- **One lint setting is a decision, not a silencing.** `wrapcheck` ignores `.RecvMsg(`: the gRPC
  client interceptor that counts calls must hand the stream's own error back unchanged, because callers
  compare `io.EOF` by identity and a wrapped end of stream reads as a failure.
- **The differential fixture is a contract between the languages.** `tests/fixtures` regenerates
  `contracts/test/fixtures/scheme-v1.json` and fails when the committed file differs; Foundry's
  `Differential.t.sol` re-verifies every tree, proof, rejection and digest in it. Proven on
  2026-10-11 by mutating Go's leaf tag: Foundry went red.
- **Every adapter carries a compile-time assertion** (`var _ notary.Ledger = (*Ledger)(nil)`).
- **Generated bindings are never committed.** `bin/generate-bindings` builds them from the forge
  artifacts into `internal/contractgen/`; the Dockerfile and `bin/ci` run it every time.

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
  `contracts/test/`, Foundry's own convention.
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
- **`bin/analyze` never accepts a high finding**, only low ones. Aderyn's `reentrancy-state-change` on
  the court was answered by moving the ledger reads into private views, not by an entry.
- **gosec G115** flags every narrowing integer conversion; `UnixTime` is a sealed type whose
  constructors check the range, so the conversions happen after a visible bound.

## Not built yet

Steps 1–5 are done: contracts v1.0.23, payment v0.1.11, notaryd v0.2.0 anchoring on Base Sepolia, and
the proof API (v0.3.0, gateway v1.0.0) live and proven end to end in the lab on 2026-10-13. Step 6a —
the statement court and a signed statement on every proof — is live: the production court
`0xE1Ffcb004Dd734F83ED0B194fdaFF5B0c7c5927B` and notaryd v0.4.0, proven in production on 2026-10-13. 6b (delivery through notification) and steps 7–14 do not exist yet. Contract
deploys stay manual, never in CI.
