package pdf

import (
	"context"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/line"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/consts/orientation"
	"github.com/johnfercher/maroto/v2/pkg/consts/pagesize"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

// Cores do certificado
var (
	primaryColor   = &props.Color{Red: 26, Green: 26, Blue: 46}    // #1a1a2e
	goldColor      = &props.Color{Red: 201, Green: 162, Blue: 39}  // #c9a227
	darkBlueColor  = &props.Color{Red: 15, Green: 52, Blue: 96}    // #0f3460
	grayColor      = &props.Color{Red: 102, Green: 102, Blue: 102} // #666666
	lightGrayColor = &props.Color{Red: 153, Green: 153, Blue: 153} // #999999
)

// MarotoGenerator implementa CertificateGenerator usando maroto v2
type MarotoGenerator struct{}

// NewMarotoGenerator cria uma nova instância do gerador de certificados
func NewMarotoGenerator() *MarotoGenerator {
	return &MarotoGenerator{}
}

// Generate gera um certificado PDF usando maroto v2
func (g *MarotoGenerator) Generate(_ context.Context, data CertificateData) ([]byte, error) {
	cfg := config.NewBuilder().
		WithPageSize(pagesize.A4).
		WithOrientation(orientation.Horizontal).
		WithLeftMargin(15).
		WithRightMargin(15).
		WithTopMargin(15).
		WithBottomMargin(15).
		Build()

	m := maroto.New(cfg)

	m.AddRows(g.buildHeader(data)...)
	m.AddRows(g.buildContent(data)...)
	m.AddRows(g.buildSignatures(data)...)
	m.AddRows(g.buildFooter()...)

	doc, err := m.Generate()
	if err != nil {
		return nil, err
	}

	return doc.GetBytes(), nil
}

func (g *MarotoGenerator) buildHeader(_ CertificateData) []core.Row {
	return []core.Row{
		// Espaço superior
		row.New(10),

		// Linha decorativa superior
		row.New(2).Add(
			col.New(2),
			line.NewCol(8, props.Line{
				Color:     goldColor,
				Thickness: 1,
			}),
			col.New(2),
		),

		// Espaço
		row.New(8),

		// Título "Certificado"
		row.New(20).Add(
			col.New(12).Add(
				text.New("CERTIFICADO", props.Text{
					Size:  36,
					Style: fontstyle.Bold,
					Align: align.Center,
					Color: primaryColor,
				}),
			),
		),

		// Linha decorativa inferior ao título
		row.New(2).Add(
			col.New(3),
			line.NewCol(6, props.Line{
				Color:     goldColor,
				Thickness: 0.5,
			}),
			col.New(3),
		),

		// Espaço
		row.New(8),
	}
}

func (g *MarotoGenerator) buildContent(data CertificateData) []core.Row {
	return []core.Row{
		// "Certificamos que"
		row.New(8).Add(
			col.New(12).Add(
				text.New("Certificamos que", props.Text{
					Size:  12,
					Align: align.Center,
					Color: grayColor,
				}),
			),
		),

		// Espaço
		row.New(5),

		// Nome do participante
		row.New(15).Add(
			col.New(12).Add(
				text.New(data.RecipientName, props.Text{
					Size:  24,
					Style: fontstyle.Bold,
					Align: align.Center,
					Color: primaryColor,
				}),
			),
		),

		// Espaço
		row.New(5),

		// Descrição
		row.New(8).Add(
			col.New(1),
			col.New(10).Add(
				text.New("participou do evento", props.Text{
					Size:  11,
					Align: align.Center,
					Color: grayColor,
				}),
			),
			col.New(1),
		),

		// Nome do evento
		row.New(10).Add(
			col.New(12).Add(
				text.New(data.EventName, props.Text{
					Size:  14,
					Style: fontstyle.Bold,
					Align: align.Center,
					Color: darkBlueColor,
				}),
			),
		),

		// Data e carga horária
		row.New(8).Add(
			col.New(1),
			col.New(10).Add(
				text.New("realizado em "+data.EventDate+", com carga horária total de "+data.Workload+".", props.Text{
					Size:  11,
					Align: align.Center,
					Color: grayColor,
				}),
			),
			col.New(1),
		),

		// Espaço
		row.New(10),

		// Detalhes em colunas - Labels
		row.New(6).Add(
			col.New(4).Add(
				text.New("DATA DO EVENTO", props.Text{
					Size:  8,
					Align: align.Center,
					Color: lightGrayColor,
				}),
			),
			col.New(4).Add(
				text.New("CARGA HORÁRIA", props.Text{
					Size:  8,
					Align: align.Center,
					Color: lightGrayColor,
				}),
			),
			col.New(4).Add(
				text.New("DATA DE EMISSÃO", props.Text{
					Size:  8,
					Align: align.Center,
					Color: lightGrayColor,
				}),
			),
		),

		// Detalhes em colunas - Valores
		row.New(8).Add(
			col.New(4).Add(
				text.New(data.EventDate, props.Text{
					Size:  10,
					Style: fontstyle.Bold,
					Align: align.Center,
					Color: primaryColor,
				}),
			),
			col.New(4).Add(
				text.New(data.Workload, props.Text{
					Size:  10,
					Style: fontstyle.Bold,
					Align: align.Center,
					Color: primaryColor,
				}),
			),
			col.New(4).Add(
				text.New(data.CertificateDate, props.Text{
					Size:  10,
					Style: fontstyle.Bold,
					Align: align.Center,
					Color: primaryColor,
				}),
			),
		),

		// Espaço antes das assinaturas
		row.New(10),
	}
}

func (g *MarotoGenerator) buildSignatures(data CertificateData) []core.Row {
	return []core.Row{
		// Linhas de assinatura
		row.New(1).Add(
			col.New(2),
			line.NewCol(3, props.Line{
				Color:     primaryColor,
				Thickness: 0.5,
			}),
			col.New(2),
			line.NewCol(3, props.Line{
				Color:     primaryColor,
				Thickness: 0.5,
			}),
			col.New(2),
		),

		// Espaço
		row.New(2),

		// Nomes dos signatários
		row.New(6).Add(
			col.New(2),
			col.New(3).Add(
				text.New(data.DirectorName, props.Text{
					Size:  9,
					Style: fontstyle.Bold,
					Align: align.Center,
					Color: primaryColor,
				}),
			),
			col.New(2),
			col.New(3).Add(
				text.New(data.CoordinatorName, props.Text{
					Size:  9,
					Style: fontstyle.Bold,
					Align: align.Center,
					Color: primaryColor,
				}),
			),
			col.New(2),
		),

		// Títulos
		row.New(5).Add(
			col.New(2),
			col.New(3).Add(
				text.New("Diretor(a)", props.Text{
					Size:  8,
					Align: align.Center,
					Color: grayColor,
				}),
			),
			col.New(2),
			col.New(3).Add(
				text.New("Coordenador(a)", props.Text{
					Size:  8,
					Align: align.Center,
					Color: grayColor,
				}),
			),
			col.New(2),
		),
	}
}

func (g *MarotoGenerator) buildFooter() []core.Row {
	return []core.Row{
		// Espaço
		row.New(8),

		// Rodapé
		row.New(5).Add(
			col.New(12).Add(
				text.New("Certificado gerado automaticamente pelo sistema Checkin Gate", props.Text{
					Size:  8,
					Align: align.Center,
					Color: lightGrayColor,
				}),
			),
		),
	}
}
