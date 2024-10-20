package main

import (
	"testing"
)

// TestNewTask tests the NewTask function

func TestNewTask(t *testing.T) {
	// Create a new task
	task, err := NewTask("Test Task")
	if err != nil {
		t.Fatalf("Error creating new task: %v", err)
	}
	// Check that the task description is correct
	if task.TaskDescription != "Test Task" {
		t.Fatalf("Task description is incorrect: %v", task.TaskDescription)
	}
}
