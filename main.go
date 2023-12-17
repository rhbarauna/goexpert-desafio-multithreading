package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	cep := "89216310"
	viaCepChannel := make(chan string)
	brasilApiChannel := make(chan string)

	go buscaBrasilAPI(cep, brasilApiChannel)
	go buscaViaCep(cep, viaCepChannel)

	select {
	case viaBrasil := <-brasilApiChannel:
		printResponse("BrasilAPI", viaBrasil)
	case viaCep := <-viaCepChannel:
		printResponse("ViaCep", viaCep)
	case <-time.After(time.Second * 13):
		log.Println("Tempo de execução expirado")
	}
}

func buscaViaCep(reqCep string, channel chan string) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://viacep.com.br/ws/%s/json/", reqCep), nil)
	if err != nil {
		log.Printf("Falha ao montar a requisição ViaCep. %s\n", err.Error())
		return
	}

	response, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Printf("Falha ao executar a requisição ViaCep. %s\n", err.Error())
		return
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Falha ao ler o body ViaCep. %s\n", err.Error())
		return
	}

	channel <- string(body)
}

func buscaBrasilAPI(reqCep string, channel chan string) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", reqCep), nil)
	if err != nil {
		log.Printf("Falha ao montar o request BrasilAPI. %s\n", err.Error())
		return
	}
	response, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Printf("Falha ao executar o request BrasilAPI. %s\n", err.Error())
		return
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Falha ao ler o body BrasilAPI. %s\n", err.Error())
		return
	}
	channel <- string(body)
}

func printResponse(provider string, response string) {
	log.Printf("Responsta recebida pelo %s: %v\n", provider, response)
}
