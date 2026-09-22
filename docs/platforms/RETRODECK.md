# RetroDECK Platform Paths

Grout runs as a non-Steam game on Steam Deck systems configured with [RetroDECK][retrodeck]; see
[Installation Guide for Steam Deck / RetroDECK](../getting-started/install-retrodeck.md). It reads
RetroDECK's own configuration to resolve paths, and writes downloaded games into RetroDECK's ES-DE
gamelist.xml so they're playable from RetroDECK's library, but Grout itself is not launched
through ES-DE.

## Configuration

RetroDECK stores its configuration in `~/.var/app/net.retrodeck.retrodeck/config/retrodeck/retrodeck.json`. Grout reads ROM, BIOS, save, and media paths from this file.

Set `RETRODECK_CFG` to an alternate `retrodeck.json` path when using a nonstandard installation or testing a separate configuration.

## Paths

| Resource | Path |
| --- | --- |
| Gamelists | `~/.var/app/net.retrodeck.retrodeck/config/ES-DE/gamelists/<system>/gamelist.xml` |
| Media | `downloaded_media_path` in `retrodeck.json` |
| Save root | `saves_path` in `retrodeck.json` |

## Save Sync

Save sync is best effort. Grout maps RetroArch-backed emulators to `retroarch/saves`, but standalone emulators can use different locations.

--8<-- "docs/_includes/cfw-links.md"
