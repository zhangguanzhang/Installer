#!/bin/bash
curl -s -H 'Accept: [*/*]' \
  -H 'User-Agent: curl/7.29.0' \
  -H 'X-Anaconda-Architecture: x86_64' \
  -H 'X-Anaconda-System-Release: CentOS Linux' \
  -H 'X-Rhn-Provisioning-Mac-0: enp61s0f0 9c:e8:95:d8:3c:cc' \
  -H 'X-Rhn-Provisioning-Mac-1: enp61s0f1 9c:e8:95:d8:3c:cd' \
  -H 'X-Rhn-Provisioning-Mac-2: enp61s0f2 9c:e8:95:d8:3c:ce' \
  -H 'X-Rhn-Provisioning-Mac-3: enp61s0f3 9c:e8:95:d8:3c:cf' \
  -H 'X-Rhn-Provisioning-Mac-4: ens1f0 3c:f5:cc:91:1f:68' \
  -H 'X-Rhn-Provisioning-Mac-5: ens1f1 3c:f5:cc:91:1f:6a' \
  -H 'X-Rhn-Provisioning-Mac-6: ens2f0 3c:f5:cc:91:1e:48' \
  -H 'X-Rhn-Provisioning-Mac-7: ens2f1 3c:f5:cc:91:1e:4a' \
  -H 'X-System-Serial-Number: 210200A00QH185002000' $@
