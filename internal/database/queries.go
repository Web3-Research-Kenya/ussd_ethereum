package database

func (s *service) CreateTable() {
	schema := `
	CREATE TABLE IF NOT EXISTS wallets (
		id integer primary key autoincrement,
		phone_number text not null unique,
		public_key text not null unique,
		pin text not null,
		keystore_path text not null
	);`
	s.db.MustExec(schema)
}

func (s *service) InsertWallet(phoneNumber, pin, publicKey, keystorePath string) error {
	stmt := `
	insert into wallets (
		phone_number,
		pin,
		public_key,
		keystore_path
	) values (
		?, ?, ?,?
	);
	`
	_, err := s.db.Exec(stmt, phoneNumber, pin, publicKey, keystorePath)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) SelectWalletByPhone(phoneNumber string) (*WalletRecord, error) {
	stmt := `
	select *
	from wallets
	where phone_number = ?;
	`

	record := WalletRecord{}

	err := s.db.Get(&record, stmt, phoneNumber)
	if err != nil {
		return nil, err
	}

	return &record, nil
}

func (s *service) UpdateKeystorePathByID(path string, id uint64) {
	stmt := `
	update wallets 
	set keystore_path = ?
	where id = ?;
	`

	s.db.MustExec(stmt, path, id)
}
