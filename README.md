# Snake Game server!

## To run use go 1.24.0

- clone this repo and cd into it
- `go run main.go`
- it should open a window and be listening on `:8080`
- next in a new terminal, cd into this repo and `player/` inside it
- check it out!

If you want to play yourself via the API, it's easy using the `player/go-snake/` module's functions or else, just call the api via `curl`:

```
alias start='curl http://localhost:8080/start'

alias up='curl -s -X POST http://localhost:8080/move -d '\''{"move":"up"}'\'' -H "Content-Type: application/json"'

alias down='curl -s -X POST http://localhost:8080/move -d '\''{"move":"down"}'\'' -H "Content-Type: application/json"'

alias left='curl -s -X POST http://localhost:8080/move -d '\''{"move":"left"}'\'' -H "Content-Type: application/json"'

alias right='curl -s -X POST http://localhost:8080/move -d '\''{"move":"right"}'\'' -H "Content-Type: application/json"'
```
