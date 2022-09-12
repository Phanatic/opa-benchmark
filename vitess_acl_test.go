package lib

import (
	"vitess.io/vitess/go/json2"
	tableaclpb "vitess.io/vitess/go/vt/proto/tableacl"

	querypb "vitess.io/vitess/go/vt/proto/query"
	"vitess.io/vitess/go/vt/tableacl"
	"vitess.io/vitess/go/vt/tableacl/simpleacl"

	"testing"
	//"fmt"
	"github.com/pkg/errors"
)

func Benchmark_Vitess_Classic_Authorized(b *testing.B) {

	if err := initAclConfig(); err != nil {
		panic(err.Error())
	}

	callerID := &querypb.VTGateCallerID{
		Username: "planetscale-reader",
		Groups:   []string{"planetscale-reader"},
	}

	for i := 0; i < b.N; i++ {
		// vttablet currently builds Authorized check this per query plan.
		// We duplicate it here per invocation, can ideally be done once per query.
		authorized := tableacl.Authorized("reminders", tableacl.READER)
		if !authorized.IsMember(callerID) {
			panic("test failed")
		}
	}
}

func Benchmark_Vitess_Classic_Prepared_Authorized(b *testing.B) {

	if err := initAclConfig(); err != nil {
		panic(err.Error())
	}

	callerID := &querypb.VTGateCallerID{
		Username: "planetscale-reader",
		Groups:   []string{"planetscale-reader"},
	}

	// vttablet currently builds Authorized check this per query plan.
	authorized := tableacl.Authorized("reminders", tableacl.READER)

	for i := 0; i < b.N; i++ {
		if !authorized.IsMember(callerID) {
			panic("test failed")
		}
	}
}

func initAclConfig() error {
	aclconfig := `{
	    "table_groups":
	    [
	        {
	            "name": "planetscale user groups",
	            "table_names_or_prefixes":
	            [
	                "reminders"
	            ],
	            "readers":
	            [
	                "planetscale-reader",
	                "planetscale-writer",
	                "planetscale-admin"
	            ],
	            "writers":
	            [
	                "planetscale-writer",
	                "planetscale-writer-only",
	                "planetscale-admin"
	            ],
	            "admins":
	            [
	                "planetscale-admin"
	            ]
	        }
	    ]
	}`

	data := []byte(aclconfig)

	config := &tableaclpb.Config{}
	if jsonErr := json2.Unmarshal(data, config); jsonErr != nil {
		return errors.Wrapf(jsonErr, "unable to unmarshal Table ACL data as JSON: %s", data)
	}

	_, err := tableacl.GetCurrentACLFactory()
	if err != nil {
		//return errors.Wrap(err, "unable to get current acl factory")
		tableacl.Register("simpleacl", &simpleacl.Factory{})
		//tableacl.Register("simpleacl", &simpleacl.Factory{})
		tableacl.SetDefaultACL("simpleacl")
	}

	if err := tableacl.InitFromProto(config); err != nil {
		return errors.Wrapf(err, "InitFromProto failed")
	}

	if err := tableacl.ValidateProto(config); err != nil {
		return errors.New("Acl config is invalid, exiting")
	}

	return nil
}
