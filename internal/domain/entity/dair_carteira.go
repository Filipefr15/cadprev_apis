package entity

import (
	"encoding/json"
	"fmt"

	"github.com/shopspring/decimal"
)

type StringOrNumber string

func (s *StringOrNumber) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch val := v.(type) {
	case float64:
		*s = StringOrNumber(fmt.Sprintf("%v", val))
	case string:
		*s = StringOrNumber(val)
	default:
		*s = ""
	}
	return nil
}

type Dair_carteira struct {
	NR_CNPJ_ENTIDADE string          `json:"nr_cnpj_entidade"`
	SG_UF            string          `json:"sg_uf"`
	NO_ENTE          string          `json:"no_ente"`
	DT_MES_BIMESTRE  StringOrNumber  `json:"dt_mes_bimestre"`
	DT_ANO           StringOrNumber  `json:"dt_ano"`
	NO_SEGMENTO      string          `json:"no_segmento"`
	NO_TIPO_ATIVO    string          `json:"no_tipo_ativo"`
	PC_CMN           string          `json:"pc_cmn"`
	ID_ATIVO         string          `json:"id_ativo"`
	NO_FUNDO         string          `json:"no_fundo"`
	QT_RPPS          string          `json:"qt_rpps"`
	VL_ATUAL_ATIVO   DecimalOrString `json:"vl_atual_ativo"`
	VL_TOTAL_ATUAL   DecimalOrString `json:"vl_total_atual"`
	PC_RPPS          string          `json:"pc_rpps"`
	VL_PATRIMONIO    DecimalOrString `json:"vl_patrimonio"`
	PC_PATRIMONIO    string          `json:"pc_patrimonio"`
}

type DecimalOrString struct {
	decimal.Decimal
}

func (d *DecimalOrString) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch val := v.(type) {
	case float64:
		d.Decimal = decimal.NewFromFloat(val)
	case string:
		dec, err := decimal.NewFromString(val)
		if err != nil {
			d.Decimal = decimal.Zero
			return nil
		}
		d.Decimal = dec
	default:
		d.Decimal = decimal.Zero
	}
	return nil
}
