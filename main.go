package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/correiostech/rastro"
	"github.com/correiostech/token"
)

var (
	urlToken  = flag.String("tk", "https://api.correios.com.br/token/v1/autentica", "url para obter token")
	urlRastro = flag.String("e", "https://api.correios.com.br/srorastro/v1/objetos/", "endpoint de rastro")
	objeto    = flag.String("o", "AA003039703BR", "objetos a sere rastreados")
	usuario   = flag.String("u", "", "nome de usuario de acesso às API dos Correios")
	senha     = flag.String("p", "", "senha do usuario de acesso às API dos Correios")
	tempo     = flag.Int("s", 15*60, "tempo para execucao, padrão 15 minutos")
)

func init() {
	*usuario = os.Getenv("USUARIO_API_CORREIOS")
	*senha = os.Getenv("SENHA_API_CORREIOS")
	go token.Start(*urlToken, *usuario, *senha)
}

func formataData(data string) (string, error) {
	layout := "2006-01-02T15:04:05"
	t, err := time.Parse(layout, data)
	if err != nil {
		return "", err
	}
	brLayout := "02/01/2006 15:04:05"
	brStr := t.Format(brLayout)
	return brStr, nil
}

func rastreia() {
	for {
		clientRastro, err := rastro.New(*urlRastro, token.Token)
		if err != nil {
			log.Println("erro: rastreia objetos 1: " + err.Error())
		}
		rastros, err := clientRastro.Rastreia(*objeto, token.Token)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		for _, o := range rastros.Objetos {
			dt, err := formataData(o.Eventos[0].DataHora)
			if err != nil {
				log.Println(err.Error())
			}
			log.Printf("%s: %s %s %s %s %s",
				o.CodigoObjeto,
				o.Eventos[0].Codigo,
				o.Eventos[0].Tipo,
				o.Eventos[0].Descricao,
				o.Eventos[0].Unidade.Nome,
				dt,
			)
		}
		time.Sleep(time.Second * time.Duration(*tempo))
	}
}
func main() {
	flag.Parse()
	clientToken, err := token.GetToken(*urlToken, *usuario, *senha)
	if err != nil {
		log.Panic(err.Error())
	}
	token.Token = clientToken.Token
	rastreia()
}
