--TEST--
Test path traversal (include)

--ENV--
AIKIDO_LOG_LEVEL=INFO
AIKIDO_BLOCK=1

--FILE--
<?php

$file = '../../bait_file.txt';

$_SERVER['HTTP_USER'] = $file;

symlink("./test/cli/bait_file.txt", $file);
    
include($file);

?>

--EXPECTREGEX--
.*Fatal error: Uncaught Exception: Aikido firewall has blocked a path traversal attack.*
