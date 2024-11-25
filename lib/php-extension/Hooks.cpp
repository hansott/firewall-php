#include "Includes.h"

unordered_map<std::string, PHP_HANDLERS> HOOKED_FUNCTIONS = {
    /* Outgoing request */
    AIKIDO_REGISTER_FUNCTION_HANDLER_WITH_POST(curl_exec),                                  // curl_exec(CurlHandle $handle): string|bool

    /* Shell execution */
    AIKIDO_REGISTER_FUNCTION_HANDLER_EX(exec, handle_shell_execution),                      // exec(string $command, array &$output = null, int &$result_code = null): string|false
    AIKIDO_REGISTER_FUNCTION_HANDLER_EX(shell_exec, handle_shell_execution),                // shell_exec(string $command): string|false|null
    AIKIDO_REGISTER_FUNCTION_HANDLER_EX(system, handle_shell_execution),                    // system(string $command, int &$result_code = null): string|false
    AIKIDO_REGISTER_FUNCTION_HANDLER_EX(passthru, handle_shell_execution),                  // passthru(string $command, int &$result_code = null): ?false
    AIKIDO_REGISTER_FUNCTION_HANDLER_EX(popen, handle_shell_execution),                     // popen(string $command, string $mode): resource|false
    AIKIDO_REGISTER_FUNCTION_HANDLER_EX(proc_open, handle_shell_execution_with_array),      // proc_open(array|string $command, array $descriptor_spec, array &$pipes, ?string $cwd = null, ?array $env_vars = null, ?array $options = null): resource|false

    /* Path access */
    AIKIDO_REGISTER_FUNCTION_HANDLER_EX(chdir, handle_pre_file_path_access),                // chdir(string $directory): bool
    AIKIDO_REGISTER_FUNCTION_HANDLER_EX(chgrp, handle_pre_file_path_access),                // chgrp(string $filename, string|int $group): bool
    AIKIDO_REGISTER_FUNCTION_HANDLER_EX(chmod, handle_pre_file_path_access),                // chmod(string $filename, int $permissions): bool
    AIKIDO_REGISTER_FUNCTION_HANDLER_EX(chown, handle_pre_file_path_access),                // chown(string $filename, string|int $user): bool
    AIKIDO_REGISTER_FUNCTION_HANDLER_EX(copy, handle_pre_file_path_access_2),               // copy(string $from, string $to, ?resource $context = null): bool
    AIKIDO_REGISTER_FUNCTION_HANDLER_EX(file_put_contents, handle_pre_file_path_access),    // file_put_contents(string $filename, mixed $data, int $flags = 0, ?resource $context = null): int|false
    AIKIDO_REGISTER_FUNCTION_HANDLER_EX(lchgrp, handle_pre_file_path_access),               // chgrp(string $filename, string|int $group): bool
    AIKIDO_REGISTER_FUNCTION_HANDLER_EX(lchown, handle_pre_file_path_access),               // lchown(string $filename, string|int $user): bool
    AIKIDO_REGISTER_FUNCTION_HANDLER_EX(link, handle_pre_file_path_access_2),               // link(string $target, string $link): bool
    AIKIDO_REGISTER_FUNCTION_HANDLER_EX(mkdir, handle_pre_file_path_access),                // mkdir(string $directory, int $permissions = 0777, bool $recursive = false, ?resource $context = null): bool
    AIKIDO_REGISTER_FUNCTION_HANDLER_EX(move_uploaded_file, handle_pre_file_path_access_2), // move_uploaded_file(string $from, string $to): bool
    AIKIDO_REGISTER_FUNCTION_HANDLER_EX(opendir, handle_pre_file_path_access),              // opendir(string $directory, ?resource $context = null): resource|false
    AIKIDO_REGISTER_FUNCTION_HANDLER_EX(parse_ini_file, handle_pre_file_path_access),       // parse_ini_file(string $filename, bool $process_sections = false, int $scanner_mode = INI_SCANNER_NORMAL): array|false
    AIKIDO_REGISTER_FUNCTION_HANDLER_EX(readfile, handle_pre_file_path_access),             // readfile(string $filename, bool $use_include_path = false, ?resource $context = null): int|false
    AIKIDO_REGISTER_FUNCTION_HANDLER_EX(readlink, handle_pre_file_path_access),             // readlink(string $path): string|false
    AIKIDO_REGISTER_FUNCTION_HANDLER_EX(rename, handle_pre_file_path_access_2),             // rename(string $from, string $to, ?resource $context = null): bool
    AIKIDO_REGISTER_FUNCTION_HANDLER_EX(rmdir, handle_pre_file_path_access),                // rmdir(string $directory, ?resource $context = null): bool
    AIKIDO_REGISTER_FUNCTION_HANDLER_EX(scandir, handle_pre_file_path_access),              // scandir(string $directory, int $sorting_order = SCANDIR_SORT_ASCENDING, ?resource $context = null): array|false
    AIKIDO_REGISTER_FUNCTION_HANDLER_EX(symlink, handle_pre_file_path_access_2),            // symlink(string $target, string $link): bool
    AIKIDO_REGISTER_FUNCTION_HANDLER_EX(touch, handle_pre_file_path_access),                // touch(string $filename, ?int $mtime = null, ?int $atime = null): bool
    AIKIDO_REGISTER_FUNCTION_HANDLER_EX(unlink, handle_pre_file_path_access),               // unlink(string $filename, ?resource $context = null): bool

    /* Path access (with post hooking handlers) */
    AIKIDO_REGISTER_FUNCTION_HANDLER_WITH_POST_EX(file, handle_pre_file_path_access, handle_post_file_path_access),              // file(string $filename, int $flags = 0, ?resource $context = null): array|false
    AIKIDO_REGISTER_FUNCTION_HANDLER_WITH_POST_EX(file_get_contents, handle_pre_file_path_access, handle_post_file_path_access), // file_get_contents(string $filename, bool $use_include_path = false, ?resource $context = null, int $offset = 0, ?int $length = null): string|false
    AIKIDO_REGISTER_FUNCTION_HANDLER_WITH_POST_EX(fopen, handle_pre_file_path_access, handle_post_file_path_access),             // fopen(string $filename, string $mode, bool $use_include_path = false, ?resource $context = null): resource|false
};

unordered_map<AIKIDO_METHOD_KEY, PHP_HANDLERS, AIKIDO_METHOD_KEY_HASH> HOOKED_METHODS = {
    AIKIDO_REGISTER_METHOD_HANDLER(pdo, query)};

void HookFunctions() {
    for (auto &it : HOOKED_FUNCTIONS) {
        zend_function *function_data = (zend_function *)zend_hash_str_find_ptr(CG(function_table), it.first.c_str(), it.first.length());
        if (!function_data) {
            AIKIDO_LOG_WARN("Function \"%s\" does not exist!\n", it.first.c_str());
            continue;
        }
        if (it.second.original_handler) {
            AIKIDO_LOG_WARN("Function \"%s\" already hooked (original handler %p)!\n", it.first.c_str(), it.second.original_handler);
            continue;
        }

        it.second.original_handler = function_data->internal_function.handler;
        function_data->internal_function.handler = aikido_generic_handler;
        AIKIDO_LOG_INFO("Hooked function \"%s\" (original handler %p)!\n", it.first.c_str(), it.second.original_handler);
    }
}

void UnhookFunctions() {
    for (auto &it : HOOKED_FUNCTIONS) {
        zend_function *function_data = (zend_function *)zend_hash_str_find_ptr(CG(function_table), it.first.c_str(), it.first.length());
        if (!function_data) {
            AIKIDO_LOG_WARN("Function \"%s\" does not exist!\n", it.first.c_str());
            continue;
        }
        if (!it.second.original_handler) {
            AIKIDO_LOG_WARN("Cannot unhook function \"%s\" without an original handler (was not previously hooked)!\n", it.first.c_str());
            continue;
        }
        function_data->internal_function.handler = it.second.original_handler;
        AIKIDO_LOG_INFO("Unhooked function \"%s\" (original handler %p)!\n", it.first.c_str(), it.second.original_handler);
        it.second.original_handler = nullptr;
    }
}

void HookMethods() {
    for (auto &it : HOOKED_METHODS) {
        zend_class_entry *class_entry = (zend_class_entry *)zend_hash_str_find_ptr(CG(class_table), it.first.class_name.c_str(), it.first.class_name.length());
        if (!class_entry) {
            AIKIDO_LOG_WARN("Class \"%s\" does not exist!\n", it.first.class_name.c_str());
            continue;
        }

        zend_function *method = (zend_function *)zend_hash_str_find_ptr(&class_entry->function_table, it.first.method_name.c_str(), it.first.method_name.length());
        if (!method) {
            AIKIDO_LOG_WARN("Method \"%s->%s\" does not exist!\n", it.first.class_name.c_str(), it.first.method_name.c_str());
            continue;
        }

        if (it.second.original_handler) {
            AIKIDO_LOG_WARN("Method \"%s->%s\" already hooked (original handler %p)!\n", it.first.class_name.c_str(), it.first.method_name.c_str(), it.second.original_handler);
            continue;
        }

        it.second.original_handler = method->internal_function.handler;
        method->internal_function.handler = aikido_generic_handler;
        AIKIDO_LOG_INFO("Hooked method \"%s->%s\" (original handler %p)!\n", it.first.class_name.c_str(), it.first.method_name.c_str(), it.second.original_handler);
    }
}

void UnhookMethods() {
    for (auto &it : HOOKED_METHODS) {
        zend_class_entry *class_entry = (zend_class_entry *)zend_hash_str_find_ptr(CG(class_table), it.first.class_name.c_str(), it.first.class_name.length());
        if (!class_entry) {
            AIKIDO_LOG_WARN("Class \"%s\" does not exist!\n", it.first.class_name.c_str());
            continue;
        }

        zend_function *method = (zend_function *)zend_hash_str_find_ptr(&class_entry->function_table, it.first.method_name.c_str(), it.first.method_name.length());
        if (!method) {
            AIKIDO_LOG_WARN("Method \"%s->%s\" does not exist!\n", it.first.class_name.c_str(), it.first.method_name.c_str());
            continue;
        }

        if (!it.second.original_handler) {
            AIKIDO_LOG_WARN("Cannot unhook method \"%s->%s\" without an original handler (was not previously hooked)!\n", it.first.class_name.c_str(), it.first.method_name.c_str());
            continue;
        }

        method->internal_function.handler = it.second.original_handler;
        AIKIDO_LOG_INFO("Unhooked method \"%s->%s\" (original handler %p)!\n", it.first.class_name.c_str(), it.first.method_name.c_str(), it.second.original_handler);
        it.second.original_handler = nullptr;
    }
}
