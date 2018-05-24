package report

import (
	"fmt"
	"os"
	"encoding/csv"
	"time"
	"clarin/unity-client/api"
	"strings"
	"encoding/json"
	"runtime"
)

type Report struct {
	anonymous_report_entities []AnonymousReportEntity
	personal_report_entities []PersonalReportEntity
}

type AnonymousReportEntity struct {
	Domain string `json:"domain"`
	Country string `json:"county"`
	LastAuthn string `json:"last_auth"`
	DaysSinceLastAuthentication int `json:"days_since_last_auth"`
	Purpose string `json:purpose`
}

func (e *AnonymousReportEntity) StringArrayHeader() ([]string) {
	return []string {"domain", "country", "last_authn", "days_since_last_authn", "purpose"}
}
func (e *AnonymousReportEntity) AsStringArray() ([]string) {
	return []string {e.Domain, e.Country, e.LastAuthn, fmt.Sprintf("%d", e.DaysSinceLastAuthentication), e.Purpose}
}

type PersonalReportEntity struct {
	FullName string `json:full_name`
	Member string `json:"member"`
	Email string `json:"email"'`
	Domain string `json:"domain"`
	Country string `json:"county"`
	LastAuthn string `json:"last_auth"`
	DaysSinceLastAuthentication int `json:"days_since_last_auth"`
	Purpose string `json:purpose`
}

func (e *PersonalReportEntity) StringArrayHeader() ([]string) {
	return []string {"full_name", "member", "email", "domain", "country", "last_authn", "days_since_last_authn", "purpose"}
}
func (e *PersonalReportEntity) AsStringArray() ([]string) {
	return []string {e.FullName, e.Member, e.Email, e.Domain, e.Country, e.LastAuthn, fmt.Sprintf("%d", e.DaysSinceLastAuthentication), e.Purpose}
}


func (r *Report) reset() {
	r.anonymous_report_entities = []AnonymousReportEntity{}
	r.personal_report_entities = []PersonalReportEntity{}
}

/**
 * Not processed:
 *
 *		e.GetAttributeValuesAsString("clarin-lr-list"),
 *		e.GetAttributeValuesAsString("clarin-motivation"),
 *		e.GetAttributeValuesAsString("cn"),
 *
 */
func (r *Report) Compute(entities []api.Entity) {
	r.reset()

	no_full_name := []string{}
	no_member := []string{}
	no_motivation := []string{}
	for _, e := range entities {

		//personal attributes
		full_name := e.GetAttributeValuesAsString("clarin-full-name")
		member := e.GetAttributeValuesAsString("member")
		email := getEmail(e)
		motivation := e.GetAttributeValuesAsString("clarin-motivation")
		//anonymous attributes
		domain := r.getEmailDomain(email)
		country := e.GetAttributeValuesAsString("clarin-country")
		last_auth := e.GetAttributeValuesAsString("sys:LastAuthentication")
		days_since_last_auth := r.getDaysSinceLastAuth(last_auth)
		purpose := e.GetAttributeValuesAsString("clarin-purpose")

		r.anonymous_report_entities = append(r.anonymous_report_entities, AnonymousReportEntity{
		 	Country: country,
		 	Domain: domain,
		 	LastAuthn: last_auth,
		 	DaysSinceLastAuthentication: days_since_last_auth,
		 	Purpose: purpose,
		 })

		r.personal_report_entities = append(r.personal_report_entities, PersonalReportEntity{
			FullName: full_name,
			Member: member,
			Email: email,
			Domain: domain,
			Country: country,
			LastAuthn: last_auth,
			DaysSinceLastAuthentication: days_since_last_auth,
			Purpose: purpose,
		})

		if full_name == "" {
			no_full_name = append(no_full_name, fmt.Sprintf("%s\t%d", email, e.Id))
		}
		if member == "" {
			no_member = append(no_member, fmt.Sprintf("%s\t%d", email, e.Id))
		}
		if motivation == "" {
			no_motivation = append(no_motivation, fmt.Sprintf("%s\t%d", email, e.Id))
		}
 	}

	fmt.Printf("No full name value:\n")
 	for _, v := range no_full_name {
 		fmt.Printf("%s\n", v)
	}

	fmt.Printf("No member value:\n")
	for _, v := range no_member {
		fmt.Printf("%s\n", v)
	}

	fmt.Printf("No motivation value:\n")
	for _, v := range no_motivation {
		fmt.Printf("%s\n", v)
	}
}

func (r *Report) getDaysSinceLastAuth(last_auth string) (int) {
	layout := "2006-01-02T15:04:05" //2018-03-05T09:45:21
	days_since_last_auth := -1
	if len(last_auth) > 0 {
		ts_last_auth, err := time.Parse(layout, last_auth)
		if err != nil {
			fmt.Printf("Failed to parse timestamp")
		} else {
			d_since_last_auth := time.Since(ts_last_auth)
			days_since_last_auth = int(d_since_last_auth.Hours() / (24.0*7))
		}
	}
	return days_since_last_auth
}

func getEmail(e api.Entity) (string) {
	email_identity := "unkown"
	for _, id := range e.Identities {
		if id.TypeId == "email" {
			email_identity = id.Value
		}
	}
	return email_identity
}

func (r *Report) getEmailDomain(email string) (string) {
	domain := "unkown"
	if email != "unkown" {
		domain = strings.Split(email, "@")[1]
	}
	return domain
}

func (r *Report) Write(kind, output_format string) {

	if kind == "anonymous" || kind == "both" {
		if output_format == "csv" || output_format == "tsv" {
			tsv := output_format == "tsv"
			r.writeAnonymouseSeparated(tsv)
		} else if output_format == "json" {

		} else if output_format == "google" {

		} else {

		}
	}

	if kind == "personal" || kind == "both" {
		if output_format == "csv" || output_format == "tsv" {
			tsv := output_format == "tsv"
			r.writePersonalSeparated(tsv)
		} else if output_format == "json" {

		} else if output_format == "google" {

		} else {

		}
	}
}

func (r *Report) writeSeparatedValues(kind, ext, sep string) {

}

func (r *Report) PrettyPrint() {
	/*
	fmt.Printf("#accounts : %d\n", r.General.Num_accounts)
	fmt.Printf("   #active: %d\n", r.General.Num_authenticated_accounts)
	r.printMap(r.Domains, "Domains")
	r.printMap(r.Countries, "Country")
	r.printMapInt64(r.Last_auths, "Last authenticated")
	*/
}

func (r *Report) addToMap(m map[string]int64, key string) (map[string]int64) {
	if _, ok := m[key]; !ok {
		m[key] = 1
	} else {
		m[key] += 1
	}
	return m
}

func (r *Report) PrettyAsJson() {
	json_bytes, err := json.MarshalIndent(r.anonymous_report_entities, "", "  ")
	if err != nil {
		fmt.Printf("Failed to marshall JSON. Error: %v\n", err)
	}
	fmt.Printf("%s\n", string(json_bytes))
}

func (r *Report) writeAsJson() {

}

func (r *Report) writePersonalSeparated(tsv bool) {
	ext := "csv"
	if tsv {
		ext = "tsv"
	}

	header := PersonalReportEntity{}
	data := [][]string{	header.StringArrayHeader() }
	for _, e := range r.personal_report_entities {
		data = append(data, e.AsStringArray())
	}
	r.writeCsv(fmt.Sprintf("%s.%s", "personal_report", ext), data, tsv)
}

func (r *Report) writeAnonymouseSeparated(tsv bool) {
	ext := "csv"
	if tsv {
		ext = "tsv"
	}

	header := AnonymousReportEntity{}
	data := [][]string{	header.StringArrayHeader() }
	for _, e := range r.anonymous_report_entities {
		data = append(data, e.AsStringArray())
	}
	r.writeCsv(fmt.Sprintf("%s.%s", "anonymous_report", ext), data, tsv)
}

func (r *Report) writeCsv(filename string, data [][]string, tsv bool) {
	file, err := r.GetFile(filename)
	if err != nil {
		fmt.Printf("Cannot create file: %v\n", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	if tsv {
		writer.Comma = '\t'
	}
	defer writer.Flush()

	for _, value := range data {
		err := writer.Write(value)
		if err != nil {
			fmt.Printf("Cannot write to file: %v\n", err)
		}
	}
}

func (r *Report) UserHomeDir() string {
	env := "HOME"
	if runtime.GOOS == "windows" {
		env = "USERPROFILE"
	} else if runtime.GOOS == "plan9" {
		env = "home"
	}
	return os.Getenv(env)
}

func (r *Report) GetFile(filename string) (*os.File, error) {
	basename := "unity_reporting"
	dir := r.UserHomeDir()
	base := fmt.Sprintf("%s/%s", dir, basename)

	if _, err := os.Stat(base); os.IsNotExist(err) {
		fmt.Printf("Creating %s\n", base)
		os.MkdirAll(base, os.ModePerm);
	}

	path := fmt.Sprintf("%s/%s", base, filename)
	file, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("Cannot create file: %v\n", err)
	}
	return file, nil
}

