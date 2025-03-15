# LLM playing a snake game via REST API

## What?! Why!? 😂 🐍 🍎 🤖 🎮

I thought it would be awesome to have an LLM drive a snake game in a container and to watch the game via VNC.
To start, I used ChatGPT to iterate on a simple golang snake game that a human could play.
Then I put the snake game in a container and played it over VNC. Noice!

But ultimately, I wanted THE LLM TO PLAY, autonomously with this EXTERNAL LIVE system over an API.

So I dropped the container and the VNC for a second and instead reworked (with the help of ChatGPT again) the snake game to run an http server that exposes a simple API of /start, /status, /move.

Then I could run the snake game via go on my machine and control it via the API which I did manually at first via `curl`.

I created functions in a dagger module called `GoSnake` to wrap these calls.

Since Dagger's LLM features allow functions to be exposed as tools for the LLM to call, the theory was that the LLM could play the LIVE running snake game on my machine via these tool calls that would hit the API.

In practice, getting the LLM to do it has been tough. I tried a bunch of local models with Ollama, and some got close, but it's tricky to even get the LLM to understand the situation and to get it to "looK" with getGameState() often enough to not just run into a wall or hallucinate a game state or try to devise a strategy or code to "solve" the problem of playing the game.

## The Challenge

SOOOOOO, I thought it would be awesome if folks tried to alter the prompt or GoSnake functions/tools or whatever to make this actually work. It's probably about finding the right combo of model, prompt, tools UX (e.g. one `makeMove($direction)` function or separate `up()`, `left()`, etc functions), as well as making the game not too slow or too fast, the game grid not too big or too small. Or maybe the representation of the game grid (currently in JSON needs to be improved??).

**A BIG solution space to explore! Will you make it work!? Will you get the high score?!**


