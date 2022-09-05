package lib

import (
	"context"
	"encoding/json"
	"github.com/open-policy-agent/opa/rego"
)
func PrepareAllowAllQuery() (rego.PreparedEvalQuery, error) {
	// This might need to be a file on the file system that is injected into the vttablet container.
	roleEnforcement := `
package role_enforcement
default allow_table_access = true
default allow_column_access = true
`
	return rego.New(
		rego.Query("data.role_enforcement.allow_table_access"),
		rego.Query("data.role_enforcement.allow_column_access"),
		rego.Module("role_enforcement.rego", roleEnforcement),
	).PrepareForEval(context.Background())
}

func PrepareQuery() (rego.PreparedEvalQuery, error) {

	schemaRoles := `
package schema_roles
role_permissions := {
	"customer_table_reader": [{"action": "read", "table": "customers"}],
	"customer_id_reader": [{"action": "read", "table": "customers", "column": "id"}],
	"customer_name_reader": [{"action": "read", "table": "customers", "column": "name"}],
}`

	// We can generate these from the custom ACL config store.
	userRoleAssignments := `
package user_role_assignments
# user-role assignments
user_roles := {"sam": ["customer_table_reader", "customer_id_reader", "customer_name_reader"]}
`
	// This might need to be a file on the file system that is injected into the vttablet container.
	roleEnforcement := `
package role_enforcement
import data.schema_roles
import data.user_role_assignments
import future.keywords.every

# logic that implements RBAC.
default allow_table_access = false

allow_table_access {
	# lookup the list of roles for the user
	roles := user_role_assignments.user_roles[input.user]

	# for each role in that list
	r := roles[_]

	# lookup the permissions list for role r
	permissions := schema_roles.role_permissions[r]

	# for each permission
	p := permissions[_]

	# check if the permission granted to r matches the user's request
	p == {"action": input.action, "table": input.table}
}

default allow_column_access = false

allow_column_access {
	# for every column in the input
	every column in input.columns {
		# lookup the list of roles for the user
		roles := user_role_assignments.user_roles[input.user]

		# for each role in that list
		r := roles[_]

		# lookup the permissions list for role r
		permissions := schema_roles.role_permissions[r]

		# for each permission
		p := permissions[_]

		# check if the permission granted to r matches the user's request
		p == {"action": input.action, "table": input.table, "column": column}
	}
}`
	return rego.New(
		rego.Query("data.role_enforcement.allow_table_access"),
		rego.Query("data.role_enforcement.allow_column_access"),
		rego.Module("schema_roles.rego", schemaRoles),
		rego.Module("user_role_assignments.rego", userRoleAssignments),
		rego.Module("role_enforcement.rego", roleEnforcement),
	).PrepareForEval(context.Background())
}

func Run() bool {

	query, err :=  PrepareQuery()

	if err != nil {
		panic(err.Error())
	}

	var input map[string]interface{}
	inputJson := `{
    "user": "sam",
    "action": "read",
    "table": "customers",
    "columns": [
        "id",
        "name"
    ]
}`
	json.Unmarshal([]byte(inputJson), &input)

	return Evaluate(query, input)
}

func Evaluate(query rego.PreparedEvalQuery, input map[string]interface{}) bool {
	return EvaluateWithOption(query, rego.EvalInput(input))
}

func EvaluateWithOption(query rego.PreparedEvalQuery, input rego.EvalOption) bool {
	ctx := context.TODO()
	results, err := query.Eval(ctx, input)
	if err != nil {
		panic(err.Error())
	}

	return results.Allowed()
}
