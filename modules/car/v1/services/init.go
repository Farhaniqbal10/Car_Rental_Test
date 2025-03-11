package services

import (
	"car_rental_test/modules/car/v1/models"
	"context"

	"github.com/jmoiron/sqlx"
)

type ICarData interface {
	//transaction
	BeginTx(ctx context.Context) (*sqlx.Tx, error)
	CommitTx(ctx context.Context, tx *sqlx.Tx) error
	RollbackTx(ctx context.Context, tx *sqlx.Tx) error
	//BOOKING
	GetBookingsByParams(ctx context.Context, tx *sqlx.Tx, query models.BookingQueryParams) ([]models.Booking, error)
	CreateBooking(ctx context.Context, tx *sqlx.Tx, input models.Booking) (int64, error)
	UpdateBooking(ctx context.Context, id int64, tx *sqlx.Tx, input models.Booking) (models.Booking, error)
	DeleteBooking(ctx context.Context, id int64) (int, error)
	GetCarRentPrice(ctx context.Context, carID int64) (int64, error) // Fungsi baru untuk mendapatkan harga sewa mobil
	DecreaseCarStock(ctx context.Context, tx *sqlx.Tx, carID int64) (int64, error)
	IncreaseCarStock(ctx context.Context, tx *sqlx.Tx, carID int64) (int64, error) // 🔹 Fungsi baru untuk menambah stok mobil
	GetBookingByID(ctx context.Context, id int64) (models.Booking, error)          // 🔹 Fungsi baru untuk mendapatkan booking lama
	//CARS
	GetCar(ctx context.Context, query models.CarQueryParams) ([]models.Car, error)
	CreateCar(ctx context.Context, tx *sqlx.Tx, input models.Car) (int64, error)
	UpdateCar(ctx context.Context, id int64, input models.Car) (models.Car, error)
	DeleteCar(ctx context.Context, id int64) (int, error)
	GetCarByID(ctx context.Context, id int64) (models.Car, error)
}

type Service struct {
	carData ICarData
}

func New(carData ICarData) *Service {
	return &Service{
		carData: carData,
	}
}
