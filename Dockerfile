FROM centos:centos7

USER root

RUN printf '[google64]\nname=Google - x86_64\nbaseurl=http://dl.google.com/linux/rpm/stable/x86_64\nenabled=1\ngpgcheck=1\ngpgkey=https://dl-ssl.google.com/linux/linux_signing_key.pub' > /etc/yum.repos.d/google.repo

RUN rpm --import http://dl.fedoraproject.org/pub/epel/RPM-GPG-KEY-EPEL-7 &&\
    yum -y update && yum -y install epel-release &&\
    yum -y install \
    Xvfb \
    google-chrome-stable \
    https://centos7.iuscommunity.org/ius-release.rpm


COPY xunya-legion .

EXPOSE 9099

CMD ["./xunya-legion"]
