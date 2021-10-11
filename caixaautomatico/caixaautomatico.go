package caixaautomatico

//  Processo CaixaAutomático: Solicita depósito, retirada
//e  consulta  de  saldo  em  conta  existente.  As  duas 
//primeiras   devem  ser  operações  garantidamente  não-
//idempotentes  (semântica  de  execução  exactely  once)
//mesmo que  ocorra algum erro na confirmação da operação
//(simular com injeção de falhas).



//cliente