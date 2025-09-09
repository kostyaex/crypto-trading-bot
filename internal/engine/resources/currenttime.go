// Время выполнения. Текущее время при выполнении в реальном режиме или время данные из базы при бэктестинге
package resources

type CurrentTime struct {
	Timestamp int64
}

func NewCurrentTime() *CurrentTime {
	return &CurrentTime{}
}
