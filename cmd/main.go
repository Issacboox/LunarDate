package main

import (
	"bam/internal/adapter/config"
	h "bam/internal/adapter/handler"
	r "bam/internal/adapter/repository/mysql/repository"
	s "bam/internal/core/service"

	"bam/internal/adapter/repository/mysql"
	"bam/internal/adapter/route"
	ut "bam/internal/core/utils"
	"context"
)

var prod = ut.GetEnv("APP_ENV", "development") == "prod"

func main() {
	config, err := config.New()
	if err != nil {
		panic(err)
	}

	db, err := mysql.NewDatabase(context.Background(), config.DB, prod)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
	}()

	if err := db.Migrate(); err != nil {
		panic(err)
	}

	if err := db.SetMaxOpenConns(10); err != nil {
		panic(err)
	}

	orbianRepository := r.NewOrbianRepository(db.DB)
	orbianService := s.NewOrbianService(orbianRepository)
	orbianHandler := h.NewOrbianHandler(orbianService)

	r, err := route.NewRouter(*orbianHandler) // Pass the handler directly, not the address
	if err != nil {
		panic(err)
	}

	r.Serve(":3000")
}
