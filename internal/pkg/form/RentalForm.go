package form

type RentalForm struct {
	CustomerUUID    string  `json:"customerUUID" validate:"required"`
	MotorcycleUUID  string  `json:"motorcycleUUID" validate:"required"`
	ReturnDatePlan  string  `json:"returnDatePlan" validate:"required"`
	PaymentMethodId int     `json:"paymentMethodId"`
	PaymentAmount   float64 `json:"paymentAmount"`
}
