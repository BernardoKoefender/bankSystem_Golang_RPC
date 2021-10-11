package main

import (
    "fmt"
    "log"
    "net/rpc"
	"os"
    "strconv"
	"bufio"
	"strings"
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
        log.Fatal("Dialing: ", err)
    }

	//Variavel para receber os resultados
	var reply float64

	//Buffer para ler do terminal
	reader := bufio.NewReader(os.Stdin)

	for {
		//Leitura de uma linha do terminal
		text, e := reader.ReadString('\n')
		if e != nil {
			log.Fatal(e)
		}
		text = strings.Replace(text, "\r\n", "", -1) //Windows
		//text = strings.Replace(text, "\n", "", -1) //Unix

		//Separa a linha pelos espacos em branco
		input := strings.Split(text, " ")

		//Converte string para float
		a, e1 := strconv.ParseFloat(input[1],64)
		if e1 != nil {
			log.Fatal(e1)
		}
		b, e2 := strconv.ParseFloat(input[2],64)
		if e2 != nil {
			log.Fatal(e2)
		}

		//Cria a struct para enviar para o servidor
		args := calc.Args{A: a, B: b}

		/*
        Call chama um metodo atrves da conexao estabelecida
        anteriormente. Os parametros devem ser:
        -Metodo a ser chamado no servidor no formato 'Tipo.Nome'.
        Este parametro deve ser uma string
        -Primeiro argumento do metodo
        -Segundo argumento do metodo(ponteiro para receber a resposta)
    */
		err = c.Call("Arith." + input[0], args, &reply)
		if err != nil {
			log.Fatal("Arith error: ", err)
		}
		fmt.Printf("Result = %f\n", reply)
	}
}