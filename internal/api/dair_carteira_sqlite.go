package api

import (
	"database/sql"

	"github.com/Filipefr15/cadprev_apis/internal/domain/entity"
)

func CreateDairCarteiraTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS dair_carteira (
		nr_cnpj_entidade TEXT,
		sg_uf TEXT,
		no_ente TEXT,
		dt_mes_bimestre TEXT,
		dt_ano TEXT,
		no_segmento TEXT,
		no_tipo_ativo TEXT,
		pc_cmn TEXT,
		id_ativo TEXT,
		no_fundo TEXT,
		qt_rpps TEXT,
		vl_atual_ativo TEXT,
		vl_total_atual TEXT,
		pc_rpps TEXT,
		vl_patrimonio TEXT,
		pc_patrimonio TEXT
	);`
	_, err := db.Exec(query)
	return err
}

func InsertDairCarteira(db *sql.DB, item entity.Dair_carteira) error {
	query := `INSERT INTO dair_carteira (
		nr_cnpj_entidade, sg_uf, no_ente, dt_mes_bimestre, dt_ano, no_segmento, no_tipo_ativo, pc_cmn, id_ativo, no_fundo, qt_rpps, vl_atual_ativo, vl_total_atual, pc_rpps, vl_patrimonio, pc_patrimonio
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	_, err := db.Exec(query,
		item.NR_CNPJ_ENTIDADE,
		item.SG_UF,
		item.NO_ENTE,
		string(item.DT_MES_BIMESTRE),
		string(item.DT_ANO),
		item.NO_SEGMENTO,
		item.NO_TIPO_ATIVO,
		item.PC_CMN,
		item.ID_ATIVO,
		item.NO_FUNDO,
		item.QT_RPPS,
		item.VL_ATUAL_ATIVO.String(),
		item.VL_TOTAL_ATUAL.String(),
		item.PC_RPPS,
		item.VL_PATRIMONIO.String(),
		item.PC_PATRIMONIO,
	)
	return err
}
