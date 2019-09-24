package nego

import (
	"errors"

	"github.com/wizcas/mudever.svc/packet"
)

// MockCommittee is for unit tests only
type MockCommittee struct {
	bad    bool
	Packet packet.Packet
	Err    error
}

// ErrMockCommit is for mocking the error on handler's committments.
var ErrMockCommit = errors.New("MOCK")

// NewMockCommittee returns a committee pointer that just for unit testing，
// with the 'bad' parameter which causes an ErrMock on any Commit().
func NewMockCommittee(bad bool) *MockCommittee {
	return &MockCommittee{bad: bad}
}

// Commit for unit test, which causes an ErrMockCommit on a bad committee
func (s *MockCommittee) Commit(p packet.Packet) error {
	if s.bad {
		return ErrMockCommit
	}
	s.Packet = p
	return nil
}

// GotError for unit tests, which records error in the Err field.
func (s *MockCommittee) GotError(err error) {
	s.Err = err
}
