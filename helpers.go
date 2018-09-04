package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/kennygrant/sanitize"
)

//OBS: por enquanto estou trabalhando com "input" global, mas o ideal é para cada passo fazer um closure em torno de uma variável para input e dar input como valor daquela variável

func marcacao(mat, i string) string {
	input := normalize(i)

	//verificar se há marcação em andamento
	marca, ok := feriasMAP[mat]
	if !ok {
		switch input {
		case "sim":
			//buscar dados de dias de férias disponíveis no sistema (agora buscar em dadosFerias.go; mudar isso para )
			feriasDisponiveis := feriasMichael

			//criar marcacao on feriasMAP com matricula como chave e incluir quantidade de dias disponíveis para férias
			feriasMAP[mat] = &passosFerias{diasDisponiveis: feriasDisponiveis.dias}

			//verificar se há dias para gozo de férias
			if feriasDisponiveis.dias == 0 {
				return "Não há dias férias disponiveis para gozo."
			}
			//se houver dias retorna para o usuário quantos dias e solicita data de início
			return fmt.Sprintf("Existem %d dias de férias para gozo. Por favor digite a data de início no seguinte formato: DD/MM/AAAA", feriasDisponiveis.dias)

		case "nao":
			return "Obrigado por utilizar o Bot para marcação de férias."

		default:
			return `Não entendi sua solicitação, digite "sim" ou "não"`
		}
	}

	//hipóteses em que já há uma marcação em curso
	if ok {
		//verificar se é formato de data, deixar no escopo de fora porque será usado outras vezes
		//layout referência usado por Go
		layout := "02/01/2006"
		//hipótese em que ainda não há data de início registrada
		if marca.regDataInicio == false {
			dataInicio, err := time.Parse(layout, input)
			if err != nil {
				//verificar os tipos de erros que tem; por enquanto só retorna para usuário
				return "Por favor verifique a data de início."
			}
			//verificar se a data solicitada é posterior a hoje
			todayT := time.Now()
			todayF, err := time.Parse(layout, todayT.Format(layout))
			if err != nil {
				return "Houve um problema com sua solicitação. Por gentileza forneça novamente a data de início."
			}
			//subtrair de data inicio o equivalente a today e dar um truncate para ver horas; aqui a difernça está em horas (float64), mais adiante uso dias (int)
			diff := dataInicio.Sub(todayF).Hours()
			if diff < 24 {
				//mandar ao usuario que tem que ser mais que um dia
				return "A data de início tem que ser superior à fornecida."
			}
			//marcar no fériasMAP a data em dataInicio e true para regDataInicio
			marca.dataDeInicio = input
			marca.regDataInicio = true
			return "Por favor indique a data de fim. Lembrando que os períodos de férias devem ser marcados em intervalos de, no mínimo, 10 dias."
		}
		//hipótese em que há data de início mas não há data de fim
		if marca.regDataFim == false {
			//transformar a data de início no feriasMAP em Time
			di, err := time.Parse(layout, marca.dataDeInicio)
			if err != nil {
				return "Houve um problema, por favor insira novamente a data de fim."
			}
			//transformar a data de fim vinda no input em Time
			df, err := time.Parse(layout, input)
			if err != nil {
				return "Houve um problema, por favor insira novamente a data de fim no formado DD/MM/AAAA."
			}
			//verificar se há uma diferença de pelo menos 10 dias entre elas; conertido de horas (float64) para dias (int)
			diasPretendidos := int(df.Sub(di).Hours())/24 + 1
			if diasPretendidos < 10 {
				return "Deve haver pelo menos 10 dias de diferença entra a data de início e fim. Por favor, indique outra data de fim."
			}
			// verificar se tem dias de férias suficientes
			if diasPretendidos > marca.diasDisponiveis {
				return fmt.Sprintf("O número de dias de férias pretendido (%d) excede o número de dias de férias disponíveis (%d)", diasPretendidos, marca.diasDisponiveis)
			}

			//TO DO: falta verificar se o período solicitado é múltiplo de 10

			//caso tenha dias de férias suficientes
			marca.dataDeFim = input
			marca.regDataFim = true
			return fmt.Sprintf(`Podemos confirmar a solicitação de %d dias de férias, entre os dias %s e %s? Digite "sim" ou "não".`, diasPretendidos, marca.dataDeInicio, marca.dataDeFim)
		}
		//hipótese em que já ha data de inicio e de fim e foi solicitada confirmação para o cliente
		if marca.sistamaOk == false {
			switch input {
			case "sim":
				//fazer requisição para o sistema de marcação de férias; em caso de sucesso deletar o registro da marcação no feriasMAP e retornar para o usuário.
				delete(feriasMAP, mat)
				return "Solicitação realizada com sucesso. Um email será enviado ao Procurador-Chefe para autorização. Deseja solicitar mais um período de férias?"
			case "nao":
				//deleta os passos de marcação de férias e manda para o usuário.
				delete(feriasMAP, mat)
				return "Essa solicitação de férias foi cancelada. Deseja iniciar outra solicitação de férias?"
			default:
				return `Não entendi sua solicitação, digite "sim" ou "não"`
			}
		}
	}
	//só vou printar marca, que será usada nos próximos passos
	fmt.Println(marca)
	return "próximo passo"
}

//func normalize recebe a string de input, retira os acentos e coloca em caixa baixa
func normalize(s string) string {
	str := strings.ToLower(sanitize.Accents(s))
	return str
}
