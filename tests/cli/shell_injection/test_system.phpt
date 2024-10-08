--TEST--
Test PHP shell injection (system)

--ENV--
AIKIDO_LOG_LEVEL=INFO
AIKIDO_BLOCK=1

--FILE--
<?php

$_SERVER['HTTP_USER'] = 'www.example`whoami`.com';

$output = array();
$return_var = 0;
system('binary --domain www.example`whoami`.com', $return_var);
echo $return_var;
echo "\n";

?>

--EXPECTREGEX--
.*Fatal error: Uncaught Exception: Aikido firewall has blocked a shell injection.*

