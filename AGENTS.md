# AGENTS.md file for smolgit
 
## Dev environment tips
- Use `make` to inspect available targets for project, like build, lint, run, test.
 
## Testing instructions
- Find the CI plan in the .github/workflows folder.
- From the package root you can just call `make test`. The commit should pass all tests before you merge.
- Fix any test or type errors until the whole suite is green.
- After moving files or changing imports, run `make lint` to be sure golangci-lint linters still pass.
- Add or update tests for the code you change, even if nobody asked.
 
## PR instructions
- Title format: [<project_name>] <Title>
- Always run `make lint` and `make test` before committing.
