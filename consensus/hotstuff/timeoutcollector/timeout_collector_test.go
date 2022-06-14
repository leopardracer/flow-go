package timeoutcollector

import (
	"errors"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/onflow/flow-go/consensus/hotstuff/helper"
	"github.com/onflow/flow-go/consensus/hotstuff/mocks"
	"github.com/onflow/flow-go/consensus/hotstuff/model"
	"github.com/onflow/flow-go/model/flow"
)

func TestTimeoutCollector(t *testing.T) {
	suite.Run(t, new(TimeoutCollectorTestSuite))
}

// TimeoutCollectorTestSuite is a test suite for testing TimeoutCollector. It stores mocked
// state internally for testing behavior.
type TimeoutCollectorTestSuite struct {
	suite.Suite

	view              uint64
	notifier          *mocks.Consumer
	processor         *mocks.TimeoutProcessor
	onNewQCDiscovered *mocks.OnNewQCDiscovered
	onNewTCDiscovered *mocks.OnNewTCDiscovered
	collector         *TimeoutCollector
}

func (s *TimeoutCollectorTestSuite) SetupTest() {
	s.view = 1000
	s.notifier = &mocks.Consumer{}
	s.processor = &mocks.TimeoutProcessor{}
	s.onNewQCDiscovered = mocks.NewOnNewQCDiscovered(s.T())
	s.onNewTCDiscovered = mocks.NewOnNewTCDiscovered(s.T())

	s.onNewQCDiscovered.On("Execute", mock.Anything).Maybe()
	s.onNewTCDiscovered.On("Execute", mock.Anything).Maybe()

	s.collector = NewTimeoutCollector(s.view, s.notifier, s.processor, s.onNewQCDiscovered.Execute, s.onNewTCDiscovered.Execute)
}

// TestView tests that `View` returns the same value that was passed in constructor
func (s *TimeoutCollectorTestSuite) TestView() {
	require.Equal(s.T(), s.view, s.collector.View())
}

// TestAddTimeout_HappyPath tests that process in happy path executed by multiple workers deliver expected results
// all operations should be successful, no errors expected
func (s *TimeoutCollectorTestSuite) TestAddTimeout_HappyPath() {
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			timeout := helper.TimeoutObjectFixture(helper.WithTimeoutObjectView(s.view))
			s.processor.On("Process", timeout).Return(nil).Once()
			err := s.collector.AddTimeout(timeout)
			require.NoError(s.T(), err)
		}()
	}

	wg.Wait()
	s.processor.AssertExpectations(s.T())
}

// TestAddTimeout_DoubleTimeout tests that submitting two different timeouts for same view ends with reporting
// double timeout to notifier which can be slashed later.
func (s *TimeoutCollectorTestSuite) TestAddTimeout_DoubleTimeout() {
	timeout := helper.TimeoutObjectFixture(helper.WithTimeoutObjectView(s.view))
	s.processor.On("Process", timeout).Return(nil).Once()
	err := s.collector.AddTimeout(timeout)
	require.NoError(s.T(), err)

	otherTimeout := helper.TimeoutObjectFixture(helper.WithTimeoutObjectView(s.view),
		helper.WithTimeoutObjectSignerID(timeout.SignerID))

	s.notifier.On("OnDoubleTimeoutDetected", timeout, otherTimeout).Once()

	err = s.collector.AddTimeout(otherTimeout)
	require.NoError(s.T(), err)
	s.notifier.AssertCalled(s.T(), "OnDoubleTimeoutDetected", timeout, otherTimeout)
	s.processor.AssertNumberOfCalls(s.T(), "Process", 1)
}

// TestAddTimeout_RepeatedTimeout checks that repeated timeouts are silently dropped without any errors.
func (s *TimeoutCollectorTestSuite) TestAddTimeout_RepeatedTimeout() {
	timeout := helper.TimeoutObjectFixture(helper.WithTimeoutObjectView(s.view))
	s.processor.On("Process", timeout).Return(nil).Once()
	err := s.collector.AddTimeout(timeout)
	require.NoError(s.T(), err)
	err = s.collector.AddTimeout(timeout)
	require.NoError(s.T(), err)
	s.processor.AssertNumberOfCalls(s.T(), "Process", 1)
}

// TestAddTimeout_TimeoutCacheException tests that submitting timeout object for view which is not designated for this
// collector results in ErrTimeoutForIncompatibleView.
func (s *TimeoutCollectorTestSuite) TestAddTimeout_TimeoutCacheException() {
	// incompatible view is an exception and not handled by timeout collector
	timeout := helper.TimeoutObjectFixture(helper.WithTimeoutObjectView(s.view + 1))
	err := s.collector.AddTimeout(timeout)
	require.ErrorIs(s.T(), err, ErrTimeoutForIncompatibleView)
	s.processor.AssertNotCalled(s.T(), "Process")
}

// TestAddTimeout_InvalidTimeout tests that sentinel errors while processing timeouts are correctly handled and reported
// to notifier, but exceptions are propagated to caller.
func (s *TimeoutCollectorTestSuite) TestAddTimeout_InvalidTimeout() {
	s.Run("invalid-timeout", func() {
		timeout := helper.TimeoutObjectFixture(helper.WithTimeoutObjectView(s.view))
		s.processor.On("Process", timeout).Return(model.NewInvalidTimeoutErrorf(timeout, "")).Once()
		s.notifier.On("OnInvalidTimeoutDetected", timeout).Once()
		err := s.collector.AddTimeout(timeout)
		require.NoError(s.T(), err)

		s.notifier.AssertCalled(s.T(), "OnInvalidTimeoutDetected", timeout)
	})
	s.Run("process-exception", func() {
		exception := errors.New("invalid-signature")
		timeout := helper.TimeoutObjectFixture(helper.WithTimeoutObjectView(s.view))
		s.processor.On("Process", timeout).Return(exception).Once()
		err := s.collector.AddTimeout(timeout)
		require.ErrorIs(s.T(), err, exception)
	})
}

// TestAddTimeout_TONotifications tests that TimeoutCollector in happy path reports the newest discovered QC and TC>
func (s *TimeoutCollectorTestSuite) TestAddTimeout_TONotifications() {
	qcCount := 100
	// generate QCs with increasing view numbers
	if s.view < uint64(qcCount) {
		s.T().Fatal("invalid test configuration")
	}

	*s.onNewQCDiscovered = mocks.OnNewQCDiscovered{}
	*s.onNewTCDiscovered = mocks.OnNewTCDiscovered{}

	var highestReportedQC *flow.QuorumCertificate
	s.onNewQCDiscovered.On("Execute", mock.Anything).Run(func(args mock.Arguments) {
		qc := args.Get(0).(*flow.QuorumCertificate)
		if highestReportedQC == nil || highestReportedQC.View < qc.View {
			highestReportedQC = qc
		}
	})

	lastViewTC := helper.MakeTC(helper.WithTCView(s.view - 1))
	s.onNewTCDiscovered.On("Execute", lastViewTC).Once()

	timeouts := make([]*model.TimeoutObject, 0, qcCount)
	for i := 0; i < qcCount; i++ {
		qc := helper.MakeQC(helper.WithQCView(uint64(i)))
		timeout := helper.TimeoutObjectFixture(func(timeout *model.TimeoutObject) {
			timeout.View = s.view
			timeout.NewestQC = qc
			timeout.LastViewTC = lastViewTC
		})
		timeouts = append(timeouts, timeout)
		s.processor.On("Process", timeout).Return(nil).Once()
	}

	expectedHighestQC := timeouts[len(timeouts)-1].NewestQC

	// shuffle timeouts in random order
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(timeouts), func(i, j int) {
		timeouts[i], timeouts[j] = timeouts[j], timeouts[i]
	})

	var wg sync.WaitGroup
	wg.Add(len(timeouts))
	for _, timeout := range timeouts {
		go func(timeout *model.TimeoutObject) {
			defer wg.Done()
			err := s.collector.AddTimeout(timeout)
			require.NoError(s.T(), err)
		}(timeout)
	}
	wg.Wait()

	s.onNewTCDiscovered.AssertCalled(s.T(), "Execute", lastViewTC)
	require.Equal(s.T(), expectedHighestQC, highestReportedQC)
}
