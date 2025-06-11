FROM docker.elastic.co/beats/filebeat:8.13.0

# Copia o filebeat.yml com permissões corretas diretamente
COPY --chown=root:root --chmod=600 filebeat.yml /usr/share/filebeat/filebeat.yml
