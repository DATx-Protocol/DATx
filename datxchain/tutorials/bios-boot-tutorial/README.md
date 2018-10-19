# datx bios tutorial

The `datx-bios-tutorial.py` script simulates the DATXOS bios boot sequence.  
# 在datx/Code/tutorials/bios-boot-tutorial目录下：

## 0、普通节点需要
    1)修改datx-normal-node-tutorial.py的peer to peer address
      添加其他节点的peer to peer address，不要写入自己的地址，且地址必须是已经在网络中存在的
    2)在my_account.json文件中填入：节点账户名称，公私钥。  
      并将账户名称、公私钥等信息提供给种子节点，由种子节点为其分配资源。  
## 1、种子节点 
     ./datx-bios-tutorial.py -a --http-server 你自己的IP
      
## 2、普通节点
     等种子节点启动后：  
     ./datx-normal-node-tutorial.py -a --http-server 你自己的IP  

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
		