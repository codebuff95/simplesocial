package databases

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

var GlobalDBM map[string]*DBManager
type DBManager struct{
	Name string
	Database string
	User string
	Password string
	Con *sql.DB
}

func InitGlobalDBM(){
	GlobalDBM = make(map[string]*DBManager)
}
func (dbm *DBManager) Open() error{
	var err error
	dbm.Con,err = sql.Open(dbm.Name, dbm.User+":"+dbm.Password+"@/"+dbm.Database)
	return err
}
