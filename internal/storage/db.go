package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "modernc.org/sqlite"
)

// DB envolve a instância nativa `*sql.DB` estendendo-a com métodos de persistência do LOTGD.
//
// Didática Go: Usamos composição por embutimento (struct embedding) de `*sql.DB` para disponibilizar
// diretamente todos os métodos do pacote padrão `database/sql` (como Exec, QueryRow, BeginTx)
// adicionando métodos utilitários específicos de domínio.
type DB struct {
	*sql.DB
}

// OpenDB abre e inicializa a conexão com o banco de dados SQLite, aplicando pragmas de alta performance e migrações.
//
// Didática Go:
// 1. Usamos o driver CGO-free `modernc.org/sqlite`, permitindo compilação cruzada sem dependências C externas.
// 2. Pragmas de Performance:
//    - `busy_timeout(5000)`: aguarda até 5 segundos para adquirir lock de escrita em ambiente concorrente.
//    - `foreign_keys(1)`: ativa a verificação rigorosa de integridade referencial.
//    - `synchronous(NORMAL)`: otimiza o flush de disco sem comprometer a integridade dos dados em modo WAL.
// 3. Pool de Conexões em BBS: `SetMaxOpenConns(1)` limita a apenas 1 conexão aberta de escrita para evitar concorrência `database is locked`.
// 4. `PRAGMA journal_mode=WAL`: grava transações no Write-Ahead Log, permitindo leituras concorrentes simultâneas.
func OpenDB(dsn string) (*DB, error) {
	pragmaParams := "_pragma=busy_timeout(5000)&_pragma=foreign_keys(1)&_pragma=synchronous(NORMAL)"
	formattedDSN := dsn
	if strings.Contains(dsn, "?") {
		formattedDSN += "&" + pragmaParams
	} else {
		formattedDSN += "?" + pragmaParams
	}

	db, err := sql.Open("sqlite", formattedDSN)
	if err != nil {
		return nil, fmt.Errorf("falha ao abrir banco de dados sqlite: %w", err)
	}

	// Perfil de execução BBS: limita o pool para 1 conexão de escrita para estabilidade no SQLite
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	// Habilita o modo WAL (Write-Ahead Logging) no cabeçalho do arquivo SQLite
	if _, err := db.Exec("PRAGMA journal_mode=WAL;"); err != nil {
		db.Close()
		return nil, fmt.Errorf("falha ao configurar modo WAL no sqlite: %w", err)
	}

	lotgdDB := &DB{DB: db}
	if err := lotgdDB.migrate(); err != nil {
		db.Close()
		return nil, fmt.Errorf("falha ao executar migrações de banco: %w", err)
	}

	return lotgdDB, nil
}

// SavePlayer persiste o estado do herói no banco usando o repositório de jogadores.
func (d *DB) SavePlayer(p Player) error {
	repo := NewPlayerRepository(d)
	return repo.Save(context.Background(), &p)
}

// AuthenticatePlayer valida as credenciais de acesso do usuário.
func (d *DB) AuthenticatePlayer(username, password string) (Player, error) {
	repo := NewPlayerRepository(d)
	p, err := repo.Authenticate(context.Background(), username, password)
	if err != nil {
		return Player{}, err
	}
	return *p, nil
}

// CreatePlayer registra um novo aventureiro no banco de dados.
func (d *DB) CreatePlayer(username, password string) (Player, error) {
	repo := NewPlayerRepository(d)
	p, err := repo.Register(context.Background(), username, password)
	if err != nil {
		return Player{}, err
	}
	return *p, nil
}

// migrate executa as migrações de DDL (Data Definition Language) de forma evolutiva e idempotente.
//
// Didática Go: Lemos o pragma `user_version` do SQLite para rastrear a versão atual do schema.
// Se a versão for inferior a 2, verificamos e adicionamos colunas pendentes (como `potions_count`)
// e atualizamos o `user_version` atomicamente.
func (d *DB) migrate() error {
	var userVersion int
	if err := d.QueryRow("PRAGMA user_version;").Scan(&userVersion); err != nil {
		return fmt.Errorf("falha ao ler user_version do sqlite: %w", err)
	}

	if userVersion < 2 {
		var tableExists int
		err := d.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='players'").Scan(&tableExists)
		if err != nil {
			return fmt.Errorf("falha ao verificar existência da tabela players: %w", err)
		}

		if tableExists > 0 {
			var columnExists int
			err := d.QueryRow("SELECT COUNT(*) FROM pragma_table_info('players') WHERE name='potions_count'").Scan(&columnExists)
			if err != nil {
				return fmt.Errorf("falha ao verificar coluna potions_count: %w", err)
			}
			if columnExists == 0 {
				if _, err := d.Exec("ALTER TABLE players ADD COLUMN potions_count INTEGER NOT NULL DEFAULT 0;"); err != nil {
					return fmt.Errorf("falha ao adicionar coluna potions_count: %w", err)
				}
			}
		}

		if _, err := d.Exec("PRAGMA user_version = 2;"); err != nil {
			return fmt.Errorf("falha ao atualizar user_version: %w", err)
		}
	}

	schema := `
	CREATE TABLE IF NOT EXISTS players (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL COLLATE NOCASE,
		password_hash TEXT NOT NULL,
		level INTEGER NOT NULL DEFAULT 1,
		experience INTEGER NOT NULL DEFAULT 0,
		gold INTEGER NOT NULL DEFAULT 50,
		bank_gold INTEGER NOT NULL DEFAULT 0,
		health INTEGER NOT NULL DEFAULT 20,
		max_health INTEGER NOT NULL DEFAULT 20,
		attack INTEGER NOT NULL DEFAULT 5,
		defense INTEGER NOT NULL DEFAULT 2,
		weapon_id TEXT NOT NULL DEFAULT 'stick',
		armor_id TEXT NOT NULL DEFAULT 'clothes',
		potions_count INTEGER NOT NULL DEFAULT 1,
		forest_fights INTEGER NOT NULL DEFAULT 15,
		dragon_kills INTEGER NOT NULL DEFAULT 0,
		last_login_day TEXT NOT NULL DEFAULT '',
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS village_state (
		day_date TEXT PRIMARY KEY,
		dragon_alive INTEGER NOT NULL DEFAULT 1,
		dragon_hp INTEGER NOT NULL DEFAULT 250,
		dragon_max_hp INTEGER NOT NULL DEFAULT 250,
		dragon_atk INTEGER NOT NULL DEFAULT 35,
		dragon_def INTEGER NOT NULL DEFAULT 20,
		dragon_gold_reward INTEGER NOT NULL DEFAULT 3000,
		slayer_name TEXT NOT NULL DEFAULT ''
	);

	CREATE TABLE IF NOT EXISTS news (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		message TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := d.Exec(schema)
	return err
}
