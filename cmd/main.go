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

	caseProjector := projection.NewCaseProjectior(db, service, slidesRepoSQL, casesRepoSQL)

	usecases := application.NewUsecases(storage)
	usecases.RegisterEventHandler(domain.EventTypeSlideCreated, caseProjector.HandleSlideCreated)
	usecases.RegisterEventHandler(domain.EventTypeSlideUpdated, caseProjector.HandleSlideUpdated)

	usecases.StartEventsProcessor(5 * time.Second)

	// err := usecases.CreateCase(context.Background())
	// if err != nil {
	// 	panic(err)
	// }

	// slideID, _ := uuid.Parse("6f52d738-53b5-461f-88b3-b26b37cdcebe")
	// err := usecases.AddSlide(context.Background(), slideID)
	// if err != nil {
	// 	panic(err)
	// }

	slideID, _ := uuid.Parse("6289620d-68bd-41a6-adac-41ede95bb16d")
	err := usecases.FinishSlide(context.Background(), slideID)
	if err != nil {
		panic(err)
	}

	time.Sleep(50 * time.Second)
}
