package main

import "time"

type Article struct {
	Id          int      `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`

	Content   string     `json:"content"`
    Summary   string     `json:"summary"`
	Tags      string     `json:"tags"`
	Category  string     `json:"cate"`

	CreatedAt   int64    `json:"created_at"`
    UpdatedAt   int64    `json:"updated_at"`
    Url         string   `json:"url"`
}

type Articles []*Article

func (a Articles) Len() int {
    return len(a)
}

func (a Articles) Swap(i, j int) {
    a[i], a[j] = a[j], a[i]
}

func (a Articles) Less(i, j int) bool {
    return a[i].CreatedAt > a[j].CreatedAt
}

type Tag struct {
    Count int      `json:"count"`
    Name  string   `json:"name"`
    Posts []*Article `json:"posts"`
    Url   string   `json:"url"`
}

type Category struct {
    Count   int         `json:"count"`
    Name    string      `json:"name"`
    Posts   []*Article  `json:"posts"`
    Url     string      `json:"url"`
}

type CollatedYear struct {
    Year   string           `json:"year"`
    Months []*CollatedMonth `json:"months"`
    months map[string]*CollatedMonth `json:"-"`
}

type CollatedMonth struct {
    Month  string   `json:"month"`
    Posts  []*Article `json:"posts"`
    month  time.Month `json:"-"`
}

type CollatedYears []*CollatedYear

func (c CollatedYears) Len() int {
    return len(c)
}

func (c CollatedYears) Swap(i, j int) {
    c[i], c[j] = c[j], c[i]
}

func (c CollatedYears) Less(i, j int) bool {
    return c[i].Year > c[j].Year
}

type CollatedMonths []*CollatedMonth

func (c CollatedMonths) Len() int {
    return len(c)
}

func (c CollatedMonths) Swap(i, j int) {
    c[i], c[j] = c[j], c[i]
}

func (c CollatedMonths) Less(i, j int) bool {
    return c[i].month > c[j].month
}
