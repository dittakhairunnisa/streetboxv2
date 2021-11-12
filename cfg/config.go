package cfg

var Config = struct {
	Env string
	Api struct {
		Host string
		Port string
	}
	JwtKey string
	Path   struct {
		Image string
		Doc   string
		Csv   string
	}
	Certificate struct {
		Dir      string
		Filename string
	}
	GoogleApplicationCredentials string
	Web                          struct {
		Host string
		Port string
	}
	Smpt struct {
		Host     string
		Port     string
		Email    string
		Password string
	}
	Postgres struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}
	Log struct {
		Filename string
		MaxSize  int
		MaxAge   int
	}
	Xendit struct {
		ApiKey            string
		CallbackQrisToken string
		ApiHost           string
		CallbackQrisUrl   string
	}
	Fcm struct {
		ServerToken1      string
		ServerToken2 			string
		ServerToken3    	string
		ProjectID    			string
	}
	Topic struct {
		OrderOnline      	string
		TrxPushNotif 			string
		VisitOrder    		string
	}
}{}
