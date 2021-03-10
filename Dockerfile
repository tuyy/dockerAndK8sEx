FROM reg.navercorp.com/base/alpine:latest
RUN mkdir -p /home1/irteam/apps/myfirst
WORKDIR /home1/irteam/apps/myfirst
COPY dist .
CMD ./first