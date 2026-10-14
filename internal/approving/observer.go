package approving

type Observer interface {
	Approved()
	Refused()
	Failed()
}
