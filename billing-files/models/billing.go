package models

type Membership struct {
	MembershipID string  `json:"membership_id"`
	Discount     float64 `json:"discount"`
}

type Promotions struct {
	PromotionID   string  `json:"promotion_id"`
	PromotionCode string  `json:"promotion_code"`
	PromotionDiscount      float64 `json:"promotion_discount"`
}
