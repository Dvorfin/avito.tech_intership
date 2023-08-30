FROM debian

RUN apt-get update && apt-get install -y nginx

#копирование конфига самого nginx
COPY nginx.conf /etc/nginx/

# копирование конфига для веб приложения
COPY tech-avito.conf /etc/nginx/sites-available/

# добавление ссылки на ip
RUN ln -s /etc/nginx/sites-available/tech-avito.conf /etc/nginx/sites-enabled/ 

# EXPOSE 7878


CMD ["nginx", "-g", "daemon off;"]




