#!/usr/bin/env bash

arr=("admin_user" "class_record" "class_record_file" "class_room" "course_config" "department" "grade_group_user" "group" "group_class_room" "group_classification" "group_course_teacher" "group_user" "liveclass" "liveclass_record" "liveclass_record_user" "liveclass_relation_id" "operation_log" "org" "org_manager" "raw_file" "record_config" "role" "role_permission" "sc_class" "sc_event_log" "sc_invite_record" "sc_student" "school_stage_config" "sys_dictionary_value" "sys_permission" "teaching_activity" "teaching_group" "teaching_group_user" "teaching_project" "url_path" "user" "user_manager" "user_student" "resource_knowledge" "org_knowledge" "org_course_version" "knowledge_config" "chapter_config" "chapter_file" "course_package" "org_chapter_config" "app_version_update" "login_refresh_token" "org_setting_switch" "quiz" "quiz_codepoint_record" "tql_codepoint_key" "liveclass_record_quiz_answer" "org_portal_column" "org_portal_categorie" "org_portal_categorier_file")

# for ((i = 0; i < ${#arr[@]}; i++)); do
#     echo ${arr[i]}
#     touch ./sql/129_${i}_add_${arr[i]}_table.sql
# done
# for name in ${arr[@]}; do
#     echo ${name}
#     touch
#     goose -dir ./sql -s create 1_29_add_${name}_table.sql sql
# done
for ((i = 0; i < ${#arr[@]}; i++)); do
    echo ${arr[i]}
    goose -dir ./sql -s create add_${arr[i]}_table.sql sql
done
