package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	gui "github.com/grupawp/warships-lightgui/v2"
)

//zaimmplementować statystyke skuteczności strzałów

var gameProperties GameProperties

func customBoard() ([]string, error) {
	return []string{
		"A1",
		"A3",
		"B9",
		"C7",
		"D1",
		"D2",
		"D3",
		"D4",
		"D7",
		"E7",
		"F1",
		"F2",
		"F3",
		"F5",
		"G5",
		"G8",
		"G9",
		"I4",
		"J4",
		"J8",
	}, nil
}

func Fire(coord string) (string, error) {

	Loop := true
	no_tries := 0
	for ok := true; ok; ok = Loop {
		client := &DefaultHTTPClient{}

		GameBoard := map[string]interface{}{
			"coord": coord,
		}
		b, err := json.Marshal(GameBoard)
		if err != nil {
			fmt.Println(err)
			return "", fmt.Errorf("failed to marshall: ", err)
		}
		jsonBody := []byte(b)
		fmt.Println("body////")
		fmt.Printf((string(jsonBody)))

		getHeaders := map[string]string{
			"X-Auth-Token": gameProperties.Token,
		}

		postResponse, err := client.Post("https://go-pjatk-server.fly.dev/api/game/fire", "application/json", jsonBody, getHeaders)
		if err != nil {
			fmt.Println("POST request failed", err)
			fmt.Errorf("post request failed", err)
		}
		if postResponse.StatusCode != http.StatusOK {
			fmt.Printf("unexpected status: %d, %s", postResponse.StatusCode, postResponse.Header.Get("message"))
			fmt.Errorf("unexpected status: %d, %s", postResponse.StatusCode, postResponse.Header.Get("message"))
		}
		if postResponse.StatusCode == http.StatusOK {
			Loop = false
		}

		defer postResponse.Body.Close()
		fmt.Println("POST ResponseBody:")
		fmt.Println(string(postResponse.Header.Get("X-Auth-Token")))
		token := postResponse.Header.Get("X-Auth-Token")
		if len(token) == 0 {
			fmt.Errorf("cannot obtain token")
		}
		// Reading response body
		postResponseBody, err := io.ReadAll(postResponse.Body)
		if err != nil {
			// fmt.Println("Failed to read POST response body:", err)
			fmt.Errorf("Failed to read POST response body", err)
		}
		fmt.Println("POST Response:")
		fmt.Println(string(postResponseBody))
		gameProperties.Token = token

		var data map[string]interface{}

		err = json.Unmarshal([]byte(postResponseBody), &data)
		if err != nil {
			return "", fmt.Errorf("failed to unmarshal response")
		}
		result := fmt.Sprintf("%s", data["result"])
		fmt.Println(result)

		if no_tries == 3 {
			Loop = false
			return "", fmt.Errorf("couldn't perform fire request")
		}
		no_tries += 1

		return result, fmt.Errorf("couldn't perform fire request")

	}
	return "", nil

}

func stringToSlice(inp string) []string {
	inp = strings.Replace(inp, "[", "", -1)
	inp = strings.Replace(inp, "]", "", -1)
	s := strings.Split(inp, " ")
	return s
}
func Board() ([]string, error) {

	Loop := true
	no_tries := 0
	for ok := true; ok; ok = Loop {
		client := &DefaultHTTPClient{}

		getHeaders := map[string]string{
			"X-Auth-Token": gameProperties.Token,
		}

		//////
		resp, err := client.Get("https://go-pjatk-server.fly.dev/api/game/board", getHeaders)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("unexpected status: %d, %s", resp.StatusCode, resp.Header.Get("message"))

			return nil, fmt.Errorf("unexpected status: %d, %s", resp.StatusCode, resp.Header.Get("message"))
		}
		if resp.StatusCode == http.StatusOK {
			Loop = false
		}
		var data map[string]interface{}

		err = json.Unmarshal([]byte(body), &data)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal response")
		}
		key := fmt.Sprintf("%s", data["board"])
		result := stringToSlice(key)
		if len(result) != 20 {
			fmt.Printf("Not enough pieces")
			return nil, fmt.Errorf("Not enough pieces")
		}
		fmt.Println("%d", len(result))

		if no_tries == 3 {
			Loop = false
			return nil, fmt.Errorf("Couldn't retrieve board")
		}
		no_tries += 1

		return result, nil

	}
	return nil, fmt.Errorf("Couldn't retrieve board")

}

//func Fire(coord string) (string, error)

type GameProperties struct {
	Token string
	Board []string
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type StatusResponse struct {
	StatusCode int
	Body       map[string]interface{}
}

// http
type DefaultHTTPClient struct{}

func (c *DefaultHTTPClient) Get(url string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return http.DefaultClient.Do(req)
}

func (c *DefaultHTTPClient) Post(url string, bodyType string, body []byte, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", bodyType)
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return http.DefaultClient.Do(req)
}

func InitGame() error {

	Loop := true
	no_tries := 0
	for ok := true; ok; ok = Loop {

		client := &DefaultHTTPClient{}
		getHeaders := map[string]string{
			"accept": "application/json",
		}
		GameBoard := map[string]interface{}{
			"coords":      gameProperties.Board,
			"desc":        "Pierwsza gra",
			"nick":        "BetonoJanusz",
			"target_nick": "",
			"wpbot":       true,
		}
		b, err := json.Marshal(GameBoard)
		if err != nil {
			fmt.Println(err)
			fmt.Errorf("failed to marshall: ", err)
		}
		jsonBody := []byte(b)

		postResponse, err := client.Post("https://go-pjatk-server.fly.dev/api/game", "application/json", jsonBody, getHeaders)
		if err != nil {
			fmt.Println("POST request failed", err)
			fmt.Errorf("post request failed", err)
		}
		if postResponse.StatusCode != http.StatusOK {
			fmt.Printf("unexpected status: %d, %s", postResponse.StatusCode, postResponse.Header.Get("message"))
			fmt.Errorf("unexpected status: %d, %s", postResponse.StatusCode, postResponse.Header.Get("message"))
		}
		if postResponse.StatusCode == http.StatusOK {
			Loop = false
		}

		defer postResponse.Body.Close()
		fmt.Println("POST ResponseBody:")
		fmt.Println(string(postResponse.Header.Get("X-Auth-Token")))
		token := postResponse.Header.Get("X-Auth-Token")
		if len(token) == 0 {
			fmt.Errorf("cannot obtain token")
		}
		// Reading response body
		postResponseBody, err := io.ReadAll(postResponse.Body)
		if err != nil {
			// fmt.Println("Failed to read POST response body:", err)
			fmt.Errorf("Failed to read POST response body", err)
		}
		fmt.Println("POST Response:")
		fmt.Println(string(postResponseBody))
		gameProperties.Token = token

		if no_tries == 3 {
			Loop = false
			return fmt.Errorf("couldn't initialize game")
		}
		no_tries += 1
	}
	return nil

}

func getStats() (*StatusResponse, error) {

	Loop := true
	no_tries := 0
	for ok := true; ok; ok = Loop {

		client := &DefaultHTTPClient{}

		getHeaders := map[string]string{
			"X-Auth-Token": gameProperties.Token,
		}

		//////
		resp, err := client.Get("https://go-pjatk-server.fly.dev/api/stats", getHeaders)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var data map[string]interface{}

		err = json.Unmarshal([]byte(body), &data)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal response")
		}

		if no_tries == 3 {
			Loop = false
			return nil, fmt.Errorf("couldn't retrieve lobby data")
		}
		no_tries += 1

		if resp.StatusCode == 200 {
			Loop = false
		}
		return &StatusResponse{
			StatusCode: resp.StatusCode,
			Body:       data,
		}, nil

	}
	return nil, fmt.Errorf("couldn't retrieve lobby data")
}

func getLobby() (*StatusResponse, error) {

	Loop := true
	no_tries := 0
	for ok := true; ok; ok = Loop {

		client := &DefaultHTTPClient{}

		getHeaders := map[string]string{}

		//////
		resp, err := client.Get("https://go-pjatk-server.fly.dev/api/lobby", getHeaders)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var data map[string]interface{}

		err = json.Unmarshal([]byte(body), &data)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal response")
		}

		if no_tries == 3 {
			Loop = false
			return nil, fmt.Errorf("couldn't retrieve lobby data")
		}
		no_tries += 1

		if resp.StatusCode == 200 {
			Loop = false
		}
		return &StatusResponse{
			StatusCode: resp.StatusCode,
			Body:       data,
		}, nil

	}
	return nil, fmt.Errorf("couldn't retrieve lobby data")
}

func Status() (*StatusResponse, error) {

	Loop := true
	no_tries := 0
	for ok := true; ok; ok = Loop {

		client := &DefaultHTTPClient{}

		getHeaders := map[string]string{
			"X-Auth-Token": gameProperties.Token,
		}

		//////
		resp, err := client.Get("https://go-pjatk-server.fly.dev/api/game", getHeaders)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var data map[string]interface{}

		err = json.Unmarshal([]byte(body), &data)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal response")
		}

		if no_tries == 3 {
			Loop = false
			return nil, fmt.Errorf("couldn't initialize game")
		}

		if resp.StatusCode == http.StatusOK {
			Loop = false
		}
		no_tries += 1
		return &StatusResponse{
			StatusCode: resp.StatusCode,
			Body:       data,
		}, nil

	}
	return nil, fmt.Errorf("couldn't get status")

}
func getInput() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(">")
	input, _ := reader.ReadString('\n')

	cleanedInput := strings.ReplaceAll(input, "\n", "")
	cleanedInput = strings.ReplaceAll(cleanedInput, "\r", "")
	return cleanedInput
}

func isValidFormat(input string) bool {
	validFormat := regexp.MustCompile(`^[a-j](?:10|[1-9])$`)
	return validFormat.MatchString(input)
}

func getCoords() (string, error) {
	fmt.Println("Podaj kordynaty")
	waitingLoop := true
	for ok := true; ok; ok = waitingLoop {
		userInput := strings.ToLower(getInput())
		if isValidFormat(userInput) {
			return userInput, nil
			waitingLoop = false
		} else {
			fmt.Println("Błędne koordynaty, spróbuj jeszcze raz")
		}
	}

	return "", nil
}

func main() {
	gameProperties.Board, _ = customBoard()
	err := InitGame()
	if err != nil {
		panic(err)
	}

	fmt.Println("token: ", gameProperties.Token)
	GameStatus, err := Status()
	if err != nil {
		panic(err)
	}

	///
	for key, value := range GameStatus.Body {
		fmt.Printf("%s: %v\n", key, value)
	}
	fmt.Printf("////////////")
	key, value := GameStatus.Body["game_status"]
	fmt.Printf("%s: %v\n", key, value)

	menuLoop := true
	for ok := true; ok; ok = menuLoop {
		fmt.Printf("Wybierz opcję: \n|start \n|exit \n|lobby\n|stats\n")
		userInp := getInput()
		fmt.Println(userInp)
		switch userInp {
		case "start":
			menuLoop = false
		case "exit":
			return
		case "lobby":
			lobby, err := getLobby()
			if err != nil {
				panic(err)
			}
			fmt.Println("Lobby: ")
			for key, value := range lobby.Body {
				fmt.Printf("%s: %v\n", key, value)
			}
		case "stats":
			stats, err := getStats()
			if err != nil {
				panic(err)
			}
			fmt.Println("stats: ")
			for key, value := range stats.Body {
				fmt.Printf("%s: %v\n", key, value)
			}
		default:
			fmt.Println("Spróbuj jeszcze raz")
		}
	}

	//gameStatus := GameStatus.Body["game_status"]
	waitingLoop := true
	for ok := true; ok; ok = waitingLoop {
		GameStatus, err = Status()
		if err != nil {
			panic(err)
		}
		if GameStatus.Body["game_status"] == "no_game" || GameStatus.Body["game_status"] == "waiting_wpbot" {

			for key, value := range GameStatus.Body {
				fmt.Printf("%s: %v\n", key, value)
			}

		} else {
			waitingLoop = false
		}
		time.Sleep(1 * time.Second)
	}
	//

	board := gui.New(gui.NewConfig())
	board.Display()
	gameLoop := true
	for ok := true; ok; ok = gameLoop {

		GameStatus, err = Status()
		if err != nil {
			panic(err)
		}
		if GameStatus.Body["should_fire"] == true {
			coord, _ := getCoords()
			Fire(coord)
			fmt.Println(GameStatus.Body["timer"])
		}

		fmt.Println(GameStatus.Body)
		//game status znika za 2 wykonaniem z jakiegos powodu

		time.Sleep(1 * time.Second)

	}

	key, err = Board()
	str := key
	fmt.Println(str)
	//Board()

}
