FROM registry.access.redhat.com/ubi8

ARG USER_ID=1000
ARG USER_NAME=galaxy
ARG USER_GROUP=galaxy
ARG COMPOSE_PROFILE
ARG LOCK_REQUIREMENTS
ARG DEV_SOURCE_PATH
ARG DEV_IMAGE_SUFFIX
ARG DEV_VOLUME_SUFFIX

ENV LANG=en_US.UTF-8 \
    PYTHONUNBUFFERED=1 \
    PULP_SETTINGS=/etc/pulp/settings.py \
    DJANGO_SETTINGS_MODULE=pulpcore.app.settings

RUN \
    dnf -y install \
        gettext \
        git \
        gcc \
        glibc-langpack-en \
        python39-devel \
        libpq \
        libpq-devel \
        pinentry \
        make \
        openldap-devel \
        skopeo \
        procps \
        net-tools \
        jq

RUN groupadd --gid "${USER_ID}" "${USER_GROUP}"
RUN useradd --uid ${USER_ID} --gid "${USER_GROUP}" "${USER_NAME}"
RUN mkdir /app && chown -R ${USER_ID}:${USER_GROUP} /app
RUN mkdir /venv && chown -R ${USER_ID}:${USER_GROUP} /venv

USER ${USER_NAME}:${USER_GROUP}

RUN python3 -m venv /venv
RUN git clone https://github.com/ansible/galaxy_ng /app/galaxy_ng
RUN /venv/bin/pip install --upgrade pip
RUN /venv/bin/pip install -r /app/galaxy_ng/requirements/requirements.common.txt
RUN /venv/bin/pip install -e /app/galaxy_ng
RUN /venv/bin/pip install django_extensions django-storages boto3
