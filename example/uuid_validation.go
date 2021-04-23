package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"gopkg.in/go-playground/validator.v9"
)

type example struct {
	U1 string `json:"u1" validate:"required,uuid4"`
	U2 string `json:"u2" validate:"required,uuid4"`
	U3 string `json:"u3" validate:"uuid4"`
}

func main() {
	fmt.Println("Hello, playground")
	u2, _ := uuid.NewUUID()
	u3, _ := uuid.NewUUID()
	s := map[string]interface{}{
		"u1": uuid.New(),
		"u2": u2.String(),
		"u3": u3,
	}
	data, err := json.Marshal(s)
	if err != nil {
		fmt.Println("here")
		fmt.Println(err)
	}
	var a example
	reader := bytes.NewReader(data)
	err = json.NewDecoder(reader).Decode(&a)
	if err != nil {
		fmt.Println("bye")
		fmt.Println(err)
	}
	var validate = validator.New()
	err = validate.Struct(a)
	if err != nil {
		fmt.Println("here")
		fmt.Printf("%v", err)
	}
	fmt.Println(a.U1)
}
