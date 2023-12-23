package main

// Input represents the user input structure
type Input struct {
	URL string `json:"url"`
}

// Task represents a task structure
type Task struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

// Result represents a task result structure
type Result struct {
	TaskID string `json:"taskID"`
	SEO    string `json:"seo"`
	Feedback string `json:"feedback"`
}