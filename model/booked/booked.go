package booked

type Booked struct {
	BookedId int `json:"booked_id"`
	MotorId int `json:"motor_id"`
	Nama string `json:"nama"`
	Alamat string `json:"alamat"`
	Harga string `json:"harga"`
}