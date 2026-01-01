package psql

import (
	"context"
	"fmt"

	"github.com/hfleury/bk_globalshot/internal/model"
	"github.com/hfleury/bk_globalshot/internal/repository"
	"github.com/hfleury/bk_globalshot/pkg/db"
)

type PostgresRoomRepository struct {
	db db.Db
}

func NewPostgresRoomRepository(db db.Db) repository.RoomRepository {
	return &PostgresRoomRepository{db: db}
}

func (r *PostgresRoomRepository) Create(ctx context.Context, room *model.Room) error {
	query := `INSERT INTO rooms (name, unit_id, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id`
	// Use GetDb() to access DbTx
	err := r.db.GetDb().QueryRowContext(ctx, query, room.Name, room.UnitID, room.CreatedAt, room.UpdatedAt).Scan(&room.ID)
	if err != nil {
		return fmt.Errorf("failed to create room: %w", err)
	}
	return nil
}

func (r *PostgresRoomRepository) FindAll(ctx context.Context, limit, offset int, unitID string) ([]*model.Room, int64, error) {
	query := `SELECT id, name, unit_id, created_at, updated_at FROM rooms WHERE 1=1`
	args := []interface{}{}
	argCounter := 1

	if unitID != "" {
		query += fmt.Sprintf(" AND unit_id = $%d", argCounter)
		args = append(args, unitID)
		argCounter++
	}

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCounter, argCounter+1)
	args = append(args, limit, offset)

	rows, err := r.db.GetDb().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list rooms: %w", err)
	}
	defer rows.Close()

	rooms := make([]*model.Room, 0)
	for rows.Next() {
		var room model.Room
		if err := rows.Scan(&room.ID, &room.Name, &room.UnitID, &room.CreatedAt, &room.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("failed to scan room: %w", err)
		}
		rooms = append(rooms, &room)
	}

	var total int64
	countQuery := `SELECT COUNT(*) FROM rooms`
	// Simple count for now, if filtering becomes heavy we should filter count too
	if unitID != "" {
		countQuery += " WHERE unit_id = $1"
		if err := r.db.GetDb().QueryRowContext(ctx, countQuery, unitID).Scan(&total); err != nil {
			return nil, 0, fmt.Errorf("failed to count rooms: %w", err)
		}
	} else {
		if err := r.db.GetDb().QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
			return nil, 0, fmt.Errorf("failed to count rooms: %w", err)
		}
	}

	return rooms, total, nil
}

func (r *PostgresRoomRepository) FindByID(ctx context.Context, id string) (*model.Room, error) {
	query := `SELECT id, name, unit_id, created_at, updated_at FROM rooms WHERE id = $1`
	var room model.Room
	err := r.db.GetDb().QueryRowContext(ctx, query, id).Scan(&room.ID, &room.Name, &room.UnitID, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to find room: %w", err)
	}
	return &room, nil
}

func (r *PostgresRoomRepository) Update(ctx context.Context, room *model.Room) error {
	query := `UPDATE rooms SET name = $1, unit_id = $2, updated_at = $3 WHERE id = $4`
	_, err := r.db.GetDb().ExecContext(ctx, query, room.Name, room.UnitID, room.UpdatedAt, room.ID)
	if err != nil {
		return fmt.Errorf("failed to update room: %w", err)
	}
	return nil
}

func (r *PostgresRoomRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM rooms WHERE id = $1`
	_, err := r.db.GetDb().ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete room: %w", err)
	}
	return nil
}
