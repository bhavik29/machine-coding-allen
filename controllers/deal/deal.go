package deal

import (
	"allen-machine-coding/db/models/deal"
	"context"
	"errors"
	"time"
)

type Deal struct {
	ID               string
	Name             string
	MaxNumberOfItems int
	Duration         time.Duration // in seconds
	IsActive         bool
}

func (d *Deal) CreateDeal(ctx context.Context) (string, error) {
	currentTime := time.Now()

	dealEntity := deal.DealSchema{
		Name:             d.Name,
		MaxNumberOfItems: d.MaxNumberOfItems,
		DealEndTime:      currentTime.Add(d.Duration * time.Second),
		IsActive:         d.IsActive,
	}

	dealID, err := deal.InsertOne(ctx, &dealEntity)
	if err != nil {
		return "", err
	}

	return dealID, nil
}

func (d *Deal) EndDeal(ctx context.Context, dealID string) error {
	dealEntity := deal.FindOne(ctx, dealID)
	if dealEntity == nil {
		return errors.New("no deal to end")
	}

	err := deal.UpdateOne(ctx, dealID, &deal.DealSchema{
		Name:             dealEntity.Name,
		MaxNumberOfItems: dealEntity.MaxNumberOfItems,
		DealEndTime:      dealEntity.DealEndTime,
		IsActive:         false,
	})

	return err
}
