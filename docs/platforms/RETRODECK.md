# RetroDECK Platform Paths

Grout runs as an ES-DE port on Steam Deck systems configured with [RetroDECK][retrodeck].

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
