#!/bin/sh

# Replace environment variables in the nginx config template
envsubst '${NGINX_DOMAIN}' < /etc/nginx/nginx.conf.template > /etc/nginx/conf.d/default.conf

# Start Nginx
nginx -g 'daemon off;'