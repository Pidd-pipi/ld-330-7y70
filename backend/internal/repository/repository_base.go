package repository

import "errors"

var (
	ErrNotFound             = errors.New("not found")
	ErrPatientAlreadyMerged = errors.New("该档案已经合并，不能重复合并")
	ErrMergeIncomplete      = errors.New("关联病历或处方未全部转入，合并已取消")
)
