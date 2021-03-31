package main

import (
	"crypto/rand"
	"math/big"
	"reflect"
	"testing"
)

func TestStore(t *testing.T) {
	expect := NewStore()

	passwordLength, err := rand.Int(rand.Reader, big.NewInt(256))
	if err != nil {
		t.Fatal(err)
	}
	password := make([]byte, passwordLength.Int64())
	n, err := rand.Read(password)
	if err != nil {
		t.Fatal(err)
	}
	if n != int(passwordLength.Int64()) {
		t.Fatalf("unexpected password length: %d", n)
	}

	data, err := WriteStore(password, expect)
	if err != nil {
		t.Fatal(err)
	}

	got, err := ReadStore(password, data)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expect.data, got.data) {
		t.Fatalf("expected %#v, got %#v", expect.data, got.data)
	}
}
