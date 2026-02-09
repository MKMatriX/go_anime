package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// AnimeInfo — структура для основных данных аниме
type AnimeInfo struct {
	ID           int      `json:"id"`
	Title        Title    `json:"title"`
	Description  string   `json:"description"`
	Format       string   `json:"format"`
	Status       string   `json:"status"`
	Episodes     *int     `json:"episodes"`
	Duration     *int     `json:"duration"`
	AverageScore *int     `json:"averageScore"`
	Genres       []string `json:"genres"`
	Season       string   `json:"season"`
	SeasonYear   *int     `json:"seasonYear"`
	StartDate    Date     `json:"startDate"`
	CoverImage   Image    `json:"coverImage"`
}

type Title struct {
	Romaji  string `json:"romaji"`
	English string `json:"english"`
	Native  string `json:"native"`
}

type Date struct {
	Year  *int `json:"year"`
	Month *int `json:"month"`
	Day   *int `json:"day"`
}

type Image struct {
	Large string `json:"large"`
}

// AniListResponse — структура ответа от AniList
type AniListResponse struct {
	Data struct {
		Media AnimeInfo `json:"Media"`
	} `json:"data"`
}

func GetAnimeInfo(title string) (*AnimeInfo, error) {
	query := `
	query ($search: String) {
	  Media(search: $search, type: ANIME, sort: SEARCH_MATCH) {
	    id
	    title {
	      romaji
	      english
	      native
	    }
	    description(asHtml: false)
	    format
	    status
	    episodes
	    duration
	    averageScore
	    genres
	    season
	    seasonYear
	    startDate {
	      year
	      month
	      day
	    }
	    coverImage {
	      large
	    }
	  }
	}`

	variables := map[string]string{
		"search": title,
	}

	requestBody := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://graphql.anilist.co", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "go_anime/0.1 (matrix-elf@yandex.ru)")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("AniList API error: %d - %s", resp.StatusCode, string(body))
	}

	var result AniListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Data.Media.ID == 0 {
		return nil, fmt.Errorf("anime not found for title: %s", title)
	}

	return &result.Data.Media, nil
}
