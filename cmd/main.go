package main

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/wintermonth2298/library-ddd/internal/catalog/application"
	"github.com/wintermonth2298/library-ddd/internal/catalog/config"
	"github.com/wintermonth2298/library-ddd/internal/catalog/domain"
	"github.com/wintermonth2298/library-ddd/internal/catalog/infra/storage/projection"
	"github.com/wintermonth2298/library-ddd/internal/catalog/infra/storage/sql/psql"
	"github.com/wintermonth2298/library-ddd/internal/pkg/psqlclient"
)

func main() {
	cfg := config.MustLoad()

	db := psqlclient.MustNew(psqlclient.Config{
		Username: cfg.PSQL.User,
		Password: cfg.PSQL.Password,
		Host:     cfg.PSQL.Host,
		Port:     cfg.PSQL.Port,
		Database: cfg.PSQL.DB,
	})
	mustMigrateUp(db)

	storage := psql.NewStorage(db)

	service := domain.NewService()

	caseProjector := projection.NewCaseProjectior(db, service)
	usecases := application.NewUsecases(storage, service)
	usecases.RegisterEventHandler(domain.EventTypeSlideCreated, caseProjector.HandleSlideCreated)
	usecases.RegisterEventHandler(domain.EventTypeSlideFinished, caseProjector.HandleSlideUpdated)

	usecases.StartEventsProcessor(5 * time.Second)

	// err := usecases.CreateCase(context.Background())
	// if err != nil {
	// 	panic(err)
	// }

	slideID, _ := uuid.Parse("122c9557-784c-4fb2-a890-6f5a333505fb")
	err := usecases.AddSlide(context.Background(), slideID)
	if err != nil {
		panic(err)
	}

	// slideID, _ := uuid.Parse("2cbf1bb7-5bb3-44a8-aae0-ef3a68f4a96c")
	// err := usecases.FinishSlide(context.Background(), slideID)
	// if err != nil {
	// 	panic(err)
	// }

	time.Sleep(50 * time.Second)
}
