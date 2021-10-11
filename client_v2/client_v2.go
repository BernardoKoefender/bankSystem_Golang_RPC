package main

import (
    "fmt"
    "log"
    "net/rpc"
    "../calc"
)

func main() {
    /*
        Inicializa o cliente na porta 4040 do localhost
        utilizando o protocolo tcp. Se o servidor estiver
        em outra maquina deve ser utilizado IP:porta no 
        segundo parametro.
	*/
    c, err := rpc.DialHTTP("tcp", "localhost:4040")
    if err != nil {
        log.Fatal("Error dialing: ", err)
    }

    //Variavel para receber os resultados
    var reply float64
    //Estrutura para enviar os numeros
    args := calc.Args{A: 3, B: 4}

    /*
        Call chama um metodo atrves da conexao estabelecida
        anteriormente. Os parametros devem ser:
        -Metodo a ser chamado no servidor no formato 'Tipo.Nome'.
        Este parametro deve ser uma string
        -Primeiro argumento do metodo
        -Segundo argumento do metodo(ponteiro para receber a resposta)
    */
    err = c.Call("Arith.Mult", args, &reply)
    if err != nil {
        log.Fatal("Arith error: ", err)
    }
    fmt.Printf("Arith: %f*%f=%f\n", args.A, args.B, reply)

    err = c.Call("Arith.Div", args, &reply)
    if err != nil {
        log.Fatal("Arith error: ", err)
    }
    fmt.Printf("Arith: %f/%f=%f\n", args.A, args.B, reply)

    err = c.Call("Arith.Sum", args, &reply)
    if err != nil {
        log.Fatal("Arith error: ", err)
    }
    fmt.Printf("Arith: %f+%f=%f\n", args.A, args.B, reply)

    err = c.Call("Arith.Sub", args, &reply)
    if err != nil {
        log.Fatal("Arith error: ", err)
    }
    fmt.Printf("Arith: %f-%f=%f\n", args.A, args.B, reply)
}
