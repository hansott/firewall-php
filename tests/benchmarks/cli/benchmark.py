import os
import subprocess
import sys
import testlib

def read_php_code(file_path):
    """
    Reads the PHP code from a file, stripping the '<?php' tag if it exists.
    """
    with open(file_path, 'r') as file:
        php_code = file.read()
    
    # Strip '<?php' from the beginning of the file if present
    php_code = php_code.replace('<?php', '', 1).strip()
    php_code = php_code.replace('?>', '', 1).strip()
    return php_code

def generate_php_script(baseline_script, php_code_to_insert):
    """
    Generates a new PHP script by replacing the placeholder with the actual PHP code to insert.
    """
    # Replace placeholder with actual PHP code
    new_script = baseline_script.replace('// <insert PHP code here>', php_code_to_insert)
    return new_script

def execute_php_script(test_name, php_script):
    """
    Executes the provided PHP script and returns its output.
    """
    # Write the script to a temporary file
    with open(f'{test_name}.php', 'w') as file:
        file.write(php_script)
        
    if testlib.is_aikido_installed():
        filename = f"{test_name}_with_aikido.txt"
    else:
        filename = f"{test_name}_without_aikido.txt"
    
    # Execute the script and capture the output
    result = subprocess.run(['php', f'{test_name}.php', filename], 
                            env=dict(os.environ, AIKIDO_LOG_LEVEL="ERROR"))
    
    # Remove the temporary file
    os.remove(f'{test_name}.php')
    

def main(directory):
    """
    Main function to iterate through all folders in a directory and run PHP scripts.
    """
    # Baseline PHP script
    
    baseline_php_script = ""
    with open("benchmark.php", "r") as f:
        baseline_php_script = f.read()

    # Iterate through each folder in the provided directory
    for folder_name in os.listdir(directory):
        folder_path = os.path.join(directory, folder_name)
        benchmark_test_name = os.path.basename(os.path.normpath(folder_name))
        
        # Ensure it is a directory
        if os.path.isdir(folder_path):
            php_code_file = os.path.join(folder_path, 'php_code_to_test.php')
            
            # Check if the file exists
            if os.path.exists(php_code_file):
                # Read the PHP code to test
                php_code_to_insert = read_php_code(php_code_file)
                
                # Generate the new PHP script
                new_php_script = generate_php_script(baseline_php_script, php_code_to_insert)
                
                # Execute the new PHP script
                execute_php_script(benchmark_test_name, new_php_script)

if __name__ == "__main__":
    main(sys.argv[1])
