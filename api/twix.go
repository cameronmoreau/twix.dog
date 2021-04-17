package handler

import (
	"encoding/xml"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

type ListBucketResult struct {
	XMLName  xml.Name  `xml:"ListBucketResult"`
	Name     string    `xml:"Name"`
	Contents []Content `xml:"Contents"`
}

type Content struct {
	XMLName xml.Name `xml:"Contents"`
	Key     string   `xml:"Key"`
}

const (
	S3_BUCKET_URL = "https://cameron-media.s3.us-east-1.amazonaws.com/"
	S3_PREFIX     = "twix"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	data, err := fetchMedia()
	if err != nil {
		panic(err)
	}

	image, err := http.Get(randomMediaUrl(data))
	if err != nil {
		panic(err)
	}

	defer image.Body.Close()
	io.Copy(w, image.Body)
}

func randomMediaUrl(data *ListBucketResult) string {
	rand.Seed(time.Now().UnixNano())
	key := rand.Intn(len(data.Contents)-1) + 1
	return S3_BUCKET_URL + data.Contents[key].Key
}

func fetchMedia() (data *ListBucketResult, err error) {
	resp, err := http.Get(S3_BUCKET_URL + "?list-type=2&prefix=" + S3_PREFIX)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var results ListBucketResult
	xml.Unmarshal(body, &results)

	return &results, nil
}
