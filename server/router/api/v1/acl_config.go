package v1

var authenticationAllowlistMethods = map[string]bool{
	"/monitor/health":     true,
	"/v1/signup":          true,
	"/v1/login":           true,
}

func isUnauthorizeAllowedMethod(fullMethodName string) bool {
	return authenticationAllowlistMethods[fullMethodName]
}
