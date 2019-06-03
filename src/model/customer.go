package model

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Nhanderu/brdoc"
)

type Customer struct {
	Cpf                 string
	Private, Incompleto int64
	DataUltimaCompra    time.Time
	TicketMedio         float64
	TicketUltimaCompra  float64
	CnpjMaisFrequente,
	CnpjUltimaCompra string
}

func DropTableCustomer(db *sql.DB) error {
	_, err := db.Exec(`DROP TABLE IF EXISTS CLIENTE`)
	if err != nil {
		log.Fatalf("Error droping table ERROR: %s", err)
		return err
	}

	return nil
}

func CreateTableCustomer(db *sql.DB) error {
	//
	DropTableCustomer(db)

	//Cria tabelas
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS CLIENTE ( id serial, cpf text, cpf_valido bool, private int, incompleto int, data_ultima_compra date, ticket_medio numeric(15,2), ticket_ultima_compra numeric(15,2), cnpj_mais_frequente text, cnpj_mais_frequente_valido bool, cnpj_ultima_compra text, cnpj_ultima_compra_valido bool)`)
	if err != nil {
		log.Fatalf("Error creating table ERROR: %s", err)
		return err
	}

	return nil
}

func InsertRowCustomer(db *sql.DB, row Customer) {

	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	sqlStatement := `INSERT INTO CLIENTE(` +
		`  cpf ,  private,  incompleto,  data_ultima_compra,  ticket_medio,  ticket_ultima_compra,  cnpj_mais_frequente,  cnpj_ultima_compra) VALUES` +
		fmt.Sprintf(" ('%v', %v, %v, '%v', %v, %v, '%v', '%v') ",
			row.Cpf, row.Private, row.Incompleto, row.DataUltimaCompra.Format("2006-01-02"), row.TicketMedio, row.TicketUltimaCompra, row.CnpjMaisFrequente, row.CnpjUltimaCompra)
	_, err := db.ExecContext(ctx, sqlStatement)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func ValidateCustomerDocs(db *sql.DB) {
	rows, err := db.Query(`SELECT id, cpf, cnpj_mais_frequente, cnpj_ultima_compra FROM CLIENTE`)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	defer rows.Close()
	for rows.Next() {
		var id int
		var cpf string
		var cpfValido bool
		var cnpjMaisFrequente string
		var cnpjMaisFrequenteValido bool
		var cnpjUltimaCompra string
		var cnpjUltimaCompraValido bool
		err = rows.Scan(&id, &cpf, &cnpjMaisFrequente, &cnpjUltimaCompra)

		if brdoc.IsCPF(cpf) {
			cpfValido = true
		}

		if brdoc.IsCNPJ(cnpjMaisFrequente) {
			cnpjMaisFrequenteValido = true
		}

		if brdoc.IsCNPJ(cnpjUltimaCompra) {
			cnpjUltimaCompraValido = true
		}

		sqlStatement := `UPDATE cliente SET cpf = $2, cnpj_mais_frequente = $3, cnpj_ultima_compra = $4, ` +
			`cpf_valido = $5, cnpj_mais_frequente_valido = $6, cnpj_ultima_compra_valido = $7` +
			`WHERE id = $1`

		_, err = db.Exec(sqlStatement, id, ClearDocs(cpf), ClearDocs(cnpjMaisFrequente), ClearDocs(cnpjUltimaCompra),
			cpfValido, cnpjMaisFrequenteValido, cnpjUltimaCompraValido)
		if err != nil {
			log.Fatalf("error: %v", err)
			fmt.Println(err)
		}

	}
	err = rows.Err() // get any error encountered ing iteration

}

func ClearDocs(doc string) string {
	doc = strings.Replace(doc, ".", "", -1)
	doc = strings.Replace(doc, "-", "", -1)
	doc = strings.Replace(doc, "/", "", -1)
	return doc
}

func UpdateCustomerDocs() {

}
