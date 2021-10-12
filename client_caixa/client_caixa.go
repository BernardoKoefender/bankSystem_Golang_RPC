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
    var reply int64
    //Estrutura para enviar dados da conta
    args := administracao.Args{Id: 1337, Name: "robocop", Cash: 0.0}

    var quit bool = false;
    for quit == false{
        fmt.Println("\n-----------------------------")
    	fmt.Printf("ATM: Choose an option:\n0: Exit\n1: Add Funds\n2: Withdraw Funds\n3: Check balance\n\n")
    	
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


    		err = c.Call("Adm.AccountExists", (&args).Id, &reply)
    		if err != nil{
    			log.Fatal("Adm: AccountExists error: ", err)
    			continue
    		}
    		if reply == 0{
    			fmt.Printf("Account %d doesn't exsists. Aborting.\n", (&args).Id)
    			continue
    		} else if reply == 1{
    			err = c.Call("Adm.AddFunds", &args, &reply)
    			if err != nil{
	    			log.Fatal("Adm: AddFunds error: ", err)
    				continue
    			} else {
    				if reply == 1 {
    					fmt.Printf("Added $%.2f to account %d.", (&args).Cash, (&args).Id)
    				} else {
    					fmt.Printf("Deposit failed. Reply code: %d\n", &reply)
    				}
    			}
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


    		err = c.Call("Adm.AccountExists", (&args).Id, &reply)
    		if err != nil{
    			log.Fatal("Adm: AccountExists error: ", err)
    			continue
    		}
    		if reply == 0{
    			fmt.Printf("Account %d doesn't exsists. Aborting.\n", (&args).Id)
    			continue
    		} else if reply == 1{
    			err = c.Call("Adm.WithdrawFunds", &args, &reply)
    			if err != nil{
	    			log.Fatal("Adm: WithdrawFunds error: ", err)
    				continue
    			} else {
    				if reply == 1 {
    					fmt.Printf("Removed $%.2f from account %d.", (&args).Cash, (&args).Id)
    				} else {
    					fmt.Printf("Withdraw failed. Reply code: %d\n", &reply)
    				}
    			}
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
    			err = c.Call("Adm.CheckFunds", &args, &reply)
    			if err != nil{
	    			log.Fatal("Adm: CheckFunds error: ", err)
    				continue
    			} else {
    				if reply == 1 {
    					fmt.Printf("Checking funds of account %d. Total is: $%.2f", (&args).Id, (&args).Cash)
    				} else {
    					fmt.Printf("Fund check failed. Reply code: %d\n", reply)
    				}
    			}
    		}
    	//////////// UNKNOWN CODE
    	} else{
    		fmt.Println("Unknown operation.")
    	}
    }
   	fmt.Printf("Shutdown: client_agencia.go\n")
}
