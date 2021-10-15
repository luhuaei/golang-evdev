package main

import (
	"context"
	"fmt"

	"evdev"
)

type manager struct {
	touchpad *evdev.InputDevice
	keyboard *evdev.InputDevice

	output *evdev.OutputDevice

	isTouch bool
}

func NewWithSelect() (*manager, error) {
	output, err := evdev.BindOutput()
	if err != nil {
		fmt.Println("ERR", err)
		return nil, err
	}

	fmt.Println("Select Your touchpad device...")
	touchpad, err := selectDevice()
	if err != nil {
		return nil, err
	}

	fmt.Println("Select Your keyboard device...")
	keyboard, err := selectDevice()
	if err != nil {
		return nil, err
	}

	return &manager{
		touchpad: touchpad,
		keyboard: keyboard,
		output:   output,
	}, nil
}

func (m *manager) Close() error {
	err := m.keyboard.File.Close()
	if err != nil {
		return err
	}

	err = m.touchpad.File.Close()
	if err != nil {
		return err
	}

	return m.output.Close()
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
				err = m.output.Emits(es)
				if err != nil {
					return err
				}
				continue
			}

			evfmt := "time %d.%-8d type %d (%s), code %-3d (%s), value %d\n"
			for _, ev := range es {

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

				if ev.Type != evdev.EV_KEY || ev.Code == evdev.KEY_LEFTCTRL {
					err = m.output.Emit(ev)
					if err != nil {
						return err
					}
					continue
				}

				err = m.output.EmitKey(evdev.KEY_LEFTCTRL, true)
				if err != nil {
					return err
				}
				err = m.output.Emit(ev)
				if err != nil {
					return err
				}

				err = m.output.EmitKey(evdev.KEY_LEFTCTRL, false)
				if err != nil {
					return err
				}
			}
		}
	}
}
