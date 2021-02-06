package main

import (
	"errors"
	"fmt"
	"github.com/jedib0t/go-pretty/table"
	"os"
	"reflect"
)

type Player struct {
	name      string
	represent string
}
type Square struct {
	x, y   int
	player *Player
}

type Board struct {
	size          int
	squares       [][]*Square
	players       [2]*Player
	currentPlayer *Player
	nextPlayer    *Player
	winner        *Player
}

type Game interface {
	reset() error
	playRound(player Player, x, y int)
}

func (board *Board) reset() error {
	// check the size of the board
	if board.size == 0 {
		return errors.New("board doesn't have size")
	}
	if board.size%2 == 0 {
		return errors.New("board size must be odd")
	}

	// now prepare an empty row of squares
	var currentRow []*Square

	// magical loop to fill the board with squares
	// will continue until the size equals to what has been desired
	for y := 0; len(board.squares) < board.size*board.size; y++ {
		currentRow = append(currentRow, &Square{
			x: len(board.squares),
			y: y,
		})

		if y == board.size-1 {
			board.squares = append(board.squares, currentRow)
			y = -1
			currentRow = []*Square{}
		}
	}

	// now ask for players
	var players [2]*Player
	for playerNumber := 0; playerNumber < 2; playerNumber++ {
		// ask for name
		var name string
		fmt.Printf("Player %d name: ", playerNumber+1)
		fmt.Scanln(&name)

		// and for representation
		var represent string
		fmt.Printf("%s represents: ", name)
		fmt.Scanln(&represent)

		// create player according to the information
		player := &Player{
			name:      name,
			represent: represent,
		}

		// push it for further use
		players[playerNumber] = player
	}

	// set the players now
	board.players = players

	// set the first player
	board.currentPlayer = board.players[0]

	// and also next player
	board.nextPlayer = board.players[1]

	return nil
}

func (board *Board) playRound(x, y int) error {
	// first need to check the current represented value
	if board.squares[x][y].player != nil {
		return errors.New("square already taken")
	}

	// and then represent the square
	board.squares[x][y].player = board.currentPlayer

	return nil
}

func (board *Board) handleNextRound() error {
	// get the location of next square
	var x, y int
	fmt.Printf("%s enter position in format: x y: ", board.currentPlayer.name)
	inputsCount, _ := fmt.Scanf("%d %d", &x, &y)

	// if inputs count doesn't match the expectation throw error
	if inputsCount != 2 {
		return errors.New("invalid number of inputs given")
	}

	// otherwise check for validity of inputs
	if x-1 < 0 || x-1 > board.size || y-1 < 0 || y-1 > board.size {
		return errors.New("invalid square selected")
	}

	err := board.playRound(x-1, y-1)

	if err != nil {
		return err
	}

	// swap players as well
	currentPlayer := board.currentPlayer
	board.currentPlayer = board.nextPlayer
	board.nextPlayer = currentPlayer

	return nil
}

func (board Board) checkResult() (finished bool, winner *Player) {
	var currentRow []*Player
	var currentColumn []*Player
	var currentDiameter []*Player
	var winnerPatterns [2][]*Player

	for i := 0; i < board.size; i++ {
		winnerPatterns[0] = append(winnerPatterns[0], board.players[0])
		winnerPatterns[1] = append(winnerPatterns[1], board.players[1])
	}

	for index1 := 0; index1 < board.size; index1++ {
		for index2 := 0; index2 < board.size; index2++ {
			currentRow = append(currentRow, board.squares[index1][index2].player)
			currentColumn = append(currentColumn, board.squares[index2][index1].player)
			currentDiameter = append(currentDiameter, board.squares[index2][index2].player)
		}
		for _, winnerPattern := range winnerPatterns {
			if reflect.DeepEqual(currentRow, winnerPattern) {
				return true, currentRow[0]
			}
			if reflect.DeepEqual(currentColumn, winnerPattern) {
				return true, currentColumn[0]
			}
			if reflect.DeepEqual(currentDiameter, winnerPattern) {
				return true, currentDiameter[0]
			}
		}
	}

	return false, nil
}

func (board *Board) render() (error, bool) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	for x := 0; x < board.size; x++ {
		var row []string
		for y := 0; y < board.size; y++ {
			player := board.squares[x][y].player
			if player != nil {
				row = append(row, player.represent)
			} else {
				row = append(row, " ")
			}
		}
		// append rows
		t.AppendRow(table.Row{row})
	}

	// now render the table
	t.Render()

	// check for result first
	finished, winner := board.checkResult()

	// if game has been finished then print the result
	if finished {
		board.winner = winner
		fmt.Println("Game finished")
		fmt.Printf("Winner %s who represent %s\n", winner.name, winner.represent)
		return nil, finished
	}

	// otherwise play round
	err := board.handleNextRound()

	return err, false
}
