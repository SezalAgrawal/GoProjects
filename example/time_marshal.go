package main

import (
	"fmt"
	"time"
)

type book struct {
	Create time.Time `json:"time" bson:"time"`
}

func test() {
	//var b book
	// a := book{
	// 	Create: time.Date(2020, 2, 2, 0, 0, 0, 0, time.Local),
	// }
	// obj, _ := json.Marshal(a)
	// fmt.Printf("%s\n", obj)
	// if err := json.Unmarshal(obj, &b); err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(b)

	// var i primitive.DateTime
	// i = 1580601600000
	// c := primitive.M{
	// 	"time": i,
	// }
	// fmt.Printf("%T", c["time"])
	// obj, _ := json.Marshal(c)
	// fmt.Printf("%s\n", obj)
	// if err := json.Unmarshal(obj, &b); err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(b)

	// var d primitive.DateTime
	// d = -62135596800000
	// t := time.Unix(int64(d)/1000, int64(d)%1000*1000000)
	// fmt.Println(t)
	// fmt.Println(t.UTC())
	// x,_ := json.Marshal(t)
	// fmt.Printf("%s", x)

	// var t1, m time.Time
	// t1 = time.Date(0001, 1, 1, 0, 0, 0, 0, time.Local)
	// obj, _ := json.Marshal(t1)
	// if err := json.Unmarshal(obj, &m); err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(m)
	// fmt.Println(t1.Zone())
	// t2 = time.Date(2001, 1, 1, 0, 0, 0, 0, time.Local)
	// obj, _ = json.Marshal(t2)
	// if err := json.Unmarshal(obj, &m); err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(m)
	// fmt.Println(t2.Zone())
	// var b book
	// var i primitive.DateTime
	// i = 1580601600000
	// c := primitive.M{
	// 	"time": i,
	//
	// obj, _ := json.Marshal(c)
	// if err := json.Unmarshal(obj, &b); err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(b)
	given := time.Date(0001, 1, 1, 0, 0, 0, 0, time.Local)
	//t1 := given.Local().UTC()
	// want := time.Date(0001, 1, 1, 0, 0, 0, 0, time.UTC)
	// equal := t1.UTC().Equal(t3)
	fmt.Println(given.UTC())
}
