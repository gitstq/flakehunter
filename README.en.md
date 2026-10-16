# 🎲 FlakeHunter

<p align="center">
  <strong>Zero-dependency, offline-first flaky test detective for CI</strong><br/>
  <em>One single binary. No cloud. No subscriptions. Just statistics.</em>
</p>

<p align="center">
  <a href="https://github.com/gitstq/flakehunter/actions/workflows/ci.yml"><img alt="CI" src="https://github.com/gitstq/flakehunter/actions/workflows/ci.yml/badge.svg"/></a>
  <a href="https://github.com/gitstq/flakehunter/releases"><img alt="Release" src="https://img.shields.io/github/v/release/gitstq/flakehunter?color=blue"/></a>
  <a href="https://github.com/gitstq/flakehunter/blob/main/LICENSE"><img alt="License" src="https://img.shields.io/badge/license-MIT-green"/></a>
  <a href="https://golang.org/"><img alt="Go" src="https://img.shields.io/badge/Go-1.24%2B-00ADD8?logo=go"/></a>
  <a href="https://github.com/gitstq/flakehunter"><img alt="Zero Dependency" src="https://img.shields.io/badge/dependencies-zero-brightgreen"/></a>
</p>

<p align="center">
  🌐 Language / 语言:
  <a href="README.md">简体中文</a> ·
  <a href="README.zh-Hant.md">繁體中文</a> ·
  <a href="README.en.md"><b>English</b></a>
</p>

---

## 🎉 Introduction

**FlakeHunter** is a zero-dependency, offline-first command-line tool built to hunt down the ghosts of your CI pipeline: **flaky tests**. It reads local test result files — nothing ever leaves your machine — and applies statistics to answer three questions: **which tests flip-flop, how often, and how confident can we be that it isn't just bad luck?**

### 😫 The pain it solves

- 💸 **Eroded trust**: three flaky tests hiding in a 100k-line repo, and every red build gets dismissed with "just re-run it" — while real regressions slip through;
- 🕵️ **Invisible patterns**: green today, red tomorrow, green again on Friday. No human can spot that pattern by scrolling through build logs;
- ☁️ **SaaS friction**: services like BuildPulse and DeFlaky require shipping your test data to a third party and paying per seat;
- 🔒 **Compliance walls**: finance and healthcare teams are forbidden from sending test artifacts off-premise — local analysis is the only option.

### ✨ How it differs from existing tools

| Tool | Shape | FlakeHunter's take |
| --- | --- | --- |
| BuildPulse / DeFlaky | SaaS dashboards | **Pure local CLI** — data stays on your machine, zero subscription |
| flakefinder | Pulls history via CI APIs | **Parses local artifacts** — works with any CI system |
| quarantine (framework-bound) | Tied to one framework | **Format-agnostic**: JUnit XML / `go test -json` / TAP / Jest / Vitest |
| Manual re-runs | Hope-based | **Wilson confidence intervals** — separates real flakes from unlucky failures |

### 🌱 Inspiration

FlakeHunter is inspired by three strong signals from the GitHub Trending charts of October 2026: **AI capabilities sinking into concrete tools** (not more frameworks, just one problem solved well), **the terminal renaissance** (lightweight, fast, customizable), and **local-first computing** (data and compute stay on your machine). It distills those forces into one tight screw: a single binary that fixes one stubborn, real engineering problem.

---

## ✨ Key features

- 🔌 **Multi-format ingestion**: auto-detects JUnit XML (pytest / Jest / Maven Surefire / Gradle / TestNG / xUnit / RSpec), `go test -json`, TAP v12/13, and Jest/Vitest `--json` reports
- 📊 **Wilson confidence scoring**: not just "failed N times" but "is this test *really* flipping?" — flake score plus a high/med/low confidence label
- 🎯 **Single-run heuristics** (no history required):
  - 🎲 `rerun-flip`: failed then passed within the same run (retry traces) — a smoking gun
  - ⏱️ `timing-spike`: duration more than 4× the run median — a timing-sensitive flake
- 🔗 **Order-dependence detection**: phi correlation finds tests that only fail when a specific predecessor runs first
- 🗂️ **Quarantine workflow**: `add` to park a flaky test → `sync` to auto-release it once stable → `patch` to generate skip snippets for pytest / Jest / Go / JUnit / xUnit
- 🚦 **CI quality gate**: `gate` decides the exit code from flake-rate and flake-count thresholds, and emits native GitHub Actions annotations (`::warning` / `::error`) plus Step Summary markdown
- 🪶 **Truly zero dependencies**: 100% Go standard library — no `go.sum`, one ~8MB binary, `go install` and go
- 🔒 **Offline-first**: history lives in a local `.flakehunter/history.jsonl` file — perfect for air-gapped and regulated environments
- 🎨 **Three output modes**: colored terminal tables, machine-readable JSON, GitHub-flavored markdown; respects `NO_COLOR`

---

## 🚀 Quick start

### 📋 Requirements

| Scenario | Requirement |
| --- | --- |
| Prebuilt binary | None! Windows / macOS / Linux, amd64 / arm64 |
| Build from source | Go 1.24+ |
| CI integration | Any CI system (GitHub Actions / GitLab CI / Jenkins / self-hosted runners) |

### 📥 Install (pick one)

**Option 1: Download a release binary (recommended)**

```bash
# Linux (amd64)
curl -LO https://github.com/gitstq/flakehunter/releases/latest/download/flakehunter-linux-amd64
chmod +x flakehunter-linux-amd64
sudo mv flakehunter-linux-amd64 /usr/local/bin/flakehunter

# macOS (Apple Silicon)
curl -LO https://github.com/gitstq/flakehunter/releases/latest/download/flakehunter-darwin-arm64
chmod +x flakehunter-darwin-arm64
sudo mv flakehunter-darwin-arm64 /usr/local/bin/flakehunter
```

> Windows users: grab `flakehunter-windows-amd64.exe` from the [Releases page](https://github.com/gitstq/flakehunter/releases) and add it to your `PATH`. Verify downloads against `checksums.txt` (SHA256).

**Option 2: go install**

```bash
go install github.com/gitstq/flakehunter@latest
```

**Option 3: Build from source**

```bash
git clone https://github.com/gitstq/flakehunter.git
cd flakehunter
go build -o flakehunter .
```

### ⚡ 30-second tour

```bash
# 1️⃣ Scan the bundled sample data (no test framework needed)
flakehunter scan testdata

# 2️⃣ Initialize the local state directory
flakehunter init

# 3️⃣ Record a run (feed it your CI artifacts)
flakehunter record ./test-results

# 4️⃣ See the flake statistics report
flakehunter report

# 5️⃣ Gate your pipeline: red if the flake rate exceeds 5%
flakehunter gate --max-flake-rate 5
```

### 🎬 Demo

<details>
<summary>📸 Terminal report preview (placeholder — full GIF coming soon)</summary>

```
🎲 FlakeHunter — 12 runs, 47 tests tracked (✅ 41 · 🎲 4 · 💥 1 · ❓ 1)

    RUNS   P/F/E    CONF  FLAKE   MEDIAN   TEST
🎲 12     9/3/0     high  25.0%   120ms    tests.test_net#test_fetch
🎲 12     8/4/0     high  33.3%   45ms     tests.test_auth#test_token_refresh
💥 12     0/12/0    -     -       8ms      tests.test_db#test_migration
❓ 2      2/0/0     -     -       12ms     tests.test_new#test_parser
✅ 12     12/0/0    -     -       3ms      tests.test_math#test_add
```

</details>

---

## 📖 Detailed usage guide

### 1️⃣ Make your test framework emit result files

FlakeHunter never runs your tests — it only reads their results:

```bash
# pytest → JUnit XML
pytest --junitxml=test-results/junit.xml

# go test → JSON event stream
go test -json ./... > test-results/gotest.jsonl

# Jest / Vitest → JSON report
jest --json --outputFile=test-results/jest.json
vitest run --reporter=json --outputFile=test-results/vitest.json

# Anything that speaks TAP (Node tap, prove, some BDD frameworks)
prove --archive t/ > test-results/tap.txt
```

### 2️⃣ Command reference

#### `flakehunter scan <paths...>` — one-off diagnosis

Parses the given files or directories (directories are walked recursively for `.xml` / `.json` / `.jsonl` / `.tap`) and prints a run overview plus flake signals. **scan always exits 0** — it is a diagnostic, not a gate.

```bash
flakehunter scan ./test-results
flakehunter scan ./test-results --record          # also append to history
flakehunter scan junit.xml --out json             # machine-readable output
flakehunter scan ./test-results --out markdown --summary "$GITHUB_STEP_SUMMARY"
```

#### `flakehunter record <paths...>` — build history

Appends one run to the local history (`.flakehunter/history.jsonl`, last 100 runs kept by default). **Record once per CI build** — statistics need a few runs to become meaningful.

```bash
flakehunter record ./test-results
```

#### `flakehunter report` — flake statistics

```bash
flakehunter report                                  # colored terminal table
flakehunter report --out json                       # JSON (pipe to jq / dashboards)
flakehunter report --out markdown --summary sum.md  # GitHub markdown
```

Status legend: ✅ `stable_pass` · 🎲 `flaky` · 💥 `stable_fail` · ❓ `insufficient` (fewer than `min_runs` runs recorded).

#### `flakehunter gate` — CI quality gate

```bash
flakehunter gate --max-flake-rate 5      # cap the flake rate at 5%
flakehunter gate --max-flaky 3           # cap the flaky count at 3
flakehunter gate --fail-on-broken        # also fail on consistently failing tests
flakehunter gate --quiet                 # only violations and annotations
```

Exit codes: `0` = pass, `1` = threshold violated. Defaults come from `.flakehunter/config.json` (5% flake rate, 0 flaky tests).

#### `flakehunter quarantine` — parking workflow

```bash
# Park a flaky test (copy its fingerprint from report output)
flakehunter quarantine add "tests.test_net#test_fetch" --reason "timeout on CI" --framework pytest

# List the quarantine manifest
flakehunter quarantine list

# Auto-release tests that stayed stable for 5 runs
flakehunter quarantine sync

# Generate a skip snippet
flakehunter quarantine patch "tests.test_net#test_fetch" --framework pytest
# Output:
# tests.test_net#test_fetch
# @pytest.mark.skip(reason="flakehunter: quarantined flaky test")
# def test_fetch(...):
```

Quarantined tests are **excluded from gate decisions** but still shown in reports marked `(quarantined)` — so nothing gets forgotten.

### 3️⃣ Full GitHub Actions integration

Upload test results as artifacts, then record and gate across builds:

```yaml
name: Tests
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: '1.24' }
      - run: go test -json ./... > test-results/gotest.jsonl
      - uses: actions/upload-artifact@v4
        with:
          name: test-results-${{ github.run_id }}
          path: test-results/

  flakehunter:
    runs-on: ubuntu-latest
    needs: test
    steps:
      - uses: actions/download-artifact@v4
        with: { merge-multiple: true, path: test-results/ }
      - uses: actions/setup-go@v5
        with: { go-version: '1.24' }
      - run: go install github.com/gitstq/flakehunter@latest
      - name: Record this run
        run: flakehunter record test-results
      - name: Flake gate
        run: |
          flakehunter gate --max-flake-rate 5 --summary "$GITHUB_STEP_SUMMARY"
```

> 💡 **Tip**: `gate` prints `::warning title=FlakeHunter::...` annotations that surface directly in the Checks tab, and `--summary` content renders in the job summary page.

### 4️⃣ Configuration

`flakehunter init` writes `.flakehunter/config.json`:

```json
{
  "max_runs": 100,
  "min_runs": 3,
  "max_flake_rate": 5,
  "max_flaky": 0,
  "framework": "pytest",
  "quarantine_stable_runs": 5
}
```

| Field | Default | Meaning |
| --- | --- | --- |
| `max_runs` | 100 | How many runs of history to keep (auto-rotated) |
| `min_runs` | 3 | Runs needed before a test can be classified |
| `max_flake_rate` | 5 | Default gate threshold for flake rate (percent) |
| `max_flaky` | 0 | Default gate threshold for flaky count |
| `framework` | pytest | Default framework for `quarantine patch` |
| `quarantine_stable_runs` | 5 | Consecutive stable runs before auto-release |

> 🔐 **Security note**: `.flakehunter/` holds only local history and config — **never commit it** (already in `.gitignore`). No tokens, no network access, anywhere in the flow.

---

## 💡 Design rationale & roadmap

### 🧭 Design principles

1. **Hide complexity inside, leave simplicity outside** — one `flakehunter record` plus one `flakehunter gate` is the whole workflow;
2. **Statistics over vibes** — Wilson score intervals tell real flakes apart from unlucky failures, so a single bad run never frames an innocent test;
3. **Format-agnostic, pipe-friendly** — the ingest layer and the engine layer are fully decoupled; a new framework means one new parser;
4. **Offline-first, compliance-friendly** — there is no network call anywhere in the codebase.

### 🛠️ Tech choices

- **Go with the standard library only**: single binary, trivial cross-compilation, millisecond startup, zero supply-chain risk;
- **JSONL append-only storage**: crash-safe, greppable, naturally rotatable;
- **Phi correlation**: a classic, tiny-to-compute statistic with strong explanatory power for order dependence;
- **Angular commits + three-platform CI**: readable history with cloud-side quality gates.

### 🗺️ Roadmap

- **v1.1**: xUnit.NET / TestNG / Cypress JSON parsers; `report --html` static report page
- **v1.2**: `flakehunter bisect` — pinpoint the commit where a test started flaking, using git history
- **v1.3**: `flakehunter mcp` — expose flake queries to AI agents as an MCP server
- **v2.0**: multi-repository aggregation view (local aggregation, nothing uploaded)

### 🤗 Contribution ideas

New framework parsers, more heuristics (e.g. resource-leak detection), report visualizations, documentation translations (Japanese, Korean, Spanish)…

---

## 📦 Packaging & deployment guide

### 🎁 Release artifacts

Every `v*` tag triggers a cloud build matrix that publishes six platform binaries plus checksums to the Release page:

| File | Platform |
| --- | --- |
| `flakehunter-linux-amd64` | Linux x86_64 |
| `flakehunter-linux-arm64` | Linux ARM64 (Raspberry Pi / Graviton) |
| `flakehunter-darwin-amd64` | macOS Intel |
| `flakehunter-darwin-arm64` | macOS Apple Silicon |
| `flakehunter-windows-amd64.exe` | Windows x64 |
| `flakehunter-windows-arm64.exe` | Windows ARM64 |
| `checksums.txt` | SHA256 checksums for all artifacts |

### 🔍 Verify a download

```bash
sha256sum -c checksums.txt --ignore-missing
```

### 🔨 Cross-compile yourself

```bash
# Linux ARM64
GOOS=linux GOARCH=arm64 go build -o flakehunter-linux-arm64 .

# Windows ARM64
GOOS=windows GOARCH=arm64 go build -o flakehunter-windows-arm64.exe .
```

### 🐳 CI deployment notes

- **GitHub Actions**: either `go install` or download a release binary — full example above;
- **GitLab CI**: `flakehunter record` in `after_script`, keep `.flakehunter/` with `artifacts:when:always`;
- **Jenkins**: archive `.flakehunter/history.jsonl` with `archiveArtifacts` to accumulate history across builds.

---

## 🤝 Contributing

### 🐛 Issue guidelines

- Bug reports: include `flakehunter version` output, the exact command, a sanitized input snippet, and expected vs. actual behavior
- Feature requests: describe the pain point, why existing commands don't cover it, and whether you'd submit a PR

### 🔧 Pull request conventions

**Angular commit convention** is enforced:

```text
feat: add xUnit.NET parser
fix: make TAP directive parsing case-insensitive
docs: add GitLab CI integration example
refactor: extract Wilson scoring helper
test: add edge cases for order-dependence detection
```

PR checklist (all enforced by CI):

1. ✅ `gofmt -l .` prints nothing
2. ✅ `go vet ./...` passes
3. ✅ `go test ./... -race` green on all three platforms
4. ✅ New features ship with unit tests and testdata fixtures
5. ✅ User-visible changes update all three README versions

### 🏗️ Local development

```bash
git clone https://github.com/gitstq/flakehunter.git
cd flakehunter
gofmt -l .          # formatting check
go vet ./...        # static analysis
go test ./... -race # run tests
go build -o flakehunter . && ./flakehunter scan testdata
```

---

## 📄 License

Released under the **MIT License** — use it, modify it, ship it, sell it; just keep the copyright notice and license text in copies. See [LICENSE](LICENSE) for the full text.

---

<p align="center">
  🎲 <b>May your CI only ever go red for real regressions.</b><br/>
  <sub>Built with the Go standard library only · No cloud · No tracking · No flaky tolerance</sub>
</p>
