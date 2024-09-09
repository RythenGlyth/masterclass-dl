package main

type CSRFResponse struct {
	Param string `json:"param"`
	Token string `json:"token"`
}

type ProfileResponse struct {
	UUID              string      `json:"uuid"`
	Slug              string      `json:"slug"`
	DisplayName       string      `json:"display_name"`
	Tagline           string      `json:"tagline"`
	City              interface{} `json:"city"`
	Country           interface{} `json:"country"`
	PersonalURL       interface{} `json:"personal_url"`
	PersonalURLType   interface{} `json:"personal_url_type"`
	Profession        interface{} `json:"profession"`
	Location          string      `json:"location"`
	DefaultProfile    bool        `json:"default_profile"`
	MatureContent     interface{} `json:"mature_content_enabled"`
	AvatarColor       string      `json:"avatar_color"`
	OnboardedAt       interface{} `json:"onboarded_at"`
	CareerGoal        interface{} `json:"career_goal"`
	ChatEnabled       bool        `json:"chat_enabled"`
	ChatOptIn         bool        `json:"chat_opt_in"`
	ProfileImage      interface{} `json:"profile_image"`
	Skills            []string    `json:"skills"`
	PrimaryCategoryID interface{} `json:"primary_category_id"`
	ID                int         `json:"id"`
	User              struct {
		Email                 string      `json:"email"`
		ProfileID             int         `json:"profile_id"`
		ProfilePhotoURL       interface{} `json:"profile_photo_url"`
		Locale                interface{} `json:"locale"`
		SubtitleEnabled       bool        `json:"subtitle_enabled"`
		SubtitleLocale        interface{} `json:"subtitle_locale"`
		MatureContent         interface{} `json:"mature_content_enabled"`
		RequiresConsent       bool        `json:"requires_consent"`
		NoEmailUser           bool        `json:"no_email_user"`
		NeedsSetup            bool        `json:"needs_setup"`
		Username              interface{} `json:"username"`
		FirstName             string      `json:"first_name"`
		LastName              string      `json:"last_name"`
		CanUpgrade            bool        `json:"can_upgrade"`
		ActiveAnnualPass      bool        `json:"active_annual_pass"`
		BogoGiftToken         interface{} `json:"bogo_gift_token"`
		BogoGiftExpiration    interface{} `json:"bogo_gift_expiration"`
		BogoGiftExpirationRaw interface{} `json:"bogo_gift_expiration_raw"`
		BogoEligible          bool        `json:"bogo_eligible"`
		GdprOptIn             bool        `json:"gdpr_opt_in"`
		Slug                  string      `json:"slug"`
		AvailableAuthTypes    []struct {
			Type        string `json:"type"`
			LastLoginAt string `json:"last_login_at"`
		} `json:"available_auth_types"`
		CurrentAuthType                 string        `json:"current_auth_type"`
		OrganizationsAdministered       []interface{} `json:"organizations_administered"`
		GatedFeature                    bool          `json:"gated_feature"`
		EmailToken                      string        `json:"email_token"`
		EnrolledCourses                 []string      `json:"enrolled_courses"`
		HasHadSubscriptions             bool          `json:"has_had_subscriptions"`
		HasCreditCard                   bool          `json:"has_credit_card"`
		BraintreePaypalEmailOfLatestSub string        `json:"braintree_paypal_email_of_latest_sub"`
		HasBraintreeCCSub               bool          `json:"has_braintree_cc_sub"`
		SponsoredSubscriptionInfo       interface{}   `json:"sponsored_subscription_info"`
		HasActiveAfterpaySubscription   bool          `json:"has_active_afterpay_subscription"`
		ID                              int           `json:"id"`
		Entitlement                     struct {
			ID int `json:"id"`
		} `json:"entitlement"`
		OrganizationsMembered []interface{} `json:"organizations_membered"`
		EnterpriseSeatToken   interface{}   `json:"enterprise_seat_token"`
	} `json:"user"`
}

type CartDataResponse struct {
	Email        string `json:"email"`
	Subscription struct {
		// Status                    string `json:"status"`
		// ProviderType              string `json:"provider_type"`
		// CanceledAt                string `json:"canceled_at"`
		// CreatedAt                 string `json:"created_at"`
		// RemainingSubscriptionDays int    `json:"remaining_subscription_days"`
		// CancelAtPeriodEnd         bool   `json:"cancel_at_period_end"`
		// PaymentGatewayUserProduct string `json:"payment_gateway_user_product"`
		// ExpirationDate            string `json:"expiration_date"`
		// OriginatorType            string `json:"originator_type"`
		// IsMonthlyPass             bool   `json:"is_monthly_pass"`
		// CanceledOnPaypalDashboard string `json:"canceled_on_paypal_dashboard"`
		// PurchaseCountry           string `json:"purchase_country"`
		// AccessType                string `json:"access_type"`
		// CurrentSubscriptionCycle  int    `json:"current_subscription_cycle"`
		// CanUpgradeToConsumer      bool   `json:"can_upgrade_to_consumer"`
		// IsGatedSampling           bool   `json:"is_gated_sampling"`
		// TrialStartsAt             string `json:"trial_starts_at"`
		// TrialEndsAt               string `json:"trial_ends_at"`
		ID int `json:"id"`
	} `json:"subscription"`
	PartnershipsData            interface{} `json:"partnerships_data"`
	EnterpriseAdmin             bool        `json:"enterprise_admin"`
	EnterpriseBusinessUser      bool        `json:"enterprise_business_user"`
	UpdatedFromSap              bool        `json:"updated_from_sap"`
	UserState                   string      `json:"user_state"`
	ConnectedToFacebook         bool        `json:"connected_to_facebook"`
	NeedConfirmResidence        bool        `json:"need_confirm_residence"`
	ProjectPlusOneOfferEligible bool        `json:"project_plus_one_offer_eligible"`
	TrialDay                    interface{} `json:"trial_day"`
	TrialDaysLeft               interface{} `json:"trial_days_left"`
	ID                          int         `json:"id"`
}

type SubscriptionResponse struct {
	ExpiresAt                             string `json:"expires_at"`
	ProviderType                          string `json:"provider_type"`
	CancelAtPeriodEnd                     bool   `json:"cancel_at_period_end"`
	StartsAt                              string `json:"starts_at"`
	RenewalPurchasePlanID                 int    `json:"renewal_purchase_plan_id"`
	ChangeType                            string `json:"change_type"`
	Sponsored                             bool   `json:"sponsored"`
	IsCurrentUnderAutoUpgrade             bool   `json:"is_current_under_auto_upgrade"`
	RemainingDays                         int    `json:"remaining_days"`
	Status                                string `json:"status"`
	Active                                bool   `json:"active"`
	IsCurrentSubscriptionUnderAutoUpgrade bool   `json:"is_current_subscription_under_auto_upgrade"`
	SponsorName                           string `json:"sponsor_name"`
	CurrentGoogleProductID                string `json:"current_google_product_id"`
	ID                                    int    `json:"id"`
	PurchasePlan                          struct {
		Slug               string `json:"slug"`
		BillingCyclePeriod int    `json:"billing_cycle_period"`
		MobileDisplayName  string `json:"mobile_display_name"`
		IsAnnualPass       bool   `json:"is_annual_pass"`
		IsMonthlyPass      bool   `json:"is_monthly_pass"`
		IsInstallments     bool   `json:"is_installments"`
		ID                 int    `json:"id"`
		Product            struct {
			ID        int    `json:"id"`
			AssetSlug string `json:"asset_slug"`
			Price     string `json:"price"`
			Pricing   struct {
				CountryCode    string `json:"country_code"`
				Currency       string `json:"currency"`
				Price          int    `json:"price"`
				BasePrice      int    `json:"base_price"`
				TaxAmount      int    `json:"tax_amount"`
				TaxInclusive   bool   `json:"tax_inclusive"`
				ApplyTax       bool   `json:"apply_tax"`
				CurrencySymbol string `json:"currency_symbol"`
			} `json:"pricing"`
			PricingMarketingText  interface{} `json:"pricing_marketing_text"`
			PricingMarketingText2 interface{} `json:"pricing_marketing_text_2"`
			VanityPrice           interface{} `json:"vanity_price"`
			ProductLTV            struct {
				Gift    float64 `json:"gift"`
				Regular float64 `json:"regular"`
			} `json:"product_ltv"`
			ProductLTVInUSD                      int         `json:"product_ltv_in_usd"`
			MonetizedFlatRate                    string      `json:"monetized_flat_rate"`
			MonetizedYearPricePerMonth           string      `json:"monetized_year_price_per_month"`
			MonetizedYearPricePerMonthDiscounted interface{} `json:"monetized_year_price_per_month_discounted"`
			FlatRate                             int         `json:"flat_rate"`
			Coupon                               interface{} `json:"coupon"`
		} `json:"product"`
	} `json:"purchase_plan"`
	RenewalPurchasePlan struct {
		Slug               string `json:"slug"`
		BillingCyclePeriod int    `json:"billing_cycle_period"`
		MobileDisplayName  string `json:"mobile_display_name"`
		IsAnnualPass       bool   `json:"is_annual_pass"`
		IsMonthlyPass      bool   `json:"is_monthly_pass"`
		IsInstallments     bool   `json:"is_installments"`
		ID                 int    `json:"id"`
		Product            struct {
			ID        int    `json:"id"`
			AssetSlug string `json:"asset_slug"`
			Price     string `json:"price"`
			Pricing   struct {
				CountryCode    string `json:"country_code"`
				Currency       string `json:"currency"`
				Price          int    `json:"price"`
				BasePrice      int    `json:"base_price"`
				TaxAmount      int    `json:"tax_amount"`
				TaxInclusive   bool   `json:"tax_inclusive"`
				ApplyTax       bool   `json:"apply_tax"`
				CurrencySymbol string `json:"currency_symbol"`
			} `json:"pricing"`
			PricingMarketingText  interface{} `json:"pricing_marketing_text"`
			PricingMarketingText2 interface{} `json:"pricing_marketing_text_2"`
			VanityPrice           interface{} `json:"vanity_price"`
			ProductLTV            struct {
				Gift    float64 `json:"gift"`
				Regular float64 `json:"regular"`
			} `json:"product_ltv"`
			ProductLTVInUSD                      int         `json:"product_ltv_in_usd"`
			MonetizedFlatRate                    string      `json:"monetized_flat_rate"`
			MonetizedYearPricePerMonth           string      `json:"monetized_year_price_per_month"`
			MonetizedYearPricePerMonthDiscounted interface{} `json:"monetized_year_price_per_month_discounted"`
			FlatRate                             int         `json:"flat_rate"`
			Coupon                               interface{} `json:"coupon"`
		} `json:"product"`
	} `json:"renewal_purchase_plan"`
	Entitlement struct {
		ID int `json:"id"`
	} `json:"entitlement"`
}
