package gamelist

import (
	"github.com/beevik/etree"
)

const (
	NameElement        = "name"
	DescElement        = "desc"
	ImageElement       = "image"
	PlayersElement     = "players"
	GenreElement       = "genre"
	PathElement        = "path"
	GameListElement    = "gameList"
	GameElement        = "game"
	FolderElement      = "folder"
	FolderLinkElement  = "folderlink"
	ReleaseDateElement = "releasedate"
	DeveloperElement   = "developer"
	PublisherElement   = "publisher"
	RatingElement      = "rating"
	MD5Element         = "md5"
	VideoElement       = "video"
	MarqueeElement     = "marquee"
	BezelElement       = "bezel"
	ManualElement      = "manual"
	FanartElement      = "fanart"
	BoxbackElement     = "boxback"
	ThumbnailElement   = "thumbnail"
	LangElement        = "lang"
	RegionElement      = "region"
	CheevosHashElement = "cheevosHash"
	CheevosIDElement   = "cheevosId"
	ScraperIDElement   = "scraperId"

	// wrapperElement is used to re-parse ES-DE gamelist files that have
	// sibling elements before <gameList> (e.g. <alternativeEmulator>).
	// Standard XML requires a single root, so we wrap the raw bytes before
	// parsing and unwrap on save.
	wrapperElement = "groutWrapper"
)

type FileName string

const (
	GameListFileName      FileName = "gamelist.xml"
	MiyooGameListFileName FileName = "miyoogamelist.xml"
)

type GameList struct {
	document *etree.Document
}

func New() *GameList {
	return &GameList{
		document: emptyGameList(),
	}
}

func emptyGameList() *etree.Document {
	document := etree.NewDocument()
	document.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	document.CreateElement(GameListElement)
	return document
}

// Parse reads a gamelist XML byte slice. If the bytes are not valid XML
// (e.g. ES-DE files with <alternativeEmulator> before <gameList>), it wraps
// the content in a synthetic root element so that all children remain
// accessible. Save() strips the wrapper back out when writing.
func (gl *GameList) Parse(b []byte) error {
	document := etree.NewDocument()
	if err := document.ReadFromBytes(b); err == nil {
		gl.document = document
		return nil
	}

	wrapped := []byte("<" + wrapperElement + ">" + string(b) + "</" + wrapperElement + ">")
	document = etree.NewDocument()
	if err := document.ReadFromBytes(wrapped); err != nil {
		return err
	}
	gl.document = document
	return nil
}

// root returns the <gameList> element, navigating through the wrapper element
// when the document was parsed from an ES-DE file with sibling elements.
func (gl *GameList) root() *etree.Element {
	if root := gl.document.SelectElement(GameListElement); root != nil {
		return root
	}
	if wrapper := gl.document.SelectElement(wrapperElement); wrapper != nil {
		return wrapper.SelectElement(GameListElement)
	}
	return nil
}

func (gl *GameList) Contains(element, value string) bool {
	root := gl.root()
	if root == nil {
		return false
	}
	games := root.SelectElements(GameElement)
	for _, game := range games {
		element := game.FindElement(element)
		if element != nil && element.Text() == value {
			return true
		}
	}
	return false
}

func (gl *GameList) GetGameElementByName(name string) *etree.Element {
	return gl.getElementByName(GameElement, name)
}

func (gl *GameList) getElementByName(elementName, name string) *etree.Element {
	root := gl.root()
	if root == nil {
		return nil
	}
	for _, entry := range root.SelectElements(elementName) {
		nameElement := entry.FindElement(NameElement)
		if nameElement != nil && nameElement.Text() == name {
			return entry
		}
	}
	return nil
}

func (gl *GameList) GameContainsElements(name string, elements []string) bool {
	e := gl.GetGameElementByName(name)
	if e == nil {
		return false
	}
	for _, element := range elements {
		if e.FindElement(element) == nil {
			return false
		}
	}
	return true
}

// Save writes the gamelist to path. When the document was parsed from an
// ES-DE file using the wrapper strategy, the wrapper element is stripped so
// the original multi-root structure is preserved.
func (gl *GameList) Save(path string) error {
	if wrapper := gl.document.SelectElement(wrapperElement); wrapper != nil {
		out := etree.NewDocument()
		for _, child := range wrapper.Child {
			switch t := child.(type) {
			case *etree.Element:
				out.AddChild(t.Copy())
			case *etree.ProcInst:
				out.CreateProcInst(t.Target, t.Inst)
			case *etree.CharData:
				out.CreateCharData(t.Data)
			case *etree.Comment:
				out.CreateComment(t.Data)
			}
		}
		out.Indent(4)
		return out.WriteToFile(path)
	}
	gl.document.Indent(4)
	return gl.document.WriteToFile(path)
}

func (gl *GameList) AddGameEntry(info map[string]string) {
	gl.addEntry(GameElement, info)
}

func (gl *GameList) addEntry(elementName string, info map[string]string) {
	root := gl.root()
	if root == nil {
		gl.document = emptyGameList()
		root = gl.root()
	}
	entry := root.CreateElement(elementName)

	for key, value := range info {
		entry.CreateElement(key).SetText(value)
	}
}

func (gl *GameList) AdddOrUpdateEntry(name string, info map[string]string) {
	gl.removeEntryByName(FolderElement, name)
	gl.addOrUpdateEntry(GameElement, name, info)
}

func (gl *GameList) AddOrUpdateFolderEntry(name string, info map[string]string) {
	gl.removeEntryByName(GameElement, name)
	gl.addOrUpdateEntry(FolderElement, name, info)
}

func (gl *GameList) removeEntryByName(elementName, name string) {
	entry := gl.getElementByName(elementName, name)
	if entry != nil {
		gl.root().RemoveChild(entry)
	}
}

func (gl *GameList) addOrUpdateEntry(elementName, name string, info map[string]string) {
	entry := gl.getElementByName(elementName, name)
	if entry == nil {
		gl.addEntry(elementName, info)
		return
	}

	for key, value := range info {
		if element := entry.FindElement(key); element != nil {
			element.SetText(value)
		} else {
			entry.CreateElement(key).SetText(value)
		}
	}
}

func (gl *GameList) SetGameID(name, id string) {
	game := gl.GetGameElementByName(name)
	if game == nil {
		return
	}

	game.CreateAttr("id", id)
}
