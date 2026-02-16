package pdf

import "context"

// CertificateData contém os dados para preencher o template do certificado
type CertificateData struct {
	RecipientName   string
	EventName       string
	EventDate       string
	Workload        string
	DirectorName    string
	CoordinatorName string
	CertificateDate string
}

// CertificateGenerator define a interface para geração de certificados em PDF
type CertificateGenerator interface {
	// Generate gera um certificado PDF a partir dos dados fornecidos
	Generate(ctx context.Context, data CertificateData) ([]byte, error)
}
