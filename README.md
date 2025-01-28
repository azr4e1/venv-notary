# venv-notary

venv-notary is a simple utility to easily manage your Python environment needs. It follows the UNIX phylosophy of:

> do one thing and do it well

venv-notary doesn't try to manage packages or do dependency resolution. It only creates and manages python environment, so that you don't need to worry about where to create them and how to activate them, or even remembering where the hell you created that environment you need now.

## Requirements

- bash (more shells coming soon)
- go version 1.23

## Install

```bash
go install github.com/azr4e1/venv-notary/vn@latest
```

Make sure your `$GOPATH` is in your `$PATH`

Once installed, run `vn completion <YOUR SHELL>`, and append the result into your shell config in order to enable completions for venv-notary.

Shells supported for completion:

- bash
- zsh
- fish
- powershell

## Usage

venv-notary distinguishes between two types of environments: local and global.

**Global** environments are environments that you would call from anywhere in the filesystem. They are meant to contain packages that you know you would need for various tasks, and not necessarily just for a single project. For example, you could have a `data-science` global environment for whenever you need to perform some eda on the fly, or a `llm` environment when you want to just create a small script calling openai or gemini APIs.

**Local** environments are environments specific to a folder. They are meant to be used as project specific environments, and are similar in concept to poetry's shell.

### Create a new environment

Create a new global environment:

```bash
vn create data-science
```


Create a new local environment:

```bash
vn create -l
````


Create a new environment with a specific Python version:

```bash
vn create python39-venv -p python3.9
```

### Activate an environment

Activate a global environment:

```bash
vn activate data-science
```

Activate a local environment:

```bash
vn activate -l
```

Activate an environment with specific version:

```bash
vn activate -l -p python3.9
```

**Note**: if the environment you want to activate hasn't been created before, it will automatically be created.

### Delete an environment

Delete a global environment:

```bash
vn delete data-science
```

Delete a local environment:

```bash
vn delete -l
```

Delete an environment with specific version:

```bash
vn delete -l -p python3.9
```

### Clean local/global environments

`clean` is like `delete` on steroid. It allows to delete environments in batches.

You must specify whether to clean the global or the local environments. **If no other flag is provided, all your local/global environments will be deleted**

You can also provide the `-p/--python` flag to delete only the environments with a specific Python version.

You can also provide the `-n/--name` flag to delete only the environments whose name match a regexp pattern.

Examples:

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
