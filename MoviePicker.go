package main

import (
	"fmt"
	"github.com/jmcvetta/napping"
//	"gopkg.in/jmcvetta/napping.v3"
	"log"
	"net/url"
)

type Film struct{
	Id 								 uint64 `json:"id"`
	Adult 						 bool	 `json:"adult"`
	GenreIDs					 []int `json:"genre_ids"`
	BackdropPath			 string `json:"backdrop_path"`
	OriginalLanguage  string `json:"original_language" `
	OriginalTitle  	 string `json:"original_title" `
	Overview					 string `json:"overview"`
	PosterPath				 string `json:"poster_path"`
	Popularity				 float32 `json:"popularity"`
	ReleaseDate 			 string `json:"release_date"`
	Title 					 	 string `json:"title"`
	Video							 bool `json:"video"`	
	VoteAverage 			 float32 `json:"vote_average"`
	VoteCount 				 int32 `json:"vote_count"`
}
type resultNum struct{
	TotalPages				int `json:"total_pages"`
	TotalResults			int `json:"total_results"`
}
	
func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}

func main() {
	//
	// Struct to hold error response
	//
	e := struct {
		Message string
	}{}
	
	// Define Search query as a URL keys and values
	search := url.Values{}
	search.Set("api_key","6442b8ee0e13c4415af27562719f67e9")
	search.Add("query","james bond")
	// Print search !
	fmt.Println(search.Encode())
	// Define Results as an array of Films !
	result := struct{
		Page 					int `json:"page"`
		Results 			[]Film `json:"results"`
		TotalPages 	int `json:"total_page"`
		TotalResults int `json:"total_results"`
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
		fmt.Printf("%#v\n", result)
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
	fmt.Println("--------------------------------------------------------------------------------")
	println("")
}
