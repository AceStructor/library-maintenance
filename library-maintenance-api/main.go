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

type SearchRequest struct {
	Artist string `json:"artist"`
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
		var req SearchRequest

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

	r.Run(":8080")
}
