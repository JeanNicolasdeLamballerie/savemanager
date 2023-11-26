package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

const VERSION = "0.1.0"

type GlobalConfiguration struct {
	Version              string `toml:"version" comment:"The current version of the installation. Used internally to check for older data !"`
	Verbose              bool   `toml:"verbose" comment:"Activating this option will log various information to the console running the save manager. Use this if you encounter issues and wish to debug (or get assistance in debugging !) a problem you're encountering."`
	ExperimentalFeatures bool   `toml:"experimentalFeatures" comment:"Activating this option will turn on experimental features. Use at your own risk, they might not be stable :)"`
	SavesLocation        string `toml:"savesLocation" comment:"The place where your saves & backups are stored. Do NOT change this manually unless you have no profiles or saves at all. Use the 'ManualMigrate' value below instead to change where your saved are backed up."`
	ManualMigrate        string `toml:"manualMigrate" comment:"[NOT IMPLEMENTED YET] Editing this value to any existing path will move all your savefile data to that location, in a directory named 'profiles', then delete them from the current directory. Do NOT edit SavesLocation manually unless you know what you're doing."`
	SeedAction           bool   `toml:"seedAction" comment:"Automatically attempts to create basic gametypes and tags. Only useful on a new database."`
	initialized          bool

	// TODO : CHECK for seed/Seed action history ? Somehow ?
}

//	os.UserConfigDir() TODO this returns appdata
//
// TODO place all files/db in appdata ?
func GetProjectDataLocation() string {

	var appDataLoc, _ = os.UserConfigDir()
	var roamingLoc = filepath.Join(appDataLoc)
	return filepath.Join(roamingLoc, "SaveManager_by_Dekharen")

}

var projectDataLoc = GetProjectDataLocation()

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	// TODO check for permission errors ?
	return false, err
}

var configPath = filepath.Join(projectDataLoc, "savemanager.config.TOML")

var CONFIGURATION = GlobalConfiguration{initialized: false}

func (config *GlobalConfiguration) Init() {
	if config.initialized {
		println("Initialization was already called and this call does not achieve anything. This might be a coding mistake ?")
		return
	}

	var hasProjectData, error = exists(projectDataLoc)
	if !hasProjectData && error == nil {
		os.Mkdir(projectDataLoc, os.ModeDir)
	}
	file, err := os.ReadFile(configPath)
	if err != nil {
		//TODO check for permission rights ?
		config.create()
	} else {

		error := toml.Unmarshal(file, &config)
		if error != nil {

			println("The toml config file seems to be corrupted or inaccessible. Contact the developer")
			log.Fatalln(error.Error())
		}
	}
	config.initialized = true
	//TODO specify updates depending on differences ?
	if config.Version != VERSION {
		log.Fatalln("AUTOMATIC UPDATE DISABLED FOR NOW. The config version and build version are not the same !")
	}
}

func (config *GlobalConfiguration) create() {
	config.Version = VERSION
	config.Verbose = true
	config.SeedAction = true
	config.ExperimentalFeatures = false
	config.SavesLocation = filepath.Join(projectDataLoc, "profiles")
	config.ManualMigrate = ""
	println("First setup ! We're setting up in %appdata% to store some files.")
	println("configuration file location :")
	println(configPath)

	config.Save()
}
func GetConfig() *GlobalConfiguration {
	if !CONFIGURATION.initialized {

		CONFIGURATION.Init()
	}
	return &CONFIGURATION
}

func (config *GlobalConfiguration) Save() {
	configFile, err := os.Create(configPath)
	if err != nil {
		log.Fatalln(err.Error())
	}
	b, err := toml.Marshal(config)
	if err != nil {
		println(err.Error())
	}
	configFile.Write(b)
}
