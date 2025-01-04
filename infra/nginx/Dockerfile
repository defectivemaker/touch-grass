FROM nginx:alpine

# Remove default Nginx config
RUN rm /etc/nginx/conf.d/default.conf

# Install envsubst
RUN apk add --no-cache gettext

# Copy Nginx config template
COPY nginx.conf.template /etc/nginx/nginx.conf.template

# Copy SSL certificates
# COPY ./ssl /etc/nginx/ssl

# Copy entrypoint script
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Expose ports
EXPOSE 80 443

# Set entrypoint
ENTRYPOINT ["/entrypoint.sh"]