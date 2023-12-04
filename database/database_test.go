package database_test

import (
	// "fmt"
	"log"
	"math/rand"
	"savemanager/config"
	"savemanager/database"
	"strconv"
	"testing"
	// "regexp"
)

// // TestHelloName calls greetings.Hello with a name, checking
// // for a valid return value.
// func TestHelloName(t *testing.T) {
//     name := "Gladys"
//     want := regexp.MustCompile(`\b`+name+`\b`)
//     msg, err := Hello("Gladys")
//     if !want.MatchString(msg) || err != nil {
//         t.Fatalf(`Hello("Gladys") = %q, %v, want match for %#q, nil`, msg, err, want)
//     }
// }

// // TestHelloEmpty calls greetings.Hello with an empty string,
// // checking for an error.
// func TestHelloEmpty(t *testing.T) {
//     msg, err := Hello("")
//     if msg != "" || err == nil {
// t.Fatalf(`Hello("") = %q, %v, want "", error`, msg, err)
//     }
// }

// Constants
const OPERATION_COUNT int = 200

// Utility functions

func UnitProfileCreation(dbName string) {
	configuration := config.GlobalConfiguration{}
	db := database.DataBase{Name: dbName}

	configuration.Init(dbName)
	db.Connect(&configuration)

	// rand.Int()
	var secondaryIndex uint = 1
	for i := 0; i < OPERATION_COUNT; i++ {
		if secondaryIndex > 7 {
			secondaryIndex = 1
		}
		typing, err := db.GetGameTypeByID(secondaryIndex)
		if err != nil {
			log.Fatalf("The db could not retrieve the gametype for index %v ", secondaryIndex)
		}
		secondaryIndex++
		I := strconv.Itoa(i)

		profile := &database.Profile{
			ProfileName: "test-" + I,
			GamePath:    "path-" + I,
			Type:        typing,
		}
		_, errDb := db.CreateProfile(profile)
		if errDb != nil {

			log.Fatalf("An error occured while creating single profiles on index %v", i)
		}
	}
}
func BatchProfileCreation(dbName string) {

	configuration := config.GlobalConfiguration{}
	db := database.DataBase{Name: dbName}

	configuration.Init(dbName)
	db.Connect(&configuration)

	var secondaryIndex uint = 1
	var gameTypes []database.GameType
	db.DB.Find(&gameTypes)
	var profiles = make([]*database.Profile, OPERATION_COUNT)
	for i := 0; i < OPERATION_COUNT; i++ {
		if secondaryIndex > 6 {
			secondaryIndex = 1
		}
		typing := gameTypes[secondaryIndex]
		secondaryIndex++
		I := strconv.Itoa(i)

		// typing := &GameType{ID:1}
		profile := &database.Profile{
			ProfileName: "test-" + I,
			GamePath:    "path-" + I,
			Type:        &typing,
			// SaveDirectories: []*SaveDirectory{{
			// 	Name: "abc",
			// 	Path: "def",
			// 	//	Profile:  &database.Profile{},
			// 	// InfoTags: []*InfoTag{},
			// }},
		}
		profiles[i] = profile

	}
	errDb := db.BatchCreateProfile(profiles)
	if errDb != nil {

		log.Fatalf("An error occured while creating batch profiles.")
	}
}

// Tests

func TestBatchProfileCreation(t *testing.T) {
	const dbName = "test.batchProfile"
	BatchProfileCreation(dbName)
}

func TestSingleProfileCreation(t *testing.T) {
	const dbName = "test.singleProfile"
	UnitProfileCreation(dbName)
}

// Benchmarks

func inactive_BenchmarkSingleProfileCreation(b *testing.B) {
	for j := 0; j < b.N; j++ {
		var dbName = "benchmark.singleProfile" + strconv.Itoa(j)
		UnitProfileCreation(dbName)
	}
}
func inactive_BenchmarkBatchProfileCreation(b *testing.B) {
	for j := 0; j < b.N; j++ {

		var dbName = "benchmark.batchProfile" + strconv.Itoa(j)
		BatchProfileCreation(dbName)
	}
}

func BenchmarkRandInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rand.Int()
	}
}
