package core

var DBConfig map[string]string = map[string]string{
	"HOST":     "localhost",
	"PORT":     "5432",
	"USER":     "postgres",
	"PASSWORD": "MONKEYSex",
	"NAME":     "CRM-DB",
	"SSL_MODE": "disable",
}

var JWTConfig map[string]string = map[string]string{
	"ACCESS_KEY":  "asdklkgfdsasdjklkjgfdsdlkgfdsdklkfdslkjhgfdskjhgfd",
	"REFRESH_KEY": "asdklkgfdsasdjklkjgfdsdlkgfdsdklkfdslkjhgfdskjhgfd",
}

var CORSConfig map[string]any = map[string]any{
	"ALLOWED_ORIGIN": []string{
		"*",
	},
	"ALLOW_CREDENTIALS": "true",
	"ALLOWED_HEADERS": []string{
		"Content-Type", "Content-Length", "Accept-Encoding",
		"X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control",
		"X-Requested-With",
	},
	"ALLOWED_METHODS": []string{
		"POST", "OPTIONS", "GET", "PUT", "PATCH",
	},
}

var ServerConfig map[string]any = map[string]any{
	"PORT": 8080,
}
