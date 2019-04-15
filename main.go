package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type attribute struct {
	name   string
	factor int
}

type Pokemon struct {
	HP         int
	turn       bool
	Name       string
	Attributes []attribute
}

func getAttributeValue(p *Pokemon, attributeName string) int {
	for _, v := range p.Attributes {
		if v.name == attributeName {
			return v.factor
		}
	}

	return 0
}

func checkAttributes(p *Pokemon, attributeName string) bool {

	for _, v := range p.Attributes {
		if v.name == attributeName {
			return true
		}
	}

	return false
}

func random(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func (p *Pokemon) heal() {
	healPoints := random(10, 20)
	p.HP = p.HP + healPoints
}

func (p *Pokemon) useAttack(Attack int) int {
	attackPoints := 0
	if Attack == 1 {
		fmt.Println(p.Name + " Used Attack 1")
		attackPoints = random(10, 25)
	} else if Attack == 2 {
		fmt.Println(p.Name + " Used Attack 2")
		attackPoints = random(1, 35)
	} else if Attack == 4 {
		fmt.Println(p.Name + " Used Attack 4")
		attackPoints = random(1, 100)
	} else if Attack == 3 {
		fmt.Println(p.Name + " Used Heal")
		p.heal()
	} else {
		println("That is not a selection. You lost your turn!")
	}
	if checkAttributes(p, "Attack") {
		attackPoints += getAttributeValue(p, "Attack")
	}
	if random(0, 10) == 5 {
		fmt.Println("The attack was a critcal hit")
		attackPoints = attackPoints * 2
	}
	fmt.Println("--------------------------------------")

	return attackPoints
}

func (p *Pokemon) reduceHealth(points int) {
	if checkAttributes(p, "Defence") {
		points += getAttributeValue(p, "Defence")
	}

	p.HP = p.HP - points
}

func flipCoin() int {
	return random(1, 3)
}

func (p *Pokemon) getHealth() int {
	health := p.HP

	return health
}

func makeAttributes(counter int) []attribute {

	attributeNames := []string{"Attack", "Defence"}
	attr := []attribute{}

	for index := 0; index < counter; index++ {
		var a attribute
		a.name = attributeNames[random(0, len(attributeNames))]
		a.factor = random(-6, 6)
		attr = append(attr, a)
		time.Sleep(100 * time.Microsecond)
	}

	return attr

}

type Response struct {
	Name    string    `json:"name"`
	Pokemon []Pokemon `json:"pokemon_entries"`
}

func getRandomName() string {
	t := strconv.Itoa(random(0, 150))

	response, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + t)
	if err != nil {
		log.Fatal(err)
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseObject Response
	json.Unmarshal(responseData, &responseObject)

	return responseObject.Name
}

func printPokemon(p Pokemon, printAttributes bool) {

	attackBonus := 0
	defenceBonus := 0

	for _, attr := range p.Attributes {
		if attr.name == "Attack" {
			attackBonus += attr.factor
		} else if attr.name == "Defence" {
			defenceBonus += attr.factor
		}
	}

	fmt.Println("Name: " + p.Name)
	fmt.Println("Current HP: " + strconv.Itoa(p.HP))
	if printAttributes {
		fmt.Println("Attack Bonus: " + strconv.Itoa(attackBonus))
		fmt.Println("Defence Bonus: " + strconv.Itoa(defenceBonus))
	}
	fmt.Println("--------------------------------------")

}

func printStatus(p Pokemon, c Pokemon) {
	fmt.Println("Player:")
	printPokemon(p, false)
	fmt.Println("Computer:")
	printPokemon(c, false)
}

func main() {

	battle := true

	//randomNames := []string{"Pickachu", "Mew", "Charizard", "Snorlax"}

	Player := Pokemon{
		100,
		false,
		getRandomName(),
		makeAttributes(random(0, 6)),
	}

	Computer := Pokemon{
		100,
		false,
		getRandomName(),
		makeAttributes(random(0, 6)),
	}

	fmt.Print("Welcome to pokemon \n")
	fmt.Println("Player Pokemon:")
	printPokemon(Player, true)
	fmt.Println("Computer Pokemon:")
	printPokemon(Computer, true)

	fmt.Println("flip a coin to see who starts 1 for heads 2 for tails")
	var coinflip int
	fmt.Scanln(&coinflip)

	if flipCoin() == coinflip {
		Player.turn = true
		switch coinflip {
		case 1:
			fmt.Println("Heads you won")

		case 2:
			fmt.Println("Tails you won")
		}
	} else {
		Computer.turn = true
		switch coinflip {
		case 1:
			fmt.Println("Heads you lost")

		case 2:
			fmt.Println("Tails you lost")
		}
	}

	i := 1

	for battle {

		if Player.turn == true {
			fmt.Println("Players Turn")
			if i == 1 {
				printPokemon(Player, false)
			}
			fmt.Println("Choose attack 1 = attack 1, 2 = attack 2, 3 = heal")
			var input int
			fmt.Scanln(&input)
			attackPoints := Player.useAttack(input)
			if attackPoints > 0 {

				Computer.reduceHealth(attackPoints)
			}
			if Computer.getHealth() <= 0 {
				fmt.Println("Players Pokemon "+Player.Name, " Wins")
				battle = false
				break
			}

			Player.turn = false
			Computer.turn = true

			printStatus(Player, Computer)
			time.Sleep(1 * time.Second)
		}

		if Computer.turn == true {
			fmt.Println("Computers Turn")
			if Computer.getHealth() > 20 || Player.getHealth() < 20 {
				attackPoints := Computer.useAttack(random(1, 3))
				if attackPoints > 0 {

					Player.reduceHealth(attackPoints)
				}

			} else {
				Player.reduceHealth(Computer.useAttack(3))
			}

			if Player.getHealth() <= 0 {
				fmt.Println("Computers Pokemon "+Computer.Name, " Wins")
				battle = false
				break
			}

			printStatus(Player, Computer)
			time.Sleep(1 * time.Second)

			Player.turn = true
			Computer.turn = false

		}
		i++
	}

}
