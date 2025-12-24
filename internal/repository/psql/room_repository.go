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

func (r *PostgresRoomRepository) FindAll(ctx context.Context, limit, offset int) ([]*model.Room, int64, error) {
	// TODO: Add filtering by unit_id if needed
	query := `SELECT id, name, unit_id, created_at, updated_at FROM rooms LIMIT $1 OFFSET $2`
	rows, err := r.db.GetDb().QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list rooms: %w", err)
	}
	defer rows.Close()

	var rooms []*model.Room
	for rows.Next() {
		var room model.Room
		if err := rows.Scan(&room.ID, &room.Name, &room.UnitID, &room.CreatedAt, &room.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("failed to scan room: %w", err)
		}
		rooms = append(rooms, &room)
	}

	var total int64
	countQuery := `SELECT COUNT(*) FROM rooms`
	if err := r.db.GetDb().QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count rooms: %w", err)
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
