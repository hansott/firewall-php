# Should block request

In order to enable the user blocking and rate limiting features, the protected app can call `\aikido\should_block_request` to obtain the blocking decision for the current request and act accordingly.

## Vanilla PHP

```php
<?php

// Start the session (if needed) to track user login status
session_start();

// Example function to simulate user authentication
function getAuthenticatedUserId()
{
    // Assume the user ID is stored in the session after login
    return isset($_SESSION['user_id']) ? $_SESSION['user_id'] : null;
}

// Check if Aikido extension is loaded
if (extension_loaded('aikido')) {
    // Get the user ID (from session or other auth system)
    $userId = getAuthenticatedUserId();

    // If the user is authenticated, set the user ID in Aikido Zen context
    if ($userId) {
        \aikido\set_user($userId);
    }

    // Check blocking decision from Aikido
    $decision = \aikido\should_block_request();

    if ($decision->block) {
        if ($decision->type == "blocked") {
            // If the user is blocked, return a 403 status code
            http_response_code(403);
            echo "Your user is blocked!";
        }
        else if ($decision->type == "ratelimited") {
            // If the rate limit is exceeded, return a 429 status code
            http_response_code(429);
            if ($decision->trigger == "user") {
                echo "Your user exceeded the rate limit for this endpoint!";
            }
            else if ($decision->trigger == "ip") {
                echo "Your IP ({$decision->ip}) exceeded the rate limit for this endpoint!";
            }
        }
        exit();
    }

    // Continue handling the request
    echo "Request successful!";
}
```

## Laravel

```php
<?php

namespace App\Http\Middleware;

use Closure;
use Illuminate\Support\Facades\Auth;

class ZenBlockDecision
{
    public function handle($request, Closure $next)
    {
        // Check if Aikido extension is loaded
        if (!extension_loaded('aikido')) {
            return $next($request);
        }
		
		// Get the authenticated user's ID from Laravel's Auth system
		$userId = Auth::id();

		// If a user is authenticated, set the user in Aikido's firewall context
		if ($userId) {
			\aikido\set_user($userId);
		}

        // Check blocking decision from Aikido
        $decision = \aikido\should_block_request();

        if ($decision->block) {
            if ($decision->type == "blocked") {
                if ($decision->trigger == "user") {
                    return response('Your user is blocked!', 403);
                }
                else if ($decision->trigger == "ip") {
                    return response("Your IP is not allowed to access this endpoint!", 403);
                }
            }
            else if ($decision->type == "ratelimited") {
                if ($decision->trigger == "user") {
                    return response('Your user exceeded the rate limit for this endpoint!', 429);
                }
                else if ($decision->trigger == "ip") {
                    return response("Your IP ({$decision->ip}) exceeded the rate limit for this endpoint!", 429);
                }
            }
        }

        // Continue to the next middleware or request handler
        return $next($request);
    }
}
```
