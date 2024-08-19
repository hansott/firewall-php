--TEST--
Test path traversal (filetype)

--ENV--
AIKIDO_LOG_LEVEL=INFO
AIKIDO_BLOCKING=1

--FILE--
<?php

$_SERVER['HTTP_USER'] = '../file';

$file = '../file/test.txt';
    
filetype($file);

?>

--EXPECTREGEX--
.*Fatal error: Uncaught Exception: Path traversal detected.*
