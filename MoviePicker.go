package main

import (
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
type OMDBFilm struct{
	Title 						string
	Year 						string
	Rated 						string
	Released 					string
	Runtime 					string
	Genre 						string
	Director 					string
	Writer 						string
	Actors 						string
	Plot 						string
	Language 					string
	Country						string
	Awards 						string
	Poster 						string
	Metascore 					string
	IMDBRating 					string `json:"imdbRating"`
	IMDBVotes					string `json:"imdbVotes"`
	IMDBID  					string `json:"imdbID"`
	Type 						string 
	TomatoMeter 				string `json:"tomatoMeter"`
  	TomatoImage 				string `json:"tomatoImage"`
  	TomatoRating 				string `json:"tomatoRating"`
    TomatoReviews 				string `json:"tomatoReviews"`
  	TomatoFresh 				string `json:"tomatoFresh"`
   	TomatoRotten 				string `json:"tomatoRotten"`
  	TomatoConsensus				string `json:"tomatoConsensus"`
  	TomatoUserMeter				string `json:"tomatoUserMeter"`
  	TomatoUserRating 			string `json:"tomatoUserRating"`
  	TomatoUserReviews 			string `json:"tomatoUserReviews"`
  	DVD							string
  	BoxOffice					string 
  	Production					string
  	Website 					string
    Response 					string
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
	log.Println(" ---- TMDB ---- ")
	searchTMDB(movie)
	log.Println("--------------------------------------------------------------------------------")
	log.Println(" ---- OMDB ---- ")
	searchOMDB(movie)
	log.Println("--------------------------------------------------------------------------------")
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
  	
  	getRESTResponse(conf.OMDB.URL,&search,&result)
	// Iterate through the results
	for index,film := range result.Results {
		log.Println("Index :",index)
		log.Println("Title :", film.Title)
		log.Println("imdbID : ",film.IMDBID)
		getOMDBDetails(film.IMDBID)
	}
}

func getOMDBDetails(IMDBID string){
	search := url.Values{}
	search.Set("i",IMDBID)
	search.Add("type","movie")
	search.Add("tomatoes","true")
	var result OMDBFilm
	getRESTResponse(conf.OMDB.URL,&search,&result)

	log.Println("Title :",result.Title)
	log.Println("* IMDB Rating : ",result.IMDBRating )
	log.Println("* RottenTomato Rating : ",result.TomatoRating)
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
		log.Println("Index :",index)
		log.Printf("Title : %s [%s]", film.Title,film.ReleaseDate)
		log.Println("Note : ",film.VoteAverage)
	}

}

func getRESTResponse(url string, p *url.Values, result interface{}){
	//
	// Struct to hold error response
	//
	e := struct {
		Message string
	}{}

	s := napping.Session{}
	resp, err := s.Get(url,p,result,e)
	if err != nil {
  		log.Fatal(err)
  	}
  	if resp.Status() == 200 {
  		return
  	} else {
  		log.Fatal("Resource could not be fetched")
		log.Fatal("\t Status:  %v\n", resp.Status())
		log.Fatal("\t Message: %v\n", e.Message)
  	}
}