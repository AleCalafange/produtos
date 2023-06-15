package main

import (
    "database/sql"
    "net/http"
	"github.com/go-chi/chi/v5"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
    "strconv"
    "encoding/json"
)

type StructProduto struct{
    id int
    produto string
    valor float64
}

type Recurso struct {
    db *sql.DB
}

func(met Recurso) CriaProduto(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()    
    produto := r.Form.Get("produto")
    valor := r.Form.Get("valor")

    valorint, err := strconv.ParseFloat(valor, 64)
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
    _, err = query.Exec(produto,valorint)
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Println("Produto cadastrado", produto,"Valor cadastrado", valor)
    w.Write([]byte("Você cadastrou corretamente"))
}

func main() {
    db, err := sql.Open("sqlite3", "produtos.db")
    if err != nil{
        fmt.Println(err)
        return
    }
    criaTabelaProduto(db)

    resource := Recurso{db}

    router := chi.NewRouter()
    router.Get("/getProduto", getProduto)
    router.Post("/criaProduto", resource.CriaProduto)
    router.Get("/listaProdutos", resource.listaProdutos)

    http.ListenAndServe(":3333", router)

}

func criaTabelaProduto(db *sql.DB){
    produtos_tabela := `CREATE TABLE IF NOT EXISTS produtos (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    produto TEXT,
    valor FLOAT);`
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

func (met Recurso) listaProdutos(w http.ResponseWriter, r *http.Request){
    listadeProdutos := []StructProduto{}
    rows, err := met.db.Query("SELECT id, produto, valor FROM produtos")
    if err != nil{
        http.Error(w, err.Error(), http.StatusInternalServerError)
        w.Header().Set("Error", err.Error())
    }

    for rows.Next(){
        var produto StructProduto
        err := rows.Scan(&produto.id, &produto.produto, &produto.valor);
        if err !=  nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            w.Header().Set("Error", err.Error())
            return;
        }
        listadeProdutos = append(listadeProdutos)
    }
    resultado, err := json.Marshal(listadeProdutos)
    if err !=  nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        w.Header().Set("Error", err.Error())
    }
    w.Write([]byte(resultado))
    fmt.Println(resultado)
}