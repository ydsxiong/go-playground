package poker

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const PlayerPrompt = "Please enter the number of players: "

type CLI struct {
	input  *bufio.Scanner
	output io.Writer
	game   Game
}

func NewPokerCLI(input io.Reader, output io.Writer, game Game) *CLI {
	return &CLI{bufio.NewScanner(input), output, game}
}

func (pc *CLI) readline() string {
	pc.input.Scan()
	return pc.input.Text()
}

func (pc *CLI) PlayPoker() {
	fmt.Fprint(pc.output, PlayerPrompt)
	numberofplayers, _ := strconv.Atoi(pc.readline())
	pc.game.Start(numberofplayers, pc.output)
	pc.game.Finish(extractWinner(pc.readline()))
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins", "", 1)
}
