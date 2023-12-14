package repository

import (
	"database/sql"
	"roomate/model/dto"
	"roomate/model/entity"
	query "roomate/utils/common"
	"time"
)

type BookingRepository interface {
	Get(id string) (entity.Booking, error)
	GetAll(limit, offset int) ([]entity.Booking, error)
	Create(booking entity.Booking) (entity.Booking, error)
	UpdateStatus(id string, isAgree bool, information string) (entity.Booking, error)
	Delete(id string) error
	GetOneDay(date string) (dto.SheetData, error)
	GetOneMonth(month, year string) ([]dto.SheetData, error)
	GetOneYear(year string) ([]dto.SheetData, error)
}

type bookingRepository struct {
	db *sql.DB
}

func (b *bookingRepository) Get(id string) (entity.Booking, error) {
	var booking entity.Booking
	err := b.db.QueryRow(query.GetBooking, id).
		Scan(
			&booking.Id,
			&booking.Night,
			&booking.CheckIn,
			&booking.CheckOut,
			&booking.UserId,
			&booking.CustomerId,
			&booking.IsAgree,
			&booking.Information,
			&booking.TotalPrice,
			&booking.CreatedAt,
			&booking.UpdatedAt,
			&booking.IsDeleted,
		)

	if err != nil {
		return booking, err
	}

	var bookingDetails []entity.BookingDetail
	rows, err := b.db.Query(query.GetAllBookingDetails, booking.Id)
	if err != nil {
		return booking, err
	}

	for rows.Next() {
		var bookingDetail entity.BookingDetail
		err := rows.Scan(
			&bookingDetail.Id,
			&bookingDetail.BookingId,
			&bookingDetail.RoomId,
			&bookingDetail.SubTotal,
			&bookingDetail.CreatedAt,
			&bookingDetail.UpdatedAt,
			&bookingDetail.IsDeleted,
		)

		if err != nil {
			return booking, err
		}

		bookingDetails = append(bookingDetails, bookingDetail)
	}

	booking.BookingDetails = bookingDetails
	return booking, nil
}

func (b *bookingRepository) GetAll(limit, offset int) ([]entity.Booking, error) {
	var bookings []entity.Booking
	rows, err := b.db.Query(query.GetAllBookings, limit, offset)
	if err != nil {
		return bookings, err
	}

	for rows.Next() {
		var booking entity.Booking
		err := rows.Scan(
			&booking.Id,
			&booking.Night,
			&booking.CheckIn,
			&booking.CheckOut,
			&booking.UserId,
			&booking.CustomerId,
			&booking.IsAgree,
			&booking.Information,
			&booking.TotalPrice,
			&booking.CreatedAt,
			&booking.UpdatedAt,
			&booking.IsDeleted,
		)

		if err != nil {
			return bookings, err
		}

		bookings = append(bookings, booking)
	}

	return bookings, nil
}

func (b *bookingRepository) Create(booking entity.Booking) (entity.Booking, error) {
	tx, err := b.db.Begin()
	if err != nil {
		return booking, err
	}

	err = tx.QueryRow(query.CreateBooking,
		booking.Night,
		booking.CheckIn,
		booking.CheckOut,
		booking.UserId,
		booking.CustomerId,
		booking.TotalPrice,
		time.Now(),
	).Scan(
		&booking.Id,
		&booking.Night,
		&booking.CheckIn,
		&booking.CheckOut,
		&booking.UserId,
		&booking.CustomerId,
		&booking.IsAgree,
		&booking.Information,
		&booking.TotalPrice,
		&booking.CreatedAt,
		&booking.UpdatedAt,
		&booking.IsDeleted,
	)

	if err != nil {
		tx.Rollback()
		return booking, err
	}

	var bookingDetails []entity.BookingDetail

	for _, i := range booking.BookingDetails {
		var bookingDetail entity.BookingDetail

		err = tx.QueryRow(query.CreateBookingDetail,
			booking.Id,
			i.RoomId,
			i.SubTotal,
			time.Now(),
		).Scan(
			&bookingDetail.Id,
			&bookingDetail.BookingId,
			&bookingDetail.RoomId,
			&bookingDetail.SubTotal,
			&bookingDetail.CreatedAt,
			&bookingDetail.UpdatedAt,
			&bookingDetail.IsDeleted,
		)

		if err != nil {
			tx.Rollback()
			return booking, err
		}

		var services []entity.BookingDetailService

		for _, j := range i.Services {
			var service entity.BookingDetailService
			err = tx.QueryRow(query.CreateBookingDetailService,
				bookingDetail.Id,
				j.ServiceId,
				j.ServiceName,
				time.Now(),
			).Scan(
				&service.Id,
				&service.BookingDetailId,
				&service.ServiceId,
				&service.ServiceName,
				&service.CreatedAt,
				&service.UpdatedAt,
				&service.IsDeleted,
			)
			if err != nil {
				tx.Rollback()
				return booking, err
			}

			services = append(services, service)
		}

		bookingDetail.Services = services
		bookingDetails = append(bookingDetails, bookingDetail)
	}

	booking.BookingDetails = bookingDetails

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return booking, err
	}

	return booking, nil
}

func (b *bookingRepository) UpdateStatus(id string, isAgree bool, information string) (entity.Booking, error) {
	var booking entity.Booking
	err := b.db.QueryRow(query.UpdateBookingStatus, id, isAgree, information).Scan(
		&booking.Id,
		&booking.Night,
		&booking.CheckIn,
		&booking.CheckOut,
		&booking.UserId,
		&booking.CustomerId,
		&booking.IsAgree,
		&booking.Information,
		&booking.TotalPrice,
		&booking.CreatedAt,
		&booking.UpdatedAt,
		&booking.IsDeleted,
	)

	if err != nil {
		return booking, err
	}

	return booking, nil
}

func (b *bookingRepository) Delete(id string) error {
	_, err := b.db.Exec(query.DeleteBooking, id)
	if err != nil {
		return err
	}

	return nil
}

func (b *bookingRepository) GetOneDay(date string) (dto.SheetData, error) {
	var booking dto.SheetData
	err := b.db.QueryRow(query.GetBookingOneDay, date).Scan(
		&booking.BookingId,
		&booking.CheckIn,
		&booking.CheckOut,
		&booking.UserName,     // disini masih user id
		&booking.CustomerName, // dinisi juga masih customer id
		&booking.IsAgree,
		&booking.Information,
		&booking.TotalPrice,
	)

	if err != nil {
		return booking, err
	}

	return booking, nil
}

func (b *bookingRepository) GetOneMonth(month, year string) ([]dto.SheetData, error) {
	var bookings []dto.SheetData
	rows, err := b.db.Query(query.GetBookingOneMonth, month, year)
	if err != nil {
		return bookings, err
	}

	for rows.Next() {
		var booking dto.SheetData
		err = rows.Scan(
			&booking.BookingId,
			&booking.CheckIn,
			&booking.CheckOut,
			&booking.UserName,     // disini masih user id
			&booking.CustomerName, // dinisi juga masih customer id
			&booking.IsAgree,
			&booking.Information,
			&booking.TotalPrice,
		)

		if err != nil {
			return bookings, err
		}

		bookings = append(bookings, booking)
	}

	return bookings, nil
}

func (b *bookingRepository) GetOneYear(year string) ([]dto.SheetData, error) {
	var bookings []dto.SheetData
	rows, err := b.db.Query(query.GetBookingOneYear, year)
	if err != nil {
		return bookings, err
	}

	for rows.Next() {
		var booking dto.SheetData
		err = rows.Scan(
			&booking.BookingId,
			&booking.CheckIn,
			&booking.CheckOut,
			&booking.UserName,     // disini masih user id
			&booking.CustomerName, // dinisi juga masih customer id
			&booking.IsAgree,
			&booking.Information,
			&booking.TotalPrice,
		)

		if err != nil {
			return bookings, err
		}

		bookings = append(bookings, booking)
	}

	return bookings, nil
}

func NewBookingRepository(db *sql.DB) BookingRepository {
	return &bookingRepository{
		db: db,
	}
}
