package auth

func GetSession() string {
	// Internally calling private function to fetch data
	return extractSession()
}

func extractSession() string {
	return "LoggedIn successfully"
}
