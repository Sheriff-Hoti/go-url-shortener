# Go URL Shortener

A simple URL shortener web application built with Go, [templ](https://templ.guide/), HTMX, Tailwind CSS, and SQLite.

## Features

- Shorten URLs with a unique code
- Instant feedback using HTMX (no page reloads)
- In-memory SQLite database (easy to switch to file-based)
- Hot reload for Go code and Tailwind CSS
- Clean UI with Tailwind CSS

## Tech Stack

- **Go** (1.24+)
- **[templ](https://templ.guide/)** for type-safe HTML components
- **HTMX** for dynamic frontend interactions
- **Tailwind CSS** for styling
- **SQLite** (via [glebarez/go-sqlite](https://github.com/glebarez/go-sqlite))
- **sqlc** for type-safe database access

## Getting Started

### Prerequisites

- Go installed (system-wide)
- [Nix](https://nixos.org/download.html) (optional, for reproducible dev environment)
- [Air](https://github.com/cosmtrek/air) for hot reload (optional)

### Development Setup

#### With Nix

```sh
nix-shell
```

#### Without Nix

Install dependencies:
- [templ](https://templ.guide/docs/getting-started/)
- [tailwindcss](https://tailwindcss.com/docs/installation)
- [air](https://github.com/cosmtrek/air)
- [sqlc](https://docs.sqlc.dev/en/latest/overview/install.html)

#### Generate templ components

```sh
make templ
```

#### Start the development server

```sh
make dev
```

This runs:
- Go server with hot reload (`air`)
- Tailwind CSS watcher
- templ watcher

Visit [http://localhost:3000](http://localhost:3000) in your browser.

## Usage

1. Enter a URL in the input field.
2. Click "Shorten".
3. The shortened URL will appear below the form.

## Project Structure

```
assets/         # Static assets (CSS, JS)
database/       # Database models and queries (sqlc)
templates/      # templ components
main.go         # Application entry point
Makefile        # Development commands
shell.nix       # Nix development environment
```