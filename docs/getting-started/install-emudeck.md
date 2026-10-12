# Installation Guide for EmuDeck

This guide will help you install Grout on a Steam Deck running [EmuDeck][emudeck]. Grout runs as an ES-DE port, from
the `Ports` system.

## Tested Devices

| Manufacturer | Device     | OS      |
|--------------|------------|---------|
| Valve        | Steam Deck | SteamOS |

_Please help verify compatibility on other devices by reporting your results!_

## Prerequisites

- EmuDeck and ES-DE installed, and ES-DE launched at least once.
- The standard EmuDeck `Emulation` folder layout for ROMs, BIOS files, and saves. Grout reads where EmuDeck keeps them
  from `~/.config/EmuDeck/settings.json`, and where ES-DE keeps downloaded media from
  `~/ES-DE/settings/es_settings.xml`.

## Installation Steps

1. Download the [latest Grout release](https://github.com/rommapp/grout/releases/latest/download/Grout-EmuDeck.zip) for
   EmuDeck.
2. Unzip the archive.
3. Copy the `Grout` folder and `Grout.sh` into EmuDeck's `ports` folder, `<Emulation folder>/roms/ports/`. On an
   SD card this is usually `/run/media/<sd-card-name>/Emulation/roms/ports/`, otherwise `~/Emulation/roms/ports/`.
4. Restart ES-DE, or reload its game list.
5. Start Grout from the `Ports` system.

!!! note
    Grout adds itself to the `Ports` game list the first time it runs.

## Controls

Grout uses the Deck's buttons as labelled, so `A` confirms and `B` goes back. Navigate with the D-pad; the sticks and
triggers do not move the cursor. On a game list, click the left stick (`L3`) to open the BIOS screen.

## Update

### In-App update (Recommended)

Grout has a built-in update mechanism. To update Grout, launch the application and navigate to the `Settings` menu. From
there, select `Check for Updates`. If a new version is available, follow the on-screen prompts to download and install
the update.

### Manual update

Download the latest release and replace the `Grout` folder and `Grout.sh` in the ports folder. Keep `config.json` in the
`Grout` folder to preserve your server settings.

## Next Steps

See the [User Guide](../usage/guide.md).

--8<-- "docs/_includes/cfw-links.md"
