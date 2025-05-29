package autoload

import "github.com/indeedhat/dotenv"

func init() {
	_ = dotenv.Load()
}
