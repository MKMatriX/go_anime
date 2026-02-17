// services/anidb_search.go
package services

import (
	"compress/gzip"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	xmlGzURL    = "https://anidb.net/api/anime-titles.xml.gz"
	cacheFile   = "anime-titles.xml"
	maxAgeHours = 24 // обновлять кэш не чаще раза в сутки
	userAgent   = "go_anime/0.1 (matrix-elf@yandex.ru)"
)

// TitleEntry — одна запись из anime-titles
type TitleEntry struct {
	AID   int
	Type  string // "main", "official", "syn", "short", "kana" и т.д.
	Lang  string
	Title string
}

// TitleList — промежуточная структура для xml-декодирования
type TitleList struct {
	XMLName xml.Name `xml:"animetitles"`
	Animes  []struct {
		AID    int `xml:"aid,attr"`
		Titles []struct {
			Type   string `xml:"type,attr"`
			Lang   string `xml:"lang,attr"`                                                // xml:lang
			NSLang string `xml:"http://www.w3.org/XML/1998/namespace lang,attr,omitempty"` // на случай если пространство имён явно
			Text   string `xml:",chardata"`
		} `xml:"title"`
	} `xml:"anime"`
}

// GetAniDBId ищет наиболее вероятный AniDB AID по названию аниме
// Возвращает AID (>0) или 0, если ничего не найдено
func GetAniDBId(animeName string) (int, error) {
	query := strings.TrimSpace(animeName)
	if query == "" {
		return 0, fmt.Errorf("AniDB error: empty search query")
	}

	entries, err := loadOrDownloadTitles()
	if err != nil {
		return 0, err
	}

	matches := search(entries, query)
	if len(matches) == 0 {
		return 0, fmt.Errorf("AniDB error: anime not found")
	}

	// // Сортируем: приоритет main (type=1), затем более свежие (больший AID)
	// sort.Slice(matches, func(i, j int) bool {
	// 	if matches[i].Type != matches[j].Type {
	// 		return matches[i].Type < matches[j].Type // main выше
	// 	}
	// 	return matches[i].AID > matches[j].AID
	// })

	// Возвращаем самый приоритетный вариант
	return matches[0].AID, nil
}

func loadOrDownloadTitles() ([]TitleEntry, error) {
	cachePath := filepath.Join(os.TempDir(), cacheFile)

	// Если кэш свежий — используем его
	if fi, err := os.Stat(cachePath); err == nil {
		if time.Since(fi.ModTime()) < maxAgeHours*time.Hour {
			f, err := os.Open(cachePath)
			if err == nil {
				defer f.Close()
				return parseAnimeTitles(f)
			}
		}
	}

	// Скачиваем
	resp, err := httpGetWithUA(xmlGzURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status %d", resp.StatusCode)
	}

	gz, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, err
	}
	defer gz.Close()

	// Сохраняем во временный файл
	tmpPath := cachePath + ".tmp"
	tmp, err := os.Create(tmpPath)
	if err != nil {
		return nil, err
	}
	defer tmp.Close()

	if _, err := io.Copy(tmp, gz); err != nil {
		return nil, err
	}
	tmp.Close()

	// Атомарное перемещение
	if err := os.Rename(tmpPath, cachePath); err != nil {
		return nil, err
	}

	f, err := os.Open(cachePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return parseAnimeTitles(f)
}
func parseAnimeTitles(r io.Reader) ([]TitleEntry, error) {
	var titles []TitleEntry

	// Сначала пробуем декодировать как полноценный XML
	decoder := xml.NewDecoder(r)
	decoder.Strict = false // AniDB иногда кладёт некорректные символы

	var list TitleList
	if err := decoder.Decode(&list); err == nil && len(list.Animes) > 0 {
		// Успешно распарсили как XML
		for _, anime := range list.Animes {
			for _, t := range anime.Titles {
				lang := t.Lang
				if lang == "" {
					lang = t.NSLang // запасной вариант
				}
				titleText := strings.TrimSpace(t.Text)
				if titleText == "" {
					continue
				}

				titles = append(titles, TitleEntry{
					AID:   anime.AID,
					Type:  t.Type,
					Lang:  lang,
					Title: titleText,
				})
			}
		}
	} else {
		return nil, fmt.Errorf("xml decode failed: %w", err)
	}

	if len(titles) == 0 {
		return nil, fmt.Errorf("не найдено ни одного тайтла в файле")
	}

	return titles, nil
}

func search(entries []TitleEntry, q string) []TitleEntry {
	qLower := strings.ToLower(q)
	var matches []TitleEntry

	for _, e := range entries {
		if strings.Contains(strings.ToLower(e.Title), qLower) {
			matches = append(matches, e)
		}
	}

	return matches
}

func httpGetWithUA(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)

	client := &http.Client{Timeout: 45 * time.Second}
	return client.Do(req)
}
