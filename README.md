# 🎲 FlakeHunter

<p align="center">
  <strong>零依赖 · 离线优先 · 跨平台的 flaky 测试侦探 CLI</strong><br/>
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
  🌐 语言 / Language:
  <a href="README.md"><b>简体中文</b></a> ·
  <a href="README.zh-Hant.md">繁體中文</a> ·
  <a href="README.en.md">English</a>
</p>

---

## 🎉 项目介绍

**FlakeHunter** 是一个零依赖、离线优先的命令行工具，专门用来捕捉 CI 流水线里最隐蔽的"幽灵"——**flaky tests（不稳定测试）**。它直接读取本地测试结果文件（无需上传任何数据到云端），用统计学方法告诉你：**哪些测试在反复横跳、跳得有多频繁、有多大概率不是偶然**。

### 😫 解决的痛点

- 💸 **信任危机**：一个 10 万行的仓库里混着 3 个 flaky 测试，每次 CI 红屏大家第一反应都是"再跑一次"，真正的回归反而被淹没；
- 🕵️ **定位困难**：测试今天绿、明天红、后天又绿，靠人眼翻 Jenkins 记录几乎不可能发现规律；
- ☁️ **SaaS 门槛**：BuildPulse、DeFlaky 等商业方案要求把测试数据传到第三方服务器，还要按席位付费；
- 🔒 **合规限制**：金融、医疗等敏感行业严禁测试数据出域，只能本地分析。

### ✨ 与现有方案的差异化

| 方案 | 形态 | FlakeHunter 的做法 |
| --- | --- | --- |
| BuildPulse / DeFlaky | SaaS 仪表盘 | **纯本地 CLI**，数据不出本机，零订阅费 |
| flakefinder | 依赖 CI API 拉历史 | **解析本地产物文件**，任何 CI 系统都能用 |
| quarantine（框架专用） | 绑定单一框架 | **格式无关**：JUnit XML / `go test -json` / TAP / Jest / Vitest 通吃 |
| 人工重跑大法 | 碰运气 | **Wilson 置信区间**统计评分，区分"真 flaky"与"偶然失败" |

### 🌱 灵感来源

本项目受 2026 年 10 月 GitHub Trending 三大趋势启发：**AI 能力下沉到具体工具**（不堆框架，专注解决一个问题）、**终端工具回潮**（轻量、低开销、高可定制）、**本地优先**（数据与计算不出本机）。FlakeHunter 把这三股力量拧成一颗螺丝钉：一个单二进制文件，解决 flaky 测试这一个真实而顽固的工程痛点。

---

## ✨ 核心特性

- 🔌 **多格式摄取**：自动识别 JUnit XML（pytest / Jest / Maven Surefire / Gradle / TestNG / xUnit / RSpec）、`go test -json`、TAP v12/13、Jest/Vitest `--json` 报告
- 📊 **Wilson 置信区间评分**：不只是"失败过几次"，而是回答"这个测试真的在翻转吗？"——输出 flake score + 置信度（high/med/low）
- 🎯 **单跑启发式检测**（无需历史）：
  - 🎲 `rerun-flip`：同一轮里先失败后通过（重试痕迹）= 实锤 flaky
  - ⏱️ `timing-spike`：耗时超过同轮中位数 4 倍以上 = 时序敏感型 flaky
- 🔗 **顺序依赖检测**：用 phi 相关系数找出"只有跟在某个测试后面才失败"的顺序敏感型 flaky
- 🗂️ **Quarantine 工作流**：`add` 临时隔离 → `sync` 稳定后自动解除 → `patch` 生成 pytest / Jest / Go / JUnit / xUnit 五框架的跳过代码片段
- 🚦 **CI 质量门禁**：`gate` 命令按 flake 率/数量阈值决定退出码，并输出 GitHub Actions 原生注解（`::warning` / `::error`）与 Step Summary Markdown
- 🪶 **真·零依赖**：100% Go 标准库，无 `go.sum`，单二进制约 8MB，`go install` 即装即用
- 🔒 **离线优先**：所有历史数据存在本地 `.flakehunter/history.jsonl`，适合内网与合规环境
- 🎨 **三种输出**：彩色终端表格、机器可读 JSON、GitHub Markdown，尊重 `NO_COLOR` 环境变量

---

## 🚀 快速开始

### 📋 环境要求

| 场景 | 要求 |
| --- | --- |
| 使用预编译二进制 | 无！Windows / macOS / Linux 均可，amd64 / arm64 |
| 从源码构建 | Go 1.24+ |
| 集成到 CI | 任何 CI 系统（GitHub Actions / GitLab CI / Jenkins / 内网 Runner） |

### 📥 安装（三选一）

**方式一：下载 Release 二进制（推荐）**

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

> Windows 用户请到 [Releases 页面](https://github.com/gitstq/flakehunter/releases)下载 `flakehunter-windows-amd64.exe`，解压后加入 `PATH` 即可。下载后可对照 `checksums.txt` 校验 SHA256。

**方式二：go install**

```bash
go install github.com/gitstq/flakehunter@latest
```

**方式三：从源码构建**

```bash
git clone https://github.com/gitstq/flakehunter.git
cd flakehunter
go build -o flakehunter .
```

### ⚡ 30 秒上手

```bash
# 1️⃣ 用仓库自带的样本数据体验一次扫描（无需任何测试框架）
flakehunter scan testdata

# 2️⃣ 初始化本地状态目录
flakehunter init

# 3️⃣ 记录一次运行（把 CI 产出的测试结果喂给它）
flakehunter record ./test-results

# 4️⃣ 查看 flake 统计报告
flakehunter report

# 5️⃣ 设置 CI 门禁：flake 率超过 5% 就让流水线红屏
flakehunter gate --max-flake-rate 5
```

### 🎬 演示

<details>
<summary>📸 终端报告效果（占位演示，正式 GIF 录制中）</summary>

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

## 📖 详细使用指南

### 1️⃣ 让测试框架产出结果文件

FlakeHunter 不跑测试，只读结果。各框架导出方式：

```bash
# pytest → JUnit XML
pytest --junitxml=test-results/junit.xml

# go test → JSON 事件流
go test -json ./... > test-results/gotest.jsonl

# Jest / Vitest → JSON 报告
jest --json --outputFile=test-results/jest.json
vitest run --reporter=json --outputFile=test-results/vitest.json

# 任何输出 TAP 的测试（Node tap、prove、部分 BDD 框架）
prove --archive t/ > test-results/tap.txt
```

### 2️⃣ 命令详解

#### `flakehunter scan <paths...>` — 单跑诊断

解析给定的文件或目录（目录会递归查找 `.xml` / `.json` / `.jsonl` / `.tap`），输出本轮概览 + flake 信号。**scan 永远返回退出码 0**，它是诊断工具而非门禁。

```bash
flakehunter scan ./test-results
flakehunter scan ./test-results --record          # 同时记入历史
flakehunter scan junit.xml --out json             # 机器可读输出
flakehunter scan ./test-results --out markdown --summary "$GITHUB_STEP_SUMMARY"
```

#### `flakehunter record <paths...>` — 记录历史

把一轮运行追加到本地历史（`.flakehunter/history.jsonl`，默认保留最近 100 轮）。**CI 里每轮构建 record 一次**，积累够了就能出统计报告。

```bash
flakehunter record ./test-results
```

#### `flakehunter report` — flake 统计报告

```bash
flakehunter report                                  # 彩色终端表格
flakehunter report --out json                       # JSON（可接 jq / 仪表盘）
flakehunter report --out markdown --summary sum.md  # GitHub Markdown
```

状态图例：✅ `stable_pass` 稳定通过 · 🎲 `flaky` 不稳定 · 💥 `stable_fail` 持续失败 · ❓ `insufficient` 数据不足（不足 `min_runs` 轮）。

#### `flakehunter gate` — CI 质量门禁

```bash
flakehunter gate --max-flake-rate 5      # flake 率上限 5%
flakehunter gate --max-flaky 3           # flaky 测试个数上限 3
flakehunter gate --fail-on-broken        # 持续失败的测试也判红
flakehunter gate --quiet                 # 只输出违规与注解
```

退出码：`0` = 通过，`1` = 触发阈值。阈值缺省读 `.flakehunter/config.json`（默认 flake 率 5%、flaky 个数 0）。

#### `flakehunter quarantine` — 隔离工作流

```bash
# 临时隔离一个 flaky 测试（从 report 输出里复制指纹）
flakehunter quarantine add "tests.test_net#test_fetch" --reason "timeout on CI" --framework pytest

# 查看隔离清单
flakehunter quarantine list

# 稳定 5 轮后自动解除隔离
flakehunter quarantine sync

# 生成跳过代码片段
flakehunter quarantine patch "tests.test_net#test_fetch" --framework pytest
# 输出：
# tests.test_net#test_fetch
# @pytest.mark.skip(reason="flakehunter: quarantined flaky test")
# def test_fetch(...):
```

被隔离的测试**不参与 gate 判定**，但 report 里仍会标注 `(quarantined)`，防止"隔离了就被遗忘"。

### 3️⃣ GitHub Actions 完整集成示例

把测试结果上传为 artifact，然后在下一个 job 里跨轮记录与门禁：

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

> 💡 **提示**：`gate` 输出的 `::warning title=FlakeHunter::...` 注解会直接出现在 Checks 面板，`--summary` 内容会渲染在 job summary 页。

### 4️⃣ 配置说明

`flakehunter init` 生成 `.flakehunter/config.json`：

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

| 字段 | 默认 | 说明 |
| --- | --- | --- |
| `max_runs` | 100 | 本地历史保留轮数（自动轮换） |
| `min_runs` | 3 | 少于该轮数的测试判为"数据不足" |
| `max_flake_rate` | 5 | gate 默认 flake 率上限（百分比） |
| `max_flaky` | 0 | gate 默认 flaky 个数上限 |
| `framework` | pytest | `quarantine patch` 的默认框架 |
| `quarantine_stable_runs` | 5 | 连续稳定多少轮后自动解除隔离 |

> 🔐 **安全说明**：`.flakehunter/` 目录只含本地历史与配置，**不应提交入库**（`.gitignore` 已默认排除）。全流程无需任何 Token 或网络访问。

---

## 💡 设计思路与迭代规划

### 🧭 设计理念

1. **把复杂度藏进内部，把简单留给用户**——一条 `flakehunter record` + 一条 `flakehunter gate` 就是完整工作流；
2. **统计先行，拒绝拍脑袋**——用 Wilson 得分区间区分"真 flaky"与"偶然失败"，2 轮 1 败不会冤枉任何测试；
3. **格式无关，管道友好**——解析层（ingest）与引擎层（engine）完全解耦，新框架只需加一个解析器；
4. **离线优先，合规友好**——没有任何网络调用代码路径，内网环境开箱即用。

### 🛠️ 技术选型

- **Go + 纯标准库**：单二进制、跨平台交叉编译、启动毫秒级、零依赖供应链风险；
- **JSONL 追加式存储**：崩溃安全、易于 `grep`/`jq` 调试、天然支持轮换；
- **Phi 相关系数**：顺序依赖检测的经典统计量，计算量极小但解释性强；
- **Angular 提交规范 + 三平台 CI**：保证历史可读、质量门禁云端化。

### 🗺️ 迭代规划

- **v1.1**：新增 xUnit.NET / TestNG / Cypress JSON 解析器；`report --html` 静态报告页
- **v1.2**：`flakehunter bisect`——结合 git 历史自动定位"从哪个 commit 开始 flaky"
- **v1.3**：`flakehunter mcp`——以 MCP Server 形态向 AI Agent 暴露 flake 查询能力
- **v2.0**：多仓库聚合视图（本地聚合，不上云）

### 🤗 社区贡献方向

欢迎从这些方向参与：新测试框架解析器、更多启发式规则（如资源泄漏检测）、报告可视化、文档翻译（日/韩/西语版）……

---

## 📦 打包与部署指南

### 🎁 Release 产物

每次打 `v*` 标签，GitHub Actions 会在云端自动构建 6 个平台产物并发布到 Release：

| 文件 | 平台 |
| --- | --- |
| `flakehunter-linux-amd64` | Linux x86_64 |
| `flakehunter-linux-arm64` | Linux ARM64（树莓派 / Graviton） |
| `flakehunter-darwin-amd64` | macOS Intel |
| `flakehunter-darwin-arm64` | macOS Apple Silicon |
| `flakehunter-windows-amd64.exe` | Windows x64 |
| `flakehunter-windows-arm64.exe` | Windows ARM64 |
| `checksums.txt` | 全部产物的 SHA256 校验和 |

### 🔍 校验下载完整性

```bash
sha256sum -c checksums.txt --ignore-missing
```

### 🔨 自行交叉编译

```bash
# Linux ARM64
GOOS=linux GOARCH=arm64 go build -o flakehunter-linux-arm64 .

# Windows ARM64
GOOS=windows GOARCH=arm64 go build -o flakehunter-windows-arm64.exe .
```

### 🐳 在 CI 中的部署建议

- **GitHub Actions**：`go install` 或下载 Release 二进制均可，见上文完整示例；
- **GitLab CI**：在 `after_script` 中 `flakehunter record`，用 `artifacts:when:always` 保留 `.flakehunter/`；
- **Jenkins**：用 `archiveArtifacts` 归档 `.flakehunter/history.jsonl`，实现跨构建累积。

---

## 🤝 贡献指南

### 🐛 Issue 反馈规则

- 报告 bug 请附：`flakehunter version` 输出、触发命令、脱敏后的输入文件片段、期望 vs 实际行为
- 提新功能建议请说明：解决的痛点、为什么现有命令无法覆盖、是否愿意提交 PR

### 🔧 PR 提交规范

采用 **Angular 提交规范**：

```text
feat: 新增 xUnit.NET 解析器
fix: 修复 TAP 指令大小写敏感问题
docs: 补充 GitLab CI 集成示例
refactor: 抽取 Wilson 评分公共函数
test: 补充顺序依赖检测边界用例
```

PR 检查清单（CI 会自动执行）：

1. ✅ `gofmt -l .` 无输出
2. ✅ `go vet ./...` 通过
3. ✅ `go test ./... -race` 全绿（三平台）
4. ✅ 新功能附带单元测试与 testdata 样本
5. ✅ 涉及用户可见变更时同步更新三语 README

### 🏗️ 本地开发

```bash
git clone https://github.com/gitstq/flakehunter.git
cd flakehunter
gofmt -l .          # 格式检查
go vet ./...        # 静态检查
go test ./... -race # 运行测试
go build -o flakehunter . && ./flakehunter scan testdata
```

---

## 📄 开源协议说明

本项目基于 **MIT License** 开源——你可以自由地使用、修改、分发和商用，唯一要求是在副本中保留版权声明与协议文本。详见 [LICENSE](LICENSE)。

---

<p align="center">
  🎲 <b>愿你的 CI 从此只因真正的回归而红屏。</b><br/>
  <sub>Built with Go standard library only · No cloud · No tracking · No flaky tolerance</sub>
</p>
