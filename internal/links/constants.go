package links

const CATCHALL = "/*"

// ===========================================================================
// == AUTH LINKS
// ===========================================================================

const AUTH_AUTH = "/auth/auth"
const AUTH_LOGIN = "/auth/login"
const AUTH_LOGOUT = "/auth/logout"
const AUTH_REGISTER = "/auth/register"

// ===========================================================================
// == ADMIN LINKS
// ===========================================================================

const ADMIN_HOME = "/admin"
const ADMIN_BLOG = ADMIN_HOME + "/blog"
const ADMIN_CHAT = ADMIN_HOME + "/chat"
const ADMIN_CHAT_RAG = ADMIN_HOME + "/chat/rag"
const ADMIN_CMS = ADMIN_HOME + "/cms"
const ADMIN_CMS_OLD = ADMIN_HOME + "/cmsold"
const ADMIN_FILE_MANAGER = ADMIN_HOME + "/file-manager"
const ADMIN_LOGS = ADMIN_HOME + "/logs"
const ADMIN_MEDIA = ADMIN_HOME + "/media"
const ADMIN_SHOP = ADMIN_HOME + "/shop"
const ADMIN_STATS = ADMIN_HOME + "/stats"
const ADMIN_TASKS = ADMIN_HOME + "/tasks"
const ADMIN_USERS = ADMIN_HOME + "/users"
const ADMIN_USERS_USER_CREATE = ADMIN_USERS + "/user-create"
const ADMIN_USERS_USER_DELETE = ADMIN_USERS + "/user-delete"
const ADMIN_USERS_USER_IMPERSONATE = ADMIN_USERS + "/user-impersonate"
const ADMIN_USERS_USER_MANAGER = ADMIN_USERS + "/user-manager"
const ADMIN_USERS_USER_UPDATE = ADMIN_USERS + "/user-update"

// ===========================================================================
// == USER LINKS
// ===========================================================================

const USER_HOME = "/user"

// User Orders
const USER_ORDERS = USER_HOME + "/orders"
const USER_ORDER_CREATE = USER_ORDERS + "/create"
const USER_ORDER_CREATE_PAYMENT_BEGIN = USER_ORDER_CREATE + "/payment-begin"
const USER_ORDER_DELETE = USER_ORDERS + "/delete"
const USER_ORDER_LIST = USER_ORDERS + "/list"

const USER_PROFILE = USER_HOME + "/profile"

// User Subscription
const USER_SUBSCRIPTION = USER_HOME + "/subscription"
const USER_SUBSCRIPTION_PLAN_SELECT = USER_SUBSCRIPTION + "/plan-select"
const USER_SUBSCRIPTION_PLAN_SELECT_AJAX = USER_SUBSCRIPTION + "/plan-select-ajax"
const USER_SUBSCRIPTION_PAYMENT_CANCELED = USER_SUBSCRIPTION + "/payment-canceled"
const USER_SUBSCRIPTION_PAYMENT_SUCCESS = USER_SUBSCRIPTION + "/payment-success"


// ===========================================================================
// == CHAT LINKS
// ===========================================================================

const CHAT_HOME = "/chat"

// ===========================================================================
// == WEBSITE LINKS
// ===========================================================================

const HOME = "/"
const BLOG = HOME + "blog"
const BLOG_POST = BLOG + "/post"
const BLOG_POST_WITH_REGEX = BLOG_POST + "/{id:[0-9]+}"
const BLOG_POST_WITH_REGEX2 = BLOG_POST + "/{id:[0-9]+}/{title}"
const BLOG_POST_01 = BLOG_POST + "/:id"
const BLOG_POST_02 = BLOG_POST + "/:id/:title"

const CONTACT = HOME + "contact"
const FILES = HOME + "files" + CATCHALL
const FLASH = HOME + "flash"
const LIVEFLUX = HOME + "liveflux"
const MEDIA = HOME + "media" + CATCHALL
const PAYPAL_CANCEL = "/paypal/cancel"
const PAYPAL_NOTIFY = "/paypal/notify"
const PAYPAL_SUCCESS = "/paypal/success"
const PAYMENT_CANCELED = "/payment/canceled"
const PAYMENT_SUCCESS = "/payment/success"
const RESOURCES = HOME + "resources" + CATCHALL

const SHOP = "/shop"
const SHOP_PRODUCT = SHOP + "/product"
const SHOP_PRODUCT_WITH_REGEX = SHOP_PRODUCT + "/{id:[0-9]+}"
const SHOP_PRODUCT_WITH_REGEX2 = SHOP_PRODUCT + "/{id:[0-9]+}/{title}"

const SITEMAPXML = HOME + "sitemap.xml"

const PRIVACY_POLICY = HOME + "privacy-policy"
const TERMS_OF_USE = HOME + "terms-of-use"
const THEME = HOME + "theme"
const THUMB = HOME + "th/{extension:[a-z]+}/{size:[0-9x]+}/{quality:[0-9]+}/*"
const WIDGET = HOME + "widget"
