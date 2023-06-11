# -*- mode: ruby -*-
# vi: set ft=ruby :

# All Vagrant configuration is done below. The "2" in Vagrant.configure
# configures the configuration version (we support older styles for
# backwards compatibility). Please don't change it unless you know what
# you're doing.
Vagrant.configure("2") do |config|
  config.vm.synced_folder ".", "/vagrant", type: "nfs", nfs_udp: false

  config.vm.define "galaxy-dev" do |dev|
    dev.vm.box = "generic/debian11"
    dev.vm.hostname = "galaxy-dev"
  end

  config.vm.provider :libvirt do |libvirt|
    libvirt.cpus = 4
    libvirt.memory = 16000
  end

  config.vm.provision "shell", inline: <<-SHELL
    set -x
    export DEBIAN_FRONTEND=noninteractive
    apt -y update
    apt -y upgrade
    apt -y install git jq python3-pip docker.io libpq-dev python3-virtualenv

    # there should be a package for this right?
    if [ ! -L /usr/local/bin/python ]; then
     ln -s /usr/bin/python3 /usr/local/bin/python
    fi
    pip3 install -U pip wheel
    which ansible || pip3 install ansible
    ansible-galaxy role install geerlingguy.docker
    cd /vagrant && ansible-playbook -i 'localhost,' docker.yml

    # pulp workarounds?
    mkdir -p /var/lib/gems
    chown -R vagrant:vagrant /var/lib/gems

    # for vagrant user + docker
    usermod -aG docker runner

    # for ghacktion
    useradd --shell=/bin/bash runner
    cp -Rp /home/vagrant /home/runner
    chown -R runner:runner /home/runner
    usermod -aG docker runner
    cp /etc/sudoers.d/vagrant /etc/sudoers.d/runner
    sed -i.bak 's:vagrant:runner:g' /etc/sudoers.d/runner
    rm -f /etc/sudoers.d/runner.bak
  SHELL

end
