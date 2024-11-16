package auth_model

import (
	"yes4all/ads-noti-api/pkg/utils/constants"

	"github.com/dgrijalva/jwt-go"
)

const (
	RoleAdsAdvertiserCamapaignEdit string = "ads_advertiser_campaign_edit"
	RoleAdsAdvertiserCamapaignView string = "ads_advertiser_campaign_view"
	RoleAdsNemoReportEdit          string = "ads_nemo_report_edit"
	RoleAdsNemoReportView          string = "ads_nemo_report_view"
	RoleAdsNemoTransactionsEdit    string = "ads_nemo_transactions_edit"
	RoleAdsNemoTransactionsView    string = "ads_nemo_transactions_view"
	RoleAdsNemoAdvBillingEdit      string = "ads_adv_billing_edit"
	RoleAdsNemoAdvBillingView      string = "ads_adv_billing_view"
	RoleAdsNemoAmazonStoresEdit    string = "ads_amazon_stores_edit"
	RoleAdsNemoAmazonStoresView    string = "ads_amazon_stores_view"
	RoleAdsNemoBrandPostsEdit      string = "ads_brand_posts_edit"
	RoleAdsNemoBrandPostsView      string = "ads_brand_posts_view"
	RoleAdsNotificationEdit        string = "ads_notification_edit"
	RolePortfolioCreate            string = "ads_portfolio_create"
	RoleProfile                    string = "ads_profile"
	RoleCampaignCreate             string = "ads_campaign_create"
	AwsPortfolioIDsEdit            string = "AWS_PORTFOLIO_IDS_EDIT"
	RolePortfolioEdit              string = "ads_portfolio_edit"
	RolePortfolioView              string = "ads_portfolio_view"
	AdminYAMS                             = "admin_yams"
	RoleAdmin                      string = "ads_admin"
	Admin                          string = "ADMIN"
	InsideUser                     string = "INSIDE_USER"
	AwsProfileID                   string = "AWS_PROFILE_ID"
	EmailAdmin                     string = "admin@yes4all"
	AwsPortfolioIDs                string = "AWS_PORTFOLIO_IDS"
	AwsProfileInfo                 string = "AWS_PROFILE_INFO"
	HttpAwsTimeout                 string = "HTTP_AWS_TIMEOUT"
	FromSource                            = "FROM_SOURCE"
	MetricFunction                        = "METRIC_FUNCTION"
)

var (
	AllRole = []string{
		RoleAdsAdvertiserCamapaignEdit,
		RoleAdsAdvertiserCamapaignView,
		RoleAdsNemoReportEdit,
		RoleAdsNemoReportView,
		RoleAdsNemoTransactionsEdit,
		RoleAdsNemoTransactionsView,
		RoleAdsNemoAdvBillingEdit,
		RoleAdsNemoAdvBillingView,
		RoleAdsNemoAmazonStoresEdit,
		RoleAdsNemoAmazonStoresView,
		RoleAdsNemoBrandPostsEdit,
		RoleAdsNemoBrandPostsView,
		RolePortfolioEdit,
		RolePortfolioView,
		RoleAdmin,
		RolePortfolioCreate,
		RoleProfile,
		RoleAdsNotificationEdit,
	}

	AllRole1 = []string{
		RoleAdsNemoReportEdit,
		RoleAdsNemoReportView,
		RoleAdsNemoTransactionsEdit,
		RoleAdsNemoTransactionsView,
		RoleAdsNemoAdvBillingEdit,
		RoleAdsNemoAdvBillingView,
		RoleAdsNemoAmazonStoresEdit,
		RoleAdsNemoAmazonStoresView,
		RoleAdsNemoBrandPostsEdit,
		RoleAdsNemoBrandPostsView,
	}

	MapRolesResponse = []RoleDetail{
		{
			RoleName:  "Create Portfolios",
			RoleValue: RolePortfolioCreate,
			Type:      constants.Portfolio,
			IsActive:  false,
		},
		{
			RoleName:  "Admin Profiles",
			RoleValue: RoleAdmin,
			Type:      RoleAdmin,
			IsActive:  false,
		},
	}

	RolesUser  = []string{RolePortfolioCreate}
	RolesAdmin = []string{RolePortfolioCreate, RoleAdmin}
)

type InsideAuthClaim struct {
	Roles             []string `json:"roles"`
	Name              string   `json:"name"`
	PreferredUsername string   `json:"preferred_username"`
	GivenName         string   `json:"given_name"`
	FamilyName        string   `json:"family_name"`
	Email             string   `json:"email"`
	Picture           string   `json:"picture"`
	EmailVerified     bool     `json:"email_verified"`
	jwt.StandardClaims
}

type RoleDetail struct {
	RoleName  string `json:"roleName"`
	RoleValue string `json:"roleValue"`
	Type      string `json:"type"`
	IsActive  bool   `json:"isActive"`
}
