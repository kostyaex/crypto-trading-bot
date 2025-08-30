package models

// для добавления полей в интерактивное поле см strategy_settings_field.templ

type StrategySettings struct {
	Symbol              string                 `mapstructure:"symbol"`   // используемая пара
	Interval            string                 `mapstructure:"interval"` // используемый интервал для выборки данных биржи
	Cluster             ClusterSettings        `mapstructure:"cluster"`
	SeriesBuilderConfig map[string]interface{} `mapstructure:"series_builder"`
	DispatcherConfig    DispatcherSettings     `mapstructure:"dispatcher"`
}
