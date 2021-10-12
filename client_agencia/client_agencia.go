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
    var reply int
    //Estrutura para enviar dados da conta
    args := administracao.Args{Id: 0, Name: "null", Cash: 0.0}

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
            err = c.Call("Adm.AccountExists", (&args).Id, &reply)
            if err != nil{
                log.Fatal("Adm: AccountExists error: ", err)
                continue
            }
            if reply == 1{
                fmt.Printf("Account %d already exsists. Aborting.\n", (&args).Id)
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
    			err = c.Call("Adm.CreateAccount", args, &reply)
    			if err != nil {
        			log.Fatal("Adm: CreateAccount error: ", err)
    			} else {
                    if reply == args.Key {
                        fmt.Printf("%6d: Account %d-%s created.\n", args.Key, args.Id, args.Name);
                    } else if reply == 0 {
                        fmt.Printf("%6d: Account %d-%s not created. Unused reply code %d.\n", args.Key, args.Id, args.Name, reply);
                    } else if reply == 1 {
                        fmt.Printf("%6d: Account %d-%s not created. Key already processed.\n", args.Key, args.Id, args.Name);
                    } else if reply == 2 {
                        fmt.Printf("%6d: Account %d-%s not created. Key not on pending list.\n", args.Key, args.Id, args.Name);
                    } else if reply == 3 {
                        fmt.Printf("%6d: Account %d-%s not created. Failed to process key.\n", args.Key, args.Id, args.Name);
                    } else {
                        fmt.Printf("%6d: Account %d-%s not created. Reply code %d unknown.\n", args.Key, args.Id, args.Name, reply);
                    }
                }
            } else {
                fmt.Printf("%6d: Invalid Key. Aborting.\n", args.Key)
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
    		
            // Check if account exists
            err = c.Call("Adm.AccountExists", (&args).Id, &reply)
            if err != nil{
                log.Fatal("Adm: AccountExists error: ", err)
                continue
            }
            if reply == 0{
                fmt.Printf("Account %d doesn't exsist. Aborting.\n", (&args).Id)
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
                err = c.Call("Adm.RemoveAccount", args, &reply)
                if err != nil {
                    log.Fatal("Adm: RemoveAccount error: ", err)
                } else {
                    if reply == args.Key {
                        fmt.Printf("%6d: Account %d-%s removed.\n", args.Key, args.Id, args.Name);
                    } else if reply == 0 {
                        fmt.Printf("%6d: Account %d-%s not removed. Account not found.\n", args.Key, args.Id, args.Name);
                    } else if reply == 1 {
                        fmt.Printf("%6d: Account %d-%s not removed. Key already processed.\n", args.Key, args.Id, args.Name);
                    } else if reply == 2 {
                        fmt.Printf("%6d: Account %d-%s not removed. Key not on pending list.\n", args.Key, args.Id, args.Name);
                    } else if reply == 3 {
                        fmt.Printf("%6d: Account %d-%s not removed. Failed to process key.\n", args.Key, args.Id, args.Name);
                    } else {
                        fmt.Printf("%6d: Account %d-%s not removed. Reply code %d unknown.\n", args.Key, args.Id, args.Name, reply);
                    }
                }
            } else {
                fmt.Printf("%6d: Invalid Key. Aborting.\n", args.Key)
                continue
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
                fmt.Printf("Invalid Key %d. Aborting.\n", args.Key)
                continue
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
