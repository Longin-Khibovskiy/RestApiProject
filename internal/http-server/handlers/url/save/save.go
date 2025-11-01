package save

import resp "github.com/Longin-Khibovskiy/RestApiProject.git/internal/lib/api/response"

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias"`
}

type URLSaver interface {
	SaveURL(URL, alias string) (int64, error)
}
