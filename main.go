package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	godotenv "github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

type Perintent struct {
	Intent string
	Text   []string
}

func read() []Perintent {

	result := make([]Perintent, 0, 1600)
	file, err := os.Open("data.txt")

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string

	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}

	for _, txtline := range txtlines {

		sliceLine := strings.Split(txtline, "\t")

		if len(sliceLine) == 2 {
			var elemPerintent Perintent
			result = append(result, elemPerintent)
			result[len(result)-1].Intent = sliceLine[1]
		}

		result[len(result)-1].Text = append(result[len(result)-1].Text, sliceLine[0])
	}

	file.Close()

	return result
}

type UtterRequest struct {
	Text     string   `json:"text"`
	Intent   string   `json:"intent"`
	Entities []string `json:"entities"`
	Traits   []string `json:"traits"`
}

func main() {
	client := &http.Client{}

	allintent := read()

	fmt.Println(allintent)

	var bearer = "Bearer " + os.Getenv("WIT_AI_TOKEN")

	for _, intentObj := range allintent {

		postBody, _ := json.Marshal(map[string]string{
			"name": intentObj.Intent,
		})

		responseBody := bytes.NewBuffer(postBody)

		req, err := http.NewRequest("POST", "https://api.wit.ai/intents?v=20200513", responseBody)

		if err != nil {
			log.Fatalf("An Error Occured %v", err)
		}

		req.Header.Add("Authorization", bearer)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)

		if err != nil {
			log.Println("Error on response.\n[ERROR] -", err)
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error while reading the response bytes:", err)
		}

		fmt.Println(string([]byte(body)))

		allUtter := make([]UtterRequest, 0)

		for _, utterance := range intentObj.Text {

			var newUtter UtterRequest
			fake := make([]string, 0)

			newUtter.Entities = fake
			newUtter.Text = utterance
			newUtter.Intent = intentObj.Intent
			newUtter.Traits = fake

			allUtter = append(allUtter, newUtter)
		}

		postBody, err = json.Marshal(allUtter)
		if err != nil {
			log.Fatalf("An Error Occured %v", err)
		}

		responseBody = bytes.NewBuffer(postBody)

		req, err = http.NewRequest("POST", "https://api.wit.ai/utterances?v=20200513", responseBody)

		if err != nil {
			log.Fatalf("An Error Occured %v", err)
		}

		req.Header.Add("Authorization", bearer)
		req.Header.Set("Content-Type", "application/json")

		resp, err = client.Do(req)

		if err != nil {
			log.Println("Error on response.\n[ERROR] -", err)
		}

		defer resp.Body.Close()

		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error while reading the response bytes:", err)
		}

		fmt.Println(string([]byte(body)))
		fmt.Println()

		time.Sleep(1000 * time.Millisecond)
	}

}
