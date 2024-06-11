package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"time"
)

func main() {
	args := os.Args[1:]
	regex := regexp.MustCompile(`\D`)
	postalCode := regex.ReplaceAllString(args[0], "")

	fmt.Println(postalCode)

	viaCepChannel := make(chan string)
	brasilAPIChannel := make(chan string)

	go viaCep(postalCode, viaCepChannel)
	go brasilAPI(postalCode, brasilAPIChannel)

	select {
	case viaCep := <-viaCepChannel:
		fmt.Printf("Via Cep: %s\n", viaCep)
	case brasilAPI := <-brasilAPIChannel:
		fmt.Printf("Brasil API: %s\n", brasilAPI)
	case <-time.After(time.Second):
		fmt.Println("Timed out")
	}
}

func viaCep(postalCode string, channel chan<- string) {
	response, err := http.Get("https://viacep.com.br/ws/" + postalCode + "/json/")
	if err != nil {
		return
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}
	channel <- string(body)
}

func brasilAPI(postalCode string, channel chan<- string) {
	response, err := http.Get("https://brasilapi.com.br/api/cep/v1/" + postalCode)
	if err != nil {
		return
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}
	channel <- string(body)
}
