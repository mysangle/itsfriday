package sqlite

import (
	"context"
	"fmt"
	"strings"

	"itsfriday/store"
)

// book

func (d *DB) CreateBook(ctx context.Context, create *store.Book) (*store.Book, error) {
    fields := []string{"`user_id`", "`title`", "`author`", "`translator`", "`pages`", "`pub_year`", "`genre`"}
	placeholder := []string{"?", "?", "?", "?", "?", "?", "?"}
	args := []any{create.UserID, create.Title, create.Author, create.Translator, create.Pages, create.PubYear, create.Genre}
	stmt := "INSERT INTO book (" + strings.Join(fields, ", ") + ") VALUES (" + strings.Join(placeholder, ", ") + ") RETURNING id, created_ts"
	if err := d.db.QueryRowContext(ctx, stmt, args...).Scan(
		&create.ID,
		&create.CreatedTs,
	); err != nil {
		return nil, err
	}

	return create, nil
}

func (d *DB) UpdateBook(ctx context.Context, update *store.UpdateBook) (*store.Book, error) {
    set, args := []string{}, []any{}
	if v := update.UserID; v != nil {
		set, args = append(set, "user_id = ?"), append(args, *v)
	}
	if v := update.Title; v != nil {
		set, args = append(set, "title = ?"), append(args, *v)
	}
	if v := update.Author; v != nil {
		set, args = append(set, "author = ?"), append(args, *v)
	}
	if v := update.Translator; v != nil {
		set, args = append(set, "translator = ?"), append(args, *v)
	}
	if v := update.Pages; v != nil {
		set, args = append(set, "pages = ?"), append(args, *v)
	}
	if v := update.PubYear; v != nil {
		set, args = append(set, "pub_year = ?"), append(args, *v)
	}
	if v := update.Genre; v != nil {
		set, args = append(set, "genre = ?"), append(args, *v)
	}
	args = append(args, update.ID)

	query := `
		UPDATE book
		SET ` + strings.Join(set, ", ") + `
		WHERE id = ?
		RETURNING id, user_id, title, author, translator, pages, pub_year, genre, created_ts
	`
	book := &store.Book{}
	if err := d.db.QueryRowContext(ctx, query, args...).Scan(
		&book.ID,
		&book.UserID,
		&book.Title,
		&book.Author,
		&book.Translator,
		&book.Pages,
		&book.PubYear,
		&book.Genre,
		&book.CreatedTs,
	); err != nil {
		return nil, err
	}

	return book, nil
}

func (d *DB) ListBooks(ctx context.Context, find *store.FindBook) ([]*store.Book, error) {
    where, args := []string{"1 = 1"}, []any{}

	if v := find.ID; v != nil {
		where, args = append(where, "id = ?"), append(args, *v)
	}
	if v := find.UserID; v != nil {
		where, args = append(where, "user_id = ?"), append(args, *v)
	}
	if v := find.Title; v != nil {
		where, args = append(where, "title = ?"), append(args, *v)
	}
	if v := find.Author; v != nil {
		where, args = append(where, "author = ?"), append(args, *v)
	}

	orderBy := []string{"created_ts DESC"}
	query := `
		SELECT 
			id,
			user_id,
			title,
			author,
			translator,
			pages,
			pub_year,
			genre,
			created_ts
		FROM book
		WHERE ` + strings.Join(where, " AND ") + ` ORDER BY ` + strings.Join(orderBy, ", ")
	if v := find.Limit; v != nil {
		query += fmt.Sprintf(" LIMIT %d", *v)
	}

	rows, err := d.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]*store.Book, 0)
	for rows.Next() {
		var book store.Book
		if err := rows.Scan(
			&book.ID,
			&book.UserID,
			&book.Title,
			&book.Author,
			&book.Translator,
			&book.Pages,
			&book.PubYear,
			&book.Genre,
			&book.CreatedTs,
		); err != nil {
			return nil, err
		}
		list = append(list, &book)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (d *DB) DeleteBook(ctx context.Context, delete *store.DeleteBook) error {
    result, err := d.db.ExecContext(ctx, `
		DELETE FROM book WHERE id = ?
	`, delete.ID)
	if err != nil {
		return err
	}
	if _, err := result.RowsAffected(); err != nil {
		return err
	}
	return nil
}

// book_review

func (d *DB) CreateBookReview(ctx context.Context, create *store.BookReview) (*store.BookReview, error) {
    fields := []string{"`user_id`", "`book_id`", "`date_read`", "`rating`", "`review`"}
	placeholder := []string{"?", "?", "?", "?", "?"}
	args := []any{create.UserID, create.BookID, create.DateRead, create.Rating, create.Review}
	stmt := "INSERT INTO book_review (" + strings.Join(fields, ", ") + ") VALUES (" + strings.Join(placeholder, ", ") + ") RETURNING id, created_ts"
	if err := d.db.QueryRowContext(ctx, stmt, args...).Scan(
		&create.ID,
		&create.CreatedTs,
	); err != nil {
		return nil, err
	}

	return create, nil
}

func (d *DB) UpdateBookReview(ctx context.Context, update *store.UpdateBookReview) (*store.BookReview, error) {
    set, args := []string{}, []any{}
	if v := update.UserID; v != nil {
		set, args = append(set, "user_id = ?"), append(args, *v)
	}
	if v := update.BookID; v != nil {
		set, args = append(set, "book_id = ?"), append(args, *v)
	}
	if v := update.DateRead; v != nil {
		set, args = append(set, "date_read = ?"), append(args, *v)
	}
	if v := update.Rating; v != nil {
		set, args = append(set, "rating = ?"), append(args, *v)
	}
	if v := update.Review; v != nil {
		set, args = append(set, "review = ?"), append(args, *v)
	}
	args = append(args, update.ID)

	query := `
		UPDATE book_review
		SET ` + strings.Join(set, ", ") + `
		WHERE id = ?
		RETURNING id, user_id, book_id, date_read, rating, review, created_ts
	`
	bookReview := &store.BookReview{}
	if err := d.db.QueryRowContext(ctx, query, args...).Scan(
		&bookReview.ID,
		&bookReview.UserID,
		&bookReview.BookID,
		&bookReview.DateRead,
		&bookReview.Rating,
		&bookReview.Review,
		&bookReview.CreatedTs,
	); err != nil {
		return nil, err
	}

	return bookReview, nil
}

func (d *DB) ListBookReviews(ctx context.Context, find *store.FindBookReview) ([]*store.BookReview, error) {
    where, args := []string{"1 = 1"}, []any{}

	if v := find.ID; v != nil {
		where, args = append(where, "id = ?"), append(args, *v)
	}
	if v := find.UserID; v != nil {
		where, args = append(where, "user_id = ?"), append(args, *v)
	}
	if v := find.BookID; v != nil {
		where, args = append(where, "book_id = ?"), append(args, *v)
	}
	if v := find.DateRead; v != nil {
		where, args = append(where, "date_read = ?"), append(args, *v)
	}
	if v := find.Rating; v != nil {
		where, args = append(where, "rating = ?"), append(args, *v)
	}

	orderBy := []string{"created_ts DESC"}
	query := `
		SELECT 
			id,
			user_id,
			book_id,
			date_read,
			rating,
			review,
			created_ts
		FROM book_review
		WHERE ` + strings.Join(where, " AND ") + ` ORDER BY ` + strings.Join(orderBy, ", ")
	if v := find.Limit; v != nil {
		query += fmt.Sprintf(" LIMIT %d", *v)
	}

	rows, err := d.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]*store.BookReview, 0)
	for rows.Next() {
		var bookReview store.BookReview
		if err := rows.Scan(
			&bookReview.ID,
			&bookReview.UserID,
			&bookReview.BookID,
			&bookReview.DateRead,
			&bookReview.Rating,
			&bookReview.Review,
			&bookReview.CreatedTs,
		); err != nil {
			return nil, err
		}
		list = append(list, &bookReview)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (d *DB) DeleteBookReview(ctx context.Context, delete *store.DeleteBookReview) error {
    result, err := d.db.ExecContext(ctx, `
		DELETE FROM book_review WHERE id = ?
	`, delete.ID)
	if err != nil {
		return err
	}
	if _, err := result.RowsAffected(); err != nil {
		return err
	}
	return nil
}

func (d *DB) ListBooksReadInYear(ctx context.Context, userID int32, year int32) ([]*store.BookRead, error) {
    where, args := []string{"1 = 1"}, []any{}

	where, args = append(where, "`book_review`.`user_id` = ?"), append(args, userID)
	start := fmt.Sprintf("%04d-01-01", year)
	where, args = append(where, "`book_review`.`date_read` >= ?"), append(args, start)
	end := fmt.Sprintf("%04d-12-31", year)
	where, args = append(where, "`book_review`.`date_read` <= ?"), append(args, end)

	orderBy := []string{"date_read ASC"}
	fields := []string{
		"`book`.`id` AS `book_id`",
		"`book_review`.`id` AS `review_id`",
		"`book`.`title` AS `title`",
		"`book`.`author` AS `author`",
		"`book`.`translator` AS `translator`",
		"`book`.`pages` AS `pages`",
		"`book`.`pub_year` AS `pub_year`",
		"`book`.`genre` AS `genre`",
		"`book_review`.`date_read` AS `date_read`",
		"`book_review`.`rating` AS `rating`",
		"`book_review`.`review` AS `review`",
	}
	query := "SELECT " + strings.Join(fields, ", ") + "FROM `book` " +
		"LEFT JOIN `book_review` ON `book`.`id` = `book_review`.`book_id`" +
		"WHERE " + strings.Join(where, " AND ") + " " +
		"ORDER BY " + strings.Join(orderBy, ", ")

	rows, err := d.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]*store.BookRead, 0)
	for rows.Next() {
		var bookRead store.BookRead
		if err := rows.Scan(
			&bookRead.BookID,
			&bookRead.ReviewID,
			&bookRead.Title,
			&bookRead.Author,
			&bookRead.Translator,
			&bookRead.Pages,
			&bookRead.PubYear,
			&bookRead.Genre,
			&bookRead.DateRead,
			&bookRead.Rating,
			&bookRead.Review,
		); err != nil {
			return nil, err
		}
		list = append(list, &bookRead)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (d *DB) ReportBook(ctx context.Context, userID int32) ([]*store.ReportBook, error) {
    where, args := []string{"1 = 1"}, []any{}

	where, args = append(where, "`book_review`.`user_id` = ?"), append(args, userID)

	orderBy := []string{"date_read ASC"}
	query := `
		SELECT
		    strftime('%Y', date_read) AS year,
			count(*) AS count
		FROM book_review
		WHERE ` + strings.Join(where, " AND ") +
		` GROUP BY strftime('%Y', date_read)` +
		` ORDER BY ` + strings.Join(orderBy, ", ")

	rows, err := d.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]*store.ReportBook, 0)
	for rows.Next() {
		var reportBook store.ReportBook
		if err := rows.Scan(
			&reportBook.Year,
			&reportBook.Count,
		); err != nil {
			return nil, err
		}
		list = append(list, &reportBook)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}
