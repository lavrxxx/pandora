package events

import (
	"github.com/spacelavr/pandora/pkg/api/env"
	"github.com/spacelavr/pandora/pkg/broker"
	"github.com/spacelavr/pandora/pkg/types"
)

// Listen listen for events
func Listen() error {
	var (
		chSendCert = make(chan *types.Certificate)
	)

	if err := env.GetBroker().Publish(broker.SubjectCertificate, chSendCert); err != nil {
		return err
	}

	for {
		select {
		case cert := <-env.ReadCert():
			chSendCert <- cert
		}
	}
}
