package service

import (
	"context"
	"math/rand"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/turngeek/myaktion-go-2023/src/banktransfer/grpc/banktransfer"
	"google.golang.org/protobuf/types/known/emptypb"
)

type BankTransferService struct {
	banktransfer.BankTransferServer
	counter int32
	queue   chan banktransfer.Transaction
}

func NewBankTransferService() *BankTransferService {
	rand.Seed(time.Now().UnixNano())
	return &BankTransferService{counter: 1, queue: make(chan banktransfer.Transaction)}
}

func (s *BankTransferService) TransferMoney(_ context.Context, transaction *banktransfer.Transaction) (*emptypb.Empty, error) {
	entry := log.WithField("transaction", transaction)
	entry.Info("Received transaction")
	s.processTransaction(transaction)
	return &emptypb.Empty{}, nil
}

func (s *BankTransferService) ProcessTransactions(stream banktransfer.BankTransfer_ProcessTransactionsServer) error {
	return func() error {
		for {
			select {
			case <-stream.Context().Done():
				log.Info("Watching transactions cancelled from the client side")
				return nil
			case transaction := <-s.queue:
				id := transaction.Id
				entry := log.WithField("transaction", transaction)
				entry.Info("Sending transaction to the client")
				if err := stream.Send(&transaction); err != nil {
					// TODO: we need to recover from this if communication is not up
					entry.WithError(err).Error("Sending transaction failed")
					return err
				}
				entry.Info("Transaction sent. Waiting for processing response")
				response, err := stream.Recv()
				if err != nil {
					// TODO: we need to recover from this if communication is not up
					entry.WithError(err).Error("Receiving processing response failed")
					return err
				}
				if response.Id != id {
					// NOTE: this is just a guard and not happening as transaction is local per connection
					entry.Error("Received processing response of a different transaction")
				} else {
					entry.Info("Processing response received")
				}
			}
		}
	}()
}

func (s *BankTransferService) Close() {
	close(s.queue)
}

func (s *BankTransferService) processTransaction(transaction *banktransfer.Transaction) {
	entry := log.WithField("transaction", transaction)
	go func(transaction banktransfer.Transaction) {
		entry.Info("Start processing transaction")
		transaction.Id = atomic.AddInt32(&s.counter, 1)
		time.Sleep(time.Duration(rand.Intn(9)+1) * time.Second)
		s.queue <- transaction
		entry.Info("Processing transaction finished")
	}(*transaction)
}