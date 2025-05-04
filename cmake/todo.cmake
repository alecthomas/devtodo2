macro(generateTodo SRC_DIR)
    file(GLOB_RECURSE TODO_DATABASES ${SRC_DIR}/*.todo ${SRC_DIR}/*.todo2)
    list(LENGTH TODO_DATABASES TODOS_LENGTH)

    if(${TODOS_LENGTH} GREATER 0)
        foreach(val ${TODO_DATABASES})
            get_filename_component(todo_srcdir ${val} DIRECTORY)
            #get_filename_component(todo_srcdir_rel ${val} DIRECTORY BASE_DIR ${SRC_DIR})
            string(REPLACE ${SRC_DIR} "{PRJ}" todo_srcdir_rel ${val})
            string(REPLACE ${SRC_DIR} ${CMAKE_BINARY_DIR} todo_src ${val})
            string(APPEND todo_src ".h")
            get_filename_component(todo_workdir ${todo_src} DIRECTORY)

            set(todo_v1 '')
            set(todo_v2 '')
            if ("${val}" STREQUAL "${todo_srcdir}/.todo2")
                set(todo_v2 "${val}")
            else ()
                set(todo_v1 "${val}")
                if (EXISTS "${todo_srcdir}/.todo2") # ignore v1
                    continue()
                endif()
            endif()

            execute_process(
                #COMMAND todo -f +children -T --database ${val}
                COMMAND /usr/local/bin/todo2 --legacy-file=${todo_v1} --file=${todo_v2} -T
                WORKING_DIRECTORY ${todo_workdir}

                OUTPUT_VARIABLE out_var
                OUTPUT_STRIP_TRAILING_WHITESPACE
                ERROR_VARIABLE out_error
                ERROR_STRIP_TRAILING_WHITESPACE
                RESULT_VARIABLE result_exit_code
            )

            if (NOT EXISTS ${todo_workdir}/TODO)
                message(NOTICE "todo out_var: ${out_var}, out_error: ${out_error}, result_exit_code: ${result_exit_code}")
                message(WARNING "todo can't generate the TODO file in dir: ${todo_workdir}, source: ${val}")
                continue()
            endif()

            file(READ ${todo_workdir}/TODO todo_data)
            string(LENGTH "${todo_data}" DATA_LEN)
            if (DATA_LEN GREATER 0)
                #find priority
                string(REGEX REPLACE "\\- ([^-]*)priority veryhigh" "- BUG: \\1priority veryhigh" todo_data "${todo_data}")
                string(REGEX REPLACE "\\- ([^-]*)priority high" "- FIXME: \\1priority high" todo_data "${todo_data}")
                string(REGEX REPLACE "\\- ([^-]*)priority medium" "- WARNING: \\1priority medium" todo_data "${todo_data}")
                string(REGEX REPLACE "\\- ([^-]*)priority (low|verylow)" "- TODO: \\1priority \\2" todo_data "${todo_data}")

                #update priority
                #string(REPLACE "- " "- TODO: " todo_data ${todo_data})
                file(WRITE ${todo_src} "/* ${todo_srcdir_rel}\n\n${todo_data}*/")

                list(APPEND TODO_SOURCES ${todo_src})
            endif()

            #set_source_files_properties(${todo_src} PROPERTIES GENERATED 1)
        endforeach()
        add_library(TODO INTERFACE)
        target_include_directories(TODO INTERFACE ${CMAKE_BINARY_DIR})
        target_sources(TODO PUBLIC ${TODO_SOURCES})
    endif()
endmacro()
