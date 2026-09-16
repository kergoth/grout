# Installation Guide for Steam Deck / EmuDeck

This guide covers installing Grout as an ES-DE port on a Steam Deck using [EmuDeck][emudeck].

## Tested Devices

| Manufacturer | Device |
| --- | --- |
| _None yet_ | _Please report your results!_ |

## Prerequisites

- EmuDeck and ES-DE installed on your Steam Deck.
- ES-DE launched at least once.

## Installation Steps

1. Download the [latest Grout release](https://github.com/rommapp/grout/releases/latest/download/Grout-EmuDeck.zip) for EmuDeck.
2. Unzip the archive.
3. Copy the `Grout` folder and `Grout.sh` to the EmuDeck ports directory. The directory is `<storagePath>/Emulation/roms/ports/`, where `storagePath` comes from `~/.config/EmuDeck/settings.json`.
   Common locations are `~/Emulation/roms/ports/` and `/run/media/<sd-card-name>/Emulation/roms/ports/`.
4. Restart ES-DE or refresh its game list.
5. Start Grout from the `Ports` system.

## Important Configuration

!!! important
    Grout reads EmuDeck's storage path from `~/.config/EmuDeck/settings.json`.
    Keep the standard `Emulation` directory layout for ROMs, BIOS files, and saves.

!!! note
    Save sync is best effort because standalone emulator save paths differ from RetroArch paths.

## Update

### In-App Update (Recommended)

Launch Grout, open `Settings`, and select `Check for Updates`.

### Manual Update

Download the latest release and replace the `Grout` folder and `Grout.sh` in the ports directory. Keep `config.json` to preserve the server configuration.

## Next Steps

See the [User Guide](../usage/guide.md).

--8<-- "docs/_includes/cfw-links.md"
