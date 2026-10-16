package notary

import "time"

type SponsorshipWindow struct {
	After UnixTime
	Until UnixTime
}

func (w SponsorshipWindow) Admit(now UnixTime) error {
	const op = "notary.SponsorshipWindow.Admit"

	if !w.After.IsZero() && now.Seconds() <= w.After.Seconds() {
		return Errf(
			KindNotDue,
			op,
			nil,
			"the court takes this call only after %s; the chain's clock reads %s",
			w.After.Time().Format(time.RFC3339),
			now.Time().Format(time.RFC3339),
		)
	}

	if !w.Until.IsZero() && now.Seconds() > w.Until.Seconds() {
		return Errf(
			KindCallRefused,
			op,
			nil,
			"the court stopped taking this call at %s; the chain's clock reads %s",
			w.Until.Time().Format(time.RFC3339),
			now.Time().Format(time.RFC3339),
		)
	}

	return nil
}
