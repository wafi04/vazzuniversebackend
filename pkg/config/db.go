package config

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
)

type Database struct {
	Main    *sqlx.DB
	Replica *sqlx.DB
}

func NewDatabase() (*Database, error) {

	mainDbURL := LoadEnv("DATABASE_MAIN_URL")
	replicaDbURL := LoadEnv("DATABASE_REPLICA_URL")
	// Connect to main database
	mainDB, err := sqlx.Open("pgx", mainDbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open main database connection: %v", err)
	}

	mainDB.SetMaxOpenConns(25)
	mainDB.SetMaxIdleConns(5)
	mainDB.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err := mainDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping main database: %v", err)
	}

	// Connect to replica database
	replicaDB, err := sqlx.Open("pgx", replicaDbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open replica database connection: %v", err)
	}

	replicaDB.SetMaxOpenConns(25)
	replicaDB.SetMaxIdleConns(5)
	replicaDB.SetConnMaxLifetime(5 * time.Minute)

	if err := replicaDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping replica database: %v", err)
	}

	return &Database{
		Main:    mainDB,
		Replica: replicaDB,
	}, nil
}

func (d *Database) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	err := d.Main.DB.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf("db down: %v", err)
		return stats
	}

	stats["status"] = "up"
	stats["message"] = "It's healthy"

	dbStats := d.Main.DB.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	if dbStats.OpenConnections > 40 {
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

// Close closes both database connections
func (d *Database) Close() {
	if d.Main != nil {
		d.Main.Close()
	}
	if d.Replica != nil {
		d.Replica.Close()
	}
}
