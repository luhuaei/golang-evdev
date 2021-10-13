package main

import (
	"context"
	"fmt"

	evdev "github.com/gvalkov/golang-evdev"
)

type manager struct {
	touchpad *evdev.InputDevice
	keyboard *evdev.InputDevice

	isTouch bool
}

func NewWithSelect() (*manager, error) {
	fmt.Println("Select Your touchpad device...")
	touchpad, err := select_device()
	if err != nil {
		return nil, err
	}

	fmt.Println("Select Your keyboard device...")
	keyboard, err := select_device()
	if err != nil {
		return nil, err
	}

	return &manager{
		touchpad: touchpad,
		keyboard: keyboard,
	}, nil
}

func (m *manager) worker() error {
	ctx, cancel := context.WithCancel(context.Background())

	terrC := make(chan error, 1)
	go func() {
		terrC <- m.touchpadWorker(ctx)
		close(terrC)
		cancel()
	}()

	err := m.keyboardWorker(ctx)
	if err != nil {
		cancel()
		return err
	}
	return <-terrC
}

func (m *manager) touchpadWorker(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			es, err := m.touchpad.Read()
			if err != nil {
				return err
			}

			var touch bool
			var finger bool
			var found bool
			for _, e := range es {
				if e.Type != evdev.EV_KEY {
					continue
				}

				switch e.Code {
				case evdev.BTN_TOUCH:
					found = true
					touch = e.Value == int32(1)
				case evdev.BTN_TOOL_FINGER:
					found = true
					finger = e.Value == int32(1)
				}
			}

			if found {
				m.isTouch = touch && finger
			}
		}
	}
}

func (m *manager) keyboardWorker(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			es, err := m.keyboard.Read()
			if err != nil {
				return err
			}

			if !m.isTouch {
				continue
			}

			evfmt := "time %d.%-8d type %d (%s), code %-3d (%s), value %d\n"
			for _, ev := range es {
				if ev.Type != evdev.EV_KEY {
					continue
				}

				var key string
				var ok bool
				key, ok = evdev.KEY[int(ev.Code)]
				if !ok {
					key, ok = evdev.BTN[int(ev.Code)]
					if !ok {
						key = "?"
					}
				}
				fmt.Printf(evfmt, ev.Time.Sec, ev.Time.Usec, ev.Type,
					evdev.EV[int(ev.Type)], ev.Code, key, ev.Value)
			}
		}
	}
}
