package internal

import (
	"time"

	ics "github.com/arran4/golang-ical"
)

type ResourceData struct {
	Links       []Link    `json:"URLs"`
	LastUpdated time.Time `json:"LastUpdated"`
}

type Link struct {
	Name     string `json:"Name"`
	URL      string `json:"URL"`
	Checksum string `json:"Checksum"`
}

type Calendar struct {
	Title    string
	Filename string
	Events   []Event
}

func (c *Calendar) ToICS() string {
	cal := ics.NewCalendar()

	cal.SetMethod(ics.MethodRequest)

	for _, v := range c.Events {
		cal.AddVEvent(v.ToISCEvent())
	}

	return cal.Serialize()
}

func (c *Calendar) GetICSEvents() []ics.VEvent {
	e := make([]ics.VEvent, len(c.Events))
	for i, v := range c.Events {
		e[i] = *v.ToISCEvent()
	}
	return e
}

type Event struct {
	ID             string
	Summary        string
	Description    string
	URL            string
	Location       string
	OrganizerEmail string
	OrganizerCN    string
	CreatedDate    time.Time
	DtStampTime    time.Time
	ModifiedAt     time.Time
	StartAt        time.Time
	EndAt          time.Time
}

func (e *Event) ToISCEvent() *ics.VEvent {
	v := ics.NewEvent(e.ID)
	v.SetSummary(e.Summary)
	v.SetDescription(e.Description)
	v.SetURL(e.URL)
	v.SetLocation(e.Location)
	v.SetOrganizer(e.OrganizerCN)
	v.SetCreatedTime(e.CreatedDate)
	v.SetDtStampTime(e.DtStampTime)
	v.SetModifiedAt(e.ModifiedAt)
	v.SetStartAt(e.StartAt)
	v.SetEndAt(e.EndAt)
	return v
}
