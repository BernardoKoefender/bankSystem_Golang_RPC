package calc

////////// IMPORTANT
// disable go modules to avoid anoying compilation errors:
// go env -w GO111MODULE=off




//Tipos exportaveis (inicio com letra maiuscula) podem ser registrados por servidores
type Args struct {
	A, B float64
}

type Arith float64

/*
	Metodos devem:
	-Pertencer a um tipo exportavel (Arith neste caso) e ser exportaveis
	-Possuir dois parametros de entrada. O primeiro pode ser qualquer tipo
	exportavel ou tipo nativo de go. O segundo de ser obrigatoriamente
	um ponteiro. O segundo argumento eh usado para o retorno do metodo.
	-Retornar um erro. Se for retornado algo alem de nil o cliente recebera 
	apenas o erro, sem o ponteiro de reply
*/
func (a *Arith) Sum(args *Args, reply *float64) error {
	*reply = args.A + args.B
	return nil
}

func (a *Arith) Sub(args *Args, reply *float64) error {
	*reply = args.A - args.B
	return nil
}

func (a *Arith) Mult(args *Args, reply *float64) error {
	*reply = args.A * args.B
	return nil
}

func (a *Arith) Div(args *Args, reply *float64) error {
	*reply = args.A / args.B
	return nil
}
