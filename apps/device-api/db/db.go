package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"smarthome/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB represents the database connection
type DB struct {
	Pool *pgxpool.Pool
}

// New creates a new DB instance
func New(connString string) (*DB, error) {
	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	// Test the connection
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return &DB{Pool: pool}, nil
}

// Close closes the database connection
func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}

// GetDevicesByHouseId retrieves all devices by house from the database
func (db *DB) GetDevicesByHouseId(ctx context.Context, houseId int) ([]models.Device, error) {
	query := `
		SELECT id, device_model_id, house_id, serial_number, name, status, last_updated, created_at
		FROM devices
		WHERE id = $1
		ORDER BY id
	`
	rows, err := db.Pool.Query(ctx, query, houseId)
	if err != nil {
		return nil, fmt.Errorf("error querying devices: %w", err)
	}
	defer rows.Close()

	var devices []models.Device
	for rows.Next() {
		var s models.Device
		err := rows.Scan(
			&s.ID,
			&s.DeviceModelID,
			&s.HouseID,
			&s.SerialNumber,
			&s.Name,
			&s.Status,
			&s.LastUpdated,
			&s.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning device row: %w", err)
		}
		devices = append(devices, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating device rows: %w", err)
	}

	return devices, nil
}

// GetDeviceByID retrieves a device by its ID
func (db *DB) GetDeviceByID(ctx context.Context, id int) (models.Device, error) {
	query := `
		SELECT id, device_model_id, house_id, serial_number, name, status, last_updated, created_at
		FROM devices
		WHERE id = $1
	`

	var s models.Device
	err := db.Pool.QueryRow(ctx, query, id).Scan(
		&s.ID,
		&s.DeviceModelID,
		&s.HouseID,
		&s.SerialNumber,
		&s.Name,
		&s.Status,
		&s.LastUpdated,
		&s.CreatedAt,
	)
	if err != nil {
		return models.Device{}, fmt.Errorf("error getting device by ID: %w", err)
	}

	return s, nil
}

// CreateDevice creates a new device in the database
func (db *DB) CreateDevice(ctx context.Context, s models.DeviceCreate) (models.Device, error) {
	query := `
		INSERT INTO devices (device_model_id, house_id, serial_number, name, status, last_updated, created_at)
		VALUES ($1, $2, $3, $4, 'active', $5, $5)
		RETURNING id, device_model_id, house_id, serial_number, name, status, last_updated, created_at
	`

	now := time.Now()
	var device models.Device
	err := db.Pool.QueryRow(ctx, query,
		s.DeviceModelID,
		s.HouseID,
		s.SerialNumber,
		s.Name,
		now,
	).Scan(
		&device.ID,
		&device.DeviceModelID,
		&device.HouseID,
		&device.SerialNumber,
		&device.Name,
		&device.Status,
		&device.LastUpdated,
		&device.CreatedAt,
	)
	if err != nil {
		return models.Device{}, fmt.Errorf("error creating device: %w", err)
	}

	return device, nil
}

// UpdateDevice updates an existing device
func (db *DB) UpdateDevice(ctx context.Context, id int, s models.DeviceUpdate) (models.Device, error) {
	// First check if the device exists
	_, err := db.GetDeviceByID(ctx, id)
	if err != nil {
		return models.Device{}, err
	}

	// Build the update query dynamically based on which fields are provided
	query := "UPDATE devices SET last_updated = $1"
	args := []interface{}{time.Now()}
	argCount := 2

	if s.Name != "" {
		query += fmt.Sprintf(", name = $%d", argCount)
		args = append(args, s.Name)
		argCount++
	}

	if s.SerialNumber != "" {
		query += fmt.Sprintf(", serial_number = $%d", argCount)
		args = append(args, s.SerialNumber)
		argCount++
	}

	query += fmt.Sprintf(", house_id = $%d", argCount)
	args = append(args, s.HouseID)
	argCount++

	query += fmt.Sprintf(", device_model_id = $%d", argCount)
	args = append(args, s.DeviceModelID)
	argCount++

	// Add the WHERE clause and RETURNING clause
	query += ` WHERE id = $` + fmt.Sprintf("%d", argCount) + `
		RETURNING id, device_model_id, house_id, serial_number, name, status, last_updated, created_at`
	args = append(args, id)

	var device models.Device
	err = db.Pool.QueryRow(ctx, query, args...).Scan(
		&device.ID,
		&device.DeviceModelID,
		&device.HouseID,
		&device.SerialNumber,
		&device.Name,
		&device.Status,
		&device.LastUpdated,
		&device.CreatedAt,
	)
	if err != nil {
		return models.Device{}, fmt.Errorf("error updating device: %w", err)
	}

	return device, nil
}

// DeleteDevice deletes a device by its ID
func (db *DB) DeleteDevice(ctx context.Context, id int) error {
	query := "DELETE FROM devices WHERE id = $1"
	result, err := db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting device: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("device not found")
	}

	return nil
}

// GetDevicesAttributes retrieves all device attribute by deviceId
func (db *DB) GetDevicesAttributes(ctx context.Context, deviceId int) ([]models.DeviceAttribute, error) {
	query := `
		SELECT id, device_id, key, value
		FROM devices_attributes
		WHERE device_id = $1
		ORDER BY id
	`
	rows, err := db.Pool.Query(ctx, query, deviceId)
	if err != nil {
		return nil, fmt.Errorf("error querying device attribute: %w", err)
	}
	defer rows.Close()

	var deviceAttrs []models.DeviceAttribute
	for rows.Next() {
		var s models.DeviceAttribute
		err := rows.Scan(
			&s.ID,
			&s.DeviceID,
			&s.Key,
			&s.Value,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning device attribute row: %w", err)
		}
		deviceAttrs = append(deviceAttrs, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating device attribute rows: %w", err)
	}

	return deviceAttrs, nil
}

// CreateDeviceAttribute creates a new device attribute in the database
func (db *DB) CreateDeviceAttribute(ctx context.Context, s models.DeviceAttribute) (models.DeviceAttribute, error) {
	query := `
		INSERT INTO devices_attributes (device_id, key, value)
		VALUES ($1, $2, $3)
		RETURNING id, device_id, key, value
	`

	var attr models.DeviceAttribute
	err := db.Pool.QueryRow(ctx, query,
		s.DeviceID,
		s.Key,
		s.Value,
	).Scan(
		&attr.ID,
		&attr.DeviceID,
		&attr.Key,
		&attr.Value,
	)
	if err != nil {
		return models.DeviceAttribute{}, fmt.Errorf("error creating device attribute: %w", err)
	}

	return attr, nil
}

// DeleteDeviceAttributes deletes a device attributes by deviceId
func (db *DB) DeleteDeviceAttributes(ctx context.Context, deviceId int) error {
	query := "DELETE FROM devices_attributes WHERE device_id = $1"
	result, err := db.Pool.Exec(ctx, query, deviceId)
	if err != nil {
		return fmt.Errorf("error deleting device attributes: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("device attributes not found")
	}

	return nil
}
