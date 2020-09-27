FROM nginx:1.15-alpine
CMD envsubst '$$APP_CONTAINER$$DOMAIN_NAME$$API_CONTAINER' < /etc/nginx/conf.d/default.conf.template > /etc/nginx/conf.d/default.conf && \
    nginx -g 'daemon off;' && \
    while :; do sleep 6h && wait ${!}; nginx -s reload; done
ADD ./data/nginx/nginx.conf /etc/nginx
