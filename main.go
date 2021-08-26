package main

import (
	"fmt"
	"github.com/99-66/go-worker-pool-example/models"
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

func main() {
	var (
		wg          sync.WaitGroup
		jobsChannel = make(chan models.Job)
	)

	// main 함수가 끝나고 소요시간을 출력한다
	defer elapsed()()

	// Worker수를 지정한다
	numberOfWorkers := runtime.NumCPU()
	if workerCount := os.Getenv("WORKER_COUNT"); workerCount != "" {
		numberOfWorkers, _ = strconv.Atoi(workerCount)
	}

	// Worker수만큼 waitGroup을 설정한다
	wg.Add(numberOfWorkers)

	var jobs []models.Job
	// 실행할 job을 생성한다
	numberOfJobs := 10
	for i := 0; i < numberOfJobs; i++ {
		jobs = append(jobs, models.Job{
			Id: fmt.Sprintf("%d", i),
		})
	}

	// Worker수 만큼 worker를 실행한다
	for i := 0; i < numberOfWorkers; i++ {
		go worker(&wg, jobsChannel)
	}

	// job을 channel을 통해 worker로 전달한다
	for _, job := range jobs {
		jobsChannel <- job
	}

	close(jobsChannel)
	wg.Wait()
}

func testFunc(job *models.Job) {
	log.Printf("running jobs! %s\n", job.Id)
}

// worker 채널로 전달된 Job을 수행한다
// 함수가 종료될 때마다 waitGroup을 1씩 감소시킨다
func worker(wg *sync.WaitGroup, jobsChannel <-chan models.Job) {
	defer wg.Done()

	// jobsChannel로 전달된 job을 실행한다
	for job := range jobsChannel {
		testFunc(&job)
	}
}

// elapsed Main 함수의 소요시간을 확인하기 위한 함수이다
func elapsed() func() {
	start := time.Now()

	return func() {
		log.Printf("elapsed time: %s\n", time.Since(start))
	}
}
