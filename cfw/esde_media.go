package cfw

import "path/filepath"

// esdeFolder names the folder ES-DE reads each kind of media from, under
// downloaded_media/<system>/. Kinds ES-DE has no folder for, such as bezels,
// are absent.
var esdeFolder = map[ArtSlot]string{
	ArtCover:   "covers",
	ArtMarquee: "marquees",
	ArtBoxback: "backcovers",
	ArtFanart:  "fanart",
	ArtVideo:   "videos",
	ArtManual:  "manuals",
}

// esdeMedia builds the media directory resolver for a firmware running ES-DE.
// ES-DE finds media by folder and file name alone, so every kind gets its own
// folder and a file takes the rom's name with no suffix. system names the ES-DE
// system a rom directory belongs to.
func esdeMedia(root func() string, system func(romDir string) string) func(slot ArtSlot, romDir string) string {
	return func(slot ArtSlot, romDir string) string {
		folder, ok := esdeFolder[slot]
		if !ok {
			return ""
		}
		return filepath.Join(root(), system(romDir), folder)
	}
}
