# open-sandbox Constitution

## Core Principles

### I. MVP First, Demo-Ready
优先保证可演示、可用。功能选择以“端到端闭环可跑”为第一标准。

### II. Single-Node, Single-Container
MVP 必须在单机/单容器内完成一体化沙箱（Browser/VNC/IDE/Jupyter/Shell/File/Code）。
本机环境是Windows，所有的构建产物和缓存以及临时文件请你留在D盘
采用TDD开发的时候请你自我闭环，可以在wsl中运行测试，wsl在d盘，允许安装任何所需要的工具链

### III. Simplicity Over Cleverness
清晰可维护、少魔法。模块职责单一，避免隐式依赖与隐式行为。

### IV. Test-First (Non-Negotiable)
严格 TDD：先写测试再实现。任何实现必须有对应测试。

### V. Safe-by-Default for MVP
安全与性能优化可延后，但不得留下明显安全洞；任何风险必须在文档中说明。

## Technical Standards

### Language & Dependencies
- 主语言为 Go；API 使用 net/http。
- 依赖最小化，避免引入重量级框架。

### Browser & VNC
- 浏览器使用有头模式。
- 必须支持 CDP 控制与 VNC 接管。

### API & Error Model
- 所有 API 必须有清晰、统一的错误模型与最小日志输出。
- 错误返回需可定位（例如包含错误码/消息/trace_id 或等效字段）。

## Documentation Requirements

文档必须齐全且可落地执行，至少包含：
- README
- 运行方式
- 端口说明
- 环境变量说明
- limitations / TODO（写清所有假设与未决点）

## Development Workflow

- 原子化开发与原子化提交（小步提交，每次变更聚焦单一目标）。
- 遵循协作规范：清晰命名、错误处理与最小日志输出。

## Testing Strategy

- 严格 TDD：先写测试再实现。
- 测试必须覆盖主要路径与边界条件。
- 任何实现都必须有对应测试。

## Governance

本宪章优先级高于其他开发惯例；任何偏离必须在文档中说明并记录原因。

**Version**: 1.0.0 | **Ratified**: 2026-01-22 | **Last Amended**: 2026-01-22
