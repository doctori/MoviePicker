package main 

import (
	"gopkg.in/redis.v3"
	"encoding/json"
	"fmt"
	)

func cacheTMDBFilm(film TMDBFilm,done chan struct{}){
	// Store the film as Json String
	fmt.Println("prout")
	jsonFilm,err := json.Marshal(film)
	if err != nil {
		panic(err)
	}
	client := redis.NewClient(&redis.Options{
		Addr:		"localhost:6379",
		Password: 	"",
		DB:			0,
		})
	err = client.Set(film.Title,jsonFilm,0).Err()
	if err != nil {
		panic(err)
	}
	Trace.Printf("[%s] saved into the cache (for later)\n",film.Title)
	close(done)

}