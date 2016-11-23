# Snaketrap

NOTE: This is a WIP. None of the interfaces are final yet.

A HipChat bot wrangler that works off HipChat's "Integration" system, meaning you need
to set up name and command, the second parameter of your message will be used
to pick a specific bot to handle your request.

For example; Switching to the next sheriff/engineer on duty.

```
/bot sheriff next
```

## Current bots

- Sheriff - Keeps a list of engineers on duty and rotates daily
- tbd