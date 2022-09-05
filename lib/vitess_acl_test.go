package lib

import (
	tableaclpb "vitess.io/vitess/go/vt/proto/tableacl"
	"vitess.io/vitess/go/json2"

	"vitess.io/vitess/go/vt/tableacl"
	"vitess.io/vitess/go/vt/tableacl/simpleacl"
	querypb "vitess.io/vitess/go/vt/proto/query"

	"testing"
	//"fmt"
	"github.com/pkg/errors"
)

func Benchmark_Vitess_Classic_ReaderGroup(b *testing.B) {
	var (
		readAcl  *tableacl.ACLResult
	)

	if err := initAclConfig(); err != nil {
		panic(err.Error())
	}

	if role, ok := tableacl.RoleByName("READER"); ok {
		readAcl = tableacl.Authorized("reminders", role)
	}

	callerID := &querypb.VTGateCallerID{
		Username: "planetscale-reader",
		Groups:   []string{"planetscale-reader"},
	}

	for i := 0; i < b.N; i++ {
		if !readAcl.IsMember(callerID) {
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
		tableacl.Register("simpleacl_phani_raj", &simpleacl.Factory{})
		//tableacl.Register("simpleacl", &simpleacl.Factory{})
		tableacl.SetDefaultACL("simpleacl_phani_raj")
	}

	if err := tableacl.InitFromProto(config); err != nil {
		return errors.Wrapf(err, "InitFromProto failed")
	}

	if err := tableacl.ValidateProto(config); err != nil {
		return errors.New("Acl config is invalid, exiting")
	}

	return nil
}
