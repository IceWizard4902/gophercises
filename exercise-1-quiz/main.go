package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// This to "generalize" the code and not hardcode it to work only with 2D arrays
type problem struct {
	question string
	answer   string
}

func parseLines(lines [][]string) []problem {
	// If we already know the length, just specify it instead of doing append
	ret := make([]problem, len(lines))

	for i, line := range lines {
		ret[i] = problem{
			question: line[0],
			// Cover the case of CSV file has some spaces, making it "impossible" to solve
			// due to the space trimming behavior of Scanf
			answer: strings.TrimSpace(line[1]),
		}
	}

	return ret
}

func main() {
	csvFileName := flag.String("csv", "problems.csv", "a csv file in the format of 'question, answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()

	file, err := os.Open(*csvFileName)
	if err != nil {
		log.Fatalf("Unable to open the CSV file: %s\n", *csvFileName)
	}
	r := csv.NewReader(file)

	// Not going to have memory problems as the size of problems.csv is not big
	lines, err := r.ReadAll()
	if err != nil {
		log.Fatalf("Failed to parse CSV file with error: %v", err)
	}

	problems := parseLines(lines)

	// Timer runs only once, Ticker runs multiple times
	// Timer here so that the setup time of the code does not penalize the player
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	// Blocks and stop the program a message is sent to the channel
	// <-timer.C
	scores := 0

problemloop:
	for i, problem := range problems {
		fmt.Printf("Problem: %d: %s =", i+1, problem.question)
		answerCh := make(chan string)
		go func() {
			var answer string
			// Scrapes spaces so it is appropriate for this task
			fmt.Scanf("%s\n", &answer)
			// Send answer to channel
			answerCh <- answer
		}()

		// See https://gobyexample.com/select
		select {
		case <-timer.C:
			fmt.Println()
			break problemloop
		case answer := <-answerCh:
			if answer == problem.answer {
				scores++
			}
		}
	}
	fmt.Printf("You scored %d out of %d\n", scores, len(problems))
}
