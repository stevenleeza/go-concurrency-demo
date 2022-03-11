package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/briandowns/spinner"
	"github.com/manifoldco/promptui"
)

const (
	STCT = "Single-threaded constant time"
	MTCT = "Multithreaded constant time"
	STRT = "Single-threaded random time"
	MTRT = "Multithreaded random time"
	EXIT = "-Exit Program-"
)

type Job struct {
	id       int
	duration time.Duration
}

func main() {
	for {
		canContinue := demo()

		if !canContinue {
			break
		}
	}
}

func demo() bool {
	userChoice, err := userPrompt()

	if err != nil || EXIT == userChoice {
		return false
	}

	start := time.Now()

	spinner := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	spinner.Color("red")
	spinner.Start()

	const numJobs = 10

	var channel = make(chan Job, numJobs)

	rand.Seed(time.Now().UnixNano())

	var isRandom bool
	var isMultithreaded bool
	var description string

	switch userChoice {
	case STCT:
		isRandom = false
		isMultithreaded = false
		description = "Jobs are executed sequentially in constant time.\nJobs are guaranteed to complete in order, but at the expense of execution time.\nNote how nothing is printed out until all the jobs have completed due to the job function calls being blocking."
	case MTCT:
		isRandom = false
		isMultithreaded = true
		description = "Jobs are executed concurrently in constant time.\nJobs are not guaranteed to complete in order as they could be distributed across multiple threads.\nThe execution time is dependent on the number of CPU cores available on the host system."
	case STRT:
		isRandom = true
		isMultithreaded = false
		description = "Jobs are executed sequentially in random time to emulate the unpredictable latencies of real-world execution.\nJobs are guaranteed to complete in order, but at the expense of execution time.\nNote how nothing is printed out until all the jobs have completed due to the job function calls being blocking."
	case MTRT:
		isRandom = true
		isMultithreaded = true
		description = "Jobs are executed concurrently in random time to emulate the unpredictable latencies of real-world execution.\nJobs are not guaranteed to complete in order as they could be distributed across multiple threads.\nThe execution time is dependent on the number of CPU cores available on the host system.\nNote how the job results are printed out immediately due to the job function calls being non-blocking.\nIf you've made it this far, you've hopefully seen the impressive gains from utilising multithreading!"
	}

	fmt.Println(description)

	for i := 0; i < numJobs; i++ {
		var variance float64

		if isRandom {
			variance = rand.Float64() * 3
		} else {
			variance = 1
		}

		duration := math.Round(float64(time.Second) * variance)

		if isMultithreaded {
			go job(i+1, time.Duration(duration), channel)
		} else {
			job(i+1, time.Duration(duration), channel)
		}
	}

	jobCompleteCounter := 0

	for job := range channel {
		jobCompleteCounter++
		spinner.Stop()
		fmt.Printf("\nJob %d:\t%dms", job.id, job.duration.Milliseconds())

		if jobCompleteCounter == numJobs {
			fmt.Print("\n\n")
			close(channel)
		}
	}

	elapsed := time.Since(start)
	fmt.Printf("Execution time: %s\n\n", elapsed)

	return true
}

func userPrompt() (string, error) {
	prompt := promptui.Select{
		Label: "Select Demo",
		Items: []string{STCT, MTCT, STRT, MTRT, EXIT},
	}

	_, userChoice, err := prompt.Run()

	if err != nil {
		return "", err
	}

	return userChoice, nil
}

func job(id int, duration time.Duration, channel chan<- Job) {
	time.Sleep(duration)

	channel <- Job{
		id:       id,
		duration: duration,
	}
}
