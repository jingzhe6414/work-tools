#!/usr/bin/env bash
#原文安装脚本为华为ECS https://support.huaweicloud.com/usermanual-ecs/ecs_03_0199.html


set -Eeo pipefail
TAG="555.58.02"
URI_PRE="https://hgcs-drivers"

#######################################
# log_info
# Globals: log_file
# Arguments: None
#######################################
function log_info() {
  echo -en "[$(date +"%Y-%m-%d %H:%M:%S")][INFO] $*" >> "${log_file}"
}
#######################################
# log_warning
# Globals: log_file
# Arguments: None
#######################################
function log_warning() {
  echo -en "[$(date +"%Y-%m-%d %H:%M:%S")][WARNING] $*" >> "${log_file}"
}
#######################################
# echo_error
# Globals: log_file
# Arguments: None
#######################################
function echo_error() {
  echo -en $*
  echo -en "[$(date +"%Y-%m-%d %H:%M:%S")][ERROR] $*" >> "${log_file}"
}
#######################################
# echo_info
# Globals: log_file
# Arguments: None
#######################################
function echo_info() {
  echo -en $*
  echo -en "[$(date +"%Y-%m-%d %H:%M:%S")][INFO] $*" >> "${log_file}"
}


#######################################
# download_software
# Globals: selected_software,uri_base,work_dir
# Arguments: None
# Outputs: None
#######################################
# function download_software() {
# }
function download_software() {
  echo_info "\n****************Download software***************\n"
    echo_info "------------------------------------------------\n"
    log_info "Downloading the driver: https://us.download.nvidia.com/XFree86/Linux-x86_64/$TAG/NVIDIA-Linux-x86_64-$TAG.run\n"
    wget -T 10 https://us.download.nvidia.com/XFree86/Linux-x86_64/$TAG/NVIDIA-Linux-x86_64-$TAG.run -P "${work_dir}"
    log_info "The driver is downloaded successfully and saved in ${work_dir}/NVIDIA-Linux-x86_64-$TAG.run\n"
}


#######################################
# install_nvidia_driver
# Globals: work_dir
# Arguments: None
# Outputs: None
#######################################
function install_nvidia_driver() {
  echo_info "------------------------------------------------\n"
  local driver
  apt install make gcc g++ -y
  driver=${work_dir}/NVIDIA-Linux-x86_64-$TAG.run
  echo_info "Installing the Nvidia GPU driver: ${driver}\n"
  if [[ -z "${driver}" ]]; then
    echo_error "The driver file is missing.\n"
    exit 1
  fi
  chmod +x "${driver}"
  if lsmod | grep -qi "nouveau"; then
    ${driver} -s --no-kernel-module --no-nouveau-check
    ${driver} -x --target "${work_dir}"/nvidia_driver
    make -k -j -C "${work_dir}"/nvidia_driver/kernel
    mv "${work_dir}"/nvidia_driver/kernel/*.ko /lib/modules/"$(uname -r)"/kernel/drivers/video/
    depmod -a
  else
    ${driver} -s
  fi
  echo_info "The Nvidia GPU driver installation is complete.\n\n"
}

#######################################
# disable_mig
# Globals: None
# Arguments: None
# Outputs: None
#######################################
function disable_mig() {
  if [[ -f "/lib/systemd/system/nvidia-mig-disable.service" ]]; then
    return
  fi
cat > "/lib/systemd/system/nvidia-mig-disable.service" <<EOF
[Unit]
Description=NVIDIA disable MIG service
After=nvidia-fabricmanager.service
[Service]
User=root
Type=oneshot
RemainAfterExit=yes
ExecStart=/usr/bin/nvidia-smi -mig 0
TimeoutStartSec=0
[Install]
WantedBy=multi-user.target
EOF
  systemctl daemon-reload
  systemctl enable nvidia-mig-disable.service
}
#######################################
# install_software
# Globals: selected_software
# Arguments: None
# Outputs: None
#######################################
function install_software() {
  echo_info "\n****************Install software****************\n"
  if [[ -n "${selected_software["ascend"]}" ]]; then
    install_ascend_driver
  else
    install_nvidia_driver
    disable_mig
  fi
}
#######################################
# disable_nouveau
# Globals: os_type
# Arguments: None
# Outputs: None
#######################################
function disable_nouveau() {
  echo_info "------------------------------------------------\n"
  echo_info "Write 'blacklist nouveau' and 'options nouveau modeset=0' to '/etc/modprobe.d/blacklist.conf'.\n"
  echo "blacklist nouveau" >> /etc/modprobe.d/blacklist.conf
  echo "options nouveau modeset=0" >> /etc/modprobe.d/blacklist.conf
  echo_info "Update initramfs...\n"
  if echo "${os_type}" | grep -qi 'centos\|HuaweiCloudEulerOS'; then
    mv /boot/initramfs-"$(uname -r)".img /boot/initramfs-"$(uname -r)"-nouveau.img
    dracut /boot/initramfs-"$(uname -r)".img "$(uname -r)"
  elif echo "${os_type}" | grep -qi 'ubuntu'; then
    update-initramfs -u
  fi
  echo_info "The nouveau driver disable complete.\n"
}
#######################################
# check_auto_install
# Globals: os_type
# Arguments: None
# Outputs: None
#######################################
function check_auto_install() {
  local i=0
  local log_str
  local pstr
  local time_str
  local begin
  local end
  local total
  

  time_str=$(head -n 2 "${log_file}" | tail -n 1)
  time_str=${time_str: 1:19}
  begin=$(date -d "${time_str}" +%s)
  end=$(date '+%s')
  total=$((end-begin))
  i=$((total/3))

  while [[ ${i} -lt 100 ]]
  do
    log_str=$(tail -n 200 "${log_file}")
    if ! ps ax | grep -v grep | grep "auto_install.sh" | grep -q "NVIDIA-Linux-"; then
      printf "%s\n\n" "${log_str}"
      break
    fi
    if echo "${log_str}" | grep -q "Auto install end"; then
      i=100
    else
      i=$((i+1))
      if [ $i -ge 99 ]; then
        i=99
      fi
    fi
    pstr=$(printf "%-${i}s" "#")
    pstr=${pstr// /#}
    clear
    printf "%s\n\n" "${log_str}"
    printf "| %-100s | %d%% " "${pstr}" "${i}"
    sleep 3
  done 
}
#######################################
# usage
# Globals: None
# Arguments: None
# Outputs: None
#######################################
function usage() {
  printf "usage: auto_install.sh [OPTIONS]...\n"
  printf "        -c fileName, Specifies the fileName of the CUDA to be installed..\n"
  printf "        -u fileName, Specifies the fileName of the cudnn to be installed.\n"
  printf "        -t fileName, Specifies the fileName of the tesla driver to be installed.\n"
  printf "        -g fileName, Specifies the fileName of the grid driver to be installed.\n"
  printf "        -f fileName, Specifies the fileName of the fabricmanager to be installed.\n"
  printf "        -q fileName, Specifies the fileName of the quadro driver to be installed.\n"
  printf "        -a fileName, Specifies the fileName of the Ascend NPU driver to be installed.\n"
  printf "        -n fileName, Specifies the fileName of the Ascend cann to be installed.\n"
}
#######################################
# error_handler
# Globals: selected_software
# Arguments: None
# Outputs: None
#######################################
function error_handler() {
  local ret=$?
  [ ${ret} -eq 0 ] && return
  echo_error "execute '${BASH_COMMAND}' failed, at function:${FUNCNAME[1]}, line:${BASH_LINENO[0]}, exit code=${ret}\n"
  echo_error "For details, see:\n"
  echo_error "For auto install: /var/log/auto-install.log\n"
  if [[ -n "${selected_software["ascend"]}" ]]; then
    echo_error "For Ascend NPU driver: /var/log/ascend_seclog/ascend_install.log\n"
    echo_error "For Ascend cann: /var/log/ascend_seclog/ascend_toolkit_install.log\n"
  else
    echo_error "For Nvidia driver: /var/log/nvidia-installer.log or /var/log/nvidia-uninstall.log\n"
    echo_error "For Nvidia cuda: /var/log/cuda-installer.log or /var/log/cuda-uninstaller.log\n"
  fi
  exit ${ret}
}


#######################################
# main
# Globals: software,candidate_software,selected_software
# Arguments: opt
#######################################
function main() {
  declare -A software   # Store software in key-value pairs.
  declare -A candidate_software  # Store candidate software in key-value pairs.
  # Store the software that will be installed.
  # selected_software["cuda"]=filename,selected_software["cudnn"]=filename,
  # selected_software["tesla"]=filename,selected_software["grid"]=filename,selected_software["quadro"]=filename
  declare -A selected_software

  local os_type
  local instance_type
  local machine_type
  local uri_base
  local work_dir
  local log_file
  local optios_num
  local release="release"
  local check="false"
  local profile_file=""
  local begin
  local end
  local total_time

  while getopts 'c:u:t:g:q:f:p:a:n:z:' opt; do
    case ${opt} in
      c)
        selected_software["cuda"]="$OPTARG" ;;
      u)
        selected_software["cudnn"]="$OPTARG" ;;
      t)
        selected_software["tesla"]="$OPTARG" ;;
      g)
        selected_software["grid"]="$OPTARG" ;;
      q)
        selected_software["quadro"]="$OPTARG" ;;
      f)
        selected_software["fabricmanager"]="$OPTARG" ;;
      a)
        selected_software["ascend"]="$OPTARG" ;;
      n)
        selected_software["cann"]="$OPTARG" ;;
      p)
        release="$OPTARG" ;;
      z)
        check="$OPTARG" ;;
      *)
        usage
        exit 1
        ;;
    esac
  done

  log_file="/var/log/auto-install.log"

  os_type="$(lsb_release -i | awk '{print $3}') $(lsb_release -r | awk '{print $2}' | awk -F. '{print $1"."$2}')"
  if echo "${os_type}" | grep -qi 'centos\|HuaweiCloudEulerOS'; then
    profile_file="/root/.bash_profile"
  elif echo "${os_type}" | grep -qi 'ubuntu'; then
    profile_file="/root/.profile"
  fi

  if [[ "${check}" == "true" ]]; then
    check_auto_install
    exit 0
  fi

  work_dir="$(pwd)/auto_install_$(date +'%Y%m%d%H%M%S')"
  mkdir -p "${work_dir}"
  echo > ${log_file}

  echo_info "***************Auto install begin***************\n"
  log_info "Options: $*\n"

  if ! python3 --version >/dev/null; then
    echo_error "Python3 is not installed. Please install it first.\n"
    exit 1
  fi


  begin=$(date '+%s')


  download_software

  install_software

  if lsmod | grep -q "nouveau"; then
    disable_nouveau
  fi

  rm -rf "${work_dir}"
  rm -rf "${BASH_SOURCE[0]}"

  sed -i '/auto_install.sh -z true/d' "${profile_file}"

  end=$(date '+%s')
  total_time=$((end-begin))

  echo_info "------------------------------------------------\n"
  echo_info "The installation total time: $total_time s.\n"
  echo_info "The installation is complete and will be reboot after 15 seconds.\n"
  echo_info "****************Auto install end****************\n"
  sleep 15
  reboot
}

trap 'error_handler' ERR

main "$@"
