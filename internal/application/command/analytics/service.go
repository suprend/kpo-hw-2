package analytics

import (
	"context"

	appanalytics "kpo-hw-2/internal/application/analytics"
	appcommand "kpo-hw-2/internal/application/command"
	"kpo-hw-2/internal/domain"
)

// Service конструирует команды для аналитических сценариев.
type Service struct {
	analytics  appanalytics.Service
	decorators Decorators
}

// NewService связывает аналитический сервис и набор декораторов.
func NewService(analytics appanalytics.Service, decorators Decorators) *Service {
	return &Service{
		analytics:  analytics,
		decorators: decorators,
	}
}

// NetTotals формирует команду, считающую агрегированные показатели по операциям.
func (s *Service) NetTotals(operations []*domain.Operation) appcommand.Command[appanalytics.Totals] {
	base := appcommand.Func[appanalytics.Totals]{
		ExecFn: func(_ context.Context) (appanalytics.Totals, error) {
			if s.analytics == nil {
				return appanalytics.Totals{}, nil
			}
			return s.analytics.NetTotals(operations)
		},
		NameFn: func() string { return "analytics.net_totals" },
	}

	return appcommand.Wrap(base, s.decorators.NetTotals...)
}

// Decorators группирует опциональные декораторы команд аналитики.
type Decorators struct {
	NetTotals []appcommand.Decorator[appanalytics.Totals]
}
