# cal.ope.cool

> A simple project for generating iCal Calendars for Sporting Events

[https://cal.ope.cool](https://cal.ope.cool)

This project is a proof of concept for generating an iCal for a list of sporting events.

Uses:
- Go
- [nhlapi](https://github.com/dword4/nhlapi)
- [oapi-codegen](https://github.com/deepmap/oapi-codegen)
- [golang-ical](https://github.com/arran4/golang-ical)

## TODO

| Goal | Progress |
| ---- | -------- |
| Basic Website to Display | Complete ✅ |
| Individual NHL Team Calendars | Complete ✅ |
| Bulk NHL Schedule for All Games | Pending |
| Define other Sports & API availability | TBD |
| Persist Cache in DB (either KV/bolt or Sqlite) | TBD |
| Setup Github Actions for daily cron | TBD |
| With persistence, setup checksum for diff checking ical | TBD |
