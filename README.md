# Snaketrap

NOTE: This is a WIP. None of the interfaces are final yet.

A HipChat bot wrangler that works off HipChat's "Integration" system, meaning you need
to set up a name and command, the second parameter of your message will be used
to pick a specific bot to handle your request.

`/bot <botname> <command> <args>...`

For example; Switching to the next sheriff/engineer on duty.

`/bot sheriff next`

## Install

- go get github.com/gerbenjacobs/snaketrap
- go run main.go

## Usage 

- copy `config.json.example` to `config.json`
- Go to hipchat.com web interface, find your room and create an Integration, 
put the generated auth key in `config.json` under `"hipchat": { "bot_auth": "key-here" }`
- (Optional for Sheriff) Go to hipchat.com web interface, go to your settings, 
under "API Access" and create a token for `Administer room` and `Send notification` scopes,
put the generated scope key in `config.json` under `"hipchat": { "scope_auth": "key-here" }`
## Current bots

- Sheriff - Keeps a list of engineers on duty and rotates daily