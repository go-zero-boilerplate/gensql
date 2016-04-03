package schema

type DialectVisitor interface {
	VisitMysql(*mysql)
	VisitSqlite(*sqlite)
	VisitPostgres(*postgres)
}

func (m *mysql) Accept(visitor DialectVisitor)    { visitor.VisitMysql(m) }
func (s *sqlite) Accept(visitor DialectVisitor)   { visitor.VisitSqlite(s) }
func (p *postgres) Accept(visitor DialectVisitor) { visitor.VisitPostgres(p) }
