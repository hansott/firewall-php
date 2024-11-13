<?php

if (extension_loaded('aikido')) {
    $decision = \aikido\should_block_request();

    if ($decision->block && $decision->type == "blocked" && $decision->trigger == "ip") {
        http_response_code(403);
        echo "Your IP address is not allowed to access this endpoint! (Your IP: {$decision->ip})";
        exit();
    }
}

?>
