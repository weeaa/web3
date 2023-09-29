package twitter

import (
	"log"
	"testing"
)

func TestFetchNitter(t *testing.T) {
	username := "weea_a"
	log.Println(FetchNitter(username))
}
