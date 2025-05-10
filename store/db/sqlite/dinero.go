package sqlite

import (
    "context"
	"strings"

    "itsfriday/store"
)

func (d *DB) CreateDineroCategory(ctx context.Context, create *store.DineroCategory) (*store.DineroCategory, error) {
    fields := []string{"`user_id`", "`name`", "`priority`"}
	placeholder := []string{"?", "?", "?"}
	args := []any{create.UserID, create.Name, create.Priority}
	stmt := "INSERT INTO expense_category (" + strings.Join(fields, ", ") + ") VALUES (" + strings.Join(placeholder, ", ") + ") RETURNING id"
	if err := d.db.QueryRowContext(ctx, stmt, args...).Scan(
		&create.ID,
	); err != nil {
		return nil, err
	}

	return create, nil
}

func (d *DB) ListDineroCategories(ctx context.Context, find *store.FindDineroCategory) ([]*store.DineroCategory, error) {
	where, args := []string{"1 = 1"}, []any{}

	if v := find.UserID; v != nil {
		where, args = append(where, "user_id = ?"), append(args, *v)
	}
	if v := find.Name; v != nil {
		where, args = append(where, "name = ?"), append(args, *v)
	}

	orderBy := []string{"priority ASC"}
	query := `
		SELECT 
		    id,
			user_id,
			name,
			priority
		FROM expense_category
		WHERE ` + strings.Join(where, " AND ") + ` ORDER BY ` + strings.Join(orderBy, ", ")

	rows, err := d.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]*store.DineroCategory, 0)
	for rows.Next() {
		var category store.DineroCategory
		if err := rows.Scan(
			&category.ID,
			&category.UserID,
			&category.Name,
			&category.Priority,
		); err != nil {
			return nil, err
		}
		list = append(list, &category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (d *DB) UpdateDineroCategory(ctx context.Context, update *store.UpdateDineroCategory) (*store.DineroCategory, error) {
    set, args := []string{}, []any{}
	if v := update.Name; v != nil {
		set, args = append(set, "name = ?"), append(args, *v)
	}
	if v := update.Priority; v != nil {
		set, args = append(set, "priority = ?"), append(args, *v)
	}
	args = append(args, update.ID)

	query := `
		UPDATE expense_category
		SET ` + strings.Join(set, ", ") + `
		WHERE id = ?
		RETURNING id, user_id, name, priority
	`
	category := &store.DineroCategory{}
	if err := d.db.QueryRowContext(ctx, query, args...).Scan(
		&category.ID,
		&category.UserID,
		&category.Name,
		&category.Priority,
	); err != nil {
		return nil, err
	}

	return category, nil
}

func (d *DB) DeleteDineroCategory(ctx context.Context, delete *store.DeleteDineroCategory) error {
    result, err := d.db.ExecContext(ctx, `
		DELETE FROM expense_category WHERE id = ?
	`, delete.ID)
	if err != nil {
		return err
	}
	if _, err := result.RowsAffected(); err != nil {
		return err
	}
	return nil
}
