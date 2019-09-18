package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ydsxiong/go-playground/poker-app/poker"
)

const dbFileName = "../../game.db.json"

func main() {

	store, closeStore, err := poker.LoadUpFileStore(dbFileName)
	if err != nil {
		log.Fatalf("Problem with loading in file store, %v", err)
	}
	defer closeStore()

	fmt.Println("Let's play some poker...")
	fmt.Println("Type {Name} wins to record a win")
	poker.NewPokerCLI(os.Stdin, os.Stdout, poker.NewGame(store, poker.BlindAlerterFunc(poker.StdOutAlerter))).PlayPoker()
}
