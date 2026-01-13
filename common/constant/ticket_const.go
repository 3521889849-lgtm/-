package constant

// 订单状态枚举值（与工单完全一致，数据库 + 业务层双约束）
const (
	OrderStatusDraft        = "DRAFT"        // 草稿
	OrderStatusPendingPay   = "PENDING_PAYMENT" // 待支付
	OrderStatusPaid         = "PAID"         // 已支付
	OrderStatusUsed         = "USED"         // 已使用
	OrderStatusRefunding    = "REFUNDING"    // 退款中
	OrderStatusRefunded     = "REFUNDED"     // 已退款
	OrderStatusCancelled    = "CANCELLED"    // 已取消
)

// 商家审核状态枚举值
const (
	MerchantAuditStatusInitial  = "INITIAL"  // 待审核
	MerchantAuditStatusApproved = "APPROVED" // 通过
	MerchantAuditStatusRejected = "REJECTED" // 驳回
)

// 支付方式枚举值
const (
	PayTypeWechat = "WECHAT" // 微信
	PayTypeAlipay = "ALIPAY" // 支付宝
)

// 门票状态枚举值
const (
	TicketStatusOnSale   = "ON_SALE"   // 在售
	TicketStatusOffSale  = "OFF_SALE"  // 下架
	TicketStatusStockOut = "STOCK_OUT" // 售罄
)

// 优惠券类型枚举值
const (
	CouponTypeFixed    = "FIXED"    // 满减券
	CouponTypeDiscount = "DISCOUNT" // 折扣券
)

// 优惠券状态枚举值
const (
	CouponStatusValid   = "VALID"   // 有效
	CouponStatusInvalid = "INVALID" // 失效
)

// 用户优惠券使用状态枚举值
const (
	UserCouponStatusUnused  = "UNUSED"  // 未使用
	UserCouponStatusUsed    = "USED"    // 已使用
	UserCouponStatusExpired = "EXPIRED" // 已过期
)

// 支付状态枚举值
const (
	PayStatusSuccess   = "SUCCESS"   // 成功
	PayStatusFail      = "FAIL"      // 失败
	PayStatusRefund    = "REFUND"    // 退款
	PayStatusRefunding = "REFUNDING" // 退款中
)

// 管理员角色枚举值
const (
	AdminRoleSuper    = "SUPER"    // 超级管理员
	AdminRoleOperator = "OPERATOR" // 运营管理员
)

// 操作类型枚举值
const (
	OperTypeMerchantAudit = "MERCHANT_AUDIT" // 商家审核
	OperTypeOrderRefund   = "ORDER_REFUND"   // 订单退款
	OperTypeCouponAdd     = "COUPON_ADD"     // 优惠券添加
	OperTypeOrderCancel   = "ORDER_CANCEL"   // 订单取消
)