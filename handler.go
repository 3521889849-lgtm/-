package main

import (
	"context"
	audit "example_shop/kitex_gen/audit"
	coupon "example_shop/kitex_gen/coupon"
)

// CouponServiceImpl implements the last service interface defined in the IDL.
type CouponServiceImpl struct{}

// Test implements the CouponServiceImpl interface.
func (s *CouponServiceImpl) Test(ctx context.Context, req *coupon.EmptyReq) (resp *coupon.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// ApplyAudit implements the AuditServiceImpl interface.
func (s *AuditServiceImpl) ApplyAudit(ctx context.Context, req *audit.ApplyAuditReq) (resp *audit.ApplyAuditResp, err error) {
	// TODO: Your code here...
	return
}

// FetchManualTasks implements the AuditServiceImpl interface.
func (s *AuditServiceImpl) FetchManualTasks(ctx context.Context, req *audit.FetchManualTasksReq) (resp *audit.FetchManualTasksResp, err error) {
	// TODO: Your code here...
	return
}

// ProcessManualAudit implements the AuditServiceImpl interface.
func (s *AuditServiceImpl) ProcessManualAudit(ctx context.Context, req *audit.ProcessManualAuditReq) (resp *audit.ProcessManualAuditResp, err error) {
	// TODO: Your code here...
	return
}

// GetAuditRecord implements the AuditServiceImpl interface.
func (s *AuditServiceImpl) GetAuditRecord(ctx context.Context, req *audit.GetAuditRecordReq) (resp *audit.GetAuditRecordResp, err error) {
	// TODO: Your code here...
	return
}
