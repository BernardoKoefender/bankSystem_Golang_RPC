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


// Typeholder for transation's variables
type Args struct
{
	Id 		int
	Name 	string
	Cash 	float64
	Key		int
	Msg		string
}

type Adm int

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

var PendingKeys []int
var ProcessedKeys []int
var NextKey int = 1000

//Check if account exists.
//	reply = 1 if exists, else 0
func (a *Adm) AccountExists(acc_id int, reply *int) error{
	for i := 0; i < len(AccountList); i++{
		if AccountList[i].Id == acc_id{
			//fmt.Printf("AccountExists: Account %d exists.\n", acc_id)
			*reply = 1
			return nil
		}
	}
	//fmt.Printf("AccountExists: Account %d doesn't exists.\n", acc_id)
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
//  reply = Key if successfull
//	reply = 0 if account not found
//	reply = 1 if key is already processed
//	reply = 2 if key is not on pending list
//	reply = 3 if failed to process key(this should never happen)
// 	reply = 4 unknown
func (a *Adm) AddFunds(args *Args, reply *int) error{
	ret_val := EvaluateTransaction(args)
	if ret_val == args.Key {
		for i := 0; i < len(AccountList); i++{
			if AccountList[i].Id == args.Id{
				AccountList[i].Cash += args.Cash;
				fmt.Printf("%6d: AddFunds:      Added $%.2f to account %d. New total is: $%.2f\n", args.Key, args.Cash, args.Id, AccountList[i].Cash)
				*reply = args.Key
				return nil
			}
		}
		fmt.Printf("%6d: AddFunds:      Error: account not found.\n", args.Key)
		*reply = 0
		return nil
	} else if ret_val == 1 {
		fmt.Printf("%6d: AddFunds:      This transaction is already processed. Aborting.\n", args.Key)
		*reply = 1
		return nil
	} else if ret_val == 2{
		fmt.Printf("%6d: AddFunds:      Key not found. Aborting.\n", args.Key)
		*reply = 2
		return nil
	} else if ret_val == 3 {
		fmt.Printf("%6d: AddFunds:      Failed to process key. Aborting.\n", args.Key)
		*reply = 3
		return nil
	}
	*reply = 4
	return nil
}

//Withdraw funds from account
//  reply = Key if successfull
//	reply = 0 if account not found
//	reply = 1 if key is already processed
//	reply = 2 if key is not on pending list
//	reply = 3 if failed to process key(this should never happen)
// 	reply = 4 unknown
func (a *Adm) WithdrawFunds(args *Args, reply *int) error{
	ret_val := EvaluateTransaction(args)
	if ret_val == args.Key {
		for i := 0; i < len(AccountList); i++{
			if AccountList[i].Id == args.Id{
				AccountList[i].Cash -= args.Cash;
				fmt.Printf("%6d: WithdrawFunds: Removed $%.2f from account %d. New total is: $%.2f\n", args.Key, args.Cash, args.Id, AccountList[i].Cash)
				*reply = args.Key
				return nil
			}
		}
		fmt.Printf("%6d: WithdrawFunds: Error: account not found.\n", args.Key)
		*reply = 0
		return nil
	} else if ret_val == 1 {
		fmt.Printf("%6d: WithdrawFunds: This transaction is already processed. Aborting.\n", args.Key)
		*reply = 1
		return nil
	} else if ret_val == 2{
		fmt.Printf("%6d: WithdrawFunds: Key not found. Aborting.\n", args.Key)
		*reply = 2
		return nil
	} else if ret_val == 3 {
		fmt.Printf("%6d: WithdrawFunds: Failed to process key. Aborting.\n", args.Key)
		*reply = 3
		return nil
	}
	*reply = 4
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

// Retrieved a new transaction key.
// The new key is stored in *reply
func (a *Adm) GetNewKey(args *Args, reply *int) error{
	NextKey = NextKey + 1
	PendingKeys = append(PendingKeys, NextKey)
	*reply = NextKey
	return nil
}


/////////////////////////////
//// Auxiliary functions ////
/////////////////////////////

// Processes transaction key, moving it from PendingKeys
//to ProcessedKeys
//	Return true  if key processed
//	Return false else
func ProcessKey(args *Args) bool{
	for i := 0; i < len(PendingKeys); i++{
		if PendingKeys[i] == args.Key{
			ProcessedKeys  = append(ProcessedKeys, PendingKeys[i])
			PendingKeys[i] = PendingKeys[len(PendingKeys)-1]
			PendingKeys    = PendingKeys[:len(PendingKeys)-1]
			return true
		}
	}
	return false
}

// Check if given key is on pending list
// Return true  if key is pending
// Return false if key not pending
func CheckPendingKey(args *Args) bool{
	for i := 0; i < len(PendingKeys); i++{
		if PendingKeys[i] == args.Key{
			return true
		}
	}
	return false
}

// Check if given key is on processed list
// Return true  if key is processed
// Return false if key is avaliable
func CheckProcessedKey(args *Args) bool{
	for i := 0; i < len(ProcessedKeys); i++{
		if ProcessedKeys[i] == args.Key{
			return true
		}
	};
	return false
}

// Check iif transaction(key) is executable
// Return key if executable
// Return 0 if *unused*
// Return 1 if key is already processed
// Return 2 if key is not on pending list
// Return 3 if failed to process key(this should never happen)
func EvaluateTransaction(args *Args) int{
	if CheckProcessedKey(args) == true{
		// Key has already been processed
		return 1
	} else {
		if CheckPendingKey(args) == true{
			if ProcessKey(args) == false{
				// Failed to process the key for whatever reason
				// This should never happen
				return 3
			}
		} else {
			// Key is not on pending keys list
			return 2
		}
	}
	return args.Key
}