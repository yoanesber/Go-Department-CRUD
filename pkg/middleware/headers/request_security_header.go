package headers

import "github.com/gin-gonic/gin"

// RequestSecurityHeaders is a middleware function that sets various security-related HTTP headers.
// These headers help protect against common web vulnerabilities and improve the security of the application.
const (
	xFrameOptions                = "X-Frame-Options"
	xFrameOptionsValue           = "DENY"
	xContentTypeOptions          = "X-Content-Type-Options"
	xContentTypeOptionsValue     = "nosniff"
	xssProtection                = "X-XSS-Protection"
	xssProtectionValue           = "1; mode=block"
	strictTransportSecurity      = "Strict-Transport-Security"
	strictTransportSecurityValue = "max-age=31536000; includeSubDomains; preload"
	referrerPolicy               = "Referrer-Policy"
	referrerPolicyValue          = "no-referrer"
	permissionsPolicy            = "Permissions-Policy"
	permissionsPolicyValue       = "geolocation=(self), microphone=()"
)

func RequestSecurityHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set(xFrameOptions, xFrameOptionsValue)
		c.Writer.Header().Set(xContentTypeOptions, xContentTypeOptionsValue)
		c.Writer.Header().Set(xssProtection, xssProtectionValue)
		c.Writer.Header().Set(strictTransportSecurity, strictTransportSecurityValue)
		c.Writer.Header().Set(referrerPolicy, referrerPolicyValue)
		c.Writer.Header().Set(permissionsPolicy, permissionsPolicyValue)

		c.Next()
	}
}
