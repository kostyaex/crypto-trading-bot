package types

type ClusterConfig struct {
	NumClusters int    `mapstructure:"num_clusters"`
	Block       int    `mapstructure:"block"`    // сколько записей торговых данных группировать в кластер
	Interval    string `mapstructure:"interval"` // результирующий интервал сгруппированных записей
}

func (c ClusterConfig) GetComponentType() string {
	return "clustering"
}
