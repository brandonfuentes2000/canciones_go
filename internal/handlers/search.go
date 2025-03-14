package handlers

import (
	"canciones/internal/models"
	"canciones/internal/storage"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	artist := r.URL.Query().Get("artist")
	album := r.URL.Query().Get("album")

	if name == "" && artist == "" && album == "" {
		http.Error(w, "At least one search parameter (name, artist, album) is required", http.StatusBadRequest)
		return
	}

	fmt.Println("Buscando en la base de datos...")
	storedSongs, err := storage.FindSongsInDB(name, artist, album)
	if err != nil {
		fmt.Println("Error al buscar en la base de datos:", err)
		http.Error(w, "Error al buscar en la base de datos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if len(storedSongs) > 0 {
		fmt.Println("Canciones encontradas en la base de datos:", len(storedSongs))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(storedSongs)
		return
	}

	fmt.Println("No se encontraron canciones en la base de datos. Consultando APIs externas...")

	iTunesChan := make(chan []map[string]string)
	chartLyricsChan := make(chan []map[string]string)
	errorChan := make(chan error, 2)

	go func() {
		iTunesResults, err := storage.SearchInITunes(name, artist, album)
		if err != nil {
			errorChan <- fmt.Errorf("Error en iTunes: %v", err)
			iTunesChan <- nil
			return
		}
		iTunesChan <- iTunesResults
	}()

	go func() {
		chartLyricsResults, err := storage.SearchInChartLyrics(name, artist)
		if err != nil {
			errorChan <- fmt.Errorf("Error en ChartLyrics: %v", err)
			chartLyricsChan <- nil
			return
		}

		chartLyricsChan <- chartLyricsResults
	}()

	var allResults []map[string]string
	for i := 0; i < 2; i++ {
		select {
		case iTunesResults := <-iTunesChan:
			if iTunesResults != nil {
				const exchangeRateUSDToGTQ = 7.8
				for _, result := range iTunesResults {
					priceInUSD := result["price"]
					priceParts := strings.Fields(priceInUSD)
					if len(priceParts) != 2 {
						fmt.Printf("Formato inesperado para el precio: '%s'\n", priceInUSD)
						continue
					}

					priceInNumber, err := strconv.ParseFloat(priceParts[1], 64)
					if err != nil {
						fmt.Printf("Error al convertir el precio '%s': %v\n", priceInUSD, err)
						continue
					}
					priceInGTQ := priceInNumber * exchangeRateUSDToGTQ

					song := models.Song{
						Name:     result["name"],
						Artist:   result["artist"],
						Duration: result["duration"],
						Album:    result["album"],
						Artwork:  result["artwork"],
						Price:    fmt.Sprintf("GTQ %.2f", priceInGTQ),
						Origin:   "itunes",
					}

					err = storage.SaveSong(song)
					if err != nil {
						fmt.Printf("Error al guardar la canción '%s': %v\n", song.Name, err)
					} else {
						fmt.Printf("Canción '%s' guardada correctamente en la base de datos.\n", song.Name)
					}
				}
				allResults = append(allResults, iTunesResults...)
			}
		case chartLyricsResults := <-chartLyricsChan:
			if chartLyricsResults != nil {
				for _, result := range chartLyricsResults {
					song := models.Song{
						Name:     result["name"],
						Artist:   result["artist"],
						Duration: "N/A",
						Album:    "N/A",
						Artwork:  "N/A",
						Price:    "N/A",
						Origin:   "chartlyrics",
					}

					err := storage.SaveSong(song)
					if err != nil {
						fmt.Printf("Error al guardar la canción '%s': %v\n", song.Name, err)
					} else {
						fmt.Printf("Canción '%s' guardada correctamente en la base de datos.\n", song.Name)
					}
				}
				allResults = append(allResults, chartLyricsResults...)
			}
		case err := <-errorChan:
			fmt.Println(err)
		}
	}

	if len(allResults) == 0 {
		http.Error(w, "No se encontraron canciones en las APIs externas", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allResults)
}
