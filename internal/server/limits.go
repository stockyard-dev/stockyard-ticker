package server

import "github.com/stockyard-dev/stockyard-ticker/internal/license"

type Limits struct{ MaxAccounts int }

var freeLimits = Limits{MaxAccounts: 3}
var proLimits = Limits{MaxAccounts: 0}

func LimitsFor(info *license.Info) Limits {
	if info != nil && info.IsPro() { return proLimits }
	return freeLimits
}
