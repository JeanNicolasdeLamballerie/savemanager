package filemanager

import (
	"encoding/json"
	"io/fs"
	"os"
	"savemanager/config"
	"slices"
	"strings"
)

//const VERBOSE = true

// Assigns a filename (e.g ER0000) to a game name (e.g Elden Ring).
// KEYS => Files, VALUES => Games
var SAVEFILES_AND_GAMES_NAMES = map[string]string{
	"ER0000":      "Elden Ring",                             // %appdata% bs
	"draks0005":   "Dark Souls PtDE",                        // C:\Users\%user\Documents\NBGI  => IN FOLDER DARKSOULS >(number?)
	"DRAKS0005":   "Dark Souls Remastered",                  // C:\Users\%user\Documents\NBGI => IN FOLDER DARK SOULS REMASTERED > number
	"DS2SOFS0000": "Dark Souls II Scholar of the First Sin", // %appdata% DarkSoulsII/number
	"DARKSII0000": "Dark Souls II",                          //idem
	"DS30000":     "Dark Souls III",                         // %appdata% DarkSoulsIII/number (and letters?)
	"S0000":       "Sekiro: Shadows Die Twice",
}

type SaveDirectoryInformations struct {
	Information map[string]string
	Tags        []string
}

func CheckGameType(path string) {
	saveDirectoryFS := os.DirFS(path)
	information := SaveDirectoryInformations{
		Information: map[string]string{},
		Tags:        []string{},
	}
	//information.Information["gameType"] = "custom"
	err := fs.WalkDir(saveDirectoryFS, ".", checkFilesInDirectory(&information))
	if err != nil {
		println(err.Error())
		return
	}

	slices.Sort[[]string](information.Tags)
	information.Tags = slices.Compact[[]string](information.Tags)
	testval, err := json.Marshal(information)
	if err != nil {
		println("ERROR MARSHAL")
		println(err.Error())
		return
	}

	//TODO check if gameType exists (or implement a "isEmpty on SaveDirectory ?" in case user opens an empty directory)
	println(string(testval))
}

func checkFilesInDirectory(information *SaveDirectoryInformations) func(path string, d fs.DirEntry, err error) error {
	return func(path string, d fs.DirEntry, err error) error {

		if err != nil {
			return err
		}
		if config.GetConfig().Verbose {

			info, _ := d.Info()
			println("NAME : ", d.Name())
			println("INFO : ", info.Mode(), info.Name())
			println("PATH : ", path)
		}
		filename := strings.Split(path, ".")[0]
		gameType, isKnown := SAVEFILES_AND_GAMES_NAMES[filename]
		if !isKnown {
			gameType = "custom"
		}
		information.Information["gameType"] = gameType
		switch gameType {

		// Check if game is ER :
		case "Elden Ring":

			// if so, check for savefile-specific mods to append in the tags :
			if strings.Contains(path, ".co2") {
				information.Tags = append(information.Tags, "seamless")
			}
			if strings.Contains(path, ".mod") {
				information.Tags = append(information.Tags, "convergence")
			}

		}
		if strings.Contains(path, ".sl2") {
			information.Tags = append(information.Tags, "retail")
		}
		return nil
	}

}
