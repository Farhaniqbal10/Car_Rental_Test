package services

import (
	"car_rental_test/modules/car/v1/models"
	"car_rental_test/pkg/errors"
	"context"
	"fmt"
	"time"
)

// Get Bookings by Parameters
func (s *Service) GetBookingsByParams(ctx context.Context, params models.BookingQueryParams) ([]models.BookingQueryParams, error) {
	fmt.Println("service1")
	bookings, err := s.carData.GetBookingsByParams(ctx, nil, params)
	if err != nil {
		return nil, errors.Wrap(err, "[SERVICE][GetBookingsByParams]")
	}
	fmt.Println("service2")
	return bookings, nil
}

// Create Booking with total cost calculation
func (s *Service) CreateBooking(ctx context.Context, booking models.Booking) error {
	// Cek validasi tanggal rental
	if booking.EndPeriod.Before(booking.StartPeriod) {
		return errors.Wrap(fmt.Errorf("tanggal kembali tidak boleh lebih awal dari tanggal rental"), "[SERVICE][CreateBooking]")
	}

	tx, err := s.carData.BeginTx(ctx)
	if err != nil {
		return errors.Wrap(err, "[SERVICE][CreateBooking]")
	}
	defer s.carData.RollbackTx(ctx, tx)

	// Ambil harga rental mobil
	rentPrice, err := s.carData.GetCarRentPrice(ctx, int64(booking.CarsID))
	if err != nil {
		return errors.Wrap(err, "[SERVICE][CreateBooking][GetCarRentPrice]")
	}

	// Hitung total biaya berdasarkan jumlah hari
	days := int64(booking.EndPeriod.Sub(booking.StartPeriod).Hours()/24) + 1
	booking.TotalCost = days * rentPrice

	today := time.Now().Truncate(24 * time.Hour) // Ambil tanggal hari ini tanpa waktu
	booking.Finished = booking.EndPeriod.Before(today)

	fmt.Println(today)

	// Kurangi stok mobil
	rowsAffected, err := s.carData.DecreaseCarStock(ctx, tx, int64(booking.CarsID))
	if err != nil {
		return errors.Wrap(err, "[SERVICE][CreateBooking][DecreaseCarStock]")
	}
	if rowsAffected == 0 {
		s.carData.RollbackTx(ctx, tx) // Batalkan transaksi jika stok habis
		return errors.Wrap(fmt.Errorf("stok mobil habis, transaksi dibatalkan"), "[SERVICE][CreateBooking]")
	}

	// Simpan booking ke database
	_, err = s.carData.CreateBooking(ctx, tx, booking)
	if err != nil {
		return errors.Wrap(err, "[SERVICE][CreateBooking][CreateBooking]")
	}

	// Commit transaksi
	err = s.carData.CommitTx(ctx, tx)
	if err != nil {
		return errors.Wrap(err, "[SERVICE][CreateBooking][CommitTx]")
	}

	return nil
}

// func (s *Service) CreateBooking(ctx context.Context, booking models.Booking) error {
// 	tx, err := s.carData.BeginTx(ctx)
// 	if err != nil {
// 		return errors.Wrap(err, "[SERVICE][CreateBooking]")
// 	}
// 	defer s.carData.RollbackTx(ctx, tx)

// 	// Hitung harga rental mobil
// 	rentPrice, err := s.carData.GetCarRentPrice(ctx, int64(booking.CarsID))
// 	if err != nil {
// 		return errors.Wrap(err, "[SERVICE][CreateBooking][GetCarRentPrice]")
// 	}

// 	// Hitung total biaya berdasarkan jumlah hari
// 	days := int64(booking.EndTime.Sub(booking.StartTime).Hours()/24) + 1
// 	booking.TotalCost = days * rentPrice

// 	// Simpan booking ke database
// 	_, err = s.carData.CreateBooking(ctx, tx, booking)
// 	if err != nil {
// 		return errors.Wrap(err, "[SERVICE][CreateBooking][CreateBooking]")
// 	}

// 	// Misalnya ada tambahan history atau proses lain, tambahkan di sini

// 	// Commit transaksi
// 	err = s.carData.CommitTx(ctx, tx)
// 	if err != nil {
// 		return errors.Wrap(err, "[SERVICE][CreateBooking][CommitTx]")
// 	}

// 	return nil
// }

// Update Booking with total cost recalculation

func (s *Service) UpdateBooking(ctx context.Context, id int64, booking models.Booking) error {
	tx, err := s.carData.BeginTx(ctx)
	if err != nil {
		return errors.Wrap(err, "[SERVICE][UpdateBooking][BeginTx]")
	}
	defer s.carData.RollbackTx(ctx, tx) // Jika ada error, rollback transaksi

	// 🔹 1. Ambil cars_id lama sebelum update
	oldBooking, errGetOldBooking := s.carData.GetBookingByID(ctx, id) // Ambil booking lama
	if errGetOldBooking != nil {
		return errors.Wrap(errGetOldBooking, "[SERVICE][UpdateBooking][GetBookingByID]")
	}

	// 🔹 2. Kembalikan stok mobil lama
	if oldBooking.CarsID != booking.CarsID { // Hanya jika ganti mobil
		_, errIncreaseStock := s.carData.IncreaseCarStock(ctx, tx, int64(oldBooking.CarsID))
		if errIncreaseStock != nil {
			return errors.Wrap(errIncreaseStock, "[SERVICE][UpdateBooking][IncreaseCarStock]")
		}

		fmt.Println(oldBooking.CarsID)

		// 🔹 3. Kurangi stok mobil baru
		rowsAffected, errDecreaseStock := s.carData.DecreaseCarStock(ctx, tx, int64(booking.CarsID))
		if errDecreaseStock != nil {
			return errors.Wrap(errDecreaseStock, "[SERVICE][UpdateBooking][DecreaseCarStock]")
		}
		if rowsAffected == 0 {
			return errors.Wrap(fmt.Errorf("stok mobil habis, transaksi dibatalkan"), "[SERVICE][UpdateBooking]")
		}
	}

	// 🔹 4. Hitung ulang total biaya berdasarkan harga rental mobil baru
	rentPrice, errGetCarRent := s.carData.GetCarRentPrice(ctx, int64(booking.CarsID))
	if errGetCarRent != nil {
		return errors.Wrap(errGetCarRent, "[SERVICE][UpdateBooking][GetCarRentPrice]")
	}

	days := int(booking.EndPeriod.Sub(booking.StartPeriod).Hours()/24) + 1
	booking.TotalCost = int64(days) * rentPrice

	// 🔹 5. Update booking di database
	_, errUpdateBooking := s.carData.UpdateBooking(ctx, id, tx, booking)
	if errUpdateBooking != nil {
		return errors.Wrap(errUpdateBooking, "[SERVICE][UpdateBooking][UpdateBooking]")
	}

	// 🔹 6. Commit transaksi jika semua berhasil
	err = s.carData.CommitTx(ctx, tx)
	if err != nil {
		return errors.Wrap(err, "[SERVICE][UpdateBooking][CommitTx]")
	}

	return nil
}

// Delete Booking
func (s *Service) DeleteBooking(ctx context.Context, id int64) (int, error) {
	rowsAffected, err := s.carData.DeleteBooking(ctx, id)
	if err != nil {
		return 0, errors.Wrap(err, "[SERVICE][DeleteBooking]")
	}
	return rowsAffected, nil
}
