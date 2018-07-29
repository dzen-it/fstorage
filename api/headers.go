package api

import (
	"net/http"
	"regexp"
)

const (
	//Example: X-FStorage-Hash-Control-MD5: abcdef...,
	XFSHashControlMD5       = "X-FStorage-Hash-Control-MD5"
	XFSHashControlSHA1      = "X-FStorage-Hash-Control-SHA1"
	XFSHashControlSHA256    = "X-FStorage-Hash-Control-SHA256"
	XFSHashControlSHA512    = "X-FStorage-Hash-Control-SHA512"
	XFSHashControlKeccak256 = "X-FStorage-Hash-Control-Keccak256"
	XFSHashControlKeccak512 = "X-FStorage-Hash-Control-Keccak512"
)

var (
	regexpMatchMD5HexString       *regexp.Regexp
	regexpMatchSHA1HexString      *regexp.Regexp
	regexpMatchSHA256HexString    *regexp.Regexp
	regexpMatchSHA512HexString    *regexp.Regexp
	regexpMatchKeccak256HexString *regexp.Regexp
	regexpMatchKeccak512HexString *regexp.Regexp
)

func init() {
	var err error
	regexpMatchMD5HexString, err = regexp.Compile("[a-f0-9]{32}")
	if err != nil {
		panic(err)
	}
	regexpMatchSHA1HexString, err = regexp.Compile("[a-f0-9]{40}")
	if err != nil {
		panic(err)
	}
	regexpMatchSHA256HexString, err = regexp.Compile("[a-f0-9]{64}")
	if err != nil {
		panic(err)
	}
	regexpMatchSHA512HexString, err = regexp.Compile("[a-f0-9]{128}")
	if err != nil {
		panic(err)
	}
	regexpMatchKeccak256HexString, err = regexp.Compile("[a-f0-9]{64}")
	if err != nil {
		panic(err)
	}
	regexpMatchKeccak512HexString, err = regexp.Compile("[a-f0-9]{128}")
	if err != nil {
		panic(err)
	}
}

func getHeadersHashControl(r *http.Request) (map[string]string, error) {
	m := make(map[string]string)

	if err := validateAndSetHashHeader(r, XFSHashControlMD5, m); err != nil {
		return nil, err
	}
	if err := validateAndSetHashHeader(r, XFSHashControlSHA1, m); err != nil {
		return nil, err
	}
	if err := validateAndSetHashHeader(r, XFSHashControlSHA256, m); err != nil {
		return nil, err
	}
	if err := validateAndSetHashHeader(r, XFSHashControlKeccak512, m); err != nil {
		return nil, err
	}
	if err := validateAndSetHashHeader(r, XFSHashControlSHA512, m); err != nil {
		return nil, err
	}
	if err := validateAndSetHashHeader(r, XFSHashControlKeccak256, m); err != nil {
		return nil, err
	}

	return m, nil
}

func validateAndSetHashHeader(r *http.Request, header string, hashMap map[string]string) error {
	hex := r.Header.Get(header)
	if len(hex) == 0 {
		return nil
	}

	var reg *regexp.Regexp

	switch header {
	case XFSHashControlMD5:
		reg = regexpMatchMD5HexString
	case XFSHashControlSHA1:
		reg = regexpMatchSHA1HexString
	case XFSHashControlSHA256:
		reg = regexpMatchSHA256HexString
	case XFSHashControlSHA512:
		reg = regexpMatchSHA512HexString
	case XFSHashControlKeccak256:
		reg = regexpMatchKeccak256HexString
	case XFSHashControlKeccak512:
		reg = regexpMatchKeccak512HexString
	default:
		return nil
	}

	if !reg.MatchString(hex) {
		return ErrInvalidHeader
	}

	hashMap[header] = hex
	return nil
}
