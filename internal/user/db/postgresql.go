package db

import (
	"context"
	"fmt"
	"reflect"
	"song-lib/internal/user"
	"song-lib/pkg/logging"
	"strings"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	client *pgxpool.Pool
	logger *logging.Logger
}

func NewRepository(client *pgxpool.Pool, logger *logging.Logger) user.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", "")
}

func (r *repository) Create(ctx context.Context, s *user.Song) error {
	q := `INSERT INTO library("group", song, releaseDate, text, link) 
		 VALUES($1, $2, $3, $4, $5) 
		 RETURNING id`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", formatQuery(q)))
	err := r.client.QueryRow(ctx, q, s.Group, s.Song, s.ReleaseDate, s.Text, s.Link).Scan(&s.ID)
	if err != nil {
		r.logger.Error(err)
		return err
	}


	return nil
}

func (r *repository) FindAll(ctx context.Context) ([]user.Song, error) {
	q := `SELECT id, "group", song FROM library`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", formatQuery(q)))

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		r.logger.Error(err)
		return nil, err
	}
	
	songs := make([]user.Song, 0)

	for rows.Next() {
		var song user.Song

		err = rows.Scan(&song.ID, &song.Group, &song.Song)
		if err != nil{
			r.logger.Error(err)
			return nil, err
		}

		songs = append(songs, song)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error(err)
		return nil, err
	}

	return songs, nil
}

func (r *repository) FindWithFilter(ctx context.Context, s *user.Song) ([]user.Song, error) {
	val := reflect.ValueOf(*s)
	var param string
	var value string
	fmt.Println(val.Kind())

	if val.Kind() == reflect.Struct {
		for i := 0; i < val.NumField(); i++ {
			if val.Field(i).String() != ""{
				param = strings.ToLower(val.Type().Field(i).Name)
				value = val.Field(i).String()
			}
		}
		q := fmt.Sprintf(`SELECT song, "group" FROM library WHERE "%s" = $1`, param)
		r.logger.Trace(fmt.Sprintf("SQL query: %s", formatQuery(q)))

		rows, err := r.client.Query(ctx, q, value)
		if err != nil {
			r.logger.Error(err)
			return nil, err
		}

		songs := make([]user.Song, 0)	

		for rows.Next() {
			var song user.Song

			if err = rows.Scan(&song.Song, &song.Group); err != nil {
				r.logger.Error(err)
				return nil, err
			}

			songs = append(songs, song)
		}

		if err = rows.Err(); err != nil {
			r.logger.Error(err)
			return nil, err
		}
		
		return songs, nil

	} else {
		return nil, fmt.Errorf("transferred object isn't a struct")
	}


}

func (r *repository) FindOne(ctx context.Context, s *user.Song) (*user.Song, error) {
    q := `SELECT releaseDate, text, link FROM library WHERE "group" = $1 and song = $2`
    r.logger.Trace(fmt.Sprintf("SQL query: %s", formatQuery(q)))

    rows, err := r.client.Query(ctx, q, s.Group, s.Song)
    if err != nil {
        r.logger.Error(err)
        return nil, err
    }
    defer rows.Close() 

    if !rows.Next() {
        return nil, fmt.Errorf("song not found")
    }

    err = rows.Scan(&s.ReleaseDate, &s.Text, &s.Link)
    if err != nil {
        r.logger.Error(err)
        return nil, err 
    }

    return s, nil
}

func (r *repository) Delete(ctx context.Context, s *user.Song) error {
	q := `DELETE FROM library WHERE "group" = $1 and song = $2`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", formatQuery(q)))

	_, err := r.client.Exec(ctx, q, s.Group, s.Song)
	if err != nil {
		r.logger.Error(err)
		return err
	}

	return nil
}

func (r *repository) Update(ctx context.Context, s *user.Song) error {
	val := reflect.ValueOf(s).Elem()
	tp := val.Type()
	if val.Kind() == reflect.Struct {
		var q string
		for i := 0; i < val.NumField(); i++ {
			if val.Field(i).String() != "" && strings.ToLower(tp.Field(i).Name) != "song" && strings.ToLower(tp.Field(i).Name) != "group"{
				q = fmt.Sprintf(`UPDATE library SET %s = $1 WHERE song = $2 AND "group" = $3`, strings.ToLower(tp.Field(i).Name))
				r.logger.Trace(fmt.Sprintf("SQL query: %s", formatQuery(q)))
				
				_, err := r.client.Exec(ctx, q, val.Field(i).String(), s.Song, s.Group)
				if err != nil {
					r.logger.Error(err)
					return err
				}
			}
		}
	} else {
		return fmt.Errorf("transferred object isn't a struct")
	}

	return nil
}