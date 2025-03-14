package storage

import (
	"canciones/internal/models"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type iTunesResponse struct {
	Results []struct {
		TrackName       string  `json:"trackName"`
		ArtistName      string  `json:"artistName"`
		CollectionName  string  `json:"collectionName"`
		TrackTimeMillis int64   `json:"trackTimeMillis"`
		ArtworkUrl100   string  `json:"artworkUrl100"`
		TrackPrice      float64 `json:"trackPrice"`
		Currency        string  `json:"currency"`
		Kind            string  `json:"kind"`
	} `json:"results"`
}

func SearchInITunes(name, artist, album string) ([]map[string]string, error) {
	term := ""

	if name != "" {
		term += name
	}
	if artist != "" {
		if term != "" {
			term += "+"
		}
		term += artist
	}
	if album != "" {
		if term != "" {
			term += "+"
		}
		term += album
	}

	if term == "" {
		return nil, fmt.Errorf("al menos uno de los parámetros (name, artist, album) debe ser proporcionado")
	}

	url := fmt.Sprintf("https://itunes.apple.com/search?term=%s", strings.ReplaceAll(term, " ", "+"))
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result iTunesResponse
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	var songs []map[string]string
	for _, r := range result.Results {
		if r.Kind != "song" {
			continue
		}
		duration := fmt.Sprintf("%02d:%02d", r.TrackTimeMillis/60000, (r.TrackTimeMillis/1000)%60)
		songs = append(songs, map[string]string{
			"name":     r.TrackName,
			"artist":   r.ArtistName,
			"album":    r.CollectionName,
			"duration": duration,
			"artwork":  r.ArtworkUrl100,
			"price":    fmt.Sprintf("%s %.2f", r.Currency, r.TrackPrice),
			"origin":   "apple",
		})
	}

	return songs, nil
}

func SearchInChartLyrics(name, artist string) ([]map[string]string, error) {
	if name == "" && artist == "" {
		return nil, fmt.Errorf("al menos uno de los parámetros (name, artist) debe ser proporcionado para ChartLyrics")
	}

	baseURL := "http://api.chartlyrics.com/apiv1.asmx/SearchLyric"
	params := url.Values{}

	if artist != "" {
		params.Add("artist", artist)
	}
	if name != "" {
		params.Add("song", name)
	}

	url := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	fmt.Println("URL de la solicitud:", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error al realizar la solicitud a ChartLyrics: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer la respuesta de ChartLyrics: %v", err)
	}

	if len(body) == 0 {
		fmt.Println("Respuesta vacía recibida de ChartLyrics.")
		return nil, fmt.Errorf("respuesta vacía de ChartLyrics")
	}

	fmt.Println("Respuesta de ChartLyrics:", string(body))

	var result struct {
		Songs []struct {
			Title  string `xml:"Song"`
			Artist string `xml:"Artist"`
			URL    string `xml:"SongUrl"`
		} `xml:"SearchLyricResult"`
	}
	if err := xml.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error al parsear la respuesta de ChartLyrics: %v", err)
	}

	var songs []map[string]string
	for _, r := range result.Songs {
		songs = append(songs, map[string]string{
			"name":   r.Title,
			"artist": r.Artist,
			"url":    r.URL,
			"origin": "chartlyrics",
		})
	}

	return songs, nil
}

func FindSongsInDB(name, artist, album string) ([]models.Song, error) {
	collection := GetCollection("songs_db", "songs")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	orConditions := []bson.M{}

	if name != "" {
		name = strings.ReplaceAll(name, "+", " ")
		orConditions = append(orConditions, bson.M{"name": bson.M{"$regex": name, "$options": "i"}})
	}
	if artist != "" {
		artist = strings.ReplaceAll(artist, "+", " ")
		orConditions = append(orConditions, bson.M{"artist": bson.M{"$regex": "(?i)" + artist}})
	}
	if album != "" {
		album = strings.ReplaceAll(album, "+", " ")
		orConditions = append(orConditions, bson.M{"album": bson.M{"$regex": "(?i)" + album}})
	}

	var filter bson.M
	if len(orConditions) > 0 {
		filter = bson.M{"$or": orConditions}
	} else {
		filter = bson.M{}
	}

	fmt.Printf("Filtro de búsqueda en la base de datos: %+v\n", filter)

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var songs []models.Song
	if err := cursor.All(ctx, &songs); err != nil {
		return nil, err
	}

	fmt.Printf("Canciones encontradas en la base de datos: %d\n", len(songs))
	return songs, nil
}
