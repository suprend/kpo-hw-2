package bootstrap

import (
	"fmt"

	"kpo-hw-2/internal/infrastructure/di"
	"kpo-hw-2/internal/tui"
	mainmenu "kpo-hw-2/internal/tui/screens/main"
)

func registerUI(container di.Container) error {
	if err := di.Register(container, func(di.Container) (tui.Screen, error) {
		screen := mainmenu.New()
		if screen == nil {
			return nil, fmt.Errorf("bootstrap: root screen is nil")
		}
		return screen, nil
	}); err != nil {
		return fmt.Errorf("bootstrap: register root screen: %w", err)
	}

	return nil
}
