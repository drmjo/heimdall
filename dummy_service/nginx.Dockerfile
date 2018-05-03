FROM nginx:1.12

RUN apt-get -y update \
    && apt-get -y install \
      vim \
      curl \
    && apt-get clean

COPY default.conf /etc/nginx/conf.d/default.conf
