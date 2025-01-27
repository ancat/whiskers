# Whiskers

Whiskers is a proof of concept tool for detecting malicious code that may be
hiding in Ruby Gems. It currently only analyzes gem upgrades. Individual gems
may be diffed via `gem-diff` and scanned using `gem-diff-scan`. Entire
`Gemfile.lock` files may be diffed and scanned using `gemfile-diff` and
`gemfile-diff-scan`, respectively. These commands take two `Gemfile.lock`
files: a before and an after.

This tool diffs files between upgrades to narrow the scope of files inspected,
and then uses Semgrep to statically analyze these diffs to identify common
malicious payloads. These rules can be found in the `semgrep-rules` directory.
These rules are not suitable for general use since they are likely to have
false positives that can be difficult to programmatically triage and dedupe in
other contexts.

```
$ ./whiskers -h

Usage:
  whiskers [command]

Available Commands:
  completion        Generate the autocompletion script for the specified shell
  gem-diff          Compare two versions of a gem
  gem-diff-scan     Compare two versions of a gem and scan for new issues
  gem-download      Download and extract a Ruby gem
  gemfile-diff      Compare two Gemfile.lock files
  gemfile-diff-scan Load a Gemfile diff and scan changed gems for new issues
  gems              List all gems in a Gemfile.lock
  help              Help about any command

Flags:
  -c, --config string   config file (default is $HOME/.whiskers.yaml)
  -h, --help            help for whiskers

Use "whiskers [command] --help" for more information about a command.
```
