package ql

import (
	"errors"
)

var (
	ErrTableNameRequire = errors.New("[ql] tableName requires")
	ErrUpdateMissWhere  = errors.New("[ql] where express requires with update")
	ErrColumnsRequire   = errors.New("[ql] columns requires")
	ErrDeleteMissWhere  = errors.New("[ql] delete sql miss where")
	ErrExecerNotSet     = errors.New("[ql] execer not set")
	ErrQueryerNotSet    = errors.New("[ql] queryer not set")
)
