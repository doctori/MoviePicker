package main

import (
	"fmt"
	"github.com/jmcvetta/napping"
	"log"
	"flag"
	"net/url"
)

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

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
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
		//
	// Struct to hold error response
	//
	e := struct {
		Message string
	}{}
	// Define Search query as a URL keys and values
	search := url.Values{}
	search.Set("s",film)
	search.Add("type","movie")
	search.Add("r","json")
	// Define Result as an array of OMDBFilmSearchResult !
	result := struct{
		Results 		[]OMDBFilmSearchResult `json:"Search"`
	}{}

	// Start Session
  	s := napping.Session{}
  	url := "https://omdbapi.com/"
  	fmt.Println("URL:>", url)
  
  	resp, err := s.Get(url,&search,&result,nil)
  	if err != nil {
  		log.Fatal(err)
  	}
  
  //
  // Process Response 
  //
  println("")
  fmt.Println("Response Status:",resp.Status())
  
  if resp.Status() == 200 {
		// fmt.Printf("Result: %s\n\n", resp.response)
		// resp.Unmarshal(&e)
		// Iterate through the results
		for index,film := range result.Results {
			fmt.Println("Index :",index)
			fmt.Println("Title :", film.Title)
			fmt.Println("imdbID : ",film.IMDBID)
		}
	} else {
		fmt.Println("Bad response status from OMDB server")
		fmt.Printf("\t Status:  %v\n", resp.Status())
		fmt.Printf("\t Message: %v\n", e.Message)
	}
}

// Find Films on TMDB

func searchTMDB(film string) {
	//
	// Struct to hold error response
	//
	e := struct {
		Message string
	}{}
	
	// Define Search query as a URL keys and values
	search := url.Values{}
	search.Set("api_key","6442b8ee0e13c4415af27562719f67e9")
	search.Add("query",film)
	// Define Results as an array of Films !
	result := struct{
		Page 				int `json:"page"`
		Results 			[]Film `json:"results"`
		TotalPages 			int `json:"total_page"`
		TotalResults 		int `json:"total_results"`
	}{}
	//Hope that works !
	
	
	
  // Start Session
  s := napping.Session{}
  url := "https://api.themoviedb.org/3/search/movie"
  fmt.Println("URL:>", url)
  
  resp, err := s.Get(url,&search,&result,nil)
  if err != nil {
  	log.Fatal(err)
  }
  
  //
  // Process Response 
  //
  println("")
  fmt.Println("Response Status:",resp.Status())
  
  if resp.Status() == 200 {
		// fmt.Printf("Result: %s\n\n", resp.response)
		// resp.Unmarshal(&e)
		// Iterate through the results
		for index,film := range result.Results {
			fmt.Println("Index :",index)
			fmt.Println("Title :", film.Title)
			fmt.Println("Note : ",film.VoteAverage)
		}
	} else {
		fmt.Println("Bad response status from TMDB server")
		fmt.Printf("\t Status:  %v\n", resp.Status())
		fmt.Printf("\t Message: %v\n", e.Message)
	}

}