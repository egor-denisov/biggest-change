package usecase

type Option func(*StatsOfChangingUseCase)

func CacheSize(cacheSize int) Option {
	return func(uc *StatsOfChangingUseCase) {
		uc.cacheSize = cacheSize
	}
}

func MaxGoroutines(maxGoroutines int) Option {
	return func(s *StatsOfChangingUseCase) {
		s.maxGoroutines = maxGoroutines
	}
}

func AverageAddressCountInBlock(averageAddressCountInBlock int) Option {
	return func(s *StatsOfChangingUseCase) {
		s.averageAddressCountInBlock = averageAddressCountInBlock
	}
}
func DefaultCountOfBlocks(countOfBlocks uint) Option {
	return func(s *StatsOfChangingUseCase) {
		s.countOfBlocks = countOfBlocks
	}
}
