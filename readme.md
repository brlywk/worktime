# Super Simple worktime calculator

## What?

This is really just a very simple command line calculator that let's you either enter a start time,
a pause duration and an end time to get the total amount of hours worked; or to enter a start time
and a pause duration to get the time at which a default number (8) of working hours is fulfilled
(we all know why 😉).

## Why?

Pure convenience. I'm from the European Union so working hours have to be tracked by order of the
law. As I'm moving around in the terminal most of the day anyway, I wrote this quick utility to
calculate the time I have to put into my time sheet for the day.

## How?

**WIP**
Right now the easiest way is to clone the repository, run `go build -o worktime` and than create an
alias in your terminal config, e.g.

```bash
# alias in .zshrc for example
alias worktime='folder-to-your-cloned-repo/worktime'
alias wt=worktime
```

Worktime supports two "modes":

- Running the program without any flags will put it into 'calculate working hours' mode, i.e. the
  total number of hours will be calculated from a start time, pause duration and end time
- Running the program with the `-t` flag will put it in 'calculate end time', i.e. using a default
  number of working hours (default: 8), the start time and pause duration, the time when that amount
  of working hours as been fulfilled will be calculated

### Notes

- Times need to be entered in 24h format: 15:04 instead for 3:04pm
- Times can be entered with or without a colon: 15:04 or 1504
- The result will automatically be copied to the clipboard (if possible)
- Some defaults will be used if parsing of the input time fails: 8:30 for start of the day, 60 min pause duration

## Open Points / (maybe) planned features

- Add ability to enter times as command line arguments
- Some way to provide defaults as settings, either as a config file or as command line variables
