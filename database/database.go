package database

import (
	"log"
	"path/filepath"
	"savemanager/config"
	"savemanager/filemanager"

	// "github.com/pelletier/go-toml/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DataBase struct {
	Name string
	DB   *gorm.DB
}

func (this *DataBase) Connect() error {
	DB, err := gorm.Open(sqlite.Open(filepath.Join(config.GetProjectDataLocation(), "main.db")), &gorm.Config{})
	if err != nil {
		return err
	}

	err = DB.AutoMigrate(&Profile{}, &SaveDirectory{}, &InfoTag{}, &GameType{})
	if err != nil {
		log.Fatal(err)
	}
	this.DB = DB
	configuration := config.GetConfig()
	println("Savemanager's current version :", configuration.Version)
	if configuration.SeedAction {
		this.AutomaticSeeding()
		configuration.SeedAction = false
		configuration.Save()

	}
	return nil
}
func (this *DataBase) AutomaticSeeding() error {
	gameTypes := []*GameType{}

	for fileName, gameName := range filemanager.SAVEFILES_AND_GAMES_NAMES {
		convergence := InfoTag{Name: "convergence", Color: "red"}
		seamless := InfoTag{Name: "seamless", Color: "blue"}
		retail := InfoTag{Name: "retail", Color: "black"}
		infoTags := []*InfoTag{&retail}
		if gameName == "Elden Ring" {
			infoTags = append(infoTags, &convergence, &seamless)
		}

		gameTypes = append(gameTypes, &GameType{Name: gameName, Filename: fileName, InfoTags: infoTags})
	}

	err := this.DB.Create(&gameTypes).Error
	if err != nil {
		println("Seeding gametypes failed. See error :")
		println(gameTypes)
	}
	return err
}

type SaveDirectory struct {
	gorm.Model
	Name      string `gorm:"unique"`
	Path      string `gorm:"notNull"`
	ProfileID uint
	//	Profile   *Profile `gorm:"foreignKey:ProfileID"`

	InfoTags []*InfoTag `gorm:"many2many:infoTag_SaveDirectories;"`
}

type InfoTag struct {
	gorm.Model
	Name   string
	Color  string
	TypeID uint
	//Type            *GameType        `gorm:"foreignKey:TypeID"`
	SaveDirectories []*SaveDirectory `gorm:"many2many:infoTag_SaveDirectories;"`
}

type Profile struct {
	gorm.Model
	ProfileName     string `gorm:"unique"`
	GamePath        string
	TypeID          uint
	Type            *GameType        `gorm:"foreignKey:TypeID"`
	SaveDirectories []*SaveDirectory `gorm:"foreignKey:ProfileID"`
}
type GameType struct {
	gorm.Model
	Name     string `gorm:"unique"`
	Filename string
	InfoTags []*InfoTag `gorm:"foreignKey:TypeID"`
}

// INSERT
func (this DataBase) CreateProfile(profile *Profile) (*uint, error) {

	result := this.DB.Create(profile)
	if result.Error != nil {
		return nil, result.Error
	}
	//print(profile, &profile)
	return &profile.ID, nil
}
func (this DataBase) CreateSaveDirectory(save *SaveDirectory) (*uint, error) {
	result := this.DB.Create(save)
	if result.Error != nil {
		return nil, result.Error
	}
	print(save, &save)
	return &save.ID, nil
}

// GET
func (this DataBase) GetProfileByID(profileID uint) (*Profile, error) {
	var profile Profile
	result := this.DB.First(&profile, profileID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &profile, nil
}
func (this DataBase) GetAllProfiles() (*[]*Profile, error) {
	var profiles []*Profile
	result := this.DB.Find(&profiles)
	if result.Error != nil {
		return nil, result.Error
	}
	return &profiles, nil
}

func (this DataBase) GetSaveDirectoryByID(saveDirectoryID uint) (*SaveDirectory, error) {
	var saveDirectory SaveDirectory
	result := this.DB.First(&saveDirectory, saveDirectoryID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &saveDirectory, nil
}

func (this DataBase) GetSavesFromProfile(profile *Profile) []SaveDirectory {
	var saves []SaveDirectory
	// err := this.DB.Model(profile).Association("SaveDirectories").Find(&saves)

	this.DB.Model(profile).Preload("SaveDirectories").Preload("GameTypes").Preload("InfoTags", "TypeID IN (?)", profile.TypeID).Find(&saves)
	// if err != nil {
	// 	println("Error fetching the saves !")
	// 	println(err.Error())
	//		return []SaveDirectory{}
	//}
	return saves
}

func (this DataBase) UpdateProfileByID(profile *Profile) error {
	result := this.DB.Save(profile)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (this DataBase) DeleteProfile(profile *Profile) error {
	result := this.DB.Delete(profile)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// example
