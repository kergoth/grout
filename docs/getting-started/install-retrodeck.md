# Installation Guide for Steam Deck / RetroDECK

This guide covers installing Grout as an ES-DE port on a Steam Deck using [RetroDECK][retrodeck].

## Tested Devices

| Manufacturer | Device |
| --- | --- |
| _None yet_ | _Please report your results!_ |

## Prerequisites

- RetroDECK installed and configured on your Steam Deck.
- RetroDECK launched at least once so it creates its configuration files.

## Installation Steps

1. Download the [latest Grout release](https://github.com/rommapp/grout/releases/latest/download/Grout-RetroDeck.zip) for RetroDECK.
2. Unzip the archive.
3. Copy the `Grout` folder and `Grout.sh` to the RetroDECK desktop directory. The default location is `~/retrodeck/roms/desktop/`.
   If RetroDECK is installed on an SD card, check `retrodeck.json` for `roms_path`.
4. Restart ES-DE inside RetroDECK or refresh its game list.
5. Start Grout from the `Desktop Applications` system.

!!! important
    Don't place Grout in the `Ports` directory. RetroDECK disables that system's
    "Shortcut or script" and "AppImage" launch commands, leaving only RetroArch
    source-port cores (e.g. ECWolf) as the default emulator, so a plain shell
    script placed there won't run.

## Important Configuration

!!! important
    Grout reads paths from `~/.var/app/net.retrodeck.retrodeck/config/retrodeck/retrodeck.json`.
    Launch RetroDECK at least once before running Grout.

## Update

### In-App Update (Recommended)

Launch Grout, open `Settings`, and select `Check for Updates`.

### Manual Update

Download the latest release and replace the `Grout` folder and `Grout.sh` in the desktop directory. Keep `config.json` to preserve the server configuration.

## Next Steps

See the [User Guide](../usage/guide.md).

--8<-- "docs/_includes/cfw-links.md"
