--TEST--
Test PHP shell injection (shell_exec)

--ENV--
AIKIDO_LOG_LEVEL=INFO
AIKIDO_BLOCKING=1

--FILE--
<?php

$_SERVER['HTTP_USER'] = 'www.example`whoami`.com';

$output = shell_exec('binary --domain www.example`whoami`.com');
echo $output;

?>

--EXPECTREGEX--
.*Fatal error: Uncaught Exception: Shell injection detected in.*