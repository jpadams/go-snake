package main

import (
	"context"
	"dagger/player/internal/dagger"
)

func New(
	snakeGameServer *dagger.Service,
) *Player {
	return &Player{
		Server: snakeGameServer,
	}
}

type Player struct {
	Server *dagger.Service
}

// Play snake game
func (m *Player) Play(ctx context.Context) {
	game := dag.GoSnake(m.Server)

	dag.Llm().
		WithGoSnake(game).
		WithPrompt(`You are playing a real-time snake game on a 20x20 board. 
Your snake moves continuously in its last direction—**this is NOT a turn-based game!** If you don't act fast, your snake will crash into a wall and DIE.
You only need to poll for getGameState about once every 1 or 2 seconds.

You will get a JSON blog with a representation of the game board as the 'grid'. It's 20 rows of 20 columns each. The 'direction' in the JSON is the current/last direction the snake is traveling.

Be Patient!! If the grid state seems stagnant, it's just becuase the snake is moving slowly. Use getGameState again to "look" again to see if there is a change.

### **Game Rules:**
- **Start the game:** Call 'startGame' immediately.
- **Check the board state:** Call 'getGameState' frequently to see the current game status.
- **Move the snake:** Use 'makeMove' to change direction and avoid walls while moving toward food.
- **Survival matters:** If you hit a wall or your own tail, you DIE. Stay inside the box!
- **Eat food:** Move the snake’s head to the same square as the food to eat it.

### **Critical Instructions:**
🚨 **ALWAYS use OpenAI-style tool calling!**  
🚨 **DO NOT experiment with other tool-calling formats.**  

### **How to Play:**
- Call 'getGameState' **as frequently as possible** to monitor the game state.
- **Quickly** decide whether to move or call 'getGameState' again.
- **DO NOT overthink or explain your reasoning.** Just take action by calling the correct tools.
- **DO NOT stop and ask for input.** Make your own decisions.
- **DO NOT write code to solve this—just play the game!**
- **Your only goal is to stay alive and eat as much food as possible.**

🔥 **Keep going! Don’t give up! Survive, eat, and win!** 🔥
`).LastReply(ctx)
}
