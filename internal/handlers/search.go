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

	// Buscar primero en la base de datos
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

	fmt.Println("No se encontraron canciones en la base de datos. Consultando iTunes...")

	searchTerm := ""
	if name != "" {
		searchTerm += name + " "
	}
	if artist != "" {
		searchTerm += artist + " "
	}
	if album != "" {
		searchTerm += album + " "
	}

	iTunesResults, err := storage.SearchInITunes(name, artist, album)

	if err != nil {
		fmt.Println("Error al buscar en iTunes:", err)
		http.Error(w, "Error en buscar en iTunes: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Canciones encontradas en iTunes:", len(iTunesResults))
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(iTunesResults)
}
