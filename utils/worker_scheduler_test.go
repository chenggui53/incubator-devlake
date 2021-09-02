package utils

import (
	"errors"
	"testing"
	"time"
)

func TestNewWorkerScheduler(t *testing.T) {
	testChannel := make(chan int, 100)
	s, _ := NewWorkerScheduler(5, 2)
	defer s.Release()
	for i := 1; i <= 5; i++ {
		t := i
		_ = s.Submit(func() error {
			testChannel <- t
			return nil
		})
	}
	time.Sleep(1200 * time.Millisecond)
	if len(testChannel) < 2 {
		t.Fatal(`worker not start`)
	}
	if len(testChannel) > 2 {
		t.Fatal(`worker run too fast`)
	}
	time.Sleep(time.Second)
	if len(testChannel) < 4 {
		t.Fatal(`worker not run after a second`)
	}
	if len(testChannel) > 4 {
		t.Fatal(`worker run too fast after a second`)
	}
	s.WaitUntilFinish()
	if len(*s.workerErrors) != 0 {
		t.Fatal(`worker got panic`)
	}
	if len(testChannel) != 5 {
		t.Fatal(`worker not wait until finish`)
	}
}

func TestNewWorkerSchedulerWithoutSecond(t *testing.T) {
	testChannel := make(chan int, 100)
	s, _ := NewWorkerScheduler(5, 0)
	defer s.Release()
	for i := 1; i <= 5; i++ {
		t := i
		_ = s.Submit(func() error {
			testChannel <- t
			return nil
		})
	}
	time.Sleep(5 * time.Millisecond)
	if len(testChannel) != 5 {
		t.Fatal(`worker not finish`)
	}
	s.WaitUntilFinish()
	if len(testChannel) != 5 {
		t.Fatal(`worker not finish`)
	}
}

func TestNewWorkerSchedulerWithPanic(t *testing.T) {
	testChannel := make(chan int, 100)
	s, _ := NewWorkerScheduler(1, 1)
	defer s.Release()
	_ = s.Submit(func() error {
		testChannel <- 1
		return errors.New(`error message`)
	})
	s.WaitUntilFinish()
	if len(*s.workerErrors) != 1 {
		t.Fatal(`worker not got panic`)
	}
}
