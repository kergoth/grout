# EmuDeck Platform Paths

Grout runs as an ES-DE port on Steam Deck systems configured with [EmuDeck][emudeck].

## Configuration

Grout reads `storagePath` from `~/.config/EmuDeck/settings.json`. The path selects the EmuDeck storage root, which may be on internal storage or an SD card.

## Paths

| Resource | Path |
| --- | --- |
| Gamelists | `~/ES-DE/gamelists/<system>/gamelist.xml` |
| Media | `MediaDirectory` in `~/ES-DE/settings/es_settings.xml` |
| Fallback media | `~/Emulation/tools/downloaded_media` |
| Save root | `<storagePath>/Emulation/saves` |

## Save Sync

Save sync is best effort. Grout maps RetroArch-backed emulators to `retroarch/saves`, but standalone emulators can use different locations.

## Known Limitations

- Standalone emulator save paths need testing on hardware.
- Grout uses ES-DE's configured media directory when present; otherwise it uses the EmuDeck fallback path.

--8<-- "docs/_includes/cfw-links.md"
