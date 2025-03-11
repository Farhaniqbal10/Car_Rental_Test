package data

import (
	"car_rental_test/modules/car/v1/models"
	"car_rental_test/pkg/errors"
	"context"
	"log"

	"github.com/jmoiron/sqlx"
)

const (
	getCarsByParams  = "GetCarsByParams"
	qGetCarsByParams = `SELECT * FROM cars 
WHERE 
    (cars_id = $1 OR $1 = 0)
    AND (name ILIKE $2 OR $2 = '%')
    AND (rent_price_daily = $3 OR $3 = 0)
    AND (stock >= $4 OR $4 = 0)
    AND (name ILIKE $5 OR $5 = '%') 
ORDER BY cars_id ASC;

`
	getCars      = "GetCars"
	qGetCars     = `SELECT cars_id, name, rent_price_daily, stock FROM cars ORDER BY cars_id ASC`
	getCarsByID  = "GetCarsByID"
	qGetCarsByID = `SELECT cars_id, name, rent_price_daily, stock FROM cars WHERE cars_id = $1`
	insertCars   = "InsertCars"
	qInsertCars  = `INSERT INTO cars (name, rent_price_daily, stock) VALUES ($1, $2, $3) RETURNING cars_id`
	updateCars   = "UpdateCars"
	qUpdateCars  = `UPDATE cars SET name = $1, rent_price_daily = $2, stock = $3 WHERE cars_id = $4`
	deleteCars   = "DeleteCars"
	qDeleteCars  = `DELETE FROM cars WHERE cars_id = $1`
)

var (
	CarsStmts = []statement{
		{getCarsByParams, qGetCarsByParams},
		{getCars, qGetCars},
		{getCarsByID, qGetCarsByID},
		{insertCars, qInsertCars},
		{updateCars, qUpdateCars},
		{deleteCars, qDeleteCars},
	}
)

func (d Data) GetCar(ctx context.Context, query models.CarQueryParams) ([]models.Car, error) {
	cars := []models.Car{}

	if query.Name != "" {
		query.Name = "%" + query.Name + "%"
	} else {
		query.Name = "%"
	}

	if query.Keyword != "" {
		query.Keyword = "%" + query.Keyword + "%"
	} else {
		query.Keyword = "%"
	}

	log.Printf("Executing query with params: CarsID=%d, Name=%s, RentPriceDaily=%d, Stock=%d, Keyword=%s",
		query.CarsID, query.Name, query.RentPriceDaily, query.Stock, query.Keyword)

	err := d.stmt[getCarsByParams].SelectContext(ctx, &cars,
		query.CarsID, query.Name, query.RentPriceDaily, query.Stock, query.Keyword)
	if err != nil {
		return cars, errors.Wrap(err, "[DATA][GetCars]")
	}
	return cars, nil

}

func (d *Data) GetCarByID(ctx context.Context, id int64) (models.Car, error) {
	var car models.Car
	err := d.stmt[getCarsByID].GetContext(ctx, &car, id)
	if err != nil {
		return car, errors.Wrap(err, "[DATA][GetBookingByID]")
	}
	return car, nil
}

func (d *Data) GetAllCar(ctx context.Context) ([]models.Car, error) {
	var car []models.Car
	err := d.stmt[getCars].SelectContext(ctx, &car)
	if err != nil {
		return car, errors.Wrap(err, "[DATA][GetAllBookings]")
	}
	return car, nil
}

func (d *Data) CreateCar(ctx context.Context, tx *sqlx.Tx, input models.Car) (int64, error) {

	var id int64
	stmt := d.stmt[insertCars]
	if tx != nil {
		stmt = tx.StmtxContext(ctx, stmt)
	}

	err := stmt.GetContext(ctx, &id,
		input.Name,
		input.RentPriceDaily,
		input.Stock,
	)

	if err != nil {
		return id, errors.Wrap(err, "[DATA][insertArea]")
	}

	return id, nil
}

func (d *Data) UpdateCar(ctx context.Context, id int64, input models.Car) (models.Car, error) {
	cars := models.Car{}

	stmt := d.stmt[updateCars]

	_, err := stmt.ExecContext(ctx,
		input.Name,
		input.RentPriceDaily,
		input.Stock,
		id)
	if err != nil {
		return cars, errors.Wrap(err, "[DATA][updateCars]")
	}

	return cars, nil
}

func (d *Data) DeleteCar(ctx context.Context, id int64) (int, error) {
	result, err := d.db.ExecContext(ctx, "DELETE FROM cars WHERE cars_id = $1", id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowsAffected), nil
}
