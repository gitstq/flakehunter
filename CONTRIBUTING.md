# 🤝 Contributing to FlakeHunter

Thank you for helping hunt flaky tests! This document explains how to contribute.
感谢参与 FlakeHunter，一起猎杀不稳定的测试！

---

## 🐛 Reporting issues / 反馈问题

Please include / 请附上：

1. `flakehunter version` output / 版本输出
2. The exact command you ran / 你执行的确切命令
3. A **sanitized** snippet of the input file (remove proprietary test names if needed) / 脱敏后的输入文件片段
4. Expected vs. actual behavior / 期望行为与实际行为

## 🔧 Pull requests / 提交 PR

### Commit convention / 提交规范

This project follows the **Angular commit convention** / 本项目采用 Angular 提交规范：

| Type | Use for |
| --- | --- |
| `feat:` | New features / 新功能 |
| `fix:` | Bug fixes / 缺陷修复 |
| `docs:` | Documentation only / 仅文档变更 |
| `refactor:` | Code changes with no behavior change / 无行为变化的重构 |
| `test:` | Adding or fixing tests / 测试增补与修复 |
| `ci:` | CI/workflow changes / CI 与工作流变更 |

Examples:

```text
feat: add xUnit.NET v2 parser
fix: ignore BOM when sniffing JUnit XML
docs: add GitLab CI integration example
test: cover phi coefficient degenerate tables
```

### Checklist / 自检清单

Every PR must pass the cloud CI gates / 每个 PR 都必须通过云端 CI 门禁：

1. `gofmt -l .` prints nothing / 格式检查无输出
2. `go vet ./...` passes / 静态检查通过
3. `go test ./... -race` is green on Linux, macOS and Windows / 三平台测试全绿
4. New features include unit tests **and** a fixture under `testdata/` / 新功能附带单元测试与样本数据
5. User-visible changes update **all three** README versions (`README.md`, `README.zh-Hant.md`, `README.en.md`) / 用户可见变更同步更新三语 README

### Local development / 本地开发

```bash
git clone https://github.com/gitstq/flakehunter.git
cd flakehunter
gofmt -l .            # formatting / 格式
go vet ./...          # static analysis / 静态分析
go test ./... -race   # tests / 测试
go build -o flakehunter .
./flakehunter scan testdata
```

## 🏗️ Architecture in one minute / 一分钟看懂架构

```text
main.go                 CLI entry & dispatch / 入口与分发
cmd/                    subcommands (init/scan/record/report/gate/quarantine/version)
internal/model/         core types + test fingerprint normalization / 核心类型与指纹归一化
internal/ingest/        format detection & parsers (junit/gotest/tap/jest) / 格式探测与解析器
internal/engine/        history store, Wilson scoring, heuristics / 历史存储、评分、启发式
internal/quarantine/    manifest + framework patch snippets / 隔离清单与补丁生成
internal/report/        terminal / JSON / markdown renderers / 三种报告渲染
internal/gate/          CI pass/fail evaluation / CI 门禁评估
```

Adding a new test framework? Write a parser in `internal/ingest/`, register it in `Detect`, add a fixture in `testdata/`, and you are done — the rest of the pipeline is format-agnostic.
想支持新的测试框架？只需在 `internal/ingest/` 加一个解析器、在 `Detect` 注册、在 `testdata/` 放样本——其余管线与格式无关。

## 📄 License / 授权协议

By contributing you agree that your contributions are licensed under the [MIT License](LICENSE).
贡献即表示你的贡献以 [MIT License](LICENSE) 授权。
