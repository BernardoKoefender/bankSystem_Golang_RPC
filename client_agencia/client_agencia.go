package main

//  Processo Agência:  Solicita abertura,  autenticação e 
//fechamento de contas e também pode  solicitar depósito,
//retirada e  consulta  de  saldo em  conta  existente. A
//abertura  de  conta, o  depósito e a retirada devem ser 
//operações  garantidamente  não-idempotentes  (semântica
//de execução exactely once);

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
    	fmt.Printf("Agency: Choose an option:\n0: Exit\n1: Create account\n2: Delete account\n3: Authenticate account\n4: Add Funds\n5: Withdraw Funds\n6: Check balance\n\n")
    	
    	input := bufio.NewScanner(os.Stdin)
    	input.Scan()
    	
    	if input.Text() == "0"{
    		quit = true

    	//////////// ACOUNT CREATION
    	} else if input.Text() == "1"{
    		fmt.Println("Creating account.")
    		fmt.Println("Enter account ID or 0 to abort:")
    		input = bufio.NewScanner(os.Stdin)
    		input.Scan()
    		if input.Text() == "0"{
    			continue
    		}
    		args.Id, err = strconv.Atoi(input.Text())

    		fmt.Println("Enter account name:")
    		input = bufio.NewScanner(os.Stdin)
    		input.Scan()
    		args.Name = input.Text()

    		args.Cash = 0.0

    		// Check if account exists
    		err = c.Call("Adm.AccountExists", args.Id, &reply)
    		if err != nil{
    			log.Fatal("Adm: AccountExists error: ", err)
    			continue
    		}
    		if reply == 0{
    			//If it doesn't exists, create
    			err = c.Call("Adm.CreateAccount", args, &reply)
    			if err != nil {
        			log.Fatal("Adm: CreateAccount error: ", err)
    			}
    			if reply == 1 {
    				fmt.Printf("Account %d-%s created.\n", args.Id, args.Name);
    			}else{
    				fmt.Printf("Account %d-%s not created, reply code: %d.\n", args.Id, args.Name, reply);
    			}
    		} else if reply == 1{
    			fmt.Printf("Account %d already exists. Aborting.", args.Id)
    			continue
    		}

    	//////////// ACOUNT DELETION
    	} else if input.Text() == "2"{
    		fmt.Println("Delete Acount")
    		fmt.Println("Enter account ID or 0 to abort:")
    		input = bufio.NewScanner(os.Stdin)
    		input.Scan()
    		if input.Text() == "0"{
    			continue
    		}
    		args.Id, err = strconv.Atoi(input.Text())
    		
    		err = c.Call("Adm.AccountExists", args.Id, &reply)
    		if err != nil{
    			log.Fatal("Adm: AccountExists error: ", err)
    			continue
    		}
    		if reply == 0{
    			fmt.Printf("Account %d doesn't exsists. Aborting.\n", args.Id)
    			continue
    		} else if reply == 1{
    		    err = c.Call("Adm.RemoveAccount", args.Id, &reply)
    			if err != nil{
	    			log.Fatal("Adm: RemoveAccount error: ", err)
    				continue
    			} else {
    				if reply == 1 {
    					fmt.Printf("Account %d removed.", args.Id)
    				} else {
    					fmt.Printf("Account %d not removed. Reply code: %d\n", args.Id, &reply)
    				}
    			}
    		}

		//////////// ACOUNT AUTHENTICATION
    	} else if input.Text() == "3"{
    		fmt.Println("Authenticate account")
    		fmt.Println("Enter account ID or 0 to abort:")
    		input = bufio.NewScanner(os.Stdin)
    		input.Scan()
    		if input.Text() == "0"{
    			continue
    		}
    		args.Id, err = strconv.Atoi(input.Text())
    		err = c.Call("Adm.AccountExists", args.Id, &reply)
    		if err != nil{
    			log.Fatal("Adm: AccountExists error: ", err)
    			continue
    		}
    		if reply == 0{
    			fmt.Printf("Account %d doesn't exsist.\n", args.Id)
    		} else if reply == 1{
    			fmt.Printf("Account %d exsist.\n", args.Id)
    		} else {
    			fmt.Printf("Account authentication: unknown reply code %d.\n", &reply)
    		}


    	} else if input.Text() == "4"{
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


    		err = c.Call("Adm.AccountExists", args.Id, &reply)
    		if err != nil{
    			log.Fatal("Adm: AccountExists error: ", err)
    			continue
    		}
    		if reply == 0{
    			fmt.Printf("Account %d doesn't exsists. Aborting.\n", args.Id)
    			continue
    		} else if reply == 1{
    			err = c.Call("Adm.AddFunds", args, &reply)
    			if err != nil{
	    			log.Fatal("Adm: AddFunds error: ", err)
    				continue
    			} else {
    				if reply == 1 {
    					fmt.Printf("Added $%.2f to account %d.", args.Cash, args.Id)
    				} else {
    					fmt.Printf("Deposit failed. Reply code: %d\n", &reply)
    				}
    			}
    		}

    	} else if input.Text() == "5"{
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


    		err = c.Call("Adm.AccountExists", args.Id, &reply)
    		if err != nil{
    			log.Fatal("Adm: AccountExists error: ", err)
    			continue
    		}
    		if reply == 0{
    			fmt.Printf("Account %d doesn't exsists. Aborting.\n", args.Id)
    			continue
    		} else if reply == 1{
    			err = c.Call("Adm.WithdrawFunds", args, &reply)
    			if err != nil{
	    			log.Fatal("Adm: WithdrawFunds error: ", err)
    				continue
    			} else {
    				if reply == 1 {
    					fmt.Printf("Removed $%.2f from account %d.", args.Cash, args.Id)
    				} else {
    					fmt.Printf("Withdraw failed. Reply code: %d\n", &reply)
    				}
    			}
    		}
    	} else if input.Text() == "6"{
    		fmt.Println("Check Funds")
			fmt.Println("Enter account ID or 0 to abort:")
    		input = bufio.NewScanner(os.Stdin)
    		input.Scan()
    		if input.Text() == "0"{
    			continue
    		}
    		args.Id, err = strconv.Atoi(input.Text())

			err = c.Call("Adm.AccountExists", args.Id, &reply)
    		if err != nil{
    			log.Fatal("Adm: AccountExists error: ", err)
    			continue
    		}
    		if reply == 0{
    			fmt.Printf("Account %d doesn't exsists. Aborting.\n", args.Id)
    			continue
    		} else if reply == 1{
    			err = c.Call("Adm.CheckFunds", &args, &reply)
    			if err != nil{
	    			log.Fatal("Adm: CheckFunds error: ", err)
    				continue
    			} else {
    				if reply == 1 {
    					fmt.Printf("Checking funds of account %d. Total is: $%.2f", args.Id, args.Cash)
    				} else {
    					fmt.Printf("Fund check failed. Reply code: %d\n", &reply)
    				}
    			}
    		}

    	} else{
    		fmt.Println("Unknown operation.")
    	}
    }
   	fmt.Printf("Shutdown: client_agencia.go\n")
}
