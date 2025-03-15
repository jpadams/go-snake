# go-snake connects to the snake game service

call it via code like in `../player`

```
Llm().WithSnakeGame(gameServerService).WithPrompt(...
```

or interact with it more via Dagger Shell!

```
go-snake/ $ dagger

go-snake ⋈ . tcp://0.0.0.0:8080 | start-game
go-snake ⋈ . tcp://0.0.0.0:8080 | make-move up 
```
TODO: currently caching is making only one call of each function possible in a session! 😭


or super fun, crazy!

```
go-snake/ $ dagger llm

gpt-4o ✱ /with go-snake tcp://0.0.0.0:8080
gpt-4o ✱ .echo paste in a prompt!
```
