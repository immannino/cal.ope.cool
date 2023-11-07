package cmd

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"cal.ope.cool/internal"
	"cal.ope.cool/pkg/nhl"
	"github.com/joho/godotenv"
	"github.com/oapi-codegen/runtime/types"
)

const (
	teams       = "teams"
	seasonStart = "2023-09-01"
	seasonEnd   = "2024-05-01"
	distDir     = "./docs"
	nhlDir      = "nhl"
)

type Application struct {
	nhlClient  *nhl.ClientWithResponses
	resources  *internal.ResourceData
	debug      bool
	debugDelay time.Duration
}

func Start() {
	must(godotenv.Load())
	client, err := nhl.NewClientWithResponses(os.Getenv("NHL_API"))
	must(err)

	app := NewApplication(client)

	must(app.FetchTeamSchedules(context.TODO()))
}

func NewApplication(nhlClient *nhl.ClientWithResponses) *Application {
	r := &internal.ResourceData{}
	r.Links = []internal.Link{}

	return &Application{nhlClient: nhlClient, debug: false, debugDelay: time.Millisecond * 250, resources: r}
}

func (a *Application) ConvertTeamScheduleToCalendar(schedule *nhl.Schedule, title, filename string) *internal.Calendar {
	c := &internal.Calendar{}
	c.Filename = filename
	c.Title = title
	c.Events = []internal.Event{}

	for _, v := range *schedule.Dates {
		for _, g := range *v.Games {
			c.Events = append(c.Events, ToSchedule(&g))
			if a.debug {
				log.Printf("Date=%s Home=%s Away=%s Location=%s", v.Date.Time.Format("2006-01-02"), *g.Teams.Home.Team.Name, *g.Teams.Away.Team.Name, *g.Venue.Name)
			}
		}
	}

	return c
}

func (a *Application) FetchTeamSchedules(ctx context.Context) error {
	nhlTeamsResponse, err := a.nhlClient.GetTeamsWithResponse(ctx, &nhl.GetTeamsParams{
		Expand: (*nhl.GetTeamsParamsExpand)(paramString([]string{"team.schedule.next"})),
	})
	if err != nil {
		log.Printf("error fetching nhl teams, %v", err)
		return err
	}
	log.Print("Fetched NHL Teams")

	nhlSeason := &internal.Calendar{Title: "2023-2024 NHL Season", Filename: "2023-2024-NHL-Schedule.ics"}
	nhlSeason.Events = []internal.Event{}

	for i, v := range *nhlTeamsResponse.JSON200.Teams {
		log.Printf("Processing (%d/%d): Team=%s ID=%d", i+1, len(*nhlTeamsResponse.JSON200.Teams), *v.Name, *v.Id)
		id := strconv.Itoa(*v.Id)
		start, err := time.Parse("2006-01-02", seasonStart)
		if err != nil {
			log.Printf("error parsing start date, %v", err)
			return err
		}
		end, err := time.Parse("2006-01-02", seasonEnd)
		if err != nil {
			log.Printf("error parsing end date, %v", err)
			return err
		}

		schedule, err := a.nhlClient.GetScheduleWithResponse(ctx, &nhl.GetScheduleParams{
			TeamId:    &id,
			StartDate: &types.Date{Time: start},
			EndDate:   &types.Date{Time: end},
		})
		if err != nil {
			log.Printf("error fetching schedule, %v", err)
			return err
		}
		log.Printf("Fetched Team Schedule, Team=%s, ID=%s, Games=%d", *v.Name, id, len(*schedule.JSON200.Dates))

		cal := a.ConvertTeamScheduleToCalendar(schedule.JSON200, fmt.Sprintf("2023 %s NHL Season", *v.Name), fmt.Sprintf("Team-Schedule-%s.ics", *v.ShortName))
		// Add events to parent calendar
		nhlSeason.Events = append(nhlSeason.Events, cal.Events...)

		a.PersistCal(cal)

		log.Printf("Fetched Team Schedule, Team=%s, ID=%s, Filepath=%s", *v.Name, id, cal.Filename)
		time.Sleep(a.debugDelay)
	}

	err = a.PersistCal(nhlSeason)
	if err != nil {
		log.Printf("error persisting nhlSeason cal, %v", err)
		return err
	}

	a.resources.LastUpdated = time.Now()
	b, err := json.Marshal(a.resources)
	if err != nil {
		log.Printf("error marshaling resource file, %v", err)
		return err
	}

	err = ioutil.WriteFile(fmt.Sprintf("./%s/resources.json", distDir), b, 0644)
	if err != nil {
		log.Printf("error saving resource file, %v", err)
		return err
	}

	return nil
}

func (a *Application) PersistCal(c *internal.Calendar) error {
	calICS := c.ToICS()
	err := ioutil.WriteFile(c.Filename, []byte(calICS), 0644)
	if err != nil {
		log.Printf("error writing ics file, %v", err)
		return err
	}

	a.resources.Links = append(a.resources.Links, internal.Link{
		Name:     c.Title,
		URL:      fmt.Sprintf("/%s/%s", nhlDir, c.Filename),
		Checksum: fmt.Sprintf("%x", sha256.Sum256([]byte(calICS))),
	})

	return nil
}

var (
	funcs = template.FuncMap{
		"splitSeason": func(s *string) string {
			str := *s
			return str[:4] + "-" + str[4:]
		},
		"mapUrl": func(v string) string {
			return "https://www.google.com/maps/search/" + strings.ReplaceAll(v, " ", "%20")
		},
		"formatDate": func(t time.Time) string {
			return t.Format("2006-01-02")
		},
	}
)

const (
	tmpl = `{{ splitSeason .Season }} NHL Season {{ formatDate .GameDate }}

{{ .Teams.Away.Team.Name }} at {{ .Teams.Home.Team.Name }}

Venue:
{{ .Venue.Name }}
{{ mapUrl .Venue.Name}}

Home Team:
{{ .Teams.Home.Team.Name }}
({{.Teams.Home.LeagueRecord.Wins}}-{{.Teams.Home.LeagueRecord.Losses}}-{{.Teams.Home.LeagueRecord.Ot}})

Away Team:
{{ .Teams.Away.Team.Name }}
({{.Teams.Away.LeagueRecord.Wins}}-{{.Teams.Away.LeagueRecord.Losses}}-{{.Teams.Away.LeagueRecord.Ot}})


This Calendar was generated by ope.cool. See more at https://cal.ope.cool 🏒

`
)

func CreateDescription(v *nhl.ScheduleGame) string {
	Create := func(name, t string) *template.Template {
		return template.Must(template.New(name).Funcs(funcs).Parse(t))
	}

	t := Create("description", tmpl)
	var tpl bytes.Buffer
	must(t.Execute(&tpl, v))

	return tpl.String()
}

func ToSchedule(d *nhl.ScheduleGame) internal.Event {
	return internal.Event{
		ID:          fmt.Sprintf("%d-%d-%d", *d.GamePk, *d.Teams.Away.Team.Id, *d.Teams.Home.Team.Id),
		Summary:     fmt.Sprintf("%s @ %s", *d.Teams.Away.Team.Name, *d.Teams.Home.Team.Name),
		Description: CreateDescription(d),
		URL:         *d.Link,
		Location:    *d.Venue.Name,
		OrganizerCN: "https://cal.ope.cool",
		CreatedDate: time.Now(),
		DtStampTime: time.Now(),
		ModifiedAt:  time.Now(),
		StartAt:     *d.GameDate,
		EndAt:       d.GameDate.Add(time.Hour * 3),
	}
}

func paramString(parts []string) *string {
	s := strings.Join(parts, "&")
	return &s
}

func must(err error) {
	if err != nil {
		log.Fatalf("error occured, %v", err)
	}
}

func prettyJson(data interface{}) string {
	bytes, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		panic(err)
	}

	return string(bytes)
}
