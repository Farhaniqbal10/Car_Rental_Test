package services

import (
	"car_rental_test/modules/car/v1/models"
	"car_rental_test/pkg/errors"
	"context"
	"fmt"
	"time"
)

// Get Bookings by Parameters
func (s *Service) GetBookingsByParams(ctx context.Context, params models.BookingQueryParams) ([]models.Booking, error) {
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

	// bookings := models.BookingInput{ // Tidak perlu di-convert karena sama
	// 	CustomerID:  booking.CustomerID, // Convert dari null.Int64 ke int64
	// 	CarsID:      booking.CarsID,
	// 	StartPeriod: booking.StartPeriod, // Convert dari null.Time ke time.Time
	// 	EndPeriod:   booking.EndPeriod,
	// 	TotalCost:   booking.TotalCost,
	// 	Finished:    booking.Finished,
	// }

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

func (s *Service) UpdateBooking(ctx context.Context, id int64, booking models.Booking) error {
	// 1️⃣ Validasi tanggal
	// if booking.EndPeriod.Before(booking.StartPeriod) {
	// 	return fmt.Errorf("[SERVICE][UpdateBooking] tanggal kembali tidak boleh lebih awal dari tanggal rental")
	// }

	tx, err := s.carData.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("[SERVICE][UpdateBooking][BeginTx]: %w", err)
	}
	defer s.carData.RollbackTx(ctx, tx)

	// 2️⃣ Ambil data booking lama
	oldBooking, errGetOldBooking := s.carData.GetBookingByID(ctx, id)
	if errGetOldBooking != nil {
		return fmt.Errorf("[SERVICE][UpdateBooking][GetBookingByID]: %w", errGetOldBooking)
	}

	// Jika CustomerID kosong, gunakan nilai lama
	if booking.CustomerID == 0 {
		booking.CustomerID = oldBooking.CustomerID
	}

	// Jika CarsID kosong, gunakan nilai lama
	if booking.CarsID == 0 {
		booking.CarsID = oldBooking.CarsID
	}

	// Jika StartPeriod kosong, gunakan nilai lama
	if booking.StartPeriod.IsZero() {
		booking.StartPeriod = oldBooking.StartPeriod
	}

	// Validasi StartPeriod tidak melebihi EndPeriod sebelumnya
	if booking.StartPeriod.After(oldBooking.EndPeriod) {
		return fmt.Errorf("[SERVICE][UpdateBooking] update invalid: tanggal mulai melebihi tanggal akhir sebelumnya")
	}

	// Jika EndPeriod kosong, gunakan nilai lama
	if booking.EndPeriod.IsZero() {
		booking.EndPeriod = oldBooking.EndPeriod
	}

	// Jika Finished tidak dikirim, gunakan nilai lama
	if !booking.Finished {
		booking.Finished = oldBooking.Finished
	}

	// 3️⃣ Update stok mobil jika CarsID berubah
	if booking.CarsID != oldBooking.CarsID {
		_, errIncreaseStock := s.carData.IncreaseCarStock(ctx, tx, oldBooking.CarsID)
		if errIncreaseStock != nil {
			return fmt.Errorf("[SERVICE][UpdateBooking][IncreaseCarStock]: %w", errIncreaseStock)
		}

		rowsAffected, errDecreaseStock := s.carData.DecreaseCarStock(ctx, tx, booking.CarsID)
		if errDecreaseStock != nil {
			return fmt.Errorf("[SERVICE][UpdateBooking][DecreaseCarStock]: %w", errDecreaseStock)
		}
		if rowsAffected == 0 {
			return fmt.Errorf("[SERVICE][UpdateBooking] stok mobil habis, transaksi dibatalkan")
		}
	}

	// 4️⃣ Hitung ulang total biaya
	rentPrice, errGetCarRent := s.carData.GetCarRentPrice(ctx, booking.CarsID)
	if errGetCarRent != nil {
		return fmt.Errorf("[SERVICE][UpdateBooking][GetCarRentPrice]: %w", errGetCarRent)
	}

	days := int(booking.EndPeriod.Sub(booking.StartPeriod).Hours()/24) + 1
	booking.TotalCost = int64(days) * rentPrice

	// bookings := models.BookingUpdate{ // Tidak perlu di-convert karena sama
	// 	CustomerID:  null.IntFrom(booking.CustomerID), // Convert dari null.Int64 ke int64
	// 	CarsID:      null.IntFrom(booking.CarsID),
	// 	StartPeriod: null.TimeFrom(booking.StartPeriod), // Convert dari null.Time ke time.Time
	// 	EndPeriod:   null.TimeFrom(booking.EndPeriod),
	// 	TotalCost:   null.IntFrom(booking.TotalCost),
	// 	Finished:    null.BoolFrom(booking.Finished),
	// }

	// 5️⃣ Update booking di database
	_, errUpdateBooking := s.carData.UpdateBooking(ctx, id, tx, booking)
	if errUpdateBooking != nil {
		return fmt.Errorf("[SERVICE][UpdateBooking][UpdateBooking]: %w", errUpdateBooking)
	}

	// 6️⃣ Commit transaksi
	err = s.carData.CommitTx(ctx, tx)
	if err != nil {
		return fmt.Errorf("[SERVICE][UpdateBooking][CommitTx]: %w", err)
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
