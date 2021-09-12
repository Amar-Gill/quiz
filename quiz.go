package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	var quizFileName string
	var timeLimit int

	flag.StringVar(&quizFileName, "csv", "problems.csv", "a csv file in the format of 'question,answer'")
	flag.IntVar(&timeLimit, "limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()

	fmt.Println("File: " + quizFileName)

	f, err := os.Open(quizFileName)

	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(f)

	lines, err := r.ReadAll()

	if err != nil {
		log.Fatal(err)
	}

	problems := parseLines(lines)

	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)

	numCorrect := 0
	for i, problem := range problems {
		fmt.Printf("Problem #%d: %s\n", i+1, problem.q)
		answerCh := make(chan string)
		go func() {
			var userInput string

			fmt.Print("Answer? ")
			fmt.Scanln(&userInput)

			answerCh <- userInput
		}()

		select {
		case <-timer.C:
			fmt.Println("You're too late! Mwahahahahahhaha.")
			fmt.Printf("Correct answers: %d / %d\n", numCorrect, len(lines))
			return
		case answer := <-answerCh:
			if answer == problem.a {
				numCorrect++
			}
		}
	}

	fmt.Printf("Correct answers: %d / %d\n", numCorrect, len(lines))

}

type Problem struct {
	q string
	a string
}

func parseLines(lines [][]string) []Problem {
	ret := make([]Problem, len(lines))

	for i, line := range lines {
		ret[i].q = line[0]
		ret[i].a = line[1]
	}

	return ret
}
