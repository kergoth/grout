# Installation Guide for Steam Deck / RetroDECK

This guide covers installing Grout as a non-Steam game alongside [RetroDECK][retrodeck] on a
Steam Deck. Grout runs independently of RetroDECK's ES-DE instance — it reads RetroDECK's
configuration to resolve ROM/BIOS/save/media paths, but launches through Steam's own game
library rather than as an ES-DE port.

## Tested Devices

| Manufacturer | Device |
| --- | --- |
| _None yet_ | _Please report your results!_ |

## Prerequisites

- RetroDECK installed and configured on your Steam Deck.
- RetroDECK launched at least once so it creates its configuration files.

## Installation Steps

1. Download the [latest Grout release](https://github.com/rommapp/grout/releases/latest/download/Grout-RetroDeck.zip) for RetroDECK.
2. Unzip the archive to a folder in your home directory, e.g. `~/grout/`. It should contain a
   `Grout` folder and a `Grout.sh` file.
3. Open Steam and add Grout as a non-Steam game (**Add a Game** → **Add a Non-Steam Game**, or
   the **+** button in Desktop Mode):
   - Target: `env`
   - Start In: `~/grout/`
   - Launch Options: `~/grout/Grout.sh`
4. Restart Steam for the new entry to appear, then launch Grout from your library.

## Important Configuration

!!! important
    Grout reads paths from `~/.var/app/net.retrodeck.retrodeck/config/retrodeck/retrodeck.json`.
    Launch RetroDECK at least once before running Grout.

!!! warning
    Grout doesn't currently play well with Steam Input. Disable Steam Input for Grout's
    non-Steam game entry (**Controller Settings** → **Disable Steam Input**), or use a
    controller Steam Input doesn't manage, until this is resolved.

## Update

### In-App Update (Recommended)

Launch Grout, open `Settings`, and select `Check for Updates`.

### Manual Update

Download the latest release and replace the `Grout` folder and `Grout.sh` in your install directory. Keep `config.json` to preserve the server configuration.

## Next Steps

See the [User Guide](../usage/guide.md).

--8<-- "docs/_includes/cfw-links.md"
