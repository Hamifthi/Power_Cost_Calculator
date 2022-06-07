# NewMotion programming assignment v2.0.3

## First things first

Please solve this assignment in Go and give us instructions so that we can run your solution running in a shell in a matter of minutes.
If you need any clarifications about the assignment or you'd like to report some inconsistencies, please open an issue on GitHub, tag one of us, or reach out directly (ask for a contact in the team).

## The assignment

You are given two files in CSV format:

- `tariffs.csv` contains a list of tariffs for an arbitrary period of time
- `sessions.csv` contains a list of sessions for the same period of time

Please write a program that will read both input CSV files and produce a `costs.csv` file whose rows have the following fields:

- the session identifier
- its total cost

Please note that your solution should be able to ingest and process large `sessions.csv` files (= hundreds of MBs).

### Tariffs

A tariff is defined by the following fields:

- a start instant
- an end instant
- an energy fee (in €/kWh)
- a parking fee (in €/hour)

The input `tariffs.csv` dataset is guaranteed to have the following properties:

- no two tariffs will overlap with each other (= there is only one _active_ tariff at any point in time)
- tariffs will be sorted in ascending order by start instant
- there are no time gaps between any two adjacent tariffs (= the end instant of a tariff is equal to the start instant of the next one)

### Sessions

A session is defined by the following fields:

- a unique identifier
- a start instant
- an end instant
- the energy consumed (in kWh)

The input `sessions.csv` dataset is guaranteed to have the following properties:

- each session happened in a time interval for which _at least_ one active tariff exists

### Cost calculation

Given:

```
tariff_cost = energy_consumed_in_tariff_interval * energy_fee + duration_in_tariff_interval * parking_fee
```

then the session's `total_cost` is calculated by:

- taking the sum of all _applicable_ tariffs' `tariff_cost`s
- multiplying it by `1.15` (= 15%, our service fee), and
- truncating the result to 3 decimal digits

**Notes:**

- duration intervals are intended as floating point numbers of hours;
- for the sake of simplicity, please assume that energy consumption is constant throughout the duration of the session;

A tariff is applicable to a session if:

- it was _active_ when the session started, or
- it became the _active_ one while the session was in progress

### A note about datetime strings

Datetime strings in CSV files are formatted according to RFC 3339, so the following examples are all valid:

- 2020-05-05T08:00:00+01:00
- 2020-05-05T08:00:00.0Z
- 2020-05-05T08:00:00.123456-01:00

### Examples of data

#### Input

`tariffs.csv`:
```
dt_start,dt_end,energy_fee,parking_fee
2020-06-03T00:00:00+02:00,2020-06-03T10:10:00+02:00,0.5,0.5
2020-06-03T10:10:00+02:00,2020-06-03T14:10:00+02:00,0.3,0.25
2020-06-03T14:10:00+02:00,2020-06-04T22:00:00+02:00,0.7,0.1
```

`sessions.csv`:
```
id,dt_start,dt_end,energy
a949d681-e12b-4d93-a3e5-e2e777e68f12,2020-06-03T10:00:00+02:00,2020-06-03T10:26:00+02:00,3.69
d467391f-0213-44ef-8d68-eea9144e9aa3,2020-06-03T14:49:00+02:00,2020-06-03T14:55:00+02:00,3
99e8fa40-e5e9-4957-9cf3-02e5aa78cf67,2020-06-04T02:52:00+02:00,2020-06-04T02:53:00+02:00,0
```

#### Output

`costs.csv`:
```
session_id,total_cost
a949d681-e12b-4d93-a3e5-e2e777e68f12,1.771
d467391f-0213-44ef-8d68-eea9144e9aa3,2.426
99e8fa40-e5e9-4957-9cf3-02e5aa78cf67,0.001
```
