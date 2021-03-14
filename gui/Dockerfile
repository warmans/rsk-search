FROM nginx

COPY dist/gui /usr/share/nginx/html/
COPY deploy/default.conf /etc/nginx/conf.d/

EXPOSE 80

CMD nginx -g 'daemon off;'
