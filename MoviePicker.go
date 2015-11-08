package main

import (
	"fmt"
	"github.com/jmcvetta/napping"
	"github.com/BurntSushi/toml"
	"log"
	"flag"
	"net/url"
)
type TmdbConfig struct{
	URL 	string
	ApiKey  string
}
type Omdbconfig struct{
	URL 	string
}
type Config struct{
	TMDB 	TmdbConfig
	OMDB 	Omdbconfig
}
type Film struct{
	Id 							uint64 `json:"id"`
	Adult 						bool `json:"adult"`
	GenreIDs 					[]int `json:"genre_ids"`
	BackdropPath 				string `json:"backdrop_path"`
	OriginalLanguage 			string `json:"original_language" `
	OriginalTitle 				string `json:"original_title" `
	Overview 					string `json:"overview"`
	PosterPath 					string `json:"poster_path"`
	Popularity 					float32 `json:"popularity"`
	ReleaseDate 				string `json:"release_date"`
	Title 						string `json:"title"`
	Video 						bool `json:"video"`	
	VoteAverage 				float32 `json:"vote_average"`
	VoteCount 					int32 `json:"vote_count"`
}
type OMDBFilmSearchResult struct{
	Title 						string
	Year 						string
	IMDBID						string
	Type 						string
	Poster 						string

}

var conf Config

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)

	if _, err := toml.DecodeFile("conf.toml", &conf); err != nil {
  	log.Fatal("Could Not Load Configuration")
  	return
  	}
}

func main() {
	movie := "Jame Bond"
	// we retrieve the movie name
	flag.StringVar(&movie,"movie","James Bond","The Name of the movie you are looking for")
	flag.Parse()
	fmt.Println(" ---- TMDB ---- ")
	searchTMDB(movie)
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println(" ---- OMDB ---- ")
	searchOMDB(movie)
	fmt.Println("--------------------------------------------------------------------------------")
	println("")
}

// Find Films on OMDB
func searchOMDB(film string) {

	// Define Search query as a URL keys and values
	search := url.Values{}
	search.Set("s",film)
	search.Add("type","movie")
	search.Add("r","json")
	// Define Result as an array of OMDBFilmSearchResult !
	result := struct{
		Results 		[]OMDBFilmSearchResult `json:"Search"`
	}{}
  	
  	url := conf.OMDB.URL
  	getRESTResponse(url,&search,&result)
	// Iterate through the results
	for index,film := range result.Results {
		fmt.Println("Index :",index)
		fmt.Println("Title :", film.Title)
		fmt.Println("imdbID : ",film.IMDBID)
	}
}

func getOMDBDetails(IMDBID string){
}
// Find Films on TMDB

func searchTMDB(film string) {
	
	// Define Search query as a URL keys and values
	search := url.Values{}
	search.Set("api_key",conf.TMDB.ApiKey)
	search.Add("query",film)
	// Define Results as an array of Films !
	result := struct{
		Page 				int `json:"page"`
		Results 			[]Film `json:"results"`
		TotalPages 			int `json:"total_page"`
		TotalResults 		int `json:"total_results"`
	}{}
	//Hope that works !
	
  	url := conf.TMDB.URL
  	getRESTResponse(url,&search,&result)
  
	// Iterate through the results
	for index,film := range result.Results {
		fmt.Println("Index :",index)
		fmt.Println("Title :", film.Title)
		fmt.Println("Note : ",film.VoteAverage)
	}

}

func getRESTResponse(url string, p *url.Values, result interface{}){
	//
	// Struct to hold error response
	//
	e := struct {
		Message string
	}{}

	log.Println("URL => ", url)
	s := napping.Session{}
	resp, err := s.Get(url,p,result,e)
	if err != nil {
  		log.Fatal(err)
  	}
  	log.Println("Response Status:",resp.Status())
  	if resp.Status() == 200 {
  		return
  	} else {
  		log.Fatal("Resource could not be fetched")
		log.Fatal("\t Status:  %v\n", resp.Status())
		log.Fatal("\t Message: %v\n", e.Message)
  	}
}