<?php
$a = aikido_set_user("12345", "Tudor");

if ($a == true) {
    echo "User set successfully\n";
} else {
    echo "User set failed\n";
}

?>
