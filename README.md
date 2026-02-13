# Installation

```
gh extension install rawnly/gh-targetprocess
```

# Configure

The first time you run the extension you will be prompted for 2 properties:

- **BASE URL**
  You can insert your targetprocess url (eg `<your-company>.tpondemand.com`)
- **TOKEN**
  You can generate an access token on your profile

# Usage

Just run the extension with the optional argument of the ID.

If your branch has the ID in the name, it'll be derived.

```
gh-targetprocess is a tool to create PRs starting from a Targetprocess ID or branch

Usage:
  gh-targetprocess [flags]
  gh-targetprocess [command]

Aliases:
  gh-targetprocess, gh-tp

Examples:
gh targetprocess 12345

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  configure   Configure the gh-targetprocess CLI
  help        Help about any command
  update      Update current PR with TargetProcess data
  view        View the current ticket

Flags:
  -a, --assign string   assign PR
  -d, --draft           mark pr as draft
      --dry-run         dry-run pr creation
  -h, --help            help for gh-targetprocess
  -l, --label string    label to add to the PR
      --no-body         skip body
  -w, --web             open pr in web browser

Use "gh-targetprocess [command] --help" for more information about a command.
```
