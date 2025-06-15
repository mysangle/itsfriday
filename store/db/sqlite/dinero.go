package sqlite

import (
    "context"
	"fmt"
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

	if v := find.ID; v != nil {
		where, args = append(where, "id = ?"), append(args, *v)
	}
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

func (d *DB) CreateDineroExpense(ctx context.Context, create *store.DineroExpense) (*store.DineroExpense, error) {
    fields := []string{"`user_id`", "`category_id`", "`date_used`", "`item`", "`price`"}
	placeholder := []string{"?", "?", "?", "?", "?"}
	args := []any{create.UserID, create.CategoryID, create.DateUsed, create.Item, create.Price}
	stmt := "INSERT INTO expense (" + strings.Join(fields, ", ") + ") VALUES (" + strings.Join(placeholder, ", ") + ") RETURNING id, created_ts"
	if err := d.db.QueryRowContext(ctx, stmt, args...).Scan(
		&create.ID,
		&create.CreatedTs,
	); err != nil {
		return nil, err
	}

	return create, nil
}

func (d *DB) ListDineroExpenses(ctx context.Context, find *store.FindDineroExpense) ([]*store.DineroExpense, error) {
	where, args := []string{"1 = 1"}, []any{}

	if v := find.ID; v != nil {
		where, args = append(where, "id = ?"), append(args, *v)
	}
	if v := find.UserID; v != nil {
		where, args = append(where, "user_id = ?"), append(args, *v)
	}
	if v := find.CategoryID; v != nil {
		where, args = append(where, "category_id = ?"), append(args, *v)
	}
	if find.Year != nil && find.Month != nil {
		start := fmt.Sprintf("%04d-%02d-01", *find.Year, *find.Month)
		where, args = append(where, "date_used >= ?"), append(args, start)
		end := fmt.Sprintf("%04d-%02d-01", *find.Year, *find.Month + 1)
		where, args = append(where, "date_used < ?"), append(args, end)
	}

	orderBy := []string{"date_used ASC"}
	query := `
		SELECT 
		    id,
			user_id,
			category_id,
			date_used,
			item,
			price,
			created_ts
		FROM expense
		WHERE ` + strings.Join(where, " AND ") + ` ORDER BY ` + strings.Join(orderBy, ", ")

	rows, err := d.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]*store.DineroExpense, 0)
	for rows.Next() {
		var expense store.DineroExpense
		if err := rows.Scan(
			&expense.ID,
			&expense.UserID,
			&expense.CategoryID,
			&expense.DateUsed,
			&expense.Item,
			&expense.Price,
			&expense.CreatedTs,
		); err != nil {
			return nil, err
		}
		list = append(list, &expense)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (d *DB) UpdateDineroExpense(ctx context.Context, update *store.UpdateDineroExpense) (*store.DineroExpense, error) {
    set, args := []string{}, []any{}
	if v := update.CategoryID; v != nil {
		set, args = append(set, "category_id = ?"), append(args, *v)
	}
	if v := update.DateUsed; v != nil {
		set, args = append(set, "date_used = ?"), append(args, *v)
	}
	if v := update.Item; v != nil {
		set, args = append(set, "item = ?"), append(args, *v)
	}
	if v := update.Price; v != nil {
		set, args = append(set, "price = ?"), append(args, *v)
	}
	args = append(args, update.ID)

	query := `
		UPDATE expense
		SET ` + strings.Join(set, ", ") + `
		WHERE id = ?
		RETURNING id, user_id, category_id, date_used, item, price, created_ts
	`
	expense := &store.DineroExpense{}
	if err := d.db.QueryRowContext(ctx, query, args...).Scan(
		&expense.ID,
		&expense.UserID,
		&expense.CategoryID,
		&expense.DateUsed,
		&expense.Item,
		&expense.Price,
		&expense.CreatedTs,
	); err != nil {
		return nil, err
	}

	return expense, nil
}

func (d *DB) DeleteDineroExpense(ctx context.Context, delete *store.DeleteDineroExpense) error {
    result, err := d.db.ExecContext(ctx, `
		DELETE FROM expense WHERE id = ?
	`, delete.ID)
	if err != nil {
		return err
	}
	if _, err := result.RowsAffected(); err != nil {
		return err
	}
	return nil
}

func (d *DB) GetTotalCostByCategory(ctx context.Context, find *store.FindDineroExpense) ([]*store.TotalCostPerCategory, error) {
	where, args := []string{"1 = 1"}, []any{}
	if v := find.UserID; v != nil {
		where, args = append(where, "`expense`.`user_id` = ?"), append(args, *v)
	}
	if find.Year != nil && find.Month != nil {
		start := fmt.Sprintf("%04d-%02d-01", *find.Year, *find.Month)
		where, args = append(where, "`expense`.`date_used` >= ?"), append(args, start)
		end := fmt.Sprintf("%04d-%02d-01", *find.Year, *find.Month + 1)
		where, args = append(where, "`expense`.`date_used` < ?"), append(args, end)
	}

	orderBy := []string{"`expense_category`.`priority` ASC"}
	fields := []string{
		"`expense_category`.`name` AS `name`",
		"sum(`expense`.`price`) AS `cost`",
	}
	query := "SELECT " + strings.Join(fields, ", ") + "FROM `expense_category` " +
		" LEFT JOIN `expense` ON `expense`.`category_id` = `expense_category`.`id`" +
		" WHERE " + strings.Join(where, " AND ") + " " +
		" GROUP BY `expense_category`.`id`" +
		" ORDER BY " + strings.Join(orderBy, ", ")

	result := append(args, args...)
	rows, err := d.db.QueryContext(ctx, query, result...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]*store.TotalCostPerCategory, 0)
	for rows.Next() {
		var totalCost store.TotalCostPerCategory
		if err := rows.Scan(
			&totalCost.Name,
			&totalCost.Cost,
		); err != nil {
			return nil, err
		}
		list = append(list, &totalCost)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}
