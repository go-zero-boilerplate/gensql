package schema

var (
	INTEGER   = &IntegerFieldType{}
	BIGINT    = &BigIntFieldType{}
	VARCHAR   = &VarcharFieldType{}
	BOOLEAN   = &BooleanFieldType{}
	REAL      = &RealFieldType{}
	BLOB      = &BlobFieldType{}
	DATETIME  = &DateTimeFieldType{}
	TIMESTAMP = &TimeStampFieldType{}
	CREATED   = &CreatedFieldType{}
	UPDATED   = &UpdatedFieldType{}
)

type FieldType interface {
	Accept(FieldTypeVisitor)
}

type FieldTypeVisitor interface {
	VisitInteger(*IntegerFieldType)
	VisitBigInt(*BigIntFieldType)
	VisitVarchar(*VarcharFieldType)
	VisitBoolean(*BooleanFieldType)
	VisitReal(*RealFieldType)
	VisitBlob(*BlobFieldType)
	VisitDateTime(*DateTimeFieldType)
	VisitTimeStamp(*TimeStampFieldType)
	VisitCreated(*CreatedFieldType)
	VisitUpdated(*UpdatedFieldType)
}

type IntegerFieldType struct{}
type BigIntFieldType struct{}
type VarcharFieldType struct{}
type BooleanFieldType struct{}
type RealFieldType struct{}
type BlobFieldType struct{}
type DateTimeFieldType struct{}
type TimeStampFieldType struct{}
type CreatedFieldType struct{}
type UpdatedFieldType struct{}

func (i *IntegerFieldType) Accept(visitor FieldTypeVisitor)   { visitor.VisitInteger(i) }
func (b *BigIntFieldType) Accept(visitor FieldTypeVisitor)    { visitor.VisitBigInt(b) }
func (v *VarcharFieldType) Accept(visitor FieldTypeVisitor)   { visitor.VisitVarchar(v) }
func (b *BooleanFieldType) Accept(visitor FieldTypeVisitor)   { visitor.VisitBoolean(b) }
func (r *RealFieldType) Accept(visitor FieldTypeVisitor)      { visitor.VisitReal(r) }
func (b *BlobFieldType) Accept(visitor FieldTypeVisitor)      { visitor.VisitBlob(b) }
func (d *DateTimeFieldType) Accept(visitor FieldTypeVisitor)  { visitor.VisitDateTime(d) }
func (t *TimeStampFieldType) Accept(visitor FieldTypeVisitor) { visitor.VisitTimeStamp(t) }
func (c *CreatedFieldType) Accept(visitor FieldTypeVisitor)   { visitor.VisitCreated(c) }
func (u *UpdatedFieldType) Accept(visitor FieldTypeVisitor)   { visitor.VisitUpdated(u) }
