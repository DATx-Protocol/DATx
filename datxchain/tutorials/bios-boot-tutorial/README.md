# datx bios tutorial

The `datx-bios-tutorial.py` script simulates the DATXOS bios boot sequence.  
、
# 在datx/Code/tutorials/bios-boot-tutorial目录下：

## 操作步骤  
## 第一步：配置修改
    1)打开datx-normal-node-tutorial.py文件
	  修改datx-normal-node-tutorial.py的peer to peer address
      添加其他节点的peer to peer address，必须添加种子节点的地址，目前种子节点地址为：172.31.3.5。
	  切记不要写入自己的地址，且地址必须是已经在datx链上其他节点存在的地址  

    2)在my_account.json文件中填入：节点账户名称，公私钥。
	  目前节点账户名称、公私钥可在accounts.json文件中producer字段下选择一个复制到my_account.json文件中
	  也可以用自己创建的节点名称，但需要将账户名、公私钥发给种子节点的操作者
      将账户名称、公私钥等信息提供给种子节点，由种子节点为其分配资源。  

## 第二步：kill掉自己的noddatx等进程  
     可用  
     killall kdatxd noddatx || true 
     命令，也可用合适的其他命令，比如：
     ./datx-bios-tutorial.py -k

## 第三步：种子节点需要 
     执行：
     ./datx-bios-tutorial.py -a 
     种子节点会部署合约，创建系统账户,部署发币、提币、dbtc,deth,deos

	或执行：
     ./datx-bios-tutorial或或-v2.py -a 
     种子节点会部署合约，创建系统账户 

## 第四步：普通节点需要
     1）执行前，确定种子节点已经启动
	 2）确定kill掉noddatx等进程
	 3）等种子节点启动后执行：  
        ./datx-normal-node-tutorial.py -a 
     
## 第五步：种子节点需要
     ./datx-vote-producer.py -a --http-server 种子节点IP
     这一步，会给普通节点投票，是其有出块的功能。

## 注意，datx-bios-tutorial.py脚本会

		1)  删除默认目录默认钱包 ,重新创建新钱包  
		    注：datx-normal-node-tutorial.py在默认目录（/root/datxos-wallet/）下删除掉上次创建钱包 my_wallet（如果存在）,重新创建新钱包
		2)  导入datxos的私钥，导入本地的私钥  
		3)  启动种子节点  
		4)  部署以下合约  
			token  
			msig  
			charg  
			system  
		5)  发行datx token 10000000000.0000  
			21000000.0000 DBTC  
			102000000.0000 DETH  
			1000000000.0000 DEOS  
		6)  创建测试用质押账户  
		7)  创建系统账户  
		8)  将某些质押账户注册为超级节点  
		9)  普通质押账户为超级节点投票
		