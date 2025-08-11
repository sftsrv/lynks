# Lynks

A CLI tool for interactively fixing links in markdown files

## Features

- Interactively view list of markdown files in a repository and links between them
- Support for markdown style links
- Basic configuration of link aliases
- Basic linting for links

## Installation

The project is still in the early stages of development, at the moment it's only possible to install using `go get`

### Using `go install`

```sh
go install github.com/sftsrv/lynks@latest
```

## Usage

### Config

Create a `lynks.config.json` file from the directory you want to run the command, It should have the following structure:

```json
{
  // root folder from which pages should be resolved
  "root": "./src/docs",
  "aliases": {
    // aliases resolve relative to the "root"
    // the key can be any value that you use within pages for linking
    "@api": "./generated/api"
  }
}
```

### Running

There are two ways to run the tool:

#### Interactive

You can run the tool interactively in order to browser files and fix links in all markdown files within the `root` as defined in the `lynks.config.json` by simply running `lynks`

```sh
lynks
```

#### Linter

The tool can also be run as a linter which will make use of the `lynks.config.json` and can be run using:

```sh
lynks lint
```

## Project Roadmap

Some things that I still want to do before considering this project complete:

- [ ] Add tests for like everything
- [ ] Link management strategies
- [ ] Flags for more specific behavior like:
  - Interactive "fix" mode
  - Better control of linting
    - Only show files with errors
    - Only show links with errors
- [ ] Help, informative errors, etc.
- [ ] Management of image and mdx links
- [ ] Ability to output relative links
- [ ] Support for index pages
- [ ] Imporove overall UX
- [ ] Support links with hashes
