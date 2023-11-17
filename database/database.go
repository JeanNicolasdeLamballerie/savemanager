package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

type DataBase struct {
	Name string
	DB   *gorm.DB
}

func (this *DataBase) Connect() error {
	DB, err := gorm.Open(sqlite.Open("main.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	err = DB.AutoMigrate(&Profile{}, &SaveDirectory{})
	if err != nil {
		log.Fatal(err)
	}

	this.DB = DB
	return nil
}

type SaveDirectory struct {
	gorm.Model
	Name      string `gorm:"unique"`
	Path      string `gorm:"notNull"`
	ProfileID uint
	Profile   Profile `gorm:"foreignKey:ProfileID"`
	TypeID    uint
	Type      GameType `gorm:"foreignKey:TypeID"`
}

type Profile struct {
	gorm.Model
	ProfileName     string `gorm:"unique"`
	GamePath        string
	SaveDirectories []SaveDirectory `gorm:"foreignKey:ProfileID"`
}
type GameType struct {
	gorm.Model
	Name string `gorm:"unique"`
}

// func main() {
// 	this.DB, err := connectToSQLite()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer func() {
// 		this.DBInstance, _ := db.DB()
// 		_ = this.DBInstance.Close()
// 	}()

// 	// Perform database migration
// 	err = this.DB.AutoMigrate(&SaveDirectory{}, &Profile{})
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Your CRUD operations go here

// }

// INSERT
func (this DataBase) CreateProfile(profile *Profile) (*uint, error) {

	result := this.DB.Create(profile)
	if result.Error != nil {
		return nil, result.Error
	}
	print(profile, &profile)
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
	err := this.DB.Model(profile).Association("SaveDirectories").Find(&saves)
	if err != nil {
		println("Error fetching the saves !")
		println(err)
		return []SaveDirectory{}
	}
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
