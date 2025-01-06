# Slopr

Slopr is a CLI specialised to interact with the [slop.sh](https://api.slop.sh) API, written in Go to prioritise efficiency and stability so your uploads happen when you want them, not when JavaScript decides you can have them.

# Installation

Currently, there are three ways to install slopr

1. Install the pre-built binary from the GitHub release
2. Build from source
3. Install from the AUR on arch-based distros

### Installing the pre-built binary

This one is probably the easiest solution if you aren't using an Arch-based distro (or Arch Linux itself like a sigma)

1. Download the tarball from the GitHub release
2. Place the binary in a directory on your PATH (e.g. $HOME/.local/bin/)
3. That's it, slopr is now installed

### Building from source

This is definitely the more difficult option, as you require basic knowledge of the terminal and Go itself.

1. Install Go if you haven't already
2. Clone the repository
3. Run `go build` inside the cloned repo
4. Copy the built binary to a directory on your path
5. Slopr should now be installed on your system

### Installing from the AUR (Sigmas only)

All you need is an AUR helper (paru or yay) or you can directly download it from the AUR and run `makepkg -si`

```
paru -S slopr
```

# Usage

It's genuinely the most simple piece of software I've ever used.

```
slopr /path/to/a/file
```

It's that simple. If you have a supported clipboard installed, the shareable link will be automatically copied to your clipboard. If not, you can manually copy it and share it with whoever you'd like.

Not that I think you'll need it, but you can always use the `--help` flag to view the usage, which will always be up to date.
