# Contributing to open-sandbox

Thanks for taking the time to contribute.

Development Principles
- MVP first: demo-ready and usable on a single node/container.
- Atomic changes: one goal per commit.
- Test-first: write tests before implementation (TDD).
- Minimal dependencies; prefer standard library.
- All code comments must be in English and follow best practices:
  - Explain intent/why, not what the code already says
  - Keep comments concise
  - Avoid redundant or obvious comments
  - No TODOs without an issue reference

Workflow
1) Create or update the spec using `/speckit.specify`, then `/speckit.plan` and `/speckit.tasks`.
2) Implement with `/speckit.implement` and keep each change small.
3) Run tests (WSL is acceptable); fix all failures before committing.
4) Commit with a clear message and push.

Code Style
- Prefer clarity over cleverness.
- Avoid hidden side effects.
- Keep modules focused and single-responsibility.

Reporting Issues
- Use GitHub issues and include reproduction steps, logs, and environment info.

License
By contributing, you agree that your contributions will be licensed under Apache-2.0.