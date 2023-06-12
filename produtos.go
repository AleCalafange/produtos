package main

import (
    "database/sql"
    "net/http"
	"github.com/go-chi/chi/v5"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
)

type Recurso struct {
    db *sql.DB
}

func(met Recurso) CriaProduto(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()    
    produto := r.Form.Get("produto")
    valor := r.Form.Get("valor")
    fmt.Println("Produto cadastrado", produto,"Valor cadastrado", valor)
    w.Write([]byte("Você cadastrou corretamente"))

    valor, err := strconv.Atoi("1230")
    if err !=  nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        w.Header().Set("Error", err.Error())
    }
    
    cadastro := `INSERT INTO produtos(produto,valor) VALUES(?,?);`
    query, err := met.db.Prepare(cadastro)
    if err != nil {
        fmt.Println(err)
        return
    }
    _, err = query.Exec(produto,valor)
    if err != nil {
        fmt.Println(err)
        return
    }
}

func main() {
    db, err := sql.Open("sqlite3", "produtos.db")
    criaTabelaProduto(db)
    if err != nil{
        fmt.Println(err)
        return
    }

    resource := Recurso{db}

    router := chi.NewRouter()
    router.Get("/getProduto", getProduto)
    router.Post("/criaProduto", resource.CriaProduto)
    router.Get("/listaProdutos", listaProdutos)

    http.ListenAndServe(":3333", router)

}

func criaTabelaProduto(db *sql.DB){
    produtos_tabela := `CREATE TABLE produtos (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "produto" TEXT,
    "valor" FLOAT);`
    query, err := db.Prepare(produtos_tabela)
    if err != nil {
       fmt.Println(err)
        return
    }
    query.Exec()
    fmt.Println("Tabela criada com sucesso")
}

func getProduto(w http.ResponseWriter, r *http.Request){
    w.Write([]byte("Mostra Produto"))
    fmt.Println("pega Produto")
}

func listaProdutos(w http.ResponseWriter, r *http.Request){
    w.Write([]byte("foi pro curl lista"))
    fmt.Println("lista de produtos")
}