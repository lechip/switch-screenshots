# Nintendo Switch Screenshot Organizer

Recursively copies `.jpg` and `.mp4` files from a Nintendo Switch SD card
Album folder and sorts them into per-game subfolders named after the game
title.

Runs on **macOS**, **Linux** (including Arch), and **Windows** — distributed
as a single static binary with no runtime dependencies.

## Installation

### Any platform via Go

```bash
go install github.com/lechip/switch-screenshots/cmd/switch-screenshots@latest
```

### Arch Linux

```bash
cd packaging/arch
makepkg -si
```

### macOS (Homebrew)

```bash
brew install lechip/tap/switch-screenshots
```

### Windows

Download the `.zip` from the
[latest GitHub Release](https://github.com/lechip/switch-screenshots/releases/latest),
extract it, and run `switch-screenshots.exe` from PowerShell or Command
Prompt. To use it from any directory, add the extracted folder to your `PATH`.

## Usage

```
switch-screenshots -i <Album folder> [-o <output folder>] [flags]
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--input` | `-i` | *(required)* | Path to the Switch SD card `Album` folder |
| `--output` | `-o` | `./screenshotsOutput` | Destination folder |
| `--dry-run` | | `false` | Print planned operations without copying |
| `--by-id` | | `false` | Use the raw 32-char game ID instead of the title |
| `--verbose` | `-v` | `false` | Log each file operation |
| `--version` | | | Print version information and exit |

Exit codes: `0` success · `1` argument error · `2` IO error.

### Examples

**macOS**
```bash
switch-screenshots -i /Volumes/SD/Nintendo/Album -o ~/Pictures/Switch
```

**Linux (Arch)**
```bash
switch-screenshots -i /run/media/user/SD/Nintendo/Album -o ~/Pictures/Switch
```

**Windows (PowerShell)**
```powershell
.\switch-screenshots.exe -i D:\Nintendo\Album -o C:\Users\me\Pictures\Switch
```

**Dry run — preview without copying**
```bash
switch-screenshots -i /Volumes/SD/Nintendo/Album --dry-run
```

## How it works

Nintendo Switch screenshots and videos are named like:

```
2018101310443100-EBDA73B0F0F1E8C57869EB26D8F65095.jpg
```

The 32-character hex segment after the `-` is a game-specific identifier.
This tool looks it up in a bundled copy of
[RenanGreca's `game_ids.json`](https://github.com/RenanGreca/Switch-Screenshots-Manager)
and copies each file into a subfolder named after the game:

```
screenshotsOutput/
├── The Legend of Zelda Breath of the Wild/
│   └── 2018101310443100-<id>.jpg
├── Super Mario Odyssey/
│   └── ...
└── EBDA73B0F0F1E8C57869EB26D8F65095/   ← fallback for unrecognised IDs
    └── ...
```

If an ID is not in the bundled map, the raw hex ID is used as the folder
name and a warning is printed at the end of the run suggesting
`make update-gameids`.

## Updating the game ID list

To refresh the bundled `game_ids.json` with the latest titles from
RenanGreca's repository:

```bash
make update-gameids
```

Then commit the updated file and open a PR.

## Building from source

```bash
git clone https://github.com/lechip/switch-screenshots.git
cd switch-screenshots
go build ./cmd/switch-screenshots
```

Cross-compile for a specific target:

```bash
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build ./cmd/switch-screenshots
```

## Acknowledgements

Game ID to title mapping provided by
[RenanGreca/Switch-Screenshots-Manager](https://github.com/RenanGreca/Switch-Screenshots-Manager)
(MIT licence). Many thanks for maintaining the list.
