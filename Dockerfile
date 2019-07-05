FROM centos

USER root

RUN printf '[google64]\nname=Google - x86_64\nbaseurl=http://dl.google.com/linux/rpm/stable/x86_64\nenabled=1\ngpgcheck=1\ngpgkey=https://dl-ssl.google.com/linux/linux_signing_key.pub' > /etc/yum.repos.d/google.repo

RUN yum -y --setopt=tsflags=nodocs install wget &&\
    wget https://dl.google.com/linux/direct/google-chrome-stable_current_x86_64.rpm && \
    yum -y install ./google-chrome-stable_current_x86_64.rpm

COPY xunya-legion .

EXPOSE 9099

CMD ["./xunya-legion"]
