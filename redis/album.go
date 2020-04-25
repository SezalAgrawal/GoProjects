package main

import (
	"errors"
	"fmt"
	"log"

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

// IncrementLikes increments like of an album
func IncrementLikes(id string) error {
	// establish connection with Redis server on deafult port 6379
	conn := pool.Get()
	defer conn.Close()

	exists, err := redis.Int(conn.Do("EXISTS", "album:"+id))
	if err != nil {
		return err
	} else if exists == 0 {
		return errNoAlbum
	}

	// start transaction to avoid race condition
	err = conn.Send("MULTI")
	if err != nil {
		return err
	}

	// increment like count
	err = conn.Send("HINCRBY", "album:"+id, "likes", 1)
	if err != nil {
		return err
	}
	// imcrement count in sorted set
	err = conn.Send("ZINCRBY", "likes", 1, id)
	if err != nil {
		return err
	}

	// execute transaction
	_, err = conn.Do("EXEC")
	if err != nil {
		return err
	}
	return nil
}

// FindTopThree finds top 3 albums based on likes
func FindTopThree() ([]Album, error) {
	// establish connection with Redis server on deafult port 6379
	conn := pool.Get()
	defer conn.Close()

	// infinite loop
	// can set number of attempts
	for {
		// watch for any changes in likes sorted set
		_, err := conn.Do("WATCH", "likes")
		if err != nil {
			return nil, err
		}

		// use ZREVRANGE to get top 3 albums
		ids, err := redis.Strings(conn.Do("ZREVRANGE", "likes", 0, 2))
		if err != nil {
			return nil, err
		}

		// start transaction to avoid race condition
		err = conn.Send("MULTI")
		if err != nil {
			return nil, err
		}

		// fetch album details
		for _, id := range ids {
			err := conn.Send("HGETALL", "album:"+id)
			if err != nil {
				return nil, err
			}
		}

		// execute transaction
		replies, err := redis.Values(conn.Do("EXEC"))
		if err == redis.ErrNil {
			// sorted set of likes got updated, while reading the values
			log.Println("trying again")
			continue
		} else if err != nil {
			return nil, err
		}

		albums := make([]Album, len(replies))

		for i, reply := range replies {
			var album Album
			err := redis.ScanStruct(reply.([]interface{}), &album)
			if err != nil {
				return nil, err
			}
			albums[i] = album
		}
		return albums, nil
	}
}
