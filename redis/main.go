package main

// Usecase: online record shop, and want to store info about albums for sale in Redis
// store info of a album as hash: fields like title, artist, likes, price
// key of the hash is album:{id}

// APIS:
// curl -i localhost:4000/album?id=2
// curl -i -L -d "id=2" localhost:4000/like
// curl -i localhost:4000/popular

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
)

func main() {
	// initialize connection pool
	pool = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/album", showAlbum)
	mux.HandleFunc("/like", addLike)
	mux.HandleFunc("/popular", listPopular)
	log.Println("Listening on 4000...")
	http.ListenAndServe(":4000", mux)
}

func showAlbum(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, http.StatusText(405), 405)
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	if _, err := strconv.Atoi(id); err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	bk, err := FindAlbum(id)
	if err == errNoAlbum {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	fmt.Fprintf(w, "%s by %s: £%.2f [%d likes] \n", bk.Title, bk.Artist, bk.Price, bk.Likes)
}

func addLike(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, http.StatusText(405), 405)
		return
	}
	id := r.PostFormValue("id")
	if id == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	if _, err := strconv.Atoi(id); err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	err := IncrementLikes(id)
	if err == errNoAlbum {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	// redirect to show album
	http.Redirect(w, r, "/album?id="+id, 303)
}

func listPopular(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, http.StatusText(405), 405)
		return
	}
	albums, err := FindTopThree()
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	for i, ab := range albums  {
		fmt.Fprintf(w, "%d) %s by %s: £%.2f [%d likes] \n", i+1, ab.Title, ab.Artist, ab.Price, ab.Likes)
	}
	
}
