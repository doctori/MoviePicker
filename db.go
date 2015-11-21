package main 

import (
	"gopkg.in/redis.v3"
	"encoding/json"
	"fmt"
	)
var (
	client *redis.Client
)
func init(){
    // Init Redis Client
    client = redis.NewClient(&redis.Options{
		Addr:		"localhost:6379",
		Password: 	"",
		DB:			0,
	})
}

func cacheTMDBFilm(film TMDBFilm,done chan struct{}){
	// Store the film as Json String
	fmt.Println("prout")
	jsonFilm,err := json.Marshal(film)
	if err != nil {
		panic(err)
	}
	err = client.Set(film.Title,jsonFilm,0).Err()
	if err != nil {
		panic(err)
	}
	Trace.Printf("[%s] saved into the cache (for later)\n",film.Title)
	close(done)
}

func getTMDBCacheFilm(film string,result chan<- TMDBFilm,found chan<- bool){
	// fetch the movie (if exist)
	var resultDetail TMDBFilm
	Trace.Println("Looking for [",film,"] in the cache");
	filmJson,err := client.Get(film).Result()
	if err == redis.Nil {
		Warning.Println(film," Not Found in Cache")
		fmt.Printf("%#v\n",err)
		found <- false
	}else if err != nil {
		panic(err)
	}else{
		err := json.Unmarshal([]byte(filmJson),&resultDetail)
		if err != nil {
			panic(err)
		}else{
			Trace.Println("Found [",film,"] in Ze Cache !")
			found <- true
			result <- resultDetail
		}
	}
	
}