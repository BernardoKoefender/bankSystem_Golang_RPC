package main

//  Processo CaixaAutomático: Solicita depósito, retirada
//e  consulta  de  saldo  em  conta  existente.  As  duas 
//primeiras   devem  ser  operações  garantidamente  não-
//idempotentes  (semântica  de  execução  exactely  once)
//mesmo que  ocorra algum erro na confirmação da operação
//(simular com injeção de falhas).

import (
	"bufio"
	"os"
    "fmt"
    "log"
    "strconv"
    "net/rpc"
    "../administracao"
)

func main() {

    c, err := rpc.DialHTTP("tcp", "localhost:4040")
    if err != nil {
        log.Fatal("Error dialing: ", err)
    }

    //Variavel para receber os resultados
    var reply int

    // Estrutura criada só pro retorno da função CheckFunds
    replyCheckFunds := administracao.ReplyCheckFunds{Cash: 0.0, Reply: 0}

    //Estrutura para enviar dados da conta
    args := administracao.Args{Id: 0, Name: "null", Cash: 0.0, Key: 0, Msg: "0"}


    var quit bool = false;
    for quit == false{
        fmt.Println("\n-----------------------------")
    	fmt.Printf("ATM: Choose an option:\n0: Exit\n1: Add Funds\n2: Withdraw Funds\n3: Check balance\n4: Test\n5: Deposit(forced error)\n\n")
    	
    	input := bufio.NewScanner(os.Stdin)
    	input.Scan()
    	
    	if input.Text() == "0"{
    		quit = true

    	} else if input.Text() == "1"{
    		//////////// ADD FUNDS
			fmt.Println("Add Funds")
			fmt.Println("Enter account ID or 0 to abort:")
    		input = bufio.NewScanner(os.Stdin)
    		input.Scan()
    		if input.Text() == "0"{
    			continue
    		}
    		args.Id, err = strconv.Atoi(input.Text())

			fmt.Println("Enter ammount to deposit:")
    		input = bufio.NewScanner(os.Stdin)
    		input.Scan()
    		args.Cash, err = strconv.ParseFloat(input.Text(), 64)

    		// Check if account exists
    		err = c.Call("Adm.AccountExists", (&args).Id, &reply)
    		if err != nil{
    			log.Fatal("Adm: AccountExists error: ", err)
    			continue
    		}
    		if reply == 0{
    			fmt.Printf("Account %d doesn't exsists. Aborting.\n", (&args).Id)
    			continue
    		}

    		// Get a new transaction key
			err = c.Call("Adm.GetNewKey", &args, &reply)
    		if err != nil{
    			log.Fatal("Adm: GetNewKey error: ", err)
    			continue
    		}
    		// Process the transaction
    		if reply >= 1000{
    			args.Key = reply
    			err = c.Call("Adm.AddFunds", &args, &reply)
    			if err != nil{
	    			log.Fatal("Adm: AddFunds error: ", err)
    				continue
    			} else {
    				if reply == args.Key {
    					fmt.Printf("%6d: Added $%.2f to account %d.\n", (&args).Key, (&args).Cash, (&args).Id)
    				} else if reply == 0 {
    					fmt.Printf("%6d: Deposit failed. Account doesn't exists.\n", (&args).Key)
    				} else if reply == 1 {
    					fmt.Printf("%6d: Deposit failed. Transaction already processed.\n", (&args).Key)
    				} else {
    					fmt.Printf("%6d: Deposit failed. Unknown Reply code: %d\n", (&args).Key, reply)
    				}
    			}
    		} else {
    			fmt.Printf("%6d: Invalid Key. Aborting.\n", args.Key)
    			continue
    		}
    	} else if input.Text() == "2"{
    		//////////// WITHDRAW FUNDS
			fmt.Println("Withdraw Funds")
			fmt.Println("Enter account ID or 0 to abort:")
    		input = bufio.NewScanner(os.Stdin)
    		input.Scan()
    		if input.Text() == "0"{
    			continue
    		}
    		args.Id, err = strconv.Atoi(input.Text())

			fmt.Println("Enter ammount to withdraw:")
    		input = bufio.NewScanner(os.Stdin)
    		input.Scan()
    		args.Cash, err = strconv.ParseFloat(input.Text(), 64)

    		// Check if account exists
    		err = c.Call("Adm.AccountExists", (&args).Id, &reply)
    		if err != nil{
    			log.Fatal("Adm: AccountExists error: ", err)
    			continue
    		}
    		if reply == 0{
    			fmt.Printf("Account %d doesn't exsists. Aborting.\n", (&args).Id)
    			continue
    		}

    		// Get a new transaction key
			err = c.Call("Adm.GetNewKey", &args, &reply)
    		if err != nil{
    			log.Fatal("Adm: GetNewKey error: ", err)
    			continue
    		}

    		if reply >= 1000{
    			args.Key = reply
    			err = c.Call("Adm.WithdrawFunds", &args, &reply)
    			if err != nil{
	    			log.Fatal("Adm: WithdrawFunds error: ", err)
    				continue
    			} else {
    				if reply == args.Key {
    					fmt.Printf("%6d: Removed $%.2f from account %d.", (&args).Key, (&args).Cash, (&args).Id)
    				} else if reply == 0 {
    					fmt.Printf("%6d: Withdraw failed. Unused error code.\n", (&args).Key)
    				} else if reply == 1 {
    					fmt.Printf("%6d: Withdraw failed. Transaction already processed.\n", (&args).Key)
    				} else {
    					fmt.Printf("%6d: Withdraw failed. Unknown Reply code: %d\n", (&args).Key, reply)
    				}
    			}
    		} else {
    			fmt.Printf("Invalid Key %d. Aborting.\n", args.Key)
    			continue
    		}
    	} else if input.Text() == "3"{
    		//////////// CHECK BALANCE
			fmt.Println("Check Funds")
			fmt.Println("Enter account ID or 0 to abort:")
    		input = bufio.NewScanner(os.Stdin)
    		input.Scan()
    		if input.Text() == "0"{
    			continue
    		}
    		args.Id, err = strconv.Atoi(input.Text())

			err = c.Call("Adm.AccountExists", (&args).Id, &reply)
    		if err != nil{
    			log.Fatal("Adm: AccountExists error: ", err)
    			continue
    		}
    		if reply == 0{
    			fmt.Printf("Account %d doesn't exsists. Aborting.\n", (&args).Id)
    			continue
    		} else if reply == 1{
    			err = c.Call("Adm.CheckFunds", &args, &replyCheckFunds)
    			if err != nil{
	    			log.Fatal("Adm: CheckFunds error: ", err)
    				continue
    			} else {
    				if replyCheckFunds.Reply == 1 {
    					fmt.Printf("Checking funds of account %d. Total is: $%.2f", (&args).Id, replyCheckFunds.Cash)
    				} else {
    					fmt.Printf("Fund check failed. Reply code: %d\n", replyCheckFunds.Reply)
    				}
    			}
    		}
    	} else if input.Text() == "4"{
    		//////////// KEY TEST
			fmt.Println("Testing keys")
			
			err = c.Call("Adm.GetNewKey", &args, &reply)
    		if err != nil{
    			log.Fatal("Adm: GetNewKey error: ", err)
    			continue
    		}
    		if reply == 0{
    			fmt.Printf("Failed to get new key(%d). Aborting.\n", reply)
    			continue
    		} else if reply >= 1000{
    			fmt.Printf("Key %d retrieved.\n", reply)
    		}

    	} else if input.Text() == "5"{
    		//////////// ADD FUNDS WITH REPEATED KEY ERROR
			fmt.Println("Deposit - repeated key error test")
			fmt.Println("Enter account ID or 0 to abort:")
    		input = bufio.NewScanner(os.Stdin)
    		input.Scan()
    		if input.Text() == "0"{
    			continue
    		}
    		args.Id, err = strconv.Atoi(input.Text())

			fmt.Println("Enter ammount to deposit:")
    		input = bufio.NewScanner(os.Stdin)
    		input.Scan()
    		args.Cash, err = strconv.ParseFloat(input.Text(), 64)

    		// Check if account exists
    		err = c.Call("Adm.AccountExists", (&args).Id, &reply)
    		if err != nil{
    			log.Fatal("Adm: AccountExists error: ", err)
    			continue
    		}
    		if reply == 0{
    			fmt.Printf("Account %d doesn't exsists. Aborting.\n", (&args).Id)
    			continue
    		}

    		// Get a new transaction key
			err = c.Call("Adm.GetNewKey", &args, &reply)
    		if err != nil{
    			log.Fatal("Adm: GetNewKey error: ", err)
    			continue
    		}
    		// Process the transaction
    		if reply >= 1000{
    			args.Key = reply
    			err = c.Call("Adm.AddFunds", &args, &reply)
    			if err != nil{
	    			log.Fatal("Adm: AddFunds error: ", err)
    				continue
    			} else {
    				if reply == args.Key {
    					fmt.Printf("%6d: Added $%.2f to account %d.\n", (&args).Key, (&args).Cash, (&args).Id)
    				} else if reply == 0 {
    					fmt.Printf("%6d: Deposit failed. Account doesn't exists.\n", (&args).Key)
    				} else if reply == 1 {
    					fmt.Printf("%6d: Deposit failed. Transaction already processed.\n", (&args).Key)
    				} else {
    					fmt.Printf("%6d: Deposit failed. Unknown Reply code: %d\n", (&args).Key, reply)
    				}
    			}
    			err = c.Call("Adm.AddFunds", &args, &reply)
    			if err != nil{
	    			log.Fatal("Adm: AddFunds error: ", err)
    				continue
    			} else {
    				if reply == args.Key {
    					fmt.Printf("%6d: Added $%.2f to account %d.\n", (&args).Key, (&args).Cash, (&args).Id)
    				} else if reply == 0 {
    					fmt.Printf("%6d: Deposit failed. Account doesn't exists.\n", (&args).Key)
    				} else if reply == 1 {
    					fmt.Printf("%6d: Deposit failed. Transaction already processed.\n", (&args).Key)
    				} else {
    					fmt.Printf("%6d: Deposit failed. Unknown Reply code: %d\n", args.Key, reply)
    				}
    			}
    		} else {
    			fmt.Printf("Invalid Key %d. Aborting.\n", args.Key)
    			continue
    		}
    	//////////// UNKNOWN CODE
    	} else{
    		fmt.Println("Unknown operation.")
    	}
    }
   	fmt.Printf("Shutdown: client_agencia.go\n")
}
