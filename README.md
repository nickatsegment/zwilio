# zwilio

A Twilio webhook server that plays Zork 1 over SMS.

![Screenshot](/screenshot.png?raw=true "Screenshot")

## Usage

```
# download zork
cp path/to/zork1.dat .
go run .
```

Starts a server on 8080. Then:
1. Forward somehow (ngrok?)
2. Point your Twilio webhook at it.
3. SMS it to play. (Any text to start)

## Caveats

- zmachine code adapted from [https://github.com/msinilo/zmachine]. Not sure how much of zmachine it actually supports.
- session handling is pretty naive
- no persistence
- no security
