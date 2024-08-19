--TEST--
Test path traversal (linkinfo)

--ENV--
AIKIDO_LOG_LEVEL=INFO
AIKIDO_BLOCKING=1

--FILE--
<?php

$_SERVER['HTTP_USER'] = '../file';

$file = '../file/test.txt';
    
linkinfo($file);

?>

--EXPECTREGEX--
.*Fatal error: Uncaught Exception: Path traversal detected.*
