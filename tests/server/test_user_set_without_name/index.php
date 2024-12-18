<?php
$a = \aikido\set_user("12345");

if ($a == true) {
    echo "User set successfully\n";
} else {
    echo "User set failed\n";
}
?>
