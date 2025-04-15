package database

type Database interface {
	Connect() error
	Disconnect() error
}

type DatabaseConnector struct {
	database Database
}

func NewDatabaseConnector(database Database) *DatabaseConnector {
	return &DatabaseConnector{
		database: database,
	}
}

func (c *DatabaseConnector) SetDataBase(database Database) {
	c.database = database
}

func (c *DatabaseConnector) Connect() error {
	return c.database.Connect()
}

func (c *DatabaseConnector) Disconnect() error {
	return c.database.Disconnect()
}
