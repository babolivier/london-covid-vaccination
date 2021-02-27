# London COVID-19 vaccination stats

This project retrieves COVID-19 vaccination data in London and exposes them
through an HTTP API and a simple visualisation graph.

I might write a blog post about the motivation behind this project soon-ish,
so watch this space.

## Miner

The miner is the component that retrieves the data from the NHS website and
stores it in the database. It runs in the background every hour, with a pretty
simple sequence:

1. Crawl the vaccination stats [page](https://www.england.nhs.uk/statistics/statistical-work-areas/covid-19-vaccinations/)
   for England on the NHS website and identifies XLSX files containing daily
   cumulative values for the number of administered doses per region.
2. Download each file that it hasn't already processed in a previous iteration
   and parse it.
3. Locate the data for London and store it.

As a legal note, this data is provided by the NHS under the
[Open Government Licence v3.0](https://www.nationalarchives.gov.uk/doc/open-government-licence/version/3/).

## API

An HTTP+JSON API is hosted to allow for the automated retrieval of this data.
A very simple visualisation graph (that uses this API) is also hosted at the
root of the HTTP server.

A live version of this API (and this visualisation graph) are available at
https://covid-vax-lon.abolivier.bzh/

### `GET /stats`

Response body:

```json
[
    {
        "date": "2021-02-26",
        "first_dose": 1927205,
        "second_dose": 74210
    },
    {
        "date": "2021-02-27",
        "first_dose": 1984318,
        "second_dose": 76897
    },
    {
        "date": "2021-02-25",
        "first_dose": 1870423,
        "second_dose": 71197
    },
    {
        "date": "2021-02-24",
        "first_dose": 1816819,
        "second_dose": 68836
    },
    {
        "date": "2021-02-23",
        "first_dose": 1769695,
        "second_dose": 67585
    },
    [...]
]
```

* `date`: The date at which these number were published, in the format `YYYY-MM-DD`.
* `first_dose`: The total number of first doses of the vaccine administered _up to_ this date.
* `second_dose`: The total number of second doses of the vaccine administered _up to_ this date.