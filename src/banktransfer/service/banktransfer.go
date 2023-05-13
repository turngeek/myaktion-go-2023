package service

import (
	"context"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/turngeek/myaktion-go-2023/src/banktransfer/grpc/banktransfer"
	"github.com/turngeek/myaktion-go-2023/src/banktransfer/kafka"
	"google.golang.org/protobuf/types/known/emptypb"
)

type BankTransferService struct {
	banktransfer.BankTransferServer
	keyGenerator      *KeyGenerator
	transactionWriter kafka.TransactionWriter
}

func NewBankTransferService() *BankTransferService {
	return &BankTransferService{keyGenerator: NewKeyGenerator()}
}

func (s *BankTransferService) TransferMoney(_ context.Context, transaction *banktransfer.Transaction) (*emptypb.Empty, error) {
	entry := log.WithField("transaction", transaction)
	entry.Info("Received transaction")
	s.processTransaction(transaction)
	return &emptypb.Empty{}, nil
}

func (s *BankTransferService) ProcessTransactions(stream banktransfer.BankTransfer_ProcessTransactionsServer) error {
	return func() error {
		r := kafka.NewTransactionReader()
		for {
			err := r.Read(stream.Context(), func(transaction *banktransfer.Transaction) error {
				id := transaction.Id
				entry := log.WithField("transaction", transaction)
				entry.Info("Sending transaction to the client")
				if err := stream.Send(transaction); err != nil {
					return fmt.Errorf("error sending transaction: %w", err)
				}
				entry.Info("Transaction sent. Waiting for processing response")
				response, err := stream.Recv()
				if err != nil {
					return fmt.Errorf("error receiving transaction response: %w", err)
				}
				if response.Id != id {
					// NOTE: this is just a guard and not happening as transaction is local per connection
					return errors.New("received processing response of a different transaction")
				}
				entry.Info("Processing response received")
				return nil
			})
			if err != nil {
				log.WithError(err).Error("Failed to read transaction")
				break
			}
		}
		if err := r.Close(); err != nil {
			log.WithError(err).Error("Failed to close transaction reader")
			return nil
		}
		log.Info("Transaction reader successfully closed")
		return nil
	}()
}

func (s *BankTransferService) Start() {
	log.Info("Starting banktransfer service")
	s.transactionWriter = kafka.NewTransactionWriter()
	log.Info("Successfully created transaction writer")
}

func (s *BankTransferService) Stop() {
	log.Info("Stopping banktransfer service")
	s.transactionWriter.Close()
	log.Info("Successfully closed connection to transaction writer")
}

func (s *BankTransferService) processTransaction(transaction *banktransfer.Transaction) {
	entry := log.WithField("transaction", transaction)
	// NOTE: we need to copy the transaction as we are using a own go routine to process the transaction
	// this produces a `go vet` warning as gRPC messages contain a mutex that is also copied but not used.
	// therefore this warning can be ignored
	go func(transaction banktransfer.Transaction) {
		entry.Info("Start processing transaction")
		transaction.Id = s.keyGenerator.getUniqueId()
		if err := s.transactionWriter.Write(&transaction); err != nil {
			entry.WithError(err).Error("Can't write transaction to transaction writer")
			return
		}
		entry.Info("Transaction forwarded to transaction writer. Processing transaction finished")
	}(*transaction)
}
