// Code generated by qtc from "alter_table.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line alter_table.qtpl:1
package migrate

//line alter_table.qtpl:1
import (
	"github.com/Alexandrhub/cli-orm-gen/infrastructure/db/scanner"

//line alter_table.qtpl:5

	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line alter_table.qtpl:5
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line alter_table.qtpl:5
func StreamAlterTable(qw422016 *qt422016.Writer, field scanner.Field) {
//line alter_table.qtpl:5
	qw422016.N().S(`
alter table `)
//line alter_table.qtpl:6
	qw422016.E().S(field.Table.Name)
//line alter_table.qtpl:6
	qw422016.N().S(`
	add `)
//line alter_table.qtpl:7
	qw422016.E().S(field.Name)
//line alter_table.qtpl:7
	qw422016.N().S(` `)
//line alter_table.qtpl:7
	qw422016.E().S(field.Type)
//line alter_table.qtpl:7
	qw422016.N().S(` `)
//line alter_table.qtpl:7
	qw422016.E().S(field.Default)
//line alter_table.qtpl:7
	qw422016.N().S(`;

`)
//line alter_table.qtpl:9
	if field.Constraint.Index {
//line alter_table.qtpl:9
		qw422016.N().S(`
    create `)
//line alter_table.qtpl:10
		if field.Constraint.Unique {
//line alter_table.qtpl:10
			qw422016.N().S(`unique `)
//line alter_table.qtpl:10
		}
//line alter_table.qtpl:10
		qw422016.N().S(`index `)
//line alter_table.qtpl:10
		qw422016.E().S(field.Table.Name)
//line alter_table.qtpl:10
		qw422016.N().S(`_`)
//line alter_table.qtpl:10
		qw422016.E().S(field.Constraint.Field.Name)
//line alter_table.qtpl:10
		qw422016.N().S(`_idx
     on `)
//line alter_table.qtpl:11
		qw422016.E().S(field.Table.Name)
//line alter_table.qtpl:11
		qw422016.N().S(` (`)
//line alter_table.qtpl:11
		qw422016.E().S(field.Constraint.Field.Name)
//line alter_table.qtpl:11
		qw422016.N().S(`);`)
//line alter_table.qtpl:11
	}
//line alter_table.qtpl:11
}

//line alter_table.qtpl:11
func WriteAlterTable(qq422016 qtio422016.Writer, field scanner.Field) {
//line alter_table.qtpl:11
	qw422016 := qt422016.AcquireWriter(qq422016)
//line alter_table.qtpl:11
	StreamAlterTable(qw422016, field)
//line alter_table.qtpl:11
	qt422016.ReleaseWriter(qw422016)
//line alter_table.qtpl:11
}

//line alter_table.qtpl:11
func AlterTable(field scanner.Field) string {
//line alter_table.qtpl:11
	qb422016 := qt422016.AcquireByteBuffer()
//line alter_table.qtpl:11
	WriteAlterTable(qb422016, field)
//line alter_table.qtpl:11
	qs422016 := string(qb422016.B)
//line alter_table.qtpl:11
	qt422016.ReleaseByteBuffer(qb422016)
//line alter_table.qtpl:11
	return qs422016
//line alter_table.qtpl:11
}
