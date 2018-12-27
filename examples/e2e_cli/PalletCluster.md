############################################################  
#                    Pallet Cluster Building Manual        #
#                    Pallet��Ⱥ��ֲ�                    #
############################################################  
This document is the basic configuration document for the private network built by PalletOne. The specific steps are as follows:
���ĵ�ΪPalletOne�˽������Ļ��������ĵ������岽��ο����£�
############################################################
#                    (1)First                              #
#                    (1)��һ����                           #
############################################################
(1)Basic Environment Configuration
(1)������������
	(1.1)Cluster planning
	(1.1)��Ⱥ�滮
		IP	            Name	    Role	    OS
		192.168.110.117	K8SNode01	Mediator0	Centos7
		192.168.110.118	K8SNode02	Mediator1	Centos7
		192.168.110.119	K8SNode03	Mediator2	Centos7
		192.168.110.120	K8SNode04	Mediator3	Centos7
	(1.2)Host configuration
	(1.2)��������
		1>Close SELinux and firewall
		1>�ر�selinux�ͷ���ǽ
		#sudo sed -i 's/SELINUX=enforcing/SELINUX=disabled/g' /etc/selinux/config
		#sudo systemctl stop firewalld.service && systemctl disable firewalld.service
		#
		2>Configure Host Time, Time Zone, System Language
		2>��������ʱ�䡢ʱ����ϵͳ����
		##Modify time zone
		##�޸�ʱ��
		#ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
		#
		##Modifying the System Language Environment
		##�޸�ϵͳ���Ի���
		#sudo echo 'LANG="en_US.UTF-8"' >> /etc/profile;source /etc/profile
		#
		##Configure host NTP time synchronization
		##��������NTPʱ��ͬ��
		#ntpdate -u cn.pool.ntp.org
		#
		3>performance tuning of Kernel
		3>Kernel���ܵ���
		#cat >> /etc/sysctl.conf<<EOF
		net.ipv4.ip_forward=1
		net.ipv4.neigh.default.gc_thresh1=4096
		net.ipv4.neigh.default.gc_thresh2=6144
		net.ipv4.neigh.default.gc_thresh3=8192
		EOF
		#sysctl -p
		#
		4>Add host information
		4>���������Ϣ��/etc/hosts
		#sudo hostnamectl set-hostname k8snode01     #Each node is modified with a different name[�����ڵ�����޸ģ����ֲ�ͬ]
		#vi /etc/hosts
		192.168.110.117 k8snode01
		192.168.110.118 k8snode02
		192.168.110.119 k8snode03
		192.168.110.120 k8snode04
	(1.3)Configure Docker
	(1.3)����Docker
		0>Installation Pre-Dependency
		0>��װǰ������
		#yum -y install expect spawn jq
		#
		1>Adding users (optional)
		1>����û�(��ѡ)
		#sudo adduser k8sdocker
		#
		2>Setting passwords for new users(optional)
		2>Ϊ���û���������(��ѡ)
		#sudo passwd k8sdocker
		#
		3>Setting passwords for new users(optional)
		3>Ϊ���û����sudoȨ��(��ѡ)
		#sudo echo 'k8sdocker ALL=(ALL) ALL' >> /etc/sudoers
		#
		4>Uninstall old version of Docker software
		4>ж�ؾɰ汾Docker���
		#sudo yum remove docker \
		              docker-client \
		              docker-client-latest \
		              docker-common \
		              docker-latest \
		              docker-latest-logrotate \
		              docker-logrotate \
		              docker-selinux \
		              docker-engine-selinux \
		              docker-engine \
		              container*
		#
		5>Define installation version
		5>���尲װ�汾
			5.1># step 1: install the necessary system tools
			5.1># step 1: ��װ��Ҫ��һЩϵͳ����
			#sudo yum update -y
			#sudo yum install -y yum-utils device-mapper-persistent-data lvm2 bash-completion
			#export docker_version=17.03.2
			#
			5.2># Step 2: Adding Software Source Information
			5.2># Step 2: ������Դ��Ϣ
			#sudo yum-config-manager --add-repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo
			#
			5.3># Step 3: Update and install Docker-CE
			5.3># Step 3: ���²���װ Docker-CE
			#sudo yum makecache all
			#version=$(yum list docker-ce.x86_64 --showduplicates | sort -r|grep ${docker_version}|awk '{print $2}')
			#sudo yum -y install --setopt=obsoletes=0 docker-ce-${version} docker-ce-selinux-${version}
			#
			5.4># Step 4: If you have installed a high-level version of Docker, you can downgrade the installation (optional)
			5.4># Step 4: ����Ѿ���װ�߰汾Docker,�ɽ��н�����װ(��ѡ)
			#yum downgrade --setopt=obsoletes=0 -y docker-ce-${version} docker-ce-selinux-${version}
			#
			5.5># Step 5: Set up boot start
			5.5># Step 5: ���ÿ�������
			#sudo systemctl enable docker && systemctl start docker.service
			#
		6>Docker configuration
		6>Docker����
		##Configure image acceleration address
		##���þ�����ٵ�ַ
		#vi /etc/docker/daemon.json
		{
		  "registry-mirrors": ["https://kri93zmv.mirror.aliyuncs.com"]
		}
		��ע��������֣�Failed to start Docker Application Container Engine����ִ��������䣺
		#rm /etc/docker/key.json 
		#rm -rf /var/lib/docker/
	(1.4)Configure Palletone
	(1.4)����Palletone
		0>Configure GOPATH environment variables
		0>����GOPATH��������
		##Install the specified go1.10.4. version
		##��װָ��go1.10.4.�汾
		#go version
		#
		##Configuring environment variables
		##���û�������
		#vi /etc/profile
		export GOROOT=/opt/go
		export GOPATH=/opt/gopath
		export PATH=$PATH:/opt/go/bin:/opt/gopath/bin
		#source /etc/profile
		#
		1>Download the code and compile it
		1>���ش��룬��������б���
		#cd $GOPATH
		#mkdir -p src/github.com/palletone
		#cd $GOPATH/src/github.com/palletone
		#go get -u github.com/palletone/adaptor
		#go get -u github.com/palletone/btc-adaptor
		#go get -u github.com/palletone/eth-adaptor
		#git clone https://github.com/palletone/go-palletone.git
		#ll
			drwxr-xr-x  3 root root  138 Dec 24 02:14 adaptor
			drwxr-xr-x  4 root root  230 Dec 24 02:15 btc-adaptor
			drwxr-xr-x  5 root root 4096 Dec 24 02:15 eth-adaptor
			drwxr-xr-x 25 root root 4096 Dec 24 20:39 go-palletone
		#
		#cd go-palletone/
		#make all
		#
		##Note: Be sure to remember the path of the gptn executable
		##��ע����ؼ�סgptn��ִ�г����·��($PROJECT/build/bin/gptn)
		#
		2>Give privileges to scripts
		2>Ϊ�ű�����Ȩ��
		#cd $GOPATH/src/github.com/palletone/go-palletone/examples/e2e_cli
		#chmod -R 775 *
		#
		##Modify the location of the gptn executable file in getAccount.sh in the scripts directory to the absolute path under build/bin
		##��scriptsĿ¼�е�getAccount.sh�е�gptn��ִ���ļ�λ�ý����޸�Ϊbuild/bin����ľ���·��
									#!/usr/bin/expect
									#!/bin/bash
									#�����˺���Ϣ
									set timeout 30
									#��Ҫ����gptn�ľ���·����ַ
									#spawn "gptn's dir" account new
									spawn /opt/gopath/src/github.com/palletone/go-palletone/build/bin/gptn account new
									expect "Passphrase:"
									send "1\r"
									expect "Repeat passphrase:"
									send "1\r"
									interact
		#
		##Modify the location of the gptn executable file in getInit.sh in the scripts directory to the absolute path under mediator0
		##��scriptsĿ¼�е�getInit.sh�е�gptn��ִ���ļ�λ�ý����޸�Ϊmediator0����ľ���·��
									#!/usr/bin/expect
									set timeout 30
									#����gptn��Mediator0�ľ���·��
									spawn /opt/gopath/src/github.com/palletone/go-palletone/examples/e2e_cli/channel-artifacts/mediator0/gptn init
									expect "Passphrase:"
									send "1\r"
									interact
		#
		3>Generating configuration path information for N nodes
		3>����N���ڵ������·����Ϣ
		#./generateArtifacts.sh 3               #The latter parameter can only be odd[����Ĳ���ֻ��Ϊ����]
		#ls -l
		##Note: After executing the script, the channel-artifacts directory is generated. If the execution fails, the directory is deleted.
		##��ע��ִ�иýű�������channel-artifactsĿ¼����ִ��ʧ�ܣ��򽫸�Ŀ¼����ɾ��
############################################################
#                    (2)Second                             #
#                    (2)�ڶ�����                           #
############################################################
(2)Configuring Super Nodes
(2)���ó����ڵ�
	(2.1)Generate Docker image [optional]
	(2.1)����Docker���񡾿�ѡ��
		1>Generating docker image related files
		1>����docker��������ļ�
		#cd $GOPATH/src/github.com/go-palletone
		#make docker
		#cd build/images/gptn
		#
		2>Generate an image and upload it to dockerhub
		2>���ɾ��񣬲��ϴ���dockerhub
		#docker build -t palletone/pallet-gptn:0.6 .
		#docker images
		#
		3>Log in to dockerhub
		3>��¼dockerhub
		#docker login
		#docker push palletone/pallet-gptn:0.6
		#
		##Query the image in dockerhub
		##��dockerhub��ѯ�þ���
		#
	(2.2)����Mediator0
	(2.2)Configure Mediator0
	  ##Enter the Mediator0 directory and configure accordingly. The specific configuration information is as follows:
		##����Mediator0Ŀ¼������Ӧ�����ã�����������Ϣ������ʾ��
		1>Enter mediator0 directory
		1>����mediator0Ŀ¼
		#cd $GOPATH/src/github.com/palletone/go-palletone/examples/e2e_cli
		#cd channel-artifacts/mediator0
		#ls -l
					-rwxr-xr-x 1 root root      291 Dec 25 01:25 getInit.sh
					-rwxr-xr-x 1 root root 49189104 Dec 25 01:25 gptn
					drwxr-xr-x 2 root root       38 Dec 25 01:25 log
					drwx------ 5 root root       68 Dec 25 02:41 palletone
					-rw-r--r-- 1 root root     3343 Dec 25 01:26 ptn-config.toml
					-rw-r--r-- 1 root root     2854 Dec 25 01:26 ptn-genesis.json
		2>Modify toml file
		2>�޸�toml�ļ�
		#vi ptn-config.toml
			[Node]
			DataDir = "/var/palletone/production"
			[Log]
			OutputPaths = ["stdout", "/var/palletone/log/all.log"]
			ErrorOutputPaths = ["stderr", "/var/palletone/log/error.log"]
			[Contract]
			ContractFileSystemPath = "/var/palletone/production/chaincodes"
			[P2P]
			MaxPeers = 25
			NoDiscovery = false
			BootstrapNodes = []
			StaticNodes=["pnode://3ea34ff09489627399bbeac8d3af93b34981afc623228210bd49c8ce11860f78c736aa3721ebb91aec76353a3b93ee6a2aadd05337ab0723a71a7c9f68947144@mediator0:30303","pnode://01f20de81a80738b30d944a756ade9f4222f95a696d45b451aed596eefa204f3c8ae98305363feceeb28f5c140a6736118f59c81716c0cdd123365cad8a528eb@mediator1:30303","pnode://2a891ee523a40961c0760871be0613551aab45ad7a4ecd23369b713601228173b6e91d4ce748a2cbb571ae0c9b4d47ce605500ad3785cbadeb9ca8ba1a412f6e@mediator2:30303"]
		##Note: Change the relative path in the configuration file to absolute path, replace the IP address of gptn in StaticNodes in p2p, and change it to [container name] if it runs locally.
		##��ע:�������ļ��е����·���޸�Ϊ����·������p2p�е�StaticNodes���滻Ϊgptn���ڵ�IP��ַ����Ϊ�������У������޸�Ϊ[��������];
		3>Modify JSON file
		3>�޸�json�ļ�
		#vi ptn-genesis.json
		#
		##Note: No modification is required for the time being.
		##��ע����ʱ�����޸�
	(2.3)����Mediator1
	(2.3)Configure Mediator1
		##Enter the Mediator1 directory and configure it accordingly. The specific configuration information is as follows:
		##����Mediator1Ŀ¼������Ӧ�����ã�����������Ϣ������ʾ��
		1>Enter mediator1 directory
		1>����mediator1Ŀ¼
		#cd $GOPATH/src/github.com/palletone/go-palletone/examples/e2e_cli
		#cd channel-artifacts/mediator1
		#ls -l
					drwxr-xr-x 2 root root       38 Dec 25 01:25 log
					drwx------ 5 root root       68 Dec 25 02:41 palletone
					-rw-r--r-- 1 root root     3343 Dec 25 01:26 ptn-config.toml
		2>Modify toml file
		2>�޸�toml�ļ�
		#vi ptn-config.toml
			[Node]
			DataDir = "/var/palletone/production"
			[Log]
			OutputPaths = ["stdout", "/var/palletone/log/all.log"]
			ErrorOutputPaths = ["stderr", "/var/palletone/log/error.log"]
			[Contract]
			ContractFileSystemPath = "/var/palletone/production/chaincodes"
			[P2P]
			MaxPeers = 25
			NoDiscovery = false
			BootstrapNodes = []
			StaticNodes=["pnode://3ea34ff09489627399bbeac8d3af93b34981afc623228210bd49c8ce11860f78c736aa3721ebb91aec76353a3b93ee6a2aadd05337ab0723a71a7c9f68947144@mediator0:30303","pnode://2a891ee523a40961c0760871be0613551aab45ad7a4ecd23369b713601228173b6e91d4ce748a2cbb571ae0c9b4d47ce605500ad3785cbadeb9ca8ba1a412f6e@mediator2:30303"]
		##Note: 
				Update pnode information to Mediator0 and mediator2;
				Modify the relative path in the configuration file to absolute path��
				Replace the IP address of gptn in StaticNodes in p2p, and change it to [container name] if it runs locally.
		##��ע:
				��pnode��Ϣ����Ϊmediator0��mediator2��
				�������ļ��е����·���޸�Ϊ����·����
				��p2p�е�StaticNodes���滻Ϊgptn���ڵ�IP��ַ����Ϊ�������У������޸�Ϊ[��������];
	(2.4)����Mediator2
	(2.4)Configure Mediator2
		##Enter the Mediator2 directory and configure it accordingly. The specific configuration information is as follows:
		##����Mediator2Ŀ¼������Ӧ�����ã�����������Ϣ������ʾ��
		1>Enter mediator2 directory
		1>����mediator2Ŀ¼
		#cd $GOPATH/src/github.com/palletone/go-palletone/examples/e2e_cli
		#cd channel-artifacts/mediator2
		#ls -l
					drwxr-xr-x 2 root root       38 Dec 25 01:25 log
					drwx------ 5 root root       68 Dec 25 02:41 palletone
					-rw-r--r-- 1 root root     3343 Dec 25 01:26 ptn-config.toml
		2>Modify toml file
		2>�޸�toml�ļ�
		#vi ptn-config.toml
			[Node]
			DataDir = "/var/palletone/production"
			[Log]
			OutputPaths = ["stdout", "/var/palletone/log/all.log"]
			ErrorOutputPaths = ["stderr", "/var/palletone/log/error.log"]
			[Contract]
			ContractFileSystemPath = "/var/palletone/production/chaincodes"
			[P2P]
			MaxPeers = 25
			NoDiscovery = false
			BootstrapNodes = []
			StaticNodes=["pnode://3ea34ff09489627399bbeac8d3af93b34981afc623228210bd49c8ce11860f78c736aa3721ebb91aec76353a3b93ee6a2aadd05337ab0723a71a7c9f68947144@mediator0:30303","pnode://2a891ee523a40961c0760871be0613551aab45ad7a4ecd23369b713601228173b6e91d4ce748a2cbb571ae0c9b4d47ce605500ad3785cbadeb9ca8ba1a412f6e@mediator1:30303"]
		##Note: 
				Update pnode information to Mediator0 and mediator1;
				Modify the relative path in the configuration file to absolute path��
				Replace the IP address of gptn in StaticNodes in p2p, and change it to [container name] if it runs locally.
		##��ע:
				��pnode��Ϣ����Ϊmediator0��mediator2��
				�������ļ��е����·���޸�Ϊ����·����
				��p2p�е�StaticNodes���滻Ϊgptn���ڵ�IP��ַ����Ϊ�������У������޸�Ϊ[��������];
	(2.5)Configure docker template
	(2.5)����dockerģ��
		##Enter the e2e_cli directory and configure it accordingly. The specific configuration information is as follows:
		##����e2e_cliĿ¼������Ӧ�����ã�����������Ϣ������ʾ��
		1>Enter e2e_cli directory
		1>����e2e_cliĿ¼
		#cd $GOPATH/src/github.com/palletone/go-palletone/examples/e2e_cli
		#
		2>Configuration template file
		2>����ģ���ļ�
		#vi docker-compose-e2e.yaml
		#
		##Note: Refer to the file under e2e_cli directory.
		##��ע:�ο�e2e_cliĿ¼����ĸ��ļ���
		3>Start container
		3>��������
		##Startup container in the background
		##��̨��������
		#docker-compose -f docker-compose-e2e.yaml up -d
		#
		##Query running container
		##��ѯ���е�����
		#docker ps
		#
		4>Stop running the container (optional)
		4>ֹͣ��������(��ѡ)
		#docker-compose -f docker-compose-e2e.yaml down
		#
############################################################
#                    (3)Third                              #
#                    (3)��������                           #
############################################################
(3)Configuration data node
(3)�������ݽڵ�
	(3.1)Configuration data node
	(3.1)�������ݽڵ�
		1>Copy mediator0 to mediator3
		1>����mediator0Ϊmediator3
		#cd $GOPATH/src/github.com/palletone/go-palletone/examples/e2e_cli/channel-artifacts
		#
		#cp -rf mediator0 mediator3
		#
		2>Modify the toml configuration file
		2>�޸�toml�����ļ�
		#cd mediator3
		#vi ptn-config.toml
			[MediatorPlugin]
			EnableStaleProduction=false
		3>Add mediator3 configuration information to the template file in the e2e_cli directory.
		3>��e2e_cliĿ¼��ģ���ļ������mediator3������Ϣ
		#vi docker-compose-e2e.yaml
		#
		4>Start the service based on the template file
		4>����ģ���ļ���������
		#docker-compose -f docker-compose-e2e.yaml up -d
		#docker ps
		#
		5>Enter the Mediator 3 directory to view the logs
		5>����mediator3Ŀ¼�鿴��־
		#cd mediator3/log
		#tail -f all.log
############################################################
#                    (4)Four                               #
#                    (4)���Ĳ���                           #
############################################################
(4)End
(4)��β