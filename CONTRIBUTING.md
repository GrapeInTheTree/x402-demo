# Contributing to x402-playground

Thanks for your interest in contributing!

## Development Setup

```bash
git clone https://github.com/GrapeInTheTree/x402-playground.git
cd x402-playground
cp .env.example .env    # Edit with your keys
make build              # Build all binaries
make test               # Run tests
```

## Running Tests

```bash
make test               # All unit tests
go test ./internal/demo/... -v   # Specific package
```

## Code Style

- Follow standard Go conventions (`gofmt`, `go vet`)
- Use `go vet ./...` before submitting
- Keep error handling explicit — avoid `_, _ =` for I/O operations
- Use named constants instead of magic numbers

## Pull Request Process

1. Fork the repository
2. Create a feature branch (`git checkout -b feat/my-feature`)
3. Write tests for new functionality
4. Ensure `make test` passes
5. Submit a pull request against `main`

## Project Structure

- `cmd/` — Entry points (facilitator, resource, client, explorer, balance)
- `internal/config/` — Environment variable loading
- `internal/demo/` — Shared protocol logic (types, balance, decoder, flow)
- `internal/facilserver/` — Facilitator HTTP handlers
- `internal/server/` — Resource server routes and handlers
- `internal/signer/` — EVM transaction signing
- `internal/tui/` — Bubbletea TUI (components, pages)

## License

By contributing, you agree that your contributions will be licensed under the [MIT License](LICENSE).
