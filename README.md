# Lynks

A CLI tool for interactively fixing links in markdown files

## Features

- Interactively view list of markdown files in a repository and links between them
- Support for markdown style links
- Basic configuration of link aliases

## Installation

The project is still in the early stages of development, at the moment it's only possible to install using `go get`

### Using `go install`

```sh
go install github.com/sftsrv/lynks
```

## Usage

Create a `lynks.config.json` file at your project root with the following structure:

```json
{
  // root folder from which pages should be resolved
  "root": "./src/docs",
  "aliases": {
    // aliases resolve relative to the "root"
    // the key can be any value that you use within pages for linking
    "@api": "./generated/api",
  }, 
}
```

## TODO

- [ ] Help, informative errors, etc.
- [ ] Link checking/linting via CLI
- [ ] Management of image and mdx links
- [ ] Ability to output relative links
- [ ] Support for index pages
- [ ] Imporove overall UX
