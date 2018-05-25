package report

import (
	"fmt"
	"os"
	"time"
	"clarin/unity-cli/api"
	"strings"
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
	Created string `json:"created"`
	Updated string `json:"updated"`
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
	Created string `json:"created"`
	Updated string `json:"updated"`
	Domain string `json:"domain"`
	Country string `json:"county"`
	LastAuthn string `json:"last_auth"`
	DaysSinceLastAuthentication int `json:"days_since_last_auth"`
	Purpose string `json:"purpose"`
	Motivation string `json:"motivation"`
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
		email, created, updated := getEmail(e)
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
			Created: r.timeToString(created),
			Updated: r.timeToString(updated),
		 })

		r.personal_report_entities = append(r.personal_report_entities, PersonalReportEntity{
			FullName: full_name,
			Member: member,
			Email: email,
			Created: r.timeToString(created),
			Updated: r.timeToString(updated),
			Domain: domain,
			Country: country,
			LastAuthn: last_auth,
			DaysSinceLastAuthentication: days_since_last_auth,
			Purpose: purpose,
			Motivation: motivation,
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
}

func (r *Report) timeToString(t *time.Time) (string) {
	result := "-"
	if t != nil {
		result = t.String()
	}
	return result
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

func getEmail(e api.Entity) (string, *time.Time, *time.Time) {
	email_identity := "unkown"
	var created *time.Time
	var updated *time.Time

	for _, id := range e.Identities {
		if id.TypeId == "email" {
			email_identity = id.Value
			tCreated := time.Unix(0, id.CreationTs*1000000)
			tUpdated := time.Unix(0, id.UpdateTs*1000000)
			created = &tCreated
			updated = &tUpdated
		}
	}
	return email_identity, created, updated
}

func (r *Report) getEmailDomain(email string) (string) {
	domain := "unkown"
	if email != "unkown" {
		domain = strings.Split(email, "@")[1]
	}
	return domain
}

func (r *Report) GetReport(kind string) (interface{}) {
	switch kind {
	case "anonymous": return r.anonymous_report_entities
	case "personal": return r.personal_report_entities
	}
	return nil
}

func (r *Report) GetReportAsArray(kind string) ([][]string) {
	if kind == "personal" {
		header := PersonalReportEntity{}
		data := [][]string{header.StringArrayHeader()}
		for _, e := range r.personal_report_entities {
			data = append(data, e.AsStringArray())
		}
		return data
	}

	if kind == "anonymous" {
		header := AnonymousReportEntity{}
		data := [][]string{header.StringArrayHeader()}
		for _, e := range r.anonymous_report_entities {
			data = append(data, e.AsStringArray())
		}
		return data
	}

	return nil
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

	dir := r.UserHomeDir()
	basename := "unity_reporting"
	subdir := time.Now().Format("2006-01-02")
	base := fmt.Sprintf("%s/%s/%s", dir, basename, subdir)

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

