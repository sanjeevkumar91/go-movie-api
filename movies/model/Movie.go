package model

type Movie struct {
	Title  string
	Year   string
	ImdbID string
	Type   string
	Poster string
}

type SearchMovieRequest struct {
	Title       string `json:"title,omitempty"`
	Type        string `json:"type,omitempty"`
	Year        string `json:"year,omitempty"`
	SearchQuery string `json:"searchText" binding:"required"`
	Page        string `json:"page,omitempty"`
}

type SearchMovieResponse struct {
	Movies   []Movie `json:"search"`
	Response string  `json:"Response"`
	Error    string  `json:"Error"`
}

type Rating struct {
	Source string
	Value  string
}

type GetMovieDetailsResponse struct {
	Title      string   `json:"Title"`
	Year       string   `json:"Year"`
	Rated      string   `json:"Rated"`
	Released   string   `json:"Released"`
	Genre      string   `json:"Genre"`
	Runtime    string   `json:"Runtime"`
	Director   string   `json:"Director"`
	Actors     string   `json:"Actors"`
	Plot       string   `json:"Plot"`
	Language   string   `json:"Language"`
	Country    string   `json:"Country"`
	Awards     string   `json:"Awards"`
	Poster     string   `json:"Poster"`
	Ratings    []Rating `json:"Ratings"`
	Metascore  string   `json:"Metascore"`
	ImdbRating string   `json:"imdbRating"`
	ImdbID     string   `json:"ImdbID"`
	Type       string   `json:"Type"`
	DVD        string   `json:"N/A"`
	BoxOffice  string   `json:"BoxOffice"`
	Production string   `json:"Production"`
	Website    string   `json:"Website"`
	Response   string   `json:"Response"`
	Error      string   `json:"Error"`
}

type GetMovieDetailsRequest struct {
	Title   string `json:"title,omitempty"`
	MovieID string `json:"movieId,omitempty"`
	Type    string `json:"type,omitempty"`
	Year    string `json:"year,omitempty"`
}

type AddMovieToCartRequest struct {
	MovieID string `json:"movieId" binding:"required"`
}

type AddMovieToCartResponse struct {
	Status string `json:"status"`
}
