package imgbase64

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// Default Image in case the requested image
// does not exist.
var df string = ""

// Get Default Image
func DefaultImage() string {
	return df
}

// Set Default Image
func SetDefaultImage(img string) {
	df = img
}

// encode is our main function for
// base64 encoding a passed []byte
func encode(bin []byte) []byte {
	e64 := base64.StdEncoding

	maxEncLen := e64.EncodedLen(len(bin))
	encBuf := make([]byte, maxEncLen)

	e64.Encode(encBuf, bin)
	return encBuf
}

// Lightweight HTTP Client to fetch the image
// Note: This will also pull webpages. @todo
// It is up to the user to use valid image urls.
func get(url string) ([]byte, string) {
	resp, err := http.Get(url)
	if err != nil {
		//Blank 1px x 1 px gif
		return []byte("R0lGODlhAQABAIAAAP///wAAACH5BAAAAAAALAAAAAABAAEAAAICRAEAOw=="), "image/gif"
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	ct := resp.Header.Get("Content-Type")

	if resp.StatusCode == 200 && len(body) > 512 {
		return body, ct
	}

	if DefaultImage() == "" {
		return []byte(""), ct
	}

	if url == DefaultImage() {
		return []byte(""), ct
	}

	return get(DefaultImage())
}

// FromRemote is a better named function that
// presently calls NewImage which will be deprecated.
// Function accepts an RFC compliant URL and returns
// a base64 encoded result.
func FromRemote(url string) string {
	image, mime := get(cleanUrl(url))
	enc := encode(image)

	out := format(enc, mime)
	return out
}

// FromBuffer accepts a buffer and returns a
// base64 encoded string.
func FromBuffer(buf bytes.Buffer) string {
	enc := encode(buf.Bytes())
	mime := http.DetectContentType(buf.Bytes())

	return format(enc, mime)
}

// FromLocal reads a local file and returns
// the base64 encoded version.
func FromLocal(fname string) string {
	var b bytes.Buffer
	_, err := os.Stat(fname)
	if err != nil {
		if os.IsNotExist(err) {
			panic("File does not exist")
		}
		panic("Error stating file")
	}

	file, err := os.Open(fname)
	if err != nil {
		panic("Error opening file")
	}

	_, err = b.ReadFrom(file)
	if err != nil {
		panic("Error reading file to buffer")
	}

	return FromBuffer(b)
}

// format is an abstraction of the mime switch to create the
// acceptable base64 string needed for browsers.
func format(enc []byte, mime string) string {
	switch mime {
	case "image/gif", "image/jpeg", "image/pjpeg", "image/png", "image/tiff":
		return fmt.Sprintf("data:%s;base64,%s", mime, enc)
	default:
	}

	return fmt.Sprintf("data:image/png;base64,%s", enc)
}

// cleanUrl converts whitespace in urls to %20
func cleanUrl(s string) string {
	return strings.Replace(s, " ", "%20", -1)
}
