package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/sheets/v4"
)

var spreadsheetId = ""

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func main() {
	// Main settings
	// Open our jsonFile
	jsonFile, err := os.Open("conf.json")
	if err != nil {
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var f interface{}
	err = json.Unmarshal(byteValue, &f)
	if err != nil {
		log.Fatalf("%s", err)
	}
	m := f.(map[string]interface{})
	for k, v := range m {
		spreadsheetId = fmt.Sprintf("%s", v)
		_ = k // This removes the error we get because k isn't used.
	}

	// Google credentials
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)
	sheetsrv, serr := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets Client %v", err)
	}

	writeRange := "A1" // Cell to write in
	var vr sheets.ValueRange
	myval := []interface{}{"OK!!!"}
	vr.Values = append(vr.Values, myval)
	_, serr = sheetsrv.Spreadsheets.Values.Update(spreadsheetId, writeRange, &vr).ValueInputOption("RAW").Do()
	if serr != nil {
		log.Fatalf("Unable to retrieve data from sheet. %v", serr)
	} else {
		fmt.Println("Successfully wrote to spreadsheet.")
	}

	// Append
	range2 := "A1"
	// This appends two rows to the xls
	values := [][]interface{}{[]interface{}{"b2", "c2", "d2"}, []interface{}{"b3", "c3", "d3"}}

	valueInputOption := "USER_ENTERED"
	rb := &sheets.ValueRange{
		MajorDimension: "ROWS",
		Values:         values,
	}
	resp, err := sheetsrv.Spreadsheets.Values.Append(spreadsheetId, range2, rb).ValueInputOption(valueInputOption).Do()
	str := fmt.Sprint(resp)
	fmt.Println("Response to append: " + str)
}
