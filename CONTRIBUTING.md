# Contributing to terraform-provider-insightfinder

Thank you for your interest in contributing to the InsightFinder Terraform Provider! We welcome contributions from the community.

## How to Contribute

### Reporting Issues

If you find a bug or have a feature request:

1. Check the [issue tracker](https://github.com/insightfinder/terraform-provider-insightfinder/issues) to see if it's already reported
2. If not, create a new issue with:
   - Clear description of the problem or feature
   - Steps to reproduce (for bugs)
   - Expected vs actual behavior
   - Terraform and provider versions
   - Relevant configuration snippets

### Submitting Pull Requests

1. **Fork the repository** and create your branch from `main`
2. **Make your changes**:
   - Write clear, concise code
   - Follow Go best practices
   - Add or update tests as needed
   - Update documentation
3. **Test your changes**:
   ```bash
   go test ./...
   make install
   # Test with actual Terraform configurations
   ```
4. **Commit your changes**:
   - Use clear, descriptive commit messages
   - Reference issues if applicable (`fixes #123`)
5. **Submit a pull request**:
   - Describe your changes
   - Link to related issues
   - Wait for review

## Development Setup

### Prerequisites

- Go 1.21 or later
- Terraform 1.0 or later
- Git
- Make

### Building the Provider

```bash
git clone https://github.com/insightfinder/terraform-provider-insightfinder
cd terraform-provider-insightfinder
go build
```

### Installing Locally

```bash
make install
```

### Running Tests

```bash
# Run unit tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -v ./internal/provider -run TestProjectResource
```

### Testing with Terraform

Create a `.terraformrc` file:

```hcl
provider_installation {
  dev_overrides {
    "insightfinder/insightfinder" = "/path/to/terraform-provider-insightfinder"
  }
  direct {}
}
```

Set the environment variable:
```bash
export TF_CLI_CONFIG_FILE=.terraformrc
```

## Code Style

- Follow standard Go conventions
- Run `go fmt` before committing
- Use meaningful variable and function names
- Add comments for complex logic
- Keep functions focused and concise

## Documentation

- Update README.md for significant changes
- Add/update examples in the `examples/` directory
- Update resource documentation in `docs/`
- Keep CHANGELOG.md updated

## Commit Message Guidelines

- Use present tense ("Add feature" not "Added feature")
- Use imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit first line to 72 characters
- Reference issues and pull requests

Examples:
```
feat: Add support for new log label types
fix: Correct ServiceNow OAuth authentication
docs: Update JWT configuration examples
chore: Update dependencies
```

## Release Process

Releases are automated via GitHub Actions:

1. Update VERSION file
2. Update CHANGELOG.md
3. Create and push a git tag:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```
4. GitHub Actions will build and publish the release

## Questions?

Feel free to open an issue for questions or reach out to the maintainers.

## License

By contributing, you agree that your contributions will be licensed under the Mozilla Public License 2.0.
