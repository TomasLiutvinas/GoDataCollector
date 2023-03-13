package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/chromedp/chromedp"
)

type Sets struct {
	Sets []CardSet `json:"Sets"`
}

type CardSet struct {
	Title string `json:"Title"`
	Value string `json:"Value"`
}

type Credentials struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}

func main() {
	//ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithDebugf(log.Printf)) // LOGGING
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	executeLogin(ctx)

	var sets Sets

	// Open our jsonFile
	jsonFile, err := os.Open("sets.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &sets)

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	for i := 0; i < len(sets.Sets); i++ {
		set := sets.Sets[i].Value
		body := executePageCount(ctx, set)
		data := []byte(body)
		ioutil.WriteFile(fmt.Sprintf("%s.json", set), data, 0)
	}
}

func readSetsFromFile(fileName string, data *Sets) {

	// Open our jsonFile
	jsonFile, err := os.Open(fileName)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &data)

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
}

func executePageCount(ctx context.Context, set string) string {

	// move iterating here
	var res string
	err := chromedp.Run(ctx, countPage(&res, ctx, set))

	if err != nil {
		log.Fatal(err)
	}

	return res
}

func countPage(res *string, ctx context.Context, set string) chromedp.Tasks {
	var setUrl = fmt.Sprintf("https://fabdb.net/api/collection?set=%s&per_page=500", set)

	return chromedp.Tasks{
		chromedp.Navigate(setUrl),
		chromedp.Text("pre", res),
	}
}

func executeLogin(ctx context.Context) {
	var res string
	err := chromedp.Run(ctx, login(&res))

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("got: `%s`", strings.TrimSpace(res))
}

func readCredentialsFromFile(fileName string, creds *Credentials) {

	// Open our jsonFile
	jsonFile, err := os.Open(fileName)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &creds)

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
}

func login(res *string) chromedp.Tasks {
	var fabLoginUrl = `https://fabdb.net/login`
	var emailSelector = `//input[@placeholder="Email address"]`
	var pwSelector = `//input[@placeholder="Enter password"]`
	var nextButtonSelector = `//button[text()=" Next "]`
	var firstDeckSelector = `//div[@class="main-body"]/div/div/table/tbody/tr[1]/td[1]/a/div/span`
	var creds Credentials
	readCredentialsFromFile("credentials.credentials", &creds)

	return chromedp.Tasks{
		chromedp.Navigate(fabLoginUrl),
		chromedp.WaitVisible(emailSelector),
		chromedp.SendKeys(emailSelector, creds.Username),
		chromedp.Click(nextButtonSelector),
		chromedp.WaitVisible(pwSelector),
		chromedp.SendKeys(pwSelector, creds.Password),
		chromedp.Click(nextButtonSelector),
		chromedp.Text(firstDeckSelector, res),
	}
}
