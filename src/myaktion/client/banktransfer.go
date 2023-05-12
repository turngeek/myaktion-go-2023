package client

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc"
)

var (
	bankTransferTarget = os.Getenv("BANKTRANSFER_CONNECT")
)

func GetBankTransferConnection(ctx context.Context) (*grpc.ClientConn, error) {
	var err error
	log.WithFields(log.Fields{
		"target": bankTransferTarget,
	}).Infoln("Connecting to banktransfer service")
	var conn *grpc.ClientConn
	conn, err = grpc.DialContext(ctx, bankTransferTarget, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	return conn, nil
}
