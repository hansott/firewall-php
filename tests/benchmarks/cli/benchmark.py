

subprocess.run(['php', test_script_name, str(php_port), str(mock_port)], 
    env=dict(os.environ, PYTHONPATH=f"{test_lib_dir}:$PYTHONPATH"),
    cwd=test_script_cwd,
    check=True, timeout=600)