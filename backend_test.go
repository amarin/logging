package logging_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/amarin/logging"
)

func TestBackend_MultipleInstances(t *testing.T) {
	maxBackends := 1000
	config := logging.CurrentConfig()
	wg := new(sync.WaitGroup)
	wg.Add(maxBackends)

	for i := 0; i < maxBackends; i++ {
		go func(done *sync.WaitGroup) {
			require.NotPanics(t, func() { new(logging.Backend).MustInit(*config) })
			done.Done()
		}(wg)
	}

	wg.Wait()
}
