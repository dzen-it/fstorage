package api

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"strings"
	"testing"

	"golang.org/x/crypto/sha3"
)

func TestPositiveControlHashFileReader(t *testing.T) {
	s := "Hello hello world"
	reader := strings.NewReader(s)

	a := md5.Sum([]byte(s))
	expMD5 := hex.EncodeToString(a[:])

	b := sha256.Sum256([]byte(s))
	expSHA256 := hex.EncodeToString(b[:])

	c := sha512.Sum512([]byte(s))
	expSHA512 := hex.EncodeToString(c[:])

	d := sha3.Sum256([]byte(s))
	expKeccak256 := hex.EncodeToString(d[:])

	e := sha3.Sum512([]byte(s))
	expKeccak512 := hex.EncodeToString(e[:])

	f := sha1.Sum([]byte(s))
	expSHA1 := hex.EncodeToString(f[:])

	err := controlHashFileReader(
		map[string]string{
			XFSHashControlMD5:       expMD5,
			XFSHashControlSHA256:    expSHA256,
			XFSHashControlSHA512:    expSHA512,
			XFSHashControlKeccak256: expKeccak256,
			XFSHashControlKeccak512: expKeccak512,
			XFSHashControlSHA1:      expSHA1,
		},
		reader,
	)

	if err != nil {
		t.Fatal("Error:", err)
	}

}

func TestNegativeControlHashFileReader(t *testing.T) {
	s := "Hello hello world"
	reader := strings.NewReader(s)

	a := md5.Sum([]byte(s))
	expMD5 := hex.EncodeToString(a[:])

	b := sha256.Sum256([]byte(s))
	expSHA256 := hex.EncodeToString(b[:])

	c := sha512.Sum512([]byte(s))
	expSHA512 := hex.EncodeToString(c[:])

	d := sha3.Sum256([]byte(s))
	expKeccak256 := hex.EncodeToString(d[:])

	e := sha3.Sum512([]byte(s))
	expKeccak512 := hex.EncodeToString(e[:])

	f := sha1.Sum([]byte(s))
	expSHA1 := hex.EncodeToString(f[:])

	err := controlHashFileReader(
		map[string]string{
			XFSHashControlMD5:       expMD5,
			XFSHashControlSHA256:    expSHA256,
			XFSHashControlSHA512:    expSHA512,
			XFSHashControlKeccak256: expKeccak256 + "fail",
			XFSHashControlKeccak512: expKeccak512,
			XFSHashControlSHA1:      expSHA1,
		},
		reader,
	)

	if err == nil {
		t.Fatal("Error:", err)
	}
}
