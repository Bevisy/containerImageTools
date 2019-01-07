package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func listDockerImagesTags(str string) string {
	url := fmt.Sprintf("https://registry.hub.docker.com/v2/repositories/library/%s/tags?page_size=1024", str)

	return url
}

func main() {
	resp, err := http.Get(listDockerImagesTags("nginx"))
	if err != nil {
		log.Fatalln("request failed! ", err)
	}
	defer resp.Body.Close()
	s, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("read failed! ", err)
	}
	fmt.Println(string(s))
}
