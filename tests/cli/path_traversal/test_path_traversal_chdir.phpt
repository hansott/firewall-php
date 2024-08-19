--TEST--
Test path traversal (chdir)

--ENV--
AIKIDO_LOG_LEVEL=INFO
AIKIDO_BLOCKING=1

--FILE--
<?php

$_SERVER['HTTP_USER'] = '../file';

$file = '../file/test.txt';
    
chdir($file, 0);

?>

--EXPECTREGEX--
.*Fatal error: Uncaught Exception: Path traversal detected.*
