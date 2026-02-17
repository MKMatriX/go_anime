package services

import (
	"encoding/xml"
	"fmt"
	"go_anime/internal/shared/models"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// https://feed.animetosho.org/rss2?only_tor=1&q=%5BErai-raws%5D+1080p&aid=19131

// ToshoRSSFeed — основная структура всего фида
type ToshoRSSFeed struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	AtomNS  string   `xml:"xmlns:atom,attr,omitempty"` // http://www.w3.org/2005/Atom
	Channel Channel  `xml:"channel"`
}

// Channel — канал с метаданными и элементами
type Channel struct {
	AtomLink      AtomLink    `xml:"http://www.w3.org/2005/Atom link"` // пространство имён atom
	Title         string      `xml:"title"`
	Link          string      `xml:"link"`
	Description   string      `xml:"description"`
	Language      string      `xml:"language"`
	TTL           int         `xml:"ttl"`
	LastBuildDate string      `xml:"lastBuildDate"` // можно time.Time, но часто строка
	Items         []ToshoItem `xml:"item"`
}

// atomLink — отдельная структура для <atom:link ... />
type AtomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
}

// ToshoItem — каждый <item> в фиде
type ToshoItem struct {
	Title       string    `xml:"title"`
	Description string    `xml:"description"` // CDATA с HTML, например "<strong>Total Size</strong>: 492.7 MB"
	Link        string    `xml:"link"`
	Comments    string    `xml:"comments"`
	Enclosure   Enclosure `xml:"enclosure"`
	Source      Source    `xml:"source"`
	PubDate     string    `xml:"pubDate"` // или time.Time `xml:"pubDate"`
	GUID        GUID      `xml:"guid"`
}

// Enclosure — торрент-файл
type Enclosure struct {
	URL    string `xml:"url,attr"`
	Type   string `xml:"type,attr"`
	Length string `xml:"length,attr"` // часто "0", т.к. размер неизвестен заранее
}

// Source — откуда взята раздача (обычно Nyaa)
type Source struct {
	URL  string `xml:"url,attr"`
	Name string `xml:",chardata"` // "Nyaa", "Anirena" и т.д.
}

// GUID — уникальный идентификатор (часто это ссылка)
type GUID struct {
	IsPermaLink bool   `xml:"isPermaLink,attr"`
	Value       string `xml:",chardata"`
}

func parseAnimeToshoRSS(data []byte) (*ToshoRSSFeed, error) {
	var feed ToshoRSSFeed
	err := xml.Unmarshal(data, &feed)
	if err != nil {
		return nil, err
	}
	return &feed, nil
}

func GetToshoEpisodes(aniDBId uint) ([]ToshoItem, error) {
	base := "https://feed.animetosho.org/rss2"

	params := url.Values{}
	params.Add("only_tor", "1")
	params.Add("q", `[Erai-raws] 1080p`) // автоматически закодирует [ ] и пробел
	params.Add("aid", strconv.FormatUint(uint64(aniDBId), 10))

	fullURL := base + "?" + params.Encode()

	req, err := http.NewRequest("GET", fullURL, http.NoBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/xml")
	req.Header.Set("User-Agent", "go_anime/0.1 (matrix-elf@yandex.ru)")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var body []byte
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AnimeTosho RSS error: %d - %s", resp.StatusCode, string(body))
	}

	parsedData, err := parseAnimeToshoRSS(body)
	if err != nil {
		return nil, err
	}

	return parsedData.Channel.Items, nil
}

// ---- to model

// ParseToshoItemToEpisode пытается преобразовать item из AnimeTosho в AnimeEpisodeModel
// Возвращает модель + bool (ok = удалось ли адекватно распарсить)
func ParseToshoItemToEpisode(item ToshoItem, animeID uint) (models.AnimeEpisodeModel, bool) {
	ep := models.AnimeEpisodeModel{
		AnimeID: animeID,
		// Anime:          заполняется отдельно (по animeID)
		// BaseModel:      заполняется в БД или где-то выше
		TorrentUrl: item.Enclosure.URL,
		// MagnetUrl:   если нужно — можно сгенерировать из хэша (но обычно нет в RSS)
		// LocalUrl:    заполняется позже, после скачивания
	}

	title := strings.TrimSpace(item.Title)

	// 1. Извлекаем номер эпизода
	epNum := extractEpisodeNumber(title)
	if epNum > 0 {
		ep.EpisodeNumber = epNum
	} else {
		return ep, false // без номера серии — считаем невалидным
	}

	// 2. Группа / переводчик (обычно в начале [Erai-raws], [Ohys-Raws] и т.д.)
	ep.Translator = extractGroup(title)

	// 3. Качество / разрешение (720p, 1080p, 4K и т.д.)
	ep.Width = extractResolution(title)

	// 4. Имя — можно оставить почти оригинал, или очистить
	ep.Name = cleanEpisodeName(title)

	// Дополнительно: можно парсить размер из Description, если очень нужно
	// size := extractSizeFromDescription(item.Description)

	return ep, true
}

// ──────────────────────────────────────────────
// Вспомогательные функции-парсеры
// ──────────────────────────────────────────────

var reEpisode = regexp.MustCompile(`(?i)(?:ep\s*|\s|-|第|episode\s*)(\d+)(?:\s*(?:v\d|end|fin|ova|special))?`)

func extractEpisodeNumber(title string) int {
	matches := reEpisode.FindStringSubmatch(title)
	if len(matches) >= 2 {
		num, _ := strconv.Atoi(matches[1])
		return num
	}
	return 0
}

var reGroup = regexp.MustCompile(`^\[([^]]+)\]`)

func extractGroup(title string) string {
	m := reGroup.FindStringSubmatch(title)
	if len(m) >= 2 {
		return strings.TrimSpace(m[1])
	}
	// fallback — ищем до первого -
	if idx := strings.Index(title, " - "); idx > 0 {
		return strings.TrimSpace(title[:idx])
	}
	return "Unknown"
}

var reResolution = regexp.MustCompile(`(?i)(480p|720p|1080p|1440p|2160p|4k)`)

func extractResolution(title string) string {
	if m := reResolution.FindString(title); m != "" {
		return m
	}
	return "Unknown"
}

func cleanEpisodeName(title string) string {
	// Убираем типичные префиксы и хэши в конце
	title = reGroup.ReplaceAllString(title, "")                           // [Группа]
	title = reResolution.ReplaceAllString(title, "")                      // 1080p
	title = regexp.MustCompile(`(?i)\[.*?\]`).ReplaceAllString(title, "") // [MultiSub], [хэш]
	title = strings.Trim(title, " -[]")
	return strings.TrimSpace(title)
}
