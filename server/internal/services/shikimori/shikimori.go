package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ShikiAnimeRu — структура с нужными русскими полями
type ShikiAnimeRu struct {
	ID                string               `json:"id"`
	Name              string               `json:"name"`    // ромадзи / основное
	Russian           string               `json:"russian"` // русское название
	English           string               `json:"english"` // английское
	Description       string               `json:"description"`
	DescriptionHtml   string               `json:"descriptionHtml"`
	DescriptionSource string               `json:"descriptionSource"`
	Score             float64              `json:"score"`
	Kind              string               `json:"kind"`
	Status            string               `json:"status"`
	Episodes          int                  `json:"episodes"`
	Poster            ShikiPoster          `json:"poster"`
	PersonRoles       []ShikiPersonRoles   `json:"personRoles"`
	ExternalLinks     []ShikiExternalLinks `json:"externalLinks"`
}

type ShikiPoster struct {
	ID          string `json:"id"`
	OriginalUrl string `json:"originalUrl"`
	MainUrl     string `json:"mainUrl"`
}

type ShikiPersonRoles struct {
	ID      string      `json:"id"`
	RolesRu []string    `json:"rolesRu"`
	RolesEn []string    `json:"rolesEu"`
	Person  ShikiPerson `json:"person"`
}

type ShikiPerson struct {
	ID     string       `json:"id"`
	Name   string       `json:"name"`
	Poster *ShikiPoster `json:"poster"`
}

type ShikiExternalLinks struct {
	ID   *string `json:"id"`
	Kind string  `json:"kind"`
	Url  string  `json:"url"`
}

// GraphQLResponse — обёртка ответа
type ShikiGraphQLResponse struct {
	Data struct {
		Animes []ShikiAnimeRu `json:"animes"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func GetShikiAnimeInfo(search string) (*ShikiAnimeRu, error) {
	const endpoint = "https://shikimori.one/api/graphql"

	query := `
	query SearchAnime($search: String!) {
		animes(search: $search, limit: 1) {
			id
			name
			russian
			english
			description
			descriptionHtml
			descriptionSource
			score
			kind
			status
			episodes
			poster {
				id
				originalUrl
				mainUrl
			}
			personRoles {
				id
				rolesRu
				rolesEn
				person { id name poster { id } }
			}
			externalLinks {
				id
				kind
				url
			}
		}
	}`

	variables := map[string]string{
		"search": search,
	}

	requestBody := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "go_anime/0.1 (matrix-elf@yandex.ru)")

	client := &http.Client{Timeout: 12 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("shikimori graphql http %d: %s", resp.StatusCode, string(body))
	}

	var result ShikiGraphQLResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("graphql error: %s", result.Errors[0].Message)
	}

	if len(result.Data.Animes) == 0 {
		return nil, fmt.Errorf("аниме не найдено по запросу: %q", search)
	}

	anime := &result.Data.Animes[0]

	return anime, nil
}
