package main

import (
	"github.com/jmcvetta/napping"
	"github.com/BurntSushi/toml"
	"log"
	"os"
	"flag"
	"net/url"
	"bytes"
	"fmt"
	"runtime"
)
type TmdbConfig struct{
	URL 		string
	MovieURL 	string
	ApiKey  	string
}
type Omdbconfig struct{
	URL 	string
}
type Config struct{
	TMDB 	TmdbConfig
	OMDB 	Omdbconfig
}
type Genre struct {
	Id 							uint64 `json:"id"`
	Name 						string `json:"name"`
}
type ProdCorp struct{
	Id 							uint64 `json:"id"`
	Name 						string `json:"name"`
}
type SpokenLanguage struct{
	Iso639 						string `json:"iso_639_1"`
	Name 						string `json:"name"`
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
	IMDBID						string `json:"imdb_id"`
}
type TMDBFilm struct{
	Id 							uint64 `json:"id"`
	Adult 						bool `json:"adult"`
	Budget 						uint64 `json:"budget"`
	Genres 						[]Genre `json:"genre"`
	Homepage 					string `json:"homepage"`
	BackdropPath 				string `json:"backdrop_path"`
	OriginalLanguage 			string `json:"original_language" `
	OriginalTitle 				string `json:"original_title" `
	Overview 					string `json:"overview"`
	PosterPath 					string `json:"poster_path"`
	Popularity 					float32 `json:"popularity"`
	ProdCorps 					[]ProdCorp `json:"production_companies"`
	ReleaseDate 				string `json:"release_date"`
	Title 						string `json:"title"`
	Video 						bool `json:"video"`	
	VoteAverage 				float32 `json:"vote_average"`
	VoteCount 					int32 `json:"vote_count"`
	IMDBID						string `json:"imdb_id"`
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
var (
	// Log Levels
    Trace   *log.Logger
    Info    *log.Logger
    Warning *log.Logger
    Error   *log.Logger
    // Configuration
    conf Config
)


func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)

	if _, err := toml.DecodeFile("conf.toml", &conf); err != nil {
  	log.Fatal("Could Not Load Configuration")
  	return
  	}
  	Trace = log.New(os.Stdout,
        "TRACE  : ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Info = log.New(os.Stdout,
        "INFO   : ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Warning = log.New(os.Stdout,
        "WARNING: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Error = log.New(os.Stderr,
        "ERROR: ",
        log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	runtime.GOMAXPROCS(2)
	// we retrieve the movie name
	movie :=  flag.String("movie","James Bond","The Name of the movie you are looking for")
	strict := flag.Bool("strict",false,"Will find the first Movie and get Detailed info on it")
	flag.Parse()
	if (*strict){
		// Open Channels 
		TmdbFilm := make(chan TMDBFilm)
		go getTMDBDetails(*movie,TmdbFilm)
		Info.Printf("Will Find this movie [%s] and no one else\n",*movie)
		film := <-TmdbFilm
		Trace.Println("IMDéBé : ",film.IMDBID)
	}else{
		Info.Println(" ---- TMDB ---- ")
		searchTMDB(*movie)
		Info.Println("--------------------------------------------------------------------------------")
		Info.Println(" ---- OMDB ---- ")
		searchOMDB(*movie)
		Info.Println("--------------------------------------------------------------------------------")
	
	}
}

// Get IMDBID from TMDB (yes)
func getTMDBDetails(film string,messages chan TMDBFilm){
	search := url.Values{}
	search.Set("api_key",conf.TMDB.ApiKey)
	search.Add("query",film)
	result := struct{
		Page 				int `json:"page"`
		Results 			[]Film `json:"results"`
		TotalPages 			int `json:"total_page"`
		TotalResults 		int `json:"total_results"`
	}{}

  	getRESTResponse(conf.TMDB.URL,&search,&result)
  	
	Trace.Printf("Title : %s [%s]", result.Results[0].Title,result.Results[0].ReleaseDate)
	Trace.Println("Note : ",result.Results[0].VoteAverage)
	search.Del("query")
	// Create TMDB URL
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprint(conf.TMDB.MovieURL,"/",result.Results[0].Id))

	movieURL := buffer.String()
	var resultDetail TMDBFilm
	getRESTResponse(movieURL,&search,&resultDetail)
	done := make(chan struct{})
	go cacheTMDBFilm(resultDetail,done)
	Trace.Printf("%v\n",resultDetail)
	<-done
	messages <- resultDetail

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
	for _,film := range result.Results {
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

	Info.Println("Title :",result.Title)
	Trace.Println("* IMDB Rating : ",result.IMDBRating )
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
		Trace.Printf("Title : %s [%s]", film.Title,film.ReleaseDate)
		Trace.Println("Note : ",film.VoteAverage)
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
  	Trace.Printf("%s Return Status [%v]",url,resp.Status())
  	if resp.Status() == 200 {
  		return
  	} else {
  		log.Fatal("Resource could not be fetched")
		log.Fatal("\t Status:  %v\n", resp.Status())
		log.Fatal("\t Message: %v\n", e.Message)
  	}
}