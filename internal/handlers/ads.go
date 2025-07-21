package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/plusha-fullstack/marketplaceAPI/internal/models"
	"github.com/plusha-fullstack/marketplaceAPI/internal/utils"
)

type AdHandler struct {
	db *sql.DB
}

func NewAdHandler(db *sql.DB) *AdHandler {
	return &AdHandler{db: db}
}

func (h *AdHandler) CreateAd(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(UserContextKey).(*utils.Claims)

	var req models.CreateAdRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := utils.ValidateAdTitle(req.Title); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := utils.ValidateAdDescription(req.Description); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := utils.ValidateImageURL(req.ImageURL); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := utils.ValidatePrice(req.Price); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var ad models.Ad
	err := h.db.QueryRow(`
        INSERT INTO ads (title, description, image_url, price, author_id, created_at) 
        VALUES ($1, $2, $3, $4, $5, NOW()) 
        RETURNING id, title, description, image_url, price, author_id, created_at`,
		req.Title, req.Description, req.ImageURL, req.Price, claims.UserID,
	).Scan(&ad.ID, &ad.Title, &ad.Description, &ad.ImageURL, &ad.Price, &ad.AuthorID, &ad.CreatedAt)

	if err != nil {
		http.Error(w, "Error creating ad", http.StatusInternalServerError)
		return
	}

	ad.AuthorLogin = claims.Login
	ad.IsOwner = true

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ad)
}

func (h *AdHandler) GetAds(w http.ResponseWriter, r *http.Request) {
	filters := parseFilters(r)

	var currentUserID int
	if claims := r.Context().Value(UserContextKey); claims != nil {
		currentUserID = claims.(*utils.Claims).UserID
	}

	query := `
        SELECT a.id, a.title, a.description, a.image_url, a.price, 
               a.author_id, u.login, a.created_at
        FROM ads a
        JOIN users u ON a.author_id = u.id
        WHERE 1=1`

	args := []interface{}{}
	argCount := 0

	if filters.MinPrice > 0 {
		argCount++
		query += " AND a.price >= $" + strconv.Itoa(argCount)
		args = append(args, filters.MinPrice)
	}

	if filters.MaxPrice > 0 {
		argCount++
		query += " AND a.price <= $" + strconv.Itoa(argCount)
		args = append(args, filters.MaxPrice)
	}

	switch filters.SortBy {
	case "price":
		query += " ORDER BY a.price"
	default:
		query += " ORDER BY a.created_at"
	}

	if strings.ToLower(filters.SortOrder) == "asc" {
		query += " ASC"
	} else {
		query += " DESC"
	}

	argCount++
	query += " LIMIT $" + strconv.Itoa(argCount)
	args = append(args, filters.Limit)

	argCount++
	query += " OFFSET $" + strconv.Itoa(argCount)
	args = append(args, (filters.Page-1)*filters.Limit)

	rows, err := h.db.Query(query, args...)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var ads []models.Ad
	for rows.Next() {
		var ad models.Ad
		err := rows.Scan(
			&ad.ID, &ad.Title, &ad.Description, &ad.ImageURL, &ad.Price,
			&ad.AuthorID, &ad.AuthorLogin, &ad.CreatedAt,
		)
		if err != nil {
			http.Error(w, "Error scanning ad", http.StatusInternalServerError)
			return
		}

		if currentUserID > 0 && ad.AuthorID == currentUserID {
			ad.IsOwner = true
		}

		ads = append(ads, ad)
	}

	if ads == nil {
		ads = []models.Ad{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ads)
}

func parseFilters(r *http.Request) models.AdFilters {
	filters := models.AdFilters{
		Page:      1,
		Limit:     10,
		SortBy:    "date",
		SortOrder: "desc",
	}

	if page, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && page > 0 {
		filters.Page = page
	}

	if limit, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && limit > 0 && limit <= 100 {
		filters.Limit = limit
	}

	if sortBy := r.URL.Query().Get("sort_by"); sortBy == "price" || sortBy == "date" {
		filters.SortBy = sortBy
	}

	if sortOrder := r.URL.Query().Get("sort_order"); sortOrder == "asc" || sortOrder == "desc" {
		filters.SortOrder = sortOrder
	}

	if minPrice, err := strconv.ParseFloat(r.URL.Query().Get("min_price"), 64); err == nil && minPrice >= 0 {
		filters.MinPrice = minPrice
	}

	if maxPrice, err := strconv.ParseFloat(r.URL.Query().Get("max_price"), 64); err == nil && maxPrice >= 0 {
		filters.MaxPrice = maxPrice
	}

	return filters
}
