package supersense

import "github.com/pkg/errors"

// Mux is a necesary struct to join different sources
type Mux struct {
	channel chan Event
	sources []Source
}

// NewMux returns a new mux
func NewMux(sources ...Source) (*Mux, error) {
	generalChannel := make(chan Event, 1)
	m := &Mux{channel: generalChannel, sources: sources}
	for _, s := range m.sources {
		go func(s Source) {
			for event := range *s.Events() {
				m.channel <- event
			}
		}(s)
	}
	return m, nil
}

// RunAllSources run all the sources at the same time
func (m *Mux) RunAllSources() error {
	for _, s := range m.sources {
		if err := s.Run(); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// Events returns the channel where arrive the all the events from the muxed sources
func (m *Mux) Events() chan Event {
	return m.channel
}
