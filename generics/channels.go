package generics

import "reflect"

func MergeDoneChannels[T any](channels ...<-chan T) <-chan T {

	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	done := make(chan T)

	go func() {
		defer close(done)
		switch len(channels) {
		case 2:
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default:
			select {
			case <-channels[0]:
			case <-channels[1]:
			case <-channels[2]:
			case <-MergeDoneChannels(append(channels[3:], done)...):
			}
		}
	}()
	return done
}

func readFromChannel[T any](chs ...<-chan T) <-chan T {
	done := make(chan T)

	go func() {
		defer close(done)

		cases := make([]reflect.SelectCase, len(chs))
		for i, ch := range chs {
			cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
		}

		for len(cases) > 0 {
			i, v, ok := reflect.Select(cases)
			if !ok {
				// The chosen channel has been closed, so remove it from the slice.
				cases = append(cases[:i], cases[i+1:]...)
				continue
			}
			done <- v.Interface().(T)
		}
	}()
	return done
}
