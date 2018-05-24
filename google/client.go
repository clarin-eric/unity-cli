package google

import (
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
	"fmt"
	"os"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"clarin/unity-cli/report"
	"sort"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
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
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	json.NewEncoder(f).Encode(token)
}

func Write(r *report.Report) {
	ctx := context.Background()

	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved client_secret.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	/**
	 * Create a new spreadsheet for this report
	 */
	date := time.Now().Format("2006-01-02")
	props := &sheets.SpreadsheetProperties{
		Title: fmt.Sprintf("CLARIN IDM Report - %s", date),
	}

	sheet_list := []*sheets.Sheet{
		&sheets.Sheet{Properties: &sheets.SheetProperties{Title: "General"}},
		&sheets.Sheet{Properties: &sheets.SheetProperties{Title: "Country"}},
		&sheets.Sheet{Properties: &sheets.SheetProperties{Title: "Domain"}},
		&sheets.Sheet{Properties: &sheets.SheetProperties{Title: "Purpose"}},
		&sheets.Sheet{Properties: &sheets.SheetProperties{Title: "Last Auth"}},
	}

	rb := &sheets.Spreadsheet{
		Properties: props,
		Sheets: sheet_list,
	}

	_, err = srv.Spreadsheets.Create(rb).Context(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}/*

	var vr sheets.ValueRange
	vr.Values = append(vr.Values, []interface{}{"Accounts", "Authenticated", "Not Authenticated"})
	vr.Values = append(vr.Values, []interface{}{fmt.Sprintf("%d", r.General.Num_accounts), fmt.Sprintf("%d", r.General.Num_authenticated_accounts), fmt.Sprintf("%d", r.General.Num_accounts-r.General.Num_authenticated_accounts)})
	writeDataToSheet(srv, resp.SpreadsheetId, "General", &vr)

	writeDataToSheet(srv, resp.SpreadsheetId, "Country", stringMaptoValueRange(r.Countries))
	writeDataToSheet(srv, resp.SpreadsheetId, "Domain", stringMaptoValueRange(r.Domains))
	writeDataToSheet(srv, resp.SpreadsheetId, "Purpose", stringMaptoValueRange(r.Purposes))
	writeDataToSheet(srv, resp.SpreadsheetId, "Last Auth", intMaptoValueRange(r.Last_auths))
	*/
}

func stringMaptoValueRange(m map[string]int64) (*sheets.ValueRange) {
	var vrCountry sheets.ValueRange
	vrCountry.Values = append(vrCountry.Values, []interface{}{"Key", "Value"})
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		val := m[key]
		vrCountry.Values = append(vrCountry.Values, []interface{}{key, val})
	}
	return &vrCountry
}

func intMaptoValueRange(m map[int]int64) (*sheets.ValueRange) {
	var vrCountry sheets.ValueRange
	vrCountry.Values = append(vrCountry.Values, []interface{}{"Key", "Value"})
	var keys []int
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, key := range keys {
		val := m[key]
		vrCountry.Values = append(vrCountry.Values, []interface{}{key, val})
	}
	return &vrCountry
}

func writeDataToSheet(srv *sheets.Service, spreadsheetId, sheetName string, vr *sheets.ValueRange) {
	writeRange := fmt.Sprintf("%s!A1", sheetName)

	_, err := srv.Spreadsheets.Values.Update(spreadsheetId, writeRange, vr).ValueInputOption("RAW").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet. %v", err)
	}
}