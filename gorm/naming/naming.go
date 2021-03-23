package naming

import (
	"crypto/sha1"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/jinzhu/inflection"
	"gorm.io/gorm/schema"
)

type TimeNamingStrategy struct {
	TablePrefix   string
	TableSuffix   string
	SingularTable bool
	NameReplacer  schema.Replacer
	NoLowerCase   bool
}

func (s TimeNamingStrategy) TableName(str string) string {
	return inflection.Plural(s.toDBName(str)) + s.TableSuffix
}

// ColumnName convert string to column name
func (s TimeNamingStrategy) ColumnName(table, column string) string {
	return s.toDBName(column)
}

// JoinTableName convert string to join table name
func (s TimeNamingStrategy) JoinTableName(str string) string {
	if !s.NoLowerCase && strings.ToLower(str) == str {
		return s.TablePrefix + str
	}

	if s.SingularTable {
		return s.TablePrefix + s.toDBName(str)
	}
	return s.TablePrefix + inflection.Plural(s.toDBName(str))
}

// RelatioshipFKName generate fk name for relation
func (s TimeNamingStrategy) RelationshipFKName(rel schema.Relationship) string {
	return s.formatName("fk", rel.Schema.Table, s.toDBName(rel.Name))
}

// CheckerName generate checker name
func (s TimeNamingStrategy) CheckerName(table, column string) string {
	return s.formatName("chk", table, column)
}

// IndexName generate index name
func (s TimeNamingStrategy) IndexName(table, column string) string {
	return s.formatName("idx", table, s.toDBName(column))
}

func (s TimeNamingStrategy) formatName(prefix, table, name string) string {
	formatedName := strings.Replace(fmt.Sprintf("%v_%v_%v", prefix, table, name), ".", "_", -1)

	if utf8.RuneCountInString(formatedName) > 64 {
		h := sha1.New()
		h.Write([]byte(formatedName))
		bs := h.Sum(nil)

		formatedName = fmt.Sprintf("%v%v%v", prefix, table, name)[0:56] + string(bs)[:8]
	}
	return formatedName
}

var (
	// https://github.com/golang/lint/blob/master/lint.go#L770
	commonInitialisms         = []string{"API", "ASCII", "CPU", "CSS", "DNS", "EOF", "GUID", "HTML", "HTTP", "HTTPS", "ID", "IP", "JSON", "LHS", "QPS", "RAM", "RHS", "RPC", "SLA", "SMTP", "SSH", "TLS", "TTL", "UID", "UI", "UUID", "URI", "URL", "UTF8", "VM", "XML", "XSRF", "XSS"}
	commonInitialismsReplacer *strings.Replacer
)

func init() {
	commonInitialismsForReplacer := make([]string, 0, len(commonInitialisms))
	for _, initialism := range commonInitialisms {
		commonInitialismsForReplacer = append(commonInitialismsForReplacer, initialism, strings.Title(strings.ToLower(initialism)))
	}
	commonInitialismsReplacer = strings.NewReplacer(commonInitialismsForReplacer...)
}

func (s TimeNamingStrategy) toDBName(name string) string {
	if name == "" {
		return ""
	}

	if s.NameReplacer != nil {
		name = s.NameReplacer.Replace(name)
	}

	if s.NoLowerCase {
		return name
	}

	var (
		value                          = commonInitialismsReplacer.Replace(name)
		buf                            strings.Builder
		lastCase, nextCase, nextNumber bool // upper case == true
		curCase                        = value[0] <= 'Z' && value[0] >= 'A'
	)

	for i, v := range value[:len(value)-1] {
		nextCase = value[i+1] <= 'Z' && value[i+1] >= 'A'
		nextNumber = value[i+1] >= '0' && value[i+1] <= '9'

		if curCase {
			if lastCase && (nextCase || nextNumber) {
				buf.WriteRune(v + 32)
			} else {
				if i > 0 && value[i-1] != '_' && value[i+1] != '_' {
					buf.WriteByte('_')
				}
				buf.WriteRune(v + 32)
			}
		} else {
			buf.WriteRune(v)
		}

		lastCase = curCase
		curCase = nextCase
	}

	if curCase {
		if !lastCase && len(value) > 1 {
			buf.WriteByte('_')
		}
		buf.WriteByte(value[len(value)-1] + 32)
	} else {
		buf.WriteByte(value[len(value)-1])
	}
	ret := buf.String()
	return ret
}
