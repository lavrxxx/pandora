package leveldb

import (
	"encoding/hex"
	"fmt"

	"github.com/golang/protobuf/proto"

	"pandora/pkg/pb"
	"pandora/pkg/utils/errors"
)

const (
	PREFIX_CERT = "cert-"
)

// PutCert put cert in leveldb
func (l *Leveldb) PutCert(cert *pb.Cert) error {
	k, _ := hex.DecodeString(fmt.Sprintf("%s%s", PREFIX_CERT, cert.ID))

	buf, err := proto.Marshal(cert)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := l.db.Put(k, buf, nil); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// LoadCert load cert from leveldb
func (l *Leveldb) LoadCert(id string) (*pb.Cert, error) {
	k, _ := hex.DecodeString(fmt.Sprintf("%s%s", PREFIX_CERT, id))

	buf, err := l.db.Get(k, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	cert := &pb.Cert{}

	if err := proto.Unmarshal(buf, cert); err != nil {
		return nil, errors.WithStack(err)
	}

	return cert, nil
}
