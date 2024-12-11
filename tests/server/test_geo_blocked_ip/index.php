<?php

$_SERVER['REMOTE_ADDR'] = '5.8.19.22';

if (extension_loaded('aikido')) {
    $decision = \aikido\should_block_request();

    if ($decision->block && $decision->type == "blocked" && $decision->trigger == "geoip") {
        http_response_code(403);
        echo "Your IP address is blocked due to geo restrictions. (Your IP: {$decision->ip})";
        exit();
    }
}

echo "Something";

?>
