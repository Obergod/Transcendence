#!/bin/sh
while ! timeout 1 bash -c "echo > /dev/tcp/wordpress/9000"; do
  echo "Waiting for PHP-FPM..."
  sleep 1
done

nginx -g 'daemon off;'