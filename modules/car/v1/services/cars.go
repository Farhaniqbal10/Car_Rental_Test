package services

import (
	"car_rental_test/modules/car/v1/models"
	"car_rental_test/pkg/errors"
	"context"
	"fmt"
)

func (s *Service) GetCar(ctx context.Context, params models.CarQueryParams) ([]models.Car, error) {
	car, err := s.carData.GetCar(ctx, params)
	if err != nil {
		return nil, errors.Wrap(err, "[SERVICE][GetCarsByParams]")
	}

	return car, nil

}

func (s *Service) CreateCar(ctx context.Context, car models.Car) (int64, error) {
	// Mulai transaksi
	tx, err := s.carData.BeginTx(ctx)
	if err != nil {
		return 0, errors.Wrap(err, "[SERVICE][CreateCar][BeginTx] gagal memulai transaksi")
	}
	defer s.carData.RollbackTx(ctx, tx) // Akan otomatis rollback jika ada error sebelum commit

	// Insert data mobil dan dapatkan ID baru
	id, err := s.carData.CreateCar(ctx, tx, car)
	if err != nil {
		return 0, errors.Wrap(err, "[SERVICE][CreateCar][CreateCar] gagal menyimpan data mobil")
	}

	// Commit transaksi jika tidak ada error
	err = s.carData.CommitTx(ctx, tx)
	if err != nil {
		return 0, errors.Wrap(err, "[SERVICE][CreateCar][CommitTx] gagal menyimpan transaksi")
	}

	return id, nil
}

func (s *Service) UpdateCar(ctx context.Context, id int64, car models.Car) error {
	// Ambil data lama berdasarkan CarsID
	oldCar, err := s.carData.GetCarByID(ctx, id)
	if err != nil {
		return fmt.Errorf("[SERVICE][UpdateCar][GetCarByID]: %w", err)
	}

	// Gunakan nilai lama jika field baru kosong atau bernilai default
	if car.Name != "" {
		oldCar.Name = car.Name
	}
	if car.RentPriceDaily != 0 {
		oldCar.RentPriceDaily = car.RentPriceDaily
	}
	if car.Stock != 0 {
		oldCar.Stock = car.Stock
	}

	// Update ke database
	_, err = s.carData.UpdateCar(ctx, id, oldCar)
	if err != nil {
		return fmt.Errorf("[SERVICE][UpdateCar]: %w", err)
	}

	return nil
}

func (s *Service) DeleteCar(ctx context.Context, id int64) (int, error) {
	rowsAffected, err := s.carData.DeleteCar(ctx, id)
	if err != nil {
		return 0, errors.Wrap(err, "[SERVICE][DeleteCar]")
	}
	return rowsAffected, nil
}
