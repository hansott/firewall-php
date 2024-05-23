<?php
// 1. exec()
echo "Using exec():\n";
$output = array();
$return_var = 0;
exec('echo "Hello from exec!"', $output, $return_var);
print_r($output);
echo "\n";

// 2. shell_exec()
echo "Using shell_exec():\n";
$output = shell_exec('echo "Hello from shell_exec!"');
echo $output;
echo "\n";

// 3. system()
echo "Using system():\n";
$return_var = 0;
system('echo "Hello from system!"', $return_var);
echo "\n";

// 4. passthru()
echo "Using passthru():\n";
passthru('echo "Hello from passthru!"');
echo "\n";

// 5. popen()
echo "Using popen():\n";
$handle = popen('echo "Hello from popen!"', 'r');
while (!feof($handle)) {
    $buffer = fgets($handle);
    echo $buffer;
}
pclose($handle);
echo "\n";

// 6. proc_open()
echo "Using proc_open():\n";
$descriptorspec = array(
    0 => array("pipe", "r"),  // stdin is a pipe that the child will read from
    1 => array("pipe", "w"),  // stdout is a pipe that the child will write to
    2 => array("pipe", "w")   // stderr is a pipe that the child will write to
);

$process = proc_open('echo "Hello from proc_open!"', $descriptorspec, $pipes);

if (is_resource($process)) {
    while ($s = fgets($pipes[1])) {
        echo $s;
    }
    fclose($pipes[1]);
    proc_close($process);
}
echo "\n";
?>
