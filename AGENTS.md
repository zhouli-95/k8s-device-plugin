# Repository Guidelines

A Kubernetes device plugin implementing the kubelet gRPC API. Supports `native` and `share` plugin modes with three device injection strategies (`native`, `cdi-cri`, `cdi-annotation`).

## Project Structure

- `cmd/device-plugin/` — Entry point (urfave/cli).
- `internal/` — Plugin interface, plugin manager, native/share plugins, CDI generation, config, CLI options.
- `pkg/resource/` — Device topology loader; reads JSON config from `/config/config.json`.
- `build/package/` — Dockerfile.
- `deployments/kustomize/` — Kustomize base and overlays (share, CDI).
- `deployments/static/` — Pre-rendered static manifest.
- `test/` — Separate Go module for e2e tests (Ginkgo v2 + Gomega).
- `scripts/` — Shell helpers for e2e and image loading.
- `examples/` — Sample device configuration YAMLs.

## Build, Test, and Development

```bash
make build       # Compile binary to out/
make docker      # Build Docker image (device-plugin:dev)
make test        # Run unit tests
make fmt         # Format code with go fmt
make cover       # Generate coverage report (out/coverage.html)
make deploy      # Deploy to cluster via kustomize
make prepare     # Install ginkgo and kind for e2e
make e2e         # Run full e2e suite
```

Run targeted e2e suites with label filters:
```bash
ginkgo --label-filter=basic test/e2e
ginkgo --label-filter="device-strategy || custom-domain-resource" test/e2e
```

## Coding Style & Conventions

- **Go 1.24.** Format with `go fmt` (tabs, standard Go style).
- CI enforces formatting via `make fmt && git diff --exit-code`.
- Follow Go naming conventions — `PascalCase` exports, `camelCase` unexported, fully capitalized acronyms (`CDI`, `GRPC`).
- Public APIs in `pkg/`, internal implementation in `internal/`.

## Testing Guidelines

- **Unit tests:** `go test`; run with `make test`.
- **E2e tests:** Ginkgo v2 in `test/e2e/`, organized by labels (`basic`, `config`, `device-strategy`, `custom-domain-resource`).
- Shared helpers in `test/utils/` (deploy, k8s client, bash wrappers).
- Coverage: `make cover` → `out/coverage.html`.

## Commit & Pull Request Guidelines

- Short, imperative commit subjects (e.g., "add ci", "go mod tidy").
- CI requires passing `build`, `fmt`, `test`, `docker`, and `e2e` jobs.
- PRs should describe what changed and why; link related issues.

## Architecture Overview

The plugin loads a JSON device topology, creates gRPC servers under `/var/lib/kubelet/device-plugins/`, and registers with kubelet. `PluginManager` selects native or share mode at startup. `Allocate` handles device injection per the configured strategy.

## Agent-Specific Instructions

- Always run `make fmt` after editing Go files.
- Run `make test` to validate unit test changes.
- The `test/` directory is a separate Go module — manage its dependencies independently.
