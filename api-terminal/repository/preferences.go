package repository

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const dbPath = "preferences.db"

func initPreferencesDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	createTableQuery := `CREATE TABLE IF NOT EXISTS preferencias (
		id INTEGER PRIMARY KEY CHECK (id = 1),
		tarjeta_a_t INTEGER,
		tarjeta_t_a INTEGER,
		tarjeta_cabezera_a_t INTEGER,
		tarjeta_cabezera_t_a INTEGER,
		relacion_a_t INTEGER,
		relacion_t_a INTEGER,
		relacion_cabezera_a_t INTEGER,
		relacion_cabezera_t_a INTEGER,
		impresora_actual TEXT
	);`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}

	// Asegurar que haya solo un registro con id = 1
	_, err = db.Exec(`INSERT INTO preferencias (id) SELECT 1 WHERE NOT EXISTS (SELECT 1 FROM preferencias WHERE id = 1);`)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func SavePreferencias(prefs map[string]interface{}) error {
	db, err := initPreferencesDB()
	if err != nil {
		return err
	}

	defer db.Close()
	query := `UPDATE preferencias SET
		tarjeta_a_t = ?,
		tarjeta_t_a = ?,
		tarjeta_cabezera_a_t = ?,
		tarjeta_cabezera_t_a = ?,
		relacion_a_t = ?,
		relacion_t_a = ?,
		relacion_cabezera_a_t = ?,
		relacion_cabezera_t_a = ?,
		impresora_actual = ?
	WHERE id = 1;`

	_, err = db.Exec(query,
		prefs["TARJETA-A-T"],
		prefs["TARJETA-T-A"],
		prefs["TARJETA-CABEZERA-A-T"],
		prefs["TARJETA-CABEZERA-T-A"],
		prefs["RELACION-A-T"],
		prefs["RELACION-T-A"],
		prefs["RELACION-CABEZERA-A-T"],
		prefs["RELACION-CABEZERA-T-A"],
		prefs["IMPRESORA-ACTUAL"],
	)

	return err
}

func GetPreferencias() (map[string]interface{}, error) {
	// Initialize the database connection
	db, err := initPreferencesDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Define the SQL query
	query := `SELECT tarjeta_a_t, tarjeta_t_a, tarjeta_cabezera_a_t, tarjeta_cabezera_t_a,
        relacion_a_t, relacion_t_a, relacion_cabezera_a_t, relacion_cabezera_t_a,
        impresora_actual FROM preferencias WHERE id = 1;`

	// Prepare variables to hold the scanned values
	var (
		tarjetaAT, tarjetaTA, tarjetaCabezeraAT, tarjetaCabezeraTA     sql.NullString
		relacionAT, relacionTA, relacionCabezeraAT, relacionCabezeraTA sql.NullString
		impresoraActual                                                sql.NullString
	)

	// Execute the query and scan the results into the variables
	row := db.QueryRow(query)
	err = row.Scan(
		&tarjetaAT,
		&tarjetaTA,
		&tarjetaCabezeraAT,
		&tarjetaCabezeraTA,
		&relacionAT,
		&relacionTA,
		&relacionCabezeraAT,
		&relacionCabezeraTA,
		&impresoraActual,
	)
	if err != nil {
		return nil, err
	}

	// Populate the preferences map with the scanned values, handling NULLs
	prefs := map[string]interface{}{
		"TARJETA-A-T":           nullOrValue(tarjetaAT),
		"TARJETA-T-A":           nullOrValue(tarjetaTA),
		"TARJETA-CABEZERA-A-T":  nullOrValue(tarjetaCabezeraAT),
		"TARJETA-CABEZERA-T-A":  nullOrValue(tarjetaCabezeraTA),
		"RELACION-A-T":          nullOrValue(relacionAT),
		"RELACION-T-A":          nullOrValue(relacionTA),
		"RELACION-CABEZERA-A-T": nullOrValue(relacionCabezeraAT),
		"RELACION-CABEZERA-T-A": nullOrValue(relacionCabezeraTA),
		"IMPRESORA-ACTUAL":      nullOrValue(impresoraActual),
	}

	return prefs, nil
}

// Helper function to check for NULL values
func nullOrValue(v sql.NullString) string {
	if v.Valid {
		return v.String
	}
	return ""
}
