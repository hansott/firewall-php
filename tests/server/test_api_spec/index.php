<?php
    
echo "Something";

if (extension_loaded('aikido')) {
    $decision = \aikido\should_block_request();
}

?>
