package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

type problem struct {
	question string
	answer   string
}

// problemPuller Function
func problemPuller(filename string) ([]problem, error) {
	// Read all the problems from the quiz.csv fil by opening the fle

	if fObj, err := os.Open(filename); err == nil {
		// Create a new reader for the file
		csvR := csv.NewReader(fObj)

		// it will be used to read the file
		if cLines, err := csvR.ReadAll(); err == nil {
			// Call the parseProblem function
			return parseProblem(cLines), nil
		} else {
			return nil, fmt.Errorf("error in data in csv format from %s file; %s", filename, err.Error())
		}
	} else {
		return nil, fmt.Errorf("error in opening the %s file, %s", filename, err.Error())
	}
}

// parseProblem Function
func parseProblem(lines [][]string) []problem {
	// Go over the lines and parse them with the problem struct
	r := make([]problem, len(lines))

	for i := 0; i < len(lines); i++ {
		r[i] = problem{
			question: lines[i][0],
			answer:   lines[i][1],
		}
	}
	return r
}

// function to exit the program
func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func main() {
	fmt.Println("Golang Quiz")
	// Input the name of the file
	fileName := flag.String("f", "quiz.csv", "path of the csv file")

	// Set Duration of the timer
	timer := flag.Int("t", 30, "timer of the quiz")
	flag.Parse()

	// Pull the problems from the file (call the problem puller function)
	problems, err := problemPuller(*fileName)

	// Handle errors
	if err != nil {
		exit(fmt.Sprintf("Something went wrong:%s", err.Error()))
	}

	// Create a variable to count our correct answers
	correctAns := 0

	// Using the duration of the timer initialize the timer
	tObj := time.NewTimer(time.Duration(*timer) * time.Second)
	ansC := make(chan string)

	// Loop through the problems, print the questions, we'll accept the answers
problemLoop:
	for i, p := range problems {
		var answer string
		fmt.Printf("Problem %d: %s=", i+1, p.question)

		// GoRoutine
		go func() {
			fmt.Scanf("%s", &answer)
			ansC <- answer
		}()

		select {
		case <-tObj.C:
			fmt.Println()
			break problemLoop
		case iAns := <-ansC:
			if iAns == p.answer {
				correctAns++
			}
			if i == len(problems)-1 {
				close(ansC)
			}
		}
	}

	// Calculate and print out the results
	fmt.Printf("Your result is %d out of %d\n", correctAns, len(problems))
	fmt.Println("Press enter to exit")
	<-ansC
}
