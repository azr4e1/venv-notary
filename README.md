venv-notary is a simple utility to easily manage your Python environment needs. It follows the UNIX phylosophy of:

> do one thing and do it well

venv-notary doesn't try to manage packages or do dependency resolution. It only creates and manages python environments, so that you don't need to worry about where to create them and how to activate them, or even remembering where the hell you created that environment you need now.

## Requirements

- bash, zsh, fish or powershell
- go version 1.23

## Install

```bash
go install github.com/azr4e1/venv-notary/vn@latest
```

Make sure your `$GOPATH` is in your `$PATH`.

Once installed, run `vn completion <YOUR SHELL>`, and append the result into your shell config in order to enable completions for venv-notary.

Shells supported for completion:

- bash
- zsh
- fish
- powershell

## Usage

venv-notary distinguishes between two types of environments: local and global.

**Global environments** are environments that you would call from anywhere in the filesystem. They are meant to contain packages that you know you would need for various tasks, and not necessarily just for a single project. For example, you could have a `data-science` global environment for whenever you need to perform some eda on the fly, or a `llm` environment for whenever you want to just create a small script calling openai or gemini APIs.

**Local environments** are environments specific to a folder. They are meant to be used as project specific environments, and are similar in nature to poetry's `shell`.

Remember to call `vn help` on any command if you're stuck:

```bash
vn help clean
```

### Create a new environment

Create a local environment (default):

```bash
vn create
```

Create a global environment:

```bash
vn create -g data-science
```

Create an environment with a specific Python version:

```bash
vn create -g python39-venv -p python3.9
```

### Activate an environment

Activate the local environment (default):

```bash
vn activate
```

Activate a global environment:

```bash
vn activate -g data-science
```

Activate with a specific Python version:

```bash
vn activate -p python3.9
```

**Note**: if the environment you want to activate doesn't exist, it will automatically be created.

### Delete an environment

Delete the local environment (default):

```bash
vn delete
```

Delete a global environment:

```bash
vn delete -g data-science
```

Delete the local environment with a specific Python version:

```bash
vn delete -p python3.9
```

### Clean local/global environments

`clean` is like `delete` on steroid. It allows to delete environments in batches.

You must specify whether to clean the global or the local environments. **If no other flag is provided, all your local/global environments will be deleted**!

You can also provide the `-p/--python` flag to delete only the environments with a specific Python version.

You can also provide the `-n/--name` flag to delete only the environments whose names match a regexp pattern.

#### Examples

Delete all global environments:

```bash
vn clean -g
```

Delete all local environments with this Python version:

```bash
vn clean -l -p python3.9
```
The argument to `-p` must be a Python executable.


Delete all local environments that match this pattern:

```bash
vn clean -l -n "data.*$"
```

### Run a command in an environment

Run a command in the local environment (default):

```bash
vn run pytest
```

Run a command in a global environment:

```bash
vn run -g data-science jupyter notebook
```

Run a command with a specific Python version:

```bash
vn run -p python3.9 python script.py
```

Arguments after the command name are passed through to the command. Use `--` to separate flags from the command if needed:

```bash
vn run -- pytest --tb=short
```

### List

Finally, you can list your local/global environments, optionally filtering by Python version.

List in interactive mode:

```bash
vn list
```

List only local environments:

```bash
vn list -l
```

List only global environment with Python version 3.12:

```bash
vn list -g -p python3.12
```

Output in JSON format:

```bash
vn list -j
```

The `-j` flag can be combined with any other `list` flags:

```bash
vn list -l -j
```
