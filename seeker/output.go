package seeker

//Output is an output struct that knows how to take json and turn it into
// the fields of itself
type Output interface {
	ReceiveJSON([]byte) error
}
