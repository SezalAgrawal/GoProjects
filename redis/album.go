package main

import (
	"errors"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

// Holds pool of Redis connections
var pool *redis.Pool

var errNoAlbum = errors.New("no album found")

// Album struct to hold album info
type Album struct {
	Title  string  `redis:"title"`
	Artist string  `redis:"artist"`
	Price  float64 `redis:"price"`
	Likes  int     `redis:"likes"`
}

// FindAlbum finds album based on id
func FindAlbum(id string) (*Album, error) {
	// establish connection with Redis server on deafult port 6379
	conn := pool.Get()
	defer conn.Close()

	// add an album
	// _, err = conn.Do("HMSET", "album:1", "title", "That's the spirit", "artist", "BMTH", "price", 4.5, "likes", 100)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("ALbum added!")

	// get one field of the album
	// title, err := redis.String(conn.Do("HGET", "album:1", "title"))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(title)

	// get all fields of album
	values, err := redis.Values(conn.Do("HGETALL", "album:"+id))
	if err != nil {
		return nil, err
	} else if len(values) == 0 {
		return nil, errNoAlbum
	}
	var album Album
	err = redis.ScanStruct(values, &album)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%+v\n", album)
	return &album, nil
}
