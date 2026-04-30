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

### tree screen

- k : Move-up
- j : Move-down
- Enter : change directory
- l : Enter log screen on selected file and directory.
- q : Back previous screen.

### log screen

- k : Move-up
- j : Move-down
- Enter : Enter rev screen
- q : Back previous screen.

### rev screen

- k : Move-up
- j : Move-down
- Enter : Enter diff screen on change directory.
- q : Back previous screen.

### diff screen

- k : Move-up
- j : Move-down
- q : Back previous screen.

## Usage

~~~shell
tuisvn [repository path or url]
~~~

If argv is omitted, repository path is set the current directory.

## Development status

Tuisvn is under development.  
Implemented status is following...

- [x] file and directory tree. (svn ls)
- [x] revision history. (svn log)
- [x] revision diff
- [ ] revision cat

## License

Copyright (C) 2022 [Yoshihide Shirai](mailto:yoshihide.shirai@gmail.com).  
Fork maintained by [ljoli](https://github.com/ljoli).

Licensed under the [MIT License](LICENSE).
