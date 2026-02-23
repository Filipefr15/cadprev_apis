package api

import (
	"database/sql"

	"github.com/Filipefr15/cadprev_apis/internal/domain/entity"
	_ "modernc.org/sqlite"
)

func OpenSQLite(dbPath string) (*sql.DB, error) {
	return sql.Open("sqlite", dbPath)
}

func QueryDairCarteira(db *sql.DB, params map[string]interface{}) ([]entity.Dair_carteira, error) {
	query := `SELECT * FROM dair_carteira WHERE dt_ano = ? AND (nr_cnpj_entidade = ? OR no_ente = ?)`
	args := []interface{}{params["dt_ano"], params["nr_cnpj_entidade"], params["no_ente"]}

	if params["dt_mes_bimestre"] != "" {
		query += " AND dt_mes_bimestre = ?"
		args = append(args, params["dt_mes_bimestre"])
	}
	if params["no_segmento"] != "" {
		query += " AND no_segmento = ?"
		args = append(args, params["no_segmento"])
	}
	if params["sg_uf"] != "" {
		query += " AND sg_uf = ?"
		args = append(args, params["sg_uf"])
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []entity.Dair_carteira
	for rows.Next() {
		var record entity.Dair_carteira
		if err := rows.Scan(
			&record.NR_CNPJ_ENTIDADE,
			&record.SG_UF,
			&record.NO_ENTE,
			&record.DT_MES_BIMESTRE,
			&record.DT_ANO,
			&record.NO_SEGMENTO,
			&record.NO_TIPO_ATIVO,
			&record.PC_CMN,
			&record.ID_ATIVO,
			&record.NO_FUNDO,
			&record.QT_RPPS,
			&record.VL_ATUAL_ATIVO,
			&record.VL_TOTAL_ATUAL,
			&record.PC_RPPS,
			&record.VL_PATRIMONIO,
			&record.PC_PATRIMONIO,
		); err != nil {
			return nil, err
		}
		results = append(results, record)
	}

	return results, nil
}
