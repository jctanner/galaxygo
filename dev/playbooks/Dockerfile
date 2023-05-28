FROM python:3
RUN apt -y update && apt install -y docker.io
RUN pip install --upgrade pip
RUN pip install ansible-core boto3 requests logzero
RUN ansible-galaxy collection install community.general
RUN ansible-galaxy collection install community.docker
RUN ansible-galaxy collection install amazon.aws
