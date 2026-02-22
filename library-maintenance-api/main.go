package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type QueryRequest struct {
	Query string `json:"query"`
}

type ArtistSearchRequest struct {
	Artist string `json:"artist"`
}

type GenreUpdateRequest struct {
	ArtistID string `json:"artist_id"`
	Genre string `json:"genre"`
}

type TrackResult struct {
	TrackID        int     `json:"track_id"`
	Artist         string  `json:"artist"`
	Album          string  `json:"album"`
	AlbumMBID      string  `json:"album_mbid"`
	Title          string  `json:"title"`
	TrackMBID      string  `json:"track_mbid"`
	YoutubeCode    *string `json:"youtube_code"`
	DownloadStatus string  `json:"download_status"`
	FilePath       *string `json:"file_path"`
}

type ArtistGenreResult struct {
	ArtistID   int			`json:"id"`
	Artist     string		`json:"name"`
	Genres     []string		`json:"genres"`
}

type GenreResult struct {
	GenreID      string		`json:"id"`
	Name      string		`json:"name"`
}

type UpdateRequest struct {
	TrackMBID     string `json:"track_mbid"`
	YoutubeCode   string `json:"youtube_code"`
	RetryDownload bool   `json:"retry_download"`
}

func initDB() *pgxpool.Pool {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found (ok if running in container)")
	}

	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	db := os.Getenv("POSTGRES_DB")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		user,
		password,
		host,
		port,
		db,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}

	log.Println("✅ Database connection established")

	return pool
}

func main() {
	dbPool := initDB()
	defer dbPool.Close()

	r := gin.Default()

	r.Use(cors.Default())

	r.POST("/youtube/search", func(c *gin.Context) {
		var req ArtistSearchRequest

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		query := `
		SELECT 
			t.id AS track_id,
			a.name AS artist,
			al.title AS album,
			al.mbid AS album_mbid,
			t.title,
			t.mbid AS track_mbid,
			t.youtube_code,
			t.download_status,
			t.file_path
		FROM tracks t 
		JOIN album_tracks alt ON alt.track_id = t.id
		JOIN albums al ON al.id = alt.album_id
		JOIN artist_tracks art ON art.track_id = t.id
		JOIN artists a ON art.artist_id = a.id
		WHERE LOWER(a.name) LIKE LOWER($1)
		ORDER BY a.name, al.title, t.title;
		`

		rows, err := dbPool.Query(context.Background(), query, "%"+req.Artist+"%")
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		defer rows.Close()

		results := []TrackResult{}

		for rows.Next() {
			var r TrackResult

			err := rows.Scan(
				&r.TrackID,
				&r.Artist,
				&r.Album,
				&r.AlbumMBID,
				&r.Title,
				&r.TrackMBID,
				&r.YoutubeCode,
				&r.DownloadStatus,
				&r.FilePath,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}

			results = append(results, r)
		}

		c.JSON(http.StatusOK, results)
	})

	r.POST("/youtube/retry", func(c *gin.Context) {
		var req UpdateRequest

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		tx, err := dbPool.Begin(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		defer tx.Rollback(context.Background())

		updateQuery := `
		UPDATE tracks
		SET youtube_code = $1,
			download_status = CASE
				WHEN $2 = true THEN 'queued'
				ELSE download_status
			END
		WHERE mbid = $3;
		`

		_, err = tx.Exec(
			context.Background(),
			updateQuery,
			req.YoutubeCode,
			req.RetryDownload,
			req.TrackMBID,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		if err := tx.Commit(context.Background()); err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "updated",
		})
	})

	r.POST("/artistgenres", func(c *gin.Context) {
		var req ArtistSearchRequest

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		tx, err := dbPool.Begin(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		defer tx.Rollback(context.Background())

		searchQuery := `
		SELECT 
			a.id AS artist_id,
			a.name,
			COALESCE(
				ARRAY_AGG(DISTINCT g.name ORDER BY g.name)
					FILTER (WHERE g.name IS NOT NULL),
				ARRAY[]::text[]
			) AS genres
		FROM artists a
		LEFT JOIN artist_genres ag ON ag.artist_id = a.id
		LEFT JOIN genres g ON g.id = ag.genre_id
		WHERE LOWER(a.name) LIKE LOWER($1)
		GROUP BY a.id, a.name
		ORDER BY a.name;
		`

		rows, err := dbPool.Query(context.Background(), searchQuery, "%"+req.Artist+"%")
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		defer rows.Close()

		results := []ArtistGenreResult{}

		for rows.Next() {
			var r ArtistGenreResult

			err := rows.Scan(
				&r.ArtistID,
				&r.Artist,
				&r.Genres,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}

			results = append(results, r)
		}

		c.JSON(http.StatusOK, results)
	})

	r.POST("/artistgenres/all", func(c *gin.Context) {
		tx, err := dbPool.Begin(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		defer tx.Rollback(context.Background())

		searchQuery := `
		SELECT 
			a.id AS artist_id,
			a.name,
			COALESCE(
				ARRAY_AGG(DISTINCT g.name ORDER BY g.name)
					FILTER (WHERE g.name IS NOT NULL),
				ARRAY[]::text[]
			) AS genres
		FROM artists a
		LEFT JOIN artist_genres ag ON ag.artist_id = a.id
		LEFT JOIN genres g ON g.id = ag.genre_id
		GROUP BY a.id, a.name
		ORDER BY a.name;
		`

		rows, err := dbPool.Query(context.Background(), searchQuery)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		defer rows.Close()

		results := []ArtistGenreResult{}

		for rows.Next() {
			var r ArtistGenreResult

			err := rows.Scan(
				&r.ArtistID,
				&r.Artist,
				&r.Genres,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}

			results = append(results, r)
		}

		c.JSON(http.StatusOK, results)
	})

	r.POST("/artistgenres/all/nogenre", func(c *gin.Context) {
		tx, err := dbPool.Begin(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		defer tx.Rollback(context.Background())

		searchQuery := `
		SELECT 
			a.id AS artist_id,
			a.name,
			ARRAY[]::text[] AS genres
		FROM artists a
		WHERE NOT EXISTS (
			SELECT 1
			FROM artist_genres ag
			WHERE ag.artist_id = a.id
		)
		ORDER BY a.name;
		`

		rows, err := dbPool.Query(context.Background(), searchQuery)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		defer rows.Close()

		results := []ArtistGenreResult{}

		for rows.Next() {
			var r ArtistGenreResult

			err := rows.Scan(
				&r.ArtistID,
				&r.Artist,
				&r.Genres,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}

			results = append(results, r)
		}

		c.JSON(http.StatusOK, results)
	})

	r.POST("/artistgenres/deletebyname", func(c *gin.Context) {
		var req GenreUpdateRequest

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		tx, err := dbPool.Begin(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		defer tx.Rollback(context.Background())

		deleteQuery := `
		DELETE FROM artist_genres ag
		USING genres g
		WHERE ag.artist_id = $1
		AND g.name = $2
		AND ag.genre_id = g.id;
		`

		_, err = tx.Exec(
			context.Background(),
			deleteQuery,
			req.ArtistID,
			req.Genre,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		if err := tx.Commit(context.Background()); err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "updated",
		})

	})

	r.POST("/artistgenres/addbyname", func(c *gin.Context) {
		var req GenreUpdateRequest

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		tx, err := dbPool.Begin(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		defer tx.Rollback(context.Background())

		addQuery := `
		WITH new_genre AS (
			INSERT INTO genres (name)
			VALUES ($2)
			ON CONFLICT (name) DO UPDATE
				SET name = EXCLUDED.name
			RETURNING id
		)
		INSERT INTO artist_genres (artist_id, genre_id)
		SELECT $1, id
		FROM new_genre
		ON CONFLICT DO NOTHING;
		`

		_, err = tx.Exec(
			context.Background(),
			addQuery,
			req.ArtistID,
			req.Genre,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		if err := tx.Commit(context.Background()); err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "updated",
		})
	})

	r.POST("/genres", func(c *gin.Context) {
		tx, err := dbPool.Begin(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		defer tx.Rollback(context.Background())

		searchQuery := `
		SELECT id, name
		FROM genres
		`

		rows, err := dbPool.Query(context.Background(), searchQuery)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		defer rows.Close()

		results := []GenreResult{}

		for rows.Next() {
			var r GenreResult

			err := rows.Scan(
				&r.GenreID,
				&r.Name,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}

			results = append(results, r)
		}

		c.JSON(http.StatusOK, results)
	})

	r.Run(":5001")
}
