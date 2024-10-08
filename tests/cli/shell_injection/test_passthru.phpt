--TEST--
Test PHP shell injection (passthru)

--ENV--
AIKIDO_LOG_LEVEL=INFO
AIKIDO_BLOCK=1

--FILE--
<?php

$_SERVER['HTTP_USER'] = 'www.example`whoami`.com';

passthru('binary --domain www.example`whoami`.com');

?>

--EXPECTREGEX--
.*Fatal error: Uncaught Exception: Aikido firewall has blocked a shell injection.*