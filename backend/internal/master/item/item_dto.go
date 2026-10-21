package item

type KatalogItemResponse struct {
	IDMotor     string `json:"id_motor"`
	NamaMotor   string `json:"nama_motor"`
	GambarColor string `json:"gambar_color"`
}

type ColorSpec struct {
	ColorName   string `json:"color_name"`
	GambarColor string `json:"gambar_color"`
}

type MotorcycleSpecResponse struct {
	IDMotor   string      `json:"id_motor"`
	NamaMotor string      `json:"nama_motor"`
	Harga     float64     `json:"harga"`
	Detail    string      `json:"detail"`
	Brosur    string      `json:"brosur"`
	Colors    []ColorSpec `json:"colors"`
	Fiturs    []string    `json:"fiturs"`
}

type PointResponse struct {
	Point float64 `json:"point"`
}
