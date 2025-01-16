#!/bin/bash
set -euo pipefail

help() {
  echo "
deploy the basic application and all integration services
this script is executed on remote machine

Environment Variables:
  USER       : Username on the remote machine (e.g., ec2-user)
  APP        : Name of the application being deployed (e.g., my-app)
  WORK_DIR   : Working directory on remote machine
  TAR_FILE   : Path to the application's tar.gz file
  DEPLOY_RUN : Determines the deployment 'app' to base application or 'intg' to update integration services
"
}

if [[ "${1:-}" == "--help" || "${1:-}" == "-h" ]]; then
  help
  exit 0
fi

# library #

stop_service() {
  local service_name="$1"

  echo "Stopping service: ${service_name}"
  if systemctl is-active --quiet "${service_name}"; then
    sudo systemctl stop "${service_name}"
    echo "Service ${service_name} stopped successfully"
  fi
}

restart_service() {
  local service_name="$1"

  echo "Restarting service: ${service_name}"
  sudo systemctl restart "${service_name}"
  sleep 2
  if ! systemctl is-active --quiet "${service_name}"; then
    echo "Failed to restart ${service_name}"
    exit 1
  fi
  echo "${service_name} restarted successfully"
}

update_systemd() {
  local app="$1"
  local sys="$2"

  echo "Updating systemd service from ${app} to ${sys}"
  sudo mv "${app}" "${sys}"
  sudo ln -s "${sys}" "${app}"
  sudo chown -R root:root "${sys}"
  sudo systemctl daemon-reload
}

update_logrotate() {
  local app="$1"
  local sys="$2"

  echo "Updating logrotate configuration from ${app} to ${sys}"
  sudo mv "${app}" "${sys}"
  sudo ln -s "${sys}" "${app}"
  sudo chown -R root:root "${sys}"
  sudo chmod 644 "${sys}"

  if sudo logrotate -d "${sys}"; then
    echo "logrotate configuration validated successfully"
  else
    echo "Error: logrotate configuration file is invalid"
    exit 1
  fi
}

update_fluentd_conf() {
  local app="$1"
  local sys="$2"

  echo "Updating fluent configuration from ${app} to ${sys}"
  sudo mv "${app}" "${sys}"
  sudo ln -s "${sys}" "${app}"
  sudo chown -R root:root "${sys}"
  sudo chmod 644 "${sys}"
}

verify_and_restart_fluentd() {
  local conf="$1"

  if fluentd --dry-run -c "${conf}"; then
    echo "fluent configuration validated successfully"
  else
    echo "Error: fluent configuration validation failed"
    exit 1
  fi
  restart_service fluentd
}

update_crontab() {
  local keyword="$1"
  local cron_job="$2"

  echo "Updating crontab for keyword: ${keyword}"
  line=$(grep -n -w "$keyword" /etc/crontab | cut -d: -f1)
  if [ -n "$line" ]; then
      sudo sed -i "${line}d" /etc/crontab
      echo "Removed existing cron job for keyword: ${keyword}"
  fi
  echo "$cron_job" | sudo tee -a /etc/crontab > /dev/null
  echo "Cron job added: ${cron_job}"
}

# main #

update_integration() {
  echo "Updating integration services..."

  # /var/log/
  echo "Creating log directory for ${APP}"
  sudo mkdir -p "/var/log/${APP}"
  sudo chmod 775 "/var/log/${APP}"
  sudo touch "/var/log/${APP}/stdout.log" "/var/log/${APP}/stderr.log"
  sudo chown "${USER}:${USER}" /var/log/${APP}
  sudo chown "${USER}:${USER}" /var/log/${APP}/*.log

  # logrotate
  update_logrotate "${WORK_DIR}/builds/logrotate.conf" "/etc/logrotate.d/${APP}.conf"

  # fluent
  update_fluentd_conf "${WORK_DIR}/builds/fluent_bq.conf" "/etc/fluent/fluent_bq.conf"
  update_fluentd_conf "${WORK_DIR}/builds/fluent_bq_buffer.conf" "/etc/fluent/fluent_bq_buffer.conf"
  update_fluentd_conf "${WORK_DIR}/builds/fluent.conf" "/etc/fluent/fluentd.conf"
  verify_and_restart_fluentd "/etc/fluent/fluentd.conf"

  # crontab
  update_crontab "cron.daily" "  0  0  * * *   root    test -x /usr/sbin/anacron || ( cd / && run-parts --report /etc/cron.daily )"
  restart_service cron

  # systemd
  update_systemd "${WORK_DIR}/builds/systemd.service" "/etc/systemd/system/${APP}.service"
}

if [ ! -f "$TAR_FILE" ]; then
  echo "Error: The file ${TAR_FILE} does not exist."
  exit 1
fi

stop_service "${APP}"

sudo mkdir -p "${WORK_DIR}"
echo "tar -xzvf ${TAR_FILE} -C ${WORK_DIR}"
sudo tar -xzvf "${TAR_FILE}" -C "${WORK_DIR}"
sudo chown -R "${USER}:${USER}" "${WORK_DIR}"

if [[ "${DEPLOY_RUN:-}" == "intg" ]]; then
  update_integration
fi

restart_service "${APP}"

sudo systemctl enable "${APP}"
