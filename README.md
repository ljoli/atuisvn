# atuisvn: Terminal user interface for Subversion

> **Fork** of [YoshihideShirai/tuisvn](https://github.com/YoshihideShirai/tuisvn) — original work by [Yoshihide Shirai](mailto:yoshihide.shirai@gmail.com).  
> This fork bumps the Go baseline to 1.22, modernizes the CI workflows and migrates the module path to `github.com/ljoli/atuisvn`.

## What is atuisvn

atuisvn is a Subversion (svn) client for the terminal.  
It functions mainly as a svn repository browser like [tig](https://github.com/jonas/tig).

![atuisvn screenshot](./docs/screenshot.png)

## System requirement

- Linux
- Installed svn command.

## Installation

~~~shell
go install github.com/ljoli/atuisvn@latest
~~~

## Key bindings

- h / ? : Toggle help menu

### main screen

- j / k : Move down / up
- Enter : Open repository
- q : Quit

### tree screen

- j / k : Move down / up
- Enter : Open directory
- l : Open log screen on selected file or directory
- q : Back / Quit

### log screen

- j / k : Move down / up
- Enter : Open revision screen
- q : Back

### rev screen

- j / k : Move down / up
- Enter : Open diff screen on changed path
- c : Open file content at selected revision (svn cat)
- q : Back

### cat screen

- j / k : Move down / up
- gg / G : Go to first / last line
- Ctrl+d / Ctrl+u : Half page down / up
- Ctrl+f / Ctrl+b : Page down / up
- { / } : Previous / next paragraph
- /word : Search forward
- ?word : Search backward
- n / N : Next / previous search result
- :n : Go to line n (example: :10)
- s : Save file locally (edit path, then Enter)
- q : Back

### diff screen

- j / k : Move down / up
- q : Back

## Usage

~~~shell
tuisvn [repository path or url]
~~~

If argv is omitted, repository path is set the current directory.

## Build

Requires [Go 1.24+](https://go.dev/dl/).

~~~shell
git clone https://github.com/ljoli/atuisvn.git
cd atuisvn
go build -o tuisvn .
~~~

The resulting `tuisvn` binary can be moved anywhere in your `$PATH`.

On Windows, building inside WSL is recommended:

~~~shell
export PATH=/usr/local/go/bin:$PATH
GOTOOLCHAIN=local go build -o tuisvn .
~~~

## Development status

Tuisvn is under development.  
Implemented status is following...

- [x] file and directory tree. (svn ls)
- [x] revision history. (svn log)
- [x] revision diff
- [x] revision cat

## License

Copyright (C) 2022 [Yoshihide Shirai](mailto:yoshihide.shirai@gmail.com).  
Fork maintained by [ljoli](https://github.com/ljoli).

Licensed under the [MIT License](LICENSE).
