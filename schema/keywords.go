package schema

// List of vendor-specific keywords
var (
	AUTO_INCREMENT = &AutoIncrementKeyword{}
	PRIMARY_KEY    = &PrimaryKeyKeyword{}
	INDEX          = &IndexKeyword{}
	UNIQUE_INDEX   = &UniqueIndexKeyword{}
	NULL           = &NullKeyword{}
	NOT_NULL       = &NotNullKeyword{}
	DEFAULT        = &DefaultKeyword{}
)

type SchemaKeyword interface {
	Accept(ShemaKeywordVisitor)
}

type ShemaKeywordVisitor interface {
	VisitAutoIncrement(*AutoIncrementKeyword)
	VisitPrimaryKey(*PrimaryKeyKeyword)
	VisitIndex(*IndexKeyword)
	VisitUniqueIndex(*UniqueIndexKeyword)
	VisitNull(*NullKeyword)
	VisitNotNull(*NotNullKeyword)
	VisitDefault(*DefaultKeyword)
}

type AutoIncrementKeyword struct{}
type PrimaryKeyKeyword struct{}
type IndexKeyword struct{}
type UniqueIndexKeyword struct{}
type NullKeyword struct{}
type NotNullKeyword struct{}
type DefaultKeyword struct{}

func (a *AutoIncrementKeyword) Accept(visitor ShemaKeywordVisitor) { visitor.VisitAutoIncrement(a) }
func (p *PrimaryKeyKeyword) Accept(visitor ShemaKeywordVisitor)    { visitor.VisitPrimaryKey(p) }
func (i *IndexKeyword) Accept(visitor ShemaKeywordVisitor)         { visitor.VisitIndex(i) }
func (u *UniqueIndexKeyword) Accept(visitor ShemaKeywordVisitor)   { visitor.VisitUniqueIndex(u) }
func (n *NullKeyword) Accept(visitor ShemaKeywordVisitor)          { visitor.VisitNull(n) }
func (n *NotNullKeyword) Accept(visitor ShemaKeywordVisitor)       { visitor.VisitNotNull(n) }
func (d *DefaultKeyword) Accept(visitor ShemaKeywordVisitor)       { visitor.VisitDefault(d) }
