package services

import (
	"fmt"
	"testing"
)

func TestMatrix(t *testing.T) {
    cliente := NewClienteMatrix("kaduhod1", PROVEDOR_PADRAO, "cX2@hDEADv5kbpE")
    cliente.SetIdDeUsuario()
    if err := cliente.Login(); err != nil {
        t.Fail()
    }
    if err := cliente.BuscarSalas(); err != nil {
        t.Fail()
    }
    fmt.Println("SALA", cliente.Salas[0])
    _, err := cliente.BuscarMensagensDeSala(cliente.Salas[0])
    if err != nil {
        fmt.Println(err)
        t.Fail()
    }
    //fmt.Println(res)
    t.Fail()
}
