package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {

	var code string
	flag.StringVar(&code, "c", "", "the code youre supposed to enter on menti.com")

	var wordfile string
	flag.StringVar(&wordfile, "w", "words.txt", "file that contains the words to be added")

	var repeat int
	flag.IntVar(&repeat, "r", 1, "how often to loop through the file")

	flag.Parse()

	var words_file []string

	f, _ := os.Open(wordfile)
	defer f.Close()
	var scanner = bufio.NewScanner(f)

	for scanner.Scan() {
		words_file = append(words_file, scanner.Text())
	}

	var words_vote []string

	min := func(a, b int) int {
		if a <= b {
			return a
		}
		return b
	}

	var vote_key = getVoteKey(code)
	var public_key, entries = getPublicKeyAndEntries(vote_key)

	for i := 0; i < len(words_file); i += entries {
		batch := words_file[i:min(i+entries, len(words_file))]
		words_vote = append(words_vote, strings.Join(batch, " "))
	}

	for i := 0; i < repeat; i++ {
		for _, words := range words_vote {
			var identifier = getIdentifier()
			makeVote(vote_key, public_key, identifier, words)
		}
	}
}

func getVoteKey(code string) string {

	client := &http.Client{}
	request, _ := http.NewRequest("GET", fmt.Sprintf("https://www.menti.com/core/vote-ids/%s/series", code), nil)

	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.82 Safari/537.36")
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Referer", "https://www.menti.com/")
	request.Header.Set("Authority", "www.menti.com")

	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Failed making request")
		os.Exit(0)
	}

	data, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	var voteKeyJson map[string]string

	json.Unmarshal([]byte(data), &voteKeyJson)

	return voteKeyJson["vote_key"]

}

func getPublicKeyAndEntries(vote_key string) (string, int) {

	client := &http.Client{}
	request, _ := http.NewRequest("GET", fmt.Sprintf("https://www.menti.com/core/vote-keys/%s/series", vote_key), nil)

	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.82 Safari/537.36")
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Referer", fmt.Sprintf("https://www.menti.com/%s", vote_key))
	request.Header.Set("Authority", "www.menti.com")

	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Failed making request")
		os.Exit(0)
	}

	data, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	//fmt.Println(string(data))
	var publicKeyJson map[string][]map[string]string

	json.Unmarshal([]byte(data), &publicKeyJson)
	entries, _ := strconv.Atoi(publicKeyJson["questions"][0]["max_nb_words"])
	publicKey := publicKeyJson["questions"][0]["public_key"]
	return publicKey, entries
}

func getIdentifier() string {
	client := &http.Client{}
	request, _ := http.NewRequest("POST", "https://www.menti.com/core/identifier", nil)

	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.82 Safari/537.36")

	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Failed making request")
		os.Exit(0)
	}

	data, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	var identifierJson map[string]string

	json.Unmarshal([]byte(data), &identifierJson)

	return identifierJson["identifier"]
}

func makeVote(vote_key string, public_key string, identifier string, words string) {

	var jsonStr = []byte(fmt.Sprintf(`{"question_type":"wordcloud","vote":"%s"}`, words))

	client := &http.Client{}
	request, _ := http.NewRequest("POST", fmt.Sprintf("https://www.menti.com/core/votes/%s", public_key), bytes.NewBuffer(jsonStr))

	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.82 Safari/537.36")
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Referer", fmt.Sprintf("https://www.menti.com/%s", vote_key))
	request.Header.Set("Authority", "www.menti.com")
	request.Header.Set("X-Identifier", identifier)

	_, err := client.Do(request)
	if err != nil {
		fmt.Println("[-] Failed adding words")
	} else {
		fmt.Printf("[+] Added words to the wordcloud (%s)\n", words)
	}
}
