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

// ttps://www.masterclass.com/course-images/attachments/Bmj4oYTzTWNpkKYPWmGi264U","duration_secs":880,"offline_enabled":false,"teaser":null,"brightcove_video_id":"5549405728001","media_uuid":"b3ad448d-b1e3-4d3e-a227-22be023cadb1","id":400,"course":{"id":13},"subchapters":[],"season":null,"chapters_instructors":[],"instructors":[]},{"abstract":"Hans discusses the importance of learning how to listen and dissect music when it works and doesn't work.","duration":"16:46","end_screen_type":"discussion","number":28,"slug":"learning-by-listening","title":"Learning by Listening","is_example_lesson":false,"video_thumb_url_alt_text":null,"updated_at":"2024-03-20T12:41:23.542-07:00","index_identifier":"Chapter|401|learning-by-listening","available_at":null,"uuid":"5795f149-182a-492b-996a-8684c3929cfe","teaser_brightcove_video_id":null,"type":"chapter","video_thumb_url":"https://www.masterclass.com/course-images/attachments/YkKofLmy5BenG1vFqLym2tEk","duration_secs":1006,"offline_enabled":false,"teaser":null,"brightcove_video_id":"5549402710001","media_uuid":"725f8aff-9920-4e24-ae4f-45e5777b57bb","id":401,"course":{"id":13},"subchapters":[],"season":null,"chapters_instructors":[],"instructors":[]},{"abstract":"All artists struggle with the challenges that come with pursuing a life in the arts. Hear Hans' advice on how to never give up and never compromise your voice.","duration":"07:07","end_screen_type":"discussion","number":29,"slug":"life-of-a-composer-part-1","title":"Life of a Composer: Part 1","is_example_lesson":false,"video_thumb_url_alt_text":null,"updated_at":"2024-03-20T12:41:23.846-07:00","index_identifier":"Chapter|402|life-of-a-composer-part-1","available_at":null,"uuid":"fc6c72a4-45dc-4207-95df-0ebb00691c13","teaser_brightcove_video_id":null,"type":"chapter","video_thumb_url":"https://www.masterclass.com/course-images/attachments/KRLt8jXbLKFkQAbz6hU9jnZL","duration_secs":427,"offline_enabled":false,"teaser":null,"brightcove_video_id":"5549394645001","media_uuid":"4c619e6c-2107-4732-a5a7-d77b535aef2c","id":402,"course":{"id":13},"subchapters":[],"season":null,"chapters_instructors":[],"instructors":[]},{"abstract":"Hans continues his discussion on an artist's life, telling you why he was inspired to pursue the life of a composer in the first place.","duration":"09:25","end_screen_type":"discussion","number":30,"slug":"life-of-a-composer-part-2","title":"Life of a Composer: Part 2","is_example_lesson":false,"video_thumb_url_alt_text":null,"updated_at":"2024-03-20T12:41:24.137-07:00","index_identifier":"Chapter|403|life-of-a-composer-part-2","available_at":null,"uuid":"152b176f-f490-472a-b7e6-269a4883082b","teaser_brightcove_video_id":null,"type":"chapter","video_thumb_url":"https://www.masterclass.com/course-images/attachments/H6cHub5pU93gEKgAow6Dbjh2","duration_secs":565,"offline_enabled":false,"teaser":null,"brightcove_video_id":"5549405729001","media_uuid":"5dc30a67-3afa-428a-8633-e60fcd310003","id":403,"course":{"id":13},"subchapters":[],"season":null,"chapters_instructors":[],"instructors":[]},{"abstract":"Listen to Hans' final words as he closes out his MasterClass and as you move forward in your career.","duration":"01:45","end_screen_type":"discussion","number":31,"slug":"closing-dfd090b1-f7a2-49d9-964d-2b8560247989","title":"Closing","is_example_lesson":true,"video_thumb_url_alt_text":null,"updated_at":"2024-06-13T01:16:44.261-07:00","index_identifier":"Chapter|392|closing-dfd090b1-f7a2-49d9-964d-2b8560247989","available_at":null,"uuid":"167cf942-a456-47ac-a366-0ce9c118402d","teaser_brightcove_video_id":null,"type":"chapter","video_thumb_url":"https://www.masterclass.com/course-images/attachments/BNVBFQ97f8VfG3sPLTTuHwYT","duration_secs":105,"offline_enabled":false,"teaser":null,"brightcove_video_id":"5549402711001","media_uuid":"54387a2f-a062-41b6-b965-f172333166ce","id":392,"course":{"id":13},"subchapters":[],"season":null,"chapters_instructors":[],"instructors":[]}],"upcoming_chapters":[],"courses_instructors":[{"course_id":13,"instructor_id":46,"position":0,"id":50,"course":{"id":13},"instructor":{"id":46}}],"instructors":[{"name":"Hans Zimmer","bio":null,"headshot_url":null,"headshot_url_alt_text":null,"id":46,"courses":[{"id":13}]}],"categories":[{"name":"Arts \u0026 Entertainment","slug":"film-tv","parent_id":null,"position":2,"cover_image":null,"cover_image_alt_text":null,"id":3,"camps":[{"id":4},{"id":41},{"id":42},{"id":43}],"courses":[{"id":458},{"id":106},{"id":61},{"id":76},{"id":1},{"id":19},{"id":352},{"id":16},{"id":18},{"id":221},{"id":153},{"id":5},{"id":58},{"id":77},{"id":205},{"id":72},{"id":9},{"id":460},{"id":59},{"id":206},{"id":222},{"id":149},{"id":432},{"id":98},{"id":104},{"id":266},{"id":309},{"id":209},{"id":92},{"id":451},{"id":90},{"id":196},{"id":7},{"id":55},{"id":86},{"id":101},{"id":462},{"id":10},{"id":2},{"id":100},{"id":256},{"id":75},{"id":11},{"id":69},{"id":107},{"id":62},{"id":397},{"id":17},{"id":67},{"id":95},{"id":57},{"id":218},{"id":145},{"id":70},{"id":56},{"id":97},{"id":82},{"id":264},{"id":12},{"id":208},{"id":73},{"id":79},{"id":194},{"id":13},{"id":265},{"id":99},{"id":80},{"id":270},{"id":454},{"id":258},{"id":192},{"id":199},{"id":213},{"id":450},{"id":201},{"id":109},{"id":147},{"id":111},{"id":189},{"id":431}]},{"name":"Music","slug":"music-entertainment","parent_id":null,"position":3,"cover_image":null,"cover_image_alt_text":null,"id":4,"camps":[{"id":4}],"courses":[{"id":2},{"id":7},{"id":12},{"id":11},{"id":55},{"id":67},{"id":79},{"id":92},{"id":98},{"id":104},{"id":192},{"id":194},{"id":196},{"id":13},{"id":111},{"id":201},{"id":213},{"id":264},{"id":265},{"id":266},{"id":270},{"id":349},{"id":388}]}],"seasons":[],"product":{"id":19},"mobile_web_gem":{"id":33},"pdfs":[{"title":"Class Guide","preview_image_1_alt_text":null,"url":"https://s3.us-west-1.amazonaws.com/www-static.masterclass.com/attachments/aj7blb6qpfhcfcbw0tigwoopd7n0?response-content-disposition=attachment%3Bfilename%3DHZ_classguide_en-US.pdf\u0026response-content-type=application%2Fpdf\u0026X-Amz-Algorithm=AWS4-HMAC-SHA256\u0026X-Amz-Credential=AKIARKYVM5WKMTWGFDV5%2F20240910%2Fus-west-1%2Fs3%2Faws4_request\u0026X-Amz-Date=20240910T080105Z\u0026X-Amz-Expires=43200\u0026X-Amz-SignedHeaders=host\u0026X-Amz-Signature=113d6882eb07a0fe08acba65376d5e6e25c0d7ae2781c5786c0a95cdf21d64f4","preview_image_1":"https://www.masterclass.com/course-images/attachments/pa6521jk5fh2r93m5b38hf2r2g2y","is_workbook":true,"id":347}],"all_pdfs":[{"title":"Class Guide","preview_image_1_alt_text":null,"url":"https://s3.us-west-1.amazonaws.com/www-static.masterclass.com/attachments/aj7blb6qpfhcfcbw0tigwoopd7n0?response-content-disposition=attachment%3Bfilename%3DHZ_classguide_en-US.pdf\u0026response-content-type=application%2Fpdf\u0026X-Amz-Algorithm=AWS4-HMAC-SHA256\u0026X-Amz-Credential=AKIARKYVM5WKMTWGFDV5%2F20240910%2Fus-west-1%2Fs3%2Faws4_request\u0026X-Amz-Date=20240910T080105Z\u0026X-Amz-Expires=43200\u0026X-Amz-SignedHeaders=host\u0026X-Amz-Signature=113d6882eb07a0fe08acba65376d5e6e25c0d7ae2781c5786c0a95cdf21d64f4","preview_image_1":"https://www.masterclass.com/course-images/attachments/pa6521jk5fh2r93m5b38hf2r2g2y","is_workbook":true,"id":347}],"welcome_survey":{"id":3,"course":{"id":13}},"primary_category":{"id":3},"upcoming_chapter":null}
type CourseResponse struct {
	Title                      string      `json:"title"`
	SeasonNote                 interface{} `json:"season_note"`
	Slug                       string      `json:"slug"`
	BrightcoveVideoID          string      `json:"brightcove_video_id"`
	MediaUUID                  string      `json:"media_uuid"`
	AudioFriendly              string      `json:"audio_friendly"`
	AapExclusive               bool        `json:"aap_exclusive"`
	MatureContent              bool        `json:"mature_content"`
	MultiInstructorName        interface{} `json:"multi_instructor_name"`
	InstructorTagline          string      `json:"instructor_tagline"`
	InstructorBio              string      `json:"instructor_bio"`
	HasRelatedInstructors      bool        `json:"has_related_instructors"`
	IsMultiInstructor          bool        `json:"is_multi_instructor"`
	MarketingOverview          string      `json:"marketing_overview"`
	MarketingCMColumn          string      `json:"marketing_cm_column"`
	Cinematic12x5AltText       interface{} `json:"cinematic_12x5_alt_text"`
	Cinematic16x9AltText       interface{} `json:"cinematic_16x9_alt_text"`
	ClassHeroIOSAltText        interface{} `json:"class_hero_ios_alt_text"`
	ClassSkillsAltText         interface{} `json:"class_skills_alt_text"`
	HeadshotAltText            interface{} `json:"headshot_alt_text"`
	HPFeaturedTileAltText      interface{} `json:"hp_featured_tile_alt_text"`
	HPTileAltText              interface{} `json:"hp_tile_alt_text"`
	NameplateAltText           interface{} `json:"nameplate_alt_text"`
	NameplateSVGAltText        interface{} `json:"nameplate_svg_alt_text"`
	ProductID                  int         `json:"product_id"`
	Primary1x1AltText          interface{} `json:"primary_1x1_alt_text"`
	Primary2x3AltText          interface{} `json:"primary_2x3_alt_text"`
	Primary16x9AltText         interface{} `json:"primary_16x9_alt_text"`
	Secondary16x9AltText       interface{} `json:"secondary_16x9_alt_text"`
	Outcome16x9AltText         interface{} `json:"outcome_16x9_alt_text"`
	TotalSeconds               int         `json:"total_seconds"`
	NumChapters                int         `json:"num_chapters"`
	UpdatedAt                  string      `json:"updated_at"`
	IndexIdentifier            string      `json:"index_identifier"`
	AvailableAt                interface{} `json:"available_at"`
	Tag                        interface{} `json:"tag"`
	UUID                       string      `json:"uuid"`
	Status                     string      `json:"status"`
	MarketingPrelaunchStartsAt interface{} `json:"marketing_prelaunch_starts_at"`
	MarketingPostlaunchEndsAt  interface{} `json:"marketing_postlaunch_ends_at"`
	ChatEnabled                bool        `json:"chat_enabled"`
	Type                       string      `json:"type"`
	IsSingleCut                bool        `json:"is_single_cut"`
	Headline                   string      `json:"headline"`
	ClassTileImage             string      `json:"class_tile_image"`
	ClassTileImageAltText      interface{} `json:"class_tile_image_alt_text"`
	SampleVideoID              string      `json:"sample_video_id"`
	Sample                     struct {
		BrightcoveVideoID string `json:"brightcove_video_id"`
		MediaUUID         string `json:"media_uuid"`
	} `json:"sample"`
	Primary1x1                  string        `json:"primary_1x1"`
	Primary2x3                  string        `json:"primary_2x3"`
	Primary16x9                 string        `json:"primary_16x9"`
	Cinematic12x5               string        `json:"cinematic_12x5"`
	Cinematic16x9               string        `json:"cinematic_16x9"`
	Nameplate                   string        `json:"nameplate"`
	NameplateSVG                string        `json:"nameplate_svg"`
	InstructorName              string        `json:"instructor_name"`
	WelcomeStatement            string        `json:"welcome_statement"`
	Skill                       string        `json:"skill"`
	Image2x3                    string        `json:"image_2x3"`
	Image2x3AltText             interface{}   `json:"image_2x3_alt_text"`
	Image16x9                   string        `json:"image_16x9"`
	Image16x9AltText            interface{}   `json:"image_16x9_alt_text"`
	Image4x3                    string        `json:"image_4x3"`
	Image4x3AltText             interface{}   `json:"image_4x3_alt_text"`
	Image256x135                string        `json:"image_256x135"`
	Image256x135AltText         interface{}   `json:"image_256x135_alt_text"`
	Outcome16x9                 string        `json:"outcome_16x9"`
	Overview                    string        `json:"overview"`
	ShortOverview               string        `json:"short_overview"`
	WorkbookCoverImages         []string      `json:"workbook_cover_images"`
	WorkbookCoverImagesAltTexts []string      `json:"workbook_cover_images_alt_texts"`
	WorkbookDescription         string        `json:"workbook_description"`
	ReleaseDescription          interface{}   `json:"release_description"`
	AutoPlayBrightcoveVideoID   interface{}   `json:"auto_play_brightcove_video_id"`
	AutoPlay                    interface{}   `json:"auto_play"`
	VanityName                  string        `json:"vanity_name"`
	EngagementRank              int           `json:"engagement_rank"`
	MerchandisingImages         []interface{} `json:"merchandising_images"`
	ID                          int           `json:"id"`
	Chapters                    []Chapter     `json:"chapters"`
	UpcomingChapters            []interface{} `json:"upcoming_chapters"`
	CoursesInstructors          []struct {
		CourseID     int `json:"course_id"`
		InstructorID int `json:"instructor_id"`
		Position     int `json:"position"`
		ID           int `json:"id"`
		Course       struct {
			ID int `json:"id"`
		} `json:"course"`
		Instructor struct {
			ID int `json:"id"`
		} `json:"instructor"`
	} `json:"courses_instructors"`
	Instructors []struct {
		Name               string      `json:"name"`
		Bio                interface{} `json:"bio"`
		HeadshotURL        interface{} `json:"headshot_url"`
		HeadshotURLAltText interface{} `json:"headshot_url_alt_text"`
		ID                 int         `json:"id"`
		Courses            []struct {
			ID int `json:"id"`
		} `json:"courses"`
	} `json:"instructors"`
	Categories []struct {
		Name              string      `json:"name"`
		Slug              string      `json:"slug"`
		ParentID          interface{} `json:"parent_id"`
		Position          int         `json:"position"`
		CoverImage        interface{} `json:"cover_image"`
		CoverImageAltText interface{} `json:"cover_image_alt_text"`
		ID                int         `json:"id"`
		Camps             []struct {
			ID int `json:"id"`
		} `json:"camps"`
		Courses []struct {
			ID int `json:"id"`
		} `json:"courses"`
	} `json:"categories"`
	Seasons []interface{} `json:"seasons"`
	Product struct {
		ID int `json:"id"`
	} `json:"product" `
	MobileWebGem struct {
		ID int `json:"id"`
	} `json:"mobile_web_gem"`
	Pdfs []struct {
		Title                string      `json:"title"`
		PreviewImage1AltText interface{} `json:"preview_image_1_alt_text"`
		URL                  string      `json:"url"`
		PreviewImage1        string      `json:"preview_image_1"`
		IsWorkbook           bool        `json:"is_workbook"`
		ID                   int         `json:"id" `
	} `json:"pdfs" `
	AllPDFs []struct {
		Title                string      `json:"title"`
		PreviewImage1AltText interface{} `json:"preview_image_1_alt_text"`
		URL                  string      `json:"url"`
		PreviewImage1        string      `json:"preview_image_1"`
		IsWorkbook           bool        `json:"is_workbook"`
		ID                   int         `json:"id" `
	} `json:"all_pdfs" `
	WelcomeSurvey struct {
		ID     int `json:"id"`
		Course struct {
			ID int `json:"id"`
		} `json:"course"`
	} `json:"welcome_survey"`
	PrimaryCategory struct {
		ID int `json:"id"`
	} `json:"primary_category"`
	UpcomingChapter interface{} `json:"upcoming_chapter"`
}

type Chapter struct {
	Abstract                string      `json:"abstract"`
	Duration                string      `json:"duration"`
	EndScreenType           string      `json:"end_screen_type"`
	Number                  int         `json:"number"`
	Slug                    string      `json:"slug"`
	Title                   string      `json:"title"`
	IsExampleLesson         bool        `json:"is_example_lesson"`
	VideoThumbURLAltText    interface{} `json:"video_thumb_url_alt_text"`
	UpdatedAt               string      `json:"updated_at"`
	IndexIdentifier         string      `json:"index_identifier"`
	AvailableAt             interface{} `json:"available_at"`
	UUID                    string      `json:"uuid"`
	TeaserBrightcoveVideoID interface{} `json:"teaser_brightcove_video_id"`
	Type                    string      `json:"type"`
	VideoThumbURL           string      `json:"video_thumb_url"`
	DurationSecs            int         `json:"duration_secs"`
	OfflineEnabled          bool        `json:"offline_enabled"`
	Teaser                  interface{} `json:"teaser"`
	BrightcoveVideoID       string      `json:"brightcove_video_id"`
	MediaUUID               string      `json:"media_uuid"`
	ID                      int         `json:"id"`
	Course                  struct {
		ID int `json:"id"`
	} `json:"course"`
	Subchapters         []interface{} `json:"subchapters"`
	Season              interface{}   `json:"season"`
	ChaptersInstructors []interface{} `json:"chapters_instructors"`
	Instructors         []interface{} `json:"instructors"`
}

type ChapterMetadataResponse struct {
	MediaUUID string `json:"media_uuid"`
	Duration  int    `json:"duration"`
	Poster    string `json:"poster"`
	Thumbnail string `json:"thumbnail"`
	Sources   []struct {
		Codec string `json:"codecs"`
		Type  string `json:"type"`
		Src   string `json:"src"`
	}
	TextTracks []struct {
		Src      string `json:"src"`
		SrcLang  string `json:"srclang"`
		Label    string `json:"label"`
		Kind     string `json:"kind"`
		Default  bool   `json:"default"`
		MimeType string `json:"mime_type"`
	} `json:"text_tracks"`
}
