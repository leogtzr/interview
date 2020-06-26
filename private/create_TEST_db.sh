#!/bin/bash
#: name: create_PROD_db.sh
set -o xtrace
set -o nounset
set -o pipefail

readonly work_dir="$(dirname "$(readlink --canonicalize-existing "${0}")")"
readonly error_schema_file_not_found=80
readonly error_schema_setup_file_not_found=81
readonly error_env_var_not_set=82
readonly schema_template_file="${work_dir}/schema.sql"
readonly schema_setup_file="${work_dir}/setup.sql"
readonly schema_test_file="${work_dir}/schema_TEST.sql"
readonly schema_name="recruitment_interviews_test"

if [[ ! -f "${schema_template_file}" ]]; then
	echo "error: ${schema_template_file} not found." >&2
	exit ${error_schema_file_not_found}
fi

if [[ ! -f "${schema_setup_file}" ]]; then
	echo "error: ${schema_setup_file} not found." >&2
	exit ${error_schema_setup_file_not_found}
fi

sed "s/recruitment_interviews/${schema_name}/g" "${schema_template_file}" > "${schema_test_file}"

# Check if environmental variables exists ... 
if [[ -z "${DB_INTERVIEW_NAME}" ]]; then
	echo "error: DB_INTERVIEW_NAME no set" >&2
	exit ${error_env_var_not_set}
fi
if [[ -z "${DB_INTERVIEW_USER}" ]]; then
	echo "error: DB_INTERVIEW_USER no set" >&2
	exit ${error_env_var_not_set}
fi
if [[ -z "${DB_INTERVIEW_PASSWORD}" ]]; then
	echo "error: DB_INTERVIEW_PASSWORD no set" >&2
	exit ${error_env_var_not_set}
fi

if mysql -u "${DB_INTERVIEW_USER}" --password="${DB_INTERVIEW_PASSWORD}" < "${schema_test_file}"; then
	echo "Database created successfully"
	if mysql -u "${DB_INTERVIEW_USER}" --password="${DB_INTERVIEW_PASSWORD}" \
			"${schema_name}" < "${schema_setup_file}" 2> /dev/null; then
		echo "Test datadata imported"
	fi
fi

exit 0
