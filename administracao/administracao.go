package administracao

//  Processo Administração: realiza abertura e fechamento
//de  contas  (para agências),  autentica  que  contas já 
//existem  (tanto para  agências e  caixas automáticos) e 
//também  executa  as  operações  de  manipulação  destas 
//contas (saques e depósitos). Deve garantir semântica de
//execução  exactely  once para  operações que sejam não-
//idempotentes;

import (
	"fmt"
)

// Account object
//	acc_id: an identifier for the account
//	name:	account's owner name
//	cash:	ammount of cash in account
type Account struct
{
	Id 		int
	Name	string
	Cash	float64
}


// Typeholder for account's variables
type Args struct
{
	Id 		int
	Name 	string
	Cash 	float64
}

var AccountList = []Account{
    {
    	Id: 12345,
        Name: "test1",
        Cash: 0.0,
    },{
    	Id: 54321,
        Name: "test2",
        Cash : 0.0,
    },
}

type Adm int

//Check if account exists.
//	reply = 1 if exists, else 0
func (a *Adm) AccountExists(acc_id int, reply *int) error{
	for i := 0; i < len(AccountList); i++{
		if AccountList[i].Id == acc_id{
			fmt.Printf("AccountExists: Account %d exists.\n", acc_id)
			*reply = 1
			return nil
		}
	}
	fmt.Printf("AccountExists: Account %d doesn't exists.\n", acc_id)
	*reply = 0
	return nil
}

//Create new account.
//	reply = 1 if created
func (a *Adm) CreateAccount(args *Args, reply *int) error {
	acc := Account{
		Id:		args.Id,
		Name:	args.Name,
		Cash:	args.Cash,
	}
	fmt.Printf("CreateAccount: Account %d-%s created.\n", args.Id, args.Name)
	AccountList = append(AccountList, acc)
	*reply = 1
	return nil
}

//Remove an account from AccountList.
//	reply = 1 if removed, else 0
func (a *Adm) RemoveAccount(acc_id int, reply *int) error {
	for i := 0; i < len(AccountList); i++{
		if AccountList[i].Id == acc_id{
			AccountList[i] = AccountList[len(AccountList)-1]
			AccountList = AccountList[:len(AccountList)-1]
			fmt.Printf("RemoveAccount: Account %d deleted.\n", acc_id)
			*reply = 1
			return nil
		}
	}
	fmt.Printf("RemoveAccount: Error: account not found.\n", acc_id)
	*reply = 0
	return nil
}

//Add funds to account
//  reply = 1 if successfull, 0 else
func (a *Adm) AddFunds(args *Args, reply *int) error{
	for i := 0; i < len(AccountList); i++{
		if AccountList[i].Id == args.Id{
			AccountList[i].Cash += args.Cash;
			fmt.Printf("AddFunds:      Added $%.2f to account %d. New total is: $%.2f\n", args.Cash, args.Id, AccountList[i].Cash)
			*reply = 1
			return nil
		}
	}
	fmt.Printf("AddFunds:      Error: account not found.\n")
	*reply = 0
	return nil
}

//Withdraw funds from account
//  reply = 1 if successfull, 0 else
func (a *Adm) WithdrawFunds(args *Args, reply *int) error{
	for i := 0; i < len(AccountList); i++{
		if AccountList[i].Id == args.Id{
			AccountList[i].Cash -= args.Cash;
			fmt.Printf("WithdraFunds: Removed $%.2f from account %d. New total is: $%.2f\n", args.Cash, args.Id, AccountList[i].Cash)
			*reply = 1
			return nil
		}
	}
	fmt.Printf("WithdrawFunds: Error: account not found.\n")
	*reply = 0
	return nil
}

//Check funds from account
//  reply = 1 if sucessfull, 0 else
func (a *Adm) CheckFunds(args *Args, reply *int) error{
	for i := 0; i < len(AccountList); i++{
		if AccountList[i].Id == args.Id{
			(*args).Cash = AccountList[i].Cash
			fmt.Printf("CheckFunds:    Checking funds of account %d. Is $%.2f\n", args.Id, args.Cash)
			*reply = 1
			return nil
		}
	}
	fmt.Printf("CheckFunds:    Error: account not found.\n", args.Id, args.Cash)
	*reply = 0
	return nil
}