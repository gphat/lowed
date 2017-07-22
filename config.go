package lowed

type Config struct {
	Delay       string  `yaml:"delay"`
	ErrorChance float64 `yaml:"error_chance"`
	Metrics     struct {
		Counters []struct {
			Name string `yaml:"name"`
		} `yaml:"counters"`
		Gauges []struct {
			Name  string `yaml:"name"`
			Range struct {
				Max int `yaml:"max"`
				Min int `yaml:"min"`
			} `yaml:"range"`
		} `yaml:"gauges"`
		Histograms []struct {
			Name  string `yaml:"name"`
			Range struct {
				Max int `yaml:"max"`
				Min int `yaml:"min"`
			} `yaml:"range"`
		} `yaml:"histograms"`
		Sets []struct {
			Name         string `yaml:"name"`
			UniqueValues int    `yaml:"unique_values"`
		} `yaml:"sets"`
		Timers []struct {
			Name  string `yaml:"name"`
			Range struct {
				Max int `yaml:"max"`
				Min int `yaml:"min"`
			} `yaml:"range"`
		} `yaml:"timers"`
	} `yaml:"metrics"`
	Protocol     string   `yaml:"protocol"`
	Services     []string `yaml:"services"`
	StatsAddress string   `yaml:"stats_address"`
}
