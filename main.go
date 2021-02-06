package main

import "fmt"

func main() {
	board := Board{
		size: 3,
	}

	if err := board.reset(); err != nil {
		fmt.Println(err)
		return
	}

	err, finished := board.render()

	for !finished {
		if err != nil {
			fmt.Print(err)
		}

		err, finished = board.render()
	}
}
