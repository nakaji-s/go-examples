package main

import (
	"log"

	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/kr/pretty"
	_ "github.com/lib/pq"
)

type Language struct {
	ID   uint `gorm:"primary_key"`
	Name string
}

type Movie struct {
	ID         uint `gorm:"primary_key"`
	Title      string
	Language   Language
	LanguageID uint
}

type Artist struct {
	ID     uint `gorm:"primary_key"`
	Name   string
	Movies []Movie `gorm:"many2many:artist_movies"`
}

func createArtists() {
	langs := []Language{{Name: "english"},
		{Name: "tamil"},
		{Name: "hindi"}}

	for i, _ := range langs {
		if err := db.Create(&langs[i]).Error; err != nil {
			log.Fatal(err)
		}
	}

	movies := []Movie{
		{Title: "Nayagan", Language: langs[1]},
		{Title: "Anbe sivam", Language: langs[1]},
		{Title: "3 idiots", Language: langs[2]},
		{Title: "Shamithab", Language: langs[2]},
		{Title: "Dark Knight", Language: langs[0]},
		{Title: "310 to Yuma", Language: langs[0]},
	}
	for i, _ := range movies {
		if err := db.Create(&movies[i]).Error; err != nil {
			log.Fatal(err)
		}
	}

	artists := []Artist{{Name: "Madhavan", Movies: []Movie{movies[1], movies[2]}},
		{Name: "Kamal Hassan", Movies: []Movie{movies[0], movies[1]}},
		{Name: "Dhanush", Movies: []Movie{movies[3]}},
		{Name: "Aamir Khan", Movies: []Movie{movies[2]}},
		{Name: "Amitabh Bachchan", Movies: []Movie{movies[3]}},
		{Name: "Christian Bale", Movies: []Movie{movies[4], movies[5]}},
		{Name: "Russell Crowe", Movies: []Movie{movies[5]}},
	}

	for i, _ := range artists {
		if err := db.Create(&artists[i]).Error; err != nil {
			log.Fatal(err)
		}
	}
}

var db *gorm.DB

func main() {
	var err error
	db, err = gorm.Open("postgres", "host=localhost user=postgres sslmode=disable password=admin")
	if err != nil {
		log.Fatal(err)
	}

	db.DropTableIfExists(new(Language), new(Movie), new(Artist))
	db.AutoMigrate(new(Language), new(Movie), new(Artist))
	db.LogMode(true)

	createArtists()

	// Get the list of all artist who acted in "english" movies
	var artists []Artist
	//if err = db.Joins("JOIN artist_movies on artist_movies.artist_id=artists.id").
	//	Joins("JOIN movies on movies.id=artist_movies.movie_id").
	//	Joins("JOIN languages on movies.language_id=languages.id").
	//	Where("languages.name=?", "tamil").
	//	Group("artists.id").
	//	Preload("Movies").
	//	Find(&artists).Error; err != nil {
	//	log.Fatal(err)
	//}
	//pretty.Println(artists)

	artists = []Artist{}
	artistIds := []uint{}
	rows, _ := db.Model(&Artist{}).Rows()
	for rows.Next() {
		artist := Artist{}
		db.ScanRows(rows, &artist)
		artistIds = append(artistIds, artist.ID)
	}
	rows.Close()
	type ArtistMovie struct {
		ArtistId uint
		MovieId  uint
	}
	artistMovies := []ArtistMovie{}
	fmt.Println(artistIds)
	db.Select("*").Table("artist_movies").Where("artist_id IN (?)", artistIds).Find(&artistMovies)
	pretty.Println(artistMovies)

	//
	//// Get the list the artists for movie "Nayagan"
	//artists = []Artist{}
	//if err = db.Joins("JOIN artist_movies on artist_movies.artist_id=artists.id").
	//	Joins("JOIN movies on artist_movies.movie_id=movies.id").Where("movies.title=?", "Nayagan").
	//	Group("artists.id").Find(&artists).Error; err != nil {
	//	log.Fatal(err)
	//}
	//
	//for _, ar := range artists {
	//	fmt.Println(ar.Name)
	//}
	//
	//// Get the list of artists for movies "3 idiots", "Shamitab" and "310 to Yuma"
	//artists = []Artist{}
	//
	//if err = db.Joins("JOIN artist_movies on artist_movies.artist_id=artists.id").
	//	Joins("JOIN movies on artist_movies.movie_id=movies.id").
	//	Where("movies.title in (?)", []string{"3 idiots", "Shamitabh", "310 to Yuma"}).
	//	Group("artists.id").Find(&artists).Error; err != nil {
	//	log.Fatal(err)
	//}
	//for _, ar := range artists {
	//	fmt.Println(ar.Name)
	//}
}
