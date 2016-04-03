package schema

var (
	INTEGER   = &IntegerFieldType{}
	VARCHAR   = &VarcharFieldType{}
	BOOLEAN   = &BooleanFieldType{}
	REAL      = &RealFieldType{}
	BLOB      = &BlobFieldType{}
	DATETIME  = &DateTimeFieldType{}
	TIMESTAMP = &TimeStampFieldType{}
)

type FieldType interface {
	Accept(FieldTypeVisitor)
}

type FieldTypeVisitor interface {
	VisitInteger(*IntegerFieldType)
	VisitVarchar(*VarcharFieldType)
	VisitBoolean(*BooleanFieldType)
	VisitReal(*RealFieldType)
	VisitBlob(*BlobFieldType)
	VisitDateTime(*DateTimeFieldType)
	VisitTimeStamp(*TimeStampFieldType)
}

type IntegerFieldType struct{}
type VarcharFieldType struct{}
type BooleanFieldType struct{}
type RealFieldType struct{}
type BlobFieldType struct{}
type DateTimeFieldType struct{}
type TimeStampFieldType struct{}

func (i *IntegerFieldType) Accept(visitor FieldTypeVisitor)   { visitor.VisitInteger(i) }
func (v *VarcharFieldType) Accept(visitor FieldTypeVisitor)   { visitor.VisitVarchar(v) }
func (b *BooleanFieldType) Accept(visitor FieldTypeVisitor)   { visitor.VisitBoolean(b) }
func (r *RealFieldType) Accept(visitor FieldTypeVisitor)      { visitor.VisitReal(r) }
func (b *BlobFieldType) Accept(visitor FieldTypeVisitor)      { visitor.VisitBlob(b) }
func (d *DateTimeFieldType) Accept(visitor FieldTypeVisitor)  { visitor.VisitDateTime(d) }
func (t *TimeStampFieldType) Accept(visitor FieldTypeVisitor) { visitor.VisitTimeStamp(t) }
