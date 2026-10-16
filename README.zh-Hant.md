# 🎲 FlakeHunter

<p align="center">
  <strong>零依賴 · 離線優先 · 跨平台的不穩定測試偵測 CLI</strong><br/>
  <em>Zero-dependency, offline-first flaky test detective for your CI pipelines.</em>
</p>

<p align="center">
  <a href="https://github.com/gitstq/flakehunter/actions/workflows/ci.yml"><img alt="CI" src="https://github.com/gitstq/flakehunter/actions/workflows/ci.yml/badge.svg"/></a>
  <a href="https://github.com/gitstq/flakehunter/releases"><img alt="Release" src="https://img.shields.io/github/v/release/gitstq/flakehunter?color=blue"/></a>
  <a href="https://github.com/gitstq/flakehunter/blob/main/LICENSE"><img alt="License" src="https://img.shields.io/badge/license-MIT-green"/></a>
  <a href="https://golang.org/"><img alt="Go" src="https://img.shields.io/badge/Go-1.24%2B-00ADD8?logo=go"/></a>
  <a href="https://github.com/gitstq/flakehunter"><img alt="Zero Dependency" src="https://img.shields.io/badge/dependencies-zero-brightgreen"/></a>
</p>

<p align="center">
  🌐 語言 / Language:
  <a href="README.md">简体中文</a> ·
  <a href="README.zh-Hant.md"><b>繁體中文</b></a> ·
  <a href="README.en.md">English</a>
</p>

---

## 🎉 專案介紹

**FlakeHunter** 是一個零依賴、離線優先的命令列工具，專門捕捉 CI 流程中最隱蔽的「幽靈」——**不穩定測試（flaky tests）**。它直接讀取本機的測試結果檔案（無需將任何資料上傳雲端），用統計方法告訴你：**哪些測試在反覆跳針、跳得有多頻繁、有多大的把握不是偶然**。

### 😫 解決的痛點

- 💸 **信任危機**：一個十萬行的專案裡混著三個不穩定測試，每次 CI 紅燈大家的第一反應都是「再跑一次」，真正的程式碼退化反而被淹沒；
- 🕵️ **定位困難**：測試今天綠、明天紅、後天又綠，靠人眼翻 Jenkins 紀錄幾乎不可能發現規律；
- ☁️ **SaaS 門檻**：BuildPulse、DeFlaky 等商業方案要求把測試資料傳到第三方伺服器，還要按人頭收費；
- 🔒 **合規限制**：金融、醫療等敏感產業嚴禁測試資料離開內網，只能在本機分析。

### ✨ 與現有方案的差異化

| 方案 | 形態 | FlakeHunter 的做法 |
| --- | --- | --- |
| BuildPulse / DeFlaky | SaaS 儀表板 | **純本機 CLI**，資料不出本機，零訂閱費用 |
| flakefinder | 依賴 CI API 拉歷史 | **解析本機產物檔案**，任何 CI 系統都能用 |
| quarantine（框架專用） | 綁定單一框架 | **格式無關**：JUnit XML / `go test -json` / TAP / Jest / Vitest 通吃 |
| 人工重跑大法 | 碰運氣 | **Wilson 信賴區間**統計評分，區分「真的不穩定」與「偶然失敗」 |

### 🌱 靈感來源

本專案受 2026 年 10 月 GitHub Trending 三大趨勢啟發：**AI 能力下沉到具體工具**（不堆框架，專注解決一個問題）、**終端機工具回潮**（輕量、低開銷、高可客製）、**本地優先**（資料與計算不出本機）。FlakeHunter 把這三股力量擰成一顆螺絲釘：一個單一二進位檔，解決不穩定測試這個真實而頑固的工程痛點。

---

## ✨ 核心特性

- 🔌 **多格式讀取**：自動辨識 JUnit XML（pytest / Jest / Maven Surefire / Gradle / TestNG / xUnit / RSpec）、`go test -json`、TAP v12/13、Jest/Vitest `--json` 報告
- 📊 **Wilson 信賴區間評分**：不只是「失敗過幾次」，而是回答「這個測試真的在跳針嗎？」——輸出不穩定分數與信心水準（high/med/low）
- 🎯 **單輪啟發式偵測**（無需歷史）：
  - 🎲 `rerun-flip`：同一輪裡先失敗後通過（重試痕跡）= 不穩定實錘
  - ⏱️ `timing-spike`：耗時超過同輪中位數四倍以上 = 時序敏感型不穩定
- 🔗 **順序依賴偵測**：用 phi 相關係數找出「只有跟在某個測試後面才失敗」的順序敏感型不穩定
- 🗂️ **隔離（Quarantine）工作流**：`add` 暫時隔離 → `sync` 穩定後自動解除 → `patch` 產生 pytest / Jest / Go / JUnit / xUnit 五框架的略過程式碼片段
- 🚦 **CI 品質閘道**：`gate` 指令依不穩定率/數量門檻決定結束碼，並輸出 GitHub Actions 原生註解（`::warning` / `::error`）與 Step Summary Markdown
- 🪶 **真的零依賴**：100% Go 標準函式庫，沒有 `go.sum`，單一二進位檔約 8MB，`go install` 即裝即用
- 🔒 **離線優先**：所有歷史資料存在本機 `.flakehunter/history.jsonl`，適合內網與合規環境
- 🎨 **三種輸出**：彩色終端機表格、機器可讀 JSON、GitHub Markdown，並尊重 `NO_COLOR` 環境變數

---

## 🚀 快速開始

### 📋 環境需求

| 情境 | 需求 |
| --- | --- |
| 使用預先編譯的二進位檔 | 無！Windows / macOS / Linux 皆可，amd64 / arm64 |
| 從原始碼建置 | Go 1.24+ |
| 整合進 CI | 任何 CI 系統（GitHub Actions / GitLab CI / Jenkins / 內網 Runner） |

### 📥 安裝（三選一）

**方式一：下載 Release 二進位檔（推薦）**

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

> Windows 使用者請到 [Releases 頁面](https://github.com/gitstq/flakehunter/releases)下載 `flakehunter-windows-amd64.exe`，解壓縮後加入 `PATH` 即可。下載後可對照 `checksums.txt` 校驗 SHA256。

**方式二：go install**

```bash
go install github.com/gitstq/flakehunter@latest
```

**方式三：從原始碼建置**

```bash
git clone https://github.com/gitstq/flakehunter.git
cd flakehunter
go build -o flakehunter .
```

### ⚡ 30 秒上手

```bash
# 1️⃣ 用專案內建的範例資料體驗一次掃描（無需任何測試框架）
flakehunter scan testdata

# 2️⃣ 初始化本機狀態目錄
flakehunter init

# 3️⃣ 記錄一輪執行（把 CI 產出的測試結果餵給它）
flakehunter record ./test-results

# 4️⃣ 查看不穩定測試統計報告
flakehunter report

# 5️⃣ 設定 CI 閘道：不穩定率超過 5% 就讓流程紅燈
flakehunter gate --max-flake-rate 5
```

### 🎬 示範

<details>
<summary>📸 終端機報告效果（佔位演示，正式 GIF 錄製中）</summary>

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

## 📖 詳細使用指南

### 1️⃣ 讓測試框架產出結果檔案

FlakeHunter 不跑測試，只讀結果。各框架的匯出方式：

```bash
# pytest → JUnit XML
pytest --junitxml=test-results/junit.xml

# go test → JSON 事件串流
go test -json ./... > test-results/gotest.jsonl

# Jest / Vitest → JSON 報告
jest --json --outputFile=test-results/jest.json
vitest run --reporter=json --outputFile=test-results/vitest.json

# 任何輸出 TAP 的測試（Node tap、prove、部分 BDD 框架）
prove --archive t/ > test-results/tap.txt
```

### 2️⃣ 指令詳解

#### `flakehunter scan <paths...>` — 單輪診斷

解析給定的檔案或目錄（目錄會遞迴尋找 `.xml` / `.json` / `.jsonl` / `.tap`），輸出本輪概覽與不穩定訊號。**scan 永遠回傳結束碼 0**，它是診斷工具而非閘道。

```bash
flakehunter scan ./test-results
flakehunter scan ./test-results --record          # 同時記入歷史
flakehunter scan junit.xml --out json             # 機器可讀輸出
flakehunter scan ./test-results --out markdown --summary "$GITHUB_STEP_SUMMARY"
```

#### `flakehunter record <paths...>` — 記錄歷史

把一輪執行追加到本機歷史（`.flakehunter/history.jsonl`，預設保留最近 100 輪）。**CI 裡每輪建置 record 一次**，累積夠多輪次就能產出統計報告。

```bash
flakehunter record ./test-results
```

#### `flakehunter report` — 不穩定測試統計報告

```bash
flakehunter report                                  # 彩色終端機表格
flakehunter report --out json                       # JSON（可接 jq / 儀表板）
flakehunter report --out markdown --summary sum.md  # GitHub Markdown
```

狀態圖例：✅ `stable_pass` 穩定通過 · 🎲 `flaky` 不穩定 · 💥 `stable_fail` 持續失敗 · ❓ `insufficient` 資料不足（不足 `min_runs` 輪）。

#### `flakehunter gate` — CI 品質閘道

```bash
flakehunter gate --max-flake-rate 5      # 不穩定率上限 5%
flakehunter gate --max-flaky 3           # 不穩定測試個數上限 3
flakehunter gate --fail-on-broken        # 持續失敗的測試也判紅
flakehunter gate --quiet                 # 只輸出違規與註解
```

結束碼：`0` = 通過，`1` = 觸發門檻。門檻缺省讀取 `.flakehunter/config.json`（預設不穩定率 5%、不穩定個數 0）。

#### `flakehunter quarantine` — 隔離工作流

```bash
# 暫時隔離一個不穩定測試（從 report 輸出裡複製指紋）
flakehunter quarantine add "tests.test_net#test_fetch" --reason "timeout on CI" --framework pytest

# 查看隔離清單
flakehunter quarantine list

# 穩定 5 輪後自動解除隔離
flakehunter quarantine sync

# 產生略過程式碼片段
flakehunter quarantine patch "tests.test_net#test_fetch" --framework pytest
# 輸出：
# tests.test_net#test_fetch
# @pytest.mark.skip(reason="flakehunter: quarantined flaky test")
# def test_fetch(...):
```

被隔離的測試**不參與 gate 判定**，但 report 裡仍會標註 `(quarantined)`，防止「隔離了就被遺忘」。

### 3️⃣ GitHub Actions 完整整合範例

把測試結果上傳為 artifact，然後在下一個 job 裡跨輪記錄與把關：

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

> 💡 **提示**：`gate` 輸出的 `::warning title=FlakeHunter::...` 註解會直接出現在 Checks 面板，`--summary` 內容會渲染在 job summary 頁。

### 4️⃣ 設定說明

`flakehunter init` 會產生 `.flakehunter/config.json`：

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

| 欄位 | 預設值 | 說明 |
| --- | --- | --- |
| `max_runs` | 100 | 本機歷史保留輪數（自動輪換） |
| `min_runs` | 3 | 少於該輪數的測試判為「資料不足」 |
| `max_flake_rate` | 5 | gate 預設不穩定率上限（百分比） |
| `max_flaky` | 0 | gate 預設不穩定個數上限 |
| `framework` | pytest | `quarantine patch` 的預設框架 |
| `quarantine_stable_runs` | 5 | 連續穩定多少輪後自動解除隔離 |

> 🔐 **安全說明**：`.flakehunter/` 目錄只含本機歷史與設定，**不應提交進版控**（`.gitignore` 已預設排除）。全流程無需任何 Token 或網路連線。

---

## 💡 設計思路與迭代規劃

### 🧭 設計理念

1. **把複雜度藏在內部，把簡單留給使用者**——一條 `flakehunter record` 加一條 `flakehunter gate` 就是完整工作流；
2. **統計先行，拒絕拍腦袋**——用 Wilson 得分區間區分「真的不穩定」與「偶然失敗」，兩輪一敗不會冤枉任何測試；
3. **格式無關，管線友好**——解析層（ingest）與引擎層（engine）完全解耦，新框架只需加一個解析器；
4. **離線優先，合規友好**——沒有任何網路呼叫的程式碼路徑，內網環境開箱即用。

### 🛠️ 技術選型

- **Go + 純標準函式庫**：單一二進位檔、跨平台交叉編譯、啟動毫秒級、零依賴供應鏈風險；
- **JSONL 追加式儲存**：崩潰安全、易於 `grep`/`jq` 除錯、天然支援輪換；
- **Phi 相關係數**：順序依賴偵測的經典統計量，計算量極小但解釋性強；
- **Angular 提交規範 + 三平台 CI**：確保歷史可讀、品質把關雲端化。

### 🗺️ 迭代規劃

- **v1.1**：新增 xUnit.NET / TestNG / Cypress JSON 解析器；`report --html` 靜態報告頁
- **v1.2**：`flakehunter bisect`——結合 git 歷史自動定位「從哪個 commit 開始不穩定」
- **v1.3**：`flakehunter mcp`——以 MCP Server 形態向 AI Agent 暴露不穩定測試查詢能力
- **v2.0**：多儲存庫聚合視圖（本機聚合，不上雲）

### 🤗 社群貢獻方向

歡迎從這些方向參與：新測試框架解析器、更多啟發式規則（如資源洩漏偵測）、報告視覺化、文件翻譯（日/韓/西語版）……

---

## 📦 打包與部署指南

### 🎁 Release 產物

每次打下 `v*` 標籤，GitHub Actions 會在雲端自動建置六個平台產物並發布到 Release：

| 檔案 | 平台 |
| --- | --- |
| `flakehunter-linux-amd64` | Linux x86_64 |
| `flakehunter-linux-arm64` | Linux ARM64（樹莓派 / Graviton） |
| `flakehunter-darwin-amd64` | macOS Intel |
| `flakehunter-darwin-arm64` | macOS Apple Silicon |
| `flakehunter-windows-amd64.exe` | Windows x64 |
| `flakehunter-windows-arm64.exe` | Windows ARM64 |
| `checksums.txt` | 全部產物的 SHA256 校驗和 |

### 🔍 校驗下載完整性

```bash
sha256sum -c checksums.txt --ignore-missing
```

### 🔨 自行交叉編譯

```bash
# Linux ARM64
GOOS=linux GOARCH=arm64 go build -o flakehunter-linux-arm64 .

# Windows ARM64
GOOS=windows GOARCH=arm64 go build -o flakehunter-windows-arm64.exe .
```

### 🐳 在 CI 中的部署建議

- **GitHub Actions**：`go install` 或下載 Release 二進位檔皆可，詳見上方完整範例；
- **GitLab CI**：在 `after_script` 中 `flakehunter record`，用 `artifacts:when:always` 保留 `.flakehunter/`；
- **Jenkins**：用 `archiveArtifacts` 歸檔 `.flakehunter/history.jsonl`，實現跨建置累積。

---

## 🤝 貢獻指南

### 🐛 Issue 回饋規則

- 回報 bug 請附上：`flakehunter version` 輸出、觸發指令、去識別化後的輸入檔案片段、期望與實際行為
- 提出新功能建議請說明：解決的痛點、為什麼現有指令無法覆蓋、是否願意提交 PR

### 🔧 PR 提交規範

採用 **Angular 提交規範**：

```text
feat: 新增 xUnit.NET 解析器
fix: 修正 TAP 指令大小寫敏感問題
docs: 補充 GitLab CI 整合範例
refactor: 抽出 Wilson 評分共用函式
test: 補充順序依賴偵測邊界案例
```

PR 檢查清單（CI 會自動執行）：

1. ✅ `gofmt -l .` 無輸出
2. ✅ `go vet ./...` 通過
3. ✅ `go test ./... -race` 全綠（三平台）
4. ✅ 新功能附上單元測試與 testdata 範例
5. ✅ 涉及使用者可見變更時同步更新三語 README

### 🏗️ 本機開發

```bash
git clone https://github.com/gitstq/flakehunter.git
cd flakehunter
gofmt -l .          # 格式檢查
go vet ./...        # 靜態檢查
go test ./... -race # 執行測試
go build -o flakehunter . && ./flakehunter scan testdata
```

---

## 📄 開源授權條款

本專案基於 **MIT License** 開源——你可以自由地使用、修改、散布和商用，唯一要求是在副本中保留版權聲明與條款文本。詳見 [LICENSE](LICENSE)。

---

<p align="center">
  🎲 <b>願你的 CI 從此只因真正的退化而紅燈。</b><br/>
  <sub>Built with Go standard library only · No cloud · No tracking · No flaky tolerance</sub>
</p>
