package jwttoken

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestGenerate(t *testing.T) {
	token, err := GetDefaultTokenGenerator().Generate("Issuer", PTRefresh, "bcf6169f-4059-4c78-93e4-a2e6efd5620e", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	println(token)
}

func TestVerify(t *testing.T) {
	var p JwtPayload
	valid, err := GetDefaultTokenGenerator().Verify("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJGaWxlREFHIiwiZXhwIjoxNjI2ODM4NTY1LCJpYXQiOjE2MjYyMzM3NjUsImp0aSI6IjljMjAzZGFjLWYyNzItNDMyMi1hYTcyLTc3MDNkMzU5YTEzYyIsInBlcm0iOiJyZWZyZXNoIiwidWlkIjoiYmNmNjE2OWYtNDA1OS00Yzc4LTkzZTQtYTJlNmVmZDU2MjBlIn0.2bmmwiU8wzEI9yxgP4istmsJo8oOZqWf1svjxUmUquo", &p)
	if err != nil {
		t.Fatal(err)
	}
	if !valid {
		println("token has expired")
	} else {
		println("token is valid")
	}
	fmt.Printf("%+v\n", p)
}

func TestCustomJWT(t *testing.T) {
	type Extend struct {
		Filed1 string
		Filed2 string
	}
	ext := &Extend{Filed1: "filed1", Filed2: "filed2"}
	data, _ := json.Marshal(ext)
	token, err := GetDefaultTokenGenerator().Generate("FileDrive", PTCustom, "bcf6169f-4059-4c78-93e4-a2e6efd5620e", "", string(data))
	if err != nil {
		t.Fatal(err)
	}
	println(token)

	var p JwtPayload
	valid, err := GetDefaultTokenGenerator().Verify(token, &p)
	if err != nil {
		t.Fatal(err)
	}
	if !valid {
		t.Fatal("token is invalid")
	}
	if p.Permission != PTCustom {
		t.Fatal("token is not custom type")
	}
	ext2 := &Extend{}
	err = json.Unmarshal([]byte(p.Extend.(string)), ext2)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(ext, ext2) {
		t.Fatal("ext is not equal ext2")
	}
}
