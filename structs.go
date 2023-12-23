package gslib

type Input struct {
	URL string `json:"url"`
}


type Task struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}


type Result struct {
	TaskID string `json:"taskID"`
	SEO    string `json:"seo"`
	Feedback string `json:"feedback"`
}