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
bin/dev bin/ci                                  # the whole gate, exactly as CI runs it
bin/dev sh -c 'cd contracts && forge test'      # contracts only
bin/dev bin/test                                # Go suite with true coverage
bin/dev go test ./tests/fixtures -update        # regenerate the differential fixtures
bin/dev down                                    # stop the chain
```

Operations on a real network have their own scripts, all through the same image:

```sh
bin/new-keystore notary-deployer|notary-admin   # encrypted keystore outside every repository
bin/new-anchorer-key prod                       # hot key straight into the SOPS store
bin/deploy-ledger base-sepolia [--broadcast]    # dry run by default; reads deploy/<network>.env
bin/inspect-ledger base-sepolia                 # every on-chain fact checked against the config
```

`deploy/<network>.env` holds only public values and is committed. Keys never are: the deployer and
admin keystores live in `KINETIX_NOTARY_KEYSTORE_DIR` (refused when inside a git work tree), the
anchorer lives in SOPS. Why the split is in PLAN.md §8.

Every tool version lives in **one** place, `bin/install-toolchain`, installed by component
(`foundry`, `abigen`, `lint`, `analysis`): Foundry and Aderyn by release tarball and SHA-256, abigen
and golangci-lint by `go install @version`, Slither from a hash-locked requirements file. The
production image installs only `foundry abigen`. GitHub CI runs `bin/ci` inside the toolchain image
itself (`sh bin/dev bin/ci`), so the laptop and GitHub run the same Python, Go and tools. The toolchain image, the GitHub job and the GitLab job all run that
script, so local and CI cannot drift. Bump a version there and nowhere else.

## Layout, and the rule it encodes

| Package | May import |
|---|---|
| `internal/notary` | the standard library. **Nothing else.** |
| `internal/verifying` | `internal/notary` only |
| `internal/keccak`, `prooffile`, `config`, `verifycli` | `internal/notary` + their own library |
| `internal/ethereum` | the above + go-ethereum + `internal/contractgen` — the only home of go-ethereum |
| `internal/wiring` | everything; the only package outside `cmd/` allowed to import dig |
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
- **`bin/ci`** is the single gate: `forge fmt --check`, `forge build` (forge-lint with
  `deny = "warnings"`), `forge test` on the CI profile, `forge snapshot --check`, binding generation,
  `gofmt`, `go mod verify`, `go vet`, the boundary check, golangci-lint, and `bin/test`. It refuses to
  run without a node, because a skipped integration test proves nothing.
- **`bin/analyze`** runs Slither and Aderyn over `contracts/src` and fails on anything not already
  accepted. Accepted, and why: Slither `timestamp` (refusing a future snapshot needs
  `block.timestamp`; a validator moves it by seconds), Aderyn `centralization-risk` (the admin and
  anchorer roles are the design; the admin key never goes online), Aderyn `unchecked-return` on the
  constructor's `_grantRole` (it cannot fail there, and the deployed ledger's source is frozen: an edit
  would make the repository disagree with the verified bytecode). An accepted entry that stops being
  reported also fails, so the list cannot rot.
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
- **Explicit `return`** in every Solidity function, named returns included, as in the Rust service.
- **One type per file** in both languages: each Solidity struct, interface, library and contract has
  its own file; each Go type in `internal/notary` has its own file.
- **Braces on the same line**; long calls one argument per line. For Solidity this is
  `forge fmt` with `multiline_func_header = "params_first_multi"`: a header that does not fit breaks
  its parameters one per line and keeps `) … {` together. The default (`attributes_first`) puts `{`
  alone on its own line, and that mismatch turned CI red on 2026-10-11 the moment the code was
  formatted by hand to the house style. Fix the formatter's configuration, never the house style.
- **Go tests live in `tests/`**, never beside the code; use `bin/test`, because plain
  `go test -cover` reports 0% for every package under `internal/`. Foundry tests live in
  `contracts/test/`, Foundry's own convention.
- **Hand-written fakes** (`tests/notarytest`), no mocking framework.
- **Errors carry a `notary.Kind`**, read back with `notary.KindOf`. Only the edge (`verifycli`) turns a
  kind into a verdict and an exit code.
- **dig is used only to compose.** Constructors are plain functions taking their dependencies;
  `internal/wiring` registers them, binding concrete types to ports with `dig.As`. No `dig.In`
  structs, no named values, no container passed around. `tests/wiring` resolves the whole graph so a
  wiring mistake fails a test rather than a start-up.

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
- **Without `always_use_create_2_factory`, a salted `new` in a broadcast is plain CREATE.** The first
  ledger went out that way; `forge test` could not notice, because a salted `new` inside a test is
  always CREATE2. The setting is on in `foundry.toml` now. Prove a deploy's address from the
  broadcast record (`to` must be the factory), not from the script's return value.
- **forge reports verification optimistically.** It printed "All (1) contracts were verified!" over
  Blockscout's `Fail - Unable to verify` (Blockscout had no solc 0.8.37). Verification goes to
  Sourcify, and `bin/inspect-ledger` asks Sourcify itself.
- **A pinned Debian base is not a patched one.** Trivy found three fixable CRITICALs in
  `debian:bookworm-slim`, including its newest digest; the final stage runs `apt-get upgrade`. Scan
  the image locally before the first `build-and-push`, never after.
- **A multi-line `sh -c '…'` typed into the agent's shell tool can lose its script** to the command
  rewriter, and the remaining lines then run on the Mac. Trial scripts go in a file that is mounted
  into the toolchain container.
- **`cast wallet new --json` changed shape in Foundry 1.8**: the wallet is under `data[0]` inside an
  envelope, not a bare list.
- **gosec G115** flags every narrowing integer conversion; `UnixTime` is a sealed type whose
  constructors check the range, so the conversions happen after a visible bound.

## Not built yet

Step 2 is done: the ledger is live on Base Sepolia at `0x4E869A55031DE3B0d030AEBECc66E523686761a2`
and verified on Sourcify (PLAN.md step 2). From step 3 on nothing exists yet: no
payment snapshot RPC, no `notaryd`, no database, no proof API, no compose entry in
`kinetix-infrastructure`. The GitHub and GitLab pipelines are written and **have never run**.

**CI/CD runs on GitHub** (the user's standing choice); GitLab only mirrors CI. GitHub still owes
`build-and-push.yml`, `release.yml`, a Slither/Aderyn job and the branch ruleset — the list and its
constraints are in PLAN.md, "Owed next to step 2". Contract deploys stay manual, never in CI.
