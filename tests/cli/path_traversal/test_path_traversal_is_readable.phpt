--TEST--
Test path traversal (is_readable)

--ENV--
AIKIDO_LOG_LEVEL=INFO
AIKIDO_BLOCKING=1

--FILE--
<?php

$_SERVER['HTTP_USER'] = '../file';

$file = '../file/test.txt';
    
is_readable($file);

?>

--EXPECTREGEX--
.*Fatal error: Uncaught Exception: Aikido firewall has blocked a path traversal attack.*
