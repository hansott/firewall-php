--TEST--
Test path traversal (file)

--ENV--
AIKIDO_LOG_LEVEL=INFO
AIKIDO_BLOCKING=1

--FILE--
<?php

$_SERVER['HTTP_USER'] = '../file';

$file = '../file/test.txt';
    
file($file);

?>

--EXPECTREGEX--
.*Fatal error: Uncaught Exception: Path traversal detected.*
