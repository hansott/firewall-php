![Zen by Aikido for PHP](./docs/banner.svg)

# Zen, in-app firewall for PHP | by Aikido

Zen, your in-app firewall for peace of mind – at runtime.

Zen is an embedded Web Application Firewall that autonomously protects your PHP apps against common and critical attacks.

Zen protects your PHP apps by preventing user input containing dangerous strings, thus protecting agains attacks like SQL injection. It runs on the same server as your PHP app for simple [installation](#Installation) and zero maintenance.

## Features

Zen will autonomously protect your PHP applications from the inside against:

* 🛡️ [NoSQL injection attacks](https://www.aikido.dev/blog/web-application-security-vulnerabilities)
* 🛡️ [SQL injection attacks]([https://www.aikido.dev/blog/web-application-security-vulnerabilities](https://owasp.org/www-community/attacks/SQL_Injection))
* 🛡️ [Command injection attacks](https://owasp.org/www-community/attacks/Command_Injection)
* 🛡️ [Path traversal attacks](https://owasp.org/www-community/attacks/Path_Traversal)
* 🛡️ [Server-side request forgery (SSRF)](./docs/ssrf.md)

Zen operates autonomously on the same server as your PHP app to:

* ✅ Secure your app like a classic web application firewall (WAF), but with none of the infrastructure or cost.

## Installation

Zen for PHP comes as a single package that needs to be installed on the system to be protected.
Prerequisites:
* Ensure you have sudo privileges on your system.
* Check that you have a supported PHP version installed (PHP version >= 7.3).
* Make sure to use the appropriate commands for your system or cloud provider.

### Manual installation

#### For Red Hat-based Systems (RHEL, CentOS, Fedora)

`rpm -Uvh https://aikido-firewall.s3.eu-west-1.amazonaws.com/aikido-php-firewall.x86_64.rpm`

#### For Debian-based Systems (Debian, Ubuntu)

`dpkg -i https://aikido-firewall.s3.eu-west-1.amazonaws.com/aikido-php-firewall.x86_64.deb`

### Cloud providers

#### AWS Elastic beanstalk

Create a new file in `.ebextensions/01_aikido_php_firewall.config` with the following content:
```
commands:
  aikido-php-firewall:
    command: "rpm -Uvh https://aikido-firewall.s3.eu-west-1.amazonaws.com/aikido-php-firewall.x86_64.rpm"
    ignoreErrors: true

files: 
  "/opt/elasticbeanstalk/tasks/bundlelogs.d/aikido-php-firewall.conf" :
    mode: "000755"
    owner: root
    group: root
    content: |
      /var/log/aikido-*/*.log

  "/opt/elasticbeanstalk/tasks/taillogs.d/aikido-php-firewall.conf" :
    mode: "000755"
    owner: root
    group: root
    content: |
      /var/log/aikido-*/*.log
```

## Supported libraries and frameworks

### PHP versions

Zen for PHP 7.3+ works for:

### Web frameworks

* ✅ [`Laravel`](https://laravel.com/)
* ✅ [`Symphony`](https://symfony.com/)
* ✅ [`CodeIgniter`](https://codeigniter.com/)

### Database drivers
* ✅ [`PDO`](https://www.php.net/manual/en/book.pdo.php)
    * ✅ [`MySQL`](https://www.php.net/manual/en/ref.pdo-mysql.php)
    * ✅ [`Oracle`](https://www.php.net/manual/en/ref.pdo-oci.php)
    * ✅ [`PostgreSQL`](https://www.php.net/manual/en/ref.pdo-pgsql.php)
    * ✅ [`ODBC and DB2`](https://www.php.net/manual/en/ref.pdo-odbc.php)
    * ✅ [`Firebird`](https://www.php.net/manual/en/ref.pdo-firebird.php)
    * ✅ [`Microsoft SQL Server`](https://www.php.net/manual/en/ref.pdo-dblib.php)
    * ✅ [`SQLite`](https://www.php.net/manual/en/ref.pdo-sqlite.php)
* 🚧 [`MySQLi`](https://www.php.net/manual/en/book.mysqli.php)
* 🚧 [`MongoDB`](https://www.php.net/manual/en/set.mongodb.php)
* 🚧 [`Oracle OCI8`](https://www.php.net/manual/en/book.oci8.php)
* 🚧 [`PostgreSQL`](https://www.php.net/manual/en/book.pgsql.php)
* 🚧 [`SQLite3`](https://www.php.net/manual/en/book.sqlite3.php)

### Outgoing requests libraries
* ✅ [`cURL`](https://www.php.net/manual/en/book.curl.php)
* ✅ [`GuzzleHttp`](https://docs.guzzlephp.org/en/stable/)
* 🚧 [`file_get_contents`](https://www.php.net/manual/en/function.file-get-contents.php)
* 🚧 [`HTTP_Request2`](https://pear.php.net/package/http_request2)
* 🚧 [`Symfony\HTTPClient`](https://symfony.com/doc/current/http_client.html)

## Reporting to your Aikido Security dashboard

> Aikido is your no nonsense application security platform. One central system that scans your source code & cloud, shows you what vulnerabilities matter, and how to fix them - fast. So you can get back to building.

Zen is a new product by Aikido. Built for developers to level up their security. While Aikido scans, get Zen for always-on protection. 

You can use some of Zen’s features without Aikido, of course. Peace of mind is just a few lines of code away.

But you will get the most value by reporting your data to Aikido.

You will need an Aikido account and a token to report events to Aikido. If you don't have an account, you can [sign up for free](https://app.aikido.dev/login).

Here's how:
* [Log in to your Aikido account](https://app.aikido.dev/login).
* Go to [Zen](https://app.aikido.dev/runtime/services).
* Go to apps.
* Click on **Add app**.
* Choose a name for your app.
* Click **Generate token**.
* Copy the token.
* Set the token as an environment variable, `AIKIDO_TOKEN`

## Running in production (blocking) mode

By default, Zen will only detect and report attacks to Aikido.

To block requests, set the `AIKIDO_BLOCKING` environment variable to `true`.

See [Reporting to Aikido](#reporting-to-your-aikido-security-dashboard) to learn how to send events to Aikido.

## Benchmarks 

<To be added>

See [benchmarks](tests/benchmarks/) folder for more.

## Bug bounty program

Our bug bounty program is public and can be found by all registered Intigriti users at: https://app.intigriti.com/researcher/programs/aikido/aikidoruntime

## Contributing

See [CONTRIBUTING.md](.github/CONTRIBUTING.md) for more information.

## Code of Conduct

See [CODE_OF_CONDUCT.md](.github/CODE_OF_CONDUCT.md) for more information.

## Security

See [SECURITY.md](.github/SECURITY.md) for more information.