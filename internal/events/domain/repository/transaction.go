package repository

import "context"

// Repositories agrupa todos os repositórios disponíveis dentro de uma transação
type Repositories struct {
	Events     EventRepository
	Activities ActivityRepository
	CheckIns   CheckInRepository
}

// TransactionProvider gerencia transações de banco de dados
type TransactionProvider interface {
	// Transact executa a função fn dentro de uma transação
	// Se fn retornar erro, faz rollback; caso contrário, faz commit
	Transact(ctx context.Context, fn func(repos Repositories) error) error
}
