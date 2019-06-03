package main

import (
	"database/sql"
	"log"
	"strings"

	"./control"
	"./db"
	"./lib/util"
	"./model"

	_ "github.com/lib/pq"
)

var connection *sql.DB // Database connection

func main() {
	filename := "./service/base_teste.txt"

	var content []string
	var row []string
	var err error
	content = nil
	content, err = util.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error reading file:", err)
	}

	//Inicia Serviço
	control.Init()

	//Conecta banco de dados / servico
	connection, err = db.Open()
	if err != nil {
		panic(err)
	}

	err = model.CreateTableCustomer(connection)
	if err != nil {
		panic(err)
	}

	for indice, linha := range content {

		if indice > 0 {
			row = strings.Fields(linha)
			customer := model.Customer{}
			customer.Cpf = row[0]
			customer.Private = util.ConvertStrToInt(row[1])
			customer.Incompleto = util.ConvertStrToInt(row[2])
			customer.DataUltimaCompra = util.ConvertStrToDate(row[3])
			customer.TicketMedio = util.ConvertStrToFloat(row[4])
			customer.TicketUltimaCompra = util.ConvertStrToFloat(row[5])
			customer.CnpjMaisFrequente = row[6]
			customer.CnpjUltimaCompra = row[7]
			model.InsertRowCustomer(connection, customer)
		}
	}

	model.ValidateCustomerDocs(connection)
}
