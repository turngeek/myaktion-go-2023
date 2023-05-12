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

const retryTime = 5 * time.Second

type BankTransferService struct {
	banktransfer.BankTransferServer
	counter    int32
	queue      chan banktransfer.Transaction
	retryQueue chan banktransfer.Transaction
	stop       chan struct{}
}

func NewBankTransferService() *BankTransferService {
	rand.Seed(time.Now().UnixNano())
	return &BankTransferService{
		counter:    0,
		queue:      make(chan banktransfer.Transaction),
		retryQueue: make(chan banktransfer.Transaction),
		stop:       make(chan struct{})}
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
					s.requeueTransaction(&transaction)
					entry.WithError(err).Warn("Sending transaction failed")
					return err
				}
				entry.Info("Transaction sent. Waiting for processing response")
				response, err := stream.Recv()
				if err != nil {
					s.requeueTransaction(&transaction)
					entry.WithError(err).Warn("Receiving processing response failed")
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

func (s *BankTransferService) Start() {
	log.Info("Starting banktransfer service")
	go func() {
		for {
			select {
			case <-s.stop:
				break
			case transaction := <-s.retryQueue:
				s.queue <- transaction
			}
		}
	}()
}

func (s *BankTransferService) Stop() {
	log.Info("Stopping banktransfer service")
	close(s.stop)
	// XXX: we leave these channels open on purpose as the go routines of processTransaction and requeueTransaction
	// might still enqueue transactions after calling Stop
	//close(s.queue)
	//close(s.retryQueue)
}

func (s *BankTransferService) processTransaction(transaction *banktransfer.Transaction) {
	entry := log.WithField("transaction", transaction)
	go func(transaction banktransfer.Transaction) {
		entry.Info("Start processing transaction")
		transaction.Id = s.getUniqueId()
		time.Sleep(time.Duration(rand.Intn(9)+1) * time.Second)
		s.queue <- transaction
		entry.Info("Processing transaction finished")
	}(*transaction)
}

// NOTE: we prevent retrying a transaction endlessly. Usually it's best practice to move
// transactions to an error queue after failing for a couple of times
func (s *BankTransferService) requeueTransaction(transaction *banktransfer.Transaction) {
	entry := log.WithField("transaction", transaction)
	go func(transaction banktransfer.Transaction) {
		entry.Infof("Requeuing transaction. Wait for %f seconds", retryTime.Seconds())
		time.Sleep(retryTime)
		s.retryQueue <- transaction
		entry.Info("Requeued transaction")
	}(*transaction)
}

func (s *BankTransferService) getUniqueId() int32 {
	return atomic.AddInt32(&s.counter, 1)
}
