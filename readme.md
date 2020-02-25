# Go-home-brew
*Go-home-brew*, *Gomebrew*, or simply *gome* is a lite homebrew client written in Go. It uses homebrew's API and database for packages but **it's not a clone of it**. It doesn't translate Ruby code to Go code.

## Why ?
I recently used `topgrade` to update my system. It updated my ruby and gems , so homebrew broke. I fixed it by using their portable-ruby solution, but now when I use the *brew prune*, it breaks again. 

> I started this project not because it's faster than Ruby, but because it provides a static binary.

Homebrew's Ruby allows its formulas to be written in Ruby classes, and developers are using a DSL to add installation instructions. Among many things, this is something gomebrew won't support in the near future.

> Gomebrew is not a replacement for homebrew. It doesn't try to be. Homebrew is a complete project with many functionalities. Gomebrew is just a very small subset of it.

Gomebrew is defensive. If it doesn't support a certain functionality, it tries to fail by printing the reason. 



## Commands

### Install

`gomebrew install <program>`

Installs program to `gome_packages` folder and creates a symbolic link of the executable at `/usr/local/bin`. Adds the prefix `gome-` to executable to prevent it mixing with current homebrew programs. E.g. `tree` becomes `gome-tree`.

### List
`gomebrew list`

lists installed programs in `gome_packages` folder.

### Info

`gomebrew info <program>`

prints returned json from API request to homebrew.

### Upgrade

`gomebrew upgrade <program>`

Upgrades program if new version is available. If `<program>` is omitted, upgrades all programs.

### Uninstall

`gomebrew uninstall <program>`

Removes program from the `gome_packages` folder and deletes the symlink.

### Prune
`gomebrew prune`

Removes `gome_packages` folder and deletes symlinks.

## Issues

- Currently only works with standalone programs. If a program has a dependency `gomebrew` will simply exit.

 - Does not support self-build. *Homebrew* has a DSL for installation instructions and `gomebrew` probably won't support this.

 - It works with manpages for simple cases.

 - There is no caching.